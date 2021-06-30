# FC Retrieval Demo

## 1. Get demo repositories

First, get all the repositories, and checkout `264-lotus-upgrade` branch.

To get all the repositories, execute:
```
git clone https://github.com/ConsenSys/fc-retrieval-itest.git
git clone https://github.com/ConsenSys/fc-retrieval-common.git
git clone https://github.com/ConsenSys/fc-retrieval-provider.git
git clone https://github.com/ConsenSys/fc-retrieval-gateway.git
git clone https://github.com/ConsenSys/fc-retrieval-client.git
git clone https://github.com/ConsenSys/fc-retrieval-register.git
```

To checkout the `264-lotus-upgrade` branch, execute:

```
cd fc-retrieval-itest/ && git checkout 264-lotus-upgrade && cd ..
cd fc-retrieval-common/ && git checkout 264-lotus-upgrade && cd ..
cd fc-retrieval-provider/ && git checkout 264-lotus-upgrade && cd ..
cd fc-retrieval-gateway/ && git checkout 264-lotus-upgrade && cd ..
cd fc-retrieval-client/ && git checkout 264-lotus-upgrade && cd ..
cd fc-retrieval-register/ && git checkout 264-lotus-upgrade && cd ..
```

## 2. Build and setup demo

### Build Provider, Gateway and Register images

To build  Provider, Gateway and Register images, execute:
```
cd fc-retrieval-provider/ && make buildlocal tag && cd ..
cd fc-retrieval-gateway/ && make buildlocal tag && cd ..
cd fc-retrieval-register/ && make buildlocal tag && cd ..
```

### Build lotus node image

To build a Lotus node image, go to Itest folder:

```
cd fc-retrieval-itest/
```

And execute the build commands:

```
make lotusbase
make lotusdaemon
make lotusfullnode
```

It can be very long, `lotusbase` needs 52 blocks to complete.

### Docker network

Create a docker network
```
docker network create shared
```

## 3. Start demo

In `fc-retrieval-itest`, checkout the `demo` branch:

```
git checkout demo
cd demo
```

Open 4 terminals on `fc-retrieval-itest/demo` folder, and execute:

<b>Terminal 1</b> will run all the services:
```
docker-compose up
```

<b>Terminal 2</b> will run a Provider admin instance, once started it need to be `init`:
```
go run provideradmin-cli/main.go
init
```

<b>Terminal 3</b> will run a Gateway admin instance, once started it need to be `init`:
```
go run gatewayadmin-cli/main.go
init
```

<b>Terminal 4</b> will run a Client instance, once started it need to be `init`:
```
go run client-cli/main.go
init
```

## 4. Use demo

### Service registration

<b>Terminal 2</b>

To register the Provider, execute:

```
register-pvd
```
it returns a [providerID], for ex: 
```
>>> register-pvd

fe33825727c649ad6187832f6575bd119e0a2408fb6d180808d92d4bcf4bc51e
```

<b>Terminal 3</b>

To register the Gateway, execute:

```
register-gw
```
it returns a [gatewayID], for example
```
>>> register-gw

4d8c17b8077633a8b0d32b941f76f806bfc1fc661b2d48820ec5e12c436bcb4c
```

<b>Terminal 4</b>

To add an active Gateway on the client, execute:
```
add-active [gatewayID]
```

For example: 

```
>>> add-active 4d8c17b8077633a8b0d32b941f76f806bfc1fc661b2d48820ec5e12c436bcb4c
>>> ls-active

Current gateways in use:
4d8c17b8077633a8b0d32b941f76f806bfc1fc661b2d48820ec5e12c436bcb4c
```


### Publish an offer

<b>Terminal 2</b>

To publish an offer from the Provider, execute:

```
publish-random-offer [providerID]
```

It creates an offer with 3 cids, for example:
```
>>> publish-random-offer fe33825727c649ad6187832f6575bd119e0a2408fb6d180808d92d4bcf4bc51e

Published offer for cid: [
85ff0543ad514c3dc7a6b98fbe16c6b9c0ff0ee3c41f06f94a37a105a802d82f
6bb876b19a2757c96a78df7e434be9122470ed8e9c8ebabc4433609ed24846cb
bc9d5549ec4c6b3807e60ba4e4137381d138f1295c5dea172e542d6e5336dcc7
] at a price of 14
```


### Get an offer

<b>Terminal 4</b>

Finally, get an offer from the client by executing:
```
find-offer [cid1]
```

For example:

```
>>> find-offer 85ff0543ad514c3dc7a6b98fbe16c6b9c0ff0ee3c41f06f94a37a105a802d82f

Find offers:
Offer 00000001d68e4f5cea395eccf97540f3c1d61e97ac6525ca3cf6498649b106b03c1f6ddd /03152e661f64ceb768e77563a7a40f3292bcb436b24dd9eb78787728f3fcbb8e01: provider-fe33825727c649ad6187832f6575bd119e0a2408fb6d180808d92d4bcf4bc51e, cid-85ff0543ad514c3dc7a6b98fbe16c6b9c0ff0ee3c41f06f94a37a105a802d82f, price-14, expiry-1624437694, qos-42
```

## 5. Stop demo

<b>Terminal 1</b>

To stop service, exit docker and execute:

```
docker-compose down
```

<b>Terminal 2</b>

To exit Provider admin, execute:

```
exit
```
<b>Terminal 3</b>

To exit Gateway admin, execute:

```
exit
```
<b>Terminal 4</b>

To exit Client, execute:

```
exit
```