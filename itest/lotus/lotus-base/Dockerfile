# Copyright (C) 2020 ConsenSys Software Inc
FROM golang:1.15.6 as builder
WORKDIR /app
RUN apt-get update && apt-get install -y ca-certificates make gcc musl-dev git mesa-opencl-icd ocl-icd-opencl-dev bzr jq pkg-config curl clang build-essential hwloc libhwloc-dev tmux
RUN git clone https://github.com/filecoin-project/lotus.git && cd lotus && git checkout v1.9.0
RUN cd /app/lotus/extern && git clone https://github.com/filecoin-project/filecoin-ffi.git && cd filecoin-ffi && git checkout 58ef68d06a44cf9c
WORKDIR /app/lotus
RUN make 2k
ENV LOTUS_SKIP_GENESIS_CHECK=_yes_
RUN ./lotus fetch-params 2048
RUN ./lotus-seed pre-seal --sector-size 2KiB --num-sectors 2
RUN ./lotus-seed genesis new localnet.json
RUN ./lotus-seed genesis add-miner localnet.json ~/.genesis-sectors/pre-seal-t01000.json
ADD setup-lotus-base.sh .
RUN chmod +x setup-lotus-base.sh
RUN ./setup-lotus-base.sh