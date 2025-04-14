package graphBetaMacOSPKGApp

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
			"installer_size_in_bytes":    types.Int64Type,
			"installer_md5_checksum":     types.StringType,
			"installer_sha256_checksum":  types.StringType,
		},
	}

	if metadata == nil {
		// Return a null object with all required fields defined
		return types.ObjectNull(objectType.AttrTypes)
	}

	// Initialize values with nulls first
	values := map[string]attr.Value{
		"installer_file_path_source": types.StringNull(),
		"installer_url_source":       types.StringNull(),
		"installer_size_in_bytes":    types.Int64Null(),
		"installer_md5_checksum":     types.StringNull(),
		"installer_sha256_checksum":  types.StringNull(),
	}

	// Then override with actual values if they exist
	if !metadata.InstallerFilePathSource.IsNull() {
		values["installer_file_path_source"] = metadata.InstallerFilePathSource
	}

	if !metadata.InstallerURLSource.IsNull() {
		values["installer_url_source"] = metadata.InstallerURLSource
	}

	if !metadata.InstallerSizeInBytes.IsNull() {
		values["installer_size_in_bytes"] = metadata.InstallerSizeInBytes
	}

	if !metadata.InstallerMD5Checksum.IsNull() {
		values["installer_md5_checksum"] = metadata.InstallerMD5Checksum
	}

	if !metadata.InstallerSHA256Checksum.IsNull() {
		values["installer_sha256_checksum"] = metadata.InstallerSHA256Checksum
	}

	objValue, diags := types.ObjectValue(objectType.AttrTypes, values)
	if diags.HasError() {
		tflog.Warn(ctx, fmt.Sprintf("Error creating AppMetadata object: %v", diags.Errors()))
		return types.ObjectNull(objectType.AttrTypes)
	}

	return objValue
}
