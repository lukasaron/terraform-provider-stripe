---
layout: "stripe"
page_title: "Stripe: stripe_webhook_endpoint"
description: |-
The Stripe Webhook Endpoint can be created, modified, configured and removed by this resource.
---

# stripe_webhook_endpoint

With this resource, you can create webhook endpoint - [Stripe API webhook endpoint documentation](https://stripe.com/docs/api/webhook_endpoints). 

## Example Usage

Arguments accepted by this resource include:

```hcl
resource "stripe_webhook_endpoint" "webhook" {
  url            = "https://webhook-url-consumer.com"
  description    = "this is the webhook endpoint for subscriptions"
  enabled_events = [
    "customer.subscription.created", 
    "customer.subscription.updated"
  ]
}
```

## Argument Reference
* `url` - (Required) String. The URL where all webhook events will be sent.
* `enabled_events` - (Required) List(String). Events which will be triggered from the Stripe. Supported: [Stripe event types](https://stripe.com/docs/api/events/types).
* `description` - (Optional) String. Description for this webhook endpoint.
* `metadata` - (Optional) Map(String). Set of key-value pairs that you can attach to an object.
* `disabled` - (Optional) Bool. Disable the webhook endpoint if set to true. Can be used only for modification already existing webhook endpoint.

## Attribute Reference

Attributes exported by this resource include:

* `id` - String. The unique identifier for the object.
* `enabled_events` - List(String). The list of events to enable for this endpoint.
* `url` - String. The URL of the webhook endpoint.
* `description` - String. An optional description of what the webhook is used for.
* `disabled` - Bool. Informs whether the webhook endpoint is disabled.
* `metadata` - Map(String). Set of key-value pairs attached to an object.
* `secret` - String. The endpoint’s secret, used to generate webhook signatures. This field is marked as `sensitive`.