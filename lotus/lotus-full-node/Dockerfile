# Copyright (C) 2020 ConsenSys Software Inc
FROM consensys/lotus-base
COPY daemon/config.toml /root/.lotus/
COPY miner/config.toml /root/.lotusminer/
WORKDIR /app/lotus
ADD config-lotus-full-node.sh .
ADD start-lotus-full-node.sh .
RUN chmod +x config-lotus-full-node.sh
RUN chmod +x start-lotus-full-node.sh
RUN ./config-lotus-full-node.sh