[![Build Status](https://travis-ci.org/GoogleCloudPlatform/gcsbeat.svg?branch=master)](https://travis-ci.org/GoogleCloudPlatform/gcsbeat) [![Go Report Card](https://goreportcard.com/badge/github.com/GoogleCloudPlatform/gcsbeat)](https://goreportcard.com/report/github.com/GoogleCloudPlatform/gcsbeat) [![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)



# GCSBeat

GCSBeat is an Elastic Beat to read logs/data from Google Cloud Storage (GCS) buckets.
The beat reads JSON objects or raw text from files in a bucket and forwards them to a [beats output](https://www.elastic.co/guide/en/beats/filebeat/current/configuring-output.html).

Example use-cases:

* Read [Stackdriver logs](https://cloud.google.com/stackdriver/) from a GCS bucket into Elastic.
* Read gzipped logs from cold-storage into Elastic.
* Restore data from an Elastic dump.
* Watch files on a local path (possibly a mounted GCS bucket) and upload them.
* Parse JSON logs and upload them.

Note: While this project is partially maintained by Google, this is not an official Google product.

## Configuration Options

See the [_meta/beat.yml file](./_meta/beat.yml) for a list of configuration options.

Make sure your user has permissions to read files on the bucket and write metadata.
GCSBeat marks objects as being processed using metadata keys.

### Example Configurations

Read archived redis logs from the local filesystem hourly and delete them after upload:

```yaml
gcsbeat:
  interval: 60m
  bucket_id: "file:///var/log/redis"
  delete: true
  file_matches: "*.log.gz"
  codec: "text"
  unpack_gzip: true
```

Read Stackdriver logs from a bucket:

```yaml
gcsbeat:
  bucket_id: my_log_bucket
  json_key_file: /path/to/key.json
  file_matches: "*.json"
  codec: "json-stream"
```

Read files into two separate Elastic clusters:

```yaml
# Cluster 1 beat
gcsbeat:
  bucket_id: my_log_bucket
  json_key_file: /region-one-key.json
  metadata_key: "region-one-beat"
  
# Cluster 2 beat
gcsbeat:
  bucket_id: my_log_bucket
  json_key_file: /disaster-recovery-key.json
  file_matches: "*.log"
  metadata_key: "disaster-recovery-beat"
```

Read data from a read-only bucket:

```yaml
gcsbeat:
  bucket_id: read_only_log_bucket
  json_key_file: /path/to/key.json
  file_matches: "*.log"
  processed_db_path: "processed_file_list.db"
```

## Getting Started with GCSBeat

You can either download GCSBeat compiled binaries or build them yourself.

### Download

1. Download a binary for your system on the [releases page](https://github.com/GoogleCloudPlatform/gcsbeat/releases).
2. Extract the archive and edit `gcsbeat.yml` for your environment. See above for examples.
3. Run it. For testing use: `./gcsbeat -e -v -c gcsbeat.yml`. 
   You can mock publishing for testing purposes using the `-N` argument.
4. (Optional) Run `./gcsbeat setup` to install predefined indexes. 

### Build it Yourself

You can build GCSBeat yourself using the instructions in the `DEVELOPING.md` file.

### Run

To run `gcsbeat` with info level logging configured:

```shell
./gcsbeat -c gcsbeat.yml -e -v
```

Normal Mode:

```shell
./gcsbeat -c gcsbeat.yml
```

### Debug

It can sometimes be difficult to figure out why the plugin is or isn't picking up particular files.
You can enable explain mode by running the beat with the `-d "Explain"` flag:

```shell
./gcsbeat -d "Explain" -e

```

Here's an example and how to read it:

```
The bucket gcsbeat-test has five files it can process.
Note that Cloud Storage lists directories as files.
INFO	[Explain]	storage/explain.go:55	Source "gs://gcsbeat-test" found 5 files
DEBUG	[Explain]	storage/explain.go:58	 - "bak-backup-log.log"
DEBUG	[Explain]	storage/explain.go:58	 - "new.log"
DEBUG	[Explain]	storage/explain.go:58	 - "old.log"
DEBUG	[Explain]	storage/explain.go:58	 - "test-folder/"
DEBUG	[Explain]	storage/explain.go:58	 - "test-folder/test.log"

Two of the files have already been marked as processed using the x-gcsbeat-processed metadata key.
DEBUG	[Explain]	storage/explain.go:32	Test: has key "x-gcsbeat-processed"?
DEBUG	[Explain]	storage/explain.go:43	 - "bak-backup-log.log" (pass)
DEBUG	[Explain]	storage/explain.go:43	 - "new.log" (pass)
DEBUG	[Explain]	storage/explain.go:46	 - "old.log" (fail)
DEBUG	[Explain]	storage/explain.go:43	 - "test-folder/" (pass)
DEBUG	[Explain]	storage/explain.go:46	 - "test-folder/test.log" (fail)
INFO	[Explain]	storage/explain.go:50	Test: has key "x-gcsbeat-processed"? passed 3 of 5 files

None of the remaining files are already in the pending queue.
DEBUG	[Explain]	storage/explain.go:32	Test: already pending?
DEBUG	[Explain]	storage/explain.go:43	 - "bak-backup-log.log" (pass)
DEBUG	[Explain]	storage/explain.go:43	 - "new.log" (pass)
DEBUG	[Explain]	storage/explain.go:43	 - "test-folder/" (pass)
INFO	[Explain]	storage/explain.go:50	Test: already pending? passed 3 of 3 files

Two of the files match the include filter.
DEBUG	[Explain]	storage/explain.go:32	Test: matches "*.log"?
DEBUG	[Explain]	storage/explain.go:43	 - "bak-backup-log.log" (pass)
DEBUG	[Explain]	storage/explain.go:43	 - "new.log" (pass)
DEBUG	[Explain]	storage/explain.go:46	 - "test-folder/" (fail)
INFO	[Explain]	storage/explain.go:50	Test: matches "*.log"? passed 2 of 3 files

The backup file matches the exclusion filter so it is skipped.
DEBUG	[Explain]	storage/explain.go:32	Test: does not match "bak-*"?
DEBUG	[Explain]	storage/explain.go:46	 - "bak-backup-log.log" (fail)
DEBUG	[Explain]	storage/explain.go:43	 - "new.log" (pass)
INFO	[Explain]	storage/explain.go:50	Test: does not match "bak-*"? passed 1 of 2 files

Exactly one file remained.
INFO	[GCS:gcsone]	beater/gcsbeat.go:131	Added 1 files to queue
```

## License

```
Copyright (c) 2018 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
