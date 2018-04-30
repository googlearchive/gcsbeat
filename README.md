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

### Tutorial

We have a quick tutorial that will walk you through setting up permissions, creating a GCS bucket and reading data from it with gcsbeat. You can read it in the `tutorial/tutorial.md` file or click below to launch an interactive version.

[![Open in Cloud Shell](http://gstatic.com/cloudssh/images/open-btn.svg)](https://console.cloud.google.com/cloudshell/open?git_repo=https%3A%2F%2Fgithub.com%2FGoogleCloudPlatform%2Fgcsbeat&page=shell&tutorial=tutorial%2Ftutorial.md)

### Download

1. Download a binary for your system on the [releases page](https://github.com/GoogleCloudPlatform/gcsbeat/releases).
2. Extract the archive and edit `gcsbeat.yml` for your environment. See above for examples.
3. Run it. For testing use: `./gcsbeat -e -v -c gcsbeat.yml`. 
   You can mock publishing for testing purposes using the `-N` argument.
4. (Optional) Run `./gcsbeat setup` to install predefined indexes. 

### Requirements

* [Golang](https://golang.org/dl/) 1.10
* `virtualenv` >= 15.1.*
* `python` 2.7.*

### Build

To build the binary for GCSBeat run the command below. This will generate a binary
in the same directory with the name `gcsbeat`.

```shell
# Clean the beat, update the docs and build it
make clean && make update && make
```

### Run

To run `gcsbeat` with info level logging configured:

```shell
./gcsbeat -c gcsbeat.yml -e -v
```

Normal Mode:

```shell
./gcsbeat -c gcsbeat.yml
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
