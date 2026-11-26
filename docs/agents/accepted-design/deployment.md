# Accepted design for deployment

Part of the [accepted-design index](../accepted-design.md).

## Application release handoff

Each application's own CircleCI pipeline opens pull requests against
`area/apps/apps.hjson` titled `release(<app>): upgraded to vX.Y.Z`. The app
pipeline uses the shared `alexfalkowski/release` Docker image. Its Dockerfile
pins `cmd/bump` via `install-go-tool`, so the pipeline uses the version packaged
in that image. When diagnosing a release, check the app's `.circleci/config.yml`,
`../docker/release/Dockerfile`, and `../docker/release/deploy` alongside the
InfraOps update workflow.

[Application deployment instructions](../../../README.md#-application-deployment)

## CircleCI validation and updates

The setup workflow uses path filtering to enable Pulumi area jobs. It also
selects when `build` runs `make build-bump` and `make build-format`. The `build`
job runs `make lint`, `make scripts-lint`, `make api-lint`, and `make sec`.
The `specs` job runs `make specs`, `make coverage`, and `make codecov-upload`.
Both jobs can skip app version-only changes.

Area previews run on non-`master` branches. Updates run on `master` and are
serialized per area. The apps update also runs `make -C area/apps verify`,
`load`, and `lint` after applying the change.

Non-`master` workflows run `sync` after `build`, `specs`, and applicable previews.
The `master` workflow runs the `version` release job after `build` and `specs`.
`wait-all` is the final aggregate gate for both paths.

[CircleCI setup](../../../.circleci/config.yml) ·
[CircleCI jobs and workflow](../../../.circleci/continue_config.yml)

## Manual Kubernetes add-ons

`area/k8s/Makefile` manages Kubernetes add-ons from an operator workstation.
Its setup passes `CIRCLECI_K8S_TOKEN` directly to Helm for the CircleCI release
agent. CircleCI area workflows operate the Pulumi stacks.

[Add-ons Makefile](../../../area/k8s/Makefile) ·
[Cluster add-ons instructions](../../../README.md#-cluster-add-ons)
