---
layout: "stripe"
page_title: "Stripe: stripe_customer"
description: |- The Stripe Customer can be created, modified, configured and removed by this resource.
---

# stripe_customer

With this resource, you can create a customer - [Stripe API customer documentation](https://stripe.com/docs/api/customers).

Customer objects allow you to perform recurring charges, and to track multiple charges, 
that are associated with the same customer.

## Example Usage

```hcl
// A customer without details
resource "stripe_customer" "customer" {}

// A customer with some details
resource "stripe_customer" "customer" {
  name        = "Lukas Aron"
  email       = "lukas@aron.com"
  phone       = "+610123456789"
  description = "Terraform Customer"
}

// A customer with address
resource "stripe_customer" "customer" {
  name    = "Lukas Aron"
  address = {
    line1       = "1 The Best Street",
    line2       = "Apartment 401",
    city        = "Sydney",
    postal_code = "2000"
    country     = "AU"
    state       = "New South Wales"
  }
}

// A customer with other details
resource "stripe_customer" "customer" {
  name                  = "Lukas Aron"
  invoice_prefix        = "LA000"
  next_invoice_sequence = 1001
  balance               = 10000 // $100

  invoice_settings = {
    footer          = "--- Lukas Aron ---"
    customFieldName = "customFieldValue"
  }

  preferred_locales = ["eng", "esp"]

  shipping = {
    name        = "Lukas Aron",
    phone       = "+610123456789"
    line1       = "1 The Best Street",
    line2       = "Apartment 401",
    city        = "Sydney",
    postal_code = "2000"
    country     = "AU"
    state       = "New South Wales"
  }
}
```

## Argument Reference

Arguments accepted by this resource include:

* `name` - (Optional) String. The customer’s full name or business name.
* `email` - (Optional) String. Customer’s email address. It’s displayed alongside the customer in your dashboard and can be useful for searching and tracking. This may be up to 512 characters.
* `description` - (Optional) String. An arbitrary string that you can attach to a customer object. It is displayed alongside the customer in the dashboard.
* `phone` - (Optional) String. The customer’s phone number.
* `address` - (Optional) Map(String). The customer’s address, for all individual fields see: [Address Fields](#address-fields). 
* `shipping` - (Optional) Map(String). Mailing and shipping address for the customer. Appears on invoices emailed to this customer. For all individual fields see: [Shipping Fields](#shipping-fields).
* `balance` - (Optional) Int. Current balance, if any, being stored on the customer. If negative, the customer has credit to apply to their next invoice. If positive, the customer has an amount owed that will be added to their next invoice. The balance does not refer to any unpaid invoices; it solely takes into account amounts that have yet to be successfully applied to any invoice. This balance is only taken into account as invoices are finalized.
* `invoice_prefix` - (Optional) String. The prefix for the customer used to generate unique invoice numbers. Must be `3–12 uppercase letters or numbers`.
* `invoice_settings` - (Optional) Map(String). Default invoice settings for this customer. For supported fields see: [Invoice Settings Fields](#invoice-settings-fields).
* `next_invoice_sequence` - (Optional) Int. The sequence to be used on the customer’s next invoice. Defaults to 1.
* `preferred_locales` - (Optional) List(String). Customer’s preferred languages, ordered by preference.
* `metadata` - (Optional) Map(String). Set of key-value pairs that you can attach to an object. This can be useful for storing additional information about the object in a structured format.

### Address fields
* `line1` - (Optional) String. Address line 1 (e.g., street, PO Box, or company name).
* `line2` - (Optional) String. Address line 2 (e.g., apartment, suite, unit, or building).
* `city` - (Optional) String. City, district, suburb, town, or village.
* `postal_code` - (Optional) String. ZIP or postal code.
* `state` - (Optional) String. State, county, province, or region.
* `country` - (Optional) String. Two-letter country code (`ISO 3166-1 alpha-2`).

### Shipping fields
* `name` - (Optional) String. Customer name.
* `phone` - (Optional) String. Customer phone (including extension).
* `line1` - (Optional) String. Address line 1 (e.g., street, PO Box, or company name).
* `line2` - (Optional) String. Address line 2 (e.g., apartment, suite, unit, or building).
* `city` - (Optional) String. City, district, suburb, town, or village.
* `postal_code` - (Optional) String. ZIP or postal code.
* `state` - (Optional) String. State, county, province, or region.
* `country` - (Optional) String. Two-letter country code (`ISO 3166-1 alpha-2`).

### Invoice Settings Fields
* `default_payment_method` - (Optional) String. ID of a payment method that’s attached to the customer, to be used as the customer’s default payment method for subscriptions and invoices.
* `footer` - (Optional) String. Default footer to be displayed on invoices for this customer.
* `.` - (Optional) String. The `.` can be replaced by any string consequently it is considered as custom field name. 

## Attribute Reference

Attributes exported by this resource include:

* `name` - String. The customer’s full name or business name.
* `email` - String. Customer’s email address.
* `description` - String. An arbitrary string that you can attach to a customer object.
* `phone` - String. The customer’s phone number.
* `address` - Map(String). The customer’s address.
* `shipping` - Map(String). Mailing and shipping address for the customer.
* `balance` - Int. Current balance, if any, being stored on the customer. 
* `invoice_prefix` - String. The prefix for the customer used to generate unique invoice numbers.
* `default_invoice_prefix` - String. The default invoice prefix generated by Stripe when not individual invoice prefix provided.
* `invoice_settings` - Map(String). Default invoice settings for this customer.
* `next_invoice_sequence` - Int. The sequence to be used on the customer’s next invoice.
* `preferred_locales` - List(String). Customer’s preferred languages.
* `metadata` - Map(String). Set of key-value pairs that you can attach to an object.