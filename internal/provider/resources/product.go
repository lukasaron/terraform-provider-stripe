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
	_ resource.Resource                = &productResource{}
	_ resource.ResourceWithConfigure   = &productResource{}
	_ resource.ResourceWithImportState = &productResource{}
)

var (
	packageDimensionsObjectType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"height": types.Float64Type,
			"length": types.Float64Type,
			"weight": types.Float64Type,
			"width":  types.Float64Type,
		},
	}

	marketingFeatureObjectType = types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"name": types.StringType,
		},
	}
)

type productModel struct {
	ID                  types.String `tfsdk:"id"`
	Active              types.Bool   `tfsdk:"active"`
	Created             types.Int64  `tfsdk:"created"`
	Description         types.String `tfsdk:"description"`
	Images              types.List   `tfsdk:"images"`
	LiveMode            types.Bool   `tfsdk:"livemode"`
	MarketingFeatures   types.List   `tfsdk:"marketing_features"`
	Metadata            types.Map    `tfsdk:"metadata"`
	Name                types.String `tfsdk:"name"`
	Object              types.String `tfsdk:"object"`
	PackageDimensions   types.Object `tfsdk:"package_dimensions"`
	Shippable           types.Bool   `tfsdk:"shippable"`
	StatementDescriptor types.String `tfsdk:"statement_descriptor"`
	TaxCode             types.String `tfsdk:"tax_code"`
	UnitLabel           types.String `tfsdk:"unit_label"`
	Updated             types.Int64  `tfsdk:"updated"`
	URL                 types.String `tfsdk:"url"`
}

func NewProductResource() resource.Resource {
	return &productResource{}
}

type productResource struct {
	client *stripe.Client
}

func (p *productResource) Configure(_ context.Context, req resource.ConfigureRequest, res *resource.ConfigureResponse) {
	// when client not set in req or already set in this resource
	if req.ProviderData == nil || p.client != nil {
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

	p.client = client
}

func (p *productResource) Metadata(_ context.Context, req resource.MetadataRequest, res *resource.MetadataResponse) {
	res.TypeName = req.ProviderTypeName + "_product"
}

func (p *productResource) Schema(_ context.Context, _ resource.SchemaRequest, res *resource.SchemaResponse) {
	res.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
					stringplanmodifier.RequiresReplace(),
				},
			},
			"active": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"description": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
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
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tax_code": schema.StringAttribute{
				Optional: true,
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
			"created": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"images": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"livemode": schema.BoolAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"marketing_features": schema.ListNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Required: true,
						},
					},
				},
			},
			"package_dimensions": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"height": schema.Float64Attribute{
						Required: true,
					},
					"length": schema.Float64Attribute{
						Required: true,
					},
					"weight": schema.Float64Attribute{
						Required: true,
					},
					"width": schema.Float64Attribute{
						Required: true,
					},
				},
			},
			"shippable": schema.BoolAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"statement_descriptor": schema.StringAttribute{
				Optional: true,
			},
			"unit_label": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"updated": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"url": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (p *productResource) Create(ctx context.Context, req resource.CreateRequest, res *resource.CreateResponse) {
	var state, plan productModel
	res.Diagnostics = req.Plan.Get(ctx, &plan)
	if res.Diagnostics.HasError() {
		return
	}

	params := &stripe.ProductCreateParams{}

	if !plan.Active.IsUnknown() {
		params.Active = stripe.Bool(plan.Active.ValueBool())
	}

	if !plan.Description.IsNull() {
		params.Description = stripe.String(plan.Description.ValueString())
	}

	if !plan.ID.IsUnknown() {
		params.ID = stripe.String(plan.ID.ValueString())
	}

	if !plan.Images.IsNull() {
		var images []string
		res.Diagnostics.Append(plan.Images.ElementsAs(ctx, &images, false)...)
		if res.Diagnostics.HasError() {
			return
		}

		params.Images = stripe.StringSlice(images)
	}

	if !plan.MarketingFeatures.IsNull() {
		var list []attr.Value
		res.Diagnostics.Append(plan.MarketingFeatures.ElementsAs(ctx, &list, false)...)
		if res.Diagnostics.HasError() {
			return
		}
		params.MarketingFeatures = make([]*stripe.ProductCreateMarketingFeatureParams, 0, len(list))
		for i, featureVal := range list {
			mFeatures := featureVal.(types.Object).Attributes()
			params.MarketingFeatures[i] = &stripe.ProductCreateMarketingFeatureParams{
				Name: stripe.String(mFeatures["name"].(types.String).ValueString()),
			}
		}
	}

	if !plan.Metadata.IsUnknown() {
		var metadata map[string]string
		res.Diagnostics.Append(plan.Metadata.ElementsAs(ctx, &metadata, false)...)
		if res.Diagnostics.HasError() {
			return
		}

		params.Metadata = metadata
	}

	if !plan.Name.IsNull() {
		params.Name = stripe.String(plan.Name.ValueString())
	}

	// TODO test
	if !plan.PackageDimensions.IsUnknown() {
		pd := plan.PackageDimensions.Attributes()
		params.PackageDimensions = &stripe.ProductCreatePackageDimensionsParams{
			Height: stripe.Float64(pd["height"].(types.Float64).ValueFloat64()),
			Length: stripe.Float64(pd["length"].(types.Float64).ValueFloat64()),
			Weight: stripe.Float64(pd["weight"].(types.Float64).ValueFloat64()),
			Width:  stripe.Float64(pd["width"].(types.Float64).ValueFloat64()),
		}
	}

	if !plan.Shippable.IsNull() {
		params.Shippable = stripe.Bool(plan.Shippable.ValueBool())
	}

	if !plan.StatementDescriptor.IsNull() {
		params.StatementDescriptor = stripe.String(plan.StatementDescriptor.ValueString())
	}

	if !plan.TaxCode.IsNull() {
		params.TaxCode = stripe.String(plan.TaxCode.ValueString())
	}

	if !plan.UnitLabel.IsNull() {
		params.UnitLabel = stripe.String(plan.UnitLabel.ValueString())
	}

	if !plan.URL.IsNull() {
		params.URL = stripe.String(plan.URL.ValueString())
	}

	product, err := p.client.V1Products.Create(ctx, params)
	if err != nil {
		res.Diagnostics.AddError("Product Create operation failed", err.Error())
		return
	}

	// Transition from Plan to State which is then stored
	state, res.Diagnostics = p.read(ctx, product.ID)
	if res.Diagnostics.HasError() {
		return
	}

	res.Diagnostics = res.State.Set(ctx, state)
}

