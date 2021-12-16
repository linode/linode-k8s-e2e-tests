package framework

import (
	"context"

	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

func (i *k8sInvocation) GetFrontendPodObject(podName string, labels map[string]string) *core.Pod {
	return &core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: i.Namespace(),
			Labels:    labels,
		},
		Spec: core.PodSpec{
			Containers: []core.Container{
				{
					Name:  "nginx",
					Image: frontendImage,
					Lifecycle: &core.Lifecycle{
						PreStop: &core.Handler{
							Exec: &core.ExecAction{
								Command: []string{"/usr/sbin/nginx", "-s", "quit"},
							},
						},
					},
				},
			},
		},
	}
}

func (i *k8sInvocation) GetBackendPodObject(podName string, labels map[string]string) *core.Pod {
	return &core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: i.Namespace(),
			Labels:    labels,
		},
		Spec: core.PodSpec{
			Containers: []core.Container{
				{
					Name:  "nginx",
					Image: backendImage,
					Ports: []core.ContainerPort{
						{
							Name:          "http",
							ContainerPort: 80,
						},
					},
				},
			},
		},
	}
}

func (i *k8sInvocation) GetPodObject(podName string, labels map[string]string) *core.Pod {
	return &core.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: i.Namespace(),
			Labels:    labels,
		},
		Spec: core.PodSpec{
			Containers: []core.Container{
				{
					Name:  "nginx",
					Image: "nginx",
					Ports: []core.ContainerPort{
						{
							Name:          "http",
							ContainerPort: 80,
						},
					},
				},
			},
		},
	}
}

func (i *k8sInvocation) CreatePod(pod *core.Pod) error {
	pod, err := i.kubeClient.CoreV1().Pods(i.Namespace()).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		return err
	}
	return i.WaitForReady(pod.ObjectMeta)

}

func (i *k8sInvocation) GetPodMetrics() (*v1beta1.PodMetricsList, error) {
	metrics, err := i.metricsClient.MetricsV1beta1().PodMetricses("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return metrics, nil
}

func (i *k8sInvocation) DeletePod(name string) error {
	return i.kubeClient.CoreV1().Pods(i.Namespace()).Delete(context.TODO(), name, *deleteInForeground())
}

func (i *k8sInvocation) GetPod(name, ns string) (*core.Pod, error) {
	return i.kubeClient.CoreV1().Pods(ns).Get(context.TODO(), name, metav1.GetOptions{})
}

func (i *k8sInvocation) WaitForReady(meta metav1.ObjectMeta) error {
	return wait.PollImmediate(i.RetryInterval, i.Timeout, func() (bool, error) {
		pod, err := i.kubeClient.CoreV1().Pods(i.Namespace()).Get(context.TODO(), meta.Name, metav1.GetOptions{})
		if pod == nil || err != nil {
			return false, nil
		}
		if pod.Status.Phase == core.PodRunning {
			return true, nil
		}
		return false, nil
	})
}
