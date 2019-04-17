package framework

import (
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (i *lbInvocation) GetNetworkPolicyObject(name string, labels map[string]string) *v1.NetworkPolicy {
	return &v1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: i.Namespace(),
		},
		Spec: v1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: labels,
			},
			PolicyTypes: []v1.PolicyType{
				v1.PolicyTypeIngress,
			},
			Ingress: []v1.NetworkPolicyIngressRule{
				{
					From: []v1.NetworkPolicyPeer{
						{
							IPBlock: &v1.IPBlock{
								CIDR: "192.168.0.0/16",
							},
						},
					},
				},
			},
		},
	}
}

func (i *lbInvocation) CreateNetworkPolicy(np *v1.NetworkPolicy) error {
	_, err := i.kubeClient.NetworkingV1().NetworkPolicies(i.Namespace()).Create(np)

	return err
}

func (i *lbInvocation) DeleteNetworkPolicy(name string) error {
	return i.kubeClient.NetworkingV1().NetworkPolicies(i.Namespace()).Delete(name, nil)
}
