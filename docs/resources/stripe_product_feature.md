---
layout: "stripe"
page_title: "Stripe: stripe_entitlements_feature"
description: |- The Stripe Product Feature can be created, configured and removed by this resource.
---

# stripe_product_feature

With this resource, you can create Product Feature - [Stripe API product feature documentation](https://docs.stripe.com/api/product-feature) 

A Product Feature represents an attachment between a feature and a product. 
When a product is purchased that has a feature attached, 
Stripe will create an entitlement to the feature for the purchasing customer.

## Example Usage

```hcl
// Product Feature attaches an Entitlement Feature to a Product
resource "stripe_product_feature" "product_feature" {
  // stripe_entitlements_feature.feature has to be created separately.
  entitlements_feature = stripe_entitlements_feature.feature.id
  // stripe.product.product has to be created separately.
  product              = stripe_product.product.id
}
```

## Argument Reference

Arguments accepted by this resource include:

* `entitlements_feature` - (Required) String. The ID of the Entitlements Feature the product will be attached to
* `product` - (Required) String. The ID of the product that this Entitlements Feature will be attached to.

## Attribute Reference

Attributes exported by this resource include:

* `id` - String. The unique identifier for the object.
* `entitlements_feature` - String. The ID of the Entitlements Feature.
* `product` - String. The ID of the product.
* `livemode` - Bool. Has the value `true` if the object exists in live mode or the value `false`
  if the object exists in test mode.