#!/bin/bash

# Pull down ASDF

ASDFPATH=$(which asdf)
if [ "$ASDFPATH" = "" ] 
then
	echo "Installing CURL"
	apk update
	apk add curl

	echo "Downloading asdf"
	curl -L https://github.com/asdf-vm/asdf/archive/v0.5.0.tar.gz > asdf.tar.gz
	tar -xzvf asdf.tar.gz
	source ./asdf-0.5.0/asdf.sh
fi

# Install all required packages

install() {
	echo "Installing $1"
	asdf plugin-add $1
	asdf install $1 $2
	asdf global $1 $2
}

install golang $GO_VERSION # 1.10.3
install python $PYTHON_VERSION # 2.7.15


# Show current packages

asdf current
