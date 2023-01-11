#!/bin/bash -e

# This script is used to build a test environment for local development.

CLUSTER_NAME="atmosphere"
KUBECTL="kubectl --context kind-${CLUSTER_NAME}"

# Create a kind cluster if one doesn't already exist
if ! kind get clusters | grep -q "${CLUSTER_NAME}"; then
  kind create cluster --name ${CLUSTER_NAME} --config hack/kind-config.yml
fi

# Install the operators which we depend on
# TODO(mnaser): Use OLM for this
${KUBECTL} apply --server-side -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml
${KUBECTL} apply --server-side -f https://github.com/prometheus-operator/prometheus-operator/raw/v0.62.0/bundle.yaml
${KUBECTL} apply --server-side -f https://raw.githubusercontent.com/percona/percona-xtradb-cluster-operator/v1.12.0/deploy/bundle.yaml
${KUBECTL} apply --server-side -f https://github.com/rabbitmq/cluster-operator/releases/download/v1.13.1/cluster-operator.yml

# Install the CRDs + charts
make charts
make install

# Install the basic dependencies that Atmosphere resources need
${KUBECTL} apply --server-side -f hack/testdata/pxc.yml

# Install a set of basic resources for a deployment
cat config/samples/{infra,openstack}_*.yaml | ${KUBECTL} apply -f -
