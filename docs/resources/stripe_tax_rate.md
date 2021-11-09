---
layout: "stripe"
page_title: "Stripe: stripe_tax_rate"
description: |-
The Stripe Tax Rate can be created, modified, configured by this resource. Deletion is not supported.
---

# stripe_tax_rate

With this resource, you can create a product - [Stripe API product documentation](https://stripe.com/docs/api/tax_rates).

Tax rates can be applied to invoices, subscriptions and Checkout Sessions to collect tax.

## Example Usage

```hcl
// the most basic tax rate
resource "stripe_tax_rate" "tax_rate" {
  display_name = "minimalist product"
  inclusive = false
  percentage = 10
}

// the full parameter tax rate
resource "stripe_tax_rate" "tax_rate" {
  active = true
  display_name = "minimalist product"
  description = ""
  inclusive = false
  percentage = 10
  jurisdiction = "JP"
  metadata = {
    key = "value"
  }
}

```

## Argument Reference

Arguments accepted by this resource include:

* `active` - (Optional) Bool. Defaults to `true`. When set to false, this tax rate cannot be used with new applications or Checkout Sessions, but will still work for subscriptions and invoices that already have it set.
* `description` - (Optional) String. An arbitrary string attached to the tax rate for your internal use only. It will not be visible to your customers.
* `display_name` - (Required) String. The display name of the tax rate, which will be shown to users.
* `inclusive` - (Required) Bool. This specifies if the tax rate is inclusive or exclusive.
* `jurisdiction` - (Optional) String. The jurisdiction for the tax rate. You can use this label field for tax reporting purposes. It also appears on your customer’s invoice.
* `metadata` - (Optional) Map(String). Set of key-value pairs that you can attach to an object. This can be useful for storing additional information about the object in a structured format.
* `percentage` - (Optional) Float. This represents the tax rate percent out of 100.

## Attribute Reference

Attributes exported by this resource include:


* `id` - String. The unique identifier for the object.
* `active` - Bool. Whether the tax-rate is currently available for set. 
* `created` - Int.  Time at which the object was created. Measured in seconds since the Unix epoch.
* `description` - String. The tax-rate's description, meant to be displayable to the customer.
* `display_name` - String. The display name of the tax rate, which will be shown to users.
* `inclusive` - Bool. Whether the tax rate is inclusive or exclusive.
* `jurisdiction` - String. The jurisdiction for the tax rate.
* `livemode` - Bool. Whether the tax-rate is currently exists in live mode. 
* `metadata` - Map(String). Set of key-value pairs that you can attach to an object. This can be useful for storing additional information about the object in a structured format.
* `percentage` - Float. This represents the tax rate percent out of 100.
