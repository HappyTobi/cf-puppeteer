#!/bin/bash

echo 'Building new `cf-puppeteer` binary...'
go install

echo 'Installing the plugin...'
cf uninstall-plugin cf-puppeteer
cf install-plugin $GOPATH/bin/cf-puppeteer
