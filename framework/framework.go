package framework

import (
	"github.com/linode/linode-k8s-e2e-tests/rand"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"
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
	restConfig    *rest.Config
	kubeConfig    string
	kubeClient    kubernetes.Interface
	metricsClient *metricsclientset.Clientset
	namespace     string
	name          string
}

func New(
	restConfig *rest.Config,
	kubeClient kubernetes.Interface,
	kubeConfig string,
	metricsClient *metricsclientset.Clientset,
) (*Framework, error) {
	suffix, errSuffix := rand.WithRandomSuffix("lke")
	if errSuffix != nil {
		return nil, errSuffix
	}

	out := &Framework{
		restConfig:    restConfig,
		kubeClient:    kubeClient,
		kubeConfig:    kubeConfig,
		metricsClient: metricsClient,
		name:          "lke-test",
		namespace:     suffix,
	}

	return out, nil
}

func (f *Framework) Invoke() (*Invocation, error) {
	suffix, errSuffix := rand.WithRandomSuffix("e2e-test")
	if errSuffix != nil {
		return nil, errSuffix
	}

	r := &rootInvocation{
		Framework: f,
		app:       suffix,
	}

	out := &Invocation{
		rootInvocation: r,
		Cluster:        &k8sInvocation{rootInvocation: r},
	}

	return out, nil
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
