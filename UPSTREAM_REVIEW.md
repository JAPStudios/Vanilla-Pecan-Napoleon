# Upstream review — Vanilla Pecan Napoleon

This is a relevance, source and licence review of repositories surfaced from the seven requested GitHub organizations. It is not a complete source-level audit of every repository in those organizations.

## Integration rule

Vanilla Pecan Napoleon is the VPN project for Sweet Shop: Internet Café. Only components that directly support VPN tunnelling, client/daemon operation, DNS resolution, DNS filtering, or privacy/security controls should be considered. Repositories are not merged wholesale merely because they belong to a requested organization. Each imported component must retain required notices and comply with its exact license and dependency licenses.

The connected GitHub account has read access but no push access to the reviewed upstream repositories. Changes therefore belong in this destination repository; upstream repositories remain separately maintained.

## Candidate map across the seven organizations

| Organization | Relevant repositories surfaced | Assessment |
|---|---|---|
| AdguardTeam | AdGuardHome, dnsproxy, AdGuardSDNSFilter, AdGuardVPNExtension, AdGuardVPNCLI, AdGuardVPNForAndroid, AdGuardVPNForWindows, AdGuardVPNForiOS, AdGuardVPNForMac, DnsLibs | DNS filtering, DNS proxying and VPN clients. Review each platform/client separately. |
| home-assistant | plugin-dns, app-dnsmasq, app-duckdns | Home Assistant add-ons, not general-purpose VPN components. No direct VPN-core integration identified. |
| ivpn | desktop-app, desktop-app-ui, desktop-app-daemon, android-app, ios-app, desktop-app-cli, wireguard-apple, moddns | VPN client/daemon and platform-specific code. GPL components and archived projects need exact compatibility review before reuse. |
| mullvad | mullvadvpn-app, wireguard-go, gotatun, dns-blocklists, encrypted-dns-profiles, wireguard-apple, wireguard-rs, wireguard-nt, hickory-dns, coredns | VPN clients, tunnel engines, platform adapters and DNS components. Mullvad wireguard-go is integrated. |
| nextdns | nextdns, blocklists | DNS-over-HTTPS client and filtering data. |
| pi-hole | FTL, dnsmasq, docker-pi-hole, pi-hole, dnsblast, massdns | DNS filtering and resolver components. FTL's reported license metadata requires further verification. |
| Windscribe | Desktop-App, wireguard, wireguard-apple, boringtun, wireguard-go, wireguard-android, wsnet, amneziawg-go, ctrld, openvpn, unbound | VPN clients, tunnel implementations, networking and DNS. |

## Repository inventory snapshot — 2026-11-26

A GitHub repository search scoped to each named organization, using `stars:>=0` and all returned pages, surfaced 495 public repositories in total:

| Organization | Repositories returned |
|---|---:|
| AdguardTeam | 123 |
| home-assistant | 131 |
| ivpn | 17 |
| mullvad | 121 |
| nextdns | 5 |
| pi-hole | 23 |
| Windscribe | 75 |
| Total | 495 |

This is an inventory of repository records returned by the GitHub search endpoint, not a source-code, dependency, or security audit of all 495 projects. The full set includes unrelated applications, websites, build tools, libraries, operating-system components, archived projects, forks and data repositories. It is not technically or legally sound to merge all of them into a single VPN codebase.

## Source and licence findings

- mullvad/wireguard-go: MIT. The integrated module is pinned to revision `86ccb9b49b8f095775109b62554efedbc1d0e3c6` (2025-02-19). Its module path is `golang.zx2c4.com/wireguard`; the destination uses an explicit module replacement to the Mullvad fork. The licence notice is recorded in `THIRD_PARTY_NOTICES.md`.
- AdguardTeam/dnsproxy: Apache-2.0; DNS proxy supporting DoH, DoT, DoQ and DNSCrypt. Its current `go.mod` declares Go 1.27.4, above the destination's Go 1.23 baseline. Relevant, but not adopted unchanged.
- AdguardTeam/AdGuardSDNSFilter: DNS filtering data repository; data format, update process and licence terms need a separate review before consumption.
- home-assistant/plugin-dns: Apache-2.0; a CoreDNS implementation packaged for Home Assistant. It is not a general-purpose VPN engine.
- ivpn/desktop-app-daemon: GPL-3.0 and archived. It is an older desktop-app service, not a safe default to import wholesale.
- ivpn/moddns: GPL-3.0; its README describes a multi-service hosted DNS filtering product and says self-hosting is not supported. A code reuse proposal needs separate GPL-3.0 compatibility and service-design review.
- mullvad/gotatun: MPL-2.0 Rust WireGuard implementation. Its current workspace declares Rust 1.95; it is a distinct alternative engine and not needed alongside the already integrated Go engine without a concrete platform requirement.
- nextdns/nextdns: MIT; DNS-over-HTTPS client. Its current `go.mod` declares Go 1.26.0, above the destination's Go 1.23 baseline.
- pi-hole/FTL: GitHub license metadata reports NOASSERTION. Do not copy or vendor it until the actual licensing terms are established.
- pi-hole/pi-hole: EUPL-1.2.
- Windscribe/boringtun: BSD-3-Clause; a fork of Cloudflare's userspace WireGuard implementation. It supplies a WireGuard protocol implementation, not the complete network/tunnel stack.
- Windscribe/wireguard-go: MIT; the repository describes itself as a mirror. Prefer the maintained, specifically selected upstream/fork rather than importing a mirror without need.
- Windscribe/ctrld: MIT; its repository is a DNS forwarding proxy supporting traditional DNS, DoH, DoT, DoQ and DoH3. Its Go module declares Go 1.23.0. It is a plausible DNS integration candidate, but its substantial dependency tree, service model and package boundaries need review before integration.
- Windscribe/Desktop-App: GPL-2.0. Its full desktop client is not suitable for wholesale incorporation without resolving the resulting distribution and licence obligations.

Where a repository's exact licence has not been verified, no licence is inferred.

## Current destination implementation

- The configuration parser/validator and CLI validation command are committed.
- Duplicate peer public keys are rejected, with regression coverage.
- Mullvad's WireGuard engine adapter is merged into `main` in commit `86fa311cb445a31949aeae95e0efb6f07f754855`.
- The Go dependency checksum file `go.sum` is committed in merge commit `737356fa9744401043086b0dd7309a4f2ddb7751`.
- The CI workflow runs `go mod tidy`, `go test ./...` and `go vet ./...`. Both PR checks for the checksum change passed. The post-merge main-branch check is being verified separately.

The engine adapter accepts a caller-supplied platform TUN device and bind. It does not create platform TUN devices, configure operating-system addresses or routes, provision users/accounts, operate VPN gateways or constitute a deployed VPN service.

## Next engineering gates

1. Complete source-level review of each relevant repository and its transitive dependencies.
2. Select platform TUN implementations and operating-system route/address configuration for supported platforms.
3. Evaluate DNS proxy and filtering integration, beginning with the compatible Go 1.23 candidate rather than importing a complete unrelated VPN client.
4. Resolve dependency versions, licences, supported platforms and security history before adopting code.
5. Implement authenticated provisioning, key rotation/revocation, DNS policy, gateway configuration, leak protection, reconnect handling and observability.
6. Add platform-specific builds and integration tests before any production-security claim.
7. Complete independent security review and deployment validation before service launch.
