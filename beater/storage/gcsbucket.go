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
	"runtime"
	"time"

	"cloud.google.com/go/storage"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/GoogleCloudPlatform/gcsbeat/config"
	"github.com/elastic/beats/libbeat/version"
)

const (
	ProcessedMetadataValue = "processed"
)

func newGcpStorageProvider(cfg *config.Config) (StorageProvider, error) {
	bucket := cfg.BucketId

	ctx := context.Background()

	options := []option.ClientOption{
		option.WithCredentialsFile(cfg.JsonKeyFile),
		option.WithUserAgent(GetUserAgent()),
	}

	client, err := storage.NewClient(ctx, options...)

	if err != nil {
		return nil, err
	}

	// TODO make sure we have appropriate permissions on the bucket
	return &gcpStorageProvider{
		ctx:            ctx,
		storageClient:  client,
		bucket:         bucket,
		processedCache: make(map[string]bool),
		metadataKey:    cfg.MetadataKey,
	}, err
}

type gcpStorageProvider struct {
	ctx            context.Context
	storageClient  *storage.Client
	bucket         string
	processedCache map[string]bool
	metadataKey    string
}

func (gsp *gcpStorageProvider) getBucket() *storage.BucketHandle {
	return gsp.storageClient.Bucket(gsp.bucket)
}

func (gsp *gcpStorageProvider) getObject(path string) *storage.ObjectHandle {
	return gsp.getBucket().Object(path)
}

func (gsp *gcpStorageProvider) getAttrs(path string) (*storage.ObjectAttrs, error) {
	return gsp.getObject(path).Attrs(gsp.ctx)
}

func (gsp *gcpStorageProvider) Read(path string) (io.ReadCloser, error) {
	return gsp.getObject(path).NewReader(gsp.ctx)
}

func (gsp *gcpStorageProvider) Remove(path string) error {
	return gsp.getObject(path).Delete(gsp.ctx)
}

func isMarkedAsProcessed(metadata map[string]string, metadataKey string) bool {
	if metadata == nil {
		return false
	}

	value, ok := metadata[metadataKey]

	return ok && value == ProcessedMetadataValue
}

func (gsp *gcpStorageProvider) WasProcessed(path string) (bool, error) {
	attrs, err := gsp.getAttrs(path)
	if err != nil {
		// True in the event of an error to avoid re-processing
		return true, err
	}

	return isMarkedAsProcessed(attrs.Metadata, gsp.metadataKey), nil
}

func (gsp *gcpStorageProvider) MarkProcessed(path string) error {
	attrs, err := gsp.getAttrs(path)
	if err != nil {
		return err
	}

	metadata := attrs.Metadata
	if metadata == nil {
		metadata = make(map[string]string)
	}

	metadata[gsp.metadataKey] = ProcessedMetadataValue

	update := storage.ObjectAttrsToUpdate{
		Metadata: metadata,
	}

	_, err = gsp.getObject(path).Update(gsp.ctx, update)
	return err
}

func (gsp *gcpStorageProvider) ListUnprocessed() ([]string, error) {
	var paths []string

	it := gsp.getBucket().Objects(gsp.ctx, nil)

	for {
		objAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}

		if err != nil {
			return paths, err
		}

		// Eliminate these early rather than polling the network again.
		if isMarkedAsProcessed(objAttrs.Metadata, gsp.metadataKey) {
			continue
		}

		paths = append(paths, objAttrs.Name)
	}

	return paths, nil
}

// GetUserAgent gets a de-facto standardish user agent string.
// It includes, OS, ARCH, build date and git commit hash.
// It uses "Elastic/GCSBeat" as the software identifier.
func GetUserAgent() string {
	os := fmt.Sprintf("(%s; %s)", runtime.GOOS, runtime.GOARCH)
	ver := fmt.Sprintf("version/%s", version.Commit())
	built := fmt.Sprintf("built/%s", version.BuildTime().Format(time.RFC3339))

	return fmt.Sprintf("Elastic/GCSBeat %s %s %s", os, ver, built)
}
