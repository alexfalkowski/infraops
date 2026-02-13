// Package config provides helpers for reading and writing HJSON configuration files
// into protobuf messages.
//
// These helpers are used by Pulumi programs and CLI tools to keep area configuration
// files normalized while preserving the file mode on disk.
//
// Read unmarshals a HJSON file into the provided protobuf message.
// Write marshals a protobuf message as HJSON, preserves the existing file mode, and
// appends a trailing newline.
package config
