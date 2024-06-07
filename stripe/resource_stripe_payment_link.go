package stripe

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/client"
	"log"
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
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				Description: "Whether the payment link’s url is active. If false, " +
					"customers visiting the URL will be shown a page saying that the link has been deactivated.",
			},
			"line_items": {
				Type:     schema.TypeList,
				Required: true,
				Description: "The line items representing what is being sold. " +
					"Each line item represents an item being sold. Up to 20 line items are supported.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"price": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The ID of the Price or Plan object",
						},
						"quantity": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The quantity of the line item being purchased.",
						},
						"adjustable_quantity": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:     schema.TypeBool,
										Required: true,
										Description: "Set to true if the quantity can be adjusted to " +
											"any non-negative Integer.",
									},
									"maximum": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  99,
										Description: "The maximum quantity the customer can purchase. " +
											"By default this value is 99. You can specify a value up to 999.",
									},
									"minimum": {
										Type:     schema.TypeInt,
										Optional: true,
										Default:  0,
										Description: "The minimum quantity the customer can purchase. " +
											"By default this value is 0. " +
											"If there is only one item in the cart then that item’s quantity " +
											"cannot go down to 0.",
									},
								},
							},
						},
					},
				},
			},
			"after_completion": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							Description: "The specified behavior after the purchase is complete. " +
								"Either redirect or hosted_confirmation.",
						},
						"custom_message": { // when hosted_confirmation type is used
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A custom message to display to the customer after the purchase is complete",
						},
						"url": { // when redirect type is used
							Type:     schema.TypeString,
							Optional: true,
							Description: "The URL the customer will be redirected to after the purchase is complete. " +
								"You can embed {CHECKOUT_SESSION_ID} into the URL to have the id of the completed " +
								"checkout session included.",
						},
					},
				},
			},
			"allow_promotion_codes": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enables user redeemable promotion codes.",
			},
			"application_fee_amount": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: "The amount of the application fee (if any) that will be requested to be applied " +
					"to the payment and transferred to the application owner’s Stripe account. " +
					"Can only be applied when there are no line items with recurring prices.",
			},
			"application_fee_percent": {
				Type:     schema.TypeFloat,
				Optional: true,
				Description: "A non-negative decimal between 0 and 100, with at most two decimal places. " +
					"This represents the percentage of the subscription invoice total that will be transferred to " +
					"the application owner’s Stripe account. There must be at least 1 line item with a recurring " +
					"price to use this field.",
			},
			"automatic_tax": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:        schema.TypeBool,
							Required:    true,
							Description: "If true, tax will be calculated automatically using the customer’s location.",
						},
						"liability": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "Type of the account referenced in the request.",
									},
									"account": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The connected account being referenced when type is account.",
									},
								},
							},
						},
					},
				},
			},
			"billing_address_collection": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "auto",
				Description: "Configuration for collecting the customer’s billing address. Defaults to auto.",
			},
			"consent_collection": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"payment_method_reuse_agreement": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"position": {
										Type:     schema.TypeString,
										Required: true,
										Description: "Determines the position and visibility of the payment method " +
											"reuse agreement in the UI. When set to auto, " +
											"Stripe’s defaults will be used. When set to hidden, " +
											"the payment method reuse agreement text will always be hidden in the UI.",
									},
								},
							},
							Description: "Determines the display of payment method reuse agreement text in the UI. " +
								"If set to hidden, it will hide legal text related to the reuse of a payment method.",
						},
						"promotions": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "If set to auto, enables the collection of customer consent for promotional " +
								"communications. The Checkout Session will determine whether to display an option " +
								"to opt into promotional communication from the merchant depending on the customer’s " +
								"locale. Only available to US merchants.",
						},
						"terms_of_service": {
							Type:     schema.TypeString,
							Optional: true,
							Description: "If set to required, it requires customers to check a terms of service " +
								"checkbox before being able to pay. There must be a valid terms of service URL " +
								"set in your Dashboard settings.",
						},
					},
				},
			},
			"currency": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Three-letter ISO currency code, in lowercase. " +
					"Must be a supported currency and supported by each line item’s price.",
			},
			"custom_fields": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
							Description: "String of your choice that your integration can use to reconcile this field. " +
								"Must be unique to this field, alphanumeric, and up to 200 characters.",
						},
						"label": {
							Type:     schema.TypeList,
							MaxItems: 1,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"custom": {
										Type:     schema.TypeString,
										Required: true,
										Description: "Custom text for the label, displayed to the customer. " +
											"Up to 50 characters.",
									},
									"type": {
										Type:        schema.TypeString,
										Optional:    true,
										Default:     "custom",
										Description: "The type of the label.",
									},
								},
							},
							Description: "The label for the field, displayed to the customer.",
						},
						"type": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The type of the field.",
						},
						"dropdown": {
							Type:        schema.TypeList,
							Description: "Configuration for type=dropdown fields.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"label": {
										Type:     schema.TypeString,
										Required: true,
										Description: "The label for the option, displayed to the customer. " +
											"Up to 100 characters.",
									},
									"value": {
										Type:     schema.TypeString,
										Required: true,
										Description: "The value for this option, not displayed to the customer, " +
											"used by your integration to reconcile the option selected by the " +
											"customer. Must be unique to this option, alphanumeric, " +
											"and up to 100 characters.",
									},
								},
							},
						},
						"numeric": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Description: "Configuration for type=numeric fields.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"maximum_length": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The maximum character length constraint for the customer’s input.",
									},
									"minimum_length": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The minimum character length requirement for the customer’s input.",
									},
								},
							},
						},
						"optional": {
							Type:     schema.TypeBool,
							Optional: true,
							Description: "Whether the customer is required to complete the field before completing " +
								"the Checkout Session. Defaults to false.",
						},
						"text": {
							Type:        schema.TypeList,
							MaxItems:    1,
							Optional:    true,
							Description: "Configuration for type=text fields",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"maximum_length": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The maximum character length constraint for the customer’s input.",
									},
									"minimum_length": {
										Type:        schema.TypeInt,
										Optional:    true,
										Description: "The minimum character length requirement for the customer’s input.",
									},
								},
							},
						},
					},
				},
			},
			"custom_text": {
				Type: schema.TypeList, // TODO continue here
			},
			"customer_creation": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Configures whether checkout sessions created by this payment link create a Customer.",
			},
			"inactive_message": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "The custom message to be displayed to a customer " +
					"when a payment link is no longer active.",
			},
			"invoice_creation": {
				// FIXME
			},
			"on_behalf_of": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The account on behalf of which to charge.",
			},
			"payment_intent_data": {
				// FIXME
			},
			"payment_method_collection": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Specify whether Checkout should collect a payment method." +
					" When set to if_required, Checkout will not collect a payment method " +
					"when the total due for the session is 0." +
					"This may occur if the Checkout Session includes a free trial or a discount." +
					"Can only be set in subscription mode. Defaults to always.",
			},
			"payment_method_types": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "The list of payment method types that customers can use. If no value is passed, " +
					"Stripe will dynamically show relevant payment methods from your payment method " +
					"settings (20+ payment methods supported).",
			},
			"phone_number_collection": {
				// FIXME
			},
			"restrictions": {
				// FIXME
			},
			"shipping_address_collection": {
				// FIXME
			},
			"shipping_options": {
				// FIXME
			},
			"submit_type": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Describes the type of transaction being performed in order to customize relevant " +
					"text on the page, such as the submit button. Changing this value will also affect the " +
					"hostname in the url property (example: donate.stripe.com).",
			},
			"subscription_data": {
				// FIXME
			},
			"tax_id_collection": {
				// FIXME
			},
			"transfer_data": {
				// FIXME
			},
		},
	}
}

func resourceStripePaymentLinkRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	var paymentLink *stripe.PaymentLink
	var err error

	params := &stripe.PaymentLinkParams{}

	err = retryWithBackOff(func() error {
		paymentLink, err = c.PaymentLinks.Get(d.Id(), params)
		return err
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return CallSet(
		d.Set("id", paymentLink.ID),
		// FIXME continue with other fields
	)
}

func resourceStripePaymentLinkCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	var paymentLink *stripe.PaymentLink
	var err error

	params := &stripe.PaymentLinkParams{}
	params.AddExpand("line_items")

	// FIXME payment link creation

	err = retryWithBackOff(func() error {
		paymentLink, err = c.PaymentLinks.New(params)
		return err
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(paymentLink.ID)
	return resourceStripePaymentLinkRead(ctx, d, m)
}

func resourceStripePaymentLinkUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	var err error

	params := &stripe.PaymentLinkParams{}

	// FIXME update fields

	err = retryWithBackOff(func() error {
		_, err = c.PaymentLinks.Update(d.Id(), params)
		return err
	})

	if err != nil {
		return diag.FromErr(err)
	}

	return resourceStripePaymentLinkRead(ctx, d, m)
}

func resourceStripePaymentLinkDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println(
		"[WARN] Stripe doesn't support deletion of payment links. " +
			"Payment link will be deactivated but not deleted and removed from the TF state")

	c := m.(*client.API)
	var err error

	params := stripe.PaymentLinkParams{
		Active: stripe.Bool(false),
	}

	err = retryWithBackOff(func() error {
		_, err = c.PaymentLinks.Update(d.Id(), &params)
		return err
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
