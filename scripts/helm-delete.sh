#!/usr/bin/env bash

CLUSTER=$1
WD=$(pwd)

kubectl delete serviceaccount tiller --kubeconfig=${CLUSTER}".conf" || true
kubectl delete clusterrolebinding tiller-cluster-rule --kubeconfig=${CLUSTER}".conf" || true
helm reset --force --kubeconfig=${WD}"/"${CLUSTER}".conf"