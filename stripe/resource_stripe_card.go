package stripe

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/client"
)

func resourceStripeCard() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceStripeCardRead,
		CreateContext: resourceStripeCardCreate,
		UpdateContext: resourceStripeCardUpdate,
		DeleteContext: resourceStripeCardDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Unique identifier for the object.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Cardholder name.",
			},
			"customer": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
				Description: "The customer that this card belongs to. ",
			},
			"number": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "The card number, as a string without any separators.",
			},
			"exp_month": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Two-digit number representing the card's expiration month.",
			},
			"exp_year": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Four-digit number representing the card's expiration year.",
			},
			"cvc": {
				Type:      schema.TypeInt,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
				Description: "Card security code. Highly recommended to always include this value, " +
					"but it's required only for accounts based in European countries.",
			},
			"address": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Address map with fields related to the address: line1, line2, city, state, " +
					"zip and country",
			},
			"address_line1_check": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "If address_line1 was provided, results of the check: pass, fail, " +
					"unavailable, or unchecked.",
			},
			"address_zip_check": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "If address_zip was provided, results of the check: pass, fail, unavailable, " +
					"or unchecked.",
			},
			"brand": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "Card brand. Can be American Express, Diners Club, Discover, JCB, MasterCard, UnionPay, " +
					"Visa, or Unknown.",
			},
			"country": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "Two-letter ISO code representing the country of the card. " +
					"You could use this attribute to get a sense of the international " +
					"breakdown of cards you’ve collected.",
			},
			"cvc_check": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "If a CVC was provided, results of the check: pass, fail, unavailable, or unchecked. " +
					"A result of unchecked indicates that CVC was provided but hasn’t been checked yet.",
			},
			"fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "Uniquely identifies this particular card number. " +
					"You can use this attribute to check whether two customers who’ve signed up with you are using " +
					"the same card number, for example. For payment methods that tokenize card information " +
					"(Apple Pay, Google Pay), the tokenized number might be provided " +
					"instead of the underlying card number.",
			},
			"funding": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Card funding type. Can be credit, debit, prepaid, or unknown.",
			},
			"last4": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last four digits of the card.",
			},
			"available_payout_methods": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Description: "A set of available payout methods for this card. " +
					"Only values from this set should be passed as the method when creating a payout.",
			},
			"tokenization_method": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "If the card number is tokenized, " +
					"this is the method that was used. Can be android_pay (includes Google Pay), apple_pay, " +
					"masterpass, visa_checkout, or null.",
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

func resourceStripeCardRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	params := &stripe.CardParams{
		Customer: stripe.String(ExtractString(d, "customer")),
	}

	card, err := c.Cards.Get(d.Id(), params)
	if err != nil {
		return diag.FromErr(err)
	}

	return CallSet(
		d.Set("name", card.Name),
		d.Set("customer", card.Customer.ID),
		d.Set("exp_month", int(card.ExpMonth)),
		d.Set("exp_year", int(card.ExpYear)),
		d.Set("address", map[string]interface{}{
			"line1":   card.AddressLine1,
			"line2":   card.AddressLine2,
			"city":    card.AddressCity,
			"state":   card.AddressState,
			"zip":     card.AddressZip,
			"country": card.AddressCountry,
		}),
		d.Set("address_line1_check", card.AddressLine1Check),
		d.Set("address_zip_check", card.AddressZipCheck),
		d.Set("brand", card.Brand),
		d.Set("country", card.Country),
		d.Set("cvc_check", card.CVCCheck),
		d.Set("fingerprint", card.Fingerprint),
		d.Set("funding", card.Funding),
		d.Set("last4", card.Last4),
		d.Set("available_payout_methods", card.AvailablePayoutMethods),
		d.Set("tokenization_method", card.TokenizationMethod),
		d.Set("metadata", card.Metadata),
	)
}

func resourceStripeCardCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	params := &stripe.CardParams{
		Customer: stripe.String(ExtractString(d, "customer")),
		Number:   stripe.String(ExtractString(d, "number")),
		ExpMonth: stripe.String(fmt.Sprintf("%02d", ExtractInt(d, "exp_month"))),
		ExpYear:  stripe.String(fmt.Sprintf("%04d", ExtractInt(d, "exp_year"))),
	}
	if name, set := d.GetOk("name"); set {
		params.Name = stripe.String(ToString(name))
	}
	if cvc, set := d.GetOk("cvc"); set {
		params.CVC = stripe.String(fmt.Sprintf("%d", ToInt(cvc)))
	}
	if address, set := d.GetOk("address"); set {
		addressMap := ToMap(address)
		for k, v := range addressMap {
			value := stripe.String(ToString(v))
			switch k {
			case "line1":
				params.AddressLine1 = value
			case "line2":
				params.AddressLine2 = value
			case "city":
				params.AddressCity = value
			case "state":
				params.AddressState = value
			case "zip":
				params.AddressZip = value
			case "country":
				params.AddressCountry = value
			}
		}
	}
	if meta, set := d.GetOk("metadata"); set {
		for k, v := range ToMap(meta) {
			params.AddMetadata(k, ToString(v))
		}
	}

	card, err := c.Cards.New(params)
	if err != nil {
		return diag.FromErr(err)
	}

	dg := CallSet(
		d.Set("number", *params.Number),
		func() error {
			// CVC is an optional field it needs to be checked before it is set to the state.
			// Other option is to set this filed to the state when params are filled.
			// This seems clearer to have all state sets in one place.
			if cvc, set := d.GetOk("cvc"); set {
				return d.Set("cvc", ToInt(cvc))
			}
			return nil
		}(),
	)
	if len(dg) > 0 {
		return dg
	}

	d.SetId(card.ID)
	return resourceStripeCardRead(ctx, d, m)
}

func resourceStripeCardUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	params := &stripe.CardParams{
		Customer: stripe.String(ExtractString(d, "customer")),
	}

	if d.HasChange("name") {
		params.Name = stripe.String(ExtractString(d, "name"))
	}
	if d.HasChange("exp_month") {
		params.ExpMonth = stripe.String(fmt.Sprintf("%02d", ExtractInt(d, "exp_month")))
	}
	if d.HasChange("exp_year") {
		params.ExpYear = stripe.String(fmt.Sprintf("%04d", ExtractInt(d, "exp_year")))
	}
	if d.HasChange("address") {
		addressMap := ExtractMap(d, "address")
		for k, v := range addressMap {
			value := stripe.String(ToString(v))
			switch k {
			case "line1":
				params.AddressLine1 = value
			case "line2":
				params.AddressLine2 = value
			case "city":
				params.AddressCity = value
			case "state":
				params.AddressState = value
			case "zip":
				params.AddressZip = value
			case "country":
				params.AddressCountry = value
			}
		}
	}
	if d.HasChange("metadata") {
		params.Metadata = nil
		UpdateMetadata(d, params)
	}

	_, err := c.Cards.Update(d.Id(), params)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceStripeCardRead(ctx, d, m)
}

func resourceStripeCardDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	params := &stripe.CardParams{
		Customer: stripe.String(ExtractString(d, "customer")),
	}

	_, err := c.Cards.Del(d.Id(), params)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
