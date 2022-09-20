#!/bin/bash

set -e

./k8sclusterupgradetool component version set \
-c=kind-k8s-cluster-upgrade-tool-test-cluster \
-o=aws-node \
-v=v1.11.1
