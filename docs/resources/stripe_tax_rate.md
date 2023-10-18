---
layout: "stripe"
page_title: "Stripe: stripe_tax_rate"
description: |- The Stripe Tax Rate can be created, modified and configured by this resource.
---

# stripe_tax_rate

With this resource, you can create a tax rate - [Stripe API tax rate  documentation](https://stripe.com/docs/api/tax_rates).

Tax rates can be applied to invoices, subscriptions and Checkout Sessions to collect tax.

## Example Usage

```hcl
resource "stripe_tax_rate" "tax_rate" {

  # Required Fields
  display_name            = "GST" 
  inclusive               = true  
  percentage              = 10.0
  active                  = true
  
  # Optional fields
  country                 = "AU"
  description             = "GST Australia"
  jurisdiction            = "AU"
  state                   = ""
  tax_type                = ""
  metadata                = {}
}
```

## Argument Reference

Arguments accepted by this resource include:

* `display_name` - (Required) String. The display name of the tax rate, which will be shown to users.
* `inclusive` - (Required) Bool. This specifies if the tax rate is inclusive or exclusive.
* `percentage ` - (Required) Float. This represents the tax rate percent out of 100.
* `active` - (Optional) Bool. Flag determining whether the tax rate is active or inactive (archived). Inactive tax rates cannot be used with new applications or Checkout Sessions, but will still work for subscriptions and invoices that already have it set.
* `country` - (Optional) String. Two-letter country code (ISO 3166-1 alpha-2).
* `description` - (Optional) String. An arbitrary string attached to the tax rate for your internal use only. It will not be visible to your customers.
* `jurisdiction` - (Optional) String. The jurisdiction for the tax rate. You can use this label field for tax reporting purposes. It also appears on your customer’s invoice.
* `metadata` - (Optional) Map(String). Set of key-value pairs that you can attach to an object. This can be useful for storing additional information about the object in a structured format. Individual keys can be unset by posting an empty value to them. All keys can be unset by posting an empty value to metadata.
* `state` - (Optional) String. ISO 3166-2 subdivision code, without country prefix. For example, “NY” for New York, United States.
* `tax_type` - (Optional) String. The high-level tax type, such as vat or sales_tax.


## Attribute Reference

Attributes exported by this resource include:

* `id` - String. The unique identifier for the object.
* `display_name` - String. The display name of the tax rate, which will be shown to users.
* `inclusive` - Bool. This specifies if the tax rate is inclusive or exclusive.
* `percentage ` - Float. This represents the tax rate percent out of 100.
* `active` - Bool. Flag determining whether the tax rate is active or inactive (archived). Inactive tax rates cannot be used with new applications or Checkout Sessions, but will still work for subscriptions and invoices that already have it set.
* `country` - String. Two-letter country code (ISO 3166-1 alpha-2).
* `description` - String. An arbitrary string attached to the tax rate for your internal use only. It will not be visible to your customers.
* `jurisdiction` - String. The jurisdiction for the tax rate. You can use this label field for tax reporting purposes. It also appears on your customer’s invoice.
* `metadata` - Map(String). Set of key-value pairs that you can attach to an object. This can be useful for storing additional information about the object in a structured format. Individual keys can be unset by posting an empty value to them. All keys can be unset by posting an empty value to metadata.
* `state` - String. ISO 3166-2 subdivision code, without country prefix. For example, “NY” for New York, United States.
* `tax_type` - String. The high-level tax type, such as vat or sales_tax.
* `object` - String. String representing the object’s type. Objects of the same type share the same value.
* `created` - Int. Time at which the object was created. Measured in seconds since the Unix epoch.
* `livemode` - Bool. Has the value true if the object exists in live mode or the value false if the object exists in test mode.
