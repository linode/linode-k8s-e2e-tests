#!/bin/bash
set -exou pipefail

helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add metrics-server https://kubernetes-sigs.github.io/metrics-server/
