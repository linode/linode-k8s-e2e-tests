package framework

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"time"

	"github.com/appscode/go/wait"

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
	glog.Info("Running command %q\n", cmd)
	return c.Run()
}

func deleteInForeground() *metav1.DeleteOptions {
	policy := metav1.DeletePropagationForeground
	return &metav1.DeleteOptions{PropagationPolicy: &policy}
}

func getHTTPResponse(link string) (bool, string, error) {
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

func WaitForHTTPResponse(link string) error {
	return wait.PollImmediate(RetryInterval, RetryTimout, func() (bool, error) {
		ok, _, err := getHTTPResponse(link)
		if err != nil {
			return false, nil
		}
		if ok {
			log.Println("Got response from the Backend Service")
			return true, nil
		}

		return false, nil
	})
}
