package provider

import (
	"context"
	"regexp"
	"testing"

	"github.com/foxcpp/go-mockdns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestAccResourceDNSValidationARecordSet(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDNSValidationARecordSet,
			},
		},
	})
}

func TestAccResourceDNSValidationARecordSet_addresses(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"dns": func() (*schema.Provider, error) {
				provider := New("")()

				provider.ConfigureContextFunc = func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
					return &mockdns.Resolver{
						Zones: map[string]mockdns.Zone{
							"terraform.fake.": {
								A: []string{"1.2.3.4"},
							},
						},
					}, nil
				}

				return provider, nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDNSValidationARecordSet_addresses,
			},
		},
	})
}

func TestAccResourceDNSValidationARecordSet_addresses_invalid(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"dns": func() (*schema.Provider, error) {
				provider := New("")()

				provider.ConfigureContextFunc = func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
					return &mockdns.Resolver{
						Zones: map[string]mockdns.Zone{
							"terraform.fake.": {
								A: []string{"1.2.3.4"},
							},
						},
					}, nil
				}

				return provider, nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceDNSValidationARecordSet_addresses_invalid,
				ExpectError: regexp.MustCompile(regexp.QuoteMeta("wrong addresses, expected: [2.3.4.5], got: [1.2.3.4]")),
			},
		},
	})
}

const testAccResourceDNSValidationARecordSet = `
resource "dns_validation_a_record_set" "terraform" {
  name = "terraform.io"
}
`

const testAccResourceDNSValidationARecordSet_addresses = `
resource "dns_validation_a_record_set" "terraform" {
	name = "terraform.fake"
	addresses = ["1.2.3.4"]
}
`

const testAccResourceDNSValidationARecordSet_addresses_invalid = `
resource "dns_validation_a_record_set" "terraform" {
	name = "terraform.fake"
	addresses = ["2.3.4.5"]

	timeouts {
		create = "1s"
	}
}
`
