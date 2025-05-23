// MapRemoteResourceStateToTerraform states the base properties of a SettingsCatalogProfileResourceModel to a Terraform state
package graphBetaMacOSPlatformScript

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
		"resourceId": remoteResource.GetId(),
	})

	data.ID = types.StringPointerValue(remoteResource.GetId())
	data.DisplayName = types.StringPointerValue(remoteResource.GetDisplayName())
	data.Description = types.StringPointerValue(remoteResource.GetDescription())
	data.CreatedDateTime = state.TimeToString(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = state.TimeToString(remoteResource.GetLastModifiedDateTime())
	data.RunAsAccount = state.EnumPtrToTypeString(remoteResource.GetRunAsAccount())
	data.FileName = types.StringPointerValue(remoteResource.GetFileName())

	data.RoleScopeTagIds = state.StringSliceToSet(ctx, remoteResource.GetRoleScopeTagIds())

	data.BlockExecutionNotifications = types.BoolPointerValue(remoteResource.GetBlockExecutionNotifications())
	data.ExecutionFrequency = state.ISO8601DurationToString(remoteResource.GetExecutionFrequency())
	data.ScriptContent = state.DecodeBase64ToString(ctx, string(remoteResource.GetScriptContent()))

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))

}
