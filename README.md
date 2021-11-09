# Stripe Terraform Provider

The Stripe Terraform provider uses the official Stripe SDK based on Golang. On top of that, the provider is developed
around the official Stripe API documentation [website](https://stripe.com/docs/api).

The Stripe Terraform Provider documentation can be found on the Terraform Provider documentation [website](https://registry.terraform.io/providers/umisora/stripe/latest).

## Usage:
```terraform
terraform {
  required_providers {
    stripe = {
      source = "umisora/stripe"
    }
  }
}

provider "stripe" {
  api_key="<api_secret_key>"
}
```

### Environmental variable support

The parameter `api_key` can be omitted when the `STRIPE_API_KEY` environmental variable is present.
