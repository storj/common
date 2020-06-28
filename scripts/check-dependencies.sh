#!/usr/bin/env bash
set -ueo pipefail
set +x

# This script verifies that we don't accidentally import specific packages.

if go list -deps -test ./... | grep -Eq "github.com/(lib/pq|jackc/pg)"; then
    echo "common must not have a dependency to postgres";
    exit -1;
fi

if go list -deps -test ./... | grep -q "redis"; then
    echo "common must not have a dependency to redis";
    exit -1;
fi

if go list -deps -test ./... | grep -q "bolt"; then
    echo "common must not have a dependency to bolt";
    exit -1;
fi

if go list -deps $(go list ./... | grep -v "test") | grep -q "testing"; then
    echo "common must not have a dependency to testing";
    exit -1;
fi
