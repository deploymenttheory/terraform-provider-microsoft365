package graphBetaWindowsPlatformScript

import (
	"context"
	"encoding/base64"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform states the base properties of a SettingsCatalogProfileResourceModel to a Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *WindowsPlatformScriptResourceModel, remoteResource graphmodels.DeviceManagementScriptable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": state.StringPtrToString(remoteResource.GetId()),
	})

	data.ID = types.StringValue(state.StringPtrToString(remoteResource.GetId()))
	data.DisplayName = types.StringValue(state.StringPtrToString(remoteResource.GetDisplayName()))
	data.Description = types.StringValue(state.StringPtrToString(remoteResource.GetDescription()))
	decoded, err := base64.StdEncoding.DecodeString(string(remoteResource.GetScriptContent()))
	if err != nil {
		tflog.Warn(ctx, "Failed to decode base64 script content", map[string]interface{}{
			"error": err.Error(),
		})
		data.ScriptContent = types.StringValue(string(remoteResource.GetScriptContent())) // Use original if decode fails
		return
	}
	data.ScriptContent = types.StringValue(string(decoded))
	data.RunAsAccount = state.EnumPtrToTypeString(remoteResource.GetRunAsAccount())
	data.EnforceSignatureCheck = types.BoolValue(state.BoolPtrToBool(remoteResource.GetEnforceSignatureCheck()))
	data.FileName = types.StringValue(state.StringPtrToString(remoteResource.GetFileName()))
	data.RunAs32Bit = types.BoolValue(state.BoolPtrToBool(remoteResource.GetRunAs32Bit()))
	data.RoleScopeTagIds = state.SliceToTypeStringSlice(remoteResource.GetRoleScopeTagIds())

	tflog.Debug(ctx, "Finished mapping remote resource state to Terraform state", map[string]interface{}{
		"resourceId": data.ID.ValueString(),
	})
}
