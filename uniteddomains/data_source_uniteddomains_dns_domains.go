package uniteddomains

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceDnsDomains() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDnsDomainsRead,
	}
}

func dataSourceDnsDomainsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Client)
	client.newRequest("GET", "/pfapi/dns/domain-list", nil)

	return nil
}
