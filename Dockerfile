# Copyright (C) 2020 ConsenSys Software Inc
FROM node:lts-alpine AS node
FROM golang:1.15-alpine

RUN apk add --no-cache make gcc musl-dev linux-headers git

WORKDIR /go/src/github.com/ConsenSys/fc-retrieval-itest
COPY . .
# Remove any cached dependancies. TODO is this really needed?
RUN go clean -modcache

RUN go mod download

COPY --from=node /usr/local/bin/ /usr/local/bin/