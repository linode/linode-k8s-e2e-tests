package framework

import (
	core "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/util/wait"
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
	pod, err := i.kubeClient.CoreV1().Pods(i.Namespace()).Create(pod)
	if err != nil {
		return err
	}
	return i.WaitForReady(pod.ObjectMeta)

}

func (i *k8sInvocation) DeletePod(name string) error {
	return i.kubeClient.CoreV1().Pods(i.Namespace()).Delete(name, deleteInForeground())
}

func (i *k8sInvocation) GetPod(name, ns string) (*core.Pod, error) {
	return i.kubeClient.CoreV1().Pods(ns).Get(name, metav1.GetOptions{})
}

func (i *k8sInvocation) WaitForReady(meta metav1.ObjectMeta) error {
	return wait.PollImmediate(RetryInterval, RetryTimout, func() (bool, error) {
		pod, err := i.kubeClient.CoreV1().Pods(i.Namespace()).Get(meta.Name, metav1.GetOptions{})
		if pod == nil || err != nil {
			return false, nil
		}
		if pod.Status.Phase == core.PodRunning {
			return true, nil
		}
		return false, nil
	})
}
