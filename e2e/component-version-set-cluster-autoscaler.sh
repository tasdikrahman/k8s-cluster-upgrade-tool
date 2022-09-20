#!/bin/bash

set -e

./k8sclusterupgradetool component version set \
-c=kind-k8s-cluster-upgrade-tool-test-cluster \
-o=cluster-autoscaler \
-v=v1.20.1
