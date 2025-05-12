package sharedStater

import (
	"context"
	"fmt"

	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/shared_models/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// MapAppMetadataStateToTerraform is a standalone function that ensures AppMetadata is properly initialized
// with all required fields to prevent type conversion errors
func MapAppMetadataStateToTerraform(ctx context.Context, metadata *sharedmodels.MobileAppMetaDataResourceModel) types.Object {
	objectType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"installer_file_path_source": types.StringType,
			"installer_url_source":       types.StringType,
		},
	}

	if metadata == nil {
		return types.ObjectNull(objectType.AttrTypes)
	}

	values := map[string]attr.Value{
		"installer_file_path_source": metadata.InstallerFilePathSource,
		"installer_url_source":       metadata.InstallerURLSource,
	}

	objValue, diags := types.ObjectValue(objectType.AttrTypes, values)
	if diags.HasError() {
		tflog.Warn(ctx, fmt.Sprintf("Error creating AppMetadata object: %v", diags.Errors()))
		return types.ObjectNull(objectType.AttrTypes)
	}

	return objValue
}
