#!/bin/bash
# Check environment

require() {
	if [ "$1" == "" ]
	then
		echo "The environment variable $2 must be present and contain $3"
		exit 4
	fi
}

section() {
	echo -e "\n============================\n= $1\n============================\n"
}

section "Checking config"
require "$SERVICE_ACCOUNT" "SERVICE_ACCOUNT" "the service account JSON gcsbeat will use"
require "$BUCKET" "BUCKET" "the bucket name to hold the test files"

section "Building gcsbeat"
source $(dirname $0)/setup-go.sh
make

section "Pulling down gcloud"
curl https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-204.0.0-linux-x86_64.tar.gz > gcloud.tar.gz
tar xfz gcloud.tar.gz

./google-cloud-sdk/bin/gcloud version

section "Setting up environment"
export TESTID=`date +%s`
export TESTDIR="/tmp/$TESTID"
export KEYPATH="$TESTDIR/key.json"

mkdir -p $TESTDIR

echo "Test ID: $TESTID"
echo $SERVICE_ACCOUNT > $KEYPATH
ls -la $TESTDIR
./google-cloud-sdk/bin/gcloud auth activate-service-account --key-file $KEYPATH


section "Setting up test files"

echo "000) Raw Text" > "$TESTDIR/raw.txt"
echo "001) Gzipped Text" | gzip > "$TESTDIR/gzipped.txt.gz"
echo -e "000) Raw Text\n001) Gzipped Text" > "$TESTDIR/expected.txt"

./google-cloud-sdk/bin/gsutil cp "$TESTDIR/raw.txt" "gs://$BUCKET/$TESTID/log-raw.txt"
./google-cloud-sdk/bin/gsutil cp "$TESTDIR/gzipped.txt.gz" "gs://$BUCKET/$TESTID/log-gzipped.txt.gz"


section "Running gcsbeat"
timeout 30s ./gcsbeat -e -v -c ci/fixtures/integration-test-config.yml

section "Checking Results"
cat "$TESTDIR/actual.txt" | sort > "$TESTDIR/actual-sorted.txt"
echo "Expected"
cat -n "$TESTDIR/expected.txt"
echo "Actual"
cat -n "$TESTDIR/actual-sorted.txt"
ls -lah "$TESTDIR"

cmp "$TESTDIR/actual-sorted.txt" "$TESTDIR/expected.txt"
TEST_RESULT=$?
echo "Result: $TEST_RESULT"


section "Tearing down environment"
./google-cloud-sdk/bin/gsutil rm -r "gs://$BUCKET/$TESTID"
rm -rf $TESTDIR

exit $TEST_RESULT
