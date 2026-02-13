// Package test provides Pulumi mocks used by unit tests.
//
// The mocks in this package implement pulumi.Mocks and are intended to be passed to
// pulumi.RunErr(..., pulumi.WithMocks(...)) when testing Pulumi programs without
// making real provider calls.
package test
