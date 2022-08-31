#!/bin/bash

set -e

mkdir ~/.k8sclusterupgradetool/
cd ~/.k8sclusterupgradetool/
touch ~/.k8sclusterupgradetool/config.yaml

cat > ~/.k8sclusterupgradetool/config.yaml <<EOF
---
# the values below are just one above from the test values being installed
components:
  aws-node: "v1.11.1"
  cluster-autoscaler: "v1.20.1"
  coredns: "1.8.4"
  kube-proxy: "v1.20.14"
clusterlist:
- ClusterName: "kind-k8s-cluster-upgrade-tool-test-cluster"
  AwsRegion: "region1"
  AwsAccount: "account1"
  AwsNodeObject:
    ObjectType: "daemonset"
    DeploymentName: "aws-node"
    ContainerName: "aws-node"
    Namespace: "kube-system"
  ClusterAutoscalerObject:
    ObjectType: "deployment"
    DeploymentName: "cluster-autoscaler"
    ContainerName: "aws-cluster-autoscaler"
    Namespace: "kube-system"
  CoreDnsObject:
    DeploymentName: "coredns"
    ObjectType: "deployment"
    ContainerName: "coredns"
    Namespace: "kube-system"
  KubeProxyObject:
    DeploymentName: "kube-proxy"
    ObjectType: "daemonset"
    ContainerName: "kube-proxy"
    Namespace: "kube-system"
EOF
