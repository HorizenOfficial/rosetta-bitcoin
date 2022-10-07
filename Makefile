.PHONY: deps build run lint mocks run-mainnet-online run-mainnet-offline run-testnet-online \
	run-testnet-offline check-comments add-license check-license shorten-lines test \
	coverage spellcheck salus build-local coverage-local format check-format

ADDLICENSE_INSTALL=go install github.com/google/addlicense@latest
ADDLICENSE_CMD=addlicense
ADDLICENSE_IGNORE=-ignore ".github/**/*" -ignore ".idea/**/*"
ADDLICENCE_SCRIPT=${ADDLICENSE_CMD} -c "Coinbase, Inc." -l "apache" -v ${ADDLICENSE_IGNORE}
SPELLCHECK_CMD=go run github.com/client9/misspell/cmd/misspell
GOLINES_INSTALL=go install github.com/segmentio/golines@latest
GOLINES_CMD=golines
GOLINT_INSTALL=go get golang.org/x/lint/golint
GOLINT_CMD=golint
GOVERALLS_INSTALL=go install github.com/mattn/goveralls@latest
GOVERALLS_CMD=goveralls
GOIMPORTS_CMD=go run golang.org/x/tools/cmd/goimports
GO_PACKAGES=./services/... ./indexer/... ./zen/... ./zend/... ./zenutil/... ./configuration/...
GO_FOLDERS=$(shell echo ${GO_PACKAGES} | sed -e "s/\.\///g" | sed -e "s/\/\.\.\.//g")
TEST_SCRIPT=go test ${GO_PACKAGES} -buildmode=pie
LINT_SETTINGS=golint,misspell,gocyclo,gocritic,whitespace,goconst,gocognit,bodyclose,unconvert,lll,unparam
PWD=$(shell pwd)
GZIP_CMD=$(shell command -v pigz || echo gzip)
NOFILE=100000
ZEND_VERSION=v3.2.0

deps:
	go get ./...

build:
	docker build --pull --build-arg "ZEN_COMMITTISH=${ZEND_VERSION}" -t rosetta-zen:latest https://github.com/HorizenOfficial/rosetta-zen

build-local:
	docker build --pull --build-arg "ZEN_COMMITTISH=${ZEND_VERSION}" -t rosetta-zen:latest .

build-release:
	# make sure to always set version with vX.X.X
	docker build --pull --no-cache --build-arg IS_RELEASE=true  --build-arg "ZEN_COMMITTISH=${ZEND_VERSION}" -t rosetta-zen:$(version) .;
	docker save rosetta-zen:$(version) | ${GZIP_CMD} > rosetta-zen-$(version).tar.gz;

run-mainnet-online:
	docker container rm rosetta-zen-mainnet-online || true
	docker run --rm -v "${PWD}/zen-data:/data" ubuntu:20.04 bash -c 'chown -R nobody:nogroup /data';
	docker run -d --name=rosetta-zen-mainnet-online --ulimit "nofile=${NOFILE}:${NOFILE}" -v "${PWD}/zen-data:/data" -e "MODE=ONLINE" -e "NETWORK=MAINNET" -e "PORT=8080" -p 8080:8080 -p 9033:9033 rosetta-zen:latest;

run-mainnet-offline:
	docker container rm rosetta-zen-mainnet-offline || true
	docker run -d --name=rosetta-zen-mainnet-offline -e "MODE=OFFLINE" -e "NETWORK=MAINNET" -e "PORT=8081" -p 8081:8081 rosetta-zen:latest

run-testnet-online:
	docker container rm rosetta-zen-testnet-online || true
	docker run --rm -v "${PWD}/zen-data:/data" ubuntu:20.04 bash -c 'chown -R nobody:nogroup /data';
	docker run -d --name=rosetta-zen-testnet-online --ulimit "nofile=${NOFILE}:${NOFILE}" -v "${PWD}/zen-data:/data" -e "MODE=ONLINE" -e "NETWORK=TESTNET" -e "PORT=8080" -p 8080:8080 -p 19033:19033 -p 18231:18231 rosetta-zen:latest;

run-testnet-offline:
	docker container rm rosetta-zen-testnet-offline || true
	docker run -d --name=rosetta-zen-testnet-offline -e "MODE=OFFLINE" -e "NETWORK=TESTNET" -e "PORT=8081" -p 8081:8081 rosetta-zen:latest

stop-mainnet-online:
	docker container stop rosetta-zen-mainnet-online

stop-mainnet-offline:
	docker container stop rosetta-zen-mainnet-offline

stop-testnet-online:
	docker container stop rosetta-zen-testnet-online

stop-testnet-offline:
	docker container stop rosetta-zen-testnet-offline

train:
	./zstd-train.sh $(network) transaction $(data-directory)

check-comments:
	${GOLINT_INSTALL}
	${GOLINT_CMD} -set_exit_status ${GO_FOLDERS} .
	go mod tidy

lint: | check-comments
	golangci-lint run --timeout 2m0s -v -E ${LINT_SETTINGS},gomnd

add-license:
	${ADDLICENSE_INSTALL}
	${ADDLICENCE_SCRIPT} .

check-license:
	${ADDLICENSE_INSTALL}
	${ADDLICENCE_SCRIPT} -check .

shorten-lines:
	${GOLINES_INSTALL}
	${GOLINES_CMD} -w --shorten-comments ${GO_FOLDERS} .

format:
	gofmt -s -w -l .
	${GOIMPORTS_CMD} -w .

check-format:
	! gofmt -s -l . | read
	! ${GOIMPORTS_CMD} -l . | read

test:
	${TEST_SCRIPT}

coverage:
	${GOVERALLS_INSTALL}
	if [ "${COVERALLS_TOKEN}" ]; then ${TEST_SCRIPT} -coverprofile=c.out -covermode=count; ${GOVERALLS_CMD} -coverprofile=c.out -repotoken ${COVERALLS_TOKEN}; fi

coverage-local:
	${TEST_SCRIPT} -cover

salus:
	docker run --rm -t -v ${PWD}:/home/repo coinbase/salus

spellcheck:
	${SPELLCHECK_CMD} -error .

mocks:
	rm -rf mocks;
	mockery --dir indexer --all --case underscore --outpkg indexer --output mocks/indexer;
	mockery --dir services --all --case underscore --outpkg services --output mocks/services;
	${ADDLICENSE_INSTALL}
	${ADDLICENCE_SCRIPT} .;
