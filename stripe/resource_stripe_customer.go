package stripe

import (
	"context"
	"github.com/stripe/stripe-go/v78"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stripe/stripe-go/v78/client"
)

func resourceStripeCustomer() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceStripeCustomerRead,
		CreateContext: resourceStripeCustomerCreate,
		UpdateContext: resourceStripeCustomerUpdate,
		DeleteContext: resourceStripeCustomerDelete,
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
				Description: "The customer’s full name or business name.",
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "Customer’s email address. " +
					"It’s displayed alongside the customer in your dashboard and can be useful for searching " +
					"and tracking. This may be up to 512 characters.",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "An arbitrary string that you can attach to a customer object. " +
					"It is displayed alongside the customer in the dashboard.",
			},
			"phone": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The customer’s phone number.",
			},
			"address": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Address map with fields related to the address: line1, line2, city, state, " +
					"postal_code and country",
			},
			"shipping": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Description: "Shipping map with fields like name, phone and fields related to the address: " +
					"line1, line2, city, state, postal_code and country. ",
			},
			"balance": {
				Type:     schema.TypeInt,
				Optional: true,
				Description: "An integer amount in cents that represents the customer’s current balance, " +
					"which affect the customer’s future invoices. " +
					"A negative amount represents a credit that decreases the amount due on an invoice; " +
					"a positive amount increases the amount due on an invoice.",
			},
			//TODO "coupon": {
			//	Type:     schema.TypeString,
			//	Optional: true,
			//	Description: "If you provide a coupon code, " +
			//		"the customer will have a discount applied on all recurring charges. " +
			//		"Charges you create through the API will not have the discount.",
			//},
			//TODO "promotion_code": {
			//	Type:     schema.TypeString,
			//	Optional: true,
			//	Description: "The API ID of a promotion code to apply to the customer. " +
			//		"The customer will have a discount applied on all recurring payments. " +
			//		"Charges you create through the API will not have the discount.",
			//},
			"default_invoice_prefix": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "The default (auto-generated) prefix for the customer used to generate unique" +
					" invoice numbers. ",
			},
			"invoice_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				Description: "The prefix for the customer used to generate unique invoice numbers. " +
					"Must be 3–12 uppercase letters or numbers.",
			},
			"invoice_settings": {
				Type:        schema.TypeMap,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Default invoice settings for this customer.",
			},
			"next_invoice_sequence": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     1,
				Description: "The sequence to be used on the customer’s next invoice. Defaults to 1.",
			},
			"preferred_locales": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Customer’s preferred languages, ordered by preference.",
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

