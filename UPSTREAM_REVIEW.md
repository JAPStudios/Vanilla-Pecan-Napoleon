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

- mullvad/wireguard-go: MIT license verified. The reviewed mullvad branch points to commit 86ccb9b49b8f095775109b62554efedbc1d0e3c6 (2025-02-19). This is a candidate tunnel-engine source, not yet imported.
- nextdns/nextdns: MIT license verified.
- pi-hole/pi-hole: EUPL-1.2 license verified.
- Windscribe/Desktop-App: GPL-2.0 license verified.
- The root license/readme paths for some AdGuard, IVPN and Home Assistant repositories were not available through the initial file lookups. Their licenses remain unverified here; no license is inferred.

## Current destination implementation

The first implementation commit adds a WireGuard configuration parser/validator, CLI validation command, tests and Go CI. No upstream source files have been copied into the destination yet. The CI run for that foundation passed. This is not yet a functioning VPN client or deployed VPN service.

## Next engineering gates

1. Complete the repository-by-repository inventory and license/dependency review.
2. Select the WireGuard tunnel engine and platform adapters, then implement tunnel lifecycle and operating-system address/route management.
3. Add authenticated provisioning, key rotation/revocation, DNS policy, gateway configuration, leak protection, reconnect handling and observability.
4. Add platform-specific builds and integration tests before any production-security claim.
5. Complete independent security review and deployment validation before service launch.
