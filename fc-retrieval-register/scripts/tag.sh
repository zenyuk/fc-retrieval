#!/bin/bash
echo "****************************************"
echo "*** Push docker image to Docker Hub  ***"
echo "****************************************"

# Docker image name on docker hub
IMAGE_NAME=consensys/fc-retrieval-register

VERSION=$1
echo "Register version: $VERSION"
IMAGE=$2
echo "v image: $IMAGE"


ITEST_BRANCH=`git rev-parse --abbrev-ref HEAD`
echo "Register repo branch: $ITEST_BRANCH"

TAG="develop-$ITEST_BRANCH"
echo "TAG: $TAG"

docker tag $IMAGE $IMAGE_NAME:$TAG
