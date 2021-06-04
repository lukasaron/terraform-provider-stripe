package stripe

//
//import (
//	"context"
//	"strconv"
//	"time"
//
//	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
//	"github.com/stripe/stripe-go/v72"
//	"github.com/stripe/stripe-go/v72/client"
//)
//
//func dataSourceStripeWebhookEndpoints() *schema.Resource {
//	return &schema.Resource{
//		ReadContext: dataSourceStripeWebhookEndpointRead,
//		Schema: map[string]*schema.Schema{
//			"webhook_endpoints": {
//				Type:     schema.TypeList,
//				Computed: true,
//				Elem: &schema.Resource{
//					Schema: map[string]*schema.Schema{
//						"id": {
//							Type:     schema.TypeString,
//							Computed: true,
//						},
//						"object": {
//							Type:     schema.TypeString,
//							Computed: true,
//						},
//						"description": {
//							Type:     schema.TypeString,
//							Computed: true,
//						},
//						"enabled_events": {
//							Type:     schema.TypeList,
//							Computed: true,
//							Elem:     &schema.Schema{Type: schema.TypeString},
//						},
//						"status": {
//							Type:     schema.TypeString,
//							Computed: true,
//						},
//						"url": {
//							Type:     schema.TypeString,
//							Computed: true,
//						},
//						"metadata": {
//							Type:     schema.TypeMap,
//							Computed: true,
//							Elem:     &schema.Schema{Type: schema.TypeString},
//						},
//						"application": {
//							Type:     schema.TypeString,
//							Computed: true,
//						},
//						"created": {
//							Type:     schema.TypeFloat,
//							Computed: true,
//						},
//						"api_version": {
//							Type:     schema.TypeString,
//							Computed: true,
//						},
//						"livemode": {
//							Type:     schema.TypeBool,
//							Computed: true,
//						},
//					},
//				},
//			},
//		},
//	}
//}
//
//func dataSourceStripeWebhookEndpointRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
//	c := m.(*client.API)
//
//	iter := c.WebhookEndpoints.List(&stripe.WebhookEndpointListParams{})
//	if iter.Err() != nil {
//		return diag.FromErr(iter.Err())
//	}
//
//	webhookEndpoints := make([]map[string]interface{}, 0)
//	for iter.Next() {
//		current := iter.Current().(*stripe.WebhookEndpoint)
//		webhookEndpoint := make(map[string]interface{})
//
//		webhookEndpoint["id"] = current.ID
//		webhookEndpoint["object"] = current.Object
//		webhookEndpoint["description"] = current.Description
//		webhookEndpoint["enabled_events"] = current.EnabledEvents
//		webhookEndpoint["status"] = current.Status
//		webhookEndpoint["url"] = current.URL
//		webhookEndpoint["metadata"] = current.Metadata
//		webhookEndpoint["application"] = current.Application
//		webhookEndpoint["created"] = current.Created
//		webhookEndpoint["api_version"] = current.APIVersion
//		webhookEndpoint["livemode"] = current.Livemode
//
//		webhookEndpoints = append(webhookEndpoints, webhookEndpoint)
//	}
//
//	if err := d.Set("webhook_endpoints", webhookEndpoints); err != nil {
//		return diag.FromErr(err)
//	}
//	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
//
//	return nil
//}
