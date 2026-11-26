# AGENTS.md

## Shared guidance

Use `bin/AGENTS.md` for shared skills and cross-repository defaults.

## Accepted design

Read the [accepted-design index](docs/agents/accepted-design.md) and the relevant
area files before proposing a change or recording a finding or audit entry.
Apply each decision to its documented scope and record owner-agreed choices in
the owning area file.

## Baseline

- This repo is a Go/Pulumi infra monorepo.
- Each Pulumi area lives in `area/<name>/` with `main.go`, `Pulumi.yaml`, and `<name>.hjson`.
- Shared implementation lives in `internal/`; config schema lives in `api/infraops/v2/service.proto`; generated Go lives in `api/infraops/v2/service.pb.go`.
- Keep package-level GoDocs in `doc.go` files. Do not scatter package comments across implementation files.

## Layout

- `area/apps`, `area/cf`, `area/do`, `area/gh`: Pulumi programs.
- `area/k8s`: local `helm`/`kubectl` add-ons, not Pulumi and not CI.
- `internal/app`, `internal/cf`, `internal/do`, `internal/gh`: provider/resource logic.
- `internal/config`: HJSON read/write helpers for protobuf messages.
- `internal/inputs`: shared Pulumi inputs.
- `internal/test`: Pulumi mocks.
- `cmd/bump`: internal app version bump helper.
- `cmd/format`: HJSON config formatter.
- `bin`: shared build tooling submodule.

## Related Repositories

- Sibling repositories owned by this org are normally checked out locally one
  directory up, at `../<repo-name>` (for example `../bezeichner`, `../standort`,
  `../web`, `../docker`, `../bin`), and are hosted at
  `https://github.com/alexfalkowski/<repo-name>`. Check for a local sibling
  checkout before treating a related repository as unavailable or reporting a
  cross-repo data gap. If no local checkout exists, use `gh` to read the
  repository on GitHub before reporting the gap:
  `gh repo clone <owner>/<repo> <scratch-dir> -- --depth 1`
  to explore multiple files, or `gh api repos/<owner>/<repo>/contents/<path>`
  to fetch one known file directly. `<owner>` here is `alexfalkowski`, the same
  account this repo's own `origin` remote resolves to
  (`git@github.com:alexfalkowski/infraops.git`); sibling repos share that owner.

## Commands

Run from the repo root unless noted:

```bash
make dep
make lint
make specs
make coverage
make sec
make api-lint
make api-breaking
make api-generate
```

Pulumi:

```bash
make pulumi-login
make area=<apps|cf|do|gh> pulumi-preview
make area=<apps|cf|do|gh> pulumi-update
make area=<apps|cf|do|gh> pulumi-cancel
make area=<apps|cf|do|gh> pulumi-delete
make area=<apps|cf|do|gh> pulumi-refresh
```

Local Kubernetes add-ons:

```bash
make -C area/k8s save-config
make -C area/k8s setup
make -C area/k8s delete
make -C area/k8s pods
```

Application helpers:

```bash
make -C area/apps save-config
make -C area/apps setup
make -C area/apps delete
make -C area/apps rollout
make -C area/apps verify
make -C area/apps load
make -C area/apps lint
```

Operator helper targets require the local CLI tools documented in `README.md`: `doctl`,
`helm`, `kubectl`, `kube-score`, `kubescape`, `curl`, and `vegeta`.

Config tools:

```bash
make build-format
make build-bump
```

The generated helper binaries are maintainer/operator tools; prefer the Make
targets and README-documented workflows before running them directly.

## Workflow Rules

- Run `make dep` before validation unless the task is strictly read-only.
- Use repository Make targets over ad hoc commands when they cover the task.
- After editing `api/infraops/v2/service.proto`, run `make api-generate` and include `service.pb.go`.
- Do not hand-edit generated code unless explicitly asked.
- Keep config comments in `service.proto` accurate; they are part of the HJSON configuration contract.
- Tests use `testify/require`; Pulumi tests use `pulumi.RunErr(..., pulumi.WithMocks(...))`.
- Use `internal/test.Stub` and `internal/test.ResourceStub` for Pulumi mocks.
- Lint config is `.golangci.yml`; editor defaults are in `.editorconfig`.

## Code Patterns

- Pulumi entrypoints follow: read HJSON, convert protobuf model, create resources.
- Conversion/resource pattern: `ConvertX(*v2.X) *X`, then `CreateX(ctx, *X) error`.
- CLI tools use `internal/log.NewLogger()`.
- `internal/config.Read` uses HJSON unmarshal into protobuf messages.
