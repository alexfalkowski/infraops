// Package app defines the Kubernetes application model and Pulumi resource creation helpers.
//
// It is used by the `area/apps` Pulumi program to:
//   - load application configuration (HJSON decoded into protobuf messages),
//   - convert protobuf messages into internal Go types, and
//   - create Kubernetes resources (Deployment, Service, Ingress, etc.) for each application.
//
// The functions in this package are intended to be deterministic given the same input
// configuration and are safe to use with Pulumi mocks in unit tests.
package app
