#!/usr/bin/env bash

CLUSTERCONF=$1

kubectl delete serviceaccount tiller --kubeconfig=${CLUSTERCONF}
kubectl delete clusterrolebinding tiller-cluster-rule --kubeconfig=${CLUSTERCONF}

helm reset --force --kubeconfig=${CLUSTERCONF}
