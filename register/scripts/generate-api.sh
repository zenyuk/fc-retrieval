#!/bin/bash

docker pull quay.io/goswagger/swagger \
&& alias swagger="docker run --rm -it  --user $(id -u):$(id -g) -e GOPATH=$HOME/go:/go -v $HOME:$HOME -w $(pwd) quay.io/goswagger/swagger" \
&& swagger version \
&& swagger generate server -f docs/swagger.yml -A register
