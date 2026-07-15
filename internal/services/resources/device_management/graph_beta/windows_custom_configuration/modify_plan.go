package graphBetaWindowsCustomConfiguration

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ModifyPlan validates the OMA settings list:
//   - each value must parse according to its odata_type and already be in the canonical form
//     the Graph API returns on read, otherwise the post-apply Read would rewrite the value and
//     fail with "Provider produced inconsistent result after apply"
//   - file_name is only allowed for file based setting types
//   - oma_uri must be unique across settings (encrypted value resolution matches by OMA-URI)
func (r *WindowsCustomConfigurationResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan WindowsCustomConfigurationResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.OmaSettings.IsNull() || plan.OmaSettings.IsUnknown() {
		return
	}

	var settingModels []OmaSettingResourceModel
	diags := plan.OmaSettings.ElementsAs(ctx, &settingModels, false)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	seenOmaUris := make(map[string]int, len(settingModels))

	for idx, settingModel := range settingModels {
		if !settingModel.OmaUri.IsNull() && !settingModel.OmaUri.IsUnknown() {
			omaUri := settingModel.OmaUri.ValueString()
			if firstIdx, seen := seenOmaUris[omaUri]; seen {
				resp.Diagnostics.AddError(
					"Duplicate OMA-URI",
					fmt.Sprintf("oma_settings[%d]: oma_uri %q is already used by oma_settings[%d]. "+
						"Each OMA setting must target a unique OMA-URI.", idx, omaUri, firstIdx),
				)
			} else {
				seenOmaUris[omaUri] = idx
			}
		}

		if settingModel.OdataType.IsNull() || settingModel.OdataType.IsUnknown() ||
			settingModel.Value.IsNull() || settingModel.Value.IsUnknown() {
			continue
		}

		odataType := settingModel.OdataType.ValueString()
		value := settingModel.Value.ValueString()

		if !settingModel.FileName.IsNull() && !settingModel.FileName.IsUnknown() &&
			odataType != "#microsoft.graph.omaSettingBase64" && odataType != "#microsoft.graph.omaSettingStringXml" {
			resp.Diagnostics.AddError(
				"Invalid OMA Setting",
				fmt.Sprintf("oma_settings[%d]: file_name is only applicable when odata_type is "+
					"#microsoft.graph.omaSettingBase64 or #microsoft.graph.omaSettingStringXml, got %s.", idx, odataType),
			)
		}

		parsedValue, err := parseOmaSettingValue(odataType, value)
		if err != nil {
			resp.Diagnostics.AddError(
				"Invalid OMA Setting Value",
				fmt.Sprintf("oma_settings[%d]: %s.", idx, err.Error()),
			)
			continue
		}

		if parsedValue.canonical != value {
			resp.Diagnostics.AddError(
				"Non-canonical OMA Setting Value",
				fmt.Sprintf("oma_settings[%d]: value %q is not in the canonical form %q that the Graph API "+
					"returns on read. Use the canonical form to avoid a persistent diff "+
					"(booleans: true/false, integers without leading zeros, floats without trailing zeros, "+
					"timestamps in UTC, e.g. 2024-01-01T00:00:00Z).", idx, value, parsedValue.canonical),
			)
		}
	}
}
