package provider

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDNSValidationARecordSet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSValidationARecordSetCreate,
		ReadContext:   resourceDNSValidationARecordSetRead,
		DeleteContext: resourceDNSValidationARecordSetDelete,

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
	}
}

func resourceDNSValidationARecordSetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	resolver := meta.(*net.Resolver)
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
		actual := schema.NewSet(schema.HashString, addresses)

		if !expected.Equal(actual) {
			return resource.RetryableError(fmt.Errorf("wrong addresses, expected: %v, got: %v", expected, actual))
		}

		return nil
	})

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(name)

	return resourceDNSValidationARecordSetRead(ctx, d, meta)
}

func resourceDNSValidationARecordSetRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	resolver := meta.(*net.Resolver)

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
		d.Set("addresses", schema.NewSet(schema.HashString, addresses))
	}

	return nil
}

func resourceDNSValidationARecordSetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return nil
}
