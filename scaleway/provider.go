package scaleway

import (
	"sync"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var mu = sync.Mutex{}

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SCALEWAY_ACCESS_KEY", nil),
				Deprecated:  "Use `token` instead.",
				Description: "The API key for Scaleway API operations.",
			},
			"token": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"SCALEWAY_TOKEN",
					"SCALEWAY_ACCESS_KEY",
				}, nil),
				Description: "The API key for Scaleway API operations.",
			},
			"organization": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("SCALEWAY_ORGANIZATION", nil),
				Description: "The Organization ID (a.k.a. 'access key') for Scaleway API operations.",
			},
			"region": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SCALEWAY_REGION", "par1"),
				Description: "The Scaleway API region to use.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"user_data":           resourceScalewayUserData(),
			"server":              resourceScalewayServer(),
			"token":               resourceScalewayToken(),
			"ssh_key":             resourceScalewaySSHKey(),
			"ip":                  resourceScalewayIP(),
			"security_group":      resourceScalewaySecurityGroup(),
			"security_group_rule": resourceScalewaySecurityGroupRule(),
			"volume":              resourceScalewayVolume(),
			"volume_attachment":   resourceScalewayVolumeAttachment(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"bootscript":     dataSourceScalewayBootscript(),
			"image":          dataSourceScalewayImage(),
			"security_group": dataSourceScalewaySecurityGroup(),
			"volume":         dataSourceScalewayVolume(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	apiKey := ""
	if v, ok := d.Get("token").(string); ok {
		apiKey = v
	} else {
		if v, ok := d.Get("access_key").(string); ok {
			apiKey = v
		}
	}

	config := Config{
		Organization: d.Get("organization").(string),
		APIKey:       apiKey,
		Region:       d.Get("region").(string),
	}

	return config.Client()
}
