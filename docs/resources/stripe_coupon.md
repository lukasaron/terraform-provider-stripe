---
layout: "stripe"
page_title: "Stripe: stripe_coupon"
description: |-
The Stripe Coupon can be created, modified, configured and removed by this resource.
---

# stripe_coupon

With this resource, you can create a coupon - [Stripe API coupon documentation](https://stripe.com/docs/api/coupons).

A coupon contains information about a percent-off or amount-off discount you might want to apply to a customer.

A coupon has either a `percent_off` or an `amount_off` and `currency`. If you set an `amount_off`, that amount will be subtracted from any invoice’s subtotal. 

For example, an invoice with a subtotal of $100 will have a final total of $0 if a coupon with an amount_off of 20000 is applied to it and an invoice with a subtotal of $300 will have a final total of $100 if a coupon with an amount_off of 20000 is applied to it.

## Example Usage

```hcl
// coupon for the amount off discount
resource "stripe_coupon" "coupon" {
  name            = "$10 amount off"
  amount_off      = 1000
  currency        = "aud"
  duration        = "once"
  max_redemptions = 10
}

// coupon for the percentage off discount
resource "stripe_coupon" "coupon" {
  name            = "33.3% discount"
  percentage_off  = 33.3
  duration        = "forever"
}

// coupon with limitation to a date and the product only
resource "stripe_coupon" "coupon" {
  name       = "applies to prod with ID 123 till a date"
  amount_off = 2000
  duration   = "once"
  redeem_by  = "2025-07-23T03:27:06+00:00"
  // the stripe_product.product has to be created separately
  applies_to = [stripe_product.product.id] 
}
```

## Argument Reference

Arguments accepted by this resource include:

* `coupon_id` - (Optional) String. Unique string of your choice that will be used to identify this coupon when applying it to a customer.
* `name` - (Optional) String. Name of the coupon displayed to customers on for instance invoices or receipts.
* `amount_off` - (Optional) Int. Amount (in the currency specified) that will be taken off the subtotal of any invoices for this customer.
* `currency` - (Optional) String. Required if `amount_off` has been set, the three-letter ISO code for the currency of the amount to take off.
* `percent_off` - (Optional) Float. Percent that will be taken off the subtotal of any invoices for this customer for the duration of the coupon. For example, a coupon with percent_off of 50 will make a $100 invoice $50 instead.
* `duration` - (Optional) String. Describes how long a customer who applies this coupon will get the discount. One of `forever`, `once`, and `repeating`.
* `max_redemptions` - (Optional) Int. Maximum number of times this coupon can be redeemed, in total, across all customers, before it is no longer valid.
* `redeem_by` - (Optional) String. Date after which the coupon can no longer be redeemed. Expected format is in the `RFC3339`.
* `applies_to` - (Optional) List(String). A list of product IDs this coupon applies to.
* `metadata` - (Optional) Map(String). Set of key-value pairs that you can attach to an object. This can be useful for storing additional information about the object in a structured format.

## Attribute Reference

Attributes exported by this resource include:

* `id` - String. The unique identifier for the object.
* `coupon_id` - String. The unique identifier for the object.
* `name` - String. Name of the coupon displayed to customers on for instance invoices or receipts.
* `amount_off` - Int. Amount (in the currency specified) that will be taken off the subtotal of any invoices for this customer.
* `currency` - String. The three-letter ISO code for the currency of the amount to take off.
* `percent_off` - Float. Percent that will be taken off the subtotal of any invoices for this customer for the duration of the coupon.
* `duration` - String. Describes how long a customer who applies this coupon will get the discount.
* `max_redemptions` - Int. Maximum number of times this coupon can be redeemed.
* `redeem_by` - String. Date after which the coupon can no longer be redeemed in the `RFC3339` format.
* `times_redeemed` - Int. Number of times this coupon has been applied to a customer.
* `applies_to` - List(String). A list of product IDs this coupon applies to.
* `valid` - Bool. Taking account of the above properties, whether this coupon can still be applied to a customer.
* `metadata` - Map(String). Set of key-value pairs attached to an object.