---
layout: "dns-validation"
page_title: "DNS Validation: dns_address_validation"
---

# dns_address_validation

Waits until the specified DNS name resolves to a set of addresses.

## Example Usage

```hcl
resource "dns_address_validation" "foo" {\
  provider = dns-validation

  name = "terraform.io"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The DNS name to look up
* `addresses` - (Optional) Set of IP addresses that the DNS name should resolve to (IP v4 and v6). If empty, any result will be allowed.

## Timeouts

`dns_address_validation` provides the following [Timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) configuration options:

* `create` - (Default `30 seconds`) Used for waiting for the desired DNS response

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resolved DNS name. This will match `name`, but the value will not be known until after apply.
