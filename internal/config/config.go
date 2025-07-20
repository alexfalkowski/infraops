package config

import (
	"os"

	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

// Read the path and return the configuration, unless an error occurs.
func Read[T proto.Message](path string, config T) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return prototext.Unmarshal(bytes, config)
}

// Write the configuration to the path, unless an error occurs.
func Write[T proto.Message](path string, config T) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	opts := prototext.MarshalOptions{
		Multiline: true,
	}
	bytes, err := opts.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(path, bytes, info.Mode())
}
