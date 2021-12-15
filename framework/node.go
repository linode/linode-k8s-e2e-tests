package framework

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/metrics/pkg/apis/metrics/v1beta1"
)

const (
	masterLabel = "node-role.kubernetes.io/master"
)

func (i *Invocation) GetNodeList() ([]string, error) {
	workers := make([]string, 0)
	nodes, err := i.kubeClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, node := range nodes.Items {
		if _, found := node.ObjectMeta.Labels[masterLabel]; !found {
			workers = append(workers, node.Name)
		}
	}
	return workers, nil
}

func (i *k8sInvocation) GetNodeMetrics() (*v1beta1.NodeMetricsList, error) {
	metrics, err := i.metricsClient.MetricsV1beta1().NodeMetricses().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return metrics, nil
}
