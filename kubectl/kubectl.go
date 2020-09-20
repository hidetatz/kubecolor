package kubectl

import (
	"context"
	"errors"
	"os/exec"
)

func Execute(ctx context.Context, args []string) ([]byte, error) {
	out, err := exec.CommandContext(ctx, "kubectl", args...).CombinedOutput()
	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return nil, errors.New(string(out))
		}

		return nil, err
	}

	return out, nil
}
