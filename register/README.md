# Filecoin Retrieval Register

Filecoin Secondary Retrieval Market Register

## Start the service

### Create a config file

Create a `.env` file, using [.env.example](./.env.example) as a reference:

```
cp .env.example .env
```

## Development

### Start the service with Docker

Start the service:

```
make dev
```

The API should be available at `http://localhost:9020`

To rebuild containers, following option can be used:

```
make dev arg=--build
```

### Start the service manually

Start the project manually:

```
go run cmd/register-server/main.go --host 0.0.0.0 --port 9020
```

The API should be available at `http://localhost:9020`

## Config

Config variables description:

| name            | description         | options       | default                     |
| --------------- | ------------------- | ------------- | --------------------------- |
| SERVICE_NAME    | service name        |               | Filecoin Retrieval Register |
| REDIS_URL       | redis url           |               | redis                       |
| REDIS_PORT      | redis port          |               | 6379                        |
| REDIS_PASSWORD  | redis password      |               | ""                          |
| LOG_LEVEL       | logging level       |               | info                        |
| LOG_TARGET      | logging target      | STDOUT / FILE | STDOUT                      |
| LOG_DIR         | logging directory   |               |                             |
| LOG_FILE        | logging file        |               |                             |
| LOG_MAX_BACKUPS | logging max backups |               |                             |
| LOG_MAX_AGE     | logging max age     |               |                             |
| LOG_MAX_SIZE    | logging max size    |               |                             |
| LOG_COMPRESS    | logging compress    | true / false  |                             |

## Usage

### Swagger documentation

Swagger Yaml can be found in `docs/swagger.yml`.

Once the service is started, Swagger Ui can be found at `<service_url>/docs`.

### Demo

#### Gateway registration

Register a Gateway

```
curl --location --request POST 'http://localhost:9020/registers/gateway' \
--header 'Content-Type: application/json' \
--data-raw '{
    "address": "f01234",
    "networkInfoClient": "127.0.0.1:9010",
    "networkInfoProvider": "127.0.0.1:9011",
    "networkInfoGateway": "127.0.0.1:9012",
    "networkInfoAdmin": "127.0.0.1:9013",
    "regionCode": "FR",
    "rootSigningKey": "0xABCDE123456789",
    "signingKey": "0x987654321EDCBA"
}'
```

Get gateway register list

```
curl --location --request GET 'http://localhost:9020/registers/gateway'
```

#### Provider registration

Register a Provider

```
curl --location --request POST 'http://localhost:9020/registers/provider' \
--header 'Content-Type: application/json' \
--data-raw '{
    "address": "f01234",
    "networkInfo": "127.0.0.1:9030",
    "regionCode": "FR",
    "rootSigningKey": "0xABCDE123456789",
    "signingKey": "0x987654321EDCBA"
}'
```

Get provider register list

```
curl --location --request GET 'http://localhost:9020/registers/provider'
```

## Development

### Generate API

#### Update and regenerate API

Routes can be changed in `docs/swagger.yml`. To generate or regenerate the project with updated routes, execute:

```
swagger generate server -f docs/swagger.yml -A register
```

### FAQ

#### Swagger command error

When generating swagger command, it the error `target must reside inside a location in the $GOPATH/src or be a module` appaers, execute:

```
go mod init github.com/ConsenSys/fc-retrieval-register
```
