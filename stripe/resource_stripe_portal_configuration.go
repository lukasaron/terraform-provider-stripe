package stripe

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/client"
)

func resourceStripePortalConfiguration() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceStripePortalConfigurationRead,
		DeleteContext: resourceStripePortalConfigurationDelete,
		Schema: map[string]*schema.Schema{
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
							Description: "A link to the business’s publicly available privacy policy.",
						},
						"terms_of_service_url": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "A link to the business’s publicly available terms of service.",
						},
					},
				},
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
										Description: "Whether to cancel subscriptions immediately or at the end of the billing period.",
										ValidateFunc: validation.StringInSlice([]string{"immediately", "at_period_end"}, false),
									},
									"proration_behavior": {
										Type:         schema.TypeString,
										Optional:     true,
										Description: "Whether to create prorations when canceling subscriptions.",
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
													Elem:     &schema.Schema{Type: schema.TypeString},
												},
												"product":{
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
			"default_return_url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The default URL to redirect customers to when they click on the portal’s " +
					"link to return to your website. This can be overriden when creating the session.",
			},
			"metadata": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Set of key-value pairs that you can attach to an object. " +
					"This can be useful for storing additional information about the object in a structured format.",
				Elem:        &schema.Schema{Type: schema.TypeString},
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

func resourceStripePortalConfigurationDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Println("[WARN] Stripe doesn't support deletion of customer portals through the API")
	d.SetId("")
	return nil
}
