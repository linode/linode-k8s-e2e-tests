module github.com/linode/linode-k8s-e2e-tests

go 1.12

require (
	github.com/codeskyblue/go-sh v0.0.0-20190412065543-76bd3d59ff27
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/onsi/ginkgo v1.10.3
	github.com/onsi/gomega v1.7.0
	github.com/pkg/errors v0.8.1
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45 // indirect
	google.golang.org/appengine v1.6.1 // indirect
	k8s.io/api v0.0.0-20190503110853-61630f889b3c
	k8s.io/apimachinery v0.0.0-20191030190112-bb31b70367b7
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/utils v0.0.0-20190506122338-8fab8cb257d5 // indirect
	kmodules.xyz/client-go v0.0.0-20191023042933-b12d1ccfaf57
)
