package provider

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDNSAddressValidation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSAddressValidationCreate,
		ReadContext:   resourceDNSAddressValidationRead,
		DeleteContext: resourceDNSAddressValidationDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The full name of the DNS record",
			},
			"addresses": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				ForceNew:    true,
				Description: "Set of expected addresses that the record should return",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Second),
		},
	}
}

func resourceDNSAddressValidationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	resolver := meta.(Resolver)
	name := d.Get("name").(string)

	err := resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		addrs, err := resolver.LookupHost(ctx, name)
		if err != nil {
			var dnsError *net.DNSError
			if errors.As(err, &dnsError) {
				if dnsError.IsTemporary || dnsError.IsNotFound {
					return resource.RetryableError(err)
				}
			}

			return resource.NonRetryableError(err)
		}

		v, ok := d.GetOk("addresses")
		if !ok {
			return nil
		}
		expected := v.(*schema.Set)

		addresses := make([]interface{}, len(addrs))
		for i, addr := range addrs {
			addresses[i] = addr
		}
		actual := schema.NewSet(expected.F, addresses)

		if !expected.Equal(actual) {
			return resource.RetryableError(fmt.Errorf("wrong addresses, expected: %s, got: %s", addressesString(expected), addressesString(actual)))
		}

		return nil
	})

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(name)

	return resourceDNSAddressValidationRead(ctx, d, meta)
}

func resourceDNSAddressValidationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	resolver := meta.(Resolver)

	name := d.Get("name").(string)
	addrs, err := resolver.LookupHost(ctx, name)
	if err != nil {
		var dnsError *net.DNSError
		if errors.As(err, &dnsError) && dnsError.IsNotFound {
			log.Printf("[WARN] Record not found: %s, will re-validate", name)
			d.SetId("")
			return nil
		}

		return diag.FromErr(err)
	}

	d.Set("name", name)

	if _, ok := d.GetOk("addresses"); ok {
		addresses := make([]interface{}, len(addrs))
		for i, addr := range addrs {
			addresses[i] = addr
		}
		d.Set("addresses", addresses)
	}

	return nil
}

func resourceDNSAddressValidationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}

func addressesString(set *schema.Set) string {
	addrs := make([]string, set.Len())
	for i, v := range set.List() {
		addrs[i] = v.(string)
	}
	return "[" + strings.Join(addrs, ",") + "]"
}
