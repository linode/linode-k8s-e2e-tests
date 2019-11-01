#!/bin/bash
set -exou pipefail

CLUSTERCONF=$1
kubectl create serviceaccount --namespace kube-system tiller --kubeconfig=${CLUSTERCONF} || true
kubectl create clusterrolebinding tiller-cluster-rule --clusterrole=cluster-admin --serviceaccount=kube-system:tiller --kubeconfig=${CLUSTERCONF} || true

helm init --service-account=tiller --kubeconfig=${CLUSTERCONF} || true
