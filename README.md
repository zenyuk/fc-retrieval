# Fc Retrieval Client JS

Filecoin Secondary Retrieval Market Javascript client library.

## Development

### Get Itest

To develop on Fc Retrieval Client JS, first, clone [fc-retrieval-itest](https://github.com/ConsenSys/fc-retrieval-itest):

```
git clone https://github.com/ConsenSys/fc-retrieval-itest.git
```

### Build images

To build images, execute:

```
make lotusbase
make lotusdaemon
make lotusfullnode
make build tag
```

Everytime `/util/util.go` is changed in `fc-retrieval-itest`, images should be rebuild.

### Start Itest

#### With hot reload

To start Itest with hot reload, execute:

```
RELOAD_JS_TESTS=yes go test -p 1 -v ./pkg/poc2js/poc2js_test.go
```

The tests will be executed, and containers will stay up.

#### Without hot reload

To start Itest without hot reload, execute:

```
go test -p 1 -v ./pkg/poc2v2/poc2v2usage_test.go
```

The tests will be executed, and containers will stop.

## Unit tests

To run unit tests, execute:

```
npm run test
```

## Local import

### With npm link

If you need to edit this package for your application, it might be helpful to import it locally.

In this repository execute:

```
npm link
```

Then in your target app, execute:

```
npm link fc-retrieval-client-js
```

Now every update you make in the local `fc-retrieval-client-js` are directly imported in your application.

### With local package in package.json

In your target app, edit `package.json` and add:

```
{
  [...]
  "dependencies": {
    "fc-retrieval-client-js": "file:<path-to-package>/fc-retrieval-client-js",
    [...]
  },
```
