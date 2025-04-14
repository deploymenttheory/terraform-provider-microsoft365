package graphBetaMacOSPKGApp

import (
	"context"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// preserveHCLValues updates specific attributes in the state that need to be preserved from hcl configuration
// This ensures values that aren't stored in the API but are part of the configuration
// are properly maintained in the Terraform state
// preserveHCLValues updates specific attributes in the state that need to be preserved from hcl configuration
// This ensures values that aren't stored in the API but are part of the configuration
// are properly maintained in the Terraform state
func preserveHCLValues(ctx context.Context, state *resource.UpdateResponse, appMetadata *sharedmodels.MobileAppMetaDataResourceModel, appIcon *sharedmodels.MobileAppIconResourceModel) {
	tflog.Debug(ctx, "Preserving specific HCL configuration values in state")

	// Preserve installer source paths from app_metadata
	if appMetadata != nil {
		state.State.SetAttribute(ctx, path.Root("app_metadata").AtName("installer_file_path_source"),
			appMetadata.InstallerFilePathSource)
		state.State.SetAttribute(ctx, path.Root("app_metadata").AtName("installer_url_source"),
			appMetadata.InstallerURLSource)

		tflog.Debug(ctx, "Preserved installer source paths", map[string]interface{}{
			"file_path_source": appMetadata.InstallerFilePathSource.ValueString(),
			"url_source":       appMetadata.InstallerURLSource.ValueString(),
		})
	}

	// Preserve icon source paths from app_icon
	if appIcon != nil {
		state.State.SetAttribute(ctx, path.Root("app_icon").AtName("icon_file_path_source"),
			appIcon.IconFilePathSource)
		state.State.SetAttribute(ctx, path.Root("app_icon").AtName("icon_url_source"),
			appIcon.IconURLSource)

		tflog.Debug(ctx, "Preserved app icon paths", map[string]interface{}{
			"icon_file_path_source": appIcon.IconFilePathSource.ValueString(),
			"icon_url_source":       appIcon.IconURLSource.ValueString(),
		})
	}
}
