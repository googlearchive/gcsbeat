// Config is put into a different package to prevent cyclic imports in case
// it is needed in several locations

package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/elastic/beats/libbeat/common"
	"github.com/gobwas/glob"
)

type Config struct {
	Interval    time.Duration `config:"interval"`
	BucketId    string        `config:"bucket_id" validate:"required"`
	JsonKeyFile string        `config:"json_key_file"`
	Delete      bool          `config:"delete"`
	Match       string        `config:"file_matches"`
	Exclude     string        `config:"file_exclude"`
	MetadataKey string        `config:"metadata_key"`

	// TODO add the ability to treat .gz files as gzipped streams

	// TODO add a flag to read stackdriver logs (JSON files with lists of event objects)
	// https://cloud.google.com/logging/docs/export/using_exported_logs
}

var DefaultConfig = Config{
	Interval:    60 * time.Second,
	BucketId:    "",
	JsonKeyFile: "",
	Delete:      false,
	Match:       "*",
	Exclude:     "",
	MetadataKey: "x-goog-meta-gcsbeat",
}

func GetAndValidateConfig(cfg *common.Config) (*Config, error) {
	c := DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("error in config file: %v", err)
	}

	// Preprocessing
	// GCS keys must not have leading or trailing whitespace
	c.MetadataKey = strings.TrimSpace(c.MetadataKey)

	// Validation
	if c.Interval <= 0 {
		return nil, errors.New("Interval must be positive.")
	}

	if _, err := glob.Compile(c.Match); err != nil {
		return nil, errors.New("The matches parameter is not a valid glob.")
	}

	if _, err := glob.Compile(c.Exclude); err != nil {
		return nil, errors.New("The exclude parameter is not a valid glob.")
	}

	if c.MetadataKey == "" {
		return nil, errors.New("The metadata key must not be blank.")
	}

	return &c, nil
}
