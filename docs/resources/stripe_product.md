---
layout: "stripe"
page_title: "Stripe: stripe_product"
description: |- The Stripe Product can be created, modified, configured and removed by this resource.
---

# stripe_product

With this resource, you can create a product - [Stripe API product documentation](https://stripe.com/docs/api/products).

Products describe the specific goods or services you offer to your customers. For example, 
you might offer a Standard and Premium version of your goods or service; each version would be a separate Product.

~> When there is at least one price used in the product, Stripe returns error when delete operation is performed. 
This provider updates the active status of the product to archived in this scenario and removes product from the 
terraform state.

## Example Usage

```hcl
// The minimal product
resource "stripe_product" "product" {
  name = "minimalist product"
}

// The product with more fields filled
resource "stripe_product" "product" {
  name        = "full product"
  unit_label  = "piece"
  description = "fantastic product"
  url         = "https://www.terraform.io"
}

```

## Argument Reference

Arguments accepted by this resource include:

* `product_id` - (Optional) String. The bespoke unique identifier for the object.
* `name` - (Required) String. The product’s name, meant to be displayable to the customer. Whenever this product is sold via a subscription, name will show up on associated invoice line item descriptions.
* `active` - (Optional) Bool. Whether the product is currently available for purchase. Defaults to `true`.
* `description` - (Optional) String. The product’s description, meant to be displayable to the customer. Use this field to optionally store a long form explanation of the product being sold for your own rendering purposes.
* `marketing_features` - (Optional) List(String). A list of up to 15 marketing features for this product. These are displayed in pricing tables.
* `images` - (Optional) List(String). A list of up to 8 URLs of images for this product, meant to be displayable to the customer.
* `package_dimensions` - (Optional) Map(Float). The dimensions of this product for shipping purposes. When used these fields are required: `height`,`length`,`width` and `weight`; the precision is 2 decimal places.
* `shippable` - (Optional) Bool. Whether this product is shipped (i.e., physical goods).
* `statement_descriptor` - (Optional) String. An arbitrary string to be displayed on your customer’s credit card or bank statement. While most banks display this information consistently, some may display it incorrectly or not at all. This may be up to 22 characters. The statement description may not include `<`,` >`, `\`, `"`, `’` characters, and will appear on your customer’s statement in capital letters. Non-ASCII characters are automatically stripped. It must contain at least one letter.
* `tax_code` - (Optional) String. A tax code ID. Supported values are listed in the TaxCode resource and at https://stripe.com/docs/tax/tax-categories.
* `unit_label` - (Optional) String. A label that represents units of this product in Stripe and on customers’ receipts and invoices. When set, this will be included in associated invoice line item descriptions.
* `url` - (Optional) String. A URL of a publicly-accessible webpage for this product.
* `metadata` - (Optional) Map(String). Set of key-value pairs that you can attach to an object. This can be useful for storing additional information about the object in a structured format.

## Attribute Reference

Attributes exported by this resource include:

* `id` - String. The unique identifier for the object.
* `product_id` - String. The unique identifier for the object.
* `name` - String. The product’s name, meant to be displayable to the customer. 
* `active` - Bool. Whether the product is currently available for purchase. 
* `description` - String. The product’s description, meant to be displayable to the customer.
* `marketing_features` - List(String). A list of marketing features for this product.
* `images` - List(String). A list of up to 8 URLs of images for this product.
* `package_dimensions` - Map(Float). The dimensions of this product for shipping purposes.
* `shippable` - Bool. Whether this product is shipped (i.e., physical goods).
* `statement_descriptor` - String. An arbitrary string to be displayed on your customer’s credit card or bank statement.
* `tax_code` - String. A tax code ID.
* `unit_label` - String. A label that represents units of this product in Stripe and on customers’ receipts and invoices. 
* `url` - String. A URL of a publicly-accessible webpage for this product.
* `metadata` - Map(String). Set of key-value pairs that you can attach to an object.