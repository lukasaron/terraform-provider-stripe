package resources

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
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
	Object        types.String `tfsdk:"object"`
	APIVersion    types.String `tfsdk:"api_version"`
	Application   types.String `tfsdk:"application"`
	Created       types.Int64  `tfsdk:"created"`
	Description   types.String `tfsdk:"description"`
	EnabledEvents types.List   `tfsdk:"enabled_events"`
	LiveMode      types.Bool   `tfsdk:"livemode"`
	Metadata      types.Map    `tfsdk:"metadata"`
	Secret        types.String `tfsdk:"secret"`
	Disabled      types.Bool   `tfsdk:"disabled"`
	URL           types.String `tfsdk:"url"`
	Connect       types.Bool   `tfsdk:"connect"`
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
			"description": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled_events": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
			},
			"object": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"metadata": schema.MapAttribute{
				Optional:    true,
				Computed:    true,
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
			"url": schema.StringAttribute{
				Required: true,
			},
			"application": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"livemode": schema.BoolAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"disabled": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
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
		},
	}
}

func (w *webhookEndpointResource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	var plan webhookEndpointModel
	diags := req.Plan.Get(ctx, &plan)
	res.Diagnostics.Append(diags...)
	if res.Diagnostics.HasError() {
		return
	}

	params := &stripe.WebhookEndpointCreateParams{}

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

	if !plan.URL.IsNull() {
		params.URL = stripe.String(plan.URL.ValueString())
	}

	if !plan.Description.IsUnknown() {
		params.Description = stripe.String(plan.Description.ValueString())
	}

	if !plan.Metadata.IsUnknown() {
		var metadata map[string]string
		res.Diagnostics.Append(plan.Metadata.ElementsAs(ctx, &metadata, false)...)
		if res.Diagnostics.HasError() {
			return
		}

		params.Metadata = metadata
	}

	if !plan.Connect.IsNull() {
		params.Connect = stripe.Bool(plan.Connect.ValueBool())
	}

	webhookEndpoint, err := w.client.V1WebhookEndpoints.Create(ctx, params)
	if err != nil {
		res.Diagnostics.AddError("WebhookEndpoint Create failed", err.Error())
		return
	}

	plan.ID = types.StringValue(webhookEndpoint.ID)
	plan.Secret = types.StringValue(webhookEndpoint.Secret)

	res.Diagnostics = w.read(ctx, &plan)
	if res.Diagnostics.HasError() {
		return
	}

	res.Diagnostics = res.State.Set(ctx, plan)
}

func (w *webhookEndpointResource) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	var state webhookEndpointModel
	diags := req.State.Get(ctx, &state)
	if diags.HasError() {
		res.Diagnostics.Append(diags...)
		return
	}

	res.Diagnostics = w.read(ctx, &state)
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
	if !state.URL.Equal(plan.URL) {
		params.URL = stripe.String(plan.URL.ValueString())
	}

	if !state.EnabledEvents.Equal(plan.EnabledEvents) {
		var enabledEvents []string
		res.Diagnostics.Append(plan.EnabledEvents.ElementsAs(ctx, &enabledEvents, false)...)
		if res.Diagnostics.HasError() {
			return
		}

		params.EnabledEvents = stripe.StringSlice(enabledEvents)
	}

	if !state.Description.Equal(plan.Description) {
		params.Description = stripe.String(plan.Description.ValueString())
	}

	if !state.Disabled.Equal(plan.Disabled) {
		params.Disabled = stripe.Bool(plan.Disabled.ValueBool())
	}

	if !state.Metadata.Equal(plan.Metadata) {
		var stateMeta, planMeta map[string]string
		res.Diagnostics.Append(plan.Metadata.ElementsAs(ctx, &planMeta, true)...)
		res.Diagnostics.Append(state.Metadata.ElementsAs(ctx, &stateMeta, true)...)
		if res.Diagnostics.HasError() {
			return
		}
		for key, value := range planMeta {
			params.AddMetadata(key, value)
		}
		for key, _ := range stateMeta {
			if _, set := params.Metadata[key]; !set {
				params.AddMetadata(key, "")
			}
		}
	}

	_, err := w.client.V1WebhookEndpoints.Update(ctx, state.ID.ValueString(), params)
	if err != nil {
		res.Diagnostics.AddError("WebhookEndpoint Update failed", err.Error())
		return
	}

	res.Diagnostics = w.read(ctx, &state)
	if res.Diagnostics.HasError() {
		return
	}

	res.Diagnostics = res.State.Set(ctx, state)
}

func (w *webhookEndpointResource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	var state webhookEndpointModel
	diags := req.State.Get(ctx, &state)
	res.Diagnostics.Append(diags...)
	if res.Diagnostics.HasError() {
		return
	}

	_, err := w.client.V1WebhookEndpoints.Delete(ctx, state.ID.ValueString(), nil)
	if err != nil {
		res.Diagnostics.AddError("WebhookEndpoint Delete failed", err.Error())
	}
}

func (w *webhookEndpointResource) ImportState(ctx context.Context, req resource.ImportStateRequest, res *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, res)
}

func (w *webhookEndpointResource) read(ctx context.Context, state *webhookEndpointModel) diag.Diagnostics {
	diags := diag.Diagnostics{}

	webhookEndpoint, err := w.client.V1WebhookEndpoints.Retrieve(ctx, state.ID.ValueString(), nil)
	if err != nil {
		diags.AddError("WebhookEndpoint Retrieve failed", err.Error())
		return diags
	}

	state.ID = types.StringValue(webhookEndpoint.ID)
	state.APIVersion = types.StringValue(webhookEndpoint.APIVersion)
	state.Description = types.StringValue(webhookEndpoint.Description)
	state.EnabledEvents, diags = types.ListValueFrom(ctx, types.StringType, webhookEndpoint.EnabledEvents)
	diags.Append(diags...)

	state.Metadata, diags = types.MapValueFrom(ctx, types.StringType, webhookEndpoint.Metadata)
	diags.Append(diags...)

	state.URL = types.StringValue(webhookEndpoint.URL)
	state.Object = types.StringValue(webhookEndpoint.Object)
	state.Application = types.StringValue(webhookEndpoint.Application)
	state.Connect = types.BoolValue(webhookEndpoint.Application != "")
	state.Created = types.Int64Value(webhookEndpoint.Created)
	state.LiveMode = types.BoolValue(webhookEndpoint.Livemode)
	state.Disabled = types.BoolValue(webhookEndpoint.Status == "disabled")

	return diags
}
