#!/bin/bash

set -e

./k8sclusterupgradetool component version set \
-c=kind-k8s-cluster-upgrade-tool-test-cluster \
-o=kube-proxy \
-v=v1.20.14
