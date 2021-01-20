# Copyright (C) 2020 ConsenSys Software Inc
FROM golang:1.15-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers git
WORKDIR /go/src/github.com/ConsenSys/fc-retrieval-gateway

# Add code to be run.
# Also grab dependancies from source directory, to improve build speed.
COPY . .


# Add the settings file needed at runtime to the bin directory so it
# can be accessed by the runtime Dockerfile.
COPY settings.json /go/bin/settings.json

# Get all dependancies, but don't install.
RUN go get -d -v github.com/ConsenSys/fc-retrieval-gateway/cmd/gateway
# Do a full compile of app and dependancies, forcing static linking.
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o /go/bin/gateway github.com/ConsenSys/fc-retrieval-gateway/cmd/gateway


# Pull build gateway into a second stage deploy alpine container
FROM alpine:latest
COPY --from=builder /go/bin/settings.json /etc/gateway/
COPY --from=builder /go/bin/gateway /gateway

# Run the binary when the container starts.
WORKDIR /
CMD ["/gateway"]
EXPOSE 8080 