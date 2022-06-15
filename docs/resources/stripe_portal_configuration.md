---
layout: "stripe"
page_title: "Stripe: stripe_portal_configuration"
description: |-
The Stripe Customer Portal Configuration can be created and modified by this resource.
---

# stripe_portal_configuration

With this resource, you can create a Customer Portal Configuration - [Stripe API portal configuration documentation](https://stripe.com/docs/api/customer_portal/configuration).

The Billing customer portal is a Stripe-hosted UI for subscription and billing management.

A portal configuration describes the functionality and features that you want to provide to your customers through the portal.

## Example Usage

```hcl
// A basic portal with payment method update disabled
resource "stripe_portal_configuration" "portal_configuration" {
  business_profile {
    privacy_policy_url   = "https://example.com/privacy"
    terms_of_service_url = "https://example.com/terms"
  }
  features {
    payment_method_update {
      enabled = false
    }
  }
}

// Disabling an existing portal configuration (Update)
resource "stripe_portal_configuration" "portal_configuration" {
  active = false
  business_profile {
    privacy_policy_url   = "https://example.com/privacy"
    terms_of_service_url = "https://example.com/terms"
  }
  features {
    payment_method_update {
      enabled = false
    }
  }
}

// A billing portal configuration with custom headline and metadata
resource "stripe_portal_configuration" "portal_configuration" {
  business_profile {
    headline             = "My special headline"
    privacy_policy_url   = "https://example.com/privacy"
    terms_of_service_url = "https://example.com/terms"
  }
  default_return_url = "https://example.com/special_headline"
  features {
    invoice_history {
      enabled = false
    }
    payment_method_update {
      enabled = true
    }
  }
  metadata = {
    campaign = "special headline"
  }
}

// A billing portal using all the available options
resource "stripe_portal_configuration" "portal_configuration" {
  business_profile {
    headline             = "My special headline"
    privacy_policy_url   = "https://example.com/privacy"
    terms_of_service_url = "https://example.com/terms"
  }
  default_return_url = "https://example.com/special_headline"
  features {
    customer_update {
      enabled         = true
      allowed_updates = ["email", "address", "shipping", "phone", "tax_id"]
    }
    invoice_history {
      enabled = true
    }
    payment_method_update {
      enabled = true
    }
    subscription_cancel {
      enabled = true
      cancellation_reason {
        enabled = true
        options = ["too_expensive", "missing_features", "switched_service", "unused", "customer_service", "too_complex", "low_quality", "other"]
      }
      mode               = "at_period_end"
      proration_behavior = "none"
    }
    subscription_pause {
      enabled = true
    }
    subscription_update {
      enabled                 = true
      default_allowed_updates = ["price", "quantity", "promotion_code"]
      proration_behavior      = "none"
      products {
        product = "my_product_id"
        prices  = ["my_price_id1", "my_price_id2"]
      }
    }
  }
  metadata = {
    foo = "bar"
  }
}
```

## Argument Reference

Arguments accepted by this resource include:

* `active` - (Optional) Bool. Whether the configuration is active and can be used to create portal sessions. (On create it is always set as `true`)
* `business_profile` - (Required) Map(String). The business information shown to customers in the portal.
* `default_return_url` - (Optional) String. The default URL to redirect customers to when they click on the portal’s link to return to your website. This can be overriden when creating the session.
* `features` - (Required) Map(String). Information about the features available in the portal.
* `metadata` - (Optional) Map(String). Set of key-value pairs that you can attach to an object. This can be useful for storing additional information about the object in a structured format.

### Buesiness Profile fields

* `headline` - (Optional) String. The messaging shown to customers in the portal.
* `privacy_policy_url` - (Required) String. A link to the business's publicly available privacy policy.
* `terms_of_service_url` - (Required) String. A link to the business's publicly available terms of service.

### Features fields

**At least one of these fields must be added.**

Once you add a parent at least one sub field must be included.
For example, if adding `subscription_cancel`, then you MUST also declare `enabled` as its required for that parent field.

* `customer_update` - (Optional) Map(String). Information about updating the customer details in the portal.
  * `enabled` - (Required) Bool. Whether the feature is enabled.
  * `allowed_updates` - (Required) List. The types of customer updates that are permitted. When empty, customers are not updateable. Possible values are `email`, `address`, `shipping`, `phone` and `tax_id`.
* `invoice_history` - (Optional) Map(String). Information about showing the billing history in the portal.
  * `enabled` - (Required) Bool. Whether the feature is enabled.
* `payment_method_update` - (Optional) Map(String). Information about updating payment methods in the portal.
  * `enabled` - (Required) Bool. Whether the feature is enabled.
* `subscription_cancel` - (Optional) Map(String). Information about canceling subscriptions in the portal.
  * `enabled` - (Required) Bool. Whether the feature is enabled.
  * `cancellation_reason` - (Optional) Map(String). Whether the cancellation reasons will be collected in the portal and which options are exposed to the customer.
    * `enabled` - (Required) Bool. Whether the feature is enabled.
    * `options` - (Required) List. Which cancellation reasons will be given as options to the customer. Requires at least 2 values, from `too_expensive`, `missing_features`, `switched_service`, `unused`, `customer_service`, `too_complex`, `low_quality` and `other`.
  * `mode` - (Optional) String. Whether to cancel subscriptions immediately or at the end of the billing period. Possible values are `immediately` - Cancel subscriptions immediately, `at_period_end` - After canceling, customers can still renew subscriptions until the billing period ends.
  * `proration_behavior` - (Optional) String. Whether to create prorations when canceling subscriptions. Possible values are `none` and `create_prorations`, which is only compatible with `mode=immediately`. No prorations are generated when canceling a subscription at the end of its natural billing period.
* `subscription_pause`- (Optional) Map(String). Information about pausing subscriptions in the portal.
  * `enabled` - (Required) Bool. Whether the feature is enabled.
* `subscription_update`- (Optional) Map(String). Information about updating subscriptions in the portal.
  * `enabled` - (Required) Bool. Whether the feature is enabled.
  * `default_allowed_updates` - (Required) List. The types of subscription updates that are supported. When empty, subscriptions are not updateable. Possible values are `price`, `quantity` and `promotion_code`.
  * `products` - (Optional) Map(String). The list of products that support subscription updates.
    * `prices` - (Required) List. The list of price IDs for the product that a subscription can be updated to.
    * `product` - (Required) String. The product id.
  * `proration_behavior` - (Optional) String. Determines how to handle prorations resulting from subscription updates. Valid values are `none`, `create_prorations`, and `always_invoice`.

## Attribute Reference

Attributes exported by this resource include:

* `id` - String. Unique identifier for the object.
* `object` - String. String representing the object's type.
* `active` - Bool. Whether the configuration is active and can be used to create portal sessions.
* `application` - String. ID of the Stripe Connect Application that created the configuration.
* `business_profile` - Map(String). The business information shown to customers in the portal.
  * `headline` - String. The messaging shown to customers in the portal.
  * `privacy_policy_url` - String. A link to the business's publicly available privacy policy.
  * `terms_of_service_url` - String. A link to the business's publicly available terms of service.
* `created` - Int. Time at which the object was created. Measured in seconds since the Unix epoch.
* `default_return_url` - String. The default URL to redirect customers to when they click on the portal’s link to return to your website. This can be overriden when creating the session.
* `features` - Map(String). Information about the features available in the portal.
  * `customer_update` - Map(String). Information about updating the customer details in the portal.
    * `enabled` - Bool. Whether the feature is enabled.
    * `allowed_updates` - List. The types of customer updates that are supported. When empty, customers are not updateable. - `email`, `address`, `shipping`, `phone` and `tax_id`.
  * `invoice_history` - Map(String). Information about showing the billing history in the portal.
    * `enabled` - Bool. Whether the feature is enabled.
  * `payment_method_update` - Map(String). Information about updating payment methods in the portal.
    * `enabled` - Bool. Whether the feature is enabled.
  * `subscription_cancel` - Map(String). Information about canceling subscriptions in the portal.
    * `enabled` - Bool. Whether the feature is enabled.
    * `cancellation_reason` - Map(String). Whether the cancellation reasons will be collected in the portal and which options are exposed to the customer.
      * `enabled` - Bool. Whether the feature is enabled.
      * `options` - List. Which cancellation reasons will be given as options to the customer. - `too_expensive`, `missing_features`, `switched_service`, `unused`, `customer_service`, `too_complex`, `low_quality`, `other`.
    * `mode` - String. Whether to cancel subscriptions immediately or at the end of the billing period. - `immediately`, `at_period_end`.
    * `proration_behavior` - String. Whether to create prorations when canceling subscriptions. - `none`, `create_prorations`.
  * `subscription_pause`- Map(String). Information about pausing subscriptions in the portal.
    * `enabled` - Bool. Whether the feature is enabled.
  * `subscription_update`- Map(String). Information about updating subscriptions in the portal.
    * `default_allowed_updates` - List. The types of subscription updates that are supported. When empty, subscriptions are not updateable. - `price`, `quantity`, `promotion_code`.
    * `enabled` - Bool. Whether the feature is enabled.
    * `products` - Map(String). The list of products that support subscription updates.
      * `prices` - List. The list of price IDs for the product that a subscription can be updated to.
      * `product` - String. The product id.
    * `proration_behavior` - String. Determines how to handle prorations resulting from subscription updates. - `none`, `create_prorations`, `always_invoice`.
* `is_default` - Bool. Whether the configuration is the default. If `true`, this configuration can be managed in the Dashboard and portal sessions will use this configuration unless it is overriden when creating the session.
* `livemode` - Bool. Has the value true if the object exists in live mode or the value false if the object exists in test mode.
* `metadata` - Map(String). Set of key-value pairs that you can attach to an object.
* `updated` - Int. Time at which the object was last updated. Measured in seconds since the Unix epoch.
