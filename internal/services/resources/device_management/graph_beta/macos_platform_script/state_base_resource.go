// MapRemoteResourceStateToTerraform states the base properties of a SettingsCatalogProfileResourceModel to a Terraform state
package graphBetaMacOSPlatformScript

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the base properties of a MacOSPlatformScriptResourceModel to a Terraform state.
func MapRemoteResourceStateToTerraform(ctx context.Context, data *MacOSPlatformScriptResourceModel, remoteResource graphmodels.DeviceShellScriptable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceName": remoteResource.GetDisplayName(),
		"resourceId":   remoteResource.GetId(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())
	data.RunAsAccount = convert.GraphToFrameworkEnum(remoteResource.GetRunAsAccount())
	data.FileName = convert.GraphToFrameworkString(remoteResource.GetFileName())
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())
	data.BlockExecutionNotifications = convert.GraphToFrameworkBool(remoteResource.GetBlockExecutionNotifications())
	data.ExecutionFrequency = convert.GraphToFrameworkISODuration(remoteResource.GetExecutionFrequency())
	data.ScriptContent = convert.GraphToFrameworkBytes(remoteResource.GetScriptContent())
	data.RetryCount = convert.GraphToFrameworkInt32(remoteResource.GetRetryCount())

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))
}
