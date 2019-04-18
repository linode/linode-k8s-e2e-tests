package framework

import (
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
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

func (i *k8sInvocation) DeleteService(name string) error {
	err := i.kubeClient.CoreV1().Services(i.Namespace()).Delete(name, nil)
	return err
}
