#!/bin/bash

var="$(kubectl get deployment coredns --namespace kube-system -o=jsonpath='{$.spec.template.spec.containers[:1].image}' | cut -d : -f 2)"

if [[ $var =~ "1.8.4" ]]; then
    echo "coredns successfully updated by command"
else
    exit 1
fi