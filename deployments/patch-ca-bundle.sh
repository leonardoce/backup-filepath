#!/bin/bash
# File originally from https://github.com/banzaicloud/admission-webhook-example/blob/blog/deployment/webhook-patch-ca-bundle.sh

ROOT=$(cd $(dirname $0)/../../; pwd)

set -o errexit
set -o nounset
set -o pipefail

current_context=$(kubectl config view --raw -o json | jq '."current-context"')
export CA_BUNDLE=$(kubectl config view --raw -o json | jq -r "(.clusters[] | select(.name == ${current_context})).cluster.\"certificate-authority-data\"")

if command -v envsubst >/dev/null 2>&1; then
    envsubst
else
    sed -e "s|\${CA_BUNDLE}|${CA_BUNDLE}|g"
fi
