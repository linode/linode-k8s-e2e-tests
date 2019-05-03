#!/bin/bash
set -exou pipefail

CLUSTER=$1
WD=$(pwd)
kubectl create serviceaccount --namespace kube-system tiller --kubeconfig=${CLUSTER}".conf" || true
kubectl create clusterrolebinding tiller-cluster-rule --clusterrole=cluster-admin --serviceaccount=kube-system:tiller --kubeconfig=${CLUSTER}".conf" || true

helm init --service-account=tiller --kubeconfig=${WD}"/"${CLUSTER}".conf" || true