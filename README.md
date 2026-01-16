![Gopher](assets/gopher.png)
[![CircleCI](https://circleci.com/gh/alexfalkowski/infraops.svg?style=shield)](https://circleci.com/gh/alexfalkowski/infraops)
[![codecov](https://codecov.io/gh/alexfalkowski/infraops/graph/badge.svg?token=U3X5JGAA8I)](https://codecov.io/gh/alexfalkowski/infraops)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexfalkowski/infraops/v2)](https://goreportcard.com/report/github.com/alexfalkowski/infraops/v2)
[![Go Reference](https://pkg.go.dev/badge/github.com/alexfalkowski/infraops/v2.svg)](https://pkg.go.dev/github.com/alexfalkowski/infraops/v2)
[![Stability: Active](https://masterminds.github.io/stability/active.svg)](https://masterminds.github.io/stability/active.html)

A place where all infrastructure is taken care of.

## Background

The following tools are used:

- <https://www.pulumi.com/>
- <https://kubernetes.io/docs/reference/kubectl/>
- <https://helm.sh/>
- <https://kube-score.com/>

## Configuration

The configuration follows the [prototext](https://protobuf.dev/reference/protobuf/textformat-spec/) format.

### Format

Build the application:

```bash
❯ make build-format
```

To format a configuration:

```bash
❯ ./format -k cf
```

## Areas

Each folder takes care of an area of infrastructure. Each area has a package that is used as the entry point, so it is a [facade](https://en.wikipedia.org/wiki/Facade_pattern).

### Areas Setup

To setup a new area follow the following:

- Run `pulumi new`.
- Choose the template you need, if in doubt choose `go`.
- The stack name should always be `prod`.

### Areas Configuration

Each area is defined by the configuration that is generated from the [protobuf](api/infraops/v1/service.proto) and the format used is the [Text Format](https://protobuf.dev/reference/protobuf/textformat-spec/).

### Applications (apps)

This consists of my open source projects <https://github.com/alexfalkowski> being deployed to kubernetes.

#### Applications Install

The above is for a new application. If you want to setup all current apps, run the following.

```bash
❯ make -C area/apps setup
```

#### Applications Delete

To remove all the apps, you need to run the following:

```bash
❯ make -C area/apps delete
```

#### Applications Configuration

Have a look at the [configuration](area/apps/apps.hjson).

#### Applications Version Update

Build the application:

```bash
❯ make build-bump
```

To update a version of an app:

```bash
❯ ./bump -n bezeichner -v 1.559.0
```

### Cloudflare (cf)

The code is bases on the package <https://www.pulumi.com/registry/packages/cloudflare/>.

#### Cloudflare Configuration

Have a look at the [configuration](area/cf/cf.hjson).

### DigitalOcean (do)

The code is bases on the package <https://www.pulumi.com/registry/packages/digitalocean/>.

#### DigitalOcean Project

Create manually a default project with a name and description, example:

| Name          | Description                           |
| ------------- | ------------------------------------- |
| lean-thoughts | All of experiments for lean-thoughts. |

#### DigitalOcean VPC

The account needs a default VPC. Create one manually under the region you would like with a name and description, example:

| Name         | Description               |
| ------------ | ------------------------- |
| default-fra1 | The default vpc for fra1. |

#### DigitalOcean Cluster Upgrade

The process is as follows:

- Patch versions can be updated at [do.go](internal/do/do.go).
- Minor/major versions should be updated in the [UI](https://docs.digitalocean.com/products/kubernetes/how-to/upgrade-cluster/) and then updated at [do.go](internal/do/do.go).

#### DigitalOcean Configuration

Have a look at the [configuration](area/do/do.hjson).

### GitHub (gh)

The code is based on the package <https://www.pulumi.com/registry/packages/github/>.

The original idea was inspired from <https://github.com/dirien/pulumi-github>.

#### GitHub Creation

There is a caveat when creating repositories, that requires a 2 step process.

##### GitHub Pages

Pages can only be created after the repository is present.

So, the first step is to disable it or  leave pages out.

```hjson
pages: {
  enabled: false
}
```

Then the second PR, we enable it:

```hjson
pages: {
  enabled: true
}
```

If this repository will be used to host site with a cname, we need to add this:

```hjson
pages: {
  enabled: true
  cname: www.yoursite.com
}
```

> [!NOTE]
> The reason for this is that there seems to be a timing issue with creating the `master` branch.

##### GitHub Collaborators

As with pages the repository needs to be present.

So, the first step is to have:

```hjson
collaborators: {
  enabled: false
}
```

Then the second PR, we set it to:

```hjson
collaborators: {
  enabled: true
}
```

> [!NOTE]
> This also seems like a timing issue, as rerunning the pipeline fixes it.

#### GitHub Configuration

Have a look at the [configuration](area/gh/gh.hjson).

### Kubernetes (k8s)

This contains all the packages our cluster needs.

> [!CAUTION]
> This needs to be run once you have a cluster in DigitalOcean.

#### Kubernetes Setup

To ge the cluster ready, you need to run the following:

```bash
❯ make -C area/k8s setup
```

#### Kubernetes Delete

To remove all the apps, you need to run the following:

```bash
❯ make -C area/k8s delete
```

### Dependencies

![Dependencies](./assets/..png)
