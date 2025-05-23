![Gopher](assets/gopher.png)
[![CircleCI](https://circleci.com/gh/alexfalkowski/infraops.svg?style=shield)](https://circleci.com/gh/alexfalkowski/infraops)
[![codecov](https://codecov.io/gh/alexfalkowski/infraops/graph/badge.svg?token=U3X5JGAA8I)](https://codecov.io/gh/alexfalkowski/infraops)
[![Go Report Card](https://goreportcard.com/badge/github.com/alexfalkowski/infraops/v2)](https://goreportcard.com/report/github.com/alexfalkowski/infraops/v2)
[![Go Reference](https://pkg.go.dev/badge/github.com/alexfalkowski/infraops/v2.svg)](https://pkg.go.dev/github.com/alexfalkowski/infraops/v2)
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

#### Configuration

Have a look at [configuration](area/apps/apps.pbtxt), the format is:

```pbtxt
version: "2.0"
applications: [
  {
    id: id
    kind: internal/external
    name: name
    namespace: namespace
    domain: domain
    version: version
    resources: {
      cpu: {
        min: "250m"
        max: "500m"
      }
      memory: {
        min: "128Mi"
        max: "256Mi"
      }
      storage: {
        min: "1Gi"
        max: "2Gi"
      }
    }
    secrets: ["secrets"]
    environments: [
      {
        name: name
        value: static value or "secret:secrets/value"
      }
    ]
  }
]
```

### Cloudflare (cf)

The code is bases on the package https://www.pulumi.com/registry/packages/cloudflare/.

#### Configuration

Have a look at [configuration](area/cf/cf.pbtxt), the format is:

```pbtxt
version: "2.0"
balancer_zones: [
  {
    name: name
    domain: domain
    record_names: ["name"]
    ipv4: ip
    ipv6: ip
  }
]
page_zones: [
  {
    name:   name
    domain: domain
    host:   host
  }
],
buckets: [
  {
    name: name
    region: region
  }
]
```

### DigitalOcean (do)

The code is bases on the package https://www.pulumi.com/registry/packages/digitalocean/.

#### Project

Create manually a default project with a name and description, example:

| Name          | Description                           |
| ------------- | ------------------------------------- |
| lean-thoughts | All of experiments for lean-thoughts. |

#### VPC

The account needs a default VPC. Create one manually under the region you would like with a name and description, example:

| Name         | Description               |
| ------------ | ------------------------- |
| default-fra1 | The default vpc for fra1. |

#### Configuration

Have a look at [configuration](area/do/do.pbtxt), the format is:

```pbtxt
version: "2.0"
clusters: [
  {
    name: name
    description: description
  }
]
```

### GitHub (gh)

The code is based on the package https://www.pulumi.com/registry/packages/github/.

The original idea was inspired from https://github.com/dirien/pulumi-github.

#### Creation

There is a caveat when creating repositories, that requires a 2 step process.

##### Pages

Pages can only be created after the repository is present.

So, the first step is to leave pages out.

Then the second PR, we set it to:

```pbtxt
pages: {}
```

If this repository will be used to host site with a cname, we need to add this:

```pbtxt
pages: {
  cname: "www.yoursite.com"
}
```

> [!NOTE]
> The reason for this is that there seems to be a timing issue with creating the `master` branch.

##### Collaborators

As with pages the repository needs to be present.

So, the first step is to have:

```pbtxt
enable_collaborators: false

```

Then the second PR, we set it to:

```pbtxt
enable_collaborators: true
```

> [!NOTE]
> This also seems like a timing issue, as rerunning the pipeline fixes it.

#### Configuration

Have a look at [configuration](area/gh/gh.pbtxt), the format is:

```pbtxt
version: "2.0"
repositories: [
  {
    name: name
    description: description
    homepage_url: homepage_url
    checks: ["check"]
    visibility: "public"
    template: {
      owner: owner
      repository: repository
    }
    pages: {}
    enable_collaborators: true
  }
]
```

### Kubernetes (k8s)

This contains all the packages our cluster needs.

> [!CAUTION]
> This needs to be run once you have a cluster in DigitalOcean.

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
