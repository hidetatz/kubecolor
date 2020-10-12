package command

import (
	"os"
	"os/exec"
	"sync"

	"github.com/dty1er/kubecolor/color"
	"github.com/dty1er/kubecolor/kubectl"
	"github.com/dty1er/kubecolor/printer"
	"github.com/mattn/go-isatty"
)

func Run(args []string, kubeColorDebug bool) error {
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

	fd := os.Stdout.Fd()
	colorize := isatty.IsTerminal(fd) || isatty.IsCygwinTerminal(fd) || kubeColorDebug

	subcommandInfo, ok := kubectl.InspectSubcommandInfo(args)
	if !ok {
		colorize = false
	}

	if !colorize {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
	}

	if err = cmd.Start(); err != nil {
		return err
	}

	if !colorize {
		cmd.Wait()
		return nil
	}

	wg := &sync.WaitGroup{}

	switch {
	case plainFlagFound: // --plain
		runAsync(wg, []func(){
			func() { printer.PrintPlain(outReader, os.Stdout) },
			func() { printer.PrintPlain(errReader, os.Stderr) },
		})
	case subcommandInfo.Help:
		runAsync(wg, []func(){
			func() { printer.PrintWithColor(outReader, os.Stdout, color.Yellow) },
			func() { printer.PrintErrorOrWarning(errReader, os.Stderr) },
		})
	case !ok:
		// given subcommand is not supported to colorize
		// so just print it in green
		runAsync(wg, []func(){
			func() { printer.PrintWithColor(outReader, os.Stdout, color.Green) },
			func() { printer.PrintErrorOrWarning(errReader, os.Stderr) },
		})
	default:
		runAsync(wg, []func(){
			func() { printer.Print(outReader, os.Stdout, subcommandInfo, darkBackground) }, // TODO fix to enable configuration for light background
			func() { printer.PrintErrorOrWarning(errReader, os.Stderr) },
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
