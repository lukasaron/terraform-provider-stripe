package stripe

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/client"
)

func resourceStripePromotionCode() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceStripePromotionCodeRead,
		CreateContext: resourceStripePromotionCodeCreate,
		UpdateContext: resourceStripePromotionCodeUpdate,
		DeleteContext: resourceStripePromotionCodeDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifier for the object.",
			},
			"coupon": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The coupon for this promotion code.",
			},
			"code": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: "The customer-facing code. Regardless of case, " +
					"this code must be unique across all active promotion codes for a specific customer. " +
					"If left blank, we will generate one automatically.",
			},
			"active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the promotion code is currently active.",
			},
			"customer": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: "The customer that this promotion code can be used by. " +
					"If not set, the promotion code can be used by all customers.",
			},
			"max_redemptions": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Description: "A positive integer specifying the number of times the promotion code can be redeemed. " +
					"If the coupon has specified a max_redemptions, " +
					"then this value cannot be greater than the coupon’s max_redemptions.",
			},
			"expires_at": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: "The timestamp at which this promotion code will expire. " +
					"If the coupon has specified a redeems_by, " +
					"then this value cannot be after the coupon’s redeems_by. Expected format is RFC3339",
			},
			"restrictions": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Settings that restrict the redemption of the promotion code.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"first_time_transaction": {
							Type:     schema.TypeBool,
							Required: true,
							Description: "A Boolean indicating if the Promotion Code should only be " +
								"redeemed for Customers without any successful payments or invoices",
						},
						"minimum_amount": {
							Type:     schema.TypeInt,
							Required: true,
							Description: "Minimum amount required to redeem this Promotion Code into a Coupon " +
								"(e.g., a purchase must be $100 or more to work).",
						},
						"minimum_amount_currency": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Three-letter ISO code for minimum_amount",
						},
					},
				},
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

func resourceStripePromotionCodeCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	params := &stripe.PromotionCodeParams{
		Coupon: stripe.String(ExtractString(d, "coupon")),
		Active: stripe.Bool(ExtractBool(d, "active")),
	}

	if code, set := d.GetOk("code"); set {
		params.Code = stripe.String(ToString(code))
	}
	if customer, set := d.GetOk("customer"); set {
		params.Customer = stripe.String(ToString(customer))
	}
	if maxRedemptions, set := d.GetOk("max_redemptions"); set {
		params.MaxRedemptions = stripe.Int64(ToInt64(maxRedemptions))
	}
	if expiresAt, set := d.GetOk("expires_at"); set {
		t, err := time.Parse(time.RFC3339, ToString(expiresAt))
		if err != nil {
			return diag.FromErr(err)
		}
		params.ExpiresAt = stripe.Int64(t.Unix())
	}
	if restrictions, set := d.GetOk("restrictions"); set {
		params.Restrictions = &stripe.PromotionCodeRestrictionsParams{}
		for k, v := range ToMap(restrictions) {
			switch k {
			case "first_time_transaction":
				params.Restrictions.FirstTimeTransaction = stripe.Bool(ToBool(v))
			case "minimum_amount":
				params.Restrictions.MinimumAmount = stripe.Int64(ToInt64(v))
			case "minimum_amount_currency":
				params.Restrictions.MinimumAmountCurrency = stripe.String(ToString(v))
			}
		}
	}
	if meta, set := d.GetOk("metadata"); set {
		for k, v := range ToMap(meta) {
			params.AddMetadata(k, ToString(v))
		}
	}

	promotionCode, err := c.PromotionCodes.New(params)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(promotionCode.ID)
	return resourceStripePromotionCodeRead(ctx, d, m)
}

func resourceStripePromotionCodeRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	promotionCode, err := c.PromotionCodes.Get(d.Id(), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return CallSet(
		func() error {
			if promotionCode.Coupon != nil {
				return d.Set("coupon", promotionCode.Coupon.ID)
			}
			return nil
		}(),
		d.Set("code", promotionCode.Code),
		d.Set("active", promotionCode.Active),
		func() error {
			if promotionCode.Customer != nil {
				return d.Set("customer", promotionCode.Customer.ID)
			}
			return nil
		}(),
		d.Set("max_redemptions", promotionCode.MaxRedemptions),
		d.Set("expires_at", time.Unix(promotionCode.ExpiresAt, 0).Format(time.RFC3339)),
		func() error {
			if promotionCode.Restrictions != nil {
				return d.Set("restrictions", []map[string]interface{}{
					{
						"first_time_transaction":  promotionCode.Restrictions.FirstTimeTransaction,
						"minimum_amount":          promotionCode.Restrictions.MinimumAmount,
						"minimum_amount_currency": promotionCode.Restrictions.MinimumAmountCurrency,
					},
				})
			}
			return nil
		}(),
		d.Set("metadata", promotionCode.Metadata),
	)
}

func resourceStripePromotionCodeUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	params := &stripe.PromotionCodeParams{}
	if d.HasChange("active") {
		params.Active = stripe.Bool(ExtractBool(d, "active"))
	}
	if d.HasChange("metadata") {
		params.Metadata = nil
		metadata := ExtractMap(d, "metadata")
		for k, v := range metadata {
			params.AddMetadata(k, ToString(v))
		}
	}
	_, err := c.PromotionCodes.Update(d.Id(), params)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceStripePromotionCodeRead(ctx, d, m)
}

func resourceStripePromotionCodeDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	log.Println("[WARN] Stripe SDK doesn't support Promotion Code deletion through API!")
	d.SetId("")
	return nil
}
