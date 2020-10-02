## kubecolor

Colorize your kubectl output

## What's this?

kubecolor colorizes your `kubectl` command output and nothing else.
kubecolor internally calls `kubectl` command and try to colorizes the output so
you can use kubecolor as a complete alternative of kubectl. It means you can write this in your .bash_profile:

```sh
alias kubectl = kubecolor
```

kubecolor is developed to colorize the output of only READ commands (get, describe...). 
So if the given subcommand was for WRITE operations (apply, edit...), it doesn't give any decorations on it.

For now, not all some subcommands are supported and will be done in the future. What is supported can be found below.
Even if what you want to do is not supported by kubecolor now, kubecolor still can just show `kubectl` output without any decorations,
so you don't need to switch kubecolor and kubectl but you always can use kubecolor.

Additionally, if `kubectl` resulted an error, kubecolor just shows the error message in red.

## Supported commands

Checked: Supported and works in current latest version
Unchecked: Will be supported but it's still under development
Not in the list: Won't be supported because it's not READ operation

### kubectl commands

- [ ] kubectl get
  - [x] pod
  - [ ] node
  - [ ] replicaset
  - [ ] deployment
  - [ ] service
  - [ ] and more...
- [x] kubectl top
  - [x] pod
  - [x] node
- [ ] kubectl describe
  - [ ] pod
  - [ ] node
  - [ ] replicaset
  - [ ] deployment
  - [ ] service
  - [ ] and more...
- [ ] kubectl cluster-info

### format options

- [x] json
- [x] wide
- [x] yaml

## Contributions

Always welcome. Just opening an issue should be also greatful.

## LICENSE

MIT
