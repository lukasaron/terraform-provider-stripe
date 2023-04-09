---
layout: "stripe"
page_title: "Stripe: stripe_portal_configuration"
description: |- The Stripe Customer Portal Configuration can be created and modified by this resource.
---

# stripe_portal_configuration

With this resource, you can create a Customer Portal Configuration - [Stripe API portal configuration documentation](https://stripe.com/docs/api/customer_portal/configuration).

The Billing customer portal is a Stripe-hosted UI for subscription and billing management.

A portal configuration describes the functionality and features that you want to provide to your customers through the portal.

~> Removal of the Customer Portal isn't supported through the Stripe SDK. The best practice, which this provider follows,
is to deactivate the Customer Portal by marking it as inactive on destroy, which indicates that resource is no longer
available.

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
* `business_profile` - (Required) List(Resource). The business information shown to customers in the portal. More details in [Business Profile section](#business-profile)
* `default_return_url` - (Optional) String. The default URL to redirect customers to when they click on the portal’s link to return to your website. This can be overriden when creating the session.
* `login_page` - (Optional) List(Resource). The hosted login page for this configuration. See details in [Login Page Section](#login-page).
* `features` - (Required) List(Resource). Information about the features available in the portal. Feature section described in [Feature section](#features)
* `metadata` - (Optional) Map(String). Set of key-value pairs that you can attach to an object. This can be useful for storing additional information about the object in a structured format.

### Business Profile

`business_profile` Supports the following arguments:

* `headline` - (Optional) String. The messaging shown to customers in the portal.
* `privacy_policy_url` - (Optional) String. A link to the business's publicly available privacy policy.
* `terms_of_service_url` - (Optional) String. A link to the business's publicly available terms of service.

### Login Page

`login_page` Includes only one option:

* `enabled` - (Required) Bool. Set to true to generate a shareable URL login_page.url that will take your customers to a hosted login page for the customer portal.

### Features

`features` Supports the following sections:

* `customer_update` - (Optional) List(Resource). Information about updating the customer details in the portal. See [Customer Update](#features-customer-update).
* `invoice_history` - (Optional) List(Resource). Information about showing the billing history in the portal. See [Invoice History](#features-invoice-history).
* `payment_method_update` - (Optional) List(Resource). Information about updating payment methods in the portal. See [Payment Method Update](#features-payment-method-update).
* `subscription_cancel` - (Optional) List(Resource). Information about canceling subscriptions in the portal. See [Subscription Cancel](#features-subscription-cancel).
* `subscription_pause`- (Optional) List(Resource). Information about pausing subscriptions in the portal. See [Subscription Pause](#features-subscription-pause).
* `subscription_update`- (Optional) List(Resource). Information about updating subscriptions in the portal. See [Subscription Update](#features-subscription-update).

### Features Customer Update

`customer_update` Consists of:

* `enabled` - (Required) Bool. Whether the feature is enabled.
* `allowed_updates` - (Optional) List(String). The types of customer updates that are supported [`name`, `email`, `address`, `shipping`, `phone`, `tax_id`]. When empty, customers are not updatable.

### Features Invoice History

`invoice_history` Includes only one option:

* `enabled` - (Required) Bool. Whether the feature is enabled.

### Features Payment Method Update

`payment_method_update` Consists of only one option:

* `enabled` - (Required) Bool. Whether the feature is enabled.

### Features Subscription Cancel

`subscription_cancel` Supports these arguments:

* `enabled` - (Required) Bool. Whether the feature is enabled.
* `mode` - (Optional) String. Whether to cancel subscriptions immediately or at the end of the billing period. Valid value is either `immediately` or `at_period_end`
* `proration_behavior` - (Optional) String. Whether to create prorations when canceling subscriptions. Possible values are `none` and `create_prorations`, which is only compatible with `mode=immediately`. No prorations are generated when canceling a subscription at the end of its natural billing period.
* `cancellation_reason` - (Optional) List(Resource). Whether the cancellation reasons will be collected in the portal and which options are exposed to the customer. Details of this field is in [Cancellation Reason](#features-subscription-cancel-cancellation-reason).

#### Features Subscription Cancel Cancellation Reason

`cancellation_reason` consumes the following arguments:

* `enabled` - (Required) Bool. Whether the feature is enabled.
* `options` - (Required) List(String). Which cancellation reasons will be given as options to the customer. Supported values are `too_expensive`, `missing_features`, `switched_service`, `unused`, `customer_service`, `too_complex`, `low_quality`, and `other`.


### Features Subscription Pause

`subscription_pause` Implements only one argument:

* `enabled` - (Required) Bool. Whether the feature is enabled.

### Features Subscription Update

`subscription_update` Consists of these arguments:

* `enabled` - (Required) Bool. Whether the feature is enabled.
* `default_allowed_updates` - (Required) List(String). The types of subscription updates that are supported. When empty, subscriptions are not updatable. Supported values are `price`, `quantity`, and `promotion_code`.
* `products` - (Required) List(Resource). The list of products that support subscription updates. See details [Products](#features-subscription-update-products).
* `proration_behavior` - (Optional) String. Determines how to handle prorations resulting from subscription updates. Valid values are `none`, `create_prorations`, and `always_invoice`.

#### Features Subscription Update Products

`products` has to be defined with following fields:

* `product` - (Required) String. The product id.
* `prices` - (Required) List(String). The list of price IDs for the product that a subscription can be updated to.


## Attribute Reference

Attributes exported by this resource include:

* `id` - String. Unique identifier for the object.
* `active` - Bool. Whether the configuration is active and can be used to create portal sessions.
* `business_profile` - Map(String). The business information shown to customers in the portal.
* `default_return_url` - String. The default URL to redirect customers to when they click on the portal’s link.
* `features` - Map(String). Information about the features available in the portal.
* `metadata` - Map(String). Set of key-value pairs that you can attach to an object.
