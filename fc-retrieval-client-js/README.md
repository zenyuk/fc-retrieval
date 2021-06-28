# Fc Retrieval Client JS

Filecoin Secondary Retrieval Market Javascript client library.

## Development

### Get Itest, Gateway and Provider

To develop on Fc Retrieval Client JS, first, clone [fc-retrieval-itest](https://github.com/ConsenSys/fc-retrieval-itest), [fc-retrieval-gateway](https://github.com/ConsenSys/fc-retrieval-gateway):, [fc-retrieval-provider](https://github.com/ConsenSys/fc-retrieval-provider):

```
git clone https://github.com/ConsenSys/fc-retrieval-itest.git
git clone https://github.com/ConsenSys/fc-retrieval-gateway.git
git clone https://github.com/ConsenSys/fc-retrieval-provider.git
```

### Build images

#### Itest images

To build images, go to `fc-retrieval-itest` and execute:

```
make lotusbase
make lotusdaemon
make lotusfullnode
```

Then build everything with:

```
make build tag
```

<i>(it can take few minutes)</i>

Everytime `/util/util.go` is changed in `fc-retrieval-itest`, all images should be rebuild.

#### Gateway image

To build Gateway image, go to `fc-retrieval-gateway` and execute:

```
make build tag
```

Everytime code is changed in `fc-retrieval-gateway`, image should be rebuild.

#### Provider image

To build Provider image, go to `fc-retrieval-provider` and execute:

```
make build tag
```

Everytime code is changed in `fc-retrieval-provider`, image should be rebuild.

### Start Itest

#### With hot reload

To start Itest with hot reload, execute:

```
RELOAD_JS_TESTS=yes go test -p 1 -v ./pkg/poc2js/poc2js_test.go -timeout=0
```

The tests will be executed, and containers will stay up.

#### Without hot reload

To start Itest without hot reload, execute:

```
go test -p 1 -v ./pkg/poc2js/poc2js_test.go
```

Now it is possible to edit Client Js `.test.ts`, save updates, and the tests will automatically rerun the `fc-retrieval-itest`.

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
