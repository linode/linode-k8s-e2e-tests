module github.com/linode/linode-k8s-e2e-tests

go 1.17

replace (
	k8s.io/api => k8s.io/api v0.0.0-20190918155943-95b840bb6a1f
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190913080033-27d36303b655
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190918160344-1fbdaa4c8d90
)

require (
	github.com/codeskyblue/go-sh v0.0.0-20190412065543-76bd3d59ff27
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/onsi/ginkgo v1.10.3
	github.com/onsi/gomega v1.7.0
	github.com/pkg/errors v0.8.1
	k8s.io/api v0.18.3
	k8s.io/apimachinery v0.18.3
	k8s.io/client-go v11.0.0+incompatible
	kmodules.xyz/client-go v0.0.0-20200818171030-24b2ce405feb
)

require (
	github.com/codegangsta/inject v0.0.0-20150114235600-33e0aa1cb7c0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/docker/spdystream v0.0.0-20160310174837-449fdfce4d96 // indirect
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/protobuf v1.3.2 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/googleapis/gnostic v0.1.0 // indirect
	github.com/hpcloud/tail v1.0.0 // indirect
	github.com/imdario/mergo v0.3.5 // indirect
	github.com/json-iterator/go v1.1.8 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/crypto v0.0.0-20200220183623-bac4c82f6975 // indirect
	golang.org/x/net v0.0.0-20191004110552-13f9640d40b9 // indirect
	golang.org/x/oauth2 v0.0.0-20190604053449-0f29369cfe45 // indirect
	golang.org/x/sys v0.0.0-20190826190057-c7b8b68b1456 // indirect
	golang.org/x/text v0.3.2 // indirect
	golang.org/x/time v0.0.0-20181108054448-85acf8d2951c // indirect
	google.golang.org/appengine v1.6.1 // indirect
	gopkg.in/fsnotify.v1 v1.4.7 // indirect
	gopkg.in/inf.v0 v0.9.0 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v2 v2.2.8 // indirect
	k8s.io/klog v1.0.0 // indirect
	k8s.io/utils v0.0.0-20200324210504-a9aa75ae1b89 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
)
