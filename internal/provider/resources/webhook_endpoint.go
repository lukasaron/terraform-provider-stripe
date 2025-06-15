package resources

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stripe/stripe-go/v82"
)

var (
	_ resource.Resource                = &webhookEndpointResource{}
	_ resource.ResourceWithConfigure   = &webhookEndpointResource{}
	_ resource.ResourceWithImportState = &webhookEndpointResource{}
)

type webhookEndpointModel struct {
	ID            types.String `tfsdk:"id"`
	APIVersion    types.String `tfsdk:"api_version"`
	Application   types.String `tfsdk:"application"`
	Connect       types.Bool   `tfsdk:"connect"`
	Created       types.Int64  `tfsdk:"created"`
	Description   types.String `tfsdk:"description"`
	Disabled      types.Bool   `tfsdk:"disabled"`
	EnabledEvents types.List   `tfsdk:"enabled_events"`
	LiveMode      types.Bool   `tfsdk:"livemode"`
	Metadata      types.Map    `tfsdk:"metadata"`
	Secret        types.String `tfsdk:"secret"`
	Object        types.String `tfsdk:"object"`
	URL           types.String `tfsdk:"url"`
}

func NewWebhookEndpointResource() resource.Resource {
	return &webhookEndpointResource{}
}

type webhookEndpointResource struct {
	client *stripe.Client
}

func (w *webhookEndpointResource) Configure(_ context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
	// when client not set in req or already set in this resource
	if req.ProviderData == nil || w.client != nil {
		return
	}

	client, ok := req.ProviderData.(*stripe.Client)
	if !ok {
		res.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *stripe.Client, got %T", req.ProviderData),
		)
		return
	}

	w.client = client
}

func (w *webhookEndpointResource) Metadata(_ context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_webhook_endpoint"
}

func (w *webhookEndpointResource) Schema(_ context.Context, _ resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"api_version": schema.StringAttribute{
				Computed: true,
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"application": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"connect": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
					boolplanmodifier.RequiresReplace(),
				},
			},
			"created": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"description": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"disabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"enabled_events": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
			},
			"livemode": schema.BoolAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"metadata": schema.MapAttribute{
				Optional:    true,
				Computed:    true,
				Default:     mapdefault.StaticValue(types.MapValueMust(types.StringType, map[string]attr.Value{})),
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.UseStateForUnknown(),
				},
			},
			"secret": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"object": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"url": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (w *webhookEndpointResource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	var state, plan webhookEndpointModel
	res.Diagnostics = req.Plan.Get(ctx, &plan)
	if res.Diagnostics.HasError() {
		return
	}

	params := &stripe.WebhookEndpointCreateParams{}

	if !plan.APIVersion.IsUnknown() {
		params.APIVersion = stripe.String(plan.APIVersion.ValueString())
	}

	if !plan.Connect.IsUnknown() {
		params.Connect = stripe.Bool(plan.Connect.ValueBool())
	}

	if !plan.Description.IsNull() {
		params.Description = stripe.String(plan.Description.ValueString())
	}

	if plan.Disabled.ValueBool() {
		res.Diagnostics.AddError(
			"disabled parameter can not be set to true",
			"disabled parameter is allowed when webhook endpoint exists")
		return
	}

	if !plan.EnabledEvents.IsNull() {
		var enabledEvents []string
		res.Diagnostics.Append(plan.EnabledEvents.ElementsAs(ctx, &enabledEvents, false)...)
		if res.Diagnostics.HasError() {
			return
		}

		params.EnabledEvents = stripe.StringSlice(enabledEvents)
	}

	if !plan.Metadata.IsUnknown() {
		var metadata map[string]string
		res.Diagnostics.Append(plan.Metadata.ElementsAs(ctx, &metadata, false)...)
		if res.Diagnostics.HasError() {
			return
		}

		params.Metadata = metadata
	}

	if !plan.URL.IsNull() {
		params.URL = stripe.String(plan.URL.ValueString())
	}

	webhookEndpoint, err := w.client.V1WebhookEndpoints.Create(ctx, params)
	if err != nil {
		res.Diagnostics.AddError("WebhookEndpoint Create operation failed", err.Error())
		return
	}

	// Transition from Plan to State which is then stored
	state, res.Diagnostics = w.read(ctx, webhookEndpoint.ID)
	if res.Diagnostics.HasError() {
		return
	}

	// Secret is visible when Webhook is created ONLY
	state.Secret = types.StringValue(webhookEndpoint.Secret)

	res.Diagnostics = res.State.Set(ctx, state)
}

