package framework

import (
	"github.com/appscode/go/crypto/rand"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	Image    = "linode/linode-cloud-controller-manager:latest"
	ApiToken = ""
)

const (
	frontendImage = "gcr.io/google-samples/hello-frontend:1.0"
	backendImage  = "gcr.io/google-samples/hello-go-gke:1.0"
)

type Framework struct {
	restConfig *rest.Config
	kubeClient kubernetes.Interface
	namespace  string
	name       string
}

func New(
	restConfig *rest.Config,
	kubeClient kubernetes.Interface,
) *Framework {
	return &Framework{
		restConfig: restConfig,
		kubeClient: kubeClient,

		name:      "cloud-controller-manager",
		namespace: rand.WithUniqSuffix("ccm"),
	}
}

func (f *Framework) Invoke() *Invocation {
	r := &rootInvocation{
		Framework: f,
		app:       rand.WithUniqSuffix("csi-driver-e2e"),
	}
	return &Invocation{
		rootInvocation: r,
		Cluster:        &k8sInvocation{rootInvocation: r},
	}
}

type Invocation struct {
	*rootInvocation
	Cluster *k8sInvocation
}

type rootInvocation struct {
	*Framework
	app string
}

type k8sInvocation struct {
	*rootInvocation
}
