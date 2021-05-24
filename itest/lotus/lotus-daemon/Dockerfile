# Copyright (C) 2020 ConsenSys Software Inc
FROM consensys/lotus-base
COPY daemon/config.toml /root/.lotus/
WORKDIR /app/lotus
ADD start-lotus-daemon.sh .
RUN chmod +x start-lotus-daemon.sh
ENTRYPOINT ./start-lotus-daemon.sh
EXPOSE 1234