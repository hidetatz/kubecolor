package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/dty1er/kubecolor/command"
	"github.com/dty1er/kubecolor/testhelper"
)

func main() {
	argumentFile, err := os.Open(testhelper.MustArgumentFilePath())
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(argumentFile)
	for scanner.Scan() {
		fmt.Println("--------------------")
		run(scanner.Text())
	}
}

func run(arg string) {
	args := strings.Split(arg, " ")
	kubectlOutput := getKubectlOutput(args)
	kubecolorOutput := getKubecolorOutput(args)

	fmt.Printf("args: %s\n", args)
	fmt.Printf("kubectl: \n%s\n", kubectlOutput)
	fmt.Printf("kubecolor: \n%s\n", kubecolorOutput)

	var decision string
	fmt.Printf("Do you want to save this data as testdata? y/n: ")
	fmt.Scan(&decision)

	if decision != "y" && decision != "Y" {
		fmt.Println("testdata is not saved")
		return
	}

	plainFilename := testhelper.ToPlainFilenamePathFromRoot(args)
	coloredFilename := testhelper.ToColoredFilenamePathFromRoot(args)

	if testhelper.FileExists(plainFilename) || testhelper.FileExists(coloredFilename) {
		var decision string
		fmt.Printf("Do you want to overwrite? y/n: ")
		fmt.Scan(&decision)

		if decision != "y" && decision != "Y" {
			fmt.Println("skipped")
			return
		}
	}

	plainFile, err := os.OpenFile(plainFilename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	coloredFile, err := os.OpenFile(coloredFilename, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}

	_, err = plainFile.Write([]byte(kubectlOutput))
	if err != nil {
		fmt.Printf("failed to write data in %s\n", plainFilename)
		os.Exit(1)
	}

	_, err = coloredFile.Write([]byte(kubecolorOutput))
	if err != nil {
		fmt.Printf("failed to write data in %s\n", coloredFilename)
		os.Exit(1)
	}
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
