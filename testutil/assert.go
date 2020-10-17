// testutil package is a utility for testing.
// This package is inspired by morikuni/failure testutil_test.go
package testutil

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func MustEqual(t testing.TB, want, got interface{}) {
	t.Helper()

	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("diff (-want +got):\n%s", diff)
	}
}
