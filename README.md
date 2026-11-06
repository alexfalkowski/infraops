![Gopher](assets/gopher.png)
[![DigitalOcean Referral Badge](https://web-platforms.sfo2.cdn.digitaloceanspaces.com/WWW/Badge%202.svg)](https://www.digitalocean.com/?refcode=b80d3d6467e1&utm_campaign=Referral_Invite&utm_medium=Referral_Program&utm_source=badge)

[![CircleCI](https://circleci.com/gh/alexfalkowski/infraops.svg?style=shield)](https://circleci.com/gh/alexfalkowski/infraops)
[![codecov](https://codecov.io/gh/alexfalkowski/infraops/graph/badge.svg?token=U3X5JGAA8I)](https://codecov.io/gh/alexfalkowski/infraops)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexfalkowski/infraops/v2)](https://goreportcard.com/report/github.com/alexfalkowski/infraops/v2)
[![Go Reference](https://pkg.go.dev/badge/github.com/alexfalkowski/infraops/v2.svg)](https://pkg.go.dev/github.com/alexfalkowski/infraops/v2)
[![Stability: Active](https://masterminds.github.io/stability/active.svg)](https://masterminds.github.io/stability/active.html)

# 🧭 InfraOps

Infrastructure configuration and deployment for personal services, built with Go and Pulumi.
Operators manage Kubernetes applications, Cloudflare, DigitalOcean, and GitHub from this repository.
For application features, runtime settings, and usage, see each application's own README.

## 🗺️ Overview

| Area | Manages | Configuration |
| ---- | ------- | ------------- |
| `apps` | Kubernetes application resources | [apps.hjson](area/apps/apps.hjson) |
| `cf` | Cloudflare zones, DNS, and R2 buckets | [cf.hjson](area/cf/cf.hjson) |
| `do` | DigitalOcean VPC and Kubernetes cluster | [do.hjson](area/do/do.hjson) |
| `gh` | GitHub repositories, settings, and branch protection | [gh.hjson](area/gh/gh.hjson) |
| `k8s` | Cluster add-ons installed manually with Helm and kubectl | [Makefile](area/k8s/Makefile), [nginx values](area/k8s/nginx/values.yaml) |

Each Pulumi area has its own `main.go`, `Pulumi.yaml`, and HJSON configuration under `area/<name>/`.
Shared resource logic lives in `internal/`.

```mermaid
flowchart LR
  HJSON["Area HJSON"] --> Config["Protobuf config"]
  Schema["service.proto"] -.-> Config
  Config --> Models["Internal Go models"]
  Models --> Pulumi["Pulumi resources"]
```

The [protobuf schema](api/infraops/v2/service.proto) owns field definitions, defaults, and deployment
conventions. The top-level HJSON `version` identifies the configuration schema. Semantic errors may
still surface during Pulumi or provider operations.

## 🚚 Bootstrap

Install Go at the version declared in [go.mod](go.mod), Make, and Ruby for the shared help/build helpers.
Initialize the `bin` tooling submodule, then download, tidy, and vendor dependencies:

```bash
git submodule sync
git submodule update --init
make dep
make help
```

Run `make dep` before local validation. Use `make help` for the command catalog and the
[CircleCI configuration](.circleci/continue_config.yml) for validation requirements. Go checks use
`betteralign`, `golangci-lint`, `gotestsum`, `govulncheck`, and Trivy; shell checks use ShellCheck;
API checks use Buf. The golangci-lint wrapper skips its check when the tool is absent from `PATH`.

When changing [service.proto](api/infraops/v2/service.proto), run `make api-generate` and include the
updated `service.pb.go`; do not edit generated Go directly. `make api-breaking` compares against
GitHub `master` and requires network access.

## 🚀 Preview and Deploy

Install the Pulumi CLI and configure access for the selected area:

| Area | Required access |
| ---- | --------------- |
| `apps` | kubectl and the intended kubeconfig context; `jq` for update selection |
| `cf` | `CLOUDFLARE_API_TOKEN` and `CLOUDFLARE_ACCOUNT_ID` |
| `do` | `DIGITALOCEAN_TOKEN` |
| `gh` | `GITHUB_TOKEN` |

From the repository root, log in and preview the area before applying it:

```bash
make pulumi-login
make area=cf pulumi-preview
make area=cf pulumi-update
```

Use `area=apps|cf|do|gh`. Targets select stack `alexfalkowski/<area>/prod` and working directory
`area/<area>`, where the program reads its HJSON file. The `k8s` add-ons have a separate manual workflow
below.

> [!WARNING]
> Update, refresh, and cancel targets change remote infrastructure or stack state and pass `--yes`.
> Review the preview before applying changes.

For stack recovery, the root Makefile also exposes `pulumi-refresh` and `pulumi-cancel`.

> [!CAUTION]
> `make area=cf pulumi-delete` runs `pulumi stack rm --force`. It removes stack state and can orphan
> resources that still exist.

## 🪣 R2 Object Expiration

A bucket can opt into object-age expiration with `expiration`. For example,
this illustrative Cloudflare configuration expires objects after seven days:

```hjson
{
  version: "2.0"
  buckets: [
    {
      name: example-archive
      region: EEUR
      expiration: 168h
    }
  ]
}
```

Use `24h` for a one-day testing example. These are examples, not an agreed einblick
retention policy. Values use [Go duration syntax](https://pkg.go.dev/time#ParseDuration)
and are parsed into `time.Duration`: use hours for days (`168h`, not `7d`). A configured
value must be a positive whole number of seconds within that type's range; empty,
malformed, overflowing, zero, negative and fractional-second values are rejected before
bucket registration. Durations such as `1.5h` are valid because they resolve to whole seconds.
Omitting the field creates no lifecycle resource and enables no object expiration.
Omitting `zone` keeps the bucket without a public custom domain.
See the [Bucket schema](api/infraops/v2/service.proto).

The rule covers every object, using the stored object's age, regardless of timestamps
inside its contents or UTC day prefixes in its key. It does not enforce strict event-age
retention. Cloudflare applies lifecycle changes asynchronously to existing objects;
deletion typically occurs within 24 hours after expiration and can take longer.
Shortening the age can make existing objects eligible for deletion; extending it cannot
restore deleted objects. See [Cloudflare object lifecycles](https://developers.cloudflare.com/r2/buckets/object-lifecycles/).

The [pinned provider's lifecycle resource](https://www.pulumi.com/registry/packages/cloudflare/api-docs/r2bucketlifecycle/)
replaces the complete bucket rule set. This capability manages all-object expiration
and retains R2's default seven-day incomplete multipart upload cleanup; it does not merge
other rules created outside infraops. Review existing rules before adopting it.
Lifecycle management requires a Cloudflare token with `Workers R2 Storage Write` permission.
The [Cloudflare lifecycle API](https://developers.cloudflare.com/api/resources/r2/subresources/buckets/subresources/lifecycle/methods/update/)
defines replacement, empty-prefix matching, and age in seconds.

> [!WARNING]
> Removing a previously applied `expiration` value removes the lifecycle resource
> from Pulumi state but **leaves the remote rules active**. The upstream provider pinned by
> Cloudflare v6.23.2 has a [no-op lifecycle Delete](https://github.com/cloudflare/terraform-provider-cloudflare/blob/9b51a582458f5c5d46d8547ed2604fcba91fb861/internal/services/r2_bucket_lifecycle/resource.go).
> To stop expiration, first disable or remove the `expire-objects` rule in Cloudflare's
> bucket lifecycle settings, then omit the field and review the preview. Do not apply the
> old configured value again: it reinstates expiration. Check the remote rules afterward;
> removal does not restore deleted objects or guarantee immediate cessation of pending deletion.
> The rule also carries incomplete multipart cleanup; retain that cleanup separately if needed.

Changing only the expiration value updates the same lifecycle resource and keeps bucket
and optional custom-domain identities. Reverting to older code/config that omits this
capability has the same remote-rule persistence risk. Pulumi mocks verify configuration
and resource inputs; live propagation and deletion require separately authorized R2 validation.

## 📦 Application Deployment

[apps.hjson](area/apps/apps.hjson) selects the applications to deploy. Each app's configuration file
lives at `area/apps/<namespace>/<app>.yaml`; Pulumi puts it in a ConfigMap and mounts it at
`/etc/<app>/<app>.yaml`. Files under `area/apps/fiction/` are examples unless a matching entry in
`apps.hjson` references them. Consult each application's README for the contents of its YAML config.

When adding an app, also extend [area/apps/release](area/apps/release) with its supported name,
verification and rollout commands, and load eligibility. Add load fixtures only when they exercise
useful application behavior. The release helpers currently support the `lean` namespace.

App NetworkPolicies limit ingress to ports `6060`, `8080`, and `9090` and leave egress open.
Set the intended `replicas` explicitly: omission means zero pods. A single replica can become
unavailable during rollout or disruption because the rollout and PodDisruptionBudget allow one
unavailable pod. See the [Application schema](api/infraops/v2/service.proto) for deployment kinds,
resource sizing, probes, and secret-reference syntax.

### 💾 Persistent Application Storage

Before enabling storage, update the shared release image to a storage-aware `cmd/bump`
version and update all consuming app pipelines to that image. Older helpers discard
storage settings when rewriting the configuration, even when releasing another app.

An application can opt into one persistent filesystem claim. Add this block to that
application's HJSON object; the values below are illustrative, not production sizing:

```hjson
storage: {
  capacity: 1Gi
  path: "/data"
  profile: example-csi
}
```

Choose capacity and a storage provisioning profile for the installation before deploying.
The Kubernetes adapter maps `profile` to an existing StorageClass. Omit
`profile` to leave class assignment to the cluster default, or set it to `""`
to request no class and bind a matching pre-provisioned PV. Capacity has no default;
Kubernetes/provider validation applies to quantities and available storage. The claim
requests `ReadWriteOnce` and `Filesystem`. See [Storage](api/infraops/v2/service.proto)
for the complete field contract.

The PVC's Kubernetes and Pulumi name is `<app>-data`, in the app's namespace. Pod
replacement, image/config updates and scaling retain that identity. All replicas
share the same claim. `ReadWriteOnce` permits multiple pods on one node; it does not
enforce one writer. Storage leaves the existing rolling update and disruption budget
unchanged. Applications requiring exclusive writers need a separately reviewed
rollout policy before activation; Recreate support is deferred to I-01b.

The entire claim mounts writable at the required absolute, clean, non-root
`path`, without a subPath. Choose a dedicated directory that does not hide
image files, `/etc/<app>` config or `/etc/secrets`. Put the database, engine recovery
files, pending exports and any temporary workspace beneath that mount; the app owns
creating subdirectories and bounding workspace use. No paths or environment variables
are set inside the app. The root filesystem remains read-only, with non-root UID/GID
`10001`, `fsGroup: 2000`, dropped capabilities and no privilege escalation. The volume
driver must support that group access. Pulumi mocks check resource settings; verify
actual non-root writes and recovery after pod replacement on the selected storage
before relying on persistence. [Kubernetes security context guidance](https://kubernetes.io/docs/tasks/configure-pod-container/security-context/)
describes the driver-dependent ownership handling.

Claims use [Pulumi protection](https://www.pulumi.com/docs/iac/concepts/resources/options/protect/).
Removing storage/an application or attempting a claim replacement fails until the
operator deliberately unprotects the claim. Back up or migrate data and review the
preview before doing so. A rename or namespace change changes claim identity and
does not migrate data. Class changes require migration; capacity increases require
class/driver expansion support, and shrinking is unsupported. Changing mount path
keeps the claim but requires aligning the app's paths. Rolling back to config/code
without storage attempts to remove the protected claim and is not an automatic
storage rollback.

Protection applies to Pulumi-managed operations. It does not prevent direct Kubernetes,
namespace or cluster deletion, including `make -C area/apps delete`. The bound PV's
reclaim policy controls whether claim deletion also deletes its backing storage;
this capability does not create a StorageClass or set PV reclaim policy. Inspect
those policies and arrange backups independently. [Kubernetes persistent-volume guidance](https://kubernetes.io/docs/concepts/storage/persistent-volumes/)
owns binding, access, expansion and reclamation behavior. Applications that omit
`storage` retain their existing config/secret-only volumes.

### 🏗️ Prepare the Cluster

With `doctl` installed, save kubeconfig if needed, then create the app namespace once:

```bash
make -C area/apps save-config
make -C area/apps setup
```

Create the Kubernetes Secrets required by each app before deployment. The apps program references
existing `<secretName>-secret` objects; it does not create them. Listed secret volumes mount at
`/etc/secrets/<secretName>`, while secret environment variables select individual keys.

For a new public hostname, preview and update the Cloudflare area before the apps update so public
readiness checks can resolve it. The [nginx Helm values](area/k8s/nginx/values.yaml) restrict load
balancer access to Cloudflare CIDRs; use Cloudflare-proxied DNS for that ingress.

### 🔐 Abgleich Secret

This repository provides a helper to create `abgleich-secret` in the existing `lean` namespace.
Use the intended kubeconfig context and three readable, nonempty credential files outside the repository:

```bash
make -C area/apps abgleich-secret \
  github_app_private_key_file=/absolute/path/to/private-key.pem \
  github_webhook_secret_file=/absolute/path/to/github-webhook-secret \
  openai_api_key_file=/absolute/path/to/openai-api-key
```

The helper validates all inputs before creating the Secret with keys `github-app-private-key.pem`,
`github-webhook-secret`, and `openai-api-key`. It fails if the Secret already exists and is separate
from `setup`. The [deployed config](area/apps/lean/abgleich.yaml) uses the mounted files under
`/etc/secrets/abgleich/`.

To replace an existing Secret, first check all three local files, then delete and recreate it:

```bash
kubectl delete secret abgleich-secret --namespace lean && \
make -C area/apps abgleich-secret \
  github_app_private_key_file=/absolute/path/to/private-key.pem \
  github_webhook_secret_file=/absolute/path/to/github-webhook-secret \
  openai_api_key_file=/absolute/path/to/openai-api-key
```

A replacement pod cannot start while the Secret is absent. Recreate it promptly, then restart the
application to pick up the credentials. The [Abgleich README](https://github.com/alexfalkowski/abgleich#readme)
owns GitHub App permissions, webhook setup, commands, cleanup policy, and AI configuration.

### 🧪 Deploy and Check

Run the update through Make so subsequent checks can select the affected apps:

```bash
make area=apps pulumi-preview
make area=apps pulumi-update
make -C area/apps verify
make -C area/apps load
make -C area/apps lint
```

```mermaid
sequenceDiagram
  actor Operator as Operator / CI
  participant Make as Make / release helper
  participant Pulumi
  participant Selection as .release-selection.json
  participant Apps as Deployed apps
  Operator->>Make: pulumi-update (area=apps)
  Make->>Selection: Invalidate previous selection
  Make->>Pulumi: Update apps stack
  Pulumi-->>Make: Resource events and successful summary
  Make->>Selection: Record affected apps and Git revision
  Operator->>Make: verify / load / lint
  Make->>Selection: Validate and read selection
  Make->>Apps: Check selected apps and resources
```

The ignored `area/apps/.release-selection.json` belongs to the last successful update at the current
Git revision. Failed updates and missing, invalid, or stale selections cause checks to fail.

| Update result | Check behavior |
| ------------- | -------------- |
| No relevant resource changes | Skip checks |
| App config or resource changes | Check affected apps |
| Shared or unlabelled resources change | Expand to supported apps and namespace-wide manifest scans |
| App deleted | Exclude it from checks |
| App scaled to zero | Lint it; skip verification and load |
| Deployment created or images changed | Include image scans; config-only updates skip image scans |
| Unknown app or namespace | Fail selection until helper support is added |

Verification uses bounded, fail-fast `curl` requests; Abgleich skips public URL verification in every
verification mode. Load checks use Vegeta and require every request to succeed; Abgleich has no load
scenario and is skipped. Linting uses kube-score and Kubescape against the live cluster. Install these
tools and `jq` before running the checks. Load results are ignored files under `area/apps/lean/*.bin`.

To check an existing deployment independently of the last update, use the explicit `*-all` targets:

```bash
make -C area/apps verify-all
make -C area/apps load-all
make -C area/apps lint-all
```

`lint-all` scans namespace manifests and every supported app's images, including unchanged images.
Per-app commands are available through the helper, for example `cd area/apps && ./release verify-web`.

The manual `make -C area/apps rollout` uses the Git release diff: it restarts version-only changed apps,
otherwise all supported apps. Set `APPS_RELEASE_BASE` and `APPS_RELEASE_HEAD` to override the default
comparison against the current commit's parent (or `origin/master` when unavailable). `rollout-all`
always restarts all supported apps. This Git comparison also supports CI's release-only build skip;
it does not select post-update checks. `APPS_RELEASE_DRY_RUN=true` prints helper actions without
contacting infrastructure; selective checks still need a valid selection.

> [!CAUTION]
> `make -C area/apps delete` deletes the entire `lean` namespace.

## ⚓ Cluster Add-ons

[area/k8s](area/k8s/Makefile) installs nginx ingress, the CircleCI release agent, and metrics-server
using Helm and kubectl. Run this manually from an operator workstation after creating the cluster.
With `doctl` installed, save kubeconfig if needed and confirm the intended context before setup:

```bash
make -C area/k8s save-config
make -C area/k8s setup
```

The default setup requires `CIRCLECI_K8S_TOKEN` for the release agent and passes it directly to Helm.
Use `make -C area/k8s pods` to inspect cluster pods.

> [!CAUTION]
> `make -C area/k8s delete` deletes the `nginx-ingress`, `circleci`, and `metrics-server` namespaces.

## 🧾 Configuration Tools

The repo includes two maintainer helpers that rewrite HJSON in place:

| Helper | Build | Example |
| ------ | ----- | ------- |
| Normalize an area's config | `make build-format` | `./format -k cf` |
| Update an app version | `make build-bump` | `./bump -n bezeichner -v 1.2.3` |

Both accept `-p <path>` to override the default config file. Format kinds are `apps|cf|do|gh` and default
to `area/<kind>/<kind>.hjson`; mismatched area fields are rejected without rewriting the file. Bump
defaults to `area/apps/apps.hjson` and expects a semantic version from automation. Use each binary's
`-h` output for its flags.

## 🛠️ Provider Notes

DigitalOcean's cluster version is pinned in [internal/do/do.go](internal/do/do.go). Apply patch updates
in code; initiate minor/major upgrades in the DigitalOcean UI, then align the code.

GitHub repositories require at least one `checks` entry for strict branch protection on `master`.
Required checks are the merge gate; zero approving reviews is intentional. Pages and collaborators
may need to be disabled for initial repository creation and enabled in a follow-up change. See the
[repository schema](api/infraops/v2/service.proto) and [implementation](internal/gh/repository.go)
for fixed settings and labels.
