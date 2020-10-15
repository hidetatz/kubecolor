package command

import (
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/dty1er/kubecolor/color"
	"github.com/dty1er/kubecolor/kubectl"
	"github.com/dty1er/kubecolor/printer"
	"github.com/mattn/go-colorable"
)

var (
	Stdout = colorable.NewColorableStdout()
	Stderr = colorable.NewColorableStderr()
)

type Printers struct {
	HelpPrinter        printer.Printer
	FullColoredPrinter printer.Printer
	ErrorPrinter       printer.Printer
}

// This is defined here to be replaced in test
var getPrinters = func(subcommandInfo *kubectl.SubcommandInfo, darkBackground bool) *Printers {
	return &Printers{
		HelpPrinter: &printer.SingleColoredPrinter{
			Color: color.Yellow,
		},
		FullColoredPrinter: &printer.KubectlOutputColoredPrinter{
			SubcommandInfo: subcommandInfo,
			DarkBackground: darkBackground,
		},
		ErrorPrinter: &printer.WithFuncPrinter{
			Fn: func(line string) color.Color {
				if strings.HasPrefix(strings.ToLower(line), "error") {
					return color.Red
				}

				return color.Yellow
			},
		},
	}
}

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

	kubectlCmd := "kubectl"
	if kc := os.Getenv("KUBECTL_COMMAND"); kc != "" {
		kubectlCmd = kc
	}
	cmd := exec.Command(kubectlCmd, args...)
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
		cmd.Stdout = Stdout
		cmd.Stderr = Stderr
	}

	if err = cmd.Start(); err != nil {
		return err
	}

	if !shouldColorize {
		cmd.Wait()
		return nil
	}

	printers := getPrinters(subcommandInfo, darkBackground)

	wg := &sync.WaitGroup{}

	if subcommandInfo.Help {
		runAsync(wg, []func(){
			func() { printers.HelpPrinter.Print(outReader, Stdout) },
			func() { printers.ErrorPrinter.Print(errReader, Stderr) },
		})

	} else {
		runAsync(wg, []func(){
			func() { printers.FullColoredPrinter.Print(outReader, Stdout) },
			func() { printers.ErrorPrinter.Print(errReader, Stderr) },
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
