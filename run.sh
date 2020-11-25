#!/bin/bash
docker build -t go-ping .
docker run -it --rm --name go-ping-runtime go-ping
