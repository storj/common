#!/usr/bin/env bash
find -name "*.go" | awk -F / '{print $2}' | sort | uniq | grep -v -E '(debug|process|metrics|db)' | xargs -IDIR check-cross-compile -compiler "go,go.min" storj.io/common/DIR/...
