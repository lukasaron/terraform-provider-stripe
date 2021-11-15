package stripe

import (
	"context"

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
		ResourcesMap: map[string]*schema.Resource{
			"stripe_webhook_endpoint": resourceStripeWebhookEndpoint(),
			"stripe_coupon":           resourceStripeCoupon(),
			"stripe_product":          resourceStripeProduct(),
			"stripe_promotion_code":   resourceStripePromotionCode(),
			"stripe_price":            resourceStripePrice(),
			"stripe_customer":         resourceStripeCustomer(),
			"stripe_tax_rate":         resourceStripeTaxRate(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	key := ExtractString(d, "api_key")
	return client.New(key, nil), nil
}
