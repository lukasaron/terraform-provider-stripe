package stripe

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/client"
)

func resourceStripePortalConfiguration() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceStripePortalConfigurationRead,
		CreateContext: resourceStripePortalConfigurationCreate,
		UpdateContext: resourceStripePortalConfigurationUpdate,
		DeleteContext: resourceStripePortalConfigurationDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifier for the object.",
			},
			"object": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "String representing the object's type.",
			},
			"active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the configuration is active and can be used to create portal sessions.",
			},
			"application": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the Connect Application that created the configuration.",
			},
			"business_profile": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The business information shown to customers in the portal.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"headline": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The messaging shown to customers in the portal.",
						},
						"privacy_policy_url": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "A link to the business's publicly available privacy policy.",
						},
						"terms_of_service_url": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "A link to the business's publicly available terms of service.",
						},
					},
				},
			},
			"created": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: "Time at which the object was created. " +
					"Measured in seconds since the Unix epoch.",
			},
			"default_return_url": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "The default URL to redirect customers to when they click on the portal's " +
					"link to return to your website. This can be overriden when creating the session.",
			},
			"features": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "Information about the features available in the portal.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"customer_update": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Information about updating the customer details in the portal.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Whether the feature is enabled.",
									},
									"allowed_updates": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "The types of customer updates that are supported. When empty, customers are not updateable.",
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringInSlice([]string{"email", "address", "shipping", "phone", "tax_id"}, false),
										},
									},
								},
							},
						},
						"invoice_history": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Information about showing the billing history in the portal.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Whether the feature is enabled.",
									},
								},
							},
						},
						"payment_method_update": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Information about updating payment methods in the portal.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Whether the feature is enabled.",
									},
								},
							},
						},
						"subscription_cancel": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Information about canceling subscriptions in the portal.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Whether the feature is enabled.",
									},
									"cancellation_reason": {
										Type:        schema.TypeList,
										Optional:    true,
										Description: "Whether the cancellation reasons will be collected in the portal and which options are exposed to the customer",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"enabled": {
													Type:        schema.TypeBool,
													Required:    true,
													Description: "Whether the feature is enabled.",
												},
												"options": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "Which cancellation reasons will be given as options to the customer.",
													Elem: &schema.Schema{
														Type:         schema.TypeString,
														ValidateFunc: validation.StringInSlice([]string{"too_expensive", "missing_features", "switched_service", "unused", "customer_service", "too_complex", "low_quality", "other"}, false),
													},
												},
											},
										},
									},
									"mode": {
										Type:         schema.TypeString,
										Optional:     true,
										Description:  "Whether to cancel subscriptions immediately or at the end of the billing period.",
										ValidateFunc: validation.StringInSlice([]string{"immediately", "at_period_end"}, false),
									},
									"proration_behavior": {
										Type:         schema.TypeString,
										Optional:     true,
										Description:  "Whether to create prorations when canceling subscriptions.",
										ValidateFunc: validation.StringInSlice([]string{"none", "create_prorations"}, false),
									},
								},
							},
						},
						"subscription_pause": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Information about pausing subscriptions in the portal.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Whether the feature is enabled.",
									},
								},
							},
						},
						"subscription_update": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "Information about updating subscriptions in the portal.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"default_allowed_updates": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "The types of subscription updates that are supported. When empty, subscriptions are not updateable.",
										Elem: &schema.Schema{
											Type:         schema.TypeString,
											ValidateFunc: validation.StringInSlice([]string{"price", "quantity", "promotion_code"}, false),
										},
									},
									"enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Whether the feature is enabled.",
									},
									"products": {
										Type:        schema.TypeSet,
										Optional:    true,
										Description: "The list of products that support subscription updates.",
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"prices": {
													Type:        schema.TypeList,
													Required:    true,
													Description: "The list of price IDs for the product that a subscription can be updated to.",
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
												"product": {
													Type:        schema.TypeString,
													Required:    true,
													Description: "The product id.",
												},
											},
										},
									},
									"proration_behavior": {
										Type:         schema.TypeString,
										Optional:     true,
										Description:  "Determines how to handle prorations resulting from subscription updates",
										ValidateFunc: validation.StringInSlice([]string{"none", "create_prorations", "always_invoice"}, false),
									},
								},
							},
						},
					},
				},
			},
			"is_default": {
				Type:     schema.TypeBool,
				Computed: true,
				Description: "Whether the configuration is the default. If true, this configuration can be " +
					"managed in the Dashboard and portal sessions will use this configuration unless it is " +
					"overriden when creating the session.",
			},
			"livemode": {
				Type:     schema.TypeBool,
				Computed: true,
				Description: "Has the value true if the object exists in live mode or the value false if the " +
					"object exists in test mode.",
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Description: "Set of key-value pairs that you can attach to an object. " +
					"This can be useful for storing additional information about the object in a structured format.",
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"updated": {
				Type:     schema.TypeInt,
				Computed: true,
				Description: "Time at which the object was last updated. " +
					"Measured in seconds since the Unix epoch.",
			},
		},
	}
}

func resourceStripePortalConfigurationRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	portal, err := c.BillingPortalConfigurations.Get(d.Id(), nil)
	if err != nil {
		return diag.FromErr(err)
	}
	return CallSet(
		d.Set("id", portal.ID),
		d.Set("object", portal.Object),
		d.Set("active", portal.Active),
		d.Set("application", portal.Application),
		d.Set("business_profile", portal.BusinessProfile),
		d.Set("created", portal.Created),
		d.Set("default_return_url", portal.DefaultReturnURL),
		d.Set("features", portal.Features),
		d.Set("is_default", portal.IsDefault),
		d.Set("livemode", portal.Livemode),
		d.Set("metadata", portal.Metadata),
		d.Set("updated", portal.Updated),
	)
}

func expandBusinessProfile(businessProfileI []interface{}) *stripe.BillingPortalConfigurationBusinessProfileParams {
	businessProfile := &stripe.BillingPortalConfigurationBusinessProfileParams{}
	for _, v := range businessProfileI {
		businessProfileMap := ToMap(v)
		if privacyPolicyURL, set := businessProfileMap["privacy_policy_url"]; set {
			businessProfile.PrivacyPolicyURL = stripe.String(ToString(privacyPolicyURL))
		}
		if termsOfServiceURL, set := businessProfileMap["terms_of_service_url"]; set {
			businessProfile.TermsOfServiceURL = stripe.String(ToString(termsOfServiceURL))
		}
		if headline, set := businessProfileMap["headline"]; set {
			businessProfile.Headline = stripe.String(ToString(headline))
		}
	}
	return businessProfile
}

func expandFeatures(featuresI []interface{}) *stripe.BillingPortalConfigurationFeaturesParams {
	features := &stripe.BillingPortalConfigurationFeaturesParams{}
	for _, v := range featuresI {
		featuresMap := ToMap(v)

		if customerUpdateSettings, set := featuresMap["customer_update"]; set {
			customerUpdate := &stripe.BillingPortalConfigurationFeaturesCustomerUpdateParams{}
			cu := ToSlice(customerUpdateSettings)
			for _, props := range cu {
				p := ToMap(props)
				if allowedUpdates, set := p["allowed_updates"]; set {
					customerUpdate.AllowedUpdates = stripe.StringSlice(ToStringSlice(ToSlice(allowedUpdates)))
				}
				if enabled, set := p["enabled"]; set {
					customerUpdate.Enabled = stripe.Bool(ToBool(enabled))
				}
			}
			features.CustomerUpdate = customerUpdate
		}

		if invoiceHistorySettings, set := featuresMap["invoice_history"]; set {
			invoiceHistory := &stripe.BillingPortalConfigurationFeaturesInvoiceHistoryParams{}
			ih := ToSlice(invoiceHistorySettings)
			for _, props := range ih {
				p := ToMap(props)
				if enabled, set := p["enabled"]; set {
					invoiceHistory.Enabled = stripe.Bool(ToBool(enabled))
				}
			}
			features.InvoiceHistory = invoiceHistory
		}

		if paymentMethodUpdateSettings, set := featuresMap["payment_method_update"]; set {
			paymentMethodUpdate := &stripe.BillingPortalConfigurationFeaturesPaymentMethodUpdateParams{}
			pmu := ToSlice(paymentMethodUpdateSettings)
			for _, props := range pmu {
				p := ToMap(props)
				if enabled, set := p["enabled"]; set {
					paymentMethodUpdate.Enabled = stripe.Bool(ToBool(enabled))
				}
			}
			features.PaymentMethodUpdate = paymentMethodUpdate
		}

		if subscriptionCancelSettings, set := featuresMap["subscription_cancel"]; set {
			subscriptionCancel := &stripe.BillingPortalConfigurationFeaturesSubscriptionCancelParams{}
			sc := ToSlice(subscriptionCancelSettings)
			for _, props := range sc {
				p := ToMap(props)
				if cancellationReason, set := p["cancellation_reason"]; set {
					subscriptionCancelReason := &stripe.BillingPortalConfigurationFeaturesSubscriptionCancelCancellationReasonParams{}
					scr := ToSlice(cancellationReason)
					for _, scrProps := range scr {
						scrP := ToMap(scrProps)
						if options, set := scrP["options"]; set {
							subscriptionCancelReason.Options = stripe.StringSlice(ToStringSlice(ToSlice(options)))
						}
						if enabled, set := scrP["enabled"]; set {
							subscriptionCancelReason.Enabled = stripe.Bool(ToBool(enabled))
						}
					}
					subscriptionCancel.CancellationReason = subscriptionCancelReason
				}

				if enabled, set := p["enabled"]; set {
					subscriptionCancel.Enabled = stripe.Bool(ToBool(enabled))
				}

				if mode, set := p["mode"]; set {
					subscriptionCancel.Mode = stripe.String(ToString(mode))
				}

				if prorationBehavior, set := p["proration_behavior"]; set {
					subscriptionCancel.ProrationBehavior = stripe.String(ToString(prorationBehavior))
				}
			}
			features.SubscriptionCancel = subscriptionCancel
		}

		if subscriptionPauseSettings, set := featuresMap["subscription_pause"]; set {
			subscriptionPause := &stripe.BillingPortalConfigurationFeaturesSubscriptionPauseParams{}
			sp := ToSlice(subscriptionPauseSettings)
			for _, props := range sp {
				p := ToMap(props)
				if enabled, set := p["enabled"]; set {
					subscriptionPause.Enabled = stripe.Bool(ToBool(enabled))
				}
			}
			features.SubscriptionPause = subscriptionPause
		}

		if subscriptionUpdateSettings, set := featuresMap["subscription_update"]; set {
			subscriptionUpdate := &stripe.BillingPortalConfigurationFeaturesSubscriptionUpdateParams{}
			sp := ToSlice(subscriptionUpdateSettings)
			for _, props := range sp {
				p := ToMap(props)
				if defaultAllowedUpdates, set := p["default_allowed_updates"]; set {
					subscriptionUpdate.DefaultAllowedUpdates = stripe.StringSlice(ToStringSlice(ToSlice(defaultAllowedUpdates)))
				}

				if enabled, set := p["enabled"]; set {
					subscriptionUpdate.Enabled = stripe.Bool(ToBool(enabled))
				}

				if products, set := p["products"]; set {
					var productsParams []*stripe.BillingPortalConfigurationFeaturesSubscriptionUpdateProductParams
					schemaSet := products.(*schema.Set)
					productsList := schemaSet.List()
					for _, i := range productsList {
						pParams := &stripe.BillingPortalConfigurationFeaturesSubscriptionUpdateProductParams{}
						finalProduct := ToMap(i)
						if product, set := finalProduct["product"]; set {
							pParams.Product = stripe.String(ToString(product))
						}

						if prices, set := finalProduct["prices"]; set {
							pParams.Prices = stripe.StringSlice(ToStringSlice(ToSlice(prices)))
						}
						productsParams = append(productsParams, pParams)
					}
					subscriptionUpdate.Products = productsParams
				}

				if prorationBehavior, set := p["proration_behavior"]; set {
					subscriptionUpdate.ProrationBehavior = stripe.String(ToString(prorationBehavior))
				}
			}
			features.SubscriptionUpdate = subscriptionUpdate
		}
	}
	return features
}

func resourceStripePortalConfigurationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	params := &stripe.BillingPortalConfigurationParams{}
	if defaultReturnURL, set := d.GetOk("default_return_url"); set {
		params.DefaultReturnURL = stripe.String(ToString(defaultReturnURL))
	}
	if businessProfile, set := d.GetOk("business_profile"); set {
		params.BusinessProfile = expandBusinessProfile(ToSlice(businessProfile))
	}
	if features, set := d.GetOk("features"); set {
		params.Features = expandFeatures(ToSlice(features))
	}
	if meta, set := d.GetOk("metadata"); set {
		for k, v := range ToMap(meta) {
			params.AddMetadata(k, ToString(v))
		}
	}

	portalConfiguration, err := c.BillingPortalConfigurations.New(params)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(portalConfiguration.ID)
	return resourceStripePortalConfigurationRead(ctx, d, m)
}

func resourceStripePortalConfigurationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	params := &stripe.BillingPortalConfigurationParams{}
	if d.HasChange("active") {
		params.Active = stripe.Bool(ExtractBool(d, "active"))
	}

	if d.HasChange("default_return_url") {
		params.DefaultReturnURL = stripe.String(ExtractString(d, "default_return_url"))
	}

	if d.HasChange("metadata") {
		params.Metadata = nil
		metadata := ExtractMap(d, "metadata")
		for k, v := range metadata {
			params.AddMetadata(k, ToString(v))
		}
	}

	if d.HasChange("business_profile") {
		_, newBusinessProfile := d.GetChange("business_profile")
		params.BusinessProfile = expandBusinessProfile(ToSlice(newBusinessProfile))
	}

	if d.HasChange("features") {
		_, newFeatures := d.GetChange("features")
		params.Features = expandFeatures(ToSlice(newFeatures))
	}

	_, err := c.BillingPortalConfigurations.Update(d.Id(), params)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceStripePortalConfigurationRead(ctx, d, m)
}

func resourceStripePortalConfigurationDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[WARN] Stripe doesn't support deletion of customer portals. Portal will be deactivated but not deleted")

	c := m.(*client.API)
	params := stripe.BillingPortalConfigurationParams{
		Active: stripe.Bool(false),
	}

	if _, err := c.BillingPortalConfigurations.Update(d.Id(), &params); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
