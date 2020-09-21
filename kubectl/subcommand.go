package kubectl

import (
	"errors"
	"fmt"
	"strings"
)

type SubcommandInfo struct {
	Subcommand   Subcommand
	Target       Target
	FormatOption FormatOption
	NoHeader     bool
}

type FormatOption int

const (
	Wide FormatOption = iota + 1
	Json
	Yaml
)

type Subcommand int

const (
	Get Subcommand = iota + 1
	Top
	Describe
)

func (s Subcommand) String() string {
	switch s {
	case Get:
		return "get"
	case Describe:
		return "describe"
	case Top:
		return "top"
	}
	return ""
}

type Target int

const (
	Pod Target = iota + 1
	Deployment
	Node
	ReplicaSet
)
