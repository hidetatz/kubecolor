# kubecolor

![test](https://github.com/dty1er/kubecolor/workflows/test/badge.svg?branch=main)

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

### Autocompletion

kubectl provides [autocompletion feature](https://kubernetes.io/docs/tasks/tools/install-kubectl/#enable-kubectl-autocompletion). If you are
already using it, you might have to configure it for kubecolor.

Basically, configuring autocompletion for `kubecolor` requires adding following line in your shell config file.

```shell
# autocomplete for kubecolor
complete -o default -F __start_kubectl kubecolor
```

If you are using an alias like `k="kubecolor"`, then just change above like:

```shell
complete -o default -F __start_kubectl k
```

Please also refer to [kubectl official doc](https://kubernetes.io/docs/reference/kubectl/cheatsheet/#kubectl-autocomplete).

### Specify what command to execute as kubectl

Sometimes, you may want to specify what to command to use as `kubectl`. For example,
when you want to use a versioned-kubectl `kubectl.1.17`, you can do that by an environment variable:

```shell
KUBECTL_COMMAND="kubectl.1.17" kubecolor get po
```

When you don't set `KUBECTL_COMMAND`, then `kubectl` is used by default.

## Supported commands

Checked: Supported and works in current latest version  
Unchecked: Will be supported but it's still under development  
Not in the list: Won't be supported because it's not READ operation

### kubectl commands

- [x] kubectl get
- [x] kubectl top
- [x] kubectl describe
- [ ] kubectl explain
- [ ] kubectl logs
- [x] kubectl api-rsources
- [x] kubectl api-versions
- [ ] kubectl version

### format options

- [x] json
- [x] wide
- [x] yaml
- [x] custom-columns

## Other features which currently unsupported but will be done in the future

- [x] make it works with -w option
- [ ] Configuring custom colors
- [ ] specifying multiple resources at once (e.g. `kubectl get pod,replicaset`)
  - This will actually work, but if you don't specify "--no-headers" it might look a bit strange.

## Contributions

Always welcome. Just opening an issue should be also greatful.

## LICENSE

MIT
