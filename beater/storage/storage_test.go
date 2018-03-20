// +build !integration

package storage

import (
	"io/ioutil"
	"path"
	"testing"

	"github.com/spf13/afero"
)

type spTestCase struct {
	name     string
	provider StorageProvider
}

func getAferoProvider() StorageProvider {
	fs := afero.NewMemMapFs()

	provider := newAferoStorageProvider(fs)

	fd, _ := fs.Create("exists.log")
	fd.WriteString("log1\nlog2")
	fd.Close()

	return provider
}

func getLocalProcessedProvider() StorageProvider {
	provider := getAferoProvider()
	tmp, _ := ioutil.TempDir("", "gcsbeattest")
	tmpFile := path.Join(tmp, "processed.db")

	base, _ := newLocalProcessedMiddlewareBase(provider, tmpFile, "test-key")

	return base
}

func setupSpTestCases() []spTestCase {
	return []spTestCase{
		{"afero", getAferoProvider()},
		{"localprocessed", getLocalProcessedProvider()},
		{"log", newLoggingStorageProvider(getAferoProvider())},
	}
}

func TestStorageProviderListUnprocessed(t *testing.T) {
	for _, sp := range setupSpTestCases() {
		t.Run(sp.name, func(t *testing.T) {
			// on init we should have 1 path
			paths, err := sp.provider.ListUnprocessed()
			if err != nil || len(paths) != 1 {
				t.Errorf("Expected no error %v, and 1 path: %v", err, paths)
			}

			for _, p := range paths {
				sp.provider.MarkProcessed(p)
			}

			// after all the paths should be consumed
			paths, err = sp.provider.ListUnprocessed()
			if err != nil || len(paths) != 0 {
				t.Errorf("Expected no error %v, and 0 paths: %v", err, paths)
			}
		})
	}
}

func TestStorageProviderRead(t *testing.T) {
	for _, sp := range setupSpTestCases() {
		t.Run(sp.name, func(t *testing.T) {
			if r, err := sp.provider.Read("exists.log"); err != nil || r == nil {
				t.Errorf("Expected no error %v and non-nil reader %v", err, r)
			}

			if _, err := sp.provider.Read("does/not/exist.log"); err == nil {
				t.Errorf("Expected error reading file that does not exist")
			}
		})
	}
}

func TestStorageProviderRemove(t *testing.T) {
	for _, sp := range setupSpTestCases() {
		t.Run(sp.name, func(t *testing.T) {
			if err := sp.provider.Remove("does/not/exist.log"); err == nil {
				t.Errorf("Expected error for removing a non-existant file")
			}

			if err := sp.provider.Remove("exists.log"); err != nil {
				t.Errorf("Expected no error for removing an existing file, got: %v", err)
			}

			// Removing a second time shouldn't be valid because the file is gone
			if err := sp.provider.Remove("exists.log"); err == nil {
				t.Errorf("Expected error for removing a file that was already deleted")
			}
		})
	}
}
