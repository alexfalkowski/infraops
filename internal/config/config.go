package config

import (
	"fmt"
	"os"
	"strings"

	v2 "github.com/alexfalkowski/infraops/v2/api/infraops/v2"
	"github.com/hjson/hjson-go/v4"
	"google.golang.org/protobuf/proto"
)

// Config is an alias for proto.Message and represents a protobuf configuration message.
type Config = proto.Message

// Read reads the HJSON file at path and unmarshals it into config, rejecting duplicate object keys.
//
// The caller is responsible for passing a concrete protobuf message pointer.
// Read returns I/O and decode errors. Duplicate-key diagnostics omit configuration values.
func Read[T Config](path string, config T) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	err = decode(bytes, config)
	// The HJSON decoder reports duplicate keys with both values and a source excerpt.
	if err != nil && strings.HasPrefix(err.Error(), "Found duplicate values (") {
		return fmt.Errorf("%s: duplicate HJSON object key", path)
	}
	return err
}

func decode(data []byte, config Config) error {
	options := hjson.DefaultDecoderOptions()
	options.DisallowDuplicateKeys = true
	if _, ok := config.(*v2.Kubernetes); ok {
		if err := validateStorageFields(data, options); err != nil {
			return err
		}
	}
	return hjson.UnmarshalWithOptions(data, config, options)
}

func validateStorageFields(data []byte, options hjson.DecoderOptions) error {
	// Read application keys before decoding into the case-insensitive struct fields.
	var document struct {
		Applications []map[string]any `json:"applications"`
	}
	if err := hjson.UnmarshalWithOptions(data, &document, options); err != nil {
		return err
	}
	for index, application := range document.Applications {
		seen := false
		for key := range application {
			if !strings.EqualFold(key, "storage") {
				continue
			}
			if seen {
				return fmt.Errorf(`applications[%d]: duplicate HJSON field "storage"`, index)
			}
			seen = true
		}
	}
	return nil
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
