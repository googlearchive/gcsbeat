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

package beater

import (
	"compress/gzip"
	"fmt"
	"strings"
	"time"

	"github.com/GoogleCloudPlatform/gcsbeat/beater/codec"
	"github.com/GoogleCloudPlatform/gcsbeat/beater/storage"
	"github.com/GoogleCloudPlatform/gcsbeat/config"
	mapset "github.com/deckarep/golang-set"
	"github.com/gobwas/glob"

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

// New is called by beats to instantiate the beat
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c, err := config.GetAndValidateConfig(cfg)
	if err != nil {
		return nil, err
	}

	bucket, err := storage.NewStorageProvider(c)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to bucket: %v", err)
	}

	bt := &Gcpstoragebeat{
		done:          make(chan struct{}),
		downloadQueue: make(chan string),
		config:        c,
		bucket:        bucket,
		logger:        logp.NewLogger("GCS:" + c.BucketId),
	}

	return bt, nil
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
	matcherExplaination := fmt.Sprintf("matches %q", bt.config.Match)
	excluder := glob.MustCompile(bt.config.Exclude)
	excluderExplaination := fmt.Sprintf("does not match %q", bt.config.Exclude)
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

		files, _ = storage.FilterAndExplain("already pending", files, func(path string) (bool, error) {
			return !pendingFiles.Contains(path), nil
		})

		files, _ = storage.FilterAndExplain(matcherExplaination, files, func(path string) (bool, error) {
			return matcher.Match(path), nil
		})

		files, _ = storage.FilterAndExplain(excluderExplaination, files, func(path string) (bool, error) {
			excluded := excluding && excluder.Match(path)
			return !excluded, nil
		})

		bt.logger.Infof("Added %d files to queue", len(files))
		for _, path := range files {
			bt.logger.Debugf(" - %q", path)

			bt.downloadQueue <- path
			pendingFiles.Add(path)
		}
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

	return bt.bucket.MarkProcessed(path)
}
