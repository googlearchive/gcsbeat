# GCSBeat Tutorial

## Introduction

This tutorial will walk you through how to use GCSBeat to read basic log files from storage.
It will take you through the steps to set up the Service Account, configure the beat, and extract events from a log file.

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
	
Grant permissions to `gcsbeat`. The most basic permissions you can use are `roles/storage.objectViewer`.
Additional functionality 

	gcloud projects add-iam-policy-binding $PROJECT_ID --member serviceAccount:$ACCOUNTEMAIL --role "roles/storage.objectAdmin"

## Configure GCSBeat

Show the contents of the configuration file:

	cat tutorial/config.conf

## Run GCSBeat

Execute `gcsbeat` to fetch the

	./gcsbeat -e -v -c tutorial/config.conf
	
Press `Ctrl+C` once you see the events read to the console to quit the beat.

Check that the beat correctly set the tags on the files:

	gsutil ls -L gs://$BUCKETNAME

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
