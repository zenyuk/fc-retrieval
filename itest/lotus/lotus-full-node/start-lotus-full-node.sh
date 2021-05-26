#!/bin/bash

set -m
./lotus daemon --lotus-make-genesis=devgen.car --genesis-template=localnet.json --bootstrap=false &
sleep 1
./lotus-miner run --nosync