#!/bin/sh
sudo docker rm -f testing
sudo make build-local
sudo docker run -d --name=testing --ulimit "nofile=100000:100000" -v "$(pwd)/zen-data:/data" -e "MODE=ONLINE" -e "NETWORK=REGTEST" -e "PORT=8080" -p 8080:8080 -p 18333:18333 rosetta-zen:latest sleep 9999
sudo docker exec -it testing bash

