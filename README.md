[![CircleCI](https://circleci.com/gh/alexfalkowski/infraops.svg?style=svg)](https://circleci.com/gh/alexfalkowski/infraops)

A place where all infrastructure is taken care of.

## Background

The code is based on https://www.pulumi.com/.

## Areas

Each folder takes care of an area of infrastructure. Each area has a package that is used as the entry point, so it is a [facade](https://en.wikipedia.org/wiki/Facade_pattern).

### Cloudflare (CF)

The code is bases on the package https://www.pulumi.com/registry/packages/cloudflare/.

### GitHub (gh)

The code is based on the package https://www.pulumi.com/registry/packages/github/.
