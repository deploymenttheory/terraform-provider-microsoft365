package graphBetaWindowsRemediationScript

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the remote DeviceHealthScript resource state to Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *DeviceHealthScriptResourceModel, remoteResource graphmodels.DeviceHealthScriptable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.Publisher = convert.GraphToFrameworkString(remoteResource.GetPublisher())
	data.RunAs32Bit = convert.GraphToFrameworkBool(remoteResource.GetRunAs32Bit())
	data.EnforceSignatureCheck = convert.GraphToFrameworkBool(remoteResource.GetEnforceSignatureCheck())
	data.Version = convert.GraphToFrameworkString(remoteResource.GetVersion())
	data.IsGlobalScript = convert.GraphToFrameworkBool(remoteResource.GetIsGlobalScript())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())
	data.HighestAvailableVersion = convert.GraphToFrameworkString(remoteResource.GetHighestAvailableVersion())
	data.RunAsAccount = convert.GraphToFrameworkEnum(remoteResource.GetRunAsAccount())
	data.DeviceHealthScriptType = convert.GraphToFrameworkEnum(remoteResource.GetDeviceHealthScriptType())
	data.DetectionScriptContent = convert.GraphToFrameworkBytes(remoteResource.GetDetectionScriptContent())
	data.RemediationScriptContent = convert.GraphToFrameworkBytes(remoteResource.GetRemediationScriptContent())
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))

}
