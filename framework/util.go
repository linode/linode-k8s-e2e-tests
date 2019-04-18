package framework

import (
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
	kmodules "kmodules.xyz/client-go/tools/exec"

	"github.com/golang/glog"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	scriptDirectory = "scripts"
	RetryInterval   = 5 * time.Second
	RetryTimout     = 1 * time.Minute
)

func RunScript(script string, args ...string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	return runCommand(path.Join(wd, scriptDirectory, script), args...)
}

func runCommand(cmd string, args ...string) error {
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	c.Env = append(c.Env, append(os.Environ())...)
	glog.Infof("Running command %q\n", cmd)
	return c.Run()
}

func deleteInForeground() *metav1.DeleteOptions {
	policy := metav1.DeletePropagationForeground
	return &metav1.DeleteOptions{PropagationPolicy: &policy}
}

func (f Framework) GetResponseFromPod(podName string, install bool) (bool, error) {
	pods, err := f.kubeClient.CoreV1().Pods(f.Namespace()).Get(podName, metav1.GetOptions{})
	if err != nil {
		return false, err
	}

	if install {
		err = installCurl(f.RestConfig(), pods)
		if err != nil {
			return false, err
		}
	}

	resp, err := curlInPod(f.RestConfig(), pods)
	if err != nil {
		return false, err
	}
	return strings.Contains(resp, "Hello"), nil
}

func installCurl(config *rest.Config, pod *v1.Pod) error {
	_, _ = kmodules.ExecIntoPod(config, pod, kmodules.Command("apt-get", "update", "-y"))
	_, err := kmodules.ExecIntoPod(config, pod, kmodules.Command("apt-get", "install", "curl", "-y"))
	return err
}

func curlInPod(config *rest.Config, pod *v1.Pod) (string, error) {
	return kmodules.ExecIntoPod(config, pod, kmodules.Command("curl", "http://hello", "-s", "-m", "10"))
}
