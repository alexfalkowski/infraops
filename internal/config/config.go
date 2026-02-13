package config

import (
	"os"

	"github.com/hjson/hjson-go/v4"
	"google.golang.org/protobuf/proto"
)

// Config is an alias for proto.Message and represents a protobuf configuration message.
type Config = proto.Message

// Read reads the HJSON file at path and unmarshals it into config.
//
// The caller is responsible for passing a concrete protobuf message pointer.
// Read returns any I/O error from reading the file or any decode error from HJSON unmarshalling.
func Read[T Config](path string, config T) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return hjson.Unmarshal(bytes, config)
}

// Write marshals config as HJSON and writes it to path.
//
// Write preserves the existing file mode (as returned by os.Stat) and always appends
// a trailing newline. It returns any I/O error from stat/write or any encode error
// from HJSON marshalling.
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
