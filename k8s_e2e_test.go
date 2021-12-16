package e2e_test

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/codeskyblue/go-sh"
	"github.com/linode/linode-k8s-e2e-tests/framework"
	"github.com/linode/linode-k8s-e2e-tests/rand"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CloudControllerManager", func() {
	var (
		err               error
		f                 *framework.Invocation
		workers           []string
		wordpressName     string
		metricsServerName string
	)

	BeforeEach(func() {
		wordpressName, err = rand.WithRandomSuffix("wordpress-")
		metricsServerName, err = rand.WithRandomSuffix("metrics-server-")
		Expect(err).NotTo(HaveOccurred())
		f, err = root.Invoke()
		Expect(err).NotTo(HaveOccurred())
		workers, err = f.GetNodeList()
		Expect(err).NotTo(HaveOccurred())
		Expect(len(workers)).Should(BeNumerically(">=", 2))
	})

	var createFrontendPodWithLabel = func(pod string, labels map[string]string) {
		p := f.Cluster.GetFrontendPodObject(pod, labels)
		err = f.Cluster.CreatePod(p)
		Expect(err).NotTo(HaveOccurred())
	}

	var createBackendPodWithLabel = func(pod string, labels map[string]string) {
		p := f.Cluster.GetBackendPodObject(pod, labels)
		err = f.Cluster.CreatePod(p)
		Expect(err).NotTo(HaveOccurred())
	}

	var createPodWithLabel = func(pod string, labels map[string]string) {
		p := f.Cluster.GetPodObject(pod, labels)
		err = f.Cluster.CreatePod(p)
		Expect(err).NotTo(HaveOccurred())
	}

	var createServiceWithSelector = func(serviceName string, selector map[string]string) {
		err = f.Cluster.CreateService(serviceName, selector, nil)
		Expect(err).NotTo(HaveOccurred())
	}

	var createService = func(serviceName string, selector, annotations map[string]string) {
		err = f.Cluster.CreateService(serviceName, selector, annotations)
		Expect(err).NotTo(HaveOccurred())
	}

	var createNetworkPolicy = func(name string, labels map[string]string) {
		np := f.Cluster.GetNetworkPolicyObject(name, labels)
		err = f.Cluster.CreateNetworkPolicy(np)
		Expect(err).NotTo(HaveOccurred())
	}

	var deletePods = func(pod string) {
		err = f.Cluster.DeletePod(pod)
		Expect(err).NotTo(HaveOccurred())
	}

	var deleteService = func(name string) {
		err = f.Cluster.DeleteService(name)
		Expect(err).NotTo(HaveOccurred())
	}

	var deleteNetworkPolicy = func(name string) {
		err = f.Cluster.DeleteNetworkPolicy(name)
		Expect(err).NotTo(HaveOccurred())
	}

	var helmInit = func() {
		err := framework.RunScript("helm-init.sh", kubeconfigFile)
		Expect(err).NotTo(HaveOccurred())
	}

	var installHelmChart = func(chartName, repoName string) {
		var out []byte
		var err error
		Eventually(func() error {
			switch chartName {
			case metricsServerName:
				out, err = sh.Command(
					"helm", "install", chartName, repoName, "--set", "args={--kubelet-insecure-tls}", "--kubeconfig", kubeconfigFile,
				).Output()

			case wordpressName:
				out, err = sh.Command(
					"helm", "install", chartName, repoName, "--set", "volumePermissions.enabled=true,mariadb.volumePermissions.enabled=true", "--kubeconfig", kubeconfigFile,
				).Output()
			default:
				err = fmt.Errorf("chart name %s not handled", chartName)
			}
			return err
		}).ShouldNot(HaveOccurred())

		log.Println(string(out))
	}

	var deleteHelmChart = func(chartName string) {
		By("Deleting Wordpress")
		out, err := sh.Command("helm", "delete", chartName, "--kubeconfig", kubeconfigFile).Output()
		log.Println(string(out))
		Expect(err).NotTo(HaveOccurred())
	}

	Describe("Test", func() {
		Context("NetworkPolicy", func() {
			Context("With Two Services", func() {
				var (
					frontendPod       string
					backendPod        string
					frontendLabels    map[string]string
					backendLabels     map[string]string
					frontendSvcName   = "frontend-svc"
					backendSvcName    = "hello"
					networkPolicyName = "test-network-policy"
				)

				BeforeEach(func() {
					frontendPod = "frontend-pod"
					backendPod = "backend-pod"

					frontendLabels = map[string]string{
						"app": "frontend",
					}
					backendLabels = map[string]string{
						"app": "backend",
					}

					By("Creating Pods")
					createBackendPodWithLabel(backendPod, backendLabels)
					createFrontendPodWithLabel(frontendPod, frontendLabels)

					By("Creating Service")
					createServiceWithSelector(backendSvcName, backendLabels)
					createServiceWithSelector(frontendSvcName, frontendLabels)
				})

				AfterEach(func() {
					By("Deleting the Pods")
					deletePods(frontendPod)
					deletePods(backendPod)

					By("Deleting the Service")
					deleteService(frontendSvcName)
					deleteService(backendSvcName)
					deleteNetworkPolicy(networkPolicyName)
				})

				It("shouldn't get response from the backend service after applying network policy", func() {
					By("Waiting for Response from the Backend Service")
					Eventually(func() bool {
						ok, _ := f.GetResponseFromPod(frontendPod, true)
						return ok
					}).Should(BeTrue())

					By("Applying NetworkPolicy")
					createNetworkPolicy(networkPolicyName, backendLabels)

					By("Checking Response form the Backend Service after Applying NetworkPolicy")
					Eventually(func() bool {
						ok, _ := f.GetResponseFromPod(frontendPod, false)
						return ok
					}).Should(BeFalse())
				})
			})
		})
	})

	Describe("Test", func() {
		Context("Deploying", func() {
			Context("a Wordpress Helm Chart (with a stateful & stateless component)", func() {
				BeforeEach(func() {
					By("Initializing Helm")
					helmInit()

					By("Installing Wordpress from Helm Chart")
					installHelmChart(wordpressName, "bitnami/wordpress")
				})

				AfterEach(func() {
					By("Deleting Wordpress Helm Chart")
					deleteHelmChart(wordpressName)
				})

				It("should successfully deploy Wordpress helm chart and check its components", func() {
					By("Getting Wordpress URL")
					url, err := f.Cluster.GetHTTPEndpoints(wordpressName)
					Expect(err).NotTo(HaveOccurred())

					By("Checking the Wordpress URL in " + url[0])
					err = f.WaitForHTTPResponse(url[0])
					Expect(err).NotTo(HaveOccurred())
				})
			})

			Context("a Metrics Server Helm Chart", func() {
				var (
					nodeMetrics *v1beta1.NodeMetricsList
				)
				BeforeEach(func() {
					By("Initializing Helm")
					helmInit()

					By("Installing Metrics Server from Helm Chart")
					installHelmChart(metricsServerName, "metrics-server/metrics-server")
				})

				AfterEach(func() {
					By("Deleting Metrics Server Helm Chart")
					deleteHelmChart(metricsServerName)
				})

				It("should successfully deploy Metrics Server helm chart and eventually reports metrics", func() {
					Eventually(func() bool {
						_, err = f.Cluster.GetPodMetrics()
						if err != nil {
							return false
						}
						nodeMetrics, err = f.Cluster.GetNodeMetrics()
						if err != nil {
							return false
						}
						Expect(len(nodeMetrics.Items)).To(BeNumerically(">", 0))
						return true
					}, f.Timeout, f.RetryInterval).Should(BeTrue())
				})
			})
		})
	})

	Describe("Test", func() {
		Context("Linode", func() {
			Context("External DNS", func() {
				var (
					serviceName = "test-service"
					podName     = "test-pod"
					timeout     = 2 * time.Hour
					labels      map[string]string
					annotations map[string]string
				)

				BeforeEach(func() {
					Expect(externalDomain).NotTo(Equal(""))

					labels = map[string]string{
						"app": "external-dns",
					}

					annotations = map[string]string{
						"external-dns.alpha.kubernetes.io/hostname": externalDomain,
					}

					By("Creating Pod")
					createPodWithLabel(podName, labels)

					By("Creating Service with External DNS")
					createService(serviceName, labels, annotations)
				})

				AfterEach(func() {
					By("Deleting Pod")
					deletePods(podName)

					By("Deleting Service with External DNS")
					deleteService(serviceName)
				})

				It("should successfully check the external dns", func() {
					var output string
					Eventually(func() bool {
						ok, out, _ := framework.GetHTTPResponse("http://" + externalDomain)
						output = out
						return ok
					}, timeout).Should(BeTrue())

					Expect(strings.Contains(output, "nginx")).Should(BeTrue())
				})
			})
		})
	})
})
