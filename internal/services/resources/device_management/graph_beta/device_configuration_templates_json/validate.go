package graphBetaDeviceConfigurationTemplatesJson

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	sharedconstructors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors/graph_beta/device_management"
)

func (r *DeviceConfigurationTemplatesJsonResource) ValidateConfig(
	ctx context.Context,
	req resource.ValidateConfigRequest,
	resp *resource.ValidateConfigResponse,
) {
	var data DeviceConfigurationTemplatesJsonResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() || data.Settings.IsNull() || data.Settings.IsUnknown() {
		return
	}
	// Metadata and assignments must have one owner in the configuration.
	if _, err := sharedconstructors.ConstructDeviceConfigurationSettings(
		data.Settings,
	); err != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("settings"),
			"Invalid device configuration settings",
			err.Error(),
		)
	}
}
