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
	"fmt"
	"io"

	"github.com/GoogleCloudPlatform/gcsbeat/config"
	"github.com/boltdb/bolt"
)

func newLocalProcessedMiddleware(inner StorageProvider, cfg *config.Config) (StorageProvider, error) {
	return newLocalProcessedMiddlewareBase(inner, cfg.ProcessedDbPath, cfg.MetadataKey)
}

func newLocalProcessedMiddlewareBase(inner StorageProvider, dbPath, metadataKey string) (StorageProvider, error) {
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(metadataKey))
		return err
	})
	if err != nil {
		return nil, err
	}

	return &localProcessedMiddleware{
		wrapped:   inner,
		bucketKey: []byte(metadataKey),
		db:        db,
	}, nil
}

// localProcessedMiddleware keeps a list of files locally synced to the local disk.
// this is useful for platforms where it's not possible to allow writes to a bucket for auditing,
// security, or ownership purposes.
type localProcessedMiddleware struct {
	wrapped   StorageProvider
	bucketKey []byte
	db        *bolt.DB
}

func (middleware *localProcessedMiddleware) ListUnprocessed() ([]string, error) {
	bucketKeys, err := middleware.wrapped.ListUnprocessed()
	if err != nil {
		return nil, err
	}

	message := fmt.Sprintf("exists in processed db %q", middleware.db.Path())
	return FilterAndExplain(message, bucketKeys, InvertFilter(middleware.WasProcessed))
}

func (middleware *localProcessedMiddleware) Read(path string) (io.ReadCloser, error) {
	return middleware.wrapped.Read(path)
}

func (middleware *localProcessedMiddleware) Remove(path string) error {
	// remove if upstream was not an error
	err := middleware.wrapped.Remove(path)

	if err != nil {
		return err
	}

	return middleware.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(middleware.bucketKey)
		return b.Delete([]byte(path))
	})
}

func (middleware *localProcessedMiddleware) WasProcessed(path string) (bool, error) {
	unprocessed, err := middleware.removeProcessed([]string{path})
	return len(unprocessed) == 0, err
}

func (middleware *localProcessedMiddleware) removeProcessed(paths []string) ([]string, error) {
	// open read only txn.
	out := []string{}

	err := middleware.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(middleware.bucketKey)

		for _, path := range paths {
			if bucket.Get([]byte(path)) == nil {
				out = append(out, path)
			}
		}
		return nil
	})

	return out, err
}

func (middleware *localProcessedMiddleware) MarkProcessed(path string) error {
	return middleware.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(middleware.bucketKey)
		return b.Put([]byte(path), []byte(ProcessedMetadataValue))
	})
}
