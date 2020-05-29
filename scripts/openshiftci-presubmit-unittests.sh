#!/bin/sh

# fail if some commands fails
set -e
# show commands
set -x

export ARTIFACTS_DIR="/tmp/artifacts"
export CUSTOM_HOMEDIR=$ARTIFACTS_DIR
export PATH=$PATH:$GOPATH/bin
# set location for golangci-lint cache
# otherwise /.cache is used, and it fails on permission denied
export GOLANGCI_LINT_CACHE="/tmp/.cache"

echo $PULL_SECRET_PATH
echo "==============="
echo $JOB_SPEC
echo "==============="
pr_num="$( jq .refs.pulls[0].num <<<"${JOB_SPEC}" )"
echo $pr_num
echo "==============="
echo $JOB_SPEC

make goget-tools
make validate
make test

# crosscompile and publish artifacts
make cross
cp -r dist $ARTIFACTS_DIR
