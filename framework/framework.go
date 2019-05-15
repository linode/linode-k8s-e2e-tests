package framework

import (
	"github.com/appscode/go/crypto/rand"
	cs "github.com/kubedb/apimachinery/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	Image          = "linode/linode-cloud-controller-manager:latest"
	ApiToken       = ""
	DockerRegistry = "kubedbci"
	DBCatalogName  = "9.6-v2"
	StorageClass   = "linode-block-storage"
)

const (
	frontendImage = "gcr.io/google-samples/hello-frontend:1.0"
	backendImage  = "gcr.io/google-samples/hello-go-gke:1.0"
)

type Framework struct {
	restConfig *rest.Config
	kubeConfig string
	kubeClient kubernetes.Interface
	extClient  cs.Interface
	namespace  string
	name       string
}

func New(
	restConfig *rest.Config,
	kubeClient kubernetes.Interface,
	extClient cs.Interface,
	kubeConfig string,
) *Framework {
	return &Framework{
		restConfig: restConfig,
		kubeClient: kubeClient,
		kubeConfig: kubeConfig,
		extClient:  extClient,
		name:       "lke-test",
		namespace:  rand.WithUniqSuffix("lke"),
	}
}

func (f *Framework) Invoke() *Invocation {
	r := &rootInvocation{
		Framework: f,
		app:       rand.WithUniqSuffix("e2e-test"),
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
