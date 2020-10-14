package command

import (
	"io"
	"os"
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

	// subcommandFound becomes false when subcommand is not found; e.g. "kubecolor --help"
	subcommandInfo, subcommandFound := kubectl.InspectSubcommandInfo(args)

	// when the given subcommand is supported AND --plain is NOT specified, then we colorize it
	shouldColorize := subcommandFound && isColoringSupported(subcommandInfo.Subcommand) && !plainFlagFound

	// when it is for Help, we exceptionally colorize it in Help color
	if subcommandInfo.Help {
		shouldColorize = true
	}

	cmd := exec.Command("kubectl", args...)
	cmd.Stdin = os.Stdin

	var outReader, errReader io.Reader
	var err error

	if shouldColorize {
		outReader, err = cmd.StdoutPipe()
		if err != nil {
			return err
		}

		errReader, err = cmd.StderrPipe()
		if err != nil {
			return err
		}
	} else {
		// when we don't colorize the output,
		// we don't need to capthre it so just set stdout/err here
		cmd.Stdout = stdout
		cmd.Stderr = stderr
	}

	if err = cmd.Start(); err != nil {
		return err
	}

	if !shouldColorize {
		cmd.Wait()
		return nil
	}

	wg := &sync.WaitGroup{}

	if subcommandInfo.Help {
		runAsync(wg, []func(){
			func() { printer.PrintWithColor(outReader, stdout, color.Yellow) },
			func() { printer.PrintErrorOrWarning(errReader, stderr) },
		})

	} else {
		runAsync(wg, []func(){
			func() { printer.Print(outReader, stdout, subcommandInfo, darkBackground) },
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

func isColoringSupported(sc kubectl.Subcommand) bool {
	// when you add something here, it won't be colorized
	unsupported := []kubectl.Subcommand{
		kubectl.Create,
		kubectl.Delete,
		kubectl.Edit,
		kubectl.Attach,
		kubectl.Replace,
		kubectl.Completion,
		kubectl.Exec,
		kubectl.Proxy,
		kubectl.Plugin,
		kubectl.Wait,
	}

	for _, u := range unsupported {
		if sc == u {
			return false
		}
	}

	return true
}
