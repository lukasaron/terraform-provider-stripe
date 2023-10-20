---
layout: "stripe"
page_title: "Stripe: stripe_shipping_rate"
description: |- The Stripe Shipping Rate can be created, modified and configured by this resource.
---

# stripe_shipping_rate

With this resource, you can create a shipping rate - [Stripe API price documentation](https://stripe.com/docs/api/shipping_rates).


Shipping rates let you display various shipping options—like standard, express, and overnight—with more accurate delivery estimates. 
Charge your customer for shipping using different Stripe products, some of which require coding.

~> Removal of the shipping rate isn't supported through the Stripe SDK. The best practice, which this provider follows,
is to archive the shipping rate by marking it as inactive on destroy, which indicates that the shipping rate is no longer
available.

## Example Usage

```hcl
// minimal shipping rate
resource "stripe_shipping_rate" "shipping_rate" {
  display_name = "minimal shipping rate"
  fixed_amount {
    amount   = 1000
    currency = "aud"
  }
}

// shipping rate with delivery estimate
resource "stripe_shipping_rate" "shipping_rate" {
  display_name = "shipping rate"
  fixed_amount {
    amount   = 1000
    currency = "aud"
  }

  delivery_estimate {
    minimum {
      unit = "hour"
      value = 24
    }
    maximum {
      unit = "day"
      value = 4
    }
  }
}

// shipping rate with currency options
// !!! Currency options have to be sorted alphabetically by the currency field
resource "stripe_shipping_rate" "shipping" {
  display_name = "shipping rate"
  fixed_amount {
    amount   = 1000
    currency = "aud"
    
    currency_option {
      currency = "eur"
      amount = 350
    }
    currency_option {
      currency = "usd"
      amount = 500
    }
  }
}

```

## Argument Reference

Arguments accepted by this resource include:

* `type` - (Optional) String. The type of calculation to use on the shipping rate. Can only be `fixed_amount` for now.
* `display_name` - (Required) String. The name of the shipping rate, meant to be displayable to the customer. 
  This will appear on CheckoutSessions.
* `fixed_amount` - (Required) List(Resource). Describes a fixed amount to charge for shipping. 
   Must be present if type is `fixed_amount`. For details of individual arguments see [Fixed Amount](#fixed-amount).
* `delivery_estimate` - (Optional) List(Resource). The estimated range for how long shipping will take, 
   meant to be displayable to the customer. This will appear on CheckoutSessions. 
   For details please see [Delivery Estimate](#delivery-estimate).
* `active` - (Optional) Bool. Whether the shipping rate is active (can't be used when creating). Defaults to `true`.
* `tax_behaviour` - (Optional) String. Specifies whether the price is considered inclusive of taxes or exclusive of
  taxes. One of `inclusive`, `exclusive`, or `unspecified`. Once specified it cannot be changed, default is `unspecified`.
* `tax_code` - (Optional) String. A tax code ID. The Shipping tax code is `txcd_92010001`.
* `metadata` - (Optional) Map(String). Set of key-value pairs that you can attach to an object. This can be useful for
  storing additional information about the object in a structured format.

### Fixed Amount

`fixed_amount` Supports the following arguments:

* `currency` - (Required) String. Three-letter ISO currency code, in lowercase - [supported currencies](https://stripe.com/docs/currencies).
* `amount` - (Required) Int. A non-negative integer in cents representing how much to charge.
* `currency_option` - (Optional) List(Resource). Please see argument details [Currency Option](#currency-option)

### Currency Option

`currency_option` Can be used multiple times within the `fixed_amount` part.  
~> When multiple currency_options are defined sorting by currency field is mandatory! 
Otherwise, the provider consider next run as a change.

Currency option support the following arguments:

* `currency` - (Required) String. Three-letter ISO currency code, in lowercase - [supported currencies](https://stripe.com/docs/currencies).
* `amount` - (Required) Int. (Required) Int. A non-negative integer in cents representing how much to charge.
* `tax_behaviour` - (Optional) String. Specifies whether the price is considered inclusive of taxes or exclusive of
  taxes. One of `inclusive`, `exclusive`, or `unspecified`. Once specified it cannot be changed, default is `unspecified`.

### Delivery Estimate

`delivery_estimate` Supports the following arguments:

* `minimum` - (Required) List(Resource). The lower bound of the estimated range. 
  Please see [Delivery Estimate Definition](#delivery-estimate-definition).
* `maximum` - (Required) List(Resource. The upper bound of the estimated range.
  Please see [Delivery Estimate Definition](#delivery-estimate-definition).

### Delivery Estimate Definition

`maximum` and `minimum` share the same definition:

* `unit` - (Required) String. A unit of time. Possible values `hour`, `day`, `business_day`, `week` and `month`.
* `value` - (Required) Int. Must be greater than 0.

## Attribute Reference

Attributes exported by this resource include:

* `id` - String. The unique identifier for the object.
* `type` - String. The type of calculation to use on the shipping rate.
* `display_name` - String. The name of the shipping rate, meant to be displayable to the customer.
* `active` - Bool. Whether the shipping rate can be used.
* `fixed_amount` - List(Resource). Describes a fixed amount to charge for shipping.
* `delivery_estimate` - List(Resource). The estimated range for how long shipping will take, meant to be displayable to the customer.
* `tax_behaviour` - String. Specifies whether the price is considered inclusive of taxes or exclusive of taxes. 
* `tax_code` - String. A tax code ID.
* `livemode` - Bool. Has the value true if the object exists in live mode or the value false if the object exists in test mode.