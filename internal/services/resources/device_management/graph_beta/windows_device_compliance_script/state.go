package graphBetaWindowsDeviceComplianceScript

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the remote DeviceComplianceScript resource state to Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *DeviceComplianceScriptResourceModel, remoteResource graphmodels.DeviceComplianceScriptable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceName": remoteResource.GetDisplayName(),
		"resourceId":   remoteResource.GetId(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.Publisher = convert.GraphToFrameworkString(remoteResource.GetPublisher())
	data.Version = convert.GraphToFrameworkString(remoteResource.GetVersion())
	data.RunAs32Bit = convert.GraphToFrameworkBool(remoteResource.GetRunAs32Bit())
	data.EnforceSignatureCheck = convert.GraphToFrameworkBool(remoteResource.GetEnforceSignatureCheck())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())
	data.RunAsAccount = convert.GraphToFrameworkEnum(remoteResource.GetRunAsAccount())
	data.DetectionScriptContent = convert.GraphToFrameworkBytes(remoteResource.GetDetectionScriptContent())
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state")
}
