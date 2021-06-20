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
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifier for the object.",
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				Description: "The product’s name, meant to be displayable to the customer. " +
					"Whenever this product is sold via a subscription, " +
					"name will show up on associated invoice line item descriptions.",
			},
			"active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the product is currently available for purchase.",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "The product’s description, meant to be displayable to the customer. " +
					"Use this field to optionally store a long form explanation of the product being " +
					"sold for your own rendering purposes.",
			},
			"images": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "A list of up to 8 URLs of images for this product, " +
					"meant to be displayable to the customer.",
			},
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A URL of a publicly-accessible webpage for this product.",
			},
			"statement_descriptor": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Extra information about a product which will appear on your customer’s credit " +
					"card statement. In the case that multiple products are billed at once, " +
					"the first statement descriptor will be used.",
			},
			"unit_label": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "A label that represents units of this product in Stripe and on customers’ receipts " +
					"and invoices. When set, this will be included in associated invoice line item descriptions.",
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
		d.Set("statement_descriptor", p.StatementDescriptor),
		d.Set("unit_label", p.UnitLabel),
		d.Set("metadata", p.Metadata),
	)
}

func resourceStripeProductCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)

	params := &stripe.ProductParams{
		Name:   stripe.String(String(d, "name")),
		Active: stripe.Bool(Bool(d, "active")),
	}

	if images, set := d.GetOk("images"); set {
		params.Images = stripe.StringSlice(ToStringSlice(images))
	}
	if description, set := d.GetOk("description"); set {
		params.Description = stripe.String(ToString(description))
	}
	if u, set := d.GetOk("url"); set {
		params.URL = stripe.String(ToString(u))
	}
	if uLabel, set := d.GetOk("unit_label"); set {
		params.UnitLabel = stripe.String(ToString(uLabel))
	}
	if sDescriptor, set := d.GetOk("statement_descriptor"); set {
		params.StatementDescriptor = stripe.String(ToString(sDescriptor))
	}
	if meta, set := d.GetOk("metadata"); set {
		for k, v := range ToMap(meta) {
			params.AddMetadata(k, ToString(v))
		}
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
	if d.HasChange("statement_descriptor") {
		params.StatementDescriptor = stripe.String(String(d, "statement_descriptor"))
	}
	if d.HasChange("unit_label") {
		params.UnitLabel = stripe.String(String(d, "unit_label"))
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
