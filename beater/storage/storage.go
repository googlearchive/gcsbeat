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
