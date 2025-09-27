## 3.3.5
* BUGFIXES:
    *  Customer creation crashed, missing condition to Address field added.

* DEPENDENCIES UPGRADE
    * Go 1.24 => 1.25
    * github.com/hashicorp/go-plugin v1.6.3 => v1.7.0 
    * github.com/hashicorp/terraform-plugin-go v0.28.0 => v0.29.0 
    * github.com/hashicorp/terraform-plugin-sdk/v2 v2.37.0 => v2.38.1 
    * github.com/hashicorp/terraform-registry-address v0.3.0 => v0.4.0 
    * github.com/zclconf/go-cty v1.16.3 => v1.17.0 
    * golang.org/x/mod v0.27.0 => v0.28.0 
    * golang.org/x/net v0.43.0 => v0.44.0 
    * golang.org/x/sync v0.16.0 => v0.17.0 
    * golang.org/x/sys v0.35.0 => v0.36.0 
    * golang.org/x/text v0.28.0 => v0.29.0 
    * golang.org/x/tools v0.36.0 => v0.37.0 
    * google.golang.org/genproto/googleapis/rpc v0.0.0-20250804133106-a7a43d27e69b => v0.0.0-20250922171735-9219d122eba9 
    * google.golang.org/grpc v1.74.2 => v1.75.1 
    * google.golang.org/protobuf v1.36.7 => v1.36.9

## 3.3.3
* BUGFIXES:
  * Remove resources from the state when it is not found.
  
* DEPENDENCIES UPGRADE
  * github.com/fatih/color v1.17.0 => v1.18.0
  * github.com/google/go-cmp v0.6.0 => v0.7.0
  * github.com/hashicorp/go-cty v1.4.1-0.20200414143053-d3edf31b6320 => v1.5.0
  * github.com/hashicorp/go-plugin v1.6.1 => v1.6.3
  * github.com/hashicorp/hcl/v2 v2.22.0 => v2.24.0
  * github.com/hashicorp/terraform-plugin-go v0.24.0 => v0.28.0
  * github.com/hashicorp/terraform-plugin-sdk/v2 v2.34.0 => v2.37.0
  * github.com/hashicorp/terraform-registry-address v0.2.3 => v0.3.0
  * github.com/mattn/go-colorable v0.1.13 => v0.1.14
  * github.com/oklog/run v1.1.0 => v1.2.0
  * github.com/zclconf/go-cty v1.15.0 => v1.16.3
  * golang.org/x/mod v0.21.0 => v0.27.0
  * golang.org/x/net v0.30.0 => v0.43.0
  * golang.org/x/sync v0.8.0 => v0.16.0
  * golang.org/x/sys v0.26.0 => v0.35.0
  * golang.org/x/text v0.19.0 => v0.28.0
  * golang.org/x/tools v0.26.0 => v0.36.0
  * google.golang.org/genproto/googleapis/rpc v0.0.0-20241007155032-5fefd90f89a9 => v0.0.0-20250804133106-a7a43d27e69b
  * google.golang.org/grpc v1.67.1 => v1.74.2
  * google.golang.org/protobuf v1.35.1 => v1.36.7

## 3.3.2
* NOTES:
  * api_key is optional to unblock CDKTF.
    * api key value is still checked for availability

## 3.3.1
* BUGFIXES:
  * Price forces recreation when  `transform_quantity` changes.

## 3.3.0
* NEW RESOURCES:
  * Meter

* DEPENDENCIES UPGRADE:
  * github.com/hashicorp/terraform-plugin-go v0.23.0 => v0.24.0
  * github.com/hashicorp/yamux v0.1.1 => v0.1.2
  * golang.org/x/mod v0.20.0 => v0.21.0
  * golang.org/x/net v0.28.0 => v0.30.0
  * golang.org/x/sys v0.24.0 => v0.26.0
  * golang.org/x/text v0.17.0 => v0.19.0
  * golang.org/x/tools v0.24.0 => v0.26.0
  * google.golang.org/genproto/googleapis/rpc v0.0.0-20240826202546-f6391c0de4c7 => v0.0.0-20241007155032-5fefd90f89a9
  * google.golang.org/grpc v1.65.0 => v1.67.1
  * google.golang.org/protobuf v1.34.2 => v1.35.1

## 3.2.2
* BUGFIXES:
  * Webhook Endpoint exposes the `application` attribute.
  * Webhook Endpoint `connect` field reflects the `application` field value.

* DEPENDENCIES UPGRADE:
  * github.com/hashicorp/hcl/v2 v2.21.0 => v2.22.0
  * google.golang.org/genproto/googleapis/rpc v0.0.0-20240823204242-4ba0660f739c => v0.0.0-20240826202546-f6391c0de4c7

