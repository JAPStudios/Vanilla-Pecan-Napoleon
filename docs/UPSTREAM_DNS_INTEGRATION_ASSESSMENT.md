# DNS upstream integration assessment

Review date: 2026-11-26

## Scope

Vanilla Pecan Napoleon is the secure VPN software project for Sweet Shop: Internet Café. DNS forwarding and encrypted resolver transport are relevant to the VPN client, but an upstream DNS product is not automatically a reusable VPN component. This assessment records source-level findings from the requested AdGuardTeam, Home Assistant, IVPN, Mullvad, NextDNS, Pi-hole and Windscribe organizations.

## Windscribe/ctrld

Repository: https://github.com/Windscribe/ctrld

- The repository is MIT-licensed.
- Its module declares Go 1.23.0 and toolchain Go 1.23.7, which is compatible with the current Go 1.23 language baseline in this repository.
- It supports classic DNS, DNS-over-HTTPS, DNS-over-TLS, DNS-over-QUIC and DNS-over-HTTP/3.
- The project is a complete DNS forwarding proxy/service. Its root package exposes configuration types, while its runtime and resolver orchestration are coupled to a substantial dependency graph and service-specific configuration.
- Its defaults and examples include Control D resolver endpoints. Those defaults must not be inherited as Vanilla Pecan Napoleon policy.
- The source tree contains operating-system service management, DNS listener setup, network discovery, interface changes, metrics and client-identification features. Importing the full application would add responsibilities and behaviour outside the VPN engine boundary.

Decision: do not vendor the complete ctrld application or add the whole module as a dependency at this stage. Continue reviewing whether a narrowly scoped, independently testable DNS transport component can be reused without importing unrelated service behaviour. Any reuse must retain the MIT copyright and licence notice and must be reviewed with its transitive dependencies.

## AdGuardTeam/dnsproxy

Repository: https://github.com/AdguardTeam/dnsproxy

The current default-branch go.mod declares Go 1.27.4. That is above this repository's Go 1.23 baseline. Its Apache-2.0 licence is permissive, but the toolchain mismatch and broad dependency set make direct adoption unsuitable without a deliberate Go baseline change and dependency review.

## NextDNS

Repository: https://github.com/nextdns/nextdns

The current default-branch go.mod declares Go 1.26.0, also above the current Go 1.23 baseline. The repository is MIT-licensed and contains a useful DoH implementation, but importing its resolver package would bring a newer Go requirement and its surrounding cache, endpoint and resolver abstractions. Reuse should be considered only after isolating a minimal component and reviewing its dependency requirements.

## Existing WireGuard engine

Mullvad's MIT-licensed wireguard-go device engine is already integrated as a pinned dependency. The current adapter accepts a caller-supplied TUN device and bind. Platform TUN creation, operating-system address and route configuration, DNS policy, provisioning and gateway operation remain unimplemented.

## Integration rules

1. Prefer narrowly scoped components with clear interfaces over importing complete desktop clients or service applications.
2. Verify the exact licence and all transitive dependency licences before copying or vendoring source.
3. Verify the upstream Go/Rust toolchain requirement against the repository's supported build matrix.
4. Do not inherit third-party resolver endpoints, account systems, telemetry or client-identification defaults without an explicit product requirement and privacy review.
5. Preserve required copyright, licence and NOTICE files for every reused component.
6. Add focused tests and CI coverage for each integration before merging it into main.

## Current decision

No new upstream DNS implementation is copied by this assessment commit. The existing Mullvad WireGuard engine remains the only integrated upstream implementation. The DNS candidates above are documented for further component-level evaluation; this is a technical integration decision, not a claim that all requested upstream repositories have been imported or merged.
