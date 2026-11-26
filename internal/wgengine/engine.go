package wgengine

import (
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/device"
	"golang.zx2c4.com/wireguard/tun"

	"github.com/JAPStudios/Vanilla-Pecan-Napoleon/internal/wgconfig"
)

// Engine owns a configured WireGuard device. The caller supplies a platform
// TUN device and bind implementation; this package does not configure OS routes.
type Engine struct {
	device *device.Device
}

// Start configures and starts a WireGuard device from a previously validated
// configuration. Platform-specific TUN creation and address/route setup remain
// the responsibility of the caller.
func Start(tunDevice tun.Device, bind conn.Bind, cfg *wgconfig.Config) (*Engine, error) {
	if tunDevice == nil {
		return nil, errors.New("TUN device is required")
	}
	if bind == nil {
		return nil, errors.New("WireGuard bind is required")
	}
	if cfg == nil {
		return nil, errors.New("WireGuard configuration is required")
	}

	uapi, err := buildUAPI(cfg)
	if err != nil {
		return nil, err
	}

	dev := device.NewDevice(tunDevice, bind, device.NewLogger(device.LogLevelError, "vanilla-pecan-napoleon"))
	if err := dev.IpcSet(uapi); err != nil {
		dev.Close()
		return nil, fmt.Errorf("configure WireGuard device: %w", err)
	}
	if err := dev.Up(); err != nil {
		dev.Close()
		return nil, fmt.Errorf("start WireGuard device: %w", err)
	}
	return &Engine{device: dev}, nil
}

// Close stops the WireGuard device and releases its resources.
func (e *Engine) Close() {
	if e != nil && e.device != nil {
		e.device.Close()
	}
}

func buildUAPI(cfg *wgconfig.Config) (string, error) {
	privateKey, err := base64KeyToHex(cfg.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("encode interface private key: %w", err)
	}

	var b strings.Builder
	fmt.Fprintf(&b, "private_key=%s\n", privateKey)
	if cfg.ListenPort > 0 {
		fmt.Fprintf(&b, "listen_port=%d\n", cfg.ListenPort)
	}
	for _, peer := range cfg.Peers {
		publicKey, err := base64KeyToHex(peer.PublicKey)
		if err != nil {
			return "", fmt.Errorf("encode peer public key: %w", err)
		}
		fmt.Fprintf(&b, "public_key=%s\n", publicKey)
		if peer.PresharedKey != "" {
			presharedKey, err := base64KeyToHex(peer.PresharedKey)
			if err != nil {
				return "", fmt.Errorf("encode peer preshared key: %w", err)
			}
			fmt.Fprintf(&b, "preshared_key=%s\n", presharedKey)
		}
		if peer.Endpoint != "" {
			fmt.Fprintf(&b, "endpoint=%s\n", peer.Endpoint)
		}
		for _, prefix := range peer.AllowedIPs {
			fmt.Fprintf(&b, "allowed_ip=%s\n", prefix)
		}
		if peer.PersistentKeepalive > 0 {
			fmt.Fprintf(&b, "persistent_keepalive_interval=%d\n", peer.PersistentKeepalive)
		}
	}
	b.WriteByte('\n')
	return b.String(), nil
}

func base64KeyToHex(value string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", errors.New("invalid base64 key")
	}
	defer func() {
		for i := range key {
			key[i] = 0
		}
	}()
	if len(key) != 32 {
		return "", errors.New("key must contain exactly 32 bytes")
	}
	return hex.EncodeToString(key), nil
}
