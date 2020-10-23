# Introducing --pretty output in kubectl

* Author: Hidetatsu Yaginuma ([@dty1er](https://github.com/dty1er))
  - Author of [kubecolor](https://github.com/dty1er/kubecolor)

## Overview

This is a design doc to describe why and how we introduce --pretty option into kubectl which
makes kubectl result colored. This doc is intended to be proposed in [kubernetes sig-cli regular meeting](https://github.com/kubernetes/community/tree/master/sig-cli#meetings) and
gather every participants' opinions widely.

## Context

On Aug 2018, an issue [Add ANSI colors to kubectl describe and other outputs](https://github.com/kubernetes/kubectl/issues/524) is opened in kubernetes/kubectl repository.
The issue author wished to colorize `kubectl describe` result to make it easier to read.
Because 20+ comments and 100+ :+1: reaction was left on the issue, we can expect coloring output will make kubectl better and users happier.

Even after 2 years since the issue is opened there were no actions about this, I wrote a tool, called [kubecolor](https://github.com/dty1er/kubecolor), which colorizes
kubectl output. I shared the tool in the issue.
[Thanks to @eddiezane](https://github.com/kubernetes/kubectl/issues/524#issuecomment-708606102) I found sig-cli meeting is actively held and started wondering if kubecolor can appear
in original kubectl implementation.

In this design doc, I will consider if we should introduce it and how we implement it.

## Do we really want kubectl result colored?

In this section, I will talk about why and why not we introduce colored output in kubectl.

### Why we should introduce colored output

* It can make it easier to read.

According to [the issue](https://github.com/kubernetes/kubectl/issues/524) I mentioned above, apparently some people are feeling current kubectl output
sometimes is not easier to read, because of its lack of color.
