#!/usr/bin/env bash
set -uo pipefail
set +x

# This script verifies that we don't accidentally import specific packages.

if ! check-dependency -ignore "cockroachutil;dbutil" -check "github.com/(lib/pq|jackc/pg)" ./...; then
    echo "we shouldn't import postgres outside of database packages";
    exit -1;
fi

if ! check-dependency -check "redis;bolt" ./...; then
    echo "common must not have a dependency to redis or bolt";
    exit -1;
fi

if ! check-dependency -ignore "test;common/process" -check "test" -except "internal/testlog" ./...; then
    echo "packages not related to test functionality should not bring in testing related things";
    exit -1;
fi

if ! check-dependency -ignore "accesslogs;cfgstruct;dbutil;debug;httpranger;metrics;migrate;process;requestid;storj/location/gen;tagsql;test;traces;version" -check "net/http" ./...; then
    echo "net/http is a huge dependency that we don't want to accidentally introduce into uplink";
    exit -1;
fi
