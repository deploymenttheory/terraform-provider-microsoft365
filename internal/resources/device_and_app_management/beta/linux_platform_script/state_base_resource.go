package graphBetaLinuxPlatformScript

import (
	"context"
	"encoding/base64"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform states the base properties of a SettingsCatalogProfileResourceModel to a Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *LinuxPlatformScriptResourceModel, remoteResource graphmodels.DeviceManagementScriptable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": state.StringPtrToString(remoteResource.GetId()),
	})

	data.ID = types.StringPointerValue(remoteResource.GetId())
	data.DisplayName = types.StringPointerValue(remoteResource.GetDisplayName())
	data.Description = types.StringPointerValue(remoteResource.GetDescription())

	// Handle base64 encoded script content
	decoded, err := base64.StdEncoding.DecodeString(string(remoteResource.GetScriptContent()))
	if err != nil {
		tflog.Warn(ctx, "Failed to decode base64 script content", map[string]interface{}{
			"error": err.Error(),
		})
		data.ScriptContent = types.StringValue(string(remoteResource.GetScriptContent()))
		return
	}
	data.ScriptContent = types.StringValue(string(decoded))

	data.RoleScopeTagIds = state.SliceToTypeStringSlice(remoteResource.GetRoleScopeTagIds())

	tflog.Debug(ctx, "Finished mapping remote resource state to Terraform state", map[string]interface{}{
		"resourceId": data.ID.ValueString(),
	})
}
