package testutil

import "github.com/MakeNowJust/heredoc"

func NewHereDoc(s string) string {
	return heredoc.Doc(s)
}
