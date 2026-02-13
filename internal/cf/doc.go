// Package cf provides helpers for managing Cloudflare resources via Pulumi.
//
// It is used by the `area/cf` Pulumi program to load area configuration (HJSON decoded into
// protobuf messages) and to create/update Cloudflare resources.
//
// Environment variables:
//
//   - CLOUDFLARE_ACCOUNT_ID: required; used as the account identifier for account-scoped resources.
package cf
