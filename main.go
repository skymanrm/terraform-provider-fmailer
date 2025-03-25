// main.go
package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: Provider,
	})
}

// Provider returns a terraform.ResourceProvider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("FMAILER_TOKEN", nil),
				Description: "Authentication token for FMailer API",
			},
			"endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("FMAILER_ENDPOINT", "https://api.fmailer.com"),
				Description: "The FMailer API endpoint",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"fmailer_domain_template": resourceDomainTemplate(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"fmailer_domain_template": dataSourceDomainTemplate(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

// providerConfigure configures the provider with the API client
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	token := d.Get("token").(string)
	endpoint := d.Get("endpoint").(string)

	var diags diag.Diagnostics

	client := NewClient(endpoint, token)
	return client, diags
}
