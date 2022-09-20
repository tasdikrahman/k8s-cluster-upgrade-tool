#!/bin/bash

var="$(kubectl get daemonset kube-proxy --namespace kube-system -o=jsonpath='{$.spec.template.spec.containers[:1].image}' | cut -d : -f 2)"

if [[ $var =~ "v1.20.14" ]]; then
    echo "kube-proxy successfully updated by command"
else
    exit 1
fi