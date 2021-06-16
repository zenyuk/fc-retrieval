# Copyright (C) 2020 ConsenSys Software Inc
FROM golang:1.16-alpine

RUN apk add --no-cache make gcc musl-dev linux-headers git nodejs npm

WORKDIR /go/src/github.com/ConsenSys/fc-retrieval-itest
COPY . .
# Remove any cached dependancies. TODO is this really needed?
RUN go clean -modcache

RUN go mod download
