package command

import "os"

type KubecolorConfig struct {
	Plain                bool
	DarkBackground       bool
	ForceColor           bool
	ShowKubecolorVersion bool
	KubectlCmd           string
}

func ResolveConfig(args []string) ([]string, *KubecolorConfig) {
	args, plainFlagFound := findAndRemoveBoolFlagIfExists(args, "--plain")
	args, lightBackgroundFlagFound := findAndRemoveBoolFlagIfExists(args, "--light-background")
	args, forceColorFlagFound := findAndRemoveBoolFlagIfExists(args, "--force-colors")
	args, kubecolorVersionFlagFound := findAndRemoveBoolFlagIfExists(args, "--kubecolor-version")

	darkBackground := !lightBackgroundFlagFound

	colorsForcedViaEnv := os.Getenv("KUBECOLOR_FORCE_COLORS") == "true"

	kubectlCmd := "kubectl"
	if kc := os.Getenv("KUBECTL_COMMAND"); kc != "" {
		kubectlCmd = kc
	}

	return args, &KubecolorConfig{
		Plain:                plainFlagFound,
		DarkBackground:       darkBackground,
		ForceColor:           forceColorFlagFound || colorsForcedViaEnv,
		ShowKubecolorVersion: kubecolorVersionFlagFound,
		KubectlCmd:           kubectlCmd,
	}
}

func findAndRemoveBoolFlagIfExists(args []string, key string) ([]string, bool) {
	for i, arg := range args {
		if arg == key {
			return append(args[:i], args[i+1:]...), true
		}
	}

	return args, false
}
