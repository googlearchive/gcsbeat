# GCSBeat

GCSBeat is an elastic Beat for Google Cloud Storage.
The beat reads objects from a specified bucket line by line, forwards them to your configured 
outputs then deletes the file or marks it as processed with a metadata attribute to avoid processing
it again.

## Getting Started with GCSBeat

### Requirements

* [Golang](https://golang.org/dl/) 1.10


### Build

To build the binary for GCSBeat run the command below. This will generate a binary
in the same directory with the name GCSBeat.

```
make
```


### Run

To run GCSBeat with debugging output enabled, run:

```
./gcsbeat -c gcsbeat.yml -e -d "*"
```


### Test

To test GCSBeat, run the following command:

```
make testsuite
```

alternatively:
```
make unit-tests
make system-tests
make integration-tests
make coverage-report
```

The test coverage is reported in the folder `./build/coverage/`

### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `fields.yml` by running the following command.

```
make update
```


### Cleanup

To clean  GCSBeat source code, run the following commands:

```
make fmt
make simplify
```

To clean up the build directory and generated artifacts, run:

```
make clean
```


### Clone

To clone GCSBeat from the git repository, run the following commands:

```
mkdir -p ${GOPATH}/src/github.com/GoogleCloudPlatform/gcsbeat
git clone https://github.com/GoogleCloudPlatform/gcsbeat ${GOPATH}/src/github.com/GoogleCloudPlatform/gcsbeat
```


For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).


## Packaging

The beat frameworks provides tools to crosscompile and package your beat for different platforms. This requires [docker](https://www.docker.com/) and vendoring as described above. To build packages of your beat, run the following command:

```
make package
```

This will fetch and create all images required for the build process. The hole process to finish can take several minutes.
