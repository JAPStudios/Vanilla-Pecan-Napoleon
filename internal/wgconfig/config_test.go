package wgconfig

import (
	"strings"
	"testing"
)

const goodConfig = `[Interface]
PrivateKey = AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=
Address = 10.8.0.2/32, fd00::2/128
DNS = 1.1.1.1, 2606:4700:4700::1111

[Peer]
PublicKey = AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=
AllowedIPs = 0.0.0.0/0, ::/0
Endpoint = vpn.example.net:51820
PersistentKeepalive = 25
`

func TestParseValidConfig(t *testing.T) {
	cfg, err := Parse(strings.NewReader(goodConfig))
	if err != nil {
		t.Fatal(err)
	}
	if len(cfg.Addresses) != 2 || len(cfg.DNS) != 2 || len(cfg.Peers) != 1 {
		t.Fatalf("unexpected parsed configuration: %+v", cfg)
	}
}

func TestRejectsInvalidPrivateKeyWithoutEchoingIt(t *testing.T) {
	input := strings.Replace(goodConfig, "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=", "not-a-key", 1)
	_, err := Parse(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected invalid key error")
	}
	if strings.Contains(err.Error(), "not-a-key") {
		t.Fatal("error exposed key material")
	}
}

func TestRejectsDuplicateSetting(t *testing.T) {
	input := strings.Replace(goodConfig, "Address = 10.8.0.2/32, fd00::2/128", "Address = 10.8.0.2/32\nAddress = 10.8.0.3/32", 1)
	if _, err := Parse(strings.NewReader(input)); err == nil {
		t.Fatal("expected duplicate setting error")
	}
}

func TestRejectsInvalidEndpoint(t *testing.T) {
	input := strings.Replace(goodConfig, "vpn.example.net:51820", "not-an-endpoint", 1)
	if _, err := Parse(strings.NewReader(input)); err == nil {
		t.Fatal("expected invalid endpoint error")
	}
}

func TestRejectsDuplicatePeerPublicKey(t *testing.T) {
	input := goodConfig + "\n[Peer]\nPublicKey = AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=\nAllowedIPs = 192.0.2.0/24\n"
	if _, err := Parse(strings.NewReader(input)); err == nil {
		t.Fatal("expected duplicate peer public key error")
	}
}
