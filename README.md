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

Start the project with Docker:

```
make dev
```

The API should be available at `http://localhost:8080`

### Start the service manually

Start the project manually:

```
go run cmd/filecoin-retrieval-register-server/main.go --host 0.0.0.0 --port 8080
```

The API should be available at `http://localhost:8080`

## Config

Config variables description:

| name           | description    | options | default                     |
| -------------- | -------------- | ------- | --------------------------- |
| SERVICE_NAME   | service name   |         | Filecoin Retrieval Register |
| REDIS_URL      | redis url      |         | redis                       |
| REDIS_PORT     | redis port     |         | 6379                        |
| REDIS_PASSWORD | redis password |         | ""                          |

## Usage

### Swagger documentation

Swagger Yaml can be found in `docs/swagger.yml`.

Once the service is started, Swagger Ui can be found at `<service_url>/docs`.

### Demo

Register a Gateway

```
curl --location --request POST 'http://localhost:8080/registers/gateway' \
--header 'Content-Type: application/json' \
--data-raw '{
    "address": "f01234",
    "networkInfo": "127.0.0.1:80",
    "regionCode": "FR",
    "rootSigningKey": "0xABCDE123456789",
    "sigingKey": "0x987654321EDCBA"
}'
```

Get gateway register list

```
curl --location --request GET 'http://localhost:8080/registers/gateway'
```


Register a Provider

```
curl --location --request POST 'http://localhost:8080/registers/gateway' \
--header 'Content-Type: application/json' \
--data-raw '{
    "address": "f01234",
    "networkInfo": "127.0.0.1:80",
    "regionCode": "FR",
    "rootSigningKey": "0xABCDE123456789",
    "sigingKey": "0x987654321EDCBA"
}'
```

Get provider register list

```
curl --location --request GET 'http://localhost:8080/registers/provider'
```

## Development

### Generate API

#### Update and regenerate API

Routes can be changed in `docs/swagger.yml`. To generate or regenerate the project with updated routes, execute:

```
swagger generate server -f docs/swagger.yml
```

### FAQ

#### Swagger command error

When generating swagger command, it the error `target must reside inside a location in the $GOPATH/src or be a module` appaers, execute:

```
go mod init github.com/ConsenSys/fc-retrieval-register
```
