package graphBetaMacOSCustomAttributeScript

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetamodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the remote DeviceCustomAttributeShellScript to the Terraform resource model
func MapRemoteStateToTerraform(ctx context.Context, data *DeviceCustomAttributeShellScriptResourceModel, remoteResource msgraphbetamodels.DeviceCustomAttributeShellScriptable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.CustomAttributeName = convert.GraphToFrameworkString(remoteResource.GetCustomAttributeName())
	data.CustomAttributeType = convert.GraphToFrameworkEnum(remoteResource.GetCustomAttributeType())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.ScriptContent = convert.GraphToFrameworkBytes(remoteResource.GetScriptContent())
	data.RunAsAccount = convert.GraphToFrameworkEnum(remoteResource.GetRunAsAccount())
	data.FileName = convert.GraphToFrameworkString(remoteResource.GetFileName())
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping remote state for resource %s with id %s", ResourceName, data.ID.ValueString()))
}
