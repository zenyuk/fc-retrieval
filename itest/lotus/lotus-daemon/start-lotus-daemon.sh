#!/bin/bash

./lotus wallet import --as-default ~/.genesis-sectors/pre-seal-t01000.key
./lotus daemon --lotus-make-genesis=devgen.car --genesis-template=localnet.json --bootstrap=false