func (p *productResource) Read(ctx context.Context, req resource.ReadRequest, res *resource.ReadResponse) {
	var state productModel
	res.Diagnostics = req.State.Get(ctx, &state)
	if res.Diagnostics.HasError() {
		return
	}

	state, res.Diagnostics = p.read(ctx, state.ID.ValueString())
	if res.Diagnostics.HasError() {
		return
	}

	res.Diagnostics = res.State.Set(ctx, state)
}

func (p *productResource) Update(ctx context.Context, req resource.UpdateRequest, res *resource.UpdateResponse) {
	var plan, state productModel
	res.Diagnostics = req.Plan.Get(ctx, &plan)
	res.Diagnostics = res.State.Set(ctx, &state)
	if res.Diagnostics.HasError() {
		return
	}

	params := &stripe.ProductUpdateParams{}

	if !state.Active.Equal(plan.Active) {
		params.Active = stripe.Bool(plan.Active.ValueBool())
	}

	if !state.Description.Equal(plan.Description) {
		params.Description = stripe.String(plan.Description.ValueString())
	}

	if !state.Name.Equal(plan.Name) {
		params.Name = stripe.String(plan.Name.ValueString())
	}

	if !state.Metadata.Equal(plan.Metadata) {
		metadata, diags := MetaData(ctx, state.Metadata, plan.Metadata)
		if diags.HasError() {
			res.Diagnostics = diags
			return
		}
		params.Metadata = metadata
	}

	if !state.TaxCode.Equal(plan.TaxCode) {
		params.TaxCode = stripe.String(plan.TaxCode.ValueString())
	}

	if !state.Images.Equal(plan.Images) {
		var images []string
		res.Diagnostics.Append(plan.Images.ElementsAs(ctx, &images, true)...)
		if res.Diagnostics.HasError() {
			return
		}

		params.Images = stripe.StringSlice(images)
	}

	if !state.MarketingFeatures.Equal(plan.MarketingFeatures) {
		var list []attr.Value
		res.Diagnostics.Append(plan.MarketingFeatures.ElementsAs(ctx, &list, true)...)
		if res.Diagnostics.HasError() {
			return
		}
		params.MarketingFeatures = make([]*stripe.ProductUpdateMarketingFeatureParams, 0, len(list))
		for i, featureVal := range list {
			mFeatures := featureVal.(types.Object).Attributes()
			params.MarketingFeatures[i] = &stripe.ProductUpdateMarketingFeatureParams{
				Name: stripe.String(mFeatures["name"].(types.String).ValueString()),
			}
		}
	}

	// TODO test
	if !state.PackageDimensions.Equal(plan.PackageDimensions) {
		pd := plan.PackageDimensions.Attributes()
		params.PackageDimensions = &stripe.ProductUpdatePackageDimensionsParams{
			Height: stripe.Float64(pd["height"].(types.Float64).ValueFloat64()),
			Length: stripe.Float64(pd["length"].(types.Float64).ValueFloat64()),
			Weight: stripe.Float64(pd["weight"].(types.Float64).ValueFloat64()),
			Width:  stripe.Float64(pd["width"].(types.Float64).ValueFloat64()),
		}
	}

	if !state.Shippable.Equal(plan.Shippable) {
		params.Shippable = stripe.Bool(plan.Shippable.ValueBool())
	}

	if !state.StatementDescriptor.Equal(plan.StatementDescriptor) {
		params.StatementDescriptor = stripe.String(plan.StatementDescriptor.ValueString())
	}

	if !state.UnitLabel.Equal(plan.UnitLabel) {
		params.UnitLabel = stripe.String(plan.UnitLabel.ValueString())
	}

	if !state.URL.Equal(plan.URL) {
		params.URL = stripe.String(plan.URL.ValueString())
	}

	product, err := p.client.V1Products.Update(ctx, state.ID.ValueString(), params)
	if err != nil {
		res.Diagnostics.AddError("Product Update operation failed", err.Error())
		return
	}

	state, res.Diagnostics = p.read(ctx, product.ID)
	if res.Diagnostics.HasError() {
		return
	}

	res.Diagnostics = res.State.Set(ctx, state)
}

