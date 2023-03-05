package stripe

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/client"
)

func resourceStripePrice() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceStripePriceRead,
		CreateContext: resourceStripePriceCreate,
		UpdateContext: resourceStripePriceUpdate,
		DeleteContext: resourceStripePriceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifier for the object.",
			},
			"currency": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Three-letter ISO currency code, in lowercase.",
			},
			"product": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the product that this price will belong to.",
			},
			"unit_amount": {
				Type:          schema.TypeInt,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"unit_amount_decimal"},
				Description:   "A positive integer in cents (or -1 for a free price) representing how much to charge.",
			},
			"unit_amount_decimal": {
				Type:          schema.TypeFloat,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"unit_amount"},
				Description: "Same as unit_amount, " +
					"but accepts a decimal value in cents with at most 12 decimal places. " +
					"Only one of unit_amount and unit_amount_decimal can be set",
			},
			"active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the price can be used for new purchases. Defaults to true.",
			},
			"nickname": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A brief description of the price, hidden from customers.",
			},
			"recurring": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "The recurring components of a price such as interval and usage_type.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"interval": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies billing frequency. Either day, week, month or year.",
						},
						"aggregate_usage": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "Specifies a usage aggregation strategy for prices of usage_type=metered. " +
								"Allowed values are sum for summing up all usage during a period, " +
								"last_during_period for using the last usage record reported within a period, " +
								"last_ever for using the last usage record ever (across period bounds) or max which " +
								"uses the usage record with the maximum reported usage during a period. ",
						},
						"interval_count": {
							Type:     schema.TypeInt,
							Optional: true,
							Description: "The number of intervals between subscription billings. " +
								"For example, interval=month and interval_count=3 bills every 3 months. " +
								"Maximum of one year interval allowed (1 year, 12 months, or 52 weeks).",
						},
						"usage_type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "licensed",
							Description: "Configures how the quantity per period should be determined. " +
								"Can be either metered or licensed. licensed automatically bills the quantity " +
								"set when adding it to a subscription. metered aggregates the total usage " +
								"based on usage records. Defaults to licensed.",
						},
					},
				},
			},
			"tiers": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Description: "Each element represents a pricing tier. " +
					"This parameter requires billing_scheme to be set to tiered. " +
					"See also the documentation for billing_scheme.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"up_to": {
							Type:     schema.TypeInt,
							Optional: true,
							Description: "Specifies the upper bound of this tier. " +
								"The lower bound of a tier is the upper bound of the previous tier adding one. " +
								"Use -1 to define a fallback tier.",
						},
						"flat_amount": {
							Type:     schema.TypeInt,
							Optional: true,
							Description: "The flat billing amount for an entire tier, " +
								"regardless of the number of units in the tier.",
						},
						"flat_amount_decimal": {
							Type:     schema.TypeFloat,
							Optional: true,
							Description: "Same as flat_amount, but accepts a decimal value representing an integer " +
								"in the minor units of the currency. " +
								"Only one of flat_amount and flat_amount_decimal can be set.",
						},
						"unit_amount": {
							Type:     schema.TypeInt,
							Optional: true,
							Description: "The per unit billing amount for each individual unit " +
								"for which this tier applies.",
						},
						"unit_amount_decimal": {
							Type:     schema.TypeFloat,
							Optional: true,
							Description: "Same as unit_amount, but accepts a decimal value in cents with " +
								"at most 12 decimal places. " +
								"Only one of unit_amount and unit_amount_decimal can be set.",
						},
					},
				},
			},
			"tiers_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: "Defines if the tiering price should be graduated or volume based. " +
					"In volume-based tiering, the maximum quantity within a period determines the per unit price, " +
					"in graduated tiering pricing can successively change as the quantity grows.",
			},
			"billing_scheme": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Description: "Describes how to compute the price per period. " +
					"Either per_unit or tiered. per_unit indicates that the fixed amount " +
					"(specified in unit_amount or unit_amount_decimal) will be charged per unit in quantity " +
					"(for prices with usage_type=licensed), or per unit of total usage " +
					"(for prices with usage_type=metered). " +
					"tiered indicates that the unit pricing will be computed using a tiering strategy " +
					"as defined using the tiers and tiers_mode attributes.",
			},
			"lookup_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A lookup key used to retrieve prices dynamically from a static string.",
			},
			"transfer_lookup_key": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				Description: "If set to true, will atomically remove the lookup key from the existing price, " +
					"and assign it to this price.",
			},
			"tax_behaviour": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Specifies whether the price is considered inclusive of taxes or exclusive of taxes. " +
					"One of inclusive, exclusive, or unspecified. " +
					"Once specified as either inclusive or exclusive, it cannot be changed.",
			},
			"transform_quantity": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Description: "Apply a transformation to the reported usage or set quantity " +
					"before computing the billed price. Cannot be combined with tiers",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"divide_by": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Divide usage by this number.",
						},
						"round": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "After division, either round the result up or down",
						},
					},
				},
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "One of one_time or recurring depending on whether the price is for a one-time purchase " +
					"or a recurring (subscription) purchase",
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

func resourceStripePriceRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	price, err := c.Prices.Get(d.Id(), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return CallSet(
		d.Set("currency", price.Currency),
		d.Set("product", price.Product.ID),
		func() error {
			// We should only every fall into one of these if statements as we have a conflict between them on the schema
			if d.HasChange("unit_amount") && !d.HasChange("unit_amount_decimal") {
				if err = d.Set("unit_amount", price.UnitAmount); err != nil {
					return err
				}
			}
			if d.HasChange("unit_amount_decimal") && !d.HasChange("unit_amount") {
				if err = d.Set("unit_amount_decimal", price.UnitAmountDecimal); err != nil {
					return err
				}
			}

			return nil
		}(),
		d.Set("active", price.Active),
		d.Set("nickname", price.Nickname),
		func() error {
			if price.Recurring != nil {
				return d.Set("recurring", []map[string]interface{}{
					{
						"interval":        price.Recurring.Interval,
						"aggregate_usage": price.Recurring.AggregateUsage,
						"interval_count":  price.Recurring.IntervalCount,
						"usage_type":      price.Recurring.UsageType,
					},
				})
			}
			return nil
		}(),
		func() error {
			if len(price.Tiers) > 0 {
				var tiers []map[string]interface{}
				for _, tier := range price.Tiers {
					t := map[string]interface{}{
						"up_to": func() int64 {
							// update the value to reflect the Terraform input
							if tier.UpTo == 0 {
								return -1
							}
							return tier.UpTo
						}(),
						"flat_amount":         tier.FlatAmount,
						"flat_amount_decimal": tier.FlatAmountDecimal,
						"unit_amount":         tier.UnitAmount,
						"unit_amount_decimal": tier.UnitAmountDecimal,
					}
					tiers = append(tiers, t)
				}
				return d.Set("tiers", tiers)
			}
			return nil
		}(),
		d.Set("tiers_mode", price.TiersMode),
		d.Set("billing_scheme", price.BillingScheme),
		d.Set("lookup_key", price.LookupKey),
		d.Set("tax_behaviour", price.TaxBehavior),
		func() error {
			if price.TransformQuantity != nil {
				return d.Set("transform_quantity", []map[string]interface{}{
					{
						"divide_by": price.TransformQuantity.DivideBy,
						"round":     price.TransformQuantity.Round,
					},
				})
			}
			return nil
		}(),
		d.Set("type", price.Type),
		d.Set("metadata", price.Metadata),
	)
}

func resourceStripePriceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	params := &stripe.PriceParams{
		Product:  stripe.String(ExtractString(d, "product")),
		Currency: stripe.String(ExtractString(d, "currency")),
		Active:   stripe.Bool(ExtractBool(d, "active")),
	}

	if unitAmount, set := d.GetOk("unit_amount"); set {
		amount := ToInt64(unitAmount)
		// amount is -1 when free price is required
		if amount < 0 {
			amount = 0
		}
		params.UnitAmount = stripe.Int64(amount)
	}
	if unitAmountDecimal, set := d.GetOk("unit_amount_decimal"); set {
		params.UnitAmountDecimal = stripe.Float64(ToFloat64(unitAmountDecimal))
	}
	if nickname, set := d.GetOk("nickname"); set {
		params.Nickname = stripe.String(ToString(nickname))
	}
	if recurring, set := d.GetOk("recurring"); set {
		params.Recurring = &stripe.PriceRecurringParams{}
		for k, v := range ToMap(recurring) {
			switch {
			case k == "interval" && ToString(v) != "":
				params.Recurring.Interval = stripe.String(ToString(v))
			case k == "aggregate_usage" && ToString(v) != "":
				params.Recurring.AggregateUsage = stripe.String(ToString(v))
			case k == "interval_count" && ToInt64(v) != 0:
				params.Recurring.IntervalCount = stripe.Int64(ToInt64(v))
			case k == "usage_type" && ToString(v) != "":
				params.Recurring.UsageType = stripe.String(ToString(v))
			}
		}
	}
	if tiers, set := d.GetOk("tiers"); set {
		for _, t := range ToSlice(tiers) {
			priceTier := &stripe.PriceTierParams{}
			for k, v := range ToMap(t) {
				switch {
				case k == "up_to" && ToInt64(v) != 0:
					upTo := ToInt64(v)
					if upTo < 0 {
						priceTier.UpToInf = stripe.Bool(true)
					} else {
						priceTier.UpTo = stripe.Int64(ToInt64(v))
					}
				case k == "flat_amount":
					priceTier.FlatAmount = OptionalInt64(v)
				case k == "flat_amount_decimal":
					priceTier.FlatAmountDecimal = OptionalFloat64(v)
				case k == "unit_amount":
					priceTier.UnitAmount = OptionalInt64(v)
				case k == "unit_amount_decimal":
					priceTier.UnitAmountDecimal = OptionalFloat64(v)
				}
			}
			params.Tiers = append(params.Tiers, priceTier)
		}
	}
	if tiersMode, set := d.GetOk("tiers_mode"); set {
		params.TiersMode = stripe.String(ToString(tiersMode))
	}
	if billingScheme, set := d.GetOk("billing_scheme"); set {
		params.BillingScheme = stripe.String(ToString(billingScheme))
	}
	if lookupKey, set := d.GetOk("lookup_key"); set {
		params.LookupKey = stripe.String(ToString(lookupKey))
	}
	if transferLookupKey, set := d.GetOk("transfer_lookup_key"); set {
		params.TransferLookupKey = stripe.Bool(ToBool(transferLookupKey))
	}
	if taxBehaviour, set := d.GetOk("tax_behaviour"); set {
		params.TaxBehavior = stripe.String(ToString(taxBehaviour))
	}
	if transformQuantity, set := d.GetOk("transform_quantity"); set {
		params.TransformQuantity = &stripe.PriceTransformQuantityParams{}
		for k, v := range ToMap(transformQuantity) {
			switch k {
			case "divide_by":
				params.TransformQuantity.DivideBy = stripe.Int64(ToInt64(v))
			case "round":
				params.TransformQuantity.Round = stripe.String(ToString(v))
			}
		}
	}
	if meta, set := d.GetOk("metadata"); set {
		for k, v := range ToMap(meta) {
			params.AddMetadata(k, ToString(v))
		}
	}

	price, err := c.Prices.New(params)
	if err != nil {
		return diag.FromErr(err)
	}

	dg := CallSet(func() error {
		if params.TransferLookupKey != nil {
			return d.Set("transfer_lookup_key", *params.TransferLookupKey)
		}
		return nil
	}())
	if len(dg) > 0 {
		return dg
	}

	d.SetId(price.ID)
	return resourceStripePriceRead(ctx, d, m)
}

func resourceStripePriceUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	params := &stripe.PriceParams{}

	if d.HasChange("active") {
		params.Active = stripe.Bool(ExtractBool(d, "active"))
	}
	if d.HasChange("nickname") {
		params.Nickname = stripe.String(ExtractString(d, "nickname"))
	}
	if d.HasChange("lookup_key") {
		params.LookupKey = stripe.String(ExtractString(d, "lookup_key"))
	}
	if d.HasChange("transfer_lookup_key") {
		params.TransferLookupKey = stripe.Bool(ExtractBool(d, "transfer_lookup_key"))
	}
	if d.HasChange("tax_behaviour") {
		params.TaxBehavior = stripe.String(ExtractString(d, "tax_behaviour"))
	}
	if d.HasChange("metadata") {
		params.Metadata = nil
		UpdateMetadata(d, params)
	}

	_, err := c.Prices.Update(d.Id(), params)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceStripePriceRead(ctx, d, m)
}

func resourceStripePriceDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	params := stripe.PriceParams{
		Active: stripe.Bool(false),
	}

	if _, err := c.Prices.Update(d.Id(), &params); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
