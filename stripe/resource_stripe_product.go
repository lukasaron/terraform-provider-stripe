package stripe

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/client"
)

func resourceStripeProduct() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceStripeProductRead,
		CreateContext: resourceStripeProductCreate,
		UpdateContext: resourceStripeProductUpdate,
		DeleteContext: resourceStripeProductDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"product_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
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
				Description: "Whether the product is currently available for purchase. Defaults to true.",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "The product’s description, meant to be displayable to the customer. " +
					"Use this field to optionally store a long form explanation of the product " +
					"being sold for your own rendering purposes.",
			},
			"features": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of up to 15 features for this product. These are displayed in pricing tables. ",
			},
			"images": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "A list of up to 8 URLs of images for this product, " +
					"meant to be displayable to the customer.",
			},
			"package_dimensions": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeFloat},
				Description: "The dimensions of this product for shipping purposes.",
			},
			"shippable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether this product is shipped (i.e., physical goods).",
			},
			"statement_descriptor": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "An arbitrary string to be displayed on your customer’s credit card or bank statement. " +
					"While most banks display this information consistently, " +
					"some may display it incorrectly or not at all. This may be up to 22 characters. " +
					"The statement description may not include <, >, \\, \", ’ characters, " +
					"and will appear on your customer’s statement in capital letters. " +
					"Non-ASCII characters are automatically stripped. It must contain at least one letter.",
			},
			"tax_code": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A tax code ID. Supported values are listed in the TaxCode resource and at https://stripe.com/docs/tax/tax-categories.",
			},
			"unit_label": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "A label that represents units of this product in Stripe and on customers’ " +
					"receipts and invoices. " +
					"When set, this will be included in associated invoice line item descriptions.",
			},
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A URL of a publicly-accessible webpage for this product.",
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
	var product *stripe.Product
	var err error

	err = retryWithBackOff(func() error {
		product, err = c.Products.Get(d.Id(), nil)
		return err
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return CallSet(
		d.Set("product_id", product.ID),
		d.Set("name", product.Name),
		d.Set("active", product.Active),
		d.Set("description", product.Description),
		func() error {
			if len(product.Features) > 0 {
				var features []string
				for _, feature := range product.Features {
					features = append(features, feature.Name)
				}
				return d.Set("features", features)
			}
			return nil
		}(),
		d.Set("images", product.Images),
		func() error {
			if product.PackageDimensions != nil {
				return d.Set("package_dimensions", map[string]interface{}{
					"height": product.PackageDimensions.Height,
					"length": product.PackageDimensions.Length,
					"weight": product.PackageDimensions.Weight,
					"width":  product.PackageDimensions.Width,
				})
			}
			return nil
		}(),
		d.Set("shippable", product.Shippable),
		d.Set("statement_descriptor", product.StatementDescriptor),
		func() error {
			if product.TaxCode != nil {
				return d.Set("tax_code", product.TaxCode.ID)
			}
			return nil
		}(),
		d.Set("unit_label", product.UnitLabel),
		d.Set("url", product.URL),
		d.Set("metadata", product.Metadata),
	)
}

func resourceStripeProductCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	var product *stripe.Product
	var err error

	params := &stripe.ProductParams{
		Name: stripe.String(ExtractString(d, "name")),
	}
	if productID, ok := d.GetOk("product_id"); ok {
		params.ID = stripe.String(ToString(productID))
	}
	if active, set := d.GetOk("active"); set {
		params.Active = stripe.Bool(ToBool(active))
	}
	if description, set := d.GetOk("description"); set {
		params.Description = stripe.String(ToString(description))
	}
	if features, set := d.GetOk("features"); set {
		for _, feature := range ToStringSlice(features) {
			params.Features = append(params.Features, &stripe.ProductFeatureParams{Name: stripe.String(feature)})
		}
	}
	if images, set := d.GetOk("images"); set {
		params.Images = stripe.StringSlice(ToStringSlice(images))
	}
	if packageDimensions, set := d.GetOk("package_dimensions"); set {
		params.PackageDimensions = &stripe.ProductPackageDimensionsParams{}
		dimensions := ToMap(packageDimensions)
		for k, v := range dimensions {
			switch k {
			case "height":
				params.PackageDimensions.Height = stripe.Float64(ToFloat64(v))
			case "length":
				params.PackageDimensions.Length = stripe.Float64(ToFloat64(v))
			case "weight":
				params.PackageDimensions.Weight = stripe.Float64(ToFloat64(v))
			case "width":
				params.PackageDimensions.Width = stripe.Float64(ToFloat64(v))
			}
		}
	}
	if shippable, set := d.GetOk("shippable"); set {
		params.Shippable = stripe.Bool(ToBool(shippable))
	}
	if statementDescriptor, set := d.GetOk("statement_descriptor"); set {
		params.StatementDescriptor = stripe.String(ToString(statementDescriptor))
	}
	if taxCode, set := d.GetOk("tax_code"); set {
		params.TaxCode = stripe.String(ToString(taxCode))
	}
	if unitLabel, set := d.GetOk("unit_label"); set {
		params.UnitLabel = stripe.String(ToString(unitLabel))
	}
	if url, set := d.GetOk("url"); set {
		params.URL = stripe.String(ToString(url))
	}
	if meta, set := d.GetOk("metadata"); set {
		for k, v := range ToMap(meta) {
			params.AddMetadata(k, ToString(v))
		}
	}

	err = retryWithBackOff(func() error {
		product, err = c.Products.New(params)
		return err
	})
	if err != nil {
		return diag.FromErr(err)
	}

	dg := CallSet()
	if len(dg) > 0 {
		return dg
	}

	d.SetId(product.ID)
	return resourceStripeProductRead(ctx, d, m)
}

func resourceStripeProductUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	var err error

	params := &stripe.ProductParams{}
	if d.HasChange("name") {
		params.Name = stripe.String(ExtractString(d, "name"))
	}
	if d.HasChange("active") {
		params.Active = stripe.Bool(ExtractBool(d, "active"))
	}
	if d.HasChange("description") {
		params.Description = stripe.String(ExtractString(d, "description"))
	}
	if d.HasChange("features") {
		for _, feature := range ExtractStringSlice(d, "features") {
			params.Features = append(params.Features, &stripe.ProductFeatureParams{Name: stripe.String(feature)})
		}
	}
	if d.HasChange("images") {
		params.Images = stripe.StringSlice(ExtractStringSlice(d, "images"))
	}
	if d.HasChange("package_dimensions") {
		params.PackageDimensions = &stripe.ProductPackageDimensionsParams{}
		dimensions := ExtractMap(d, "package_dimensions")
		for k, v := range dimensions {
			switch k {
			case "height":
				params.PackageDimensions.Height = stripe.Float64(ToFloat64(v))
			case "length":
				params.PackageDimensions.Length = stripe.Float64(ToFloat64(v))
			case "weight":
				params.PackageDimensions.Weight = stripe.Float64(ToFloat64(v))
			case "width":
				params.PackageDimensions.Width = stripe.Float64(ToFloat64(v))
			}
		}
	}
	if d.HasChange("shippable") {
		params.Shippable = stripe.Bool(ExtractBool(d, "shippable"))
	}
	if d.HasChange("statement_descriptor") {
		params.StatementDescriptor = stripe.String(ExtractString(d, "statement_descriptor"))
	}
	if d.HasChange("tax_code") {
		params.TaxCode = stripe.String(ExtractString(d, "tax_code"))
	}
	if d.HasChange("unit_label") {
		params.UnitLabel = stripe.String(ExtractString(d, "unit_label"))
	}
	if d.HasChange("url") {
		params.URL = stripe.String(ExtractString(d, "url"))
	}
	if d.HasChange("metadata") {
		params.Metadata = nil
		UpdateMetadata(d, params)
	}

	err = retryWithBackOff(func() error {
		_, err = c.Products.Update(d.Id(), params)
		return err
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceStripeProductRead(ctx, d, m)
}

func resourceStripeProductDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	var err error

	err = retryWithBackOff(func() error {
		_, err = c.Products.Del(d.Id(), nil)
		if err != nil {
			stripeErr := toStripeError(err)
			/*
				When ErrorTypeInvalidRequest error is returned on delete endpoint (price resources were attached)
				the product is de-activated (archived).
			*/
			if stripeErr.Type == stripe.ErrorTypeInvalidRequest {
				params := &stripe.ProductParams{Active: stripe.Bool(false)}
				_, err = c.Products.Update(d.Id(), params)
			}
		}
		return err
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