func (p *productResource) Delete(ctx context.Context, req resource.DeleteRequest, res *resource.DeleteResponse) {
	var state productModel
	res.Diagnostics = req.State.Get(ctx, &state)
	if res.Diagnostics.HasError() {
		return
	}

	_, err := p.client.V1Products.Delete(ctx, state.ID.ValueString(), nil)
	if err != nil {
		res.Diagnostics.AddError("Product Delete operation failed", err.Error())
	}
}

func (p *productResource) ImportState(ctx context.Context, req resource.ImportStateRequest, res *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, res)
}

// ---------------------------------------------------------------------------------------------------------------------
func (p *productResource) read(ctx context.Context, id string) (productModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	product, err := p.client.V1Products.Retrieve(ctx, id, nil)
	if err != nil {
		diags.AddError("Product Read operation failed", err.Error())
		return productModel{}, diags
	}

	return productModel{
		ID:                  types.StringValue(product.ID), // TODO continue
		Active:              types.BoolValue(product.Active),
		Description:         types.StringValue(product.Description),
		Name:                types.StringValue(product.Name),
		Object:              types.StringValue(product.Object),
		Created:             types.Int64Value(product.Created),
		LiveMode:            types.BoolValue(product.Livemode),
		Shippable:           types.BoolValue(product.Shippable),
		StatementDescriptor: types.StringValue(product.StatementDescriptor),
		UnitLabel:           types.StringValue(product.UnitLabel),
		URL:                 types.StringValue(product.URL),
		Updated:             types.Int64Value(product.Updated),
		Images: func() types.List {
			if product.Images != nil {
				images, d := types.ListValueFrom(ctx, types.StringType, product.Images)
				if d.HasError() {
					diags.Append(d...)
				}
				return images
			}
			return types.ListNull(types.StringType)
		}(),
		TaxCode: func() types.String {
			if product.TaxCode != nil {
				return types.StringValue(product.TaxCode.ID)
			}
			return types.StringNull()
		}(),
		Metadata: func() types.Map {
			meta, d := types.MapValueFrom(ctx, types.StringType, product.Metadata)
			if d.HasError() {
				diags.Append(d...)
			}
			return meta
		}(),
		PackageDimensions: func() types.Object {
			if product.PackageDimensions != nil {
				pd := map[string]attr.Value{
					"height": types.Float64Value(product.PackageDimensions.Height),
					"length": types.Float64Value(product.PackageDimensions.Length),
					"weight": types.Float64Value(product.PackageDimensions.Weight),
					"width":  types.Float64Value(product.PackageDimensions.Width),
				}
				packageDimensions, d := types.ObjectValue(packageDimensionsObjectType.AttrTypes, pd)
				if d.HasError() {
					diags.Append(d...)
				}
				return packageDimensions
			}
			return types.ObjectNull(packageDimensionsObjectType.AttrTypes)
		}(),
		MarketingFeatures: func() types.List {
			if product.MarketingFeatures != nil {
				mfValues := make([]attr.Value, len(product.MarketingFeatures))
				for i, f := range product.MarketingFeatures {
					mfObject, d := types.ObjectValue(
						marketingFeatureObjectType.AttrTypes,
						map[string]attr.Value{
							"name": types.StringValue(f.Name),
						},
					)
					if d.HasError() {
						diags.Append(d...)
					}
					mfValues[i] = mfObject
				}
				marketingFeatures, d := types.ListValue(types.ObjectType{AttrTypes: marketingFeatureObjectType.AttrTypes}, mfValues)
				if d.HasError() {
					diags.Append(d...)
				}
				return marketingFeatures
			}
			return types.ListNull(types.ObjectType{AttrTypes: marketingFeatureObjectType.AttrTypes})
		}(),
	}, diags
}
