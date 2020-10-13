package e2etest

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/dty1er/kubecolor/command"
)

var testArguments = []string{
	"get nodes",
	"get no",
	"get node",
	"get nodes -o json",
	"get nodes -o=json",
	"get nodes -ojson",
	"get nodes -o yaml",
	"get nodes -o=yaml",
	"get nodes -oyaml",
	"get nodes -o wide",
	"get nodes -o=wide",
	"get nodes -owide",
	"get nodes -o",
	"get deployments",
	"get replicaset",
	"get pods",
	"get",
	"-h",
	"--help",
	"get -h",
	"get --help",
	"top pods",
	"top node",
	"describe pod",
	"describe rs",
}

// TestVisual is a test for human. This test never fails but it's intended to
// be seen by developers and find any problem
// The idea is inspired by fatih/color
// https://github.com/fatih/color/blob/master/color_test.go#L160
func TestVisual(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping e2etest")
	}

	for _, arg := range testArguments {
		fmt.Println("--------------------")
		run(arg)
	}
}

func run(arg string) {
	args := strings.Split(arg, " ")
	kubectlOutput := getKubectlOutput(args)
	kubecolorOutput := getKubecolorOutput(args)

	fmt.Printf("args: %s\n", args)
	fmt.Printf("kubectl: \n%s\n", kubectlOutput)
	fmt.Printf("kubecolor: \n%s\n", kubecolorOutput)
}

func getKubectlOutput(args []string) string {
	b, _ := exec.Command("kubectl", args...).CombinedOutput()
	return string(b)
}

func getKubecolorOutput(args []string) string {
	// kubecolor writes output to stdout/err directly so
	// capture them using os.Pipe
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	stdout := os.Stdout
	stderr := os.Stderr
	// use fake stdout/err
	os.Stdout = w
	os.Stderr = w

	err = command.Run(args)
	if err != nil {
		panic(err)
	}

	// recover to original stdout/err
	os.Stdout = stdout
	os.Stderr = stderr
	w.Close()

	var buff bytes.Buffer
	io.Copy(&buff, r)

	return buff.String()
}
