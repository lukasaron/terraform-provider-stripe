# Stripe Terraform Provider

The Stripe Terraform provider uses the official Stripe SDK based on Golang. On top of that, the provider is developed
around the official Stripe API documentation [website](https://stripe.com/docs/api).

The Stripe Terraform Provider documentation can be found on the Terraform Provider documentation [website](https://registry.terraform.io/providers/lukasaron/stripe/latest).

## Usage:
```
terraform {
  required_providers {
    stripe = {
      source = "lukasaron/stripe"
    }
  }
}

provider "stripe" {
  api_key="<api_secret_key>"
}
```

### Environmental variable support

The parameter `api_key` can be omitted when the `STRIPE_API_KEY` environmental variable is present.

---

### Local Debugging
* Build the provider with `go build main.go`
* Move the final binary to the `mv main ~/.terraform.d/plugins/local/lukasaron/stripe/100/[platform]/terraform-provider-stripe_v100` where [platform] is `darwin_arm64` for Mac Apple chip for example.
* Create an HCL code with the following header:
 ```
terraform {
  required_providers {
    stripe = {
      source  = "local/lukasaron/stripe"
      version = "100"
    }
  }
}
```

* Run the solution from the code with the program argument `--debug`
* Copy the `TF_REATTACH_PROVIDERS` value.
* `export TF_REATTACH_PROVIDERS=[value]`
* Put breakpoints in the code
* Remove .terraform folder where the HCL code is.
* Run `terraform init` & `terraform plan` & `terraform apply`