---
layout: "stripe"
page_title: "Stripe: stripe_file"
description: |- 
  The Stripe File can be created by this resource.
---

# stripe_file

This object represents files hosted on Stripe's servers - [Stripe API file documentation](https://stripe.com/docs/api/files).
You can upload files with the create file request (for example, when uploading dispute evidence).

Stripe File upload [guide](https://stripe.com/docs/file-upload#uploading-a-file)

~> Update or removal of the file isn't supported through the Stripe SDK.

## Example Usage

```hcl
// minimal file field
resource "stripe_file" "logo" {
  filename = "logo.jpg"
  purpose  = "business_logo"
  base64content  = filebase64("${HOME}/logo.jpg")
}

// file with file link
resource "stripe_file" "logo" {
  filename = "logo.jpg"
  purpose  = "business_logo"
  base64content  = filebase64("${HOME}/logo.jpg")
  link_data {
    create     = true
    expires_at = 1826659124
  }
}
```

## Argument Reference

Arguments accepted by this resource include:

* `filename` - (Required) String. The suitable name for saving the file to a filesystem.
* `purpose` - (Required) String. The purpose of the uploaded file. One of these values are accepted: `account_requirement`,
  `additional_verification`, `business_icon`, `business_logo`, `customer_signature`, `dispute_evidence`,
  `document_provider_identity_document`, `finance_report_run`, `identity_document`, `identity_document_downloadable`,
  `pci_document`, `selfie`, `sigma_scheduled_query`, `tax_document_user_upload`, `terminal_reader_splashscreen`
* `base64content` (Required) String. A content file to upload encoded by Base64, 
   ideally use Terraform function [filebase64](https://developer.hashicorp.com/terraform/language/functions/filebase64) .
* `file_link_data` - (Optional) List(Resource). Parameter that automatically create a file link for the newly created file.
   Please see details [File Link Data](#file-link-data).

### File Link Data

`file_link_data` Supports the following arguments:

* `create` - (Required) Bool. Set this to `true` to create a file link for the newly created file. 
   Creating a link is only possible when the file’s purpose is one of the following: `business_icon`, `business_logo`, 
   `customer_signature`, `dispute_evidence`, `pci_document`, `tax_document_user_upload`, or `terminal_reader_splashscreen`.
* `expires_at` - (Optional) Int. The link isn’t available after this future timestamp.
* `metadata` - (Optional) Map(String). Set of key-value pairs that you can attach to an object. 
   This can be useful for storing additional information about the object in a structured format.

## Attribute Reference

Attributes exported by this resource include:

* `id` - String. The unique identifier for the object.
* `type` - String. The returned file type (for example, `csv`, `pdf`, `jpg`, or `png`).
* `filename` - String. The suitable name for saving the file to a filesystem.
* `base64content` - String. Content of the file encoded by Base64

* `purpose` - String. The purpose of the uploaded file.
* `object` - String. String representing the object’s type. Objects of the same type share the same value.
* `created` - Int. Time at which the object was created. Measured in seconds since the Unix epoch.
* `expires_at` - Int. The file expires and isn’t available at this time in epoch seconds.
* `size` - Int. The size of the file object in bytes.
* `url` - String. Use your live secret API key to download the file from this URL.
* `links` - List(Resource). A list of [file links](https://stripe.com/docs/api/files/object#file_links) that point at this file.
   Please see details of [links](#links).

### Links

`links` exports these resources:

* `id` - String. Unique identifier for the object.
* `object` - String. String representing the object’s type. Objects of the same type share the same value.
* `created` - String. Time at which the object was created. Measured in seconds since the Unix epoch.
* `expired` - Bool. Returns if the link is already expired.
* `expires_at` - Int. Time that the link expires.
* `livemode` - Bool. Has the value `true` if the object exists in live mode or the value `false` 
   if the object exists in test mode.
* `metadata` - Map(String). Set of key-value pairs that you can attach to an object.
* `url` - String. The publicly accessible URL to download the file.

## Import

Import is supported using the following syntax:

```shell
$ terraform import stripe_file.file <file_id>
```