# Upstream review — Vanilla Pecan Napoleon

This is a preliminary relevance and license review of repositories surfaced from the seven requested GitHub organizations. It is not a complete inventory of every repository in those organizations.

## Integration rule

Vanilla Pecan Napoleon is the VPN project for Sweet Shop: Internet Café. Only components that directly support VPN tunnelling, client/daemon operation, DNS resolution, DNS filtering, or privacy/security controls should be considered. Repositories are not merged wholesale merely because they belong to a requested organization. Each imported component must retain required notices and comply with its exact license and dependency licenses.

The connected GitHub account has read access but no push access to the reviewed upstream repositories. Changes therefore belong in this destination repository; they cannot be pushed into the upstream organizations.

## Preliminary candidates

| Organization | Candidate repositories | Relevance |
|---|---|---|
| AdguardTeam | AdGuardHome, dnsproxy, AdGuardVPNExtension | DNS resolution/filtering and VPN client integration |
| home-assistant | plugin-dns, app-dnsmasq, app-duckdns | DNS add-ons for Home Assistant; not a VPN engine |
| ivpn | desktop-app, android-app, ios-app, moddns | VPN client/daemon and DNS components |
| mullvad | mullvadvpn-app, wireguard-go, gotatun, dns-blocklists, hickory-dns | VPN clients, WireGuard engine, tunnel and DNS components |
| nextdns | nextdns, blocklists | DNS-over-HTTPS client and filtering data |
| pi-hole | pi-hole, FTL, dnsmasq, docker-pi-hole | DNS filtering and resolver components |
| Windscribe | Desktop-App, wireguard-go, boringtun, amneziawg-go, wsnet, ctrld, unbound | VPN clients, tunnel implementations, networking and DNS |

## License findings verified so far

- mullvad/wireguard-go: MIT license verified. The reviewed Mullvad branch points to commit 86ccb9b49b8f095775109b62554efedbc1d0e3c6 (2025-02-19). This is a candidate tunnel-engine source, not yet imported.
- nextdns/nextdns: MIT license verified.
- pi-hole/pi-hole: EUPL-1.2 license verified.
- Windscribe/Desktop-App: GPL-2.0 license verified.
- The root license/readme paths for some AdGuard, IVPN and Home Assistant repositories were not available through the initial file lookups. Their licenses remain unverified here; no license is inferred.

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

## Additional candidate checks

- AdguardTeam/dnsproxy: Apache-2.0; DNS proxy supporting DoH, DoT, DoQ and DNSCrypt. Its current go.mod declares Go 1.27.4, while this destination currently declares Go 1.23. It is a relevant DNS component but cannot be adopted unchanged without resolving the Go toolchain and dependency compatibility.
- home-assistant/plugin-dns: Apache-2.0; a CoreDNS implementation packaged for Home Assistant. It is not a general-purpose VPN engine and is not a direct fit for the current VPN core.
- ivpn/desktop-app-daemon: GPL-3.0 and archived. It is an older desktop-app service, not a safe default to import wholesale.
- ivpn/moddns: GPL-3.0; a multi-service hosted DNS filtering product. Its README says self-hosting is not supported. Reusing its architecture or code would require a separate GPL-3.0 compatibility and service-design review.
- mullvad/wireguard-go: MIT; Go implementation of WireGuard with Mullvad-specific DAITA changes. Its module path is golang.zx2c4.com/wireguard, so dependency pinning and the intended fork must be handled explicitly.
- nextdns/nextdns: MIT; a DNS-over-HTTPS client. Its current go.mod declares Go 1.26.0, above this destination's current Go 1.23 declaration.
- pi-hole/FTL: the GitHub API reports NOASSERTION for its license. Do not copy or vendor it until the actual licensing terms are established.
- Windscribe/boringtun: BSD-3-Clause; a fork of Cloudflare's userspace WireGuard implementation. It supplies the WireGuard protocol implementation, not the complete network/tunnel stack.
- Windscribe/wireguard-go: MIT; the repository describes itself as a mirror and points to the official wireguard-go repository. Prefer a maintained upstream or a specifically justified fork rather than importing a mirror without need.
- Windscribe/ctrld: MIT; a DNS forwarding proxy. Its exact role and dependencies would need evaluation before integration.

## Current destination implementation

The first implementation commit adds a WireGuard configuration parser/validator, CLI validation command, tests and Go CI. No upstream source files have been copied into the destination yet. The CI run for that foundation passed. This is not yet a functioning VPN client or deployed VPN service.

## Next engineering gates

1. Complete source-level review of relevant repositories, not just repository metadata.
2. Select one tunnel engine and platform adapters; avoid importing multiple competing implementations without a concrete requirement.
3. Resolve Go/Rust toolchain versions, dependency trees, supported platforms, license obligations and security history before adopting code.
4. Add authenticated provisioning, key rotation/revocation, DNS policy, gateway configuration, leak protection, reconnect handling and observability.
5. Add platform-specific builds and integration tests before any production-security claim.
6. Complete independent security review and deployment validation before service launch.
