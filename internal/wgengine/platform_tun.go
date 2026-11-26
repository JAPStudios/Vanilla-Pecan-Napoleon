//go:build darwin || linux || windows

package wgengine

import (
	"errors"
	"fmt"

	"golang.zx2c4.com/wireguard/conn"
	"golang.zx2c4.com/wireguard/tun"

	"github.com/JAPStudios/Vanilla-Pecan-Napoleon/internal/wgconfig"
)

// StartWithPlatformTUN creates a native WireGuard TUN device and starts the
// engine with the platform's default WireGuard bind.
//
// The caller remains responsible for configuring interface addresses,
// operating-system routes, DNS settings and any required permissions.
func StartWithPlatformTUN(name string, mtu int, cfg *wgconfig.Config) (*Engine, error) {
	if name == "" {
		return nil, errors.New("TUN interface name is required")
	}
	if mtu <= 0 {
		return nil, errors.New("TUN MTU must be greater than zero")
	}
	if cfg == nil {
		return nil, errors.New("WireGuard configuration is required")
	}

	tunDevice, err := tun.CreateTUN(name, mtu)
	if err != nil {
		return nil, fmt.Errorf("create platform TUN device: %w", err)
	}

	engine, err := Start(tunDevice, conn.NewDefaultBind(), cfg)
	if err != nil {
		return nil, fmt.Errorf("start platform WireGuard engine: %w", err)
	}
	return engine, nil
}
