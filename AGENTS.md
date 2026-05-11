# AGENTS.md

## Baseline

- Use `bin/skills/coding-standards` for code changes, fixes, refactors, reviews, tests, linting, docs, PR summaries, commits, Makefile changes, CI validation, and verification.
- This repo is a Go/Pulumi infra monorepo. Each Pulumi area lives in `area/<name>/` with `main.go`, `Pulumi.yaml`, and `<name>.hjson`.
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

Config tools:

```bash
make build-format
./format -k <apps|cf|do|gh>

make build-bump
./bump -n <appName> -v <version>
```

## Workflow Rules

- Run `make dep` before validation unless the task is strictly read-only.
- Use repository Make targets over ad hoc commands when they cover the task.
- After editing `api/infraops/v2/service.proto`, run `make api-generate` and include `service.pb.go`.
- Do not hand-edit generated code unless explicitly asked.
- Keep config comments in `service.proto` accurate; they are part of the HJSON configuration contract.
- Tests use `testify/require`; Pulumi tests use `pulumi.RunErr(..., pulumi.WithMocks(...))`.
- Use `internal/test.Stub` and `internal/test.ErrStub` for Pulumi mocks.
- Lint config is `.golangci.yml`; editor defaults are in `.editorconfig`.

## Code Patterns

- Pulumi entrypoints follow: read HJSON, convert protobuf model, create resources.
- Conversion/resource pattern: `ConvertX(*v2.X) *X`, then `CreateX(ctx, *X) error`.
- CLI tools use `internal/log.NewLogger()`.
- `internal/config.Read` uses HJSON unmarshal into protobuf messages.
- `internal/config.Write` preserves destination file mode and appends a trailing newline.

## Intentional Assumptions

- `cmd/bump` is an internal automation helper. `-v` is expected to be a semantic version supplied by automation.
- `area/k8s/Makefile` is for manual operator workstation runs, not CI. `CIRCLECI_K8S_TOKEN` is passed directly to Helm in that local workflow.
- GitHub branch protection intentionally requires zero approving PR reviews because this is a solo-maintainer workflow; required status checks are the primary merge gate.
- Apps `NetworkPolicy` resources select app pods, limit ingress to app ports, and intentionally keep egress open until per-app traffic flows are modeled. Do not tighten egress generically without checking DNS, ingress, probes, telemetry, and outbound service calls.
- Apps env var secret references use `secret:<secretName>/<key>` and become Kubernetes `SecretKeyRef`s with Secret name `<secretName>-secret` and key `<key>`. The helper code does not pre-validate the shape; malformed references fail during Pulumi/Kubernetes application.
- `Application.secrets` is an app-level dependency list; secret env vars reference specific keys. Keep this behavior aligned with `service.proto`.
- `Application.resource` falls back to `"small"` for unknown values. Current small profile: cpu `125m-250m`, memory `64Mi-128Mi`, ephemeral-storage `1Gi-2Gi`.
- Cloudflare resources require `CLOUDFLARE_ACCOUNT_ID`; `internal/cf/cf.go` reads it at package scope.
- Apps config maps read `<namespace>/<app>.yaml` relative to the Pulumi working directory. The root Makefile runs Pulumi with `--cwd area/apps`.
- GitHub Pages and collaborators may need to be disabled on first repository creation and enabled in a follow-up change.

## CI

- CircleCI setup config uses path filtering to enable area workflows.
- The main `build` job runs `make lint`, `make api-lint`, `make api-breaking`, `make sec`, `make specs`, `make coverage`, and Codecov upload.
- Area preview jobs run on non-`master`; update jobs run only on `master`.
- Apps updates also run `make -C area/apps verify`, `load`, and `lint`.
