#!/bin/bash

./lotus wallet import --as-default ~/.genesis-sectors/pre-seal-t01000.key
./lotus daemon --lotus-make-genesis=devgen.car --genesis-template=localnet.json --bootstrap=false &
sleep 1
./lotus-miner init --genesis-miner --actor=t01000 --sector-size=2KiB --pre-sealed-sectors=~/.genesis-sectors --pre-sealed-metadata=~/.genesis-sectors/pre-seal-t01000.json --nosync
sleep 30 # waiting for init to be complete
./lotus-miner run --nosync