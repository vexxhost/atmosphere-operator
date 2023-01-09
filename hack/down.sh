#!/bin/bash -e

# This script is used to teardown a test environment for local development.

CLUSTER_NAME="atmosphere"

# Delete the kind cluster
kind delete cluster --name ${CLUSTER_NAME}
