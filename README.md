# Filecoin Secondary Retrieval Market

## prerequisite
Go - https://golang.org/dl/  
Docker - https://www.docker.com/products/docker-desktop

## clone
```
git clone git@github.com:ConsenSys/fc-retrieval.git
```

## build Docker images
```
cd fc-retrieval
docker build -t consensys/fc-retrieval/gateway -f gateway/Dockerfile .
docker build -t consensys/fc-retrieval/provider -f provider/Dockerfile .
docker build -t consensys/fc-retrieval/register -f register/Dockerfile .
docker build -t consensys/fc-retrieval/itest -f itest/Dockerfile .
```

## run tests
### build Lotus Full Node for end-to-end tests
```
cd fc-retrieval/itest/lotus/lotus-base
docker build -t consensys/lotus-base .


cd fc-retrieval/itest/lotus/lotus-full-node
docker build -t consensys/lotus-full-node .
```

### execute end-to-end tests 
```
cd fc-retrieval/itest
go test -p 1 -v ./...
```
