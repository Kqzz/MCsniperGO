#!/bin/sh
goreleaser check
git tag -a $1 -m "release $1"
git push origin $1