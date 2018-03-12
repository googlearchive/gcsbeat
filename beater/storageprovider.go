package beater

import (
	"io"

	"cloud.google.com/go/storage"
	"github.com/spf13/afero"
	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/GoogleCloudPlatform/gcsbeat/config"
	"github.com/elastic/beats/libbeat/logp"
)

const (
	// TODO allow the metadata key to be user-defined so multiple beats/jobs can operate
	// on it simultaneously
	ProcessedMetadataKey   = "x-goog-meta-gcsbeat"
	ProcessedMetadataValue = "processed"
)

type StorageProvider interface {
	ListUnprocessed() (files []string, err error)
	Read(path string) (reader io.ReadCloser, err error)
	Remove(path string) error
	WasProcessed(path string) (bool, error)
	MarkProcessed(path string) error
}

func NewLoggingStorageProvider(inner StorageProvider) StorageProvider {
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

func NewAferoStorageProvider(fs afero.Fs) StorageProvider {
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
		if ! wasProcessed {
			out = append(out, f.Name())
		}
	}

	return out, nil
}

func (asp *aferoStorageProvider) Read(path string) (io.ReadCloser, error) {
	return asp.fs.Open(path)
}

func (asp *aferoStorageProvider) Remove(path string) error {
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

func newGcpStorageProvider(cfg *config.Config) (StorageProvider, error) {
	bucket := cfg.BucketId

	ctx := context.Background()

	options := []option.ClientOption{
		option.WithCredentialsFile(cfg.JsonKeyFile),
		option.WithUserAgent(UserAgent),
	}

	client, err := storage.NewClient(ctx, options...)

	if err != nil {
		return nil, err
	}

	return &gcpStorageProvider{ctx, client, bucket, make(map[string]bool)}, err
}

type gcpStorageProvider struct {
	ctx            context.Context
	storageClient  *storage.Client
	bucket         string
	processedCache map[string]bool
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

func isMarkedAsProcessed(metadata map[string]string) bool {
	if metadata == nil {
		return false
	}

	value, ok := metadata[ProcessedMetadataKey]

	return ok && value == ProcessedMetadataValue
}

func (gsp *gcpStorageProvider) WasProcessed(path string) (bool, error) {
	attrs, err := gsp.getAttrs(path)
	if err != nil {
		// True in the event of an error to avoid re-processing
		return true, err
	}

	return isMarkedAsProcessed(attrs.Metadata), nil
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

	metadata[ProcessedMetadataKey] = ProcessedMetadataValue

	update := storage.ObjectAttrsToUpdate{
		Metadata: metadata,
	}

	_, err = gsp.getObject(path).Update(gsp.ctx, update)
	return err
}

func (gsp *gcpStorageProvider) ListUnprocessed() ([]string, error) {
	var paths []string

	it := gsp.getBucket().Objects(gsp.ctx, nil)

	for objAttrs, err := it.Next(); err != iterator.Done; {

		if err != nil {
			return paths, err
		}

		// Eliminate these early rather than polling the network again.
		if isMarkedAsProcessed(objAttrs.Metadata) {
			continue
		}

		paths = append(paths, objAttrs.Name)
	}

	return paths, nil
}
