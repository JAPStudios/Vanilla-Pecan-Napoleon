package main

import (
	"fmt"
	"os"

	"github.com/JAPStudios/Vanilla-Pecan-Napoleon/internal/wgconfig"
)

func main() {
	if len(os.Args) != 3 || os.Args[1] != "validate" {
		fmt.Fprintln(os.Stderr, "usage: vpnctl validate <wireguard-config>")
		os.Exit(2)
	}
	f, err := os.Open(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, "cannot open configuration:", err)
		os.Exit(1)
	}
	defer f.Close()

	cfg, err := wgconfig.Parse(f)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid WireGuard configuration:", err)
		os.Exit(1)
	}
	fmt.Printf("Valid WireGuard configuration: %d interface address(es), %d peer(s).\n", len(cfg.Addresses), len(cfg.Peers))
}
