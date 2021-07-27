---
layout: "stripe"
page_title: "Stripe: stripe_card"
description: |-
The Stripe Card can be created, modified, configured and removed by this resource.
---

# stripe_card

With this resource, you can create a card - [Stripe API card documentation](https://stripe.com/docs/api/cards).

You can store multiple cards on a customer in order to charge the customer later. You can also store multiple debit cards on a recipient in order to transfer to those cards later.

~> Passing your cardholder’s full credit card number to Stripe’s API isn't a recommended approach. In rare cases, you may have to continue handling full credit card information directly. If this applies to you, you can enable unsafe processing in your [dashboard](https://dashboard.stripe.com/settings/integration).

## Example Usage

```hcl
// card for the customer
resource "stripe_card" "card" {
  customer  = "cus_Ju7..."
  number    = "4242424242424242"
  name      = "Lukas Aron"
  cvc       = 123
  exp_month = 8
  exp_year  = 2030
}

// card for the customer with address
resource "stripe_card" "card" {
  customer  = "cus_Ju7..."
  number    = "4242424242424242"
  name      = "Lukas Aron"
  cvc       = 123
  exp_month = 8
  exp_year  = 2030
  address = {
    line1   = "1 The Best Street",
    line2   = "Apartment 401",
    city    = "Sydney",
    state   = "NSW",
    zip     = "2000",
    country = "Australia"
  }
}
```

## Argument Reference

Arguments accepted by this resource include:

* `customer` - (Required) String. The customer that this card belongs to.
* `number` - (Required) String. The card number, as a string without any separators.
* `exp_month` - (Required) Int. Number representing the card's expiration month.
* `exp_year` - (Required) Int. Four-digit number representing the card's expiration year.
* `cvc` - (Optional) Int. Card security code. Highly recommended to always include this value, but it's required only for accounts based in European countries.
* `name` - (Optional) String. Cardholder name.
* `address` - (Optional) Map(String). Address map with fields related to the address: `line1`, `line2`, `city`, `state`, `zip` and `country`.
* `metadata` - (Optional) Map(String). Set of key-value pairs that you can attach to an object. This can be useful for storing additional information about the object in a structured format.

## Attribute Reference

Attributes exported by this resource include:

* `id` - String. The unique identifier for the object.
* `customer` - String. The customer that this card belongs to.
* `number` - String. The card number, as a string without any separators. This field is marked as `sensitive`.
* `exp_month` - Int. Number representing the card's expiration month.
* `exp_year` - Int. Four-digit number representing the card's expiration year.
* `cvc` - Int. Card security code. This field is marked as `sensitive`.
* `name` - String. Cardholder name.
* `address` - Map(String). Address map with fields related to the address.
* `address_line1_check` - String. If address `line1` was provided, results of the check: `pass`, `fail`, `unavailable`, or `unchecked`.
* `address_zip_check` - String. If address `zip` was provided, results of the check: `pass`, `fail`, `unavailable`, or `unchecked`.
* `cvc_check` - String. If a `cvc` was provided, results of the check: `pass`, `fail`, `unavailable`, or `unchecked`. A result of `unchecked` indicates that CVC was provided but hasn’t been checked yet
* `brand` - String. Card brand. Can be `American Express`, `Diners Club`, `Discover`, `JCB`, `MasterCard`, `UnionPay`, `Visa`, or `Unknown`.
* `fingerprint` - String. Uniquely identifies this particular card number. You can use this attribute to check whether two customers who’ve signed up with you are using the same card number, for example. For payment methods that tokenize card information (Apple Pay, Google Pay), the tokenized number might be provided instead of the underlying card number.
* `tokenization_method` - String. If the card number is tokenized, this is the method that was used. Can be `android_pay` (includes Google Pay), `apple_pay`, `masterpass`, `visa_checkout`, or `null`.
* `funding` - String. Card funding type. Can be `credit`, `debit`, `prepaid`, or `unknown`.
* `last4` - String. The last four digits of the card.
* `available_payout_methods` - List(String). A set of available payout methods for this card. Only values from this set should be passed as the method when creating a payout.
* `country` - String. Two-letter ISO code representing the country of the card. You could use this attribute to get a sense of the international breakdown of cards you’ve collected.
* `metadata` - Map(String). Set of key-value pairs that you can attach to an object.