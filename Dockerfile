# Run ping application
FROM golang:1.15

WORKDIR /go/src/app
COPY . .

RUN go install -v ./cmd/ping.go

CMD ["ping"]
