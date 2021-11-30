# Linode Kubernetes End-to-end (e2e) tests

This repository contains e2e tests for Linode Kubernetes Engine (LKE)

## How to run these tests

Install the following packages (macOS examples)

```shell
brew install terraform # >= v1.0.0
brew install golang # >= 1.17.0
brew install kubectl
brew install hg
```

Add the following environment variables to your shell rc

```
export LINODE_API_TOKEN=<your linode API token>

export GOPATH=$HOME/go
export PATH=$HOME/go/bin:$PATH
export GO111MODULE=on 
```

If you need a Linode API token visit this page:
https://cloud.linode.com/profile/tokens

Then, `go get` this repo
`go get github.com/linode/linode-k8s-e2e-tests`

That may fail, if it does, navigate to the directory that was created and run `go mod tidy`:

```
cd ~/go/src/github.com/linode/linode-k8s-e2e-tests
go mod tidy
```

By default the tests use $HOME/.ssh/id\_rsa.pub as the public key used to provision the cluster, so it needs to be added to your agent.

```
ssh-add $HOME/.ssh/id_rsa
```

Then, run the tests

```
make test
```
