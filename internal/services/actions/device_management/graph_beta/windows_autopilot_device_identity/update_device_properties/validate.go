package graphBetaUpdateDeviceProperties

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

// ValidateConfig validates the entire action configuration.
func (a *UpdateDevicePropertiesAction) ValidateConfig(ctx context.Context, req action.ValidateConfigRequest, resp *action.ValidateConfigResponse) {
	var data UpdateDevicePropertiesActionModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Check that at least one optional property is configured
	hasOptionalProperty := !data.UserPrincipalName.IsNull() ||
		!data.AddressableUserName.IsNull() ||
		!data.GroupTag.IsNull() ||
		!data.DisplayName.IsNull()

	if !hasOptionalProperty {
		resp.Diagnostics.AddError(
			"Missing Required Configuration",
			"At least one of the following optional attributes must be configured: "+
				"user_principal_name, addressable_user_name, group_tag, or display_name. "+
				"This action requires at least one property to update.",
		)
		return
	}

	// If user_principal_name is configured, warn if addressable_user_name is not configured
	if !data.UserPrincipalName.IsNull() && !data.UserPrincipalName.IsUnknown() {
		if data.AddressableUserName.IsNull() {
			resp.Diagnostics.AddAttributeWarning(
				path.Root("addressable_user_name"),
				"Missing Recommended Configuration",
				"When configuring user_principal_name, it is recommended to also configure "+
					"addressable_user_name for a complete user assignment. "+
					"The action may return unexpected results without both values.",
			)
		}
	}
}
