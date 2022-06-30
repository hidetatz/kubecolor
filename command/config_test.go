package command

import (
	"os"
	"testing"

	"github.com/hidetatz/kubecolor/testutil"
)

func Test_ResolveConfig(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		kubectlCommand string
		expectedArgs   []string
		expectedConf   *KubecolorConfig
	}{
		{
			name:         "no config",
			args:         []string{"get", "pods"},
			expectedArgs: []string{"get", "pods"},
			expectedConf: &KubecolorConfig{
				Plain:          false,
				DarkBackground: true,
				ForceColor:     false,
				PlainHierarchy: true,
				KubectlCmd:     "kubectl",
			},
		},
		{
			name:         "plain, dark, force",
			args:         []string{"get", "pods", "--plain", "--light-background", "--force-colors"},
			expectedArgs: []string{"get", "pods"},
			expectedConf: &KubecolorConfig{
				Plain:          true,
				DarkBackground: false,
				ForceColor:     true,
				PlainHierarchy: true,
				KubectlCmd:     "kubectl",
			},
		},
		{
			name:           "KUBECTL_COMMAND exists",
			args:           []string{"get", "pods", "--plain"},
			kubectlCommand: "kubectl.1.19",
			expectedArgs:   []string{"get", "pods"},
			expectedConf: &KubecolorConfig{
				Plain:          true,
				DarkBackground: true,
				ForceColor:     false,
				PlainHierarchy: true,
				KubectlCmd:     "kubectl.1.19",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			if tt.kubectlCommand != "" {
				os.Setenv("KUBECTL_COMMAND", tt.kubectlCommand)
				defer os.Unsetenv("KUBECTL_COMMAND")
			}

			args, conf := ResolveConfig(tt.args)
			testutil.MustEqual(t, tt.expectedArgs, args)
			testutil.MustEqual(t, tt.expectedConf, conf)
		})
	}
}
