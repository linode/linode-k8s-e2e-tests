package e2e_test

import (
	"github.com/appscode-cloud/linode-k8s-e2e-tests/framework"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CloudControllerManager", func() {
	var (
		err     error
		f       *framework.Invocation
		workers []string
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

	var createServiceWithSelector = func(serviceName string, selector map[string]string, isFrontend bool) {
		err = f.Cluster.CreateService(serviceName, selector, isFrontend)
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
					serviceURL        string
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
					createServiceWithSelector(backendSvcName, backendLabels, false)
					createServiceWithSelector(frontendSvcName, frontendLabels, true)

					By("Retrieving Service Endpoints")
					eps, err := f.Cluster.GetHTTPEndpoints(frontendSvcName)
					Expect(err).NotTo(HaveOccurred())
					Expect(len(eps)).Should(BeNumerically(">=", 1))
					serviceURL = eps[0]
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
					err = framework.WaitForHTTPResponse(serviceURL)
					Expect(err).NotTo(HaveOccurred())

					By("Applying NetworkPolicy")
					createNetworkPolicy(networkPolicyName, backendLabels)

					By("Checking Response form the Backend Service")
					err = framework.WaitForHTTPResponse(serviceURL)
					Expect(err).To(HaveOccurred())
				})
			})
		})
	})
})
