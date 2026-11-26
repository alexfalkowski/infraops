# Accepted design for applications

Part of the [accepted-design index](../accepted-design.md).

## Network policy

Each application `NetworkPolicy` selects the application's pods and limits
ingress to its ports. Egress stays open until per-app traffic flows are modeled.
Before proposing tighter egress, check DNS, ingress, probes, telemetry, and
outbound service calls for the affected apps.

[Network policy implementation](../../../internal/app/network.go)

## Secret references

Environment variable values of the form `secret:<secretName>/<key>` become
Kubernetes `SecretKeyRef`s with Secret name `<secretName>-secret` and key
`<key>`. Shape validation is deferred to Pulumi and Kubernetes.

`Application.secrets` lists app-level dependencies on logical secrets and
supplies Secret volumes where supported. Each environment variable reference
selects its own key. Keep both behaviors aligned with the schema.

[Environment variables](../../../internal/app/environment.go) ·
[Secret volumes](../../../internal/app/secret.go) ·
[Application schema](../../../api/infraops/v2/service.proto)

## Resource profile fallback

`Application.resource` falls back to `small` for unknown values. The profile
sets CPU requests and limits to `125m` and `250m`, memory to `64Mi` and `128Mi`,
and ephemeral storage to `1Gi` and `2Gi`.

[Resource profiles](../../../internal/app/resource.go) ·
[Application schema](../../../api/infraops/v2/service.proto)

## ConfigMap sources

App ConfigMaps read `<namespace>/<app>.yaml` relative to the Pulumi working
directory. The root Makefile previews the apps stack with `--cwd area/apps`.
Updates go through the apps release helper, which uses the same directory.

[ConfigMap implementation](../../../internal/app/config.go) ·
[Root Makefile](../../../Makefile) ·
[Pulumi update dispatcher](../../../area/pulumi-update) ·
[Apps release helper](../../../area/apps/release)
