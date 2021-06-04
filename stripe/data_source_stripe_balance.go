package stripe

import (
	"context"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/client"
)

func dataSourceStripeBalance() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceStripeBalanceRead,
		Schema: map[string]*schema.Schema{
			"livemode": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Has the value true if the object exists in live mode or the value false if the object exists in test mode.",
			},
			"available": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Funds that are available to be transferred or paid out, whether automatically by Stripe or explicitly via the Transfers API or Payouts API. The available balance for each currency and payment type can be found in the source_types property",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"amount": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Balance amount.",
						},
						"currency": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Three-letter ISO currency code, in lowercase. Must be a supported currency.",
						},
						"source_types": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Breakdown of balance by source types.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"pending": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "Funds that are not yet available in the balance, due to the 7-day rolling pay cycle. The pending balance for each currency, and for each payment type, can be found in the source_types property.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"amount": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Balance amount.",
						},
						"currency": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Three-letter ISO currency code, in lowercase. Must be a supported currency.",
						},
						"source_types": {
							Type:        schema.TypeMap,
							Computed:    true,
							Description: "Breakdown of balance by source types.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceStripeBalanceRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	balance, err := c.Balance.Get(&stripe.BalanceParams{})
	if err != nil {
		return diag.FromErr(err)
	}

	dg := CallSet(
		d.Set("livemode", balance.Livemode),
		d.Set("available", flattenAmountSlice(balance.Available)),
		d.Set("pending", flattenAmountSlice(balance.Pending)),
	)
	if len(dg) > 0 {
		return dg
	}

	d.SetId("balance")
	return nil
}

func flattenAmountSlice(amounts []*stripe.Amount) []map[string]interface{} {
	m := make([]map[string]interface{}, 0)
	for _, amount := range amounts {
		mItem := make(map[string]interface{})
		mItem["amount"] = strconv.FormatInt(amount.Value, 10)
		mItem["currency"] = amount.Currency
		mItem["source_types"] = flattenSourceTypes(amount.SourceTypes)
		m = append(m, mItem)
	}

	return m
}

func flattenSourceTypes(st map[stripe.BalanceSourceType]int64) map[string]string {
	m := make(map[string]string)
	for k, v := range st {
		m[string(k)] = strconv.FormatInt(v, 10)
	}
	return m
}
