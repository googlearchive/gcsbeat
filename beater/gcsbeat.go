package beater

import (
	"compress/gzip"
	"fmt"
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/gcsbeat/beater/codec"
	"github.com/GoogleCloudPlatform/gcsbeat/beater/storage"
	"github.com/GoogleCloudPlatform/gcsbeat/config"
	"github.com/deckarep/golang-set"
	"github.com/gobwas/glob"
	"github.com/spf13/afero"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
)

type Gcpstoragebeat struct {
	done          chan struct{}
	downloadQueue chan string
	config        *config.Config
	client        beat.Client
	bucket        storage.StorageProvider
	logger        *logp.Logger
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c, err := config.GetAndValidateConfig(cfg)
	if err != nil {
		return nil, err
	}

	bucket, err := connectToBucket(c)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to bucket: %v", err)
	}

	// TODO check that we have update permission so we can tag the files as updated

	bt := &Gcpstoragebeat{
		done:          make(chan struct{}),
		downloadQueue: make(chan string),
		config:        c,
		bucket:        storage.NewLoggingStorageProvider(bucket),
		logger:        logp.NewLogger("GCS:" + c.BucketId),
	}

	return bt, nil
}

func connectToBucket(cfg *config.Config) (storage.StorageProvider, error) {
	if strings.HasPrefix(cfg.BucketId, "file://") {
		basePath := cfg.BucketId[7:]
		fs := afero.NewBasePathFs(afero.NewOsFs(), basePath)
		return storage.NewAferoStorageProvider(fs), nil
	}

	// connect to GCP
	return storage.NewGcpStorageProvider(cfg)
}

func (bt *Gcpstoragebeat) Run(b *beat.Beat) error {
	bt.logger.Info("GCP storage beat is running! Hit CTRL-C to stop it.")
	bt.logger.Infof("Version: %q", storage.GetUserAgent())
	bt.logger.Infof("Config: %+v", bt.config)

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	// Start background jobs
	go bt.fileChangeWatcher()
	go bt.fileDownloader()

	ticker := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
			bt.logger.Infof("Pending Downloads: %v", len(bt.downloadQueue))
		}
	}
}

func (bt *Gcpstoragebeat) Stop() {
	bt.client.Close()
	close(bt.done)
}

func (bt *Gcpstoragebeat) fileChangeWatcher() {
	pendingFiles := mapset.NewSet()
	ticker := time.NewTicker(bt.config.Interval)
	matcher := glob.MustCompile(bt.config.Match)
	excluder := glob.MustCompile(bt.config.Exclude)
	excluding := bt.config.Exclude != ""

	for {
		select {
		case <-bt.done:
			return
		case <-ticker.C:
		}

		files, err := bt.bucket.ListUnprocessed()
		if err != nil {
			continue
		}

		alreadyPending := 0
		excluded := 0
		queued := 0
		for _, path := range files {
			if pendingFiles.Contains(path) {
				alreadyPending++
				continue
			}

			isIncluded := matcher.Match(path)
			isExcluded := excluding && excluder.Match(path)

			if !isIncluded || isExcluded {
				excluded++
				continue
			}

			bt.downloadQueue <- path
			pendingFiles.Add(path)
			queued++

			bt.logger.Infof("Found file: %q to import", path)
		}

		bt.logger.Infof("Found %d files, already pending: %d, regex excluded: %d, new: %d", len(files), alreadyPending, excluded, queued)
	}
}

func (bt *Gcpstoragebeat) fileDownloader() {
	for {
		select {
		case <-bt.done:
			return
		case path := <-bt.downloadQueue:
			bt.downloadFile(path)
		}
	}
}

func (bt *Gcpstoragebeat) downloadFile(path string) {
	bt.logger.Infof("Starting to download and parse: %q", path)

	input, err := bt.bucket.Read(path)

	if err != nil {
		return
	}

	defer input.Close()

	if bt.config.UnpackGzip && strings.HasSuffix(path, ".gz") {
		gzReader, err := gzip.NewReader(input)

		if err != nil {
			bt.logger.Errorf("Error parsing file %q: %v", path, err)
			return
		}

		defer gzReader.Close()
		input = gzReader
	}

	codec, err := codec.NewCodec(bt.config.Codec, path, input)
	if err != nil {
		bt.logger.Errorf("Error parsing file %q: %v", path, err)
		return
	}

	for codec.Next() {
		event := beat.Event{
			Timestamp: time.Now(),
			Fields:    codec.Value(),
		}

		bt.client.Publish(event)
	}

	if err := codec.Err(); err != nil {
		bt.logger.Errorf("Error parsing file %q: %v", path, err)
		return
	}

	bt.closeOutFile(path)
}

func (bt *Gcpstoragebeat) closeOutFile(path string) error {
	if bt.config.Delete {
		return bt.bucket.Remove(path)
	}

	// TODO add option (middleware?) to save a processed list locally
	return bt.bucket.MarkProcessed(path)

	// TODO add options to back up to another bucket or save files locally
}
