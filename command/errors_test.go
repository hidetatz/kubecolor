package command

import "testing"

func TestKubectlError(t *testing.T) {
	e := &KubectlError{1}
	if e.Error() != "kubectl error: 1" {
		t.Errorf("unexpected error")
	}
}
