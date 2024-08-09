---
layout: "stripe"
page_title: "Stripe: stripe_entitlements_feature"
description: |- The Stripe Entitlements Feature can be created, modified, and configured by this resource.
---

# stripe_entitlements_feature

With this resource, you can create Entitlements Feature for your products - [Stripe API entitlements feature documentation](https://docs.stripe.com/api/entitlements/feature) 

Entitlements Features can be assigned to products, and when those products are purchased, 
Stripe will create an entitlement to the feature for the purchasing customer.

~> Removal of the Entitlements Feature isn't supported through the Stripe SDK. Consequently, deactivation is applied instead.

## Example Usage

```hcl
// Entitlements Feature
resource "stripe_entitlements_feature" "feature" {
  name       = "feature"
  lookup_key = "key"
}
```

## Argument Reference

Arguments accepted by this resource include:

* `name` - (Required) String. The feature’s name, for your own purpose, not meant to be displayable to the customer.
* `lookup_key` - (Required) String. A unique key you provide as your own system identifier. This may be up to 80 characters.
* `metadata` - (Optional) Map(String). Set of key-value pairs that you can attach to an object.
    This can be useful for storing additional information about the object in a structured format.

## Attribute Reference

Attributes exported by this resource include:

* `id` - String. The unique identifier for the object.
* `name` - String. The feature’s name.
* `lookup_key` - String. A unique key you provide as your own system identifier.
* `metadata` - Map(String). Set of key-value pairs that you can attach to an object.
* `active` - Inactive features cannot be attached to new products.
* `livemode` - Bool. Has the value `true` if the object exists in live mode or the value `false`
  if the object exists in test mode.