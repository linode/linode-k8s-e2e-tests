package framework

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/rest"
	kmodules "kmodules.xyz/client-go/tools/exec"

	"github.com/golang/glog"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
)

const (
	scriptDirectory = "scripts"
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
	pods, err := f.kubeClient.CoreV1().Pods(f.Namespace()).Get(context.TODO(), podName, metav1.GetOptions{})
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

func GetHTTPResponse(link string) (bool, string, error) {
	resp, err := http.Get(link)
	if err != nil {
		return false, "", err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, "", err
	}

	return resp.StatusCode == 200, string(bodyBytes), nil
}

func (f *Invocation) WaitForHTTPResponse(link string) error {
	return wait.PollImmediate(f.rootInvocation.RetryInterval, f.rootInvocation.Timeout, func() (bool, error) {
		ok, resp, err := GetHTTPResponse(link)
		if err != nil {
			return false, nil
		}
		if ok {
			log.Println("Got response from " + link)
			if strings.Contains(resp, "Hello world") {
				return true, nil
			}
			return true, fmt.Errorf("the response didn't have Hello world post in it")
		}

		return false, nil
	})
}

func (f Framework) ApplyManifest(manifestPath string) error {
	args := []string{"apply", "--kubeconfig", f.kubeConfig, "-f", manifestPath}
	cmd := exec.Command("kubectl", args...)
	out, err := cmd.CombinedOutput()
	fmt.Println(string(out))
	if err != nil {
		return err
	}
	return nil
}
