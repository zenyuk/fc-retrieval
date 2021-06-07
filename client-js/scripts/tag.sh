#!/bin/bash
echo "****************************************"
echo "*** Tag docker image with branch name  ***"
echo "****************************************"

# Docker image name on docker hub
IMAGE_NAME=consensys/fc-retrieval-client-js

VERSION=$1
echo "Client-js version: $VERSION"
IMAGE=$2
echo "Client-js image: $IMAGE"


CLIENT_JS_BRANCH=`git rev-parse --abbrev-ref HEAD`
echo "Client-js repo branch: $CLIENT_JS_BRANCH"

TAG="develop-$CLIENT_JS_BRANCH"
echo "TAG: $TAG"
docker tag $IMAGE $IMAGE_NAME:$TAG
