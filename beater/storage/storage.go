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

package storage

import (
	"io"
	"strings"

	"github.com/GoogleCloudPlatform/gcsbeat/config"
)

type StorageProvider interface {
	ListUnprocessed() (files []string, err error)
	Read(path string) (reader io.ReadCloser, err error)
	Remove(path string) error
	WasProcessed(path string) (bool, error)
	MarkProcessed(path string) error
}

func NewStorageProvider(cfg *config.Config) (StorageProvider, error) {
	provider, err := newBaseStorageProvider(cfg)

	if err != nil {
		return nil, err
	}

	return wrapWithMiddleware(provider, cfg)
}

func newBaseStorageProvider(cfg *config.Config) (StorageProvider, error) {
	if strings.HasPrefix(cfg.BucketId, "file://") {
		return newAferoBucketProvider(cfg.BucketId), nil
	}

	// connect to GCP
	return newGcpStorageProvider(cfg)
}

func wrapWithMiddleware(provider StorageProvider, cfg *config.Config) (StorageProvider, error) {
	var err error

	if cfg.ProcessedDbPath != "" {
		provider, err = newLocalProcessedMiddleware(provider, cfg)

		if err != nil {
			return nil, err
		}
	}

	return newLoggingStorageProvider(provider), nil
}
