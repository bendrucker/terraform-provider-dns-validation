# Terraform Provider: DNS Validation

[![tests workflow status](https://github.com/bendrucker/terraform-provider-dns-validation/workflows/tests/badge.svg?branch=master)](https://github.com/bendrucker/terraform-provider-dns-validation/actions?query=workflow%3Atests) [![terraform registry](https://img.shields.io/badge/terraform-registry-623CE4)](https://registry.terraform.io/providers/bendrucker/dns-validation)

## Requirements

*	[Terraform](https://www.terraform.io/downloads.html) >= 0.12.x
*	[Go](https://golang.org/doc/install) >= 1.15

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command: 

```sh
go install
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

In order to run the full suite of Acceptance tests, run `make testacc`. The acceptance tests do not require any real resources but do require an internet connection to run.

```sh
make testacc
```
