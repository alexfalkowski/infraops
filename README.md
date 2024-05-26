[![CircleCI](https://circleci.com/gh/alexfalkowski/infraops.svg?style=svg)](https://circleci.com/gh/alexfalkowski/infraops)

A place where all infrastructure is taken care of.

## Background

The following tools are used:
- https://www.pulumi.com/
- https://kubernetes.io/docs/reference/kubectl/
- https://helm.sh/
- https://kube-score.com/


## Areas

Each folder takes care of an area of infrastructure. Each area has a package that is used as the entry point, so it is a [facade](https://en.wikipedia.org/wiki/Facade_pattern).

### Setup

To setup a new area follow the following:
- Run `pulumi new`.
- Choose the template you need, if in doubt choose `go`.
- The stack name should always be `prod`.

### Applications (apps)

This consists of my open source projects https://github.com/alexfalkowski being deployed to kubernetes.

#### Setup

To have an app ready as an example we will use `example`, you need to run the following:

```bash
❯ make -C area/apps setup-example
```

#### Install

The above is for a new application. If you want to setup all current apps, run the following.

```bash
❯ make -C area/apps setup-konfig

❯ make -C area/apps setup-standort

❯ make -C area/apps setup-bezeichner
```

#### Delete

To remove all the apps, you need to run the following:

```bash
❯ make -C area/apps delete
```

### Cloudflare (cf)

The code is bases on the package https://www.pulumi.com/registry/packages/cloudflare/.

### DigitalOcean (do)

The code is bases on the package https://www.pulumi.com/registry/packages/digitalocean/.

### GitHub (gh)

The code is based on the package https://www.pulumi.com/registry/packages/github/.

The original idea was inspired from https://github.com/dirien/pulumi-github.

### Kubernetes (k8s)

This contains all the packages our cluster needs.

#### Setup

To ge the cluster ready, you need to run the following:

```bash
❯ make -C area/k8s setup update
```

#### Delete

To remove all the apps, you need to run the following:

```bash
❯ make -C area/k8s delete
```
