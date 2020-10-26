#!/bin/bash

set -euo pipefail

if [ "$(stat -c '%u:%g' /data)" != "65534:65534" ]; then
  echo "Folder mounted at /data is not owned by nobody:nogroup, please change it's permissions on the host with 'sudo chown -R 65534:65534 path/to/zen-data'."
  exit 1
fi

/app/fetch-params.sh

exec "$@"
