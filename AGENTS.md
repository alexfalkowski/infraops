# AGENTS.md

## Repository overview
This repository is a Go-based infra “monorepo” driven by **Pulumi** programs.

Key concepts:
- Each infrastructure “area” lives under `area/<name>/` and has a Pulumi entrypoint `area/<name>/main.go`.
- Shared implementation lives under `internal/` (e.g. `internal/cf`, `internal/do`, `internal/gh`, `internal/app`).
- Area configuration is stored as **HJSON** (e.g. `area/cf/cf.hjson`) and maps to protobuf types defined in `api/infraops/v2/service.proto`.
- Package-level GoDocs are centralized in per-package `doc.go` files (not scattered across arbitrary `*.go` files).

## Quick orientation (where things live)
- `area/`
  - `apps/`, `cf/`, `do/`, `gh/`: Pulumi programs (each has `Pulumi.yaml`, `*.hjson`, `main.go`, `main_test.go`).
  - `k8s/`: cluster add-ons installed via `helm`/`kubectl` (Makefile-driven).
- `internal/`
  - `app/`: Kubernetes application deployment logic (Pulumi Kubernetes resources).
  - `cf/`: Cloudflare provisioning logic (Pulumi Cloudflare provider).
  - `do/`: DigitalOcean provisioning logic (Pulumi DigitalOcean provider).
  - `gh/`: GitHub provisioning logic (Pulumi GitHub provider).
  - `config/`: HJSON read/write helpers for protobuf messages (`internal/config/config.go`).
  - `inputs/`: shared Pulumi input constants (`internal/inputs/inputs.go`).
  - `log/`: minimal slog logger used by CLI tools (`internal/log/log.go`).
  - `test/`: Pulumi mocks used in unit tests (`internal/test/stub.go`).
- `api/`: protobuf definitions + Buf config; Go codegen output is under `api/infraops/v2/`.
- `cmd/`
  - `bump/`: CLI to bump an app version in `area/apps/apps.hjson`.
  - `format/`: CLI to normalize/format HJSON config files.
- `bin/`: git submodule containing shared build tooling used by the top-level `Makefile`.

## Rule / instruction files
No agent-specific rule files were found (e.g. `.cursor/rules`, `.cursorrules`, `CLAUDE.md`, `.github/copilot-instructions.md`).

## Essential commands (observed)
### Go deps, lint, tests, security
From the repo root:

```bash
make dep          # go mod download/tidy/vendor (see bin/build/make/go.mak)
make lint         # field-alignment + golangci-lint (uses bin submodule tooling)
make specs        # gotestsum + go test (writes junit/coverage under test/reports)
make coverage     # generates HTML + func coverage in test/reports
make sec          # govulncheck -test ./...
```

Notes:
- The `make specs` target runs tests with `-mod vendor` (see `bin/build/make/go.mak`). In CI, `make dep` runs before `make specs`.
- Lint configuration is in `.golangci.yml`.

### API (protobuf / Buf)
From the repo root:

```bash
make api-lint
make api-breaking
make api-generate
```

Or directly:

```bash
make -C api lint
make -C api breaking
make -C api generate
```

Proto sources: `api/infraops/v2/service.proto`.
Generated Go output: `api/infraops/v2/service.pb.go` (regenerate via `make api-generate` rather than editing).

Notes:
- Protobuf message/field comments are treated as part of the configuration contract and are kept detailed/implementation-accurate.
- After modifying `service.proto`, regenerate code with `make api-generate` (or `make -C api generate`) and run tests.

### Pulumi (preview/update per area)
From the repo root:

```bash
make pulumi-login
make area=<apps|cf|do|gh> pulumi-preview
make area=<apps|cf|do|gh> pulumi-update
make area=<apps|cf|do|gh> pulumi-cancel
make area=<apps|cf|do|gh> pulumi-delete
```

These targets run Pulumi with:
- stack: `alexfalkowski/<area>/prod`
- cwd: `area/<area>`

(See `Makefile:10-27`.)

### Kubernetes add-ons (helm/kubectl)
From `area/k8s/Makefile`:

```bash
make -C area/k8s save-config
make -C area/k8s setup
make -C area/k8s delete
make -C area/k8s pods
```

This Makefile references environment variables:
- `CIRCLECI_K8S_TOKEN` (used when installing the CircleCI release agent)
- `BETTER_STACK_COLLECTOR_SECRET` (used when installing Better Stack)

### Applications area helpers
From `area/apps/Makefile`:

```bash
make -C area/apps save-config
make -C area/apps setup
make -C area/apps delete
make -C area/apps rollout
make -C area/apps verify
make -C area/apps load
make -C area/apps lint   # runs kube-score + kubescape
```

