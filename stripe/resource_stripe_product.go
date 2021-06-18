package stripe

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/client"
)

func resourceStripeProduct() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceStripeProductRead,
		CreateContext: resourceStripeProductCreate,
		UpdateContext: resourceStripeProductUpdate,
		DeleteContext: resourceStripeProductDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"active": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"images": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceStripeProductRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	p, err := c.Products.Get(d.Id(), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return CallSet(
		d.Set("name", p.Name),
		d.Set("active", p.Active),
		d.Set("description", p.Description),
		d.Set("images", p.Images),
		d.Set("url", p.URL),
		d.Set("metadata", p.Metadata),
	)
}

func resourceStripeProductCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)

	params := &stripe.ProductParams{
		Name:   stripe.String(String(d, "name")),
		Active: stripe.Bool(Bool(d, "active")),
	}

	description := String(d, "description")
	if len(description) > 0 {
		params.Description = stripe.String(description)
	}

	images := StringSlice(d, "images")
	if len(images) > 0 {
		params.Images = stripe.StringSlice(images)
	}

	u := String(d, "url")
	if len(u) > 0 {
		params.URL = stripe.String(u)
	}

	for k, v := range Map(d, "metadata") {
		params.AddMetadata(k, v.(string))
	}

	product, err := c.Products.New(params)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(product.ID)
	return resourceStripeProductRead(ctx, d, m)
}

func resourceStripeProductUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)

	params := &stripe.ProductParams{}

	if d.HasChange("name") {
		params.Name = stripe.String(String(d, "name"))
	}

	if d.HasChange("active") {
		params.Active = stripe.Bool(Bool(d, "active"))
	}

	if d.HasChange("description") {
		params.Description = stripe.String(String(d, "description"))
	}

	if d.HasChange("images") {
		params.Images = stripe.StringSlice(StringSlice(d, "images"))
	}

	if d.HasChange("url") {
		params.URL = stripe.String(String(d, "url"))
	}

	if d.HasChange("metadata") {
		params.Metadata = nil
		metadata := Map(d, "metadata")
		for k, v := range metadata {
			params.AddMetadata(k, v.(string))
		}
	}

	_, err := c.Products.Update(d.Id(), params)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceStripeProductRead(ctx, d, m)
}

func resourceStripeProductDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)

	_, err := c.Products.Del(d.Id(), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
