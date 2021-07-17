---
layout: "stripe"
page_title: "Provider: Stripe"
description: |-
The Stripe provider gives the power to manage few Stripe resources.
---

# Stripe Provider

The Stripe provider uses the official Stripe SDK to call it's APIs with all the benefits and downsides. More details
about specific resources are described on the [official Stripe API reference website](https://stripe.com/docs/api).

## Example Usage

```hcl
provider "stripe" {
  api_key = "<api-key>"
}
```

~> Hard-coded credentials into any Terraform configuration is not a recommended approach.
See [Environment Variables](#environment-variables) for a better alternative.

## Argument Reference

* `api-key` - (Required) Your Stripe client secret API key. This can be omitted when the environment variable `STRIPE_API_KEY` is set.

## Environment Variables

You can provide your `api-key` through the `STRIPE_API_KEY` environment variable.

```hcl
provider "stripe" {}
```

Usage:

```bash
$ export STRIPE_API_KEY="<api-key>"
$ terraform plan
```