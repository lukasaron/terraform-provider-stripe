package stripe

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/client"
)

func resourceStripeWebhookEndpoint() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceStripeWebhookEndpointRead,
		CreateContext: resourceStripeWebhookEndpointCreate,
		UpdateContext: resourceStripeWebhookEndpointUpdate,
		DeleteContext: resourceStripeWebhookEndpointDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled_events": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"secret": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"livemode": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceStripeWebhookEndpointRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)

	webhookEndpoint, err := c.WebhookEndpoints.Get(d.Id(), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return CallSet(
		d.Set("id", webhookEndpoint.ID),
		d.Set("enabled_events", webhookEndpoint.EnabledEvents),
		d.Set("url", webhookEndpoint.URL),
		d.Set("description", webhookEndpoint.Description),
		d.Set("status", webhookEndpoint.Status),
		d.Set("created", time.Unix(webhookEndpoint.Created, 0).Format(time.RFC3339)),
		d.Set("livemode", webhookEndpoint.Livemode),
	)
}

func resourceStripeWebhookEndpointCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)

	if Bool(d, "disabled") {
		return diag.FromErr(errors.New("disabled can be set when updating existing webhook only"))
	}

	webhookEndpoint, err := c.WebhookEndpoints.New(&stripe.WebhookEndpointParams{
		URL:           stripe.String(String(d, "url")),
		EnabledEvents: stripe.StringSlice(StringSlice(d, "enabled_events")),
		Description:   stripe.String(String(d, "description")),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	dg := CallSet(
		d.Set("secret", webhookEndpoint.Secret),
	)
	if len(dg) > 0 {
		return dg
	}

	d.SetId(webhookEndpoint.ID)
	return resourceStripeWebhookEndpointRead(ctx, d, m)
}

func resourceStripeWebhookEndpointUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	params := &stripe.WebhookEndpointParams{}

	if d.HasChange("enabled_events") {
		params.EnabledEvents = stripe.StringSlice(StringSlice(d, "enabled_events"))
	}

	if d.HasChange("url") {
		params.URL = stripe.String(String(d, "url"))
	}

	if d.HasChange("description") {
		params.Description = stripe.String(String(d, "description"))
	}

	if d.HasChange("disabled") {
		params.Disabled = stripe.Bool(Bool(d, "disabled"))
	}

	_, err := c.WebhookEndpoints.Update(d.Id(), params)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceStripeWebhookEndpointRead(ctx, d, m)
}

func resourceStripeWebhookEndpointDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	_, err := c.WebhookEndpoints.Del(d.Id(), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
