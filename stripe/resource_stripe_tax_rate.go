package stripe

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/client"
)

func resourceStripeTaxRate() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceStripeTaxRateRead,
		CreateContext: resourceStripeTaxRateCreate,
		UpdateContext: resourceStripeTaxRateUpdate,
		DeleteContext: resourceStripeTaxRateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifier for the object.",
			},
			"display_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The display name of the tax rate, which will be shown to users.",
			},
			"inclusive": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: "This specifies if the tax rate is inclusive or exclusive.",
			},
			"percentage": {
				Type:        schema.TypeFloat,
				Required:    true,
				Description: "This represents the tax rate percent out of 100.",
			},
			"active": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				Description: "Flag determining whether the tax rate is active or inactive (archived). " +
					"Inactive tax rates cannot be used with new applications or Checkout Sessions, " +
					"but will still work for subscriptions and invoices that already have it set.",
			},
			"country": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Two-letter country code (ISO 3166-1 alpha-2).",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "An arbitrary string attached to the tax rate for your internal use only. " +
					"It will not be visible to your customers.",
			},
			"jurisdiction": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "The jurisdiction for the tax rate. " +
					"You can use this label field for tax reporting purposes." +
					"It also appears on your customer’s invoice.",
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Set of key-value pairs that you can attach to an object. " +
					"This can be useful for storing additional information about the object in a structured format. " +
					"Individual keys can be unset by posting an empty value to them. " +
					"All keys can be unset by posting an empty value to metadata.",
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "ISO 3166-2 subdivision code, without country prefix. " +
					"For example, “NY” for New York, United States.",
			},
			"tax_type": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "The high-level tax type, " +
					"such as vat or sales_tax.",
			},
			"object": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "String representing the object’s type. " +
					"Objects of the same type share the same value.",
			},
			"created": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: "Time at which the object was created. " +
					"Measured in seconds since the Unix epoch.",
			},
			"livemode": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Has the value true if the object exists in live mode or the value false if the object exists in test mode.",
			},
		},
	}
}

func resourceStripeTaxRateRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	var taxRate *stripe.TaxRate
	var err error

	err = retryWithBackOff(func() error {
		taxRate, err = c.TaxRates.Get(d.Id(), nil)
		return err
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return CallSet(
		d.Set("object", taxRate.Object),
		d.Set("active", taxRate.Active),
		d.Set("country", taxRate.Country),
		d.Set("created", taxRate.Created),
		d.Set("description", taxRate.Description),
		d.Set("display_name", taxRate.DisplayName),
		d.Set("inclusive", taxRate.Inclusive),
		d.Set("jurisdiction", taxRate.Jurisdiction),
		d.Set("livemode", taxRate.Livemode),
		d.Set("metadata", taxRate.Metadata),
		d.Set("percentage", taxRate.Percentage),
		d.Set("state", taxRate.State),
		d.Set("tax_type", taxRate.TaxType),
	)
}

func resourceStripeTaxRateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	var taxRate *stripe.TaxRate
	var err error

	params := &stripe.TaxRateParams{
		DisplayName: stripe.String(ExtractString(d, "display_name")),
		Inclusive:   stripe.Bool(ExtractBool(d, "inclusive")),
		Percentage:  stripe.Float64(ExtractFloat64(d, "percentage")),
	}

	if active, set := d.GetOk("active"); set {
		params.Active = stripe.Bool(ToBool(active))
	}

	if country, set := d.GetOk("country"); set {
		params.Country = stripe.String(ToString(country))
	}

	if description, set := d.GetOk("description"); set {
		params.Description = stripe.String(ToString(description))
	}

	if jurisdiction, set := d.GetOk("jurisdiction"); set {
		params.Jurisdiction = stripe.String(ToString(jurisdiction))
	}

	if meta, set := d.GetOk("metadata"); set {
		for k, v := range ToMap(meta) {
			params.AddMetadata(k, ToString(v))
		}
	}

	if state, set := d.GetOk("state"); set {
		params.State = stripe.String(ToString(state))
	}

	if taxType, set := d.GetOk("tax_type"); set {
		params.TaxType = stripe.String(ToString(taxType))
	}

	err = retryWithBackOff(func() error {
		taxRate, err = c.TaxRates.New(params)
		return err
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(taxRate.ID)
	return resourceStripeTaxRateRead(ctx, d, m)
}

func resourceStripeTaxRateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	var err error

	params := &stripe.TaxRateParams{}

	if d.HasChange("active") {
		params.Active = stripe.Bool(ExtractBool(d, "active"))
	}

	if d.HasChange("country") {
		params.Country = stripe.String(ExtractString(d, "country"))
	}

	if d.HasChange("description") {
		params.Description = stripe.String(ExtractString(d, "description"))
	}
	if d.HasChange("display_name") {
		params.DisplayName = stripe.String(ExtractString(d, "display_name"))
	}
	if d.HasChange("jurisdiction") {
		params.Jurisdiction = stripe.String(ExtractString(d, "jurisdiction"))
	}
	if d.HasChange("state") {
		params.State = stripe.String(ExtractString(d, "state"))
	}
	if d.HasChange("metadata") {
		params.Metadata = nil
		UpdateMetadata(d, params)
	}
	if d.HasChange("tax_type") {
		params.TaxType = stripe.String(ExtractString(d, "tax_type"))
	}

	err = retryWithBackOff(func() error {
		_, err = c.TaxRates.Update(d.Id(), params)
		return err
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceStripeTaxRateRead(ctx, d, m)
}

func resourceStripeTaxRateDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	log.Println("[WARN] Stripe SDK doesn't support Tax Rate deletion through API!")
	d.SetId("")
	return nil
}
