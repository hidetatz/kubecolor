package color

import (
	"testing"
)

func TestApply(t *testing.T) {
	val := "test"
	expected := "\x1b[31mtest\x1b[0m"
	if applied := Apply(val, Red); applied != expected {
		t.Fatalf("failed: %v", applied)
	}
}
