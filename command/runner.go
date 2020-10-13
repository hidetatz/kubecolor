package command

import (
	"os/exec"
	"sync"

	"github.com/dty1er/kubecolor/color"
	"github.com/dty1er/kubecolor/kubectl"
	"github.com/dty1er/kubecolor/printer"
	"github.com/mattn/go-colorable"
)

var (
	stdout = colorable.NewColorableStdout()
	stderr = colorable.NewColorableStderr()
)

func Run(args []string) error {
	args, plainFlagFound := removePlainFlagIfExists(args)
	args, lightBackgroundFlagFound := removeLightBackgroundFlagIfExists(args)
	darkBackground := !lightBackgroundFlagFound

	cmd := exec.Command("kubectl", args...)

	outReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	errReader, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	// --plain
	if plainFlagFound {
		cmd.Stdout = stdout
		cmd.Stderr = stderr
	}

	if err = cmd.Start(); err != nil {
		return err
	}

	if plainFlagFound {
		cmd.Wait()
		return nil
	}

	subcommandInfo, ok := kubectl.InspectSubcommandInfo(args)

	wg := &sync.WaitGroup{}

	switch {
	case subcommandInfo.Help:
		runAsync(wg, []func(){
			func() { printer.PrintWithColor(outReader, stdout, color.Yellow) },
			func() { printer.PrintErrorOrWarning(errReader, stderr) },
		})
	case !ok:
		// given subcommand is not supported to colorize
		// so just print it in green
		runAsync(wg, []func(){
			func() { printer.PrintWithColor(outReader, stdout, color.Green) },
			func() { printer.PrintErrorOrWarning(errReader, stderr) },
		})
	default:
		runAsync(wg, []func(){
			func() { printer.Print(outReader, stdout, subcommandInfo, darkBackground) }, // TODO fix to enable configuration for light background
			func() { printer.PrintErrorOrWarning(errReader, stderr) },
		})
	}

	wg.Wait()
	cmd.Wait()

	return nil
}

func runAsync(wg *sync.WaitGroup, tasks []func()) {
	wg.Add(len(tasks))
	for _, task := range tasks {
		task := task
		go func() {
			task()
			wg.Done()
		}()
	}
}

func removePlainFlagIfExists(args []string) ([]string, bool) {
	for i, arg := range args {
		if arg == "--plain" {
			return append(args[:i], args[i+1:]...), true
		}
	}

	return args, false
}

func removeLightBackgroundFlagIfExists(args []string) ([]string, bool) {
	for i, arg := range args {
		if arg == "--light-background" {
			return append(args[:i], args[i+1:]...), true
		}
	}

	return args, false
}
