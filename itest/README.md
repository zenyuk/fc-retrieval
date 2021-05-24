# fc-retrieval-itest

Filecoin Secondary Retrieval Market Integration Tests

See also Filecoin Secondary Retrieval Market
[client library](https://github.com/ConsenSys/fc-retrieval-client)
[gateway](https://github.com/ConsenSys/fc-retrieval-gateway) and
[retrieval provider](https://github.com/ConsenSys/fc-retrieval-provider) repositories.

## Run tests on a local Devnet

###Build images

```
make lotusbase lotusfullnode lotusdaemon build tag
```

### Execute tests

#### Clean mode:
```
go test -p 1 ./...
```

#### Debug mode:
```
go test -p 1 -v ./...
```
