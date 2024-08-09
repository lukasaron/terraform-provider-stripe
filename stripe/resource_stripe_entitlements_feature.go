package stripe

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/client"
)

func resourceStripeEntitlementsFeature() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceStripeEntitlementsFeatureRead,
		CreateContext: resourceStripeEntitlementsFeatureCreate,
		UpdateContext: resourceStripeEntitlementsFeatureUpdate,
		DeleteContext: resourceStripeEntitlementsFeatureDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifier for the object.",
			},
			"lookup_key": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "A unique key you provide as your own system identifier. This may be up to 80 characters.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The feature’s name, for your own purpose, not meant to be displayable to the customer.",
			},
			"object": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "String representing the object’s type. Objects of the same type share the same value.",
			},
			"active": {
				Type:     schema.TypeBool,
				Computed: true,
				Description: "Inactive features cannot be attached to new products and will not be returned from " +
					"the features list endpoint.",
			},
			"livemode": {
				Type:     schema.TypeBool,
				Computed: true,
				Description: "Has the value true if the object exists in live mode or the value false " +
					"if the object exists in test mode",
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

func resourceStripeEntitlementsFeatureRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	var entitlementsFeature *stripe.EntitlementsFeature
	var err error

	err = retryWithBackOff(func() error {
		entitlementsFeature, err = c.EntitlementsFeatures.Get(d.Id(), nil)
		return err
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return CallSet(
		d.Set("lookup_key", entitlementsFeature.LookupKey),
		d.Set("name", entitlementsFeature.Name),
		d.Set("active", entitlementsFeature.Active),
		d.Set("object", entitlementsFeature.Object),
		d.Set("livemode", entitlementsFeature.Livemode),
		d.Set("metadata", entitlementsFeature.Metadata),
	)
}

func resourceStripeEntitlementsFeatureCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	var entitlementsFeature *stripe.EntitlementsFeature
	var err error

	params := &stripe.EntitlementsFeatureParams{
		LookupKey: stripe.String(ExtractString(d, "lookup_key")),
		Name:      stripe.String(ExtractString(d, "name")),
	}

	if meta, set := d.GetOk("metadata"); set {
		for k, v := range ToMap(meta) {
			params.AddMetadata(k, ToString(v))
		}
	}

	err = retryWithBackOff(func() error {
		entitlementsFeature, err = c.EntitlementsFeatures.New(params)
		return err
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(entitlementsFeature.ID)
	return resourceStripeEntitlementsFeatureRead(ctx, d, m)
}

func resourceStripeEntitlementsFeatureUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	var err error

	params := &stripe.EntitlementsFeatureParams{}

	if d.HasChange("name") {
		params.Name = stripe.String(ExtractString(d, "name"))
	}

	if d.HasChange("metadata") {
		params.Metadata = nil
		UpdateMetadata(d, params)
	}

	err = retryWithBackOff(func() error {
		_, err = c.EntitlementsFeatures.Update(d.Id(), params)
		return err
	})
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceStripeEntitlementsFeatureRead(ctx, d, m)
}

func resourceStripeEntitlementsFeatureDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[WARN] Stripe SDK doesn't support Entitlements Feature deletion through API! " +
		"Entitlements Feature will be deactivated but not deleted")

	c := m.(*client.API)
	var err error

	params := stripe.EntitlementsFeatureParams{
		Active: stripe.Bool(false),
	}

	err = retryWithBackOff(func() error {
		_, err = c.EntitlementsFeatures.Update(d.Id(), &params)
		return err
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
