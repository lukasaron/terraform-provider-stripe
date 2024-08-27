---
layout: "stripe"
page_title: "Stripe: stripe_promotion_code"
description: |- 
  The Stripe Promotion Code can be created, modified and configured by this resource.
---

# stripe_promotion_code

With this resource, you can create a promotion code - [Stripe API promotion code documentation](https://stripe.com/docs/api/promotion_codes).

A Promotion Code represents a customer-redeemable code for a coupon. It can be used to create multiple codes for a single coupon.

~> Removal of the promotion code isn't supported through the Stripe SDK.

## Example Usage

```hcl
// promotion code for the coupon
resource "stripe_promotion_code" "code" {
  // coupon needs to be defined
  coupon = stripe_coupon.coupon.id
  code   = "FREE"
}

// promotion code for the coupon with limitations
resource "stripe_promotion_code" "code" {
  // coupon needs to be defined
  coupon          = stripe_coupon.coupon.id
  code            = "FREE"
  max_redemptions = 5
  expires_at      = "2025-08-03T08:37:18+00:00"
}

// promotion code for the coupon to customer
resource "stripe_promotion_code" "code" {
  // coupon needs to be defined
  coupon   = stripe_coupon.coupon.id
  code     = "FREE"
  customer = "cus..."
}

// promotion code for the coupon with restrictions
resource "stripe_promotion_code" "code" {
  // coupon needs to be defined
  coupon = stripe_coupon.coupon.id
  code   = "FREE"
  
  restrictions {
    first_time_transaction  = true
    minimum_amount          = 100
    minimum_amount_currency = "aud"
  }
}
```

## Argument Reference

Arguments accepted by this resource include:

* `coupon` - (Required) String. The coupon for this promotion code.
* `code` - (Optional) String. The customer-facing code. Regardless of case, this code must be unique across all active promotion codes for a specific customer. If left blank, we will generate one automatically.
* `active` - (Optional) Bool. Whether the promotion code is currently active. Defaults to `true`.
* `customer` - (Optional) String. The customer that this promotion code can be used by. If not set, the promotion code can be used by all customers.
* `max_redemptions` - (Optional) Int. A positive integer specifying the number of times the promotion code can be redeemed. If the coupon has specified a `max_redemptions`, then this value cannot be greater than the coupon’s `max_redemptions`.
* `expires_at` - (Optional) String. The timestamp at which this promotion code will expire. If the coupon has specified a `redeems_by`, then this value cannot be after the coupon’s `redeems_by`. Expected format is `RFC3339`.
* `restrictions` - (Optional) List(Resource). Settings that restrict the redemption of the promotion code. For details of individual arguments see [Restrictions](#restrictions).   
* `metadata` - (Optional) Map(String). Set of key-value pairs that you can attach to an object. This can be useful for storing additional information about the object in a structured format.

### Restrictions

`restrictions` Supports the following arguments:

* `first_time_transaction` - (Required) Bool. A Boolean indicating if the Promotion Code should only be redeemed for Customers without any successful payments or invoices.
* `minimum_amount` - (Optional) Int. Minimum amount required to redeem this Promotion Code into a Coupon (e.g., a purchase must be $100 or more to work).
* `minimum_amount_currency` - (Optional) String. Three-letter ISO code for `minimum_amount`.

## Attribute Reference

Attributes exported by this resource include:

* `id` - String. The unique identifier for the object.
* `coupon` - String. The coupon for this promotion code.
* `code` - String. The customer-facing code. 
* `active` - Bool. Whether the promotion code is currently active.
* `customer` - String. The customer that this promotion code can be used by.
* `max_redemptions` - Int. A positive integer specifying the number of times the promotion code can be redeemed. 
* `expires_at` - String. The timestamp at which this promotion code will expire.
* `restrictions` - List. Settings that restrict the redemption of the promotion code - `first_time_transaction`, `minimum_amount` and `minimum_amount_currency`.
* `metadata` - Map(String). Set of key-value pairs that you can attach to an object. 

## Import

Import is supported using the following syntax:

```shell
$ terraform import stripe_promotion_code.code <promotion_code_id>
```