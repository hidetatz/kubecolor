# Introducing colored output in kubectl

* Author: Hidetatsu Yaginuma ([@dty1er](https://github.com/dty1er))
  - Author of [kubecolor](https://github.com/dty1er/kubecolor)

## Overview

This is a design doc to describe why and how we introduce colored output feature into kubectl.
This doc is intended to be proposed in [kubernetes sig-cli regular meeting](https://github.com/kubernetes/community/tree/master/sig-cli#meetings) and
gather every participants' opinions widely.

## Context

On Aug 2018, an issue [Add ANSI colors to kubectl describe and other outputs](https://github.com/kubernetes/kubectl/issues/524) is opened in kubernetes/kubectl repository.
The issue author wished to colorize `kubectl describe` result to make it easier to read.
Because 20+ comments and 100+ :+1:  reactions was left on the issue, we can expect coloring output will make kubectl better and users happier.

Even after 2 years since the issue is opened there were no actions about this, I wrote a tool, called [kubecolor](https://github.com/dty1er/kubecolor), which colorizes
kubectl output. I shared the tool in the issue.
[Thanks to @eddiezane](https://github.com/kubernetes/kubectl/issues/524#issuecomment-708606102) I found sig-cli meeting is actively held and started wondering if kubecolor can appear
in original kubectl implementation.

In this design doc, I will consider if we should introduce it and how we implement it.

## Do we really want kubectl result colored?

In this section, I will talk about why and why not we introduce colored output in kubectl.
While I'm an author of kubecolor, I tried to be fair as much as possible, and 
I tried to write down all the things I came up with. However I must be missing some parts so
the readers giving some feedback to me would be appreciated a lot.

### Why we should introduce colored output

Because it can make it easier to read.
According to [the issue](https://github.com/kubernetes/kubectl/issues/524) I mentioned above, apparently some people are feeling current kubectl output
sometimes is not easier to read, because of its lack of color.
I have written [an article in medium](https://medium.com/@dty1er/colorize-kubectl-output-by-kubecolor-2c222af3163a) to compare the results.

Some kubernetes tools like [kubectx and kubens](https://github.com/ahmetb/kubectx#kubectx--kubens-power-tools-for-kubectl) shows colored output as well.

### Why we should not introduce colored output

It makes a breaking change. For example, if a uses is doing something like below, we will break it.

```shell
# a command to notify stopping pod to Slack
$ kubectl get pod | grep CrashLoopBackOff | awk '{print $1}' | send_to_slack
```

However, this appears as a problem when we introduced colored output as **default behavior**. In other words, if we make it "opt-in" feature,
it won't make a breaking changes.

## Goal and non-goal

### Goal

* Users can "specify" if they want kubectl to show the result in colors.
* Applying color never changes the original result.

### Non-goal

* It makes colored output default behavior.
* It makes breaking changes.

## Implementation overview

Because we don't want colored output default behavior of kubectl, we introduce a new command line option "--pretty" which is scoped global.

```shell
  -p, --pretty=false: If true, the output will be colored
```

## --pretty specification

`--pretty` option can be used with every kubectl subcommand. It means `--pretty` should appear in the result of `kubectl options` .
When `--pretty` is specified, the behavior will be like below.

### 1. For WRITE command

If the subcommand is for WRITE operation, and it was successfully finished, it shows its result in green. It indicates that
the operation is finished successfully.
If it results any errors, the output should be in red. It indicates there were errors.

WRITE commands: apply, scale, patch etc.

### 2. For READ

If the subcommand is for READ operation, and it was successfully finished, it shows its result in various colors. What colors are chosen depends on
the subcommand. Because kubectl outputs the result in various format, we should support all of them.

Example (this is how current kubecolor outputs the result):

* kubectl get

![image](https://user-images.githubusercontent.com/60682957/97545993-26a34d00-1a0f-11eb-8836-09ff764b6ead.png)

* kubectl describe

![image](https://user-images.githubusercontent.com/60682957/97545997-286d1080-1a0f-11eb-9a9b-8f04c3431c7c.png)

* kubectl get -o=json

![image](https://user-images.githubusercontent.com/60682957/97546004-2a36d400-1a0f-11eb-9ff4-7bf745aa1643.png)

* kubectl version --short

![image](https://user-images.githubusercontent.com/60682957/97546126-57838200-1a0f-11eb-8872-2efa39c48c5e.png)

If it results any errors, the output should be in red. It indicates there were errors.

READ commands: get, describe, top, version etc.

### 3. For interactive operation

If the subcommand is for interactive operation, it does no colorizing.

Interactive operations: edit, attach, exec -it etc.

### 4. For --help

If the command is asking to show help, it shows its result in yellow. This indicates the output is for help.

## Implementation overview

TODOOOOOOO

* Introduce new io.Writer "colorWriter"
* Introduce a middleware which colorizes output
