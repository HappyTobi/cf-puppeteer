#!/bin/bash

echo 'Building new `cf-puppeteer` binary...'
go install

echo 'Installing the plugin...'
cf uninstall-plugin Cf-puppeteer
cf install-plugin $GOPATH/bin/cf-puppeteer
