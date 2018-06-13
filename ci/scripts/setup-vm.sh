#!/bin/bash

# install gsutil
curl https://storage.googleapis.com/pub/gsutil.tar.gz > gsutil.tar.gz
tar xfz gsutil.tar.gz

echo $SERVICE_ACCOUNT > sa.json
./gsutil/gsutil config -e sa.json
