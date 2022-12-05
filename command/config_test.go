package command

import (
	"os"
	"testing"

	"github.com/hidetatz/kubecolor/testutil"
)

func Test_ResolveConfig(t *testing.T) {
	tests := []struct {
		name         string
		args         []string
		env          map[string]string
		expectedArgs []string
		expectedConf *KubecolorConfig
	}{
		{
			name:         "no config",
			args:         []string{"get", "pods"},
			expectedArgs: []string{"get", "pods"},
			expectedConf: &KubecolorConfig{
				Plain:          false,
				DarkBackground: true,
				ForceColor:     false,
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
				KubectlCmd:     "kubectl",
			},
		},
		{
			name:         "KUBECTL_COMMAND exists",
			args:         []string{"get", "pods", "--plain"},
			env:          map[string]string{"KUBECTL_COMMAND": "kubectl.1.19"},
			expectedArgs: []string{"get", "pods"},
			expectedConf: &KubecolorConfig{
				Plain:          true,
				DarkBackground: true,
				ForceColor:     false,
				KubectlCmd:     "kubectl.1.19",
			},
		},
		{
			name:         "KUBECOLOR_FORCE_COLORS env var",
			args:         []string{"get", "pods"},
			env:          map[string]string{"KUBECOLOR_FORCE_COLORS": "true"},
			expectedArgs: []string{"get", "pods"},
			expectedConf: &KubecolorConfig{
				Plain:          false,
				DarkBackground: true,
				ForceColor:     true,
				KubectlCmd:     "kubectl",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.env {
				os.Setenv(k, v)
			}
			defer func() {
				for k, _ := range tt.env {
					os.Unsetenv(k)
				}
			}()

			args, conf := ResolveConfig(tt.args)
			testutil.MustEqual(t, tt.expectedArgs, args)
			testutil.MustEqual(t, tt.expectedConf, conf)
		})
	}
}
