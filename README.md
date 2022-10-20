Cloud Foundry UAA Terraform Provider
================================

# Overview

This Terraform provider plugin allows you to configure a Cloud Foundry [User Account and Authentication](https://github.com/cloudfoundry/uaa) (UAA) service declaratively using [HCL](https://github.com/hashicorp/hcl).

## Usage

**Requires** [terraform **>=0.13**](https://www.terraform.io/downloads.html)

    terraform {
      required_providers {
        uaa = {
          source  = "TBD/uaa"
          version = "latest"
        }
      }
      required_version = ">= 0.13"
    }

# Documentation

You can find documentation at https://registry.terraform.io/providers/TBD/uaa/latest/docs

# Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.19+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and copy the binary into the `~/.terraform.d/plugins/` directory.

    make build

# Testing the Provider

### Test Containers

The provider tests can be run against a dedicated containerized UAA server by running:

    go test -v -timeout 10m ./test --tags=containerized

This will spin up a docker container and a UAA container dedicated to the tests.  It will set all of the required environment variables, run the tests, and then destroy the containers afterwards.

### Existing UAA Server

To test the provider against an existing UAA server, set the following environment variables:

    export UAA_LOGIN_URL=http://localhost:8080
    export UAA_AUTH_URL=http://localhost:8080
    export UAA_CLIENT_ID=admin
    export UAA_CLIENT_SECRET=adminsecret

Once the env variables are set, the tests can be run with the following command.

    go test -v -timeout 10m ./test/...

# Debugging

You can export the following environment variables to enable detail debug logs.

    export UAA_DEBUG=true
    export UAA_TRACE=debug.log
