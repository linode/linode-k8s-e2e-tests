package e2e_test

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/codeskyblue/go-sh"
	"github.com/linode/linode-k8s-e2e-tests/framework"
	"github.com/linode/linode-k8s-e2e-tests/rand"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CloudControllerManager", func() {
	var (
		err       error
		f         *framework.Invocation
		workers   []string
		chartName string
	)

	BeforeEach(func() {
		chartName, err = rand.WithRandomSuffix("wordpress-")
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

	var installHelmChart = func() {
		var out []byte

		Eventually(func() error {
			out, err = sh.Command("helm", "install", "stable/wordpress", "--name", chartName, "--kubeconfig", kubeconfigFile).Output()
			return err
		}).ShouldNot(HaveOccurred())

		log.Println(string(out))
	}

	var deleteHelmChart = func() {
		By("Deleting Wordpress")
		out, err := sh.Command("helm", "delete", chartName, "--purge", "--kubeconfig", kubeconfigFile).Output()
		log.Println(string(out))
		Expect(err).NotTo(HaveOccurred())

		By("Resetting Helm")
		err = framework.RunScript("helm-delete.sh", kubeconfigFile)
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
			Context("a Complex Helm Chart", func() {
				BeforeEach(func() {
					By("Initializing Helm & Tiller")
					helmInit()

					By("Installing Wordpress from Helm Chart")
					installHelmChart()
				})

				AfterEach(func() {
					By("Deleting Wordpress Helm Chart")
					deleteHelmChart()
				})

				It("should successfully deploy Wordpress helm chart and check its stateful & stateless component", func() {
					By("Getting Wordpress URL")
					url, err := f.Cluster.GetHTTPEndpoints(chartName)
					Expect(err).NotTo(HaveOccurred())

					time.Sleep(2 * time.Minute)
					fmt.Println(url[0])

					By("Checking the Wordpress URL")
					err = framework.WaitForHTTPResponse(url[0])
					Expect(err).NotTo(HaveOccurred())
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
					domain      = "getappscode.com"
				)

				BeforeEach(func() {
					labels = map[string]string{
						"app": "external-dns",
					}

					annotations = map[string]string{
						"external-dns.alpha.kubernetes.io/hostname": domain,
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
						ok, out, _ := framework.GetHTTPResponse("http://" + domain)
						output = out
						return ok
					}, timeout).Should(BeTrue())

					Expect(strings.Contains(output, "nginx")).Should(BeTrue())
				})
			})
		})
	})
})
