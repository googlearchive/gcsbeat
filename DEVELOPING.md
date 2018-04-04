# Developing

Thank you for thinking about adding new features to this GCP beat!

In general the steps for developing are: 

1. Read the [contributing file](./CONTRIBUTING.md) file. 
   We follow standard conventions for Github development.
2. Sign the [Google CLA](https://cla.developers.google.com) if you haven't already.
3. Write your code and tests.
4. Run `make pre-commit` before submitting pull-requests.
   This formats your code, tests it and updates the documentation.
5. Create a pull request.
6. That's it!

## Building

### Requirements

We support the most recent version of `go`. 
If you need a feature from a newer version listed here changes are welcome.

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

We do not use the standard beats release process to build packages.
Use the `release` make target instead.

```
make release
```

This process is faster and doesn't require docker to build releases.
