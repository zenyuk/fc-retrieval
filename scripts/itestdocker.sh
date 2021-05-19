#!/bin/bash
echo "****************************************"
echo "***       Running itestdocker        ***"
echo "****************************************"

ITEST_BRANCH=`git rev-parse --abbrev-ref HEAD`

if [ $ITEST_BRANCH == "main" ]; then 
    echo "CircleCI is running in a main, itestdocker runs."
    go test -v -p=1 --count=1 ./...
else
    # Check if this branch is in an active PR
    if curl \
        -H "Accept: application/vnd.github.v3+json" \
        https://api.github.com/repos/consensys/fc-retrieval-itest/pulls | jq '.[].head.ref' | grep $ITEST_BRANCH; then
        echo "CircleCI is running in a branch that is open for PR, itestdocker runs."
        go test -v -p=1 --count=1 ./...
    else
        echo "CircleCI is running in a branch that isn't open for PR, itestdocker does not run."
    fi
fi