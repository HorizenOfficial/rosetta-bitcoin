#!/bin/bash

set -eEuox pipefail

ZEND_VERISON_OLD="${ZEND_VERISON_OLD:-4.1.0-rc1}"
ZEND_VERSION_NEW="${ZEND_VERSION_NEW:-4.1.0}"

bsd_sed=""
[ "$(uname -s)" = "Darwin" ] && bsd_sed=" ''"

# bump zend version
sed -i${bsd_sed} "s/ZEN_COMMITTISH=v${ZEND_VERISON_OLD}/ZEN_COMMITTISH=v${ZEND_VERSION_NEW}/g" .travis.yml Dockerfile
sed -i${bsd_sed} "s/ZEN_COMMITTISH?=v${ZEND_VERISON_OLD}/ZEN_COMMITTISH?=v${ZEND_VERSION_NEW}/g" Makefile
sed -i${bsd_sed} "s/\"${ZEND_VERISON_OLD}\"/\"${ZEND_VERSION_NEW}\"/g" services/network_service_test.go services/types.go

