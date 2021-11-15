---
layout: "stripe"
page_title: "Stripe: stripe_price"
description: |- The Stripe Price can be created, modified and configured by this resource.
---

# stripe_price

With this resource, you can create a price - [Stripe API price documentation](https://stripe.com/docs/api/prices).

Prices define the unit cost, currency, and (optional) billing cycle for both recurring and one-time purchases of
products. Products help you track inventory or provisioning, and prices help you track payment terms.

Different physical goods or levels of service should be represented by products, and pricing options should be
represented by prices. This approach lets you change prices without having to change your provisioning scheme.

For example, you might have a single "gold" product that has prices for $10/month, $100/year, and €9 once.

~> Removal of the price isn't supported through the Stripe SDK.

## Example Usage

```hcl
// basic price for the product
resource "stripe_price" "price" {
  // product needs to be defined
  product     = stripe_product.product.id
  currency    = "aud"
  unit_amount = 100
}

// basic free price for the product
resource "stripe_price" "price" {
  // product needs to be defined
  product     = stripe_product.product.id
  currency    = "aud"
  unit_amount = -1
}

// recurring price for the product
resource "stripe_price" "price" {
  // product needs to be defined
  product        = stripe_product.product.id
  currency       = "aud"
  billing_scheme = "per_unit"
  unit_amount    = "100"

  recurring {
    interval       = "week"
    interval_count = 1
  }
}

// tiered price for the product
resource "stripe_price" "price" {
  // product needs to be defined
  product        = stripe_product.product.id
  currency       = "aud"
  billing_scheme = "tiered"
  tiers_mode     = "graduated"

  # free up to ten
  tiers {
    up_to       = 10
    unit_amount = 0
  }
  
  tiers {
    up_to       = 100
    unit_amount = 300
  }
  
  tiers {
    up_to       = -1
    unit_amount = 100.5
  }

  recurring {
    interval        = "week"
    aggregate_usage = "sum"
    interval_count  = 2
    usage_type      = "metered"
  }
}

```

## Argument Reference

Arguments accepted by this resource include:

* `currency` - (Required) String. Three-letter ISO currency code, in lowercase.
* `product` - (Required) String. The ID of the product that this price will belong to.
* `unit_amount` - (Required unless `billing_scheme = tiered`) Int. A positive integer in cents (or `-1` for a free
  price) representing how much to charge.
* `unit_amount_decimal` - (Optional) Float. Same as `unit_amount`, but accepts a decimal value in cents with at most 12
  decimal places. Only one of `unit_amount` and `unit_amount_decimal` can be set.
* `active` - (Optional) Bool. Whether the price can be used for new purchases. Defaults to `true`.
* `nickname` - (Optional) String. A brief description of the price, hidden from customers.
* `recurring` - (Optional) List(Resource). The recurring components of a price such as `interval` and `usage_type`. For
  details of individual arguments see [Recurring](#recurring).
* `tiers` - (Optional) List(Resource). Each element represents a pricing tier. This parameter requires `billing_scheme`
  to be set to `tiered`. See also the documentation for `billing_scheme`. For details of individual arguments
  see [Tiers](#tiers).
* `tiers_mode` - (Required if `billing_scheme = tiered`) String. Defines if the tiering price should be `graduated`
  or `volume` based. In `volume`-based tiering, the maximum quantity within a period determines the per-unit price,
  in `graduated` tiering pricing can successively change as the quantity grows.
* `billing_scheme` - (Optional) String. Describes how to compute the price per period. Either `per_unit` or `tiered`
  . `per_unit` indicates that the fixed amount (specified in `unit_amount` or `unit_amount_decimal`) will be charged per
  unit in quantity (for prices with `usage_type=licensed`), or per unit of total usage (for prices
  with `usage_type=metered`). `tiered` indicates that the unit pricing will be computed using a tiering strategy as
  defined using the `tiers` and `tiers_mode` attributes.
* `lookup_key` - (Optional) String. A lookup key used to retrieve prices dynamically from a static string.
* `transfer_lookup_key` - (Optional) Bool. If set to `true`, will atomically remove the lookup key from the existing
  price, and assign it to this price.
* `tax_behaviour` - (Optional) String. Specifies whether the price is considered inclusive of taxes or exclusive of
  taxes. One of `inclusive`, `exclusive`, or `unspecified`. Once specified as either `inclusive` or `exclusive`, it
  cannot be changed.
* `transform_quantity` - (Optional) List(Resource). Apply a transformation to the reported usage or set quantity before
  computing the billed price. Cannot be combined with `tiers`. For details of individual arguments
  see [Transform Quantity](#transform-quantity).
* `metadata` - (Optional) Map(String). Set of key-value pairs that you can attach to an object. This can be useful for
  storing additional information about the object in a structured format.

### Recurring

`recurring` Supports the following arguments:

* `interval` - (Required) String. Specifies billing frequency. Either `day`, `week`, `month` or `year`.
* `aggregate_usage` - (Optional) String. Specifies a usage of aggregation strategy for prices of `usage_type=metered`.
  Allowed values are `sum` for summing up all usage during a period, `last_during_period` for using the last usage
  record reported within a period, `last_ever` for using the last usage record ever (across period bounds) or `max`
  which uses the usage record with the maximum reported usage during a period.
* `interval_count` - (Optional) Int. The number of intervals between subscription billings. For
  example, `interval=month` and `interval_count=3` bills every 3 months. Maximum of one year interval allowed (1 year,
  12 months, or 52 weeks).
* `usage_type` - (Optional) String. Configures how the quantity per period should be determined. Can be either `metered`
  or `licensed`. `licensed` automatically bills the quantity set when adding it to a subscription. `metered` aggregates
  the total usage based on usage records. Defaults to `licensed`.

### Tiers

`tiers` Can be used multiple times within the Price resource and supports the following arguments:

* `up_to` - (Required) Int. Specifies the upper bound of this tier. The lower bound of a tier is the upper bound of the
  previous tier adding one. Use `-1` to define a fallback tier.
* `flat_amount` - (Optional) Int. The flat billing amount for an entire tier, regardless of the number of units in the
  tier.
* `flat_amount_decimal` - (Optional) Float. Same as `flat_amount`, but accepts a decimal value representing an integer
  in the minor units of the currency. Only one of `flat_amount` and `flat_amount_decimal` can be set.
* `unit_amount` - (Optional) Int. The per-unit billing amount for each individual unit for which this tier applies.
* `unit_amount_decimal` - (Optional) Float. Same as `unit_amount`, but accepts a decimal value in cents with at most 12
  decimal places. Only one of `unit_amount` and `unit_amount_decimal` can be set.

### Transform Quantity

`transform_quantity` Supports the following arguments:

* `divide_by` - (Required) Int. Divide usage by this number.
* `round` - (Required) String. After division, either round the result `up` or `down`.

## Attribute Reference

Attributes exported by this resource include:

* `id` - String. The unique identifier for the object.
* `currency` - String. Three-letter ISO currency code.
* `product` - String. The ID of the product that this price will belong to.
* `unit_amount` - Int. A positive integer in cents (or `-1` for a free price) representing how much to charge.
* `unit_amount_decimal` - Float. Same as `unit_amount`, but accepts a decimal value in cents with at most 12 decimal
  places.
* `active` - Bool. Whether the price can be used for new purchases. Defaults to `true`.
* `nickname` - String. A brief description of the price, hidden from customers.
* `recurring` - List(Resource). The recurring components of a price such as `interval` and `usage_type`.
* `tiers` - List(Resource). Each element represents a pricing tier.
* `tiers_mode` - String. Defines if the tiering price should be `graduated` or `volume` based.
* `billing_scheme` - String. Describes how to compute the price per period.
* `lookup_key` - String. A lookup key used to retrieve prices dynamically from a static string.
* `transfer_lookup_key` - Bool. `true` when lookup key was transferred.
* `tax_behaviour` - String. Specifies whether the price is considered inclusive of taxes or exclusive of taxes.
* `transform_quantity` - List(Resource). Apply a transformation to the reported usage or set quantity before computing
  the billed price.
* `type` - String. One of `one_time` or `recurring` depending on whether the price is for a one-time purchase or a
  recurring (subscription) purchase.
* `metadata` - Map(String). Set of key-value pairs that you can attach to an object.