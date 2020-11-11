package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceDNSValidationARecordSet(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDNSValidationARecordSet,
			},
		},
	})
}

const testAccResourceDNSValidationARecordSet = `
resource "dns_validation_a_record_set" "terraform" {
  name = "terraform.io"
}
`
