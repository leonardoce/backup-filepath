#!/usr/bin/env bash
set -eu

cd "$(dirname "$0")/.." || exit

if [ -f .env ]; then
    source .env
fi

# Deploy this resource in K8s
current_context=$(kubectl config view --raw -o json | jq -r '."current-context"' | sed "s/kind-//")
kind load docker-image --name=${current_context} volume_injector:${VERSION:-latest}

kubectl apply -f deployments/namespace.yaml

if [ "$(kubectl get secret -n volumeinjector -o name)" != "secret/volumeinjector-secret" ]; then
    ./deployments/create-cert.sh --service volumeinjector-service --secret volumeinjector-secret --namespace volumeinjector
fi
if [ ! -f "./deployments/mutatingwebhook.yaml" ]; then
    cat ./deployments/mutatingwebhook-template.yaml | ./deployments/patch-ca-bundle.sh > ./deployments/mutatingwebhook.yaml
    kubectl apply -f deployments/mutatingwebhook.yaml
fi

kubectl apply -f deployments/deployment.yaml
kubectl apply -f deployments/service.yaml
