package testutil

import "github.com/MakeNowJust/heredoc"

func NewHereDoc(s string) string {
	return heredoc.Doc(s)
}

func NewHereDocf(s string, args ...interface{}) string {
	return heredoc.Docf(s, args...)
}
