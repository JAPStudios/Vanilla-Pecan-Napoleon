package wgengine

import (
	"strings"
	"net/netip"
	"testing"

	"github.com/JAPStudios/Vanilla-Pecan-Napoleon/internal/wgconfig"
)

func TestBuildUAPI(t *testing.T) {
	cfg := &wgconfig.Config{
		PrivateKey: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
		ListenPort: 51820,
		Peers: []wgconfig.Peer{{
			PublicKey:           "AQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQE=",
			AllowedIPs:          mustPrefixes(t, "0.0.0.0/0", "::/0"),
			Endpoint:            "vpn.example.net:51820",
			PersistentKeepalive: 25,
		}},
	}

	got, err := buildUAPI(cfg)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"private_key=0000000000000000000000000000000000000000000000000000000000000000",
		"listen_port=51820",
		"public_key=0101010101010101010101010101010101010101010101010101010101010101",
		"allowed_ip=0.0.0.0/0",
		"allowed_ip=::/0",
		"endpoint=vpn.example.net:51820",
		"persistent_keepalive_interval=25",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("UAPI configuration missing %q", want)
		}
	}
	if !strings.HasSuffix(got, "\n\n") {
		t.Fatal("UAPI configuration must terminate with a blank line")
	}
}

func TestBase64KeyToHexRejectsInvalidLength(t *testing.T) {
	if _, err := base64KeyToHex("AQ=="); err == nil {
		t.Fatal("expected invalid key length error")
	}
}

func mustPrefixes(t *testing.T, values ...string) []netip.Prefix {
	t.Helper()
	prefixes := make([]netip.Prefix, 0, len(values))
	for _, value := range values {
		prefix, err := netip.ParsePrefix(value)
		if err != nil {
			t.Fatal(err)
		}
		prefixes = append(prefixes, prefix)
	}
	return prefixes
}
