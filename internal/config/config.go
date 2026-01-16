package config

import (
	"os"

	"go.yaml.in/yaml/v3"
	"google.golang.org/protobuf/proto"
)

// Config is just an alias to proto.Message.
type Config = proto.Message

// Read the path and return the configuration, unless an error occurs.
func Read[T Config](path string, config T) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(bytes, config)
}

// Write the configuration to the path, unless an error occurs.
func Write[T Config](path string, config T) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	out, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(path, out, info.Mode())
}
