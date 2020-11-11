---
layout: "dns-validation"
page_title: "Provider: DNS Validaton"
---

# DNS Validation Provider

The DNS validation provider does not create actual resources. It performs DNS validation, querying records to ensure they exist and resolve the expected records. This allows Terraform to create DNS records and await their propagation before using them in dependent resources.

## Example Usage

```hcl
provider "dns-validation" {}

resource "dns_address_validation" "example" {
  provider = dns-validation
  # ...
}
```
