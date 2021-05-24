# Copyright (C) 2020 ConsenSys Software Inc
FROM consensys/lotus-base
COPY daemon/config.toml /root/.lotus/
COPY miner/config.toml /root/.lotusminer/
WORKDIR /app/lotus
ADD start-lotus-full-node.sh .
RUN chmod +x start-lotus-full-node.sh
ENTRYPOINT ./start-lotus-full-node.sh
EXPOSE 1234 2345