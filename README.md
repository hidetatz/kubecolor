# kubecolor

![test](https://github.com/dty1er/kubecolor/workflows/test/badge.svg?branch=main)
[![Go Report Card](https://goreportcard.com/badge/github.com/dty1er/kubecolor)](https://goreportcard.com/report/github.com/dty1er/kubecolor)
[![codecov](https://codecov.io/gh/dty1er/kubecolor/branch/main/graph/badge.svg?token=k6ysAa5ghD)](https://codecov.io/gh/dty1er/kubecolor/)

Colorize your kubectl output

* get pods

![image](https://user-images.githubusercontent.com/60682957/95733375-04929680-0cbd-11eb-82f3-adbcfecf4a3e.png)

* describe pods

![image](https://user-images.githubusercontent.com/60682957/95733389-08beb400-0cbd-11eb-983b-cf5138277fe3.png)

* something wrong

![image](https://user-images.githubusercontent.com/60682957/95733397-0a887780-0cbd-11eb-8875-bb1000e0e597.png)

* You can change color theme for light-backgrounded environment

![image](https://user-images.githubusercontent.com/60682957/95733403-0c523b00-0cbd-11eb-9ff9-abc5469e97ca.png)

## What's this?

kubecolor colorizes your `kubectl` command output and does nothing else.
kubecolor internally calls `kubectl` command and try to colorizes the output so
you can use kubecolor as a complete alternative of kubectl. It means you can write this in your .bash_profile:

```sh
alias kubectl="kubecolor"
```

kubecolor is developed to colorize the output of only READ commands (get, describe...). 
So if the given subcommand was for WRITE operations (apply, edit...), it doesn't give great decorations on it.

For now, not all subcommands are supported and will be done in the future. What is supported can be found below.
Even if what you want to do is not supported by kubecolor now, kubecolor still can just show `kubectl` output without any decorations,
so you don't need to switch kubecolor and kubectl but you always can use kubecolor.

Additionally, if `kubectl` resulted an error, kubecolor just shows the error message in red or yellow.

## Installation

### Mac and Linux users via Homebrew

```sh
brew install dty1er/tap/kubecolor
```

### Manually via go command

```sh
go install github.com/dty1er/kubecolor/cmd/kubecolor@latest
```

If you are not using module mode (or if just above doesn't work), try below:

```sh
go get -u github.com/dty1er/kubecolor/cmd/kubecolor
```

## Usage

kubecolor understands every subcommands and options which are available for `kubectl`. What you have to do is just using `kubecolor`
instead of `kubectl` like:

```sh
kubecolor --context=your_context get pods -o json
```

If you want to make the colorized kubectl default on your shell, just add this line into your shell configuration file:

```sh
alias kubectl="kubecolor"
```

### Flags

* `--plain`

When you don't want to colorize output, you can specify `--plain`. Kubecolor understands this option and
outputs the result without colorizing. Of course, given `--plain` will never be passed to `kubectl`.
This option will help you when you want to save the output onto a file and edit them by editors.

* `--light-background`

When your terminal's background color is something light (e.g white), default color preset might look too bright and not readable.
If so, specify `--light-background` as a command line argument. kubecolor will use a color preset for light-backgrounded environment.

* `--force-colors`

By default, kubecolor never output the result in colors when the tty is not a terminal.
If you want to force kubecolor to show the result in colors for non-terminal tty, you can specify this flag.
For example, when you want to pass kubecolor result to grep (`kubecolor get pods | grep pod_name`), this option is useful.

### Autocompletion

kubectl provides [autocompletion feature](https://kubernetes.io/docs/tasks/tools/install-kubectl/#enable-kubectl-autocompletion). If you are
already using it, you might have to configure it for kubecolor.
Please also refer to [kubectl official doc for kubectl autorcomplete](https://kubernetes.io/docs/reference/kubectl/cheatsheet/#kubectl-autocomplete).

#### Bash

For Bash, configuring autocompletion requires adding following line in your shell config file.

```shell
# autocomplete for kubecolor
complete -o default -F __start_kubectl kubecolor
```

If you are using an alias like `k="kubecolor"`, then just change above like:

```shell
complete -o default -F __start_kubectl k
```

#### Zsh

For zsh make sure these lines are present in your zsh config file:
```shell
# get zsh complete kubectl
source <(kubectl completion zsh)
alias kubectl=kubecolor
# make completion work with kubecolor
compdef kubecolor=kubectl
```

#### fish

Fish completion for `kubectl` is not officially not supported.

For a workaround, you can use an unofficial completion set (**WARNING: long-term support not guaranteed**):

```shell
fisher install awinecki/fish-kubecolor-completions
```

or

```shell
mkdir -p ~/.config/fish/completions
cd ~/.config/fish
git clone https://github.com/awinecki/fish-kubecolor-completions
ln -s ../fish-kubecolor-completions/completions/kubecolor.fish completions/
```

For an alias, add `alias kc=kubecolor` to your `config.fish` file.

### Specify what command to execute as kubectl

Sometimes, you may want to specify which command to use as `kubectl` internally in kubecolor. For example, when you want to use a versioned-kubectl `kubectl.1.19`, you can do that by an environment variable:

```shell
KUBECTL_COMMAND="kubectl.1.19" kubecolor get po
```

When you don't set `KUBECTL_COMMAND`, then `kubectl` is used by default.

## Supported kubectl version

Because kubecolor internally calls `kubectl` command, if you are using unsupported kubectl version, it's also not supported by kubecolor.
Kubernetes version support policy can be found in [official doc](https://kubernetes.io/docs/setup/release/version-skew-policy/).

## Contributions

Always welcome. Just opening an issue should be also greatful.

## LICENSE

MIT
