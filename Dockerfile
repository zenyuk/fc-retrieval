# Copyright (C) 2020 ConsenSys Software Inc
FROM golang:1.15-alpine as builder
RUN apk update && apk add --no-cache make gcc musl-dev linux-headers git
WORKDIR /go/src/github.com/ConsenSys/fc-retrieval-register/
COPY . .

# Get all dependancies, but don't install.
RUN go get -d -v github.com/ConsenSys/fc-retrieval-register/cmd/register-server
# Do a full compile of app and dependancies, forcing static linking.
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o /go/bin/register-server ./cmd/register-server

# Pull build register into a second stage deploy alpine container
FROM alpine:latest
COPY --from=builder /go/bin/register-server /register-server

# Run the binary when the container starts.
WORKDIR /
CMD ["/register-server", "--host", "0.0.0.0", "--port", "8080"]
EXPOSE 8080