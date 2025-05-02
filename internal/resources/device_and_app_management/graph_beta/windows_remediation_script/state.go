package graphBetaWindowsRemediationScript

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the remote DeviceHealthScript resource state to Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *DeviceHealthScriptResourceModel, remoteResource graphmodels.DeviceHealthScriptable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	data.ID = types.StringPointerValue(remoteResource.GetId())
	data.DisplayName = types.StringPointerValue(remoteResource.GetDisplayName())
	data.Description = types.StringPointerValue(remoteResource.GetDescription())
	data.Publisher = types.StringPointerValue(remoteResource.GetPublisher())
	data.RunAs32Bit = types.BoolPointerValue(remoteResource.GetRunAs32Bit())
	data.EnforceSignatureCheck = types.BoolPointerValue(remoteResource.GetEnforceSignatureCheck())
	data.Version = types.StringPointerValue(remoteResource.GetVersion())
	data.IsGlobalScript = types.BoolPointerValue(remoteResource.GetIsGlobalScript())
	data.CreatedDateTime = state.TimeToString(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = state.TimeToString(remoteResource.GetLastModifiedDateTime())
	data.HighestAvailableVersion = types.StringPointerValue(remoteResource.GetHighestAvailableVersion())
	data.RunAsAccount = state.EnumPtrToTypeString(remoteResource.GetRunAsAccount())
	data.DeviceHealthScriptType = state.EnumPtrToTypeString(remoteResource.GetDeviceHealthScriptType())
	data.DetectionScriptContent = state.BytesToString(remoteResource.GetDetectionScriptContent())
	data.RemediationScriptContent = state.BytesToString(remoteResource.GetRemediationScriptContent())
	data.RoleScopeTagIds = state.StringSliceToSet(ctx, remoteResource.GetRoleScopeTagIds())

	tflog.Debug(ctx, "Finished mapping remote resource state to Terraform state")
}
