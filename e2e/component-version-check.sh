#!/bin/bash

set -e

./k8sclusterupgradetool component version check -c=kind-k8s-cluster-upgrade-tool-test-cluster
