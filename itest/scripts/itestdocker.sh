#!/bin/bash
echo "****************************************"
echo "***       Running itestdocker        ***"
echo "****************************************"

ITEST_BRANCH=`git rev-parse --abbrev-ref HEAD`

if [ $ITEST_BRANCH == "main" ]; then 
    echo "Itest is running in main, itestdocker runs."
    go test -v -p=1 --count=1 ./...
else
    # Check if this branch is in an active PR
    if curl \
        -H "Accept: application/vnd.github.v3+json" \
        https://api.github.com/repos/consensys/fc-retrieval/itest/pulls | jq '.[].head.ref' | grep $ITEST_BRANCH; then
        echo "Itest is running in a branch that is open for PR, itestdocker runs."
        go test -v -p=1 --count=1 ./...
    elif curl \
        -H "Accept: application/vnd.github.v3+json" \
        https://api.github.com/repos/consensys/fc-retrieval/provider/pulls | jq '.[].head.ref' | grep $ITEST_BRANCH; then
        echo "Itest is running in a branch that is open for PR in provider (Test new provider), itestdocker runs."
        go test -v -p=1 --count=1 ./...
    elif curl \
        -H "Accept: application/vnd.github.v3+json" \
        https://api.github.com/repos/consensys/fc-retrieval/gateway/pulls | jq '.[].head.ref' | grep $ITEST_BRANCH; then
        echo "Itest is running in a branch that is open for PR in gateway (Test new gateway), itestdocker runs."
        go test -v -p=1 --count=1 ./...
    elif curl \
        -H "Accept: application/vnd.github.v3+json" \
        https://api.github.com/repos/consensys/fc-retrieval/register/pulls | jq '.[].head.ref' | grep $ITEST_BRANCH; then
        echo "Itest is running in a branch that is open for PR in register (Test new register), itestdocker runs."
        go test -v -p=1 --count=1 ./...
    else
        echo "Itest is running in a branch that isn't open for PR, and we don't have a newer provider/gateway/register, itestdocker does not run."
    fi
fi