package stripe

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/client"
)

func resourceStripePaymentLink() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceStripePaymentLinkRead,
		CreateContext: resourceStripePaymentLinkCreate,
		UpdateContext: resourceStripePaymentLinkUpdate,
		DeleteContext: resourceStripePaymentLinkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifier for the object.",
			},
			"active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the payment link's url is active.",
			},
			"allow_promotion_codes": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enables user redeemable promotion codes. Defaults to false.",
			},
			"automatic_tax": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Configuration for automatic tax collection. Defaults to false.",
			},
			"currency": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Three-letter ISO currency code, in lowercase.",
			},
			"custom_text": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "Display additional text for your customers using custom text.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"shipping_address_message": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Custom text that should be displayed alongside shipping address collection. Text may be up to 1000 characters in length.",
						},
						"submit_message": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Custom text that should be displayed alongside the payment confirmation button. Text may be up to 1000 characters in length.",
						},
					},
				},
			},
			"line_item": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The line items that will be displayed on the payment link.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"price": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the [Price](https://stripe.com/docs/api/prices) or [Plan](https://stripe.com/docs/api/plans) object.",
						},
						"quantity": {
							Type:        schema.TypeInt,
							Optional:    true,
							Description: "The quantity of the line item being purchased.",
						},
					},
				},
			},
			"payment_method_types": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of payment method types that customers can use. Pass an empty string to enable automatic payment methods that use your [payment method settings](https://dashboard.stripe.com/settings/payment_methods).",
			},
			"subscription_data": {
				Type:        schema.TypeList,
				Optional:    true,
				ForceNew:    true,
				MaxItems:    1,
				Description: "When creating a subscription, the specified configuration data will be used. There must be at least one line item with a recurring price to use `subscription_data`.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The subscription's description, meant to be displayable to the customer. Use this field to optionally store an explanation of the subscription.",
						},
						"trial_period_days": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
							Description: "Integer representing the number of trial period days before the customer is charged for the first time. Has to be at least 1.",
						},
					},
				},
			},
			"url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The public URL that can be shared with customers.",
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

func resourceStripePaymentLinkRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	paymentLink, err := c.PaymentLinks.Get(d.Id(), nil)
	if err != nil {
		return diag.FromErr(err)
	}

	return CallSet(
		d.Set("active", paymentLink.Active),
		d.Set("allow_promotion_codes", paymentLink.AllowPromotionCodes),
		d.Set("automatic_tax", paymentLink.AutomaticTax.Enabled),
		func() error {
			if paymentLink.AutomaticTax != nil {
				return d.Set("automatic_tax", paymentLink.AutomaticTax.Enabled)
			}
			return nil
		}(),
		d.Set("currency", paymentLink.Currency),
		func() error {
			if paymentLink.CustomText != nil {
				customText := []map[string]interface{}{{}}
				if paymentLink.CustomText.ShippingAddress != nil {
					customText[0]["shipping_address_message"] = paymentLink.CustomText.ShippingAddress.Message
				}
				if paymentLink.CustomText.Submit != nil {
					customText[0]["submit_message"] = paymentLink.CustomText.Submit.Message
				}
				if len(customText[0]) > 0 {
					return d.Set("custom_text", customText)
				}
			}
			return nil
		}(),
		func() error {
			if paymentLink.LineItems != nil {
				if len(paymentLink.LineItems.Data) > 0 {
					var lineItems []map[string]interface{}
					for _, lineItem := range paymentLink.LineItems.Data {
						l := map[string]interface{}{
							"price":    lineItem.Price.ID,
							"quantity": lineItem.Quantity,
						}
						lineItems = append(lineItems, l)
					}
					return d.Set("line_item", lineItems)
				}
			}
			return nil
		}(),
		func() error {
			if paymentLink.PaymentMethodTypes != nil {
				var paymentMethodTypes []string
				for _, methodType := range paymentLink.PaymentMethodTypes {
					paymentMethodTypes = append(paymentMethodTypes, string(methodType))
				}
				return d.Set("payment_method_types", paymentMethodTypes)
			}
			return nil
		}(),
		func() error {
			if paymentLink.SubscriptionData != nil {
				subscriptionData := []map[string]interface{}{{}}
				if paymentLink.SubscriptionData.Description != "" {
					subscriptionData[0]["description"] = paymentLink.SubscriptionData.Description
				}
				if paymentLink.SubscriptionData.TrialPeriodDays != 0 {
					subscriptionData[0]["trial_period_days"] = paymentLink.SubscriptionData.TrialPeriodDays
				}
				if len(subscriptionData[0]) > 0 {
					return d.Set("subscription_data", subscriptionData)
				}
			}
			return nil
		}(),
		d.Set("url", paymentLink.URL),
		d.Set("metadata", paymentLink.Metadata),
	)
}

func resourceStripePaymentLinkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	params := &stripe.PaymentLinkParams{}

	if allowPromotionCodes, set := d.GetOk("allow_promotion_codes"); set {
		params.AllowPromotionCodes = stripe.Bool(ToBool(allowPromotionCodes))
	}

	if automaticTax, set := d.GetOk("automatic_tax"); set {
		params.AutomaticTax = &stripe.PaymentLinkAutomaticTaxParams{
			Enabled: stripe.Bool(ToBool(automaticTax)),
		}
	}

	if currency, set := d.GetOk("currency"); set {
		params.Currency = stripe.String(ToString(currency))
	}

	if customText, set := d.GetOk("custom_text"); set {
		params.CustomText = &stripe.PaymentLinkCustomTextParams{}
		for k, v := range ToMap(customText) {
			switch {
			case k == "shipping_address" && ToString(v) != "":
				params.CustomText.ShippingAddress = &stripe.PaymentLinkCustomTextShippingAddressParams{
					Message: stripe.String(ToString(v)),
				}
			case k == "submit_message" && ToString(v) != "":
				params.CustomText.Submit = &stripe.PaymentLinkCustomTextSubmitParams{
					Message: stripe.String(ToString(v)),
				}
			}
		}
	}

	if lineItems, set := d.GetOk("line_item"); set {
		for _, li := range ToSlice(lineItems) {
			lineItem := &stripe.PaymentLinkLineItemParams{}
			for k, v := range ToMap(li) {
				switch {
				case k == "price" && ToString(v) != "":
					lineItem.Price = stripe.String(ToString(v))
				case k == "quantity" && ToInt64(v) != 0:
					lineItem.Quantity = stripe.Int64(ToInt64(v))
				}
			}
			params.LineItems = append(params.LineItems, lineItem)
		}
	}

	if paymentMethodTypes, set := d.GetOk("payment_method_types"); set {
		for _, pm := range ToSlice(paymentMethodTypes) {
			params.PaymentMethodTypes = append(params.PaymentMethodTypes, stripe.String(ToString(pm)))
		}
	}

	if subscriptionData, set := d.GetOk("subscription_data"); set {
		params.SubscriptionData = &stripe.PaymentLinkSubscriptionDataParams{}
		for k, v := range ToMap(subscriptionData) {
			switch {
			case k == "description" && ToString(v) != "":
				params.SubscriptionData.Description = stripe.String(ToString(v))
			case k == "trial_period_days" && ToInt64(v) != 0:
				params.SubscriptionData.TrialPeriodDays = stripe.Int64(ToInt64(v))
			}
		}
	}

	if meta, set := d.GetOk("metadata"); set {
		for k, v := range ToMap(meta) {
			params.AddMetadata(k, ToString(v))
		}
	}

	paymentLink, err := c.PaymentLinks.New(params)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(paymentLink.ID)
	return resourceStripePaymentLinkRead(ctx, d, m)
}

func resourceStripePaymentLinkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	params := &stripe.PaymentLinkParams{}

	if d.HasChange("active") {
		params.Active = stripe.Bool(ExtractBool(d, "active"))
	}

	if d.HasChange("allow_promotion_codes") {
		params.AllowPromotionCodes = stripe.Bool(ExtractBool(d, "allow_promotion_codes"))
	}

	if d.HasChange("automatic_tax") {
		params.AutomaticTax = &stripe.PaymentLinkAutomaticTaxParams{
			Enabled: stripe.Bool(ExtractBool(d, "automatic_tax")),
		}
	}

	if d.HasChange("custom_text") {
		params.CustomText = &stripe.PaymentLinkCustomTextParams{}
		for k, v := range ExtractMap(d, "custom_text") {
			switch {
			case k == "shipping_address" && ToString(v) != "":
				params.CustomText.ShippingAddress = &stripe.PaymentLinkCustomTextShippingAddressParams{
					Message: stripe.String(ToString(v)),
				}
			case k == "submit_message" && ToString(v) != "":
				params.CustomText.Submit = &stripe.PaymentLinkCustomTextSubmitParams{
					Message: stripe.String(ToString(v)),
				}
			}
		}
	}

	if d.HasChange("line_item") {
		params.LineItems = nil
		for _, li := range ExtractMapSlice(d, "line_item") {
			lineItem := &stripe.PaymentLinkLineItemParams{}
			for k, v := range li {
				switch {
				case k == "price" && ToString(v) != "":
					lineItem.Price = stripe.String(ToString(v))
				case k == "quantity" && ToInt64(v) != 0:
					lineItem.Quantity = stripe.Int64(ToInt64(v))
				}
			}
			params.LineItems = append(params.LineItems, lineItem)
		}
	}

	if d.HasChange("payment_method_types") {
		params.PaymentMethodTypes = nil
		for _, pm := range ExtractStringSlice(d, "payment_method_types") {
			params.PaymentMethodTypes = append(params.PaymentMethodTypes, stripe.String(ToString(pm)))
		}
	}

	if d.HasChange("metadata") {
		params.Metadata = nil
		UpdateMetadata(d, params)
	}

	_, err := c.PaymentLinks.Update(d.Id(), params)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceStripePaymentLinkRead(ctx, d, m)
}

func resourceStripePaymentLinkDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	params := stripe.PaymentLinkParams{
		Active: stripe.Bool(false),
	}

	if _, err := c.PaymentLinks.Update(d.Id(), &params); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
