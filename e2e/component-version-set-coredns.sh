#!/bin/bash

set -e

./k8sclusterupgradetool component version set \
-c=kind-k8s-cluster-upgrade-tool-test-cluster \
-o=coredns \
-v=1.8.4
