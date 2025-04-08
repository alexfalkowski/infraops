![Gopher](assets/gopher.png)
[![CircleCI](https://circleci.com/gh/alexfalkowski/infraops.svg?style=shield)](https://circleci.com/gh/alexfalkowski/infraops)
[![codecov](https://codecov.io/gh/alexfalkowski/infraops/graph/badge.svg?token=U3X5JGAA8I)](https://codecov.io/gh/alexfalkowski/infraops)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexfalkowski/infraops)](https://goreportcard.com/report/github.com/alexfalkowski/infraops)
[![Go Reference](https://pkg.go.dev/badge/github.com/alexfalkowski/infraops.svg)](https://pkg.go.dev/github.com/alexfalkowski/infraops)
[![Stability: Active](https://masterminds.github.io/stability/active.svg)](https://masterminds.github.io/stability/active.html)

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

### Configuration

Each area is defined by the configuration that is generated from the [protobuf](api/infraops/v1/service.proto) and the format used is the [Text Format](https://protobuf.dev/reference/protobuf/textformat-spec/).

### Applications (apps)

This consists of my open source projects https://github.com/alexfalkowski being deployed to kubernetes.

#### Install

The above is for a new application. If you want to setup all current apps, run the following.

```bash
❯ make -C area/apps setup
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

#### Project

Create manually a default project for me it is with name *lean-thoughts* and description *All of experiments for lean-thoughts.*.

#### VPC

The account needs a default VPC. Create one manually under the region you would like for me it was FRA1 with name *default-fra1* and description *The default vpc for fra1*.

### GitHub (gh)

The code is based on the package https://www.pulumi.com/registry/packages/github/.

The original idea was inspired from https://github.com/dirien/pulumi-github.

#### Creation

There is a caveat when creating repositories, that requires a 2 step process.

The first step is to have:

```go
enable_pages: false
```

Then the second PR, we set it to:

```go
enable_pages: true
```

The reason for this is that there seems to be a timing issue with creating the `master` branch.

### Kubernetes (k8s)

This contains all the packages our cluster needs.

#### Setup

To ge the cluster ready, you need to run the following:

```bash
❯ make -C area/k8s setup
```

#### Delete

To remove all the apps, you need to run the following:

```bash
❯ make -C area/k8s delete
```

### Dependencies

![Dependencies](./assets/..png)