## 3.2.1
* BUGFIXES:
  * Price is fetching/populating `tiers` values to TF state in the read function.

* DEPENDENCIES UPGRADE:
  *  google.golang.org/genproto v0.0.0-20240823204242-4ba0660f739c
  *  golang.org/x/sys v0.23.0 => v0.24.0
  *  google.golang.org/genproto/googleapis/rpc v0.0.0-20240805194559-2c9e96a0b5d4 => v0.0.0-20240823204242-4ba0660f739c

## 3.2.0
* NEW RESOURCES:
  * Entitlements Feature
  * Product Feature
 
* DEPENDENCIES UPGRADE:
  * golang.org/x/mod v0.19.0 => v0.20.0
  * golang.org/x/net v0.27.0 => v0.28.0
  * golang.org/x/sync v0.7.0 => v0.8.0
  * golang.org/x/sys v0.22.0 => v0.23.0
  * golang.org/x/text v0.16.0 => v0.17.0
  * golang.org/x/tools v0.23.0 => v0.24.0
  * google.golang.org/genproto/googleapis/rpc v0.0.0-20240711142825-46eb208f015d => v0.0.0-20240805194559-2c9e96a0b5d4

## 3.1.0

* BUGFIXES:
  * Price has `custom_unit_amount` block implementation.

* DEPENDENCIES UPGRADE:
  * github.com/hashicorp/hcl/v2 v2.20.1 => v2.21.0
  * github.com/stripe/stripe-go/v78 v78.11.0 => v78.12.0
  * github.com/zclconf/go-cty v1.14.4 => v1.15.0
  * golang.org/x/mod v0.18.0 => v0.19.0
  * golang.org/x/net v0.26.0 => v0.27.0
  * golang.org/x/sys v0.21.0 => v0.22.0
  * golang.org/x/tools v0.22.0 => v0.23.0
  * google.golang.org/genproto/googleapis/rpc v0.0.0-20240610135401-a8a62080eff3 => v0.0.0-20240711142825-46eb208f015d
  * google.golang.org/grpc v1.64.0 => v1.65.0

## 3.0.2

* BUGFIXES:
  * Code Promotion: `expires_at` and `restrictions.minimum_amount` fields are misinterpreted as 0 when they are nil.

## 3.0.1

* BUGFIXES:
  * Code Promotion restrictions `minimum_amount` and `minimum_amount_currency` are optional.

## 3.0.0

* BREAKING CHANGES:
  * Remove `SubscriptionPause` from `BillingPortalConfigurationFeatures` and `BillingPortalConfigurationFeaturesParams` as the feature to pause subscription on the portal has been deprecated.
  * Rename `Features` to `MarketingFeatures` on `ProductCreateOptions`, `ProductUpdateOptions`, and `Product`

* DEPENDENCIES UPGRADE:
  * github.com/stripe/stripe-go/v76 v76.25.0 => v78.11.0
  * github.com/hashicorp/go-version v1.6.0 => v1.7.0
  * golang.org/x/mod v0.17.0 => v0.18.0
  * golang.org/x/net v0.25.0 => v0.26.0
  * golang.org/x/sys v0.20.0 => v0.21.0
  * golang.org/x/text v0.15.0 => v0.16.0
  * golang.org/x/tools v0.21.0 => v0.22.0
  * google.golang.org/genproto/googleapis/rpc v0.0.0-20240515191416-fc5f0ca64291 => v0.0.0-20240610135401-a8a62080eff3
  * google.golang.org/protobuf v1.34.1 => v1.34.2


## 2.0.0

* BUGFIXES:
  * `tax_behaviour` was renamed to `tax_behavior` to be consistent with the Stripe language. This is the breaking change!

* DEPENDENCIES UPGRADE:
    * github.com/hashicorp/go-plugin v1.6.0 => v1.6.1
    * github.com/hashicorp/terraform-plugin-go v0.22.1 => v0.23.0
    * github.com/stripe/stripe-go/v76 v76.23.0 => v76.25.0
    * golang.org/x/mod v0.16.0 => v0.17.0
    * golang.org/x/net v0.23.0 => v0.25.0
    * golang.org/x/sync v0.6.0 => v0.7.0
    * golang.org/x/sys v0.18.0 => v0.20.0
    * golang.org/x/text v0.14.0 => v0.15.0
    * golang.org/x/tools v0.19.0 => v0.21.0
    * google.golang.org/genproto/googleapis/rpc v0.0.0-20240401170217-c3f982113cda => v0.0.0-20240509183442-62759503f434
    * google.golang.org/grpc v1.63.0 => v1.63.2
    * google.golang.org/protobuf v1.33.0 => v1.34.1

## 1.XX.XX

* NOTES:
    * Production ready solution which is already widely in use.
