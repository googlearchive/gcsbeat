// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gobwas/glob"
	"github.com/GoogleCloudPlatform/gcsbeat/beater/codec"
	"github.com/elastic/beats/libbeat/common"
)

type Config struct {
	Interval    time.Duration `config:"interval"`
	BucketId    string        `config:"bucket_id" validate:"required"`
	JsonKeyFile string        `config:"json_key_file"`
	Delete      bool          `config:"delete"`
	Match       string        `config:"file_matches"`
	Exclude     string        `config:"file_exclude"`
	MetadataKey string        `config:"metadata_key"`
	Codec       string        `config:"codec"`
	UnpackGzip  bool          `config:"unpack_gzip"`
	ProcessedDbPath string    `config:"processed_db_path"`
}

var DefaultConfig = Config{
	Interval:    60 * time.Second,
	BucketId:    "",
	JsonKeyFile: "",
	Delete:      false,
	Match:       "*",
	Exclude:     "",
	MetadataKey: "x-goog-meta-gcsbeat",
	Codec:       "text",
	UnpackGzip:  false,
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

	if !codec.IsValidCodec(c.Codec) {
		msg := fmt.Sprintf("%q is an invalid codec. Use one of: %v", c.Codec, codec.ValidCodecs())
		return nil, errors.New(msg)
	}

	return &c, nil
}
