# Copyright 2020 Coinbase, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Build bitcoind
FROM ubuntu:18.04 as bitcoind-builder

RUN mkdir -p /app \
  && chown -R nobody:nogroup /app
WORKDIR /app

# Source: https://github.com/bitcoin/bitcoin/blob/master/doc/build-unix.md#ubuntu--debian
RUN apt-get update && apt-get install -y make gcc g++ autoconf autotools-dev bsdmainutils build-essential git libboost-all-dev \
  libcurl4-openssl-dev libdb++-dev libevent-dev libssl-dev libtool pkg-config python python-pip libzmq3-dev wget



RUN apt-get -qqy update \
  && apt-get -qqy --no-install-recommends install \
    bzip2 \
    ca-certificates \
    sudo \
    unzip \
    wget \
  && rm -rf /var/lib/apt/lists/* /var/cache/apt/*


RUN apt-get -qqy update \
  && apt-get -qqy --no-install-recommends install \
    curl \
    bc \
    vim \
    git \
    binutils \
    xz-utils \
    python3-pip \
    build-essential \
    pkg-config \
    libc6-dev \
    m4 \
    g++-multilib \
    autoconf \
    libtool \
    ncurses-dev \
    zlib1g-dev \
    bsdmainutils \
    automake \
    ssh-client \
  && rm -rf /var/lib/apt/lists/* \
  && apt-get -qyy clean




# VERSION: Bitcoin Core 0.20.1
RUN git clone https://github.com/HorizenOfficial/zen.git

RUN cd zen \
  && zcutil/build.sh -j2

RUN mv zen/src/zend /app/zend \
  && rm -rf zen

# Build Rosetta Server Components
FROM ubuntu:18.04 as rosetta-builder

RUN mkdir -p /app \
  && chown -R nobody:nogroup /app
WORKDIR /app

RUN apt-get update && apt-get install -y curl make gcc g++ git
ENV GOLANG_VERSION 1.15.2
ENV GOLANG_DOWNLOAD_SHA256 b49fda1ca29a1946d6bb2a5a6982cf07ccd2aba849289508ee0f9918f6bb4552
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz

RUN curl -fsSL "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
  && echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c - \
  && tar -C /usr/local -xzf golang.tar.gz \
  && rm golang.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
# Use native remote build context to build in any directory
COPY . src
RUN cd src \
  && go build \
  && cd .. \
  && mv src/rosetta-zen /app/rosetta-zen \
  && mv src/assets/* /app \
  && rm -rf src

## Build Final Image
FROM ubuntu:18.04

RUN apt-get update && \
  apt-get install --no-install-recommends -y libevent-dev libboost-system-dev libboost-filesystem-dev libboost-test-dev libboost-thread-dev net-tools vim git && \
  apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN mkdir -p /app \
  && chown -R nobody:nogroup /app \
  && mkdir -p /data \
  && chown -R nobody:nogroup /data \
  && mkdir -p /data/.zcash-params \
  && chown -R nobody:nogroup /data/.zcash-params

WORKDIR /app

# Copy binary from bitcoind-builder
#COPY --from=bitcoind-builder /app/zend /app/zend

# Copy binary from rosetta-builder
COPY --from=rosetta-builder /app/* /app/
RUN cd /root \
  && ln -sf /data/.zcash-params .zcash-params

# Set permissions for everything added to /app
RUN chmod -R 755 /app/*

CMD ["/app/rosetta-zen"]