func (w *webhookEndpointResource) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	var state webhookEndpointModel
	res.Diagnostics = req.State.Get(ctx, &state)
	if res.Diagnostics.HasError() {
		return
	}

	secret := state.Secret // save secret for overriding in the next operation
	state, res.Diagnostics = w.read(ctx, state.ID.ValueString())
	state.Secret = secret // put secret back
	if res.Diagnostics.HasError() {
		return
	}

	res.Diagnostics = res.State.Set(ctx, state)
}

func (w *webhookEndpointResource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	var plan, state webhookEndpointModel

	res.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	res.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if res.Diagnostics.HasError() {
		return
	}

	params := &stripe.WebhookEndpointUpdateParams{}

	if !state.Description.Equal(plan.Description) {
		params.Description = stripe.String(plan.Description.ValueString())
	}

	if !state.Disabled.Equal(plan.Disabled) {
		params.Disabled = stripe.Bool(plan.Disabled.ValueBool())
	}

	if !state.EnabledEvents.Equal(plan.EnabledEvents) {
		var enabledEvents []string
		res.Diagnostics.Append(plan.EnabledEvents.ElementsAs(ctx, &enabledEvents, false)...)
		if res.Diagnostics.HasError() {
			return
		}

		params.EnabledEvents = stripe.StringSlice(enabledEvents)
	}

	if !state.Metadata.Equal(plan.Metadata) {
		metadata, diags := MetaData(ctx, state.Metadata, plan.Metadata)
		if diags.HasError() {
			res.Diagnostics = diags
			return
		}
		params.Metadata = metadata
	}

	if !state.URL.Equal(plan.URL) {
		params.URL = stripe.String(plan.URL.ValueString())
	}

	webhookEndpoint, err := w.client.V1WebhookEndpoints.Update(ctx, state.ID.ValueString(), params)
	if err != nil {
		res.Diagnostics.AddError("WebhookEndpoint Update operation failed", err.Error())
		return
	}

	secret := state.Secret // save secret for overriding in the next operation
	state, res.Diagnostics = w.read(ctx, webhookEndpoint.ID)
	state.Secret = secret // put secret back
	if res.Diagnostics.HasError() {
		return
	}

	res.Diagnostics = res.State.Set(ctx, state)
}

func (w *webhookEndpointResource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	var state webhookEndpointModel
	res.Diagnostics = req.State.Get(ctx, &state)
	if res.Diagnostics.HasError() {
		return
	}

	_, err := w.client.V1WebhookEndpoints.Delete(ctx, state.ID.ValueString(), nil)
	if err != nil {
		res.Diagnostics.AddError("WebhookEndpoint Delete operation failed", err.Error())
	}
}

func (w *webhookEndpointResource) ImportState(ctx context.Context, req resource.ImportStateRequest, res *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, res)
}

// ---------------------------------------------------------------------------------------------------------------------
func (w *webhookEndpointResource) read(ctx context.Context, id string) (webhookEndpointModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	webhookEndpoint, err := w.client.V1WebhookEndpoints.Retrieve(ctx, id, nil)
	if err != nil {
		diags.AddError("WebhookEndpoint Read operation failed", err.Error())
		return webhookEndpointModel{}, diags
	}

	return webhookEndpointModel{
			ID:          types.StringValue(webhookEndpoint.ID),
			APIVersion:  types.StringValue(webhookEndpoint.APIVersion),
			Application: types.StringValue(webhookEndpoint.Application),
			Connect:     types.BoolValue(webhookEndpoint.Application != ""),
			Created:     types.Int64Value(webhookEndpoint.Created),
			Disabled:    types.BoolValue(webhookEndpoint.Status == "disabled"),
			LiveMode:    types.BoolValue(webhookEndpoint.Livemode),
			Object:      types.StringValue(webhookEndpoint.Object),
			URL:         types.StringValue(webhookEndpoint.URL),
			Description: func() types.String {
				if webhookEndpoint.Description == "" {
					return types.StringNull()
				}
				return types.StringValue(webhookEndpoint.Description)
			}(),
			EnabledEvents: func() types.List {
				events, d := types.ListValueFrom(ctx, types.StringType, webhookEndpoint.EnabledEvents)
				if d.HasError() {
					diags.Append(d...)
				}
				return events
			}(),
			Metadata: func() types.Map {
				meta, d := types.MapValueFrom(ctx, types.StringType, webhookEndpoint.Metadata)
				if d.HasError() {
					diags.Append(d...)
				}
				return meta
			}(),
		},
		diags
}
