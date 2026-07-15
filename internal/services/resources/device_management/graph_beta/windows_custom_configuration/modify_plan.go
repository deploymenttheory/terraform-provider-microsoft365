package graphBetaWindowsCustomConfiguration

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ModifyPlan validates that each OMA setting value can be converted to the data type
// implied by its odata_type, and that file_name is only used with file based setting types.
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

	for idx, settingModel := range settingModels {
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

		switch odataType {
		case "#microsoft.graph.omaSettingInteger":
			if _, err := strconv.ParseInt(value, 10, 32); err != nil {
				resp.Diagnostics.AddError(
					"Invalid OMA Setting Value",
					fmt.Sprintf("oma_settings[%d]: value %q is not a valid integer.", idx, value),
				)
			}
		case "#microsoft.graph.omaSettingBoolean":
			if _, err := strconv.ParseBool(value); err != nil {
				resp.Diagnostics.AddError(
					"Invalid OMA Setting Value",
					fmt.Sprintf("oma_settings[%d]: value %q is not a valid boolean.", idx, value),
				)
			}
		case "#microsoft.graph.omaSettingDateTime":
			if _, err := time.Parse(time.RFC3339, value); err != nil {
				resp.Diagnostics.AddError(
					"Invalid OMA Setting Value",
					fmt.Sprintf("oma_settings[%d]: value %q is not a valid RFC3339 timestamp (e.g. 2024-01-01T00:00:00Z).", idx, value),
				)
			}
		case "#microsoft.graph.omaSettingFloatingPoint":
			if _, err := strconv.ParseFloat(value, 32); err != nil {
				resp.Diagnostics.AddError(
					"Invalid OMA Setting Value",
					fmt.Sprintf("oma_settings[%d]: value %q is not a valid floating point number.", idx, value),
				)
			}
		}
	}
}
