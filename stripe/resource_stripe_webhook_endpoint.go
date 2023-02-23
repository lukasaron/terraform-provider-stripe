package stripe

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/client"
)

func resourceStripeWebhookEndpoint() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceStripeWebhookEndpointRead,
		CreateContext: resourceStripeWebhookEndpointCreate,
		UpdateContext: resourceStripeWebhookEndpointUpdate,
		DeleteContext: resourceStripeWebhookEndpointDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifier for the object.",
			},
			"enabled_events": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "The list of events to enable for this endpoint. " +
					"[’*’] indicates that all events are enabled, except those that require explicit selection.",
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The URL of the webhook endpoint.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An optional description of what the webhook is used for.",
			},
			"secret": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The endpoint’s secret, used to generate webhook signatures. Only returned at creation.",
			},
			"disabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Disable the webhook endpoint if set to true.",
			},
			"connect": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
				Description: "Whether this endpoint should receive events from connected accounts (true), " +
					"or from your account (false). Defaults to false",
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Set of key-value pairs that you can attach to an object. " +
					"This can be useful for storing additional information about the object in a structured format.",
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
		d.Set("enabled_events", webhookEndpoint.EnabledEvents),
		d.Set("url", webhookEndpoint.URL),
		d.Set("description", webhookEndpoint.Description),
		d.Set("disabled", webhookEndpoint.Status != "enabled"),
		// TODO revisit this part in the future - now hardcoded the value from the state
		d.Set("connect", ExtractBool(d, "connect")),
		d.Set("metadata", webhookEndpoint.Metadata),
	)
}

func resourceStripeWebhookEndpointCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	params := &stripe.WebhookEndpointParams{
		URL:           stripe.String(ExtractString(d, "url")),
		EnabledEvents: stripe.StringSlice(ExtractStringSlice(d, "enabled_events")),
	}
	if description, set := d.GetOk("description"); set {
		params.Description = stripe.String(ToString(description))
	}
	if connect, set := d.GetOk("connect"); set {
		params.Connect = stripe.Bool(ToBool(connect))
	}
	if meta, set := d.GetOk("metadata"); set {
		for k, v := range ToMap(meta) {
			params.AddMetadata(k, ToString(v))
		}
	}
	webhookEndpoint, err := c.WebhookEndpoints.New(params)
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
		params.EnabledEvents = stripe.StringSlice(ExtractStringSlice(d, "enabled_events"))
	}
	if d.HasChange("url") {
		params.URL = stripe.String(ExtractString(d, "url"))
	}
	if d.HasChange("description") {
		params.Description = stripe.String(ExtractString(d, "description"))
	}
	if d.HasChange("disabled") {
		params.Disabled = stripe.Bool(ExtractBool(d, "disabled"))
	}
	if d.HasChange("metadata") {
		params.Metadata = nil
		metadata := ExtractMap(d, "metadata")
		for k, v := range metadata {
			params.AddMetadata(k, ToString(v))
		}
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
