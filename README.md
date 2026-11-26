# Vanilla Pecan Napoleon

Vanilla Pecan Napoleon is the secure VPN software project for Sweet Shop: Internet Café. It is intended to provide VPN connectivity across the services and technologies operated by Sweet Shop: Internet Café.

This repository is the implementation project, not the VPN network itself. A production service requires reviewed client software, authenticated provisioning, key lifecycle management, VPN gateways, routing, DNS, monitoring, abuse controls, platform-specific tunnel integration, deployment automation and independent security testing.

## Current implementation

The initial code provides a WireGuard configuration parser and validator. It checks interface and peer configuration without printing private or preshared keys. It does not yet establish tunnels, provision accounts, operate gateways or claim production security.

## Validate a configuration

```sh
go run ./cmd/vpnctl validate /path/to/wireguard.conf
go test ./...
```

The validator accepts standard WireGuard interface and peer fields. It rejects duplicate fields, malformed keys, invalid addresses and invalid peer endpoints.

## Upstream source review

Potential upstream components are being evaluated individually. No upstream repository is modified by this project, and no source code is imported until its exact license, notices, dependencies, platform support and security implications have been reviewed. Upstream projects remain separately maintained.

Initial candidate areas:
- WireGuard tunnel implementation and platform adapters.
- VPN client and daemon architecture.
- DNS-over-HTTPS and resolver integration.
- DNS filtering and blocklist formats.

Home Assistant components are not part of the VPN core merely because they appear in the requested source-organization list.

## Security status

This repository is an early development foundation. The validator is not a VPN client or server and must not be represented as a deployed or independently audited security product.
