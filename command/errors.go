package command

import "fmt"

type KubectlError struct {
	ExitCode int
}

func (ke *KubectlError) Error() string {
	return fmt.Sprintf("kubectl error: %d", ke.ExitCode)
}
