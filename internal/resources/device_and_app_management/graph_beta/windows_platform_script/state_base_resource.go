package graphBetaWindowsPlatformScript

import (
	"context"
	"encoding/base64"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteResourceStateToTerraform(ctx context.Context, data *WindowsPlatformScriptResourceModel, remoteResource graphmodels.DeviceManagementScriptable) {
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
	data.EnforceSignatureCheck = types.BoolPointerValue(remoteResource.GetEnforceSignatureCheck())
	data.FileName = types.StringPointerValue(remoteResource.GetFileName())
	data.RunAs32Bit = types.BoolPointerValue(remoteResource.GetRunAs32Bit())

	var roleScopeTagIds []attr.Value
	for _, v := range state.SliceToTypeStringSlice(remoteResource.GetRoleScopeTagIds()) {
		roleScopeTagIds = append(roleScopeTagIds, v)
	}
	data.RoleScopeTagIds = types.ListValueMust(
		types.StringType,
		roleScopeTagIds,
	)

	tflog.Debug(ctx, "Finished mapping remote resource state to Terraform state", map[string]interface{}{
		"resourceId": data.ID.ValueString(),
	})
}