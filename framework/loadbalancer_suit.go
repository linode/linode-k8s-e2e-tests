package framework

func (i *k8sInvocation) GetHTTPEndpoints(name string) ([]string, error) {
	return i.getLoadBalancerURLs(name)
}