func resourceStripeCustomerRead(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	var customer *stripe.Customer
	var err error

	err = retryWithBackOff(func() error {
		customer, err = c.Customers.Get(d.Id(), nil)
		return err
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return CallSet(
		d.Set("name", customer.Name),
		d.Set("email", customer.Email),
		d.Set("description", customer.Description),
		d.Set("phone", customer.Phone),
		func() error {
			addressMap := make(map[string]interface{})
			if customer.Address.Line1 != "" {
				addressMap["line1"] = customer.Address.Line1
			}
			if customer.Address.Line2 != "" {
				addressMap["line2"] = customer.Address.Line2
			}
			if customer.Address.City != "" {
				addressMap["city"] = customer.Address.City
			}
			if customer.Address.State != "" {
				addressMap["state"] = customer.Address.State
			}
			if customer.Address.PostalCode != "" {
				addressMap["postal_code"] = customer.Address.PostalCode
			}
			if customer.Address.Country != "" {
				addressMap["country"] = customer.Address.Country
			}

			return d.Set("address", addressMap)
		}(),
		func() error {
			if customer.Shipping != nil {
				shippingMap := make(map[string]interface{})
				if customer.Shipping.Name != "" {
					shippingMap["name"] = customer.Shipping.Name
				}
				if customer.Shipping.Phone != "" {
					shippingMap["phone"] = customer.Shipping.Phone
				}
				if customer.Shipping.Address.Line1 != "" {
					shippingMap["line1"] = customer.Shipping.Address.Line1
				}
				if customer.Shipping.Address.Line2 != "" {
					shippingMap["line2"] = customer.Shipping.Address.Line2
				}
				if customer.Shipping.Address.City != "" {
					shippingMap["city"] = customer.Shipping.Address.City
				}
				if customer.Shipping.Address.State != "" {
					shippingMap["state"] = customer.Shipping.Address.State
				}
				if customer.Shipping.Address.PostalCode != "" {
					shippingMap["postal_code"] = customer.Shipping.Address.PostalCode
				}
				if customer.Shipping.Address.Country != "" {
					shippingMap["country"] = customer.Shipping.Address.Country
				}
				return d.Set("shipping", shippingMap)
			}
			return nil
		}(),
		d.Set("balance", customer.Balance),
		func() error {
			if _, set := d.GetOk("invoice_prefix"); set {
				return d.Set("invoice_prefix", customer.InvoicePrefix)
			}
			return nil
		}(),
		d.Set("default_invoice_prefix", customer.InvoicePrefix),
		func() error {
			if customer.InvoiceSettings != nil {
				invoiceSettingsMap := make(map[string]interface{})
				if customer.InvoiceSettings.Footer != "" {
					invoiceSettingsMap["footer"] = customer.InvoiceSettings.Footer
				}
				if customer.InvoiceSettings.DefaultPaymentMethod != nil {
					invoiceSettingsMap["default_payment_method"] = customer.InvoiceSettings.DefaultPaymentMethod.ID
				}
				for _, field := range customer.InvoiceSettings.CustomFields {
					invoiceSettingsMap[field.Name] = field.Value
				}

				return d.Set("invoice_settings", invoiceSettingsMap)
			}
			return nil
		}(),
		d.Set("next_invoice_sequence", customer.NextInvoiceSequence),
		d.Set("preferred_locales", customer.PreferredLocales),
		d.Set("metadata", customer.Metadata),
	)
}

func resourceStripeCustomerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	var customer *stripe.Customer
	var err error

	params := &stripe.CustomerParams{}

	if name, set := d.GetOk("name"); set {
		params.Name = stripe.String(ToString(name))
	}
	if email, set := d.GetOk("email"); set {
		params.Email = stripe.String(ToString(email))
	}
	if description, set := d.GetOk("description"); set {
		params.Description = stripe.String(ToString(description))
	}
	if phone, set := d.GetOk("phone"); set {
		params.Phone = stripe.String(ToString(phone))
	}
	if address, set := d.GetOk("address"); set {
		params.Address = &stripe.AddressParams{}
		addressMap := ToMap(address)
		for k, v := range addressMap {
			value := stripe.String(ToString(v))
			switch k {
			case "line1":
				params.Address.Line1 = value
			case "line2":
				params.Address.Line2 = value
			case "city":
				params.Address.City = value
			case "state":
				params.Address.State = value
			case "postal_code":
				params.Address.PostalCode = value
			case "country":
				params.Address.Country = value
			}
		}
	}
	if shipping, set := d.GetOk("shipping"); set {
		params.Shipping = &stripe.CustomerShippingParams{
			Address: &stripe.AddressParams{},
		}
		shippingMap := ToMap(shipping)
		for k, v := range shippingMap {
			value := stripe.String(ToString(v))
			switch k {
			case "name":
				params.Shipping.Name = value
			case "phone":
				params.Shipping.Phone = value
			case "line1":
				params.Shipping.Address.Line1 = value
			case "line2":
				params.Shipping.Address.Line2 = value
			case "city":
				params.Shipping.Address.City = value
			case "state":
				params.Shipping.Address.State = value
			case "postal_code":
				params.Shipping.Address.PostalCode = value
			case "country":
				params.Shipping.Address.Country = value
			}
		}
	}
	if balance, set := d.GetOk("balance"); set {
		params.Balance = stripe.Int64(ToInt64(balance))
	}
	if invoicePrefix, set := d.GetOk("invoice_prefix"); set {
		params.InvoicePrefix = stripe.String(ToString(invoicePrefix))
	}
	if invoiceSettings, set := d.GetOk("invoice_settings"); set {
		params.InvoiceSettings = &stripe.CustomerInvoiceSettingsParams{}
		invoiceSettingsMap := ToMap(invoiceSettings)
		for k, v := range invoiceSettingsMap {
			value := stripe.String(ToString(v))
			switch k {
			case "default_payment_method":
				params.InvoiceSettings.DefaultPaymentMethod = value
			case "footer":
				params.InvoiceSettings.Footer = value
			default:
				params.InvoiceSettings.CustomFields = append(params.InvoiceSettings.CustomFields,
					&stripe.CustomerInvoiceSettingsCustomFieldParams{
						Name:  stripe.String(k),
						Value: value,
					})
			}
		}
	}
	if nextInvoiceSequence, set := d.GetOk("next_invoice_sequence"); set {
		params.NextInvoiceSequence = stripe.Int64(ToInt64(nextInvoiceSequence))
	}
	if preferredLocales, set := d.GetOk("preferred_locales"); set {
		params.PreferredLocales = stripe.StringSlice(ToStringSlice(preferredLocales))
	}
	if meta, set := d.GetOk("metadata"); set {
		for k, v := range ToMap(meta) {
			params.AddMetadata(k, ToString(v))
		}
	}

	err = retryWithBackOff(func() error {
		customer, err = c.Customers.New(params)
		return err
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(customer.ID)
	return resourceStripeCustomerRead(ctx, d, m)
}

func resourceStripeCustomerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	var err error

	params := &stripe.CustomerParams{}

	if d.HasChange("name") {
		params.Name = stripe.String(ExtractString(d, "name"))
	}
	if d.HasChange("email") {
		params.Email = stripe.String(ExtractString(d, "email"))
	}
	if d.HasChange("description") {
		params.Description = stripe.String(ExtractString(d, "description"))
	}
	if d.HasChange("phone") {
		params.Phone = stripe.String(ExtractString(d, "phone"))
	}
	if d.HasChange("address") {
		params.Address = &stripe.AddressParams{}
		addressMap := ExtractMap(d, "address")
		for k, v := range addressMap {
			value := stripe.String(ToString(v))
			switch k {
			case "line1":
				params.Address.Line1 = value
			case "line2":
				params.Address.Line2 = value
			case "city":
				params.Address.City = value
			case "state":
				params.Address.State = value
			case "postal_code":
				params.Address.PostalCode = value
			case "country":
				params.Address.Country = value
			}
		}
	}
	if d.HasChange("shipping") {
		params.Shipping = &stripe.CustomerShippingParams{
			Address: &stripe.AddressParams{},
		}
		shippingMap := ExtractMap(d, "shipping")
		for k, v := range shippingMap {
			value := stripe.String(ToString(v))
			switch k {
			case "name":
				params.Shipping.Name = value
			case "phone":
				params.Shipping.Phone = value
			case "line1":
				params.Shipping.Address.Line1 = value
			case "line2":
				params.Shipping.Address.Line2 = value
			case "city":
				params.Shipping.Address.City = value
			case "state":
				params.Shipping.Address.State = value
			case "postal_code":
				params.Shipping.Address.PostalCode = value
			case "country":
				params.Shipping.Address.Country = value
			}
		}
	}
	if d.HasChange("balance") {
		params.Balance = stripe.Int64(ExtractInt64(d, "balance"))
	}
	if d.HasChange("invoice_prefix") {
		params.InvoicePrefix = stripe.String(ExtractString(d, "invoice_prefix"))
	}
	if d.HasChange("invoice_settings") {
		params.InvoiceSettings = &stripe.CustomerInvoiceSettingsParams{}
		invoiceSettingsMap := ExtractMap(d, "invoice_settings")
		for k, v := range invoiceSettingsMap {
			value := stripe.String(ToString(v))
			switch k {
			case "default_payment_method":
				params.InvoiceSettings.DefaultPaymentMethod = value
			case "footer":
				params.InvoiceSettings.Footer = value
			default:
				params.InvoiceSettings.CustomFields = append(params.InvoiceSettings.CustomFields,
					&stripe.CustomerInvoiceSettingsCustomFieldParams{
						Name:  stripe.String(k),
						Value: value,
					})
			}
		}
	}
	if d.HasChange("next_invoice_sequence") {
		params.NextInvoiceSequence = stripe.Int64(ExtractInt64(d, "next_invoice_sequence"))
	}
	if d.HasChange("preferred_locales") {
		params.PreferredLocales = stripe.StringSlice(ExtractStringSlice(d, "preferred_locales"))
	}
	if d.HasChange("metadata") {
		params.Metadata = nil
		UpdateMetadata(d, params)
	}

	err = retryWithBackOff(func() error {
		_, err = c.Customers.Update(d.Id(), params)
		return err
	})
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceStripeCustomerRead(ctx, d, m)

}

func resourceStripeCustomerDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*client.API)
	var err error

	err = retryWithBackOff(func() error {
		_, err = c.Customers.Del(d.Id(), nil)
		return err
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
