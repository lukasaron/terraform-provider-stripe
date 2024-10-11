---
layout: "stripe"
page_title: "Stripe: stripe_meter"
description: |-
  The Stripe Meter can be created, modified and configured by this resource.
---

# stripe_meter

With this resource, you can create a billing meter - [Stripe API billing meter documentation](https://docs.stripe.com/api/billing/meter).

A billing meter is a resource that allows you to track usage of a particular event. For example, you might create a billing meter to track the number of API calls made by a particular user. You can then attach the billing meter to a price and attach the price to a subscription to charge the user for the number of API calls they make.

Related guide: [Usage based billing](https://docs.stripe.com/billing/subscriptions/usage-based)


## Example Usage

```hcl
// meter
resource "stripe_meter" "sample_meter" {
  display_name = "A Sample meter"
  event_name   = "sample_meter"

  customer_mapping {
    event_payload_key = "stripe_customer_id"
    type              = "by_id"
  }

  default_aggregation {
    formula = "sum"
  }

  value_settings {
    event_payload_key = "value"
  }
}


resource "stripe_price" "price" {
  // product needs to be defined
  product        = stripe_product.product.id
  currency       = "aud"
  billing_scheme = "tiered"
  tiers_mode     = "graduated"

  # free up to ten
  tiers {
    up_to       = 10
    unit_amount = 0 // can be omitted
  }

  tiers {
    up_to       = 100
    unit_amount = 300
  }

  tiers {
    up_to               = -1
    unit_amount_decimal = 100.5
  }

  recurring {
    interval        = "week"
    interval_count  = 2
    usage_type      = "metered"
    meter          = stripe_meter.sample_meter.id
  }
}

```

## Argument Reference

Arguments accepted by this resource include:

* `display_name` - (Required) String. The display name of the meter.
* `event_name` - (Required) String. The name of the meter event to record usage for. Corresponds with the `event_name` field on meter events.
* `event_time_window` - (Optional) String. The time window to pre-aggregate meter events for, if any. Possible values are:
 - `day` - Events are pre-aggregated in daily buckets
 - `hour` - Events are pre-aggregated in hourly buckets

### Default Aggregation

`default_aggregation` Supports the following arguments:

* `formula` - (Required) String. Specifies how events are aggregated. Allowed values are `count` to count the number of events and `sum` to sum each event’s value.

### Customer Mapping

`customer_mapping` Supports the following arguments:

* `event_payload_key` - (Required) String. The key in the event payload to use for customer mapping.
* `type` - (Required) String. The method for mapping a meter event to a customer. Must be `by_id`.

### Value Settings

`value_settings` Supports the following arguments:

`event_payload_key` - (Required) String. The key in the usage event payload to use as the value for this meter. For example, if the event payload contains usage on a bytes_used field, then set the event_payload_key to “bytes_used”.

## Attribute Reference

Attributes exported by this resource include:

* `id` - String. The unique identifier for the object.
* `display_name` - String. The display name of the meter.
* `event_name` - String. The name of the meter event to record usage for. Corresponds with the `event_name` field on meter events.
* `event_time_window` - String. The time window to pre-aggregate meter events for, if any.
* `default_aggregation` - Resource. The default settings to aggregate a meter’s events with. Fields that specify how to aggregate a meter event such as `formula`.
* `customer_mapping` - Resource. Fields that specify how to map a meter event to a customer such as `event_payload_key` and `type`.
* `value_settings` - Resource. Fields that specify how to calculate a meter event’s value such as `event_payload_key`.

## Note on updating meters

Once created, you can update the `display_name`.

Other attribute edits will trigger a destroy action (archival) and creation of a new meter entry.

## Import

Import is supported using the following syntax:

```shell
$ terraform import stripe_price.meter <meter_id>
```

