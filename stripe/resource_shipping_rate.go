package stripe

import (
	"context"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/client"
)

func resourceStripeShippingRate() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceStripeShippingRateRead,
		CreateContext: resourceStripeShippingRateCreate,
		UpdateContext: resourceStripeShippingRateUpdate,
		DeleteContext: resourceStripeShippingRateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifier for the object.",
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "fixed_amount",
				Description: "The type of calculation to use on the shipping rate. " +
					"Can only be fixed_amount for now",
			},
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: "The name of the shipping rate, meant to be displayable to the customer. " +
					"This will appear on CheckoutSessions.",
			},
			"active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the shipping rate can be used for new purchases. Defaults to true.",
			},
			"fixed_amount": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "Describes a fixed amount to charge for shipping. Must be present for now.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"amount": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "A non-negative integer in cents representing how much to charge.",
						},
						"currency": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Three-letter ISO currency code, in lowercase. Must be a supported currency.",
						},
						"currency_option": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Description: "Shipping rates defined in each available currency option. " +
								"Each key must be a three-letter ISO currency code and a supported currency. " +
								"For example, to get your shipping rate in eur, " +
								"fetch the value of the eur key in currency_options. " +
								"This field is not included by default. " +
								"To include it in the response, expand the currency_options field.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"currency": {
										Type:        schema.TypeString,
										Required:    true,
										ForceNew:    true,
										Description: "Three-letter ISO currency code, in lowercase. Must be a supported currency.",
									},
									"amount": {
										Type:        schema.TypeInt,
										Required:    true,
										ForceNew:    true,
										Description: "A non-negative integer in cents representing how much to charge.",
									},
									"tax_behavior": {
										Type:     schema.TypeString,
										Optional: true,
										ForceNew: true,
										Default:  stripe.PriceTaxBehaviorUnspecified,
										Description: "Specifies whether the rate is considered inclusive of taxes or " +
											"exclusive of taxes. One of inclusive, exclusive, or unspecified. ",
									},
								},
							},
						},
					},
				},
			},
			"delivery_estimate": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"minimum": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"unit": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
										Description: "The lower bound of the estimated range. " +
											"If empty, represents no lower bound.",
									},
									"value": {
										Type:        schema.TypeInt,
										Required:    true,
										ForceNew:    true,
										Description: "Must be greater than 0.",
									},
								},
							},
						},
						"maximum": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"unit": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
										Description: "The upper bound of the estimated range. " +
											"If empty, represents no lower bound.",
									},
									"value": {
										Type:        schema.TypeInt,
										ForceNew:    true,
										Required:    true,
										Description: "Must be greater than 0.",
									},
								},
							},
						},
					},
				},
			},
			"tax_behavior": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  stripe.PriceTaxBehaviorUnspecified,
				Description: "Specifies whether the rate is considered inclusive of taxes or " +
					"exclusive of taxes. One of inclusive, exclusive, or unspecified. ",
			},
			"livemode": {
				Type:     schema.TypeBool,
				Computed: true,
				Description: "Has the value true if the object exists in live mode or the value false " +
					"if the object exists in test mode.",
			},
			"tax_code": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "A tax code ID. The Shipping tax code is txcd_92010001.",
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

func resourceStripeShippingRateRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	var shippingRate *stripe.ShippingRate
	var err error

	params := &stripe.ShippingRateParams{}
	params.AddExpand("fixed_amount.currency_options")
	err = retryWithBackOff(func() error {
		shippingRate, err = c.ShippingRates.Get(d.Id(), params)
		return err
	})
	switch {
	case isNotFoundErr(err):
		d.SetId("") // remove when resource does not exist
		return nil
	case err != nil:
		return diag.FromErr(err)
	}

	return CallSet(
		d.Set("type", shippingRate.Type),
		d.Set("display_name", shippingRate.DisplayName),
		d.Set("active", shippingRate.Active),
		func() error {
			if shippingRate.FixedAmount != nil {
				fixedAmount := map[string]interface{}{
					"amount":   shippingRate.FixedAmount.Amount,
					"currency": shippingRate.FixedAmount.Currency,
				}

				if len(shippingRate.FixedAmount.CurrencyOptions) > 0 {
					var options []map[string]interface{}
					for currency, currencyOptions := range shippingRate.FixedAmount.CurrencyOptions {
						if currency == string(shippingRate.FixedAmount.Currency) {
							continue // don't add the same currency into the currency options
						}
						options = append(options, map[string]interface{}{
							"currency":     currency,
							"amount":       currencyOptions.Amount,
							"tax_behavior": currencyOptions.TaxBehavior,
						})
					}
					sort.Slice(options, func(i, j int) bool {
						return ToString(options[i]["currency"]) < ToString(options[j]["currency"])
					})
					fixedAmount["currency_option"] = options
				}

				return d.Set("fixed_amount", []map[string]interface{}{fixedAmount})
			}
			return nil
		}(),
		func() error {
			if shippingRate.DeliveryEstimate != nil {
				deliveryEstimate := make(map[string]interface{})
				if shippingRate.DeliveryEstimate.Minimum != nil {
					deliveryEstimate["minimum"] = []map[string]interface{}{
						{
							"unit":  shippingRate.DeliveryEstimate.Minimum.Unit,
							"value": shippingRate.DeliveryEstimate.Minimum.Value,
						},
					}
				}

				if shippingRate.DeliveryEstimate.Maximum != nil {
					deliveryEstimate["maximum"] = []map[string]interface{}{
						{
							"unit":  shippingRate.DeliveryEstimate.Maximum.Unit,
							"value": shippingRate.DeliveryEstimate.Maximum.Value,
						},
					}
				}
				return d.Set("delivery_estimate", []map[string]interface{}{deliveryEstimate})
			}
			return nil
		}(),
		d.Set("tax_behavior", shippingRate.TaxBehavior),
		d.Set("livemode", shippingRate.Livemode),
		d.Set("tax_code", shippingRate.TaxCode),
		d.Set("metadata", shippingRate.Metadata),
	)
}

func resourceStripeShippingRateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	var shippingRate *stripe.ShippingRate
	var err error

	params := &stripe.ShippingRateParams{
		Type:        stripe.String(ExtractString(d, "type")),
		DisplayName: stripe.String(ExtractString(d, "display_name")),
	}

	if fixedAmount, set := d.GetOk("fixed_amount"); set {
		fixedAmountMap := ToMap(fixedAmount)
		params.FixedAmount = &stripe.ShippingRateFixedAmountParams{
			Amount:   stripe.Int64(ToInt64(fixedAmountMap["amount"])),
			Currency: stripe.String(ToString(fixedAmountMap["currency"])),
		}
		if _, set := fixedAmountMap["currency_option"]; set {
			params.FixedAmount.CurrencyOptions = make(map[string]*stripe.ShippingRateFixedAmountCurrencyOptionsParams)
			for _, options := range ToMapSlice(fixedAmountMap["currency_option"]) {
				params.FixedAmount.CurrencyOptions[ToString(options["currency"])] = &stripe.ShippingRateFixedAmountCurrencyOptionsParams{
					Amount:      stripe.Int64(ToInt64(options["amount"])),
					TaxBehavior: stripe.String(ToString(options["tax_behavior"])),
				}
			}
		}
	}

	if deliveryEstimate, set := d.GetOk("delivery_estimate"); set {
		params.DeliveryEstimate = &stripe.ShippingRateDeliveryEstimateParams{}
		delivery := ToMap(deliveryEstimate)
		if minimumDelivery, set := delivery["minimum"]; set {
			minimum := ToMap(minimumDelivery)
			params.DeliveryEstimate.Minimum = &stripe.ShippingRateDeliveryEstimateMinimumParams{
				Value: stripe.Int64(ToInt64(minimum["value"])),
			}
			if _, set := minimum["unit"]; set {
				params.DeliveryEstimate.Minimum.Unit = stripe.String(ToString(minimum["unit"]))
			}

		}
		if maximumDelivery, set := delivery["maximum"]; set {
			maximum := ToMap(maximumDelivery)
			params.DeliveryEstimate.Maximum = &stripe.ShippingRateDeliveryEstimateMaximumParams{
				Value: stripe.Int64(ToInt64(maximum["value"])),
			}
			if _, set := maximum["unit"]; set {
				params.DeliveryEstimate.Maximum.Unit = stripe.String(ToString(maximum["unit"]))
			}
		}
	}

	if taxBehavior, set := d.GetOk("tax_behavior"); set {
		params.TaxBehavior = stripe.String(ToString(taxBehavior))
	}
	if taxCode, set := d.GetOk("tax_code"); set {
		params.TaxCode = stripe.String(ToString(taxCode))
	}
	if meta, set := d.GetOk("metadata"); set {
		for k, v := range ToMap(meta) {
			params.AddMetadata(k, ToString(v))
		}
	}

	err = retryWithBackOff(func() error {
		shippingRate, err = c.ShippingRates.New(params)
		return err
	})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(shippingRate.ID)
	return resourceStripeShippingRateRead(ctx, d, m)
}

func resourceStripeShippingRateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	var err error

	params := &stripe.ShippingRateParams{}
	if d.HasChange("active") {
		params.Active = stripe.Bool(ExtractBool(d, "active"))
	}
	if d.HasChange("tax_behavior") {
		params.TaxBehavior = stripe.String(ExtractString(d, "tax_behavior"))
	}
	if d.HasChange("metadata") {
		params.Metadata = nil
		UpdateMetadata(d, params)
	}

	err = retryWithBackOff(func() error {
		_, err = c.ShippingRates.Update(d.Id(), params)
		return err
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceStripeShippingRateRead(ctx, d, m)
}

func resourceStripeShippingRateDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	var err error

	params := &stripe.ShippingRateParams{
		Active: stripe.Bool(false),
	}

	err = retryWithBackOff(func() error {
		_, err = c.ShippingRates.Update(d.Id(), params)
		return err
	})

	d.SetId("")
	return nil
}
