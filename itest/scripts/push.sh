#!/bin/bash
echo "****************************************"
echo "*** Push docker image to Docker Hub  ***"
echo "****************************************"

# Docker image name on docker hub
IMAGE_NAME=consensys/fc-retrieval-itest

VERSION=$1
echo "Integration Tests version: $VERSION"
IMAGE=$2
echo "v image: $IMAGE"


ITEST_BRANCH=`git rev-parse --abbrev-ref HEAD`
echo "Integration Tests repo branch: $ITEST_BRANCH"

if [ $ITEST_BRANCH != "main" ]; then 
    TAG="develop-$ITEST_BRANCH"
else
    TAG="develop-$VERSION"
fi
echo "TAG: $TAG"

docker tag $IMAGE $IMAGE_NAME:$TAG
echo docker push $IMAGE_NAME:$TAG
