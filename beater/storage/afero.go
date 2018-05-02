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

	"github.com/spf13/afero"
)

func newAferoBucketProvider(bucket string) StorageProvider {
	// strip the file:// prefix
	basePath := bucket[7:]
	fs := afero.NewBasePathFs(afero.NewOsFs(), basePath)
	return newAferoStorageProvider(fs)
}

func newAferoStorageProvider(fs afero.Fs) StorageProvider {
	return &aferoStorageProvider{fs, make(map[string]bool)}
}

// aferoStorageProvider implements StorageProvider using an afero FS
// it can be useful for testing locally or unit-testing with in-memory filesystems.
type aferoStorageProvider struct {
	fs        afero.Fs
	processed map[string]bool
}

func (asp *aferoStorageProvider) ListUnprocessed() ([]string, error) {
	files, err := afero.ReadDir(asp.fs, ".")
	if err != nil {
		return nil, err
	}

	var out []string
	for _, f := range files {
		out = append(out, f.Name())
	}

	explainFoundFiles(asp.fs.Name(), out)

	return FilterAndExplain("exists in processed cache", out, InvertFilter(asp.WasProcessed))
}

func (asp *aferoStorageProvider) Read(path string) (io.ReadCloser, error) {
	return asp.fs.Open(path)
}

func (asp *aferoStorageProvider) Remove(path string) error {
	asp.processed[path] = false
	return asp.fs.Remove(path)
}

func (asp *aferoStorageProvider) WasProcessed(path string) (bool, error) {
	_, ok := asp.processed[path]
	return ok, nil
}

func (asp *aferoStorageProvider) MarkProcessed(path string) error {
	asp.processed[path] = true
	return nil
}
