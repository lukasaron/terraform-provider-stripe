package stripe

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/client"
)

func resourceStripePrice() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceStripePriceRead,
		CreateContext: resourceStripePriceCreate,
		UpdateContext: resourceStripePriceRead,
		DeleteContext: resourceStripePriceDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"currency": {
				Type:     schema.TypeString,
				Required: true,
			},
			"product_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"unit_amount": {
				Type:          schema.TypeInt,
				Optional:      true,
				ConflictsWith: []string{"unit_amount_decimal"},
			},
			"unit_amount_decimal": {
				Type:          schema.TypeFloat,
				Optional:      true,
				ConflictsWith: []string{"unit_amount"},
			},
			"active": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"nickname": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"recurring": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				//Elem: &schema.Resource{
				//	Schema: map[string]*schema.Schema{
				//		"interval": {
				//			Type:     schema.TypeString,
				//			Required: true,
				//			//ExactlyOneOf: []string{
				//			//	string(stripe.PriceRecurringIntervalDay),
				//			//	string(stripe.PriceRecurringIntervalWeek),
				//			//	string(stripe.PriceRecurringIntervalMonth),
				//			//	string(stripe.PriceRecurringIntervalYear),
				//			//},
				//		},
				//		"aggregate_usage": {
				//			Type:     schema.TypeString,
				//			Optional: true,
				//			Default:  string(stripe.PriceRecurringAggregateUsageSum),
				//			//ExactlyOneOf: []string{
				//			//	string(stripe.PriceRecurringAggregateUsageSum),
				//			//	string(stripe.PriceRecurringAggregateUsageLastDuringPeriod),
				//			//	string(stripe.PriceRecurringAggregateUsageLastEver),
				//			//	string(stripe.PriceRecurringAggregateUsageMax),
				//			//},
				//		},
				//		"interval_count": {
				//			Type:     schema.TypeInt,
				//			Optional: true,
				//			Default:  1,
				//		},
				//		"usage_type": {
				//			Type:     schema.TypeString,
				//			Optional: true,
				//			//ExactlyOneOf: []string{
				//			//	string(stripe.PriceRecurringUsageTypeMetered),
				//			//	string(stripe.PriceRecurringUsageTypeLicensed),
				//			//},
				//			Default: string(stripe.PriceRecurringUsageTypeLicensed),
				//		},
				//	},
				//},
			},
			"billing_scheme": {
				Type:     schema.TypeString,
				Optional: true,
				//ExactlyOneOf: []string{
				//	string(stripe.PriceBillingSchemePerUnit),
				//	string(stripe.PriceBillingSchemeTiered),
				//},
				Default: string(stripe.PriceBillingSchemePerUnit),
			},
			"tiers": {
				Type:         schema.TypeList,
				Optional:     true,
				RequiredWith: []string{"tiers_mode"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"up_to": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"flat_amount": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"flat_amount_decimal": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"flat_unit_amount": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"unit_amount_decimal": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
					},
				},
			},
			"tiers_mode": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"tiers"},
				//ExactlyOneOf: []string{
				//	string(stripe.PriceTiersModeVolume),
				//	string(stripe.PriceTiersModeGraduated),
				//},
			},
			"transform_quantity": {
				Type:          schema.TypeSet,
				Optional:      true,
				ConflictsWith: []string{"tiers"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"divide_by": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"round": {
							Type:     schema.TypeString,
							Required: true,
							//ExactlyOneOf: []string{
							//	string(stripe.PriceTransformQuantityRoundUp),
							//	string(stripe.PriceTransformQuantityRoundDown),
							//},
						},
					},
				},
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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

	diags := CallSet(
		d.Set("currency", string(price.Currency)),
		d.Set("product_id", price.Product.ID),
		d.Set("unit_amount", int(price.UnitAmount)),
		d.Set("unit_amount_decimal", price.UnitAmountDecimal),
		d.Set("active", price.Active),
		d.Set("nickname", price.Nickname),
		d.Set("billing_scheme", string(price.BillingScheme)),
		d.Set("tiers_mode", string(price.TiersMode)),
		d.Set("metadata", price.Metadata),
	)
	if len(diags) > 0 {
		return diags
	}

	if price.Recurring != nil {
		err := d.Set("recurring", flattenPriceRecurring(price.Recurring))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if len(price.Tiers) > 0 {
		err := d.Set("tiers", flattenPriceTiers(price.Tiers))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if price.TransformQuantity != nil {
		err := d.Set("transform_quantity", flattenPriceTransformQuantity(price.TransformQuantity))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return nil
}

func flattenPriceRecurring(r *stripe.PriceRecurring) map[string]interface{} {
	m := make(map[string]interface{})
	m["interval"] = string(r.Interval)
	m["aggregate_usage"] = string(r.AggregateUsage)
	m["interval_count"] = strconv.FormatInt(r.IntervalCount, 10)
	m["usage_type"] = string(r.UsageType)
	return m
}

func flattenPriceTiers(ts []*stripe.PriceTier) []map[string]interface{} {
	res := make([]map[string]interface{}, len(ts), len(ts))
	for i, t := range ts {
		m := make(map[string]interface{})
		m["up_to"] = int(t.UpTo)
		m["flat_amount"] = int(t.FlatAmount)
		m["flat_amount_decimal"] = t.FlatAmountDecimal
		m["unit_amount"] = int(t.UnitAmount)
		m["unit_amount_decimal"] = t.UnitAmountDecimal
		res[i] = m
	}
	return res
}

func flattenPriceTransformQuantity(tq *stripe.PriceTransformQuantity) map[string]interface{} {
	m := make(map[string]interface{})
	m["divide_by"] = int(tq.DivideBy)
	m["round"] = string(tq.Round)
	return m
}

func resourceStripePriceCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)

	params := &stripe.PriceParams{
		Currency: stripe.String(String(d, "currency")),
		Product:  stripe.String(String(d, "product_id")),
		Active:   stripe.Bool(Bool(d, "active")),
		Nickname: stripe.String(String(d, "nickname")),
	}
	unitAmountSet := false
	if unitAmountDecimal, set := d.GetOk("unit_amount_decimal"); set {
		params.UnitAmountDecimal = stripe.Float64(ToFloat64(unitAmountDecimal))
		unitAmountSet = true
	}
	if !unitAmountSet {
		params.UnitAmount = stripe.Int64(Int64(d, "unit_amount"))
	}
	if tiersMode, set := d.GetOk("tiers_mode"); set {
		params.TiersMode = stripe.String(ToString(tiersMode))
	}

	if recurring, set := d.GetOk("recurring"); set {
		m := ToMap(recurring)
		params.Recurring = &stripe.PriceRecurringParams{}
		if interval, set := m["interval"]; set {
			params.Recurring.Interval = stripe.String(ToString(interval))
		}
		if aggregateUsage, set := m["aggregate_usage"]; set {
			params.Recurring.AggregateUsage = stripe.String(ToString(aggregateUsage))
		}
		if intervalCount, set := m["interval_count"]; set {
			iCnt, _ := strconv.ParseInt(intervalCount.(string), 10, 64)
			params.Recurring.IntervalCount = stripe.Int64(iCnt)
		}
		if usageType, set := m["usage_type"]; set {
			params.Recurring.UsageType = stripe.String(ToString(usageType))
		}
	}
	if tiers, set := d.GetOk("tiers"); set {
		ms := ToMapSlice(tiers)
		for _, m := range ms {
			t := &stripe.PriceTierParams{
				UpTo:              stripe.Int64(ToInt64(m["up_to"])),
				FlatAmount:        stripe.Int64(ToInt64(m["flat_amount"])),
				FlatAmountDecimal: stripe.Float64(ToFloat64(m["flat_amount_decimal"])),
				UnitAmount:        stripe.Int64(ToInt64(m["unit_amount"])),
				UnitAmountDecimal: stripe.Float64(ToFloat64(m["unit_amount_decimal"])),
			}
			params.Tiers = append(params.Tiers, t)
		}
	}
	if transformQuantity, set := d.GetOk("transform_quantity"); set {
		m := ToMap(transformQuantity)
		params.TransformQuantity = &stripe.PriceTransformQuantityParams{
			DivideBy: stripe.Int64(ToInt64(m["divide_by"])),
			Round:    stripe.String(ToString(m["round"])),
		}
	}

	price, err := c.Prices.New(params)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(price.ID)
	return resourceStripeProductRead(ctx, d, m)
}

func resourceStripePriceDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	_, err := c.Prices.Update(d.Id(), &stripe.PriceParams{
		Active: stripe.Bool(false),
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return nil
}
