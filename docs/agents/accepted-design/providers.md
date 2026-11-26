# Accepted design for providers

Part of the [accepted-design index](../accepted-design.md).

## Cloudflare account

Cloudflare resources require `CLOUDFLARE_ACCOUNT_ID`. The Cloudflare package
reads it at package scope so account-scoped resource constructors can share the
same input.

[Cloudflare account input](../../../internal/cf/cf.go)

## GitHub branch protection

The solo-maintainer workflow requires zero approving pull request reviews on
`master`. Required status checks are the primary merge gate.

[GitHub branch protection](../../../internal/gh/protection.go)

## GitHub repository creation

GitHub Pages and collaborators can be enabled in the initial repository
configuration. Vulnerability alerts, Pages, and collaborators consume repository
outputs so Pulumi waits for repository creation within the same update.

The live workflow was verified by creating an auto-initialized repository with
Pages and a collaborator invitation, then removing its configuration and
applying the update that deleted the repository.

[Repository settings](../../../internal/gh/repository.go) ·
[Collaborators](../../../internal/gh/collaborator.go)
