// Package inputs defines shared Pulumi input constants used across infrastructure packages.
//
// The values in this package are typed Pulumi inputs (pulumi.String, pulumi.Int, pulumi.Bool, etc.)
// intended to reduce repetitive conversions and keep resource argument construction consistent
// across Pulumi programs.
//
// These constants are small, opinionated helpers; prefer using them when a resource argument
// expects a Pulumi input type and the value is reused throughout the repository.
package inputs
