# Accepted design

This index records InfraOps' owner-accepted design choices and their scope.
The area files carry the same weight as [AGENTS.md](../../AGENTS.md).

Read the relevant area files before proposing a change or recording a finding
or audit entry. Revisit a decision when the owner asks or new evidence changes
the documented contract or tradeoff. Record owner-agreed choices in the owning
area file.

- [Applications](accepted-design/applications.md) covers network policy,
  secret references, resource defaults, and ConfigMap sources.
- [Providers](accepted-design/providers.md) covers Cloudflare account selection,
  GitHub branch protection, and repository creation dependencies.
- [Configuration](accepted-design/configuration.md) covers helper CLI test scope,
  HJSON writes, and protobuf validation.
- [Deployment](accepted-design/deployment.md) covers release handoff, CircleCI
  gates, and manual Kubernetes add-ons.
