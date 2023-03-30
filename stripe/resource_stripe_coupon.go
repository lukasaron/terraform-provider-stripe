package stripe

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/client"
)

func resourceStripeCoupon() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceStripeCouponRead,
		CreateContext: resourceStripeCouponCreate,
		UpdateContext: resourceStripeCouponUpdate,
		DeleteContext: resourceStripeCouponDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"coupon_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Unique string of your choice that will be used to identify this coupon when applying it to a customer. If you don’t want to specify a particular code, you can leave the ID blank and we’ll generate a random code for you.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Name of the coupon displayed to customers on for instance invoices or receipts.",
			},
			"amount_off": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"percent_off"},
				Description: "Amount (in the currency specified) that will be taken off the subtotal of any invoices " +
					"for this customer.",
			},
			"currency": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: "If amount_off has been set, " +
					"the three-letter ISO code for the currency of the amount to take off.",
			},
			"percent_off": {
				Type:          schema.TypeFloat,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"amount_off", "currency"},
				Description: "Percent that will be taken off the subtotal of any invoices for this customer " +
					"for the duration of the coupon. " +
					"For example, a coupon with percent_off of 50 will make a $100 invoice $50 instead.",
			},
			"duration": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "once",
				Description: "One of forever, once, and repeating. " +
					"Describes how long a customer who applies this coupon will get the discount.",
			},
			"duration_in_months": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Description: "If duration is repeating, the number of months the coupon applies. " +
					"Null if coupon duration is forever or once.",
			},
			"max_redemptions": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Description: "Maximum number of times this coupon can be redeemed, " +
					"in total, across all customers, before it is no longer valid.",
			},
			"redeem_by": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Date after which the coupon can no longer be redeemed. Expected format is RFC3339",
			},
			"times_redeemed": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Number of times this coupon has been applied to a customer.",
			},
			"applies_to": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "A list of product IDs this coupon applies to",
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Set of key-value pairs that you can attach to an object. " +
					"This can be useful for storing additional information about the object in a structured format.",
			},
			"valid": {
				Type:     schema.TypeBool,
				Computed: true,
				Description: "Taking account of the above properties, " +
					"whether this coupon can still be applied to a customer.",
			},
		},
	}
}

func resourceStripeCouponCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	params := &stripe.CouponParams{}

	if couponID, set := d.GetOk("coupon_id"); set {
		params.ID = stripe.String(ToString(couponID))
	}
	if name, set := d.GetOk("name"); set {
		params.Name = stripe.String(ToString(name))
	}
	if amountOff, set := d.GetOk("amount_off"); set {
		params.AmountOff = stripe.Int64(ToInt64(amountOff))
	}
	if currency, set := d.GetOk("currency"); set {
		params.Currency = stripe.String(ToString(currency))
	}
	if percentOff, set := d.GetOk("percent_off"); set {
		params.PercentOff = stripe.Float64(ToFloat64(percentOff))
	}
	if duration, set := d.GetOk("duration"); set {
		params.Duration = stripe.String(ToString(duration))
	}
	if durationInMonths, set := d.GetOk("duration_in_months"); set {
		params.DurationInMonths = stripe.Int64(ToInt64(durationInMonths))
	}
	if maxRedemptions, set := d.GetOk("max_redemptions"); set {
		params.MaxRedemptions = stripe.Int64(ToInt64(maxRedemptions))
	}
	if redeemBy, set := d.GetOk("redeem_by"); set {
		t, err := time.Parse(time.RFC3339, ToString(redeemBy))
		if err != nil {
			return diag.FromErr(err)
		}
		params.RedeemBy = stripe.Int64(t.Unix())
	}
	if appliesTo, set := d.GetOk("applies_to"); set {
		params.AppliesTo = &stripe.CouponAppliesToParams{
			Products: stripe.StringSlice(ToStringSlice(appliesTo)),
		}
	}
	if d.HasChange("metadata") {
		params.Metadata = nil
		UpdateMetadata(d, params)
	}

	coupon, err := c.Coupons.New(params)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(coupon.ID)
	return resourceStripeCouponRead(ctx, d, m)
}

func resourceStripeCouponRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	p := &stripe.CouponParams{}
	p.AddExpand("applies_to")

	coupon, err := c.Coupons.Get(d.Id(), p)
	if err != nil {
		return diag.FromErr(err)
	}

	var appliesTo []string
	if coupon.AppliesTo != nil {
		appliesTo = coupon.AppliesTo.Products
	}

	var redeemBy string
	if coupon.RedeemBy != 0 {
		redeemBy = time.Unix(coupon.RedeemBy, 0).Format(time.RFC3339)
	}

	return CallSet(
		d.Set("coupon_id", coupon.ID),
		d.Set("name", coupon.Name),
		d.Set("amount_off", coupon.AmountOff),
		d.Set("currency", coupon.Currency),
		d.Set("percent_off", coupon.PercentOff),
		d.Set("duration", coupon.Duration),
		d.Set("duration_in_months", coupon.DurationInMonths),
		d.Set("max_redemptions", coupon.MaxRedemptions),
		d.Set("redeem_by", redeemBy),
		d.Set("times_redeemed", coupon.TimesRedeemed),
		d.Set("applies_to", appliesTo),
		d.Set("metadata", coupon.Metadata),
		d.Set("valid", coupon.Valid),
	)
}

func resourceStripeCouponUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	params := &stripe.CouponParams{}

	if d.HasChange("name") {
		params.Name = stripe.String(ExtractString(d, "name"))
	}
	if d.HasChange("metadata") {
		params.Metadata = nil
		metadata := ExtractMap(d, "metadata")
		for k, v := range metadata {
			params.AddMetadata(k, ToString(v))
		}
	}

	_, err := c.Coupons.Update(d.Id(), params)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceStripeCouponRead(ctx, d, m)
}

func resourceStripeCouponDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)

	_, err := c.Coupons.Del(d.Id(), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
