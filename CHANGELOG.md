## 3.0.0

* BREAKING CHANGES:
  * Remove `SubscriptionPause` from `BillingPortalConfigurationFeatures` and `BillingPortalConfigurationFeaturesParams` as the feature to pause subscription on the portal has been deprecated.
  * Rename `Features` to `MarketingFeatures` on `ProductCreateOptions`, `ProductUpdateOptions`, and `Product`

* DEPENDENCIES UPGRADE:
  * github.com/stripe/stripe-go/v76 v76.25.0 => v78.11.0

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
