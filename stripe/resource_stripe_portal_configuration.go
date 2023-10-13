package stripe

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/client"
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
			"active": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the configuration is active and can be used to create portal sessions.",
			},
			"business_profile": {
				Type:        schema.TypeList,
				MaxItems:    1,
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
							Optional:    true,
							Description: "A link to the business's publicly available privacy policy.",
						},
						"terms_of_service_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "A link to the business's publicly available terms of service.",
						},
					},
				},
			},
			"default_return_url": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "The default URL to redirect customers to when they click on the portal's " +
					"link to return to your website. This can be overriden when creating the session.",
			},
			"login_page": {
				Type:        schema.TypeList,
				MaxItems:    1,
				Optional:    true,
				Computed:    true,
				Description: "The hosted login page for this configuration.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Description: "Set to true to generate a shareable URL login_page.url that will take your" +
								" customers to a hosted login page for the customer portal.",
						},
						"url": {
							Type:     schema.TypeString,
							Computed: true,
							Description: "A shareable URL to the hosted portal login page. " +
								"Your customers will be able to log in with their email and receive a link to their customer portal.",
						},
					},
				},
			},
			"features": {
				Type:        schema.TypeList,
				Required:    true,
				MaxItems:    1,
				Description: "Information about the features available in the portal.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"customer_update": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
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
										Optional:    true,
										Description: "The types of customer updates that are supported. When empty, customers are not updatable.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"invoice_history": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							MaxItems:    1,
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
							Computed:    true,
							MaxItems:    1,
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
							Computed:    true,
							MaxItems:    1,
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
										MaxItems:    1,
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
													Elem:        &schema.Schema{Type: schema.TypeString},
												},
											},
										},
									},
									"mode": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Whether to cancel subscriptions immediately or at the end of the billing period.",
									},
									"proration_behavior": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Whether to create prorations when canceling subscriptions.",
									},
								},
							},
						},
						"subscription_pause": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Information about pausing subscriptions in the portal.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"enabled": {
										Type:        schema.TypeBool,
										Optional:    true,
										Description: "Whether the feature is enabled.",
									},
								},
							},
						},
						"subscription_update": {
							Type:        schema.TypeList,
							Optional:    true,
							Computed:    true,
							Description: "Information about updating subscriptions in the portal.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"default_allowed_updates": {
										Type:        schema.TypeList,
										Required:    true,
										Description: "The types of subscription updates that are supported. When empty, subscriptions are not updateable.",
										Elem:        &schema.Schema{Type: schema.TypeString},
									},
									"enabled": {
										Type:        schema.TypeBool,
										Required:    true,
										Description: "Whether the feature is enabled.",
									},
									"products": {
										Type:        schema.TypeList,
										Required:    true,
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
										Type:        schema.TypeString,
										Optional:    true,
										Description: "Determines how to handle prorations resulting from subscription updates",
									},
								},
							},
						},
					},
				},
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				Description: "Set of key-value pairs that you can attach to an object. " +
					"This can be useful for storing additional information about the object in a structured format.",
				Elem: &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceStripePortalConfigurationRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	var portal *stripe.BillingPortalConfiguration
	var err error

	params := &stripe.BillingPortalConfigurationParams{}
	params.AddExpand("features.subscription_update.products")

	err = retryWithBackOff(func() error {
		portal, err = c.BillingPortalConfigurations.Get(d.Id(), params)
		return err
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return CallSet(
		d.Set("id", portal.ID),
		d.Set("active", portal.Active),
		d.Set("business_profile", func() []map[string]interface{} {
			if portal.BusinessProfile != nil {
				return []map[string]interface{}{
					{
						"headline":             portal.BusinessProfile.Headline,
						"privacy_policy_url":   portal.BusinessProfile.PrivacyPolicyURL,
						"terms_of_service_url": portal.BusinessProfile.TermsOfServiceURL,
					},
				}
			}
			return nil
		}()),
		d.Set("default_return_url", portal.DefaultReturnURL),
		d.Set("login_page", func() []map[string]interface{} {
			if portal.LoginPage != nil {
				return []map[string]interface{}{
					{
						"enabled": portal.LoginPage.Enabled,
						"url":     portal.LoginPage.URL,
					},
				}
			}
			return nil
		}()),
		d.Set("features", func() []map[string]interface{} {
			if portal.Features != nil {
				featureMap := make(map[string]interface{})
				if portal.Features.CustomerUpdate != nil {
					featureMap["customer_update"] = []map[string]interface{}{
						{
							"enabled":         portal.Features.CustomerUpdate.Enabled,
							"allowed_updates": portal.Features.CustomerUpdate.AllowedUpdates,
						},
					}
				}
				if portal.Features.InvoiceHistory != nil {
					featureMap["invoice_history"] = []map[string]interface{}{
						{
							"enabled": portal.Features.InvoiceHistory.Enabled,
						},
					}
				}
				if portal.Features.PaymentMethodUpdate != nil {
					featureMap["payment_method_update"] = []map[string]interface{}{
						{
							"enabled": portal.Features.PaymentMethodUpdate.Enabled,
						},
					}
				}
				if portal.Features.SubscriptionCancel != nil {
					subsCancelMap := map[string]interface{}{
						"enabled":            portal.Features.SubscriptionCancel.Enabled,
						"mode":               portal.Features.SubscriptionCancel.Mode,
						"proration_behavior": portal.Features.SubscriptionCancel.ProrationBehavior,
					}
					if portal.Features.SubscriptionCancel.CancellationReason != nil {
						subsCancelMap["cancellation_reason"] = []map[string]interface{}{
							{
								"enabled": portal.Features.SubscriptionCancel.CancellationReason.Enabled,
								"options": portal.Features.SubscriptionCancel.CancellationReason.Options,
							},
						}
					}
					featureMap["subscription_cancel"] = []map[string]interface{}{
						subsCancelMap,
					}
				}
				if portal.Features.SubscriptionPause != nil {
					featureMap["subscription_pause"] = []map[string]interface{}{
						{
							"enabled": portal.Features.SubscriptionPause.Enabled,
						},
					}
				}
				if portal.Features.SubscriptionUpdate != nil {
					subsUpdateMap := map[string]interface{}{
						"enabled":                 portal.Features.SubscriptionUpdate.Enabled,
						"default_allowed_updates": portal.Features.SubscriptionUpdate.DefaultAllowedUpdates,
						"proration_behavior":      portal.Features.SubscriptionUpdate.ProrationBehavior,
					}
					var products []map[string]interface{}
					for _, p := range portal.Features.SubscriptionUpdate.Products {
						products = append(products, map[string]interface{}{
							"prices":  p.Prices,
							"product": p.Product,
						})
					}
					if products != nil {
						subsUpdateMap["products"] = products
					}

					featureMap["subscription_update"] = []map[string]interface{}{
						subsUpdateMap,
					}
				}
				return []map[string]interface{}{
					featureMap,
				}
			}
			return nil
		}()),
		d.Set("metadata", portal.Metadata),
	)
}

func resourceStripePortalConfigurationCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	var portal *stripe.BillingPortalConfiguration
	var err error

	params := &stripe.BillingPortalConfigurationParams{}

	if businessProfile, set := d.GetOk("business_profile"); set {
		for _, businessProfileMap := range ToMapSlice(businessProfile) {
			params.BusinessProfile = &stripe.BillingPortalConfigurationBusinessProfileParams{}
			for k, v := range businessProfileMap {
				switch k {
				case "headline":
					params.BusinessProfile.Headline = NonZeroString(v)
				case "privacy_policy_url":
					params.BusinessProfile.PrivacyPolicyURL = NonZeroString(v)
				case "terms_of_service_url":
					params.BusinessProfile.TermsOfServiceURL = NonZeroString(v)
				}
			}
		}
	}

	if defaultReturnURL, set := d.GetOk("default_return_url"); set {
		params.DefaultReturnURL = stripe.String(ToString(defaultReturnURL))
	}

	if loginPage, set := d.GetOk("login_page"); set {
		for _, loginPageMap := range ToMapSlice(loginPage) {
			for k, v := range loginPageMap {
				switch k {
				case "enabled":
					params.LoginPage = &stripe.BillingPortalConfigurationLoginPageParams{Enabled: stripe.Bool(ToBool(v))}
				}
			}
		}
	}

	if features, set := d.GetOk("features"); set {
		for _, featureMap := range ToMapSlice(features) {
			params.Features = &stripe.BillingPortalConfigurationFeaturesParams{}
			for k, v := range featureMap {
				switch k {
				case "customer_update":
					for _, customerUpdateMap := range ToMapSlice(v) {
						params.Features.CustomerUpdate = &stripe.BillingPortalConfigurationFeaturesCustomerUpdateParams{}
						for k, v := range customerUpdateMap {
							switch k {
							case "enabled":
								params.Features.CustomerUpdate.Enabled = stripe.Bool(ToBool(v))
							case "allowed_updates":
								params.Features.CustomerUpdate.AllowedUpdates = stripe.StringSlice(ToStringSlice(v))
							}
						}
					}
				case "invoice_history":
					for _, invoiceHistoryMap := range ToMapSlice(v) {
						params.Features.InvoiceHistory = &stripe.BillingPortalConfigurationFeaturesInvoiceHistoryParams{}
						for k, v := range invoiceHistoryMap {
							switch k {
							case "enabled":
								params.Features.InvoiceHistory.Enabled = stripe.Bool(ToBool(v))
							}
						}
					}
				case "payment_method_update":
					for _, paymentMethodUpdateMap := range ToMapSlice(v) {
						params.Features.PaymentMethodUpdate = &stripe.BillingPortalConfigurationFeaturesPaymentMethodUpdateParams{}
						for k, v := range paymentMethodUpdateMap {
							switch k {
							case "enabled":
								params.Features.PaymentMethodUpdate.Enabled = stripe.Bool(ToBool(v))
							}
						}
					}
				case "subscription_cancel":
					for _, subsCancelMap := range ToMapSlice(v) {
						params.Features.SubscriptionCancel = &stripe.BillingPortalConfigurationFeaturesSubscriptionCancelParams{}
						for k, v := range subsCancelMap {
							switch k {
							case "enabled":
								params.Features.SubscriptionCancel.Enabled = stripe.Bool(ToBool(v))
							case "mode":
								params.Features.SubscriptionCancel.Mode = NonZeroString(v)
							case "proration_behavior":
								params.Features.SubscriptionCancel.ProrationBehavior = NonZeroString(v)
							case "cancellation_reason":
								for _, cancellationReasonMap := range ToMapSlice(v) {
									params.Features.SubscriptionCancel.CancellationReason = &stripe.BillingPortalConfigurationFeaturesSubscriptionCancelCancellationReasonParams{}
									for k, v := range cancellationReasonMap {
										switch k {
										case "enabled":
											params.Features.SubscriptionCancel.CancellationReason.Enabled = stripe.Bool(ToBool(v))
										case "options":
											params.Features.SubscriptionCancel.CancellationReason.Options = stripe.StringSlice(ToStringSlice(v))
										}
									}
								}
							}
						}
					}
				case "subscription_pause":
					for _, subsPauseMap := range ToMapSlice(v) {
						params.Features.SubscriptionPause = &stripe.BillingPortalConfigurationFeaturesSubscriptionPauseParams{}
						for k, v := range subsPauseMap {
							switch k {
							case "enabled":
								params.Features.SubscriptionPause.Enabled = stripe.Bool(ToBool(v))
							}
						}
					}
				case "subscription_update":
					for _, subsUpdateMap := range ToMapSlice(v) {
						params.Features.SubscriptionUpdate = &stripe.BillingPortalConfigurationFeaturesSubscriptionUpdateParams{}
						for k, v := range subsUpdateMap {
							switch k {
							case "enabled":
								params.Features.SubscriptionUpdate.Enabled = stripe.Bool(ToBool(v))
							case "default_allowed_updates":
								params.Features.SubscriptionUpdate.DefaultAllowedUpdates = stripe.StringSlice(ToStringSlice(v))
							case "proration_behavior":
								params.Features.SubscriptionUpdate.ProrationBehavior = NonZeroString(v)
							case "products":
								for _, productMap := range ToMapSlice(v) {
									product := &stripe.BillingPortalConfigurationFeaturesSubscriptionUpdateProductParams{}
									for k, v := range productMap {
										switch k {
										case "product":
											product.Product = stripe.String(ToString(v))
										case "prices":
											product.Prices = stripe.StringSlice(ToStringSlice(v))
										}
									}
									params.Features.SubscriptionUpdate.Products = append(params.Features.SubscriptionUpdate.Products, product)
								}
							}
						}
					}
				}
			}
		}
	}

	if meta, set := d.GetOk("metadata"); set {
		for k, v := range ToMap(meta) {
			params.AddMetadata(k, ToString(v))
		}
	}

	err = retryWithBackOff(func() error {
		portal, err = c.BillingPortalConfigurations.New(params)
		return err
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(portal.ID)
	return resourceStripePortalConfigurationRead(ctx, d, m)
}

func resourceStripePortalConfigurationUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	var err error

	params := &stripe.BillingPortalConfigurationParams{}
	if d.HasChange("active") {
		params.Active = stripe.Bool(ExtractBool(d, "active"))
	}

	if d.HasChange("business_profile") {
		for _, businessProfileMap := range ExtractMapSlice(d, "business_profile") {
			params.BusinessProfile = &stripe.BillingPortalConfigurationBusinessProfileParams{}
			for k, v := range businessProfileMap {
				switch k {
				case "headline":
					params.BusinessProfile.Headline = NonZeroString(v)
				case "privacy_policy_url":
					params.BusinessProfile.PrivacyPolicyURL = NonZeroString(v)
				case "terms_of_service_url":
					params.BusinessProfile.TermsOfServiceURL = NonZeroString(v)
				}
			}
		}
	}

	if d.HasChange("default_return_url") {
		params.DefaultReturnURL = stripe.String(ExtractString(d, "default_return_url"))
	}

	if d.HasChange("login_page") {
		for _, loginPageMap := range ExtractMapSlice(d, "login_page") {
			for k, v := range loginPageMap {
				switch k {
				case "enabled":
					params.LoginPage = &stripe.BillingPortalConfigurationLoginPageParams{Enabled: stripe.Bool(ToBool(v))}
				}
			}
		}
	}

	if d.HasChange("features") {
		for _, featureMap := range ExtractMapSlice(d, "features") {
			params.Features = &stripe.BillingPortalConfigurationFeaturesParams{}
			for k, v := range featureMap {
				switch k {
				case "customer_update":
					for _, customerUpdateMap := range ToMapSlice(v) {
						params.Features.CustomerUpdate = &stripe.BillingPortalConfigurationFeaturesCustomerUpdateParams{}
						for k, v := range customerUpdateMap {
							switch k {
							case "enabled":
								params.Features.CustomerUpdate.Enabled = stripe.Bool(ToBool(v))
							case "allowed_updates":
								params.Features.CustomerUpdate.AllowedUpdates = stripe.StringSlice(ToStringSlice(v))
							}
						}
					}
				case "invoice_history":
					for _, invoiceHistoryMap := range ToMapSlice(v) {
						params.Features.InvoiceHistory = &stripe.BillingPortalConfigurationFeaturesInvoiceHistoryParams{}
						for k, v := range invoiceHistoryMap {
							switch k {
							case "enabled":
								params.Features.InvoiceHistory.Enabled = stripe.Bool(ToBool(v))
							}
						}
					}
				case "payment_method_update":
					for _, paymentMethodUpdateMap := range ToMapSlice(v) {
						params.Features.PaymentMethodUpdate = &stripe.BillingPortalConfigurationFeaturesPaymentMethodUpdateParams{}
						for k, v := range paymentMethodUpdateMap {
							switch k {
							case "enabled":
								params.Features.PaymentMethodUpdate.Enabled = stripe.Bool(ToBool(v))
							}
						}
					}
				case "subscription_cancel":
					for _, subsCancelMap := range ToMapSlice(v) {
						params.Features.SubscriptionCancel = &stripe.BillingPortalConfigurationFeaturesSubscriptionCancelParams{}
						for k, v := range subsCancelMap {
							switch k {
							case "enabled":
								params.Features.SubscriptionCancel.Enabled = stripe.Bool(ToBool(v))
							case "mode":
								params.Features.SubscriptionCancel.Mode = NonZeroString(v)
							case "proration_behavior":
								params.Features.SubscriptionCancel.ProrationBehavior = NonZeroString(v)
							case "cancellation_reason":
								for _, cancellationReasonMap := range ToMapSlice(v) {
									params.Features.SubscriptionCancel.CancellationReason = &stripe.BillingPortalConfigurationFeaturesSubscriptionCancelCancellationReasonParams{}
									for k, v := range cancellationReasonMap {
										switch k {
										case "enabled":
											params.Features.SubscriptionCancel.CancellationReason.Enabled = stripe.Bool(ToBool(v))
										case "options":
											params.Features.SubscriptionCancel.CancellationReason.Options = stripe.StringSlice(ToStringSlice(v))
										}
									}
								}
							}
						}
					}
				case "subscription_pause":
					for _, subsPauseMap := range ToMapSlice(v) {
						params.Features.SubscriptionPause = &stripe.BillingPortalConfigurationFeaturesSubscriptionPauseParams{}
						for k, v := range subsPauseMap {
							switch k {
							case "enabled":
								params.Features.SubscriptionPause.Enabled = stripe.Bool(ToBool(v))
							}
						}
					}
				case "subscription_update":
					for _, subsUpdateMap := range ToMapSlice(v) {
						params.Features.SubscriptionUpdate = &stripe.BillingPortalConfigurationFeaturesSubscriptionUpdateParams{}
						for k, v := range subsUpdateMap {
							switch k {
							case "enabled":
								params.Features.SubscriptionUpdate.Enabled = stripe.Bool(ToBool(v))
							case "default_allowed_updates":
								params.Features.SubscriptionUpdate.DefaultAllowedUpdates = stripe.StringSlice(ToStringSlice(v))
							case "proration_behavior":
								params.Features.SubscriptionUpdate.ProrationBehavior = NonZeroString(v)
							case "products":
								for _, productMap := range ToMapSlice(v) {
									product := &stripe.BillingPortalConfigurationFeaturesSubscriptionUpdateProductParams{}
									for k, v := range productMap {
										switch k {
										case "product":
											product.Product = stripe.String(ToString(v))
										case "prices":
											product.Prices = stripe.StringSlice(ToStringSlice(v))
										}
									}
									params.Features.SubscriptionUpdate.Products = append(params.Features.SubscriptionUpdate.Products, product)
								}
							}
						}
					}
				}
			}
		}
	}

	if d.HasChange("metadata") {
		params.Metadata = nil
		UpdateMetadata(d, params)
	}

	err = retryWithBackOff(func() error {
		_, err = c.BillingPortalConfigurations.Update(d.Id(), params)
		return err
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceStripePortalConfigurationRead(ctx, d, m)
}

func resourceStripePortalConfigurationDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[WARN] Stripe doesn't support deletion of customer portals. Portal will be deactivated but not deleted")

	c := m.(*client.API)
	var err error

	params := stripe.BillingPortalConfigurationParams{
		Active: stripe.Bool(false),
	}

	err = retryWithBackOff(func() error {
		_, err = c.BillingPortalConfigurations.Update(d.Id(), &params)
		return err
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
