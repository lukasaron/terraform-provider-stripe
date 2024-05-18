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
	// FIXME fill parameters

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
