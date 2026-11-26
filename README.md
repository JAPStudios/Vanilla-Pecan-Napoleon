# Vanilla Pecan Napoleon

Vanilla Pecan Napoleon is the secure VPN software project for Sweet Shop: Internet Café. It is intended to provide VPN connectivity across the services and technologies operated by Sweet Shop: Internet Café.

This repository is the implementation project, not the VPN network itself. A production service requires reviewed client software, authenticated provisioning, key lifecycle management, VPN gateways, routing, DNS, monitoring, abuse controls, platform-specific tunnel integration, deployment automation and independent security testing.

## Current implementation on main

The project provides a WireGuard configuration parser and validator. It checks interface and peer configuration without printing private or preshared keys. It rejects duplicate peer public keys and validates configuration fields.

The project also includes an adapter around Mullvad's MIT-licensed wireguard-go device engine. The adapter accepts a caller-supplied platform TUN device and bind, translates validated configuration to WireGuard's UAPI, and starts or closes the engine. The pinned module revision and dependency checksums are committed in go.mod and go.sum.

A platform startup helper creates a native WireGuard TUN device on Linux, macOS and Windows and starts the engine with the platform default bind. Address assignment, operating-system routes, DNS settings and platform permissions remain separate responsibilities.

The new internal DoH transport sends DNS wire-format messages to a configured HTTPS resolver using RFC 8484 POST requests. It validates the endpoint, request and response sizes, HTTP status and response media type. It does not select or hard-code a DNS provider, send device identifiers, change system DNS settings or log DNS queries. Its design was informed by source review of Windscribe's MIT-licensed ctrld resolver; no ctrld source files are copied or vendored.

The adapter does not provision accounts, operate gateways or claim production security.

## Validate a configuration

```sh
go run ./cmd/vpnctl validate /path/to/wireguard.conf
go test ./...
```

The validator accepts standard WireGuard interface and peer fields. It rejects duplicate fields, malformed keys, invalid addresses and invalid peer endpoints.

## Upstream source review

Potential upstream components are evaluated individually. No upstream repository is modified by this project. Mullvad's wireguard-go is integrated as a pinned Go module dependency at revision `86ccb9b49b8f095775109b62554efedbc1d0e3c6` after reviewing its MIT license, module requirements and device/UAPI interfaces. Its notice is recorded in `THIRD_PARTY_NOTICES.md`, and dependency checksums are committed in `go.sum`.

Other components remain subject to individual source, licence, dependency, platform-support and security review before integration. Upstream projects remain separately maintained. The requested organizations contain many unrelated repositories; importing all of them indiscriminately would not produce a coherent VPN implementation.

Initial candidate areas:
- WireGuard tunnel implementation and platform adapters.
- VPN client and daemon architecture.
- DNS-over-HTTPS and resolver integration.
- DNS filtering and blocklist formats.

Home Assistant DNS add-ons are not part of the VPN core merely because they appear in the requested source-organization list.

## Security status

This repository is an early development foundation. The validator, WireGuard engine adapter, native TUN startup and DoH transport are not a complete VPN client or server and must not be represented as a deployed or independently audited security product.
