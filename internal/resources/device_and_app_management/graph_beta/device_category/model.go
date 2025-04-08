// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-shared-devicecategory?view=graph-rest-beta
package graphBetaDeviceCategory

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DeviceCategoryResourceModel struct {
	ID              types.String   `tfsdk:"id"`
	DisplayName     types.String   `tfsdk:"display_name"`
	Description     types.String   `tfsdk:"description"`
	RoleScopeTagIds types.List     `tfsdk:"role_scope_tag_ids"`
	Timeouts        timeouts.Value `tfsdk:"timeouts"`
}
