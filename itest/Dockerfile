# Copyright (C) 2020 ConsenSys Software Inc
FROM golang:1.15-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers git

WORKDIR /go/src/github.com/ConsenSys/fc-retrieval-itest
COPY . .
# Remove any cached dependancies. TODO is this really needed?
RUN go clean -modcache

# Get all dependancies, but don't install.
RUN go get -d -v github.com/ConsenSys/fc-retrieval-itest/cmd/itest
# Do a full compile of app and dependancies, forcing static linking.
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o /go/bin/itest github.com/ConsenSys/fc-retrieval-itest/cmd/itest


# Pull build gateway into a second stage deploy alpine container
FROM alpine:latest
COPY --from=builder /go/bin/itest /itest

# Run the binary when the container starts.
WORKDIR /
CMD ["/itest"]
