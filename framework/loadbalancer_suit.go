package framework

func (i *lbInvocation) GetHTTPEndpoints(name string) ([]string, error) {
	return i.getLoadBalancerURLs(name)
}
