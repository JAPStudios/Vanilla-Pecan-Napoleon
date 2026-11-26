# Vanilla Pecan Napoleon

Vanilla Pecan Napoleon is the secure VPN software project for Sweet Shop: Internet Café. It is intended to provide VPN connectivity across the services and technologies operated by Sweet Shop: Internet Café.

This repository is the implementation project, not the VPN network itself. A production service requires reviewed client software, authenticated provisioning, key lifecycle management, VPN gateways, routing, DNS, monitoring, abuse controls, platform-specific tunnel integration, deployment automation and independent security testing.

## Current implementation

The initial implementation provides a WireGuard configuration parser and validator. It checks interface and peer configuration without printing private or preshared keys. The development branch also adds an adapter around Mullvad's MIT-licensed wireguard-go device engine. The adapter accepts a caller-supplied platform TUN device and bind, translates validated configuration to WireGuard's UAPI, and starts or closes the engine. It does not create platform TUN devices, configure operating-system addresses or routes, provision accounts, operate gateways or claim production security.

## Validate a configuration

```sh
go run ./cmd/vpnctl validate /path/to/wireguard.conf
go test ./...
```

The validator accepts standard WireGuard interface and peer fields. It rejects duplicate fields, malformed keys, invalid addresses and invalid peer endpoints.

## Upstream source review

Potential upstream components are evaluated individually. No upstream repository is modified by this project. Mullvad's wireguard-go is now pinned as a Go module dependency at revision 86ccb9b49b8f095775109b62554efedbc1d0e3c6 after reviewing its MIT license, module requirements and device/UAPI interfaces. Its notice is recorded in THIRD_PARTY_NOTICES.md. Other components remain subject to individual license, dependency, platform-support and security review before integration. Upstream projects remain separately maintained.

Initial candidate areas:
- WireGuard tunnel implementation and platform adapters.
- VPN client and daemon architecture.
- DNS-over-HTTPS and resolver integration.
- DNS filtering and blocklist formats.

Home Assistant components are not part of the VPN core merely because they appear in the requested source-organization list.

## Security status

This repository is an early development foundation. The validator is not a VPN client or server and must not be represented as a deployed or independently audited security product.
