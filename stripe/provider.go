package stripe

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stripe/stripe-go/v72/client"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:        schema.TypeString,
				Description: "The Stripe secret API key",
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("STRIPE_API_KEY", nil),
			},
		},
		ResourcesMap:         nil,
		DataSourcesMap:       map[string]*schema.Resource{},
		ProviderMetaSchema:   nil,
		ConfigureContextFunc: providerConfigure,
		TerraformVersion:     "",
	}
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	log.Print("configuration called")
	time.Sleep(time.Second * 30)
	key := d.Get("api_key").(string)
	return client.New(key, nil), nil
}
