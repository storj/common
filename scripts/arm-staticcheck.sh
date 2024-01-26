#!/usr/bin/env bash
set -x
export GOOS=linux
export GOARCH=arm
find -name "*.go" | awk -F / '{print $2}' | sort | uniq | grep -v -E '(http|sqlite)' | xargs -IDIR staticcheck storj.io/common/DIR
