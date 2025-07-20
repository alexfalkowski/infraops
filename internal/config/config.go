package config

import (
	"os"

	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

// Read the path and return the configuration.
func Read(path string, config any) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return prototext.Unmarshal(bytes, config.(proto.Message))
}
