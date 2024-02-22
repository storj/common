#!/usr/bin/env bash
set -ueo pipefail
set +x

# This script verifies that we don't accidentally import specific packages.



for package in $(find -name "*.go" | xargs -n1 dirname | sort | uniq | \
     grep -v test | awk '{print $1}' | \
     grep -v -E "cockroachutil|dbutil"); do
    if go list -deps $package |  grep -Eq "github.com/(lib/pq|jackc/pg)"; then
        echo "postgres dependencies are allowed only in white-listed packages, but not in $package";
        exit -1;
    fi
done

if go list -deps -test ./... | grep -q "redis"; then
    echo "common must not have a dependency to redis";
    exit -1;
fi

if go list -deps -test ./... | grep -q "bolt"; then
    echo "common must not have a dependency to bolt";
    exit -1;
fi

for package in $(find -name "*.go" | xargs -n1 dirname | sort | uniq | \
     grep -v test | awk '{print $1}' | \
     grep -v -E "process|scripts"); do
    if go list -deps $package | grep -q "^test$"; then
        echo "test dependency is allowed only in white-listed packages, but not in $package";
        exit -1;
    fi
done

for package in $(find -name "*.go" | xargs -n1 dirname | sort | uniq | \
     grep -v test | awk '{print $1}' | \
     grep -v -E "httpranger|requestid|test|gen|cfgstruct|dbutil|debug|identity|metrics|migrate|process|tagsql|traces|version"); do
    if go list -deps $package | grep -q "net/http"; then
        echo "net/http dependency is allowed only in white-listed packages, but not in $package";
        exit -1;
    fi
done
