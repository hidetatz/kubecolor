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

func InspectSubcommand(command string) (Subcommand, bool) {
	switch command {
	case "get":
		return Get, true
	case "describe":
		return Describe, true
	case "top":
		return Top, true

	default:
		return Subcommand(0), false
	}
}

func InspectTarget(target string) (Target, bool) {
	switch target {
	case "po", "pod", "pods":
		return Pod, true
	case "no", "node", "nodes":
		return Node, true
	case "deploy", "deployment", "deployments":
		return Deployment, true
	case "rs", "replicaset", "replicasets":
		return ReplicaSet, true
	default:
		return Target(0), false
	}
}

func CollectCommandlineOptions(args []string, info *SubcommandInfo) {
	for i := range args {
		if args[i] == "--output" {
			if len(args)-1 > i {
				formatOption := args[i+1]
				switch formatOption {
				case "json":
					info.FormatOption = Json
				case "yaml":
					info.FormatOption = Yaml
				case "wide":
					info.FormatOption = Wide
				default:
					// custom-columns, go-template, etc are currently not supported
				}
			}
		} else if strings.HasPrefix(args[i], "-o") {
			switch args[i] {
			// both '-ojson' and '-o=json' works
			case "-ojson", "-o=json":
				info.FormatOption = Json
			case "-oyaml", "-o=yaml":
				info.FormatOption = Yaml
			case "-owide", "-o=wide":
				info.FormatOption = Wide
			default:
				// otherwise, look for next arg because '-o json' also works
				if len(args)-1 > i {
					formatOption := args[i+1]
					switch formatOption {
					case "json":
						info.FormatOption = Json
					case "yaml":
						info.FormatOption = Yaml
					case "wide":
						info.FormatOption = Wide
					default:
						// custom-columns, go-template, etc are currently not supported
					}
				}

			}
		} else if args[i] == "--no-headers" {
			info.NoHeader = true
		}
	}
}
