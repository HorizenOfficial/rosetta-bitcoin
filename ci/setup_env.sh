#!/bin/bash

set -eo pipefail

if [ ! -z "${TRAVIS_TAG}" ]; then
  GNUPGHOME="$(mktemp -d 2>/dev/null || mktemp -d -t 'GNUPGHOME')"
  export GNUPGHOME
  echo "Tagged build, fetching maintainer keys."
  read -a keys <<< "$ROSETTA_MAINTAINER_KEYS"
  for key in "${keys[@]}"; do
    gpg -v --batch --keyserver keyserver.ubuntu.com --recv "$key" ||
    gpg -v --batch --keyserver hkps://keys.openpgp.org --recv "$key" ||
    gpg -v --batch --keyserver hkp://p80.pool.sks-keyservers.net:80 --recv-keys "$key" ||
    gpg -v --batch --keyserver hkp://ipv4.pool.sks-keyservers.net --recv-keys "$key" ||
    gpg -v --batch --keyserver hkp://pgp.mit.edu:80 --recv-keys "$key"
  done
  if git verify-tag -v "${TRAVIS_TAG}"; then
    echo "Valid signed tag"
    export version="${TRAVIS_TAG}"
  fi
fi
