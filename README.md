[![CircleCI](https://circleci.com/gh/alexfalkowski/infraops.svg?style=svg)](https://circleci.com/gh/alexfalkowski/infraops)

A place where all infrastructure is taken care of.

## Background

The code is based on https://www.pulumi.com/.

## Areas

Each folder takes care of an area of infrastructure. Each area has a package that is used as the entry point, so it is a [facade](https://en.wikipedia.org/wiki/Facade_pattern).

### Cloudflare (CF)

The code is bases on the package https://www.pulumi.com/registry/packages/cloudflare/.

### DigitalOcean (DO)

The code is bases on the package https://www.pulumi.com/registry/packages/digitalocean/.

### GitHub (gh)

The code is based on the package https://www.pulumi.com/registry/packages/github/.

## Setup

To setup a new area follow the following:

```bash
❯ pulumi new
 Would you like to create a project from a template or using a Pulumi AI prompt? template
Please choose a template (227 total):
 go                                 A minimal Go Pulumi program
This command will walk you through creating a new Pulumi project.

Enter a value or leave blank to accept the (default), and press <ENTER>.
Press ^C at any time to quit.

project name (cwd): # This will default to the current working directory we are running it from.
project description (A minimal Go Pulumi program): A place to manage <Name of the provider>
Created project 'cwd'

Please enter your desired stack name.
To create a stack in an organization, use the format <org-name>/<stack-name> (e.g. `acmecorp/dev`).
stack name (dev): prod
Created stack 'prod'

Installing dependencies...

Finished installing dependencies

Your new project is ready to go!

To perform an initial deployment, run `pulumi up`

```
