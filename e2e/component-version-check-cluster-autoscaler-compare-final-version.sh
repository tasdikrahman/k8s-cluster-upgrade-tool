#!/bin/bash

var="$(kubectl get deployment cluster-autoscaler --namespace kube-system -o=jsonpath='{$.spec.template.spec.containers[:1].image}' | cut -d : -f 2)"

if [[ $var =~ "v1.20.1" ]]; then
    echo "Cluster autoscaler successfully updated by command"
else
    exit 1
fi