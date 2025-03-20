#!/bin/bash
if [[ "$1" == "" ]]; then
    echo "Usage: $0 <file>"
    exit 1
fi
PORT="${SIMPLESAMPLESERVER_PORT:-8080}"
curl -F "sha256=$(sha256sum "$1" | cut -c 1-64)" -F "file=@\"$1\";filename=\"$1\"" "localhost:${PORT}/upload"