## Configuration format and flow
- Config files are HJSON (e.g. `area/cf/cf.hjson`, `area/do/do.hjson`, `area/gh/gh.hjson`, `area/apps/apps.hjson`).
- The schema is defined in `api/infraops/v2/service.proto`.
- Config loading/writing is done via `internal/config.Read` / `internal/config.Write` (`internal/config/config.go`).
  - `Read` uses `hjson.Unmarshal` into a protobuf message.
  - `Write` preserves file mode (`os.Stat(...).Mode()`), marshals via `hjson.Marshal`, and appends a trailing newline.

### Config formatting tool
Build and run:

```bash
make build-format
./format -k <apps|cf|do|gh>
```

- `-p` can override the path; default is `area/<kind>/<kind>.hjson` (see `cmd/format/format.go`).

### App version bump tool
Build and run:

```bash
make build-bump
./bump -n <appName> -v <version>
```

- `-p` can override the path; default is `area/apps/apps.hjson` (see `cmd/bump/bump.go`).
- Implementation updates `config.GetApplications()[i].Version` (see `internal/app/version/version.go`).

## Code patterns and conventions
### Pulumi entrypoints
Each area entrypoint follows the same pattern:
- `pulumi.Run(func(ctx *pulumi.Context) error { ... })`
- read `*.hjson` via `internal/<area>.ReadConfiguration`
- convert protobuf types into internal types via `Convert*`
- create resources via `Create*`

Examples:
- `area/apps/main.go`
- `area/cf/main.go`
- `area/do/main.go`
- `area/gh/main.go`

### “Convert then Create”
A typical flow is:
- `ConvertX(*v2.X) *X`
- `CreateX(ctx, *X) error`

See:
- `internal/gh/gh.go` (`ConvertRepository`, `CreateRepository`)
- `internal/do/do.go` (`ConvertCluster`, `CreateCluster`)
- `internal/app/app.go` (`ConvertApplication`, `CreateApplication`)

### Logging
The CLI tools (`cmd/format`, `cmd/bump`) use `internal/log.NewLogger()` which returns a `slog` text logger writing to stdout.

### Testing
- Uses `testify/require`.
- Pulumi programs are tested with `pulumi.RunErr(..., pulumi.WithMocks(...))`.
- Mocks are provided by `internal/test.Stub` and `internal/test.ErrStub`.

See `area/*/main_test.go` and `internal/test/stub.go`.

## Linting / formatting expectations
- `.editorconfig` sets:
  - Go: tabs, size 4
  - Makefiles: tabs
  - General text: spaces, size 2
- `.golangci.yml` enables many linters and formatters; line length is configured via `lll` at 130.

## CI (CircleCI) behavior (observed)
- `.circleci/config.yml` is a setup config using `circleci/path-filtering` to set pipeline parameters per area.
- `.circleci/continue_config.yml` runs:
  - `build` (always): `make lint`, `make api-lint`, `make api-breaking`, `make sec`, `make specs`, `make coverage`, `make codecov-upload`.
  - per-area workflows:
    - `*_preview` runs on non-`master` branches (Pulumi preview).
    - `*_update` runs only on `master` (Pulumi update), plus extra steps for `apps` (`verify`, `load`, `lint`).

If reproducing locally, mirror the same Make targets used in CI.

## Gotchas / non-obvious behavior
- **Cloudflare requires `CLOUDFLARE_ACCOUNT_ID`**: `internal/cf/cf.go` reads it from the environment (`os.Getenv`) and stores it as a package-level Pulumi string.
- **Apps config maps read files relative to the Pulumi working directory**:
  - `internal/app/config.go` reads `<namespace>/<app>.yaml` using `os.Getwd()` + `filepath.Join(wd, ns, file)`.
  - When running Pulumi via `make area=apps pulumi-*`, Pulumi runs with `--cwd area/apps`, so ensure the expected files exist under `area/apps/<namespace>/`.
- **Apps env var secret references**:
  - `EnvVar.value` supports the convention `secret:<secretName>/<key>`.
  - At deploy time this is converted into a Kubernetes `SecretKeyRef` with Secret name `<secretName>-secret` and key `<key>` (see `internal/app/environment.go` and `internal/app/secret.go`).
- **`Application.secrets` vs secret env vars**:
  - `Application.secrets` is an application-level dependency list used by the deployment implementation to provision and/or wire Secret resources (for example volumes/attachments).
  - Secret references in `env_vars` (`secret:<secretName>/<key>`) target specific keys; the `secrets` list does not specify keys.
  - These behaviors are documented in `api/infraops/v2/service.proto` and should be kept in sync with implementation.
- **Apps resource sizing**:
  - `Application.resource` selects a resource profile; unknown values fall back to `"small"` (see `internal/app/resource.go`).
  - Current mapping includes `"small"`: cpu 125m-250m, memory 64Mi-128Mi, ephemeral-storage 1Gi-2Gi.
- **GitHub repository creation has a 2-step process for some features** (from `README.md`):
  - Pages and collaborators may need to be disabled on first creation and enabled in a follow-up change to avoid timing issues.
