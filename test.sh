#!/bin/bash
set -o allexport
source .env
go test  -coverpkg=./... -cover  -covermode=atomic  -coverprofile=coverage.out  ./tests/...
