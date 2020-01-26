#!/usr/bin/env bash

set +eu

pushd src
go build
popd
mv src/simple-sync .