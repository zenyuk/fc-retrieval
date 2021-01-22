# fc-retrieval-provider
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

### Start the service manually

Start the project manually:

```
make start-dev
```

The server should be available at `http://localhost:8080`

## Config

Config variables description:

| name            | description         | options       | default                     |
| --------------- | ------------------- | ------------- | --------------------------- |
| LOG_LEVEL       | logging level       |               | INFO                        |
| LOG_TARGET      | logging target      | STDOUT / FILE | STDOUT                      |
| LOG_DIR         | logging directory   |               |                             |
| LOG_FILE        | logging file        |               |                             |
| LOG_MAX_BACKUPS | logging max backups |               |                             |
| LOG_MAX_AGE     | logging max age     |               |                             |
| LOG_MAX_SIZE    | logging max size    |               | INFO                        |
| LOG_COMPRESS    | logging compress    | true / false  |                             |
| SERVICE_HOST    | service host        |               | provider                    |
| SERVICE_PORT    | service port        |               | 8080                        |
| SERVICE_SCHEME  | service scheme      |               | http                        |
