#!/bin/zsh

set -e

export GOOS=
export GOARCH=

cd $(eval git rev-parse --show-toplevel)
pwd
go generate ./...
go run internal/pages/main.go --local=false
