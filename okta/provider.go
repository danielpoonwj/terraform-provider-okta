package okta

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"org_name": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OKTA_ORG_NAME", nil),
				Description: "The organization to manage in Okta.",
			},
			"api_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OKTA_API_TOKEN", nil),
				Description: "API Token granting privileges to Okta API.",
			},
			"base_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "okta.com",
				Description: "The Okta url. (Use 'oktapreview.com' for Okta testing)",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"okta_users":  resourceUsers(),
			"okta_groups": resourceGroups(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	log.Println("[INFO] Initializing Okta client")
	config := Config{
		orgName:  d.Get("org_name").(string),
		domain:   d.Get("base_url").(string),
		apiToken: d.Get("api_token").(string),
	}
	if err := config.loadAndValidate(); err != nil {
		return nil, err
	}
	return &config, nil
}
