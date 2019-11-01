#!/usr/bin/env bash

CLUSTERCONF=$1

kubectl delete serviceaccount tiller --kubeconfig=${CLUSTERCONF} || true
kubectl delete clusterrolebinding tiller-cluster-rule --kubeconfig=${CLUSTERCONF} || true

helm reset --force --kubeconfig=${CLUSTERCONF}
