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

	"github.com/elastic/beats/libbeat/logp"
)

func newLoggingStorageProvider(inner StorageProvider) StorageProvider {
	return &loggingStorageProvider{
		wrapped: inner,
		logger:  logp.NewLogger("StorageProvider"),
	}
}

type loggingStorageProvider struct {
	wrapped StorageProvider
	logger  *logp.Logger
}

func (lsp *loggingStorageProvider) ListUnprocessed() ([]string, error) {
	lsp.logger.Infof("Fetching file list from server")

	files, err := lsp.wrapped.ListUnprocessed()

	if err != nil {
		lsp.logger.Errorf("Could not fetch list of files from the server: %v", err)
	}

	return files, err
}

func (lsp *loggingStorageProvider) Read(path string) (io.ReadCloser, error) {
	lsp.logger.Infof("Reading file: %q", path)

	file, err := lsp.wrapped.Read(path)

	if err != nil {
		lsp.logger.Errorf("Error reading file: %v", err)
	}

	return file, err
}

func (lsp *loggingStorageProvider) Remove(path string) error {
	lsp.logger.Infof("Deleting file: %q", path)

	err := lsp.wrapped.Remove(path)

	if err != nil {
		lsp.logger.Errorf("Error deleting file %q: %v", path, err)
	}

	return err
}

func (lsp *loggingStorageProvider) WasProcessed(path string) (bool, error) {
	lsp.logger.Infof("Checking if file %q was processed already.", path)

	processed, err := lsp.wrapped.WasProcessed(path)

	if err != nil {
		lsp.logger.Errorf("Error checking if file %q was processed: %v", path, err)
	}

	return processed, err
}

func (lsp *loggingStorageProvider) MarkProcessed(path string) error {
	lsp.logger.Infof("Marking file %q as processed.", path)

	err := lsp.wrapped.MarkProcessed(path)

	if err != nil {
		lsp.logger.Errorf("Error marking file %q as processed: %v", path, err)
	}

	return err
}
