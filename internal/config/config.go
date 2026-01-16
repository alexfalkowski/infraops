package config

import (
	"os"

	"github.com/hjson/hjson-go/v4"
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

	return hjson.Unmarshal(bytes, config)
}

// Write the configuration to the path, unless an error occurs.
func Write[T Config](path string, config T) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	out, err := hjson.Marshal(config)
	if err != nil {
		return err
	}

	out = append(out, "\n"...)
	return os.WriteFile(path, out, info.Mode())
}
