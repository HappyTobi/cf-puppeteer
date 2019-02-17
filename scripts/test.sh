#!/bin/bash
# vim: set ft=sh

set -e

cd $GOPATH/src/github.com/happytobi/cf-puppeteer

govendor install github.com/happytobi/cf-puppeteer/vendor/github.com/onsi/ginkgo/ginkgo

ginkgo -r "$@"
