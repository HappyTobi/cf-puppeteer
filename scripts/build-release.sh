#!/bin/bash

set -e

SANDBOX=$(mktemp -d)

echo "building linux..."
CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o $SANDBOX/cf-puppeteer-linux github.com/happytobi/cf-puppeteer

echo "building os x..."
CGO_ENABLED=0 GOARCH=amd64 GOOS=darwin go build -o $SANDBOX/cf-puppeteer-darwin github.com/happytobi/cf-puppeteer

echo "building windows..."
CGO_ENABLED=0 GOARCH=amd64 GOOS=windows go build -o $SANDBOX/cf-puppeteer.exe github.com/happytobi/cf-puppeteer

echo

find $SANDBOX -type f -exec file {} \;

echo
echo "binaries are in $SANDBOX!"
