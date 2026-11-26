//go:build darwin || linux || windows

package wgengine

import (
	"testing"

	"github.com/JAPStudios/Vanilla-Pecan-Napoleon/internal/wgconfig"
)

func TestStartWithPlatformTUNRejectsInvalidArguments(t *testing.T) {
	validConfig := &wgconfig.Config{}

	tests := []struct {
		name string
		iface string
		mtu int
		cfg *wgconfig.Config
		want string
	}{
		{name: "missing interface name", iface: "", mtu: 1420, cfg: validConfig, want: "TUN interface name is required"},
		{name: "invalid MTU", iface: "vpn0", mtu: 0, cfg: validConfig, want: "TUN MTU must be greater than zero"},
		{name: "missing configuration", iface: "vpn0", mtu: 1420, cfg: nil, want: "WireGuard configuration is required"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := StartWithPlatformTUN(tt.iface, tt.mtu, tt.cfg)
			if err == nil || err.Error() != tt.want {
				t.Fatalf("StartWithPlatformTUN() error = %v, want %q", err, tt.want)
			}
		})
	}
}
