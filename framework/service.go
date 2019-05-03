package framework

import (
	"fmt"
	"net/url"
	"time"

	"github.com/pkg/errors"
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/wait"
)

func (i *k8sInvocation) CreateService(serviceName string, selector map[string]string) error {
	_, err := i.kubeClient.CoreV1().Services(i.Namespace()).Create(&core.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: i.Namespace(),
		},
		Spec: core.ServiceSpec{
			Ports: []core.ServicePort{
				{
					Port:       80,
					TargetPort: intstr.FromInt(80),
					Protocol:   "TCP",
				},
			},
			Selector: selector,
		},
	})

	return err
}

func (i *k8sInvocation) GetHTTPEndpoints(name string) ([]string, error) {
	var serverAddr []string

	svc, err := i.GetServiceWithLoadBalancerStatus(name, "default")
	if err != nil {
		return serverAddr, err
	}

	ips := make([]string, 0)
	for _, ingress := range svc.Status.LoadBalancer.Ingress {
		ips = append(ips, ingress.IP)
	}

	var ports []int32
	if len(svc.Spec.Ports) > 0 {
		for _, port := range svc.Spec.Ports {
			if port.NodePort > 0 {
				ports = append(ports, port.Port)
			}
		}
	}

	for _, port := range ports {
		for _, ip := range ips {
			u, err := url.Parse(fmt.Sprintf("http://%s:%d", ip, port))
			if err != nil {
				return nil, err
			}
			serverAddr = append(serverAddr, u.String())
		}
	}

	return serverAddr, nil
}

func (i *k8sInvocation) GetServiceWithLoadBalancerStatus(name, namespace string) (*core.Service, error) {
	var (
		svc *core.Service
		err error
	)
	err = wait.PollImmediate(2*time.Second, 20*time.Minute, func() (bool, error) {
		svc, err = i.kubeClient.CoreV1().Services(namespace).Get(name, metav1.GetOptions{})
		if err != nil || len(svc.Status.LoadBalancer.Ingress) == 0 { // retry
			return false, nil
		} else {
			return true, nil
		}
	})
	if err != nil {
		return nil, errors.Errorf("failed to get Status.LoadBalancer.Ingress for service %s/%s", name, namespace)
	}
	return svc, nil
}

func (i *k8sInvocation) DeleteService(name string) error {
	err := i.kubeClient.CoreV1().Services(i.Namespace()).Delete(name, nil)
	return err
}
