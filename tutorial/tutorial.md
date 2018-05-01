# GCSBeat Tutorial

## Introduction

This tutorial will walk you through how to use GCSBeat to read basic log files from storage.
It will take you through the steps to set up the Service Account, configure the beat, and extract events from a log file.

### Prerequisites

This tutorial expects you have the `gcsbeat` code pulled down, and have the `gcloud`, `gsutil` and `go` commands installed. If you're running in a CloudShell environment you should be all set.

## Setup

### Install the beat

Get and install the beat with the following command:

	go install github.com/GoogleCloudPlatform/gcsbeat

### Setup the environment

First set up your project ID in an environment variable:

	PROJECT_ID=my-project-id

If you're using CloudShell you can run:

	PROJECT_ID=$DEVSHELL_PROJECT_ID

## Create a Bucket

Bucket names must be universally unique, so choose something unique and set the `BUCKETNAME` environment variable:

	BUCKETNAME=gcsbeat-tutorial

Create the bucket:

	gsutil mb gs://$BUCKETNAME

Copy the demo log files in to the bucket:

	gsutil cp tutorial/*.log gs://$BUCKETNAME

List the files to make sure everything is working:

	gsutil ls gs://$BUCKETNAME

## Create the Service Account to run `gcsbeat`

Create the service account:

	gcloud iam service-accounts create gcsbeat-tutorial

Create the environment variable for its email address:

	ACCOUNTEMAIL=gcsbeat-tutorial@$PROJECT_ID.iam.gserviceaccount.com

Create new credentials to let `gcsbeat` authenticate:

	gcloud iam service-accounts keys create gcsbeat-tutorial-key.json --iam-account $ACCOUNTEMAIL
	
Grant permissions to `gcsbeat`. The most basic role you can use is `roles/storage.objectViewer`.
Here we use the `objectAdmin` role because we want to track which files have been processed by 
placing metadata tags on the files.

	gcloud projects add-iam-policy-binding $PROJECT_ID --member serviceAccount:$ACCOUNTEMAIL --role "roles/storage.objectAdmin"

The plugin also supports deleting the files after they've been read or keeping a local database of processed files.

## Configure GCSBeat

Show the contents of the configuration file:

	cat tutorial/config.conf
	
Here's a description of the variables you'll see in there and what they do:

 * `interval` - how frequently to check for new files
 * `bucket_id` - the name of the bucket you want to pull from. The `${...}` syntax refers to an environment variable
 * `json_key_file` - the authentication information `gcsbeat` will use to connect to your bucket
 * `delete` - whether or not to delete the file after it has been processed
 * `file_matches` - a glob pattern, files matching this will be processed
 * `file_exclude` - a glob pattern, files matching this will be excluded
 * `metadata_key` - if set, this metadata key will be written to the objects in your bucket that have been processed by `gcsbeat` this allows you to keep track of processing state in the event of a crash
 * `codec` - how to process the files in the bucket. `text` means one line at a time
 * `unpack_gzip` - whether or not to treat files ending with `.gz` as gzipped

## Run GCSBeat

Execute `gcsbeat` to fetch the contents of the matching log file. After starting up you will see several JSON documents show up on the screen with the contents of the included log file line by line.

	./gcsbeat -e -v -c tutorial/config.conf
	
The lines will look something like the following. Note that the text is stored in the `event` field and the `line` corresponds to the line number it came from.

	{
	  "@timestamp": "2018-05-01T15:44:50.318Z",
	  "@metadata": {
	    "beat": "gcsbeat",
	    "type": "doc",
	    "version": "7.0.0-alpha1"
	  },
	  "file": "my-log.log",
	  "line": 1,
	  "beat": {...},
	  "event": "log line 1"
	}

	
Press `Ctrl+C` once you see the events read to the console to quit the beat.

Check that the beat correctly set the tags on the files. 
You should see the key `x-goog-meta-gcsbeat-tutorial` show up under the `my-log.log` entry but not the `bak-log.log` entry.

	gsutil ls -L gs://$BUCKETNAME

Example, note the `x-goog-meta-gcsbeat-tutorial` metadata tag:

	gs://gcsbeat-tutorial/my-log.log:
	    Creation time:          Tue, 01 May 2018 15:00:00 GMT
	    Update time:            Tue, 01 May 2018 15:01:00 GMT
	    Storage class:          REGIONAL
	    Content-Length:         60
	    Content-Type:           text/x-log
	    Metadata:               
		x-goog-meta-gcsbeat-tutorial:processed


## Tear down the environment

Remove the bucket:

	gsutil rm -r gs://$BUCKETNAME

Remove the service account:

	gcloud iam service-accounts delete $ACCOUNTEMAIL

## Additional Information

Run the following for more details about configuration options:

	cat _meta/beat.yml

Visit the [project on Github](https://github.com/GoogleCloudPlatform/gcsbeat) for 
additional example configurations.
