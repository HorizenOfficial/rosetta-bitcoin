<p align="center">
  <a href="https://www.rosetta-api.org">
    <img width="90%" alt="Rosetta" src="https://www.rosetta-api.org/img/rosetta_header.png">
  </a>
</p>
<h3 align="center">
   Rosetta Zen
</h3>



## Overview
`rosetta-zen` provides an implementation of the Rosetta API for
Hotizen network in Golang. If you haven't heard of the Rosetta API, you can find more
information [here](https://rosetta-api.org).

## Usage
As specified in the [Rosetta API Principles](https://www.rosetta-api.org/docs/automated_deployment.html),
all Rosetta implementations must be deployable via Docker and support running via either an
[`online` or `offline` mode](https://www.rosetta-api.org/docs/node_deployment.html#multiple-modes).

**YOU MUST INSTALL DOCKER FOR THE FOLLOWING INSTRUCTIONS TO WORK. YOU CAN DOWNLOAD
DOCKER [HERE](https://www.docker.com/get-started).**

### Install
Running the following commands will create a Docker image called `rosetta-zen:latest`.

#### From GitHub
To download the pre-built Docker image from the latest release, run:
```text
curl -sSfL https://raw.githubusercontent.com/HorizenOfficial/rosetta-zen/blob/mo/Data_API/install.sh | sh -s
```
_Do not try to install rosetta-zen using GitHub Packages!_

#### From Source
After cloning this repository, run:
```text
make build-local
```

### Run
Running the following commands will start a Docker container in
[detached mode](https://docs.docker.com/engine/reference/run/#detached--d) with
a data directory at `<working directory>/zen-data` and the Rosetta API accessible
at port `8080`.

#### Mainnet:Online
```text
docker run -d --rm --ulimit "nofile=100000:100000" -v "$(pwd)/zen-data:/data" -e "MODE=ONLINE" -e "NETWORK=MAINNET" -e "PORT=8080" -p 8080:8080 -p 8333:8333 rosetta-zen:latest
```
_If you cloned the repository, you can run `make run-mainnet-online`._

#### Mainnet:Offline
```text
docker run -d --rm -e "MODE=OFFLINE" -e "NETWORK=MAINNET" -e "PORT=8081" -p 8081:8081 rosetta-zen:latest
```
_If you cloned the repository, you can run `make run-mainnet-offline`._

#### Testnet:Online
```text
docker run -d --rm --ulimit "nofile=100000:100000" -v "$(pwd)/zen-data:/data" -e "MODE=ONLINE" -e "NETWORK=TESTNET" -e "PORT=8080" -p 8080:8080 -p 18333:18333 rosetta-zen:latest
```
_If you cloned the repository, you can run `make run-testnet-online`._

#### Testnet:Offline
```text
docker run -d --rm -e "MODE=OFFLINE" -e "NETWORK=TESTNET" -e "PORT=8081" -p 8081:8081 rosetta-zen:latest
```
_If you cloned the repository, you can run `make run-testnet-offline`._


### Network Settings
To increase the load `rosetta-zen` can handle, it is recommended to tune your OS
settings to allow for more connections. On a linux-based OS, you can run the following
commands ([source](http://www.tweaked.io/guide/kernel)):
```text
sysctl -w net.ipv4.tcp_tw_reuse=1
sysctl -w net.core.rmem_max=16777216
sysctl -w net.core.wmem_max=16777216
sysctl -w net.ipv4.tcp_max_syn_backlog=10000
sysctl -w net.core.somaxconn=10000
sysctl -p (when done)
```
_We have not tested `rosetta-zen` with `net.ipv4.tcp_tw_recycle` and do not recommend
enabling it._

You should also modify your open file settings to `100000`. This can be done on a linux-based OS
with the command: `ulimit -n 100000`.

### Memory-Mapped Files
`rosetta-bitcoin` uses [memory-mapped files](https://en.wikipedia.org/wiki/Memory-mapped_file) to
persist data in the `indexer`. As a result, you **must** run `rosetta-bitcoin` on a 64-bit
architecture (the virtual address space easily exceeds 100s of GBs) and may need to increase the
allocated size of swap space to avoid a kernel OOM.

## Architecture
`rosetta-zen` uses the `syncer`, `storage`, `parser`, and `server` package
from [`rosetta-sdk-go`](https://github.com/coinbase/rosetta-sdk-go) instead
of a new Zen-specific implementation of packages of similar functionality. Below
you can find a high-level overview of how everything fits together:
```text
                               +------------------------------------------------------------------+
                               |                                                                  |
                               |                 +--------------------------------------+         |
                               |                 |                                      |         |
                               |                 |                 indexer              |         |
                               |                 |                                      |         |
                               |                 | +--------+                           |         |
                               +-------------------+ pruner <----------+                |         |
                               |                 | +--------+          |                |         |
                         +-----v----+            |                     |                |         |
                         |   zend   |            |              +------+--------+       |         |
                         +-----+----+            |     +--------> block_storage <----+  |         |
                               |                 |     |        +---------------+    |  |         |
                               |                 | +---+----+                        |  |         |
                               +-------------------> syncer |                        |  |         |
                                                 | +---+----+                        |  |         |
                                                 |     |        +--------------+     |  |         |
                                                 |     +--------> coin_storage |     |  |         |
                                                 |              +------^-------+     |  |         |
                                                 |                     |             |  |         |
                                                 +--------------------------------------+         |
                                                                       |             |            |
+-------------------------------------------------------------------------------------------+     |
|                                                                      |             |      |     |
|         +------------------------------------------------------------+             |      |     |
|         |                                                                          |      |     |
|         |                     +---------------------+-----------------------+------+      |     |
|         |                     |                     |                       |             |     |
| +-------+---------+   +-------+---------+   +-------+-------+   +-----------+----------+  |     |
| | account_service |   | network_service |   | block_service |   | construction_service +--------+
| +-----------------+   +-----------------+   +---------------+   +----------------------+  |
|                                                                                           |
|                                         server                                            |
|                                                                                           |
+-------------------------------------------------------------------------------------------+
```

### Optimizations
* Automatically prune bitcoind while indexing blocks
* Reduce sync time with concurrent block indexing
* Use [Zstandard compression](https://github.com/facebook/zstd) to reduce the size of data stored on disk
without needing to write a manual byte-level encoding

#### Concurrent Block Syncing
To speed up indexing, `rosetta-zen` uses concurrent block processing
with a "wait free" design (using channels instead of sleeps to signal
which threads are unblocked). This allows `rosetta-zen` to fetch
multiple inputs from disk while it waits for inputs that appeared
in recently processed blocks to save to disk.
```text
                                                   +----------+
                                                   |   zend   |
                                                   +-----+----+
                                                         |
                                                         |
          +---------+ fetch block data / unpopulated txs |
          | block 1 <------------------------------------+
          +---------+                                    |
       +-->   tx 1  |                                    |
       |  +---------+                                    |
       |  |   tx 2  |                                    |
       |  +----+----+                                    |
       |       |                                         |
       |       |           +---------+                   |
       |       |           | block 2 <-------------------+
       |       |           +---------+                   |
       |       +----------->   tx 3  +--+                |
       |                   +---------+  |                |
       +------------------->   tx 4  |  |                |
       |                   +---------+  |                |
       |                                |                |
       | retrieve previously synced     |   +---------+  |
       | inputs needed for future       |   | block 3 <--+
       | blocks while waiting for       |   +---------+
       | populated blocks to save to    +--->   tx 5  |
       | disk                               +---------+
       +------------------------------------>   tx 6  |
       |                                    +---------+
       |
       |
+------+--------+
|  coin_storage |
+---------------+
```

## Testing with rosetta-cli
To validate `rosetta-zen`, [install `rosetta-cli`](https://github.com/coinbase/rosetta-cli#install)
and run one of the following commands:
* `rosetta-cli check:data --configuration-file rosetta-cli-conf/zen_testnet.json`
* `rosetta-cli check:construction --configuration-file rosetta-cli-conf/zen_testnet.json`
* `rosetta-cli check:data --configuration-file rosetta-cli-conf/zen_mainnet.json`
* `rosetta-cli check:construction --configuration-file rosetta-cli-conf/zen_mainnet.json`

## Development
* `make deps` to install dependencies
* `make test` to run tests
* `make lint` to lint the source code
* `make salus` to check for security concerns
* `make build-local` to build a Docker image from the local context
* `make coverage-local` to generate a coverage report

## License
This project is available open source under the terms of the [Apache 2.0 License](https://opensource.org/licenses/Apache-2.0).

© 2020 Coinbase
© 2020 Zen Blockchain Foundation
