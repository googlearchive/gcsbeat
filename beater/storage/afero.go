package storage

import (
	"github.com/spf13/afero"
	"io"
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

		wasProcessed, _ := asp.WasProcessed(f.Name())
		if !wasProcessed {
			out = append(out, f.Name())
		}
	}

	return out, nil
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
