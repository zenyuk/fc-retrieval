# fc-retrieval-itest

Filecoin Secondary Retrieval Market Integration and End-To-End Tests

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

#### POC2 v2 ClientJS watch(hot reload e2e tests) mode:
```
# minimal network
RELOAD_JS_TESTS=yes go test -p 1 -v ./pkg/poc2js/poc2js_test.go
```
