#!/bin/sh

go clean -testcache
go test -v ./test/...