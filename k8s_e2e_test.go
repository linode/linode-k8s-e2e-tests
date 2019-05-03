package e2e_test

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/appscode-cloud/linode-k8s-e2e-tests/framework"
	"github.com/codeskyblue/go-sh"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CloudControllerManager", func() {
	var (
		err       error
		f         *framework.Invocation
		workers   []string
		chartName = "test-chart"
	)

	BeforeEach(func() {
		f = root.Invoke()
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

	var createServiceWithSelector = func(serviceName string, selector map[string]string) {
		err = f.Cluster.CreateService(serviceName, selector)
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
		err := framework.RunScript("helm-init.sh", ClusterName)
		Expect(err).NotTo(HaveOccurred())
	}

	var getCurrentKubeConfig = func() string {
		wd, err := os.Getwd()
		Expect(err).NotTo(HaveOccurred())
		return wd + "/" + ClusterName + ".conf"
	}

	var installHelmChart = func() {
		var out []byte

		Eventually(func() error {
			out, err = sh.Command("helm", "install", "stable/wordpress", "--name", chartName, "--kubeconfig", getCurrentKubeConfig()).Output()
			return err
		}).ShouldNot(HaveOccurred())

		log.Println(string(out))
	}

	var deleteHelmChart = func() {
		By("Deleting Wordpress")
		out, err := sh.Command("helm", "delete", chartName, "--purge", "--kubeconfig", getCurrentKubeConfig()).Output()
		log.Println(string(out))
		Expect(err).NotTo(HaveOccurred())

		By("Resetting Helm")
		err = framework.RunScript("helm-delete.sh", ClusterName)
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
					url, err := f.Cluster.GetHTTPEndpoints(chartName + "-wordpress")
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
})
