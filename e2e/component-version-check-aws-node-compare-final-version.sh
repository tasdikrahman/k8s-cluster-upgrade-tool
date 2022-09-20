#!/bin/bash

var="$(kubectl get daemonset aws-node --namespace kube-system -o=jsonpath='{$.spec.template.spec.containers[:1].image}' | cut -d : -f 2)"

if [[ $var =~ "v1.11.1" ]]; then
    echo "aws-node successfully updated by command"
else
    exit 1
fi