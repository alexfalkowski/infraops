[![CircleCI](https://circleci.com/gh/alexfalkowski/infraops.svg?style=svg)](https://circleci.com/gh/alexfalkowski/infraops)

A place where all infrastructure is taken care of.

## Background

The following tools are used:
- https://www.pulumi.com/
- https://kubernetes.io/docs/reference/kubectl/
- https://helm.sh/
- https://docs.kubelinter.io/
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
❯ make -C area/apps namespace=example setup-otlp
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
