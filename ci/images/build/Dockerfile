FROM alpine:latest

RUN apk add --no-cache curl bash coreutils tar gzip git make go python2 py2-virtualenv
RUN curl -L https://github.com/asdf-vm/asdf/archive/v0.5.0.tar.gz > asdf.tar.gz && \
    tar -xzvf asdf.tar.gz && \
    mv asdf-0.5.0 /asdf

ENV PATH="/asdf/shims:/asdf/bin:${PATH}"

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH
