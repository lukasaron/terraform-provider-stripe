package stripe

//import (
//	"context"
//
//	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
//	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
//	"github.com/stripe/stripe-go/v72"
//	"github.com/stripe/stripe-go/v72/client"
//)
//
//func dataSourceStripeBalanceTransactions() *schema.Resource {
//	return &schema.Resource{
//		ReadContext: dataSourceStripeBalanceTransactionsRead,
//		Schema: map[string]*schema.Schema{
//			"type": {
//				Type:     schema.TypeString,
//				Optional: true,
//			},
//			"payout": {
//				Type:     schema.TypeString,
//				Optional: true,
//			},
//			"available": {
//				Type:     schema.TypeList,
//				Optional: true,
//				MaxItems: 1,
//				Elem: &schema.Resource{
//					Schema: map[string]*schema.Schema{
//						"condition": {
//							Type:         schema.TypeString,
//							Required:     true,
//							ExactlyOneOf: []string{"gt", "gte", "lt", "lte"},
//						},
//						"on": {
//							Type:     schema.TypeString,
//							Required: true,
//						},
//					},
//				},
//			},
//			"created": {
//				Type:     schema.TypeList,
//				Optional: true,
//				Elem: &schema.Resource{
//					Schema: map[string]*schema.Schema{
//						"condition": {
//							Type:         schema.TypeString,
//							Required:     true,
//							ExactlyOneOf: []string{"gt", "gte", "lt", "lte"},
//						},
//						"on": {
//							Type:     schema.TypeString,
//							Required: true,
//						},
//					},
//				},
//			},
//			"currency": {
//				Type:     schema.TypeString,
//				Optional: true,
//			},
//			"limit":
//		},
//	}
//}
//
//func dataSourceStripeBalanceTransactionsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
//	c := m.(*client.API)
//
//	c.BalanceTransaction.List(&stripe.BalanceTransactionListParams{
//		ListParams:       stripe.ListParams{},
//		AvailableOn:      nil,
//		AvailableOnRange: nil,
//		Created:          nil,
//		CreatedRange:     nil,
//		Currency:         nil,
//		Payout:           nil,
//		Source:           nil,
//		Type:             nil,
//	})
//}
