#!/bin/bash

set -e

SANDBOX=$(mktemp -d)

echo "building os x..."
CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -o $SANDBOX/cf-puppeteer-darwin github.com/happytobi/cf-puppeteer

echo
echo "binaries are in $SANDBOX"

cf install-plugin -f $SANDBOX/cf-puppeteer-darwin
