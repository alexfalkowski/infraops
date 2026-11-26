# Accepted design for configuration

Part of the [accepted-design index](../accepted-design.md).

## Helper commands and test scope

`cmd/bump` is an internal automation helper. Automation supplies a semantic
version through `-v`, which the helper writes as supplied.

Other testing environments verify the CLI behavior of `cmd/bump` and
`cmd/format`. Existing in-repo tests cover version updates and rejection of a
mismatched format kind. Limit test-gap findings to behavior with a concrete
coverage gap beyond the absence of an in-repo command-level suite.

[Bump command](../../../cmd/bump/bump.go) ·
[Version tests](../../../internal/app/version/version_test.go) ·
[Format tests](../../../cmd/format/format_test.go)

## HJSON writes

`internal/config.Write` preserves the existing file mode and appends a trailing
newline. Standalone tests for those two narrow behaviors are outside the
accepted test scope.

[HJSON writer](../../../internal/config/config.go)

## Protobuf validation

CircleCI keeps the `make api-breaking` step commented out for the intentional
Application field renumbering and named storage schema change. Retain the
step for re-enabling against the updated `master` baseline. `make api-breaking`
remains available for local comparison, and `make api-lint` runs in CI.

Schema changes require `make api-generate` and the resulting `service.pb.go` in
the same change. Treat the absence of a CI regeneration-and-diff step as a
generated-code freshness gap only when there is additional evidence.

[Configuration schema](../../../api/infraops/v2/service.proto) ·
[CircleCI validation](../../../.circleci/continue_config.yml)
