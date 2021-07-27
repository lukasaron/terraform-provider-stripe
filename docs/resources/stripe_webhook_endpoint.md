---
layout: "stripe"
page_title: "Stripe: stripe_webhook_endpoint"
description: |-
The Stripe Webhook Endpoint can be created, modified, configured and removed by this resource.
---

# stripe_webhook_endpoint

With this resource, you can create webhook endpoint - [Stripe API webhook endpoint documentation](https://stripe.com/docs/api/webhook_endpoints).

You can configure webhook endpoints via the API to be notified about events that happen in your Stripe account or connected accounts.

## Example Usage

```hcl
resource "stripe_webhook_endpoint" "webhook" {
  url            = "https://webhook-url-consumer.com"
  description    = "example of webhook"
  enabled_events = [
    "customer.subscription.created", 
    "customer.subscription.updated"
  ]
}
```

## Argument Reference

Arguments accepted by this resource include:

* `url` - (Required) String. The URL of the webhook endpoint.
* `enabled_events` - (Required) List(String). The list of events to enable for this endpoint. `[*]` indicates that all events are enabled, except those that require explicit selection. All supported events listed here: [Stripe event types](https://stripe.com/docs/api/events/types).
* `description` - (Optional) String. Description of what the webhook is used for.
* `metadata` - (Optional) Map(String). Set of key-value pairs that you can attach to an object. This can be useful for storing additional information about the object in a structured format.
* `disabled` - (Optional) Bool. Disable the webhook endpoint if set to `true`. Can be used only for modification already existing webhook endpoint.

## Attribute Reference

Attributes exported by this resource include:

* `id` - String. The unique identifier for the object.
* `enabled_events` - List(String). The list of events to enable for this endpoint.
* `url` - String. The URL of the webhook endpoint.
* `description` - String. An optional description of what the webhook is used for.
* `disabled` - Bool. Informs whether the webhook endpoint is disabled.
* `metadata` - Map(String). Set of key-value pairs attached to an object.
* `secret` - String. The endpoint’s secret, used to generate webhook signatures. This field is marked as `sensitive`.