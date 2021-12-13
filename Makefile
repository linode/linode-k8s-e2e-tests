export GO111MODULE=on

$(GOPATH)/bin/goimports:
	GO111MODULE=off go get golang.org/x/tools/cmd/goimports

$(GOPATH)/bin/ginkgo:
	GO111MODULE=off go get -u github.com/onsi/ginkgo/ginkgo

vet:
	go vet -composites=false ./

fmt: vet $(GOPATH)/bin/goimports
	# goimports runs a gofmt
	goimports -w *.go

.PHONY: testi[
test: $(GOPATH)/bin/ginkgo
	@if [ -z "${LINODE_API_TOKEN}" ]; then\
		echo "Skipping Test, LINODE_API_TOKEN is not set";\
	else \
		go list -m; \
		ginkgo -r --v --progress --trace --cover -- --v=3; \
	fi

test-existing: $(GOPATH)/bin/ginkgo
	@if [ -z "${LINODE_API_TOKEN}" ]; then\
		echo "Skipping Test, LINODE_API_TOKEN is not set";\
	else \
		go list -m; \
		ginkgo -r --v --progress --trace --cover -- --use-existing --kubeconfig="${TEST_KUBECONFIG}" --v=3; \
	fi

install-terraform:
	sudo apt-get install wget unzip
	wget https://releases.hashicorp.com/terraform/0.11.13/terraform_0.11.13_linux_amd64.zip
	unzip terraform_1.0.3_linux_amd64.zip
	sudo mv terraform /usr/local/bin/
