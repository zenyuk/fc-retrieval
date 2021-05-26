# Copyright (C) 2020 ConsenSys Software Inc
FROM consensys/lotus-base
COPY daemon/config.toml /root/.lotus/
WORKDIR /app/lotus
ENTRYPOINT ./lotus daemon --lotus-make-genesis=devgen.car --genesis-template=localnet.json --bootstrap=false
EXPOSE 1234