![Gopher](assets/gopher.png)
[![DigitalOcean Referral Badge](https://web-platforms.sfo2.cdn.digitaloceanspaces.com/WWW/Badge%202.svg)](https://www.digitalocean.com/?refcode=b80d3d6467e1&utm_campaign=Referral_Invite&utm_medium=Referral_Program&utm_source=badge)

[![CircleCI](https://circleci.com/gh/alexfalkowski/infraops.svg?style=shield)](https://circleci.com/gh/alexfalkowski/infraops)
[![codecov](https://codecov.io/gh/alexfalkowski/infraops/graph/badge.svg?token=U3X5JGAA8I)](https://codecov.io/gh/alexfalkowski/infraops)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexfalkowski/infraops/v2)](https://goreportcard.com/report/github.com/alexfalkowski/infraops/v2)
[![Go Reference](https://pkg.go.dev/badge/github.com/alexfalkowski/infraops/v2.svg)](https://pkg.go.dev/github.com/alexfalkowski/infraops/v2)
[![Stability: Active](https://masterminds.github.io/stability/active.svg)](https://masterminds.github.io/stability/active.html)

# 🧭 InfraOps

A Go-based infrastructure “monorepo” powered by **Pulumi**, with configuration stored as **HJSON** and decoded into a **protobuf schema**.

> [!NOTE]
> This repository is the operational source for several personal infrastructure areas. Prefer previews and narrow area changes before applying updates.

## 🗺️ Overview

This repository manages multiple infrastructure “areas”:

- `area/apps`: Kubernetes applications deployed to a cluster.
- `area/cf`: Cloudflare resources (zones, DNS records, R2 buckets/custom domains).
- `area/do`: DigitalOcean resources (VPC + Kubernetes clusters).
- `area/gh`: GitHub repositories and settings (branch protection, Pages, collaborators).
- `area/k8s`: Cluster add-ons installed via `helm`/`kubectl` (Makefile driven, not Pulumi).

Each Pulumi area has:

- a Pulumi entrypoint at `area/<name>/main.go`
- a config file at `area/<name>/<name>.hjson`
- a Pulumi project file `area/<name>/Pulumi.yaml`

Shared implementation lives under `internal/` (e.g. `internal/app`, `internal/cf`, `internal/do`, `internal/gh`).

## 🧰 Tooling

### 🧱 Development tools

Development and CI-style commands use:

- Go (version from `go.mod`: `1.25.8`)
- `make`
- Ruby, used by shared lint helpers under `bin/`
- Pulumi CLI
- `buf` for protobuf linting, breaking-change checks, and code generation
- `fieldalignment` for `make lint`
- `golangci-lint` for `make lint` when installed in `PATH`
- `gotestsum` for `make specs`
- `govulncheck` and Trivy for `make sec`

### 🛠️ Operator tools

Infrastructure operation also uses:

- Pulumi: <https://www.pulumi.com/>
- kubectl: <https://kubernetes.io/docs/reference/kubectl/>
- Helm: <https://helm.sh/>
- kube-score: <https://kube-score.com/>
- doctl: <https://docs.digitalocean.com/reference/doctl/>
- kubescape: <https://kubescape.io/>
- curl: <https://curl.se/>
- vegeta: <https://github.com/tsenart/vegeta>

Operator-only helper targets use these tools as needed:

- `make -C area/apps save-config` and `make -C area/k8s save-config`: `doctl`
- `make -C area/apps setup|delete|rollout`: `kubectl`
- `make -C area/apps lint`: `kube-score`, `kubescape`, `kubectl`
- `make -C area/apps verify`: `curl`
- `make -C area/apps load`: `vegeta`
- `make -C area/k8s setup|delete|pods`: `helm`, `kubectl`

## 🧾 Configuration (HJSON + Protobuf Schema)

All area configuration files are **HJSON** (`*.hjson`). They are decoded into protobuf messages defined in:

- `api/infraops/v2/service.proto`

Those protobuf messages are then converted into internal Go types and used to provision resources.

The protobuf types provide the expected configuration shape. Semantic checks that are not modeled in the schema may still fail later during Pulumi or provider operations.

### 🧹 Format and Normalize Config

This repo includes a small CLI to normalize/format config files by:

1. decoding the HJSON into the appropriate protobuf message
2. writing it back out in a canonical form

Build:

```bash
make build-format
```

Format a config kind (`apps|cf|do|gh`), using the default location `area/<kind>/<kind>.hjson`:

```bash
./format -k cf
```

Override the path:

```bash
./format -k cf -p area/cf/cf.hjson
```

### 📝 Config Schema Notes

A few conventions are implemented by the Go code and are worth knowing when editing HJSON:

#### 🔐 `EnvVar.value` Secret References (apps)

Environment variables support literal values, and a secret reference format:

- `secret:<secretName>/<key>`

At deploy time, this becomes a Kubernetes `SecretKeyRef`:

- Secret name: `<secretName>-secret`
- Secret key: `<key>`

The format is intentionally not pre-validated by the helper code; malformed references are expected
to fail during Pulumi/Kubernetes application.

Example:

```hjson
env_vars: [
  { name: "DATABASE_URL", value: "secret:db/url" }
]
```

#### 🧩 `Application.secrets` vs Secret Env Vars (apps)

- `Application.secrets` is an **application-level dependency list** used by the deployment implementation to wire existing Kubernetes Secrets as volumes.
- Secret references in `env_vars` (the `secret:<secretName>/<key>` format) reference **specific keys** in those secrets.
- They often use the same `<secretName>` values, but they serve different purposes.
- The app program does not create Secret objects or define Secret keys. For internal apps, each listed secret is expected to exist as `<secretName>-secret` and is mounted at `/etc/secrets/<secretName>`.

#### 📏 `Application.resource` Sizing (apps)

`Application.resource` selects a resource profile. Current mapping:

- `"small"` (default): cpu `125m-250m`, memory `64Mi-128Mi`, ephemeral-storage `1Gi-2Gi`

Unknown values fall back to `"small"`.

#### 🛡️ App NetworkPolicy Baseline

The apps Pulumi program creates a `NetworkPolicy` that selects each app's pods. Ingress is limited
to the ports exposed by the app kind: external apps expose HTTP on `8080`, while internal apps expose
debug on `6060`, HTTP on `8080`, and gRPC on `9090`. Egress currently remains open because outbound
traffic flows are not modeled per app yet. Future egress restrictions should be introduced per
namespace/app after the required flows are known.

## 🔁 Common Workflows

### 🚚 Checkout / Bootstrap

The root `Makefile` includes shared build tooling from the `bin` submodule. Initialize it before running Make targets from a fresh checkout:

```bash
git submodule sync
git submodule update --init
```

### ✅ Dependencies, Linting, Tests, and Security

> [!IMPORTANT]
> Run `make dep` before local validation after checkout, dependency changes, or generated/vendor-state changes.

From the repository root:

```bash
make dep          # download/tidy/vendor deps
make lint         # lint (including field alignment)
make sec          # govulncheck + Trivy
make specs        # gotestsum + go test (junit/coverage under test/reports)
make coverage     # HTML + function coverage under test/reports
```

### 🧬 Protobuf / API

> [!IMPORTANT]
> Do not edit generated Go code under `api/infraops/v2/*.pb.go` directly.

Instead:

```bash
make api-lint
make api-breaking
make api-generate
```

(Or: `make -C api lint|breaking|generate`.)

## 🚀 Pulumi: Preview/Update Per Area

Pulumi is typically run via Makefile targets from the repo root.

Login:

```bash
make pulumi-login
```

### 🔑 Provider Credentials

Pulumi also needs the provider credentials for the area being previewed or updated:

| Area | Required local access |
| ---- | --------------------- |
| `apps` | Kubernetes access through the current kubeconfig context |
| `cf` | `CLOUDFLARE_API_TOKEN` and `CLOUDFLARE_ACCOUNT_ID` |
| `do` | `DIGITALOCEAN_TOKEN` |
| `gh` | `GITHUB_TOKEN` |

Preview/update:

```bash
make area=cf pulumi-preview
make area=cf pulumi-update
```

Supported areas for these targets:

- `apps`, `cf`, `do`, `gh`

The Makefile runs Pulumi with:

- stack: `alexfalkowski/<area>/prod`
- working directory: `area/<area>`

That working directory matters because the programs read `<area>.hjson` via a relative path.

Other stack operations:

```bash
make area=cf pulumi-refresh
make area=cf pulumi-cancel
make area=cf pulumi-delete
```

> [!WARNING]
> `pulumi-update`, `pulumi-refresh`, and `pulumi-cancel` affect remote infrastructure or stack state. Run a preview first unless you are recovering from a known failed operation.

> [!CAUTION]
> `pulumi-delete` runs `pulumi stack rm --force`; it removes Pulumi stack state and can orphan managed resources if they still exist.

## 🧭 Areas

### 📦 Applications (`area/apps`)

Deploys Kubernetes applications described in `area/apps/apps.hjson`.

#### ⚙️ Configure

See:

- `area/apps/apps.hjson`

This file uses the `Kubernetes` message in `api/infraops/v2/service.proto`.

Internal apps also need an application config file at:

- `area/apps/<namespace>/<app>.yaml`

For example, the `bezeichner` app in namespace `lean` reads `area/apps/lean/bezeichner.yaml`.
The apps Pulumi program turns that file into a Kubernetes ConfigMap entry named `<app>.yaml`,
then mounts it into the container at `/etc/<app>/<app>.yaml`.

#### 🏗️ Install / Setup

Save the cluster kubeconfig if your local context is not already configured:

```bash
make -C area/apps save-config
```

Prepare the local app namespace/helper resources:

```bash
make -C area/apps setup
```

Apply the Pulumi resources described by `apps.hjson`:

```bash
make area=apps pulumi-update
```

#### 🗑️ Delete

> [!CAUTION]
> This deletes the `lean` namespace through `kubectl`.

```bash
make -C area/apps delete
```

#### ⬆️ Update an Application Version (bump tool)

Build:

```bash
make build-bump
```

Update a single app version in config. This is an internal helper used by automation, and the
`-v` value is expected to be a semantic version:

```bash
./bump -n bezeichner -v 1.559.0
```

By default it edits `area/apps/apps.hjson`. Override path:

```bash
./bump -n bezeichner -v 1.559.0 -p area/apps/apps.hjson
```

> [!TIP]
> Run `./format -k apps` after edits to keep config normalized.

### ☁️ Cloudflare (`area/cf`)

Manages Cloudflare resources using Pulumi’s Cloudflare provider:

- <https://www.pulumi.com/registry/packages/cloudflare/>

Config:

- `area/cf/cf.hjson`

#### 🔑 Required Environment Variables

Cloudflare Pulumi runs require:

- `CLOUDFLARE_API_TOKEN` for provider authentication
- `CLOUDFLARE_ACCOUNT_ID` for account-scoped resources like R2 buckets

### 🌊 DigitalOcean (`area/do`)

Manages DigitalOcean resources using Pulumi’s DigitalOcean provider:

- <https://www.pulumi.com/registry/packages/digitalocean/>

Config:

- `area/do/do.hjson`

#### 📐 Cluster Defaults

The DigitalOcean program creates the cluster VPC and then provisions the cluster with fixed
operational defaults:

- Region: `fra1`
- Kubernetes version: pinned in code
- Node count: `2`
- Maintenance window: any day at `23:00`
- Associated resources: destroyed with the cluster

Cluster `resource` values map to node capacity:

- `"small"` (default): 2 vCPU / 4 GB node
- `"medium"`: 4 vCPU / 8 GB node

Unknown or empty values fall back to `"small"`.

#### 🧱 Manual Prerequisites (DigitalOcean UI)

Some items may be created manually depending on account setup:

- A default project (example):

| Name          | Description                           |
| ------------- | ------------------------------------- |
| lean-thoughts | All experiments for lean-thoughts. |

The Pulumi program creates the VPC used by the cluster and attaches the cluster to it.

#### ⬆️ Kubernetes Cluster Upgrades

Cluster version is pinned in code:

- `internal/do/do.go`

Guidance:

- Patch versions can be updated in code.
- Minor/major upgrades should be initiated via the DigitalOcean UI (per DO guidance), then aligned in code.

### 🐙 GitHub (`area/gh`)

Manages GitHub resources using Pulumi’s GitHub provider:

- <https://www.pulumi.com/registry/packages/github/>

Config:

- `area/gh/gh.hjson`

This area was inspired by:

- <https://github.com/dirien/pulumi-github>

#### 🪜 Repository Creation Caveat (2-step enablement)

Some repository features may require a two-step approach: create the repository first, then enable features in a follow-up change. This avoids timing issues around initial default branch creation.

##### 📄 GitHub Pages

First change: disable Pages (or omit pages config):

```hjson
pages: { enabled: false }
```

Second change: enable Pages:

```hjson
pages: { enabled: true }
```

Optional CNAME:

```hjson
pages: {
  enabled: true
  cname: www.yoursite.com
}
```

##### 👥 Collaborators

When enabled, collaborator management grants `admin` permission to `lean-thoughts-ci`
on `alexfalkowski/<repository>`.

First change:

```hjson
collaborators: { enabled: false }
```

Second change:

```hjson
collaborators: { enabled: true }
```

If the pipeline fails due to timing, a rerun often succeeds.

### ⚓ Kubernetes Add-ons (`area/k8s`)

This is not a Pulumi area. It contains cluster add-ons installed via `helm`/`kubectl`.
These targets are intended to be run manually from an operator workstation, not from CI.

> [!CAUTION]
> Run this only after you have a Kubernetes cluster, a kubeconfig context for that cluster, and the required local secrets.

Save the cluster kubeconfig if your local context is not already configured:

```bash
make -C area/k8s save-config
```

Setup:

```bash
make -C area/k8s setup
```

Delete:

> [!CAUTION]
> This deletes the `nginx-ingress`, `circleci`, and `metrics-server` namespaces through `kubectl`.

```bash
make -C area/k8s delete
```

Useful debugging:

```bash
make -C area/k8s pods
```

#### 🔑 Required Environment Variables

The default `make -C area/k8s setup` path installs the CircleCI release agent along with nginx ingress and metrics-server.

> [!IMPORTANT]
> `CIRCLECI_K8S_TOKEN` is required when running the default k8s setup target.

- `CIRCLECI_K8S_TOKEN` (CircleCI release agent)

Because the Makefile is a local-operator workflow, the CircleCI release-agent token is passed
directly to Helm rather than through a CI secret-handling path.

## 🗂️ Repository Structure

- `area/`: Pulumi programs and k8s add-ons
- `internal/`: shared implementation (convert + create patterns per area)
- `api/`: protobuf schema and generated code
- `cmd/`: small helper CLIs (`format`, `bump`)
- `bin/`: shared build tooling (git submodule)
