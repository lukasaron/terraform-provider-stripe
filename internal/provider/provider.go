package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/lukasaron/terraform-provider-stripe/internal/provider/resources"
	"github.com/stripe/stripe-go/v82"
	"os"
)

var (
	_ provider.Provider = &stripeProvider{}
)

type stripeProviderModel struct {
	APIKey types.String `tfsdk:"api_key"`
}
type stripeProvider struct{}

func New() func() provider.Provider {
	return func() provider.Provider {
		return &stripeProvider{}
	}
}

func (s stripeProvider) Metadata(ctx context.Context, req provider.MetadataRequest, res *provider.MetadataResponse) {
	res.TypeName = "stripe"
}

func (s stripeProvider) Schema(_ context.Context, req provider.SchemaRequest, res *provider.SchemaResponse) {
	res.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (s stripeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, res *provider.ConfigureResponse) {
	var providerModel stripeProviderModel
	diags := req.Config.Get(ctx, &providerModel)
	res.Diagnostics.Append(diags...)
	if res.Diagnostics.HasError() {
		return
	}

	apiKey := os.Getenv("STRIPE_API_KEY")
	if !providerModel.APIKey.IsNull() {
		apiKey = providerModel.APIKey.ValueString()
	}

	if apiKey == "" {
		res.Diagnostics.AddAttributeError(
			path.Root("api_key"),
			"The provider cannot create Stripe Client without API key",
			"Set the STRIPE_API_KEY environment variable or set the api_key attribute in the provider configuration.",
		)
	}

	if res.Diagnostics.HasError() {
		return
	}

	client := stripe.NewClient(apiKey, nil)

	res.ResourceData = client
	res.DataSourceData = client
}

func (s stripeProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (s stripeProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewWebhookEndpointResource,
	}
}
