package wgconfig

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/netip"
	"strconv"
	"strings"
)

type Peer struct {
	PublicKey           string
	PresharedKey        string
	AllowedIPs          []netip.Prefix
	Endpoint            string
	PersistentKeepalive int
}

type Config struct {
	PrivateKey string
	Addresses  []netip.Prefix
	DNS        []netip.Addr
	ListenPort int
	Peers      []Peer
}

func validKey(value string) bool {
	b, err := base64.StdEncoding.DecodeString(value)
	return err == nil && len(b) == 32
}

func parsePrefixes(value string) ([]netip.Prefix, error) {
	parts := strings.Split(value, ",")
	out := make([]netip.Prefix, 0, len(parts))
	for _, part := range parts {
		p, err := netip.ParsePrefix(strings.TrimSpace(part))
		if err != nil {
			return nil, fmt.Errorf("invalid IP prefix")
		}
		out = append(out, p.Masked())
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("at least one IP prefix is required")
	}
	return out, nil
}

func parseDNS(value string) ([]netip.Addr, error) {
	parts := strings.Split(value, ",")
	out := make([]netip.Addr, 0, len(parts))
	for _, part := range parts {
		a, err := netip.ParseAddr(strings.TrimSpace(part))
		if err != nil {
			return nil, fmt.Errorf("invalid DNS address")
		}
		out = append(out, a)
	}
	return out, nil
}

func parsePort(value string, allowZero bool) (int, error) {
	n, err := strconv.Atoi(strings.TrimSpace(value))
	minimum := 1
	if allowZero {
		minimum = 0
	}
	if err != nil || n < minimum || n > 65535 {
		return 0, fmt.Errorf("port must be between %d and 65535", minimum)
	}
	return n, nil
}

func Parse(r io.Reader) (*Config, error) {
	cfg := &Config{}
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 4096), 1024*1024)
	section := ""
	interfaceSeen := false
	fields := map[string]bool{}
	lineNo := 0

	for scanner.Scan() {
		lineNo++
		line := strings.TrimSpace(scanner.Text())
		if i := strings.IndexAny(line, "#;"); i >= 0 {
			line = strings.TrimSpace(line[:i])
		}
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			name := strings.TrimSpace(line[1 : len(line)-1])
			switch name {
			case "Interface":
				if interfaceSeen {
					return nil, fmt.Errorf("line %d: duplicate Interface section", lineNo)
				}
				interfaceSeen = true
				section = "Interface"
				fields = map[string]bool{}
			case "Peer":
				if !interfaceSeen {
					return nil, fmt.Errorf("line %d: Peer section must follow Interface", lineNo)
				}
				section = "Peer"
				cfg.Peers = append(cfg.Peers, Peer{})
				fields = map[string]bool{}
			default:
				return nil, fmt.Errorf("line %d: unsupported section", lineNo)
			}
			continue
		}
		if section == "" {
			return nil, fmt.Errorf("line %d: setting appears outside a section", lineNo)
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			return nil, fmt.Errorf("line %d: expected key=value", lineNo)
		}
		key, value = strings.TrimSpace(key), strings.TrimSpace(value)
		if key == "" || value == "" {
			return nil, fmt.Errorf("line %d: empty key or value", lineNo)
		}
		if fields[key] {
			return nil, fmt.Errorf("line %d: duplicate %s setting", lineNo, key)
		}
		fields[key] = true

		if section == "Interface" {
			switch key {
			case "PrivateKey":
				if !validKey(value) {
					return nil, fmt.Errorf("line %d: PrivateKey must be a base64-encoded 32-byte key", lineNo)
				}
				cfg.PrivateKey = value
			case "Address":
				v, err := parsePrefixes(value)
				if err != nil {
					return nil, fmt.Errorf("line %d: Address: %w", lineNo, err)
				}
				cfg.Addresses = v
			case "DNS":
				v, err := parseDNS(value)
				if err != nil {
					return nil, fmt.Errorf("line %d: DNS: %w", lineNo, err)
				}
				cfg.DNS = v
			case "ListenPort":
				v, err := parsePort(value, true)
				if err != nil {
					return nil, fmt.Errorf("line %d: ListenPort: %w", lineNo, err)
				}
				cfg.ListenPort = v
			case "MTU", "Table", "SaveConfig", "PreUp", "PostUp", "PreDown", "PostDown":
				// Accepted as standard wg-quick interface directives. This parser does not execute them.
			default:
				return nil, fmt.Errorf("line %d: unsupported Interface setting %q", lineNo, key)
			}
			continue
		}

		peer := &cfg.Peers[len(cfg.Peers)-1]
		switch key {
		case "PublicKey":
			if !validKey(value) {
				return nil, fmt.Errorf("line %d: PublicKey must be a base64-encoded 32-byte key", lineNo)
			}
			peer.PublicKey = value
		case "PresharedKey":
			if !validKey(value) {
				return nil, fmt.Errorf("line %d: PresharedKey must be a base64-encoded 32-byte key", lineNo)
			}
			peer.PresharedKey = value
		case "AllowedIPs":
			v, err := parsePrefixes(value)
			if err != nil {
				return nil, fmt.Errorf("line %d: AllowedIPs: %w", lineNo, err)
			}
			peer.AllowedIPs = v
		case "Endpoint":
			host, port, err := net.SplitHostPort(value)
			if err != nil || strings.TrimSpace(host) == "" {
				return nil, fmt.Errorf("line %d: invalid Endpoint", lineNo)
			}
			if _, err := parsePort(port, false); err != nil {
				return nil, fmt.Errorf("line %d: invalid Endpoint port", lineNo)
			}
			peer.Endpoint = value
		case "PersistentKeepalive":
			n, err := strconv.Atoi(value)
			if err != nil || n < 0 || n > 65535 {
				return nil, fmt.Errorf("line %d: PersistentKeepalive must be between 0 and 65535", lineNo)
			}
			peer.PersistentKeepalive = n
		default:
			return nil, fmt.Errorf("line %d: unsupported Peer setting %q", lineNo, key)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read configuration: %w", err)
	}
	if !interfaceSeen || cfg.PrivateKey == "" || len(cfg.Addresses) == 0 {
		return nil, fmt.Errorf("configuration requires one Interface section with PrivateKey and Address")
	}
	if len(cfg.Peers) == 0 {
		return nil, fmt.Errorf("configuration requires at least one Peer section")
	}
	seenPeerKeys := make(map[string]int, len(cfg.Peers))
	for i, p := range cfg.Peers {
		if p.PublicKey == "" || len(p.AllowedIPs) == 0 {
			return nil, fmt.Errorf("Peer %d requires PublicKey and AllowedIPs", i+1)
		}
		if previous, exists := seenPeerKeys[p.PublicKey]; exists {
			return nil, fmt.Errorf("Peer %d duplicates the public key from Peer %d", i+1, previous)
		}
		seenPeerKeys[p.PublicKey] = i + 1
	}
	return cfg, nil
}
