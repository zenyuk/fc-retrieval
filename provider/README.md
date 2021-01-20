# Filecoin Retrieval Pro

Filecoin Secondary Retrieval Market Provider 

## Start the service

### Create a config file

Create a `.env` file, using [.env.example](./.env.example) as a reference:

```
cp .env.example .env
```

### Start the service with Docker

Start the project with Docker:

```
make start
```

The server should be available at `http://localhost:8080`

## Development

### Start the service with Docker

Start the project with Docker:

```
make dev
```

### Start the service manually

Start the project manually:

```
make start-dev
```


## Config

Config variables description:

| name           | description    | options | default                     |
| -------------- | -------------- | ------- | --------------------------- |
| LOG_LEVEL      | logging level  |         | INFO                        |
| SERVICE_HOST   | service host   |         | provider                    |
| SERVICE_PORT   | service port   |         | 8080                        |
| SERVICE_SCHEME | service scheme |         | http                        |
