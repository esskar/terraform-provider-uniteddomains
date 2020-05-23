package uniteddomains

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

// Provider returns a schema.Provider for UnitedDomains.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"email": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("UNTDDMNS_EMAIL", nil),
				Description: "E-mail address to be used to authenticate.",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("UNTDDMNS_PASSWORD", nil),
				Description: "Password to be used to authenticate.",
			},
			"insecure_https": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("UNTDDMNS_INSECURE_HTTPS", false),
				Description: "Disable verification of the UnitedDomains server's TLS certificate",
			},
			"ca_certificate": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("UNTDDMNS_CACERT", ""),
				Description: "Content or path of a Root CA to be used to verify UnitedDomains's SSL certificate",
			},
		},
		/*
			ResourcesMap: map[string]*schema.Resource{
				"powerdns_zone":   resourcePDNSZone(),
				"powerdns_record": resourcePDNSRecord(),
			},*/

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(data *schema.ResourceData) (interface{}, error) {
	config := Config{
		Email:         data.Get("email").(string),
		Password:      data.Get("password").(string),
		InsecureHTTPS: data.Get("insecure_https").(bool),
		CACertificate: data.Get("ca_certificate").(string),
	}

	return config.Client()
}
