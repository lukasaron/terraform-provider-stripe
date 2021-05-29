
terraform {
  required_providers {
    stripe = {
      source  = "local/lukasaron/stripe"
      version = "0.0.5"
    }
  }
}

provider "stripe" {}
