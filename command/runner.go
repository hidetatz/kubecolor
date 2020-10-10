package command

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"sync"

	"github.com/dty1er/kubecolor/color"
	"github.com/dty1er/kubecolor/kubectl"
	"github.com/dty1er/kubecolor/printer"
)

func Run(args []string) error {
	args, plainFlagFound := removePlaneFlagIfExists(args)

	cmd := exec.Command("kubectl", args...)

	outReader, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	errReader, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	var bufout, buferr bytes.Buffer
	bufoutReader := io.TeeReader(outReader, &bufout)
	buferrReader := io.TeeReader(errReader, &buferr)

	if err = cmd.Start(); err != nil {
		return err
	}

	subcommandInfo, ok := kubectl.InspectSubcommandInfo(args)

	wg := &sync.WaitGroup{}

	switch {
	case plainFlagFound: // --plain
		runAsync(wg, []func(){
			func() { printer.PrintPlain(bufoutReader, os.Stdout) },
			func() { printer.PrintPlain(buferrReader, os.Stderr) },
		})
	case subcommandInfo.Watch:
		runAsync(wg, []func(){
			func() { printer.PrintPlain(bufoutReader, os.Stdout) },
			func() { printer.PrintPlain(buferrReader, os.Stderr) },
		})
	case subcommandInfo.Help:
		runAsync(wg, []func(){
			func() { printer.PrintWithColor(bufoutReader, os.Stdout, color.Yellow) },
			func() { printer.PrintErrorOrWarning(buferrReader, os.Stderr) },
		})
	case !ok:
		// given subcommand is not supported to colorize
		// so just print it in green
		runAsync(wg, []func(){
			func() { printer.PrintWithColor(bufoutReader, os.Stdout, color.Green) },
			func() { printer.PrintErrorOrWarning(buferrReader, os.Stderr) },
		})
	default:
		runAsync(wg, []func(){
			func() { printer.Print(bufoutReader, os.Stdout, subcommandInfo) },
			func() { printer.PrintErrorOrWarning(buferrReader, os.Stderr) },
		})
	}

	cmd.Wait()
	wg.Wait()

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

func removePlaneFlagIfExists(args []string) ([]string, bool) {
	for i, arg := range args {
		if arg == "--plain" {
			return append(args[:i], args[i+1:]...), true
		}
	}

	return args, false
}
