package resources

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func MetaData(ctx context.Context, stateMetadata, planMetadata types.Map) (map[string]string, diag.Diagnostics) {
	var stateMeta, planMeta map[string]string
	var diags diag.Diagnostics

	diags = planMetadata.ElementsAs(ctx, &planMeta, true)
	if diags.HasError() {
		return nil, diags
	}
	diags = stateMetadata.ElementsAs(ctx, &stateMeta, true)
	if diags.HasError() {
		return nil, diags
	}

	for key := range stateMeta {
		if _, set := planMeta[key]; !set {
			planMeta[key] = ""
		}
	}

	return planMeta, diags
}
