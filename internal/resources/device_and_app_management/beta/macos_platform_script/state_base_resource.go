// MapRemoteResourceStateToTerraform states the base properties of a SettingsCatalogProfileResourceModel to a Terraform state
package graphBetaMacOSPlatformScript

import (
	"context"
	"encoding/base64"
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
		"resourceId": state.StringPtrToString(remoteResource.GetId()),
	})

	data.ID = types.StringValue(state.StringPtrToString(remoteResource.GetId()))
	data.DisplayName = types.StringValue(state.StringPtrToString(remoteResource.GetDisplayName()))
	data.Description = types.StringValue(state.StringPtrToString(remoteResource.GetDescription()))
	data.CreatedDateTime = state.TimeToString(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = state.TimeToString(remoteResource.GetLastModifiedDateTime())
	data.RunAsAccount = state.EnumPtrToTypeString(remoteResource.GetRunAsAccount())
	data.FileName = types.StringValue(state.StringPtrToString(remoteResource.GetFileName()))
	data.RoleScopeTagIds = state.SliceToTypeStringSlice(remoteResource.GetRoleScopeTagIds())
	data.BlockExecutionNotifications = types.BoolValue(state.BoolPtrToBool(remoteResource.GetBlockExecutionNotifications()))

	// ExecutionFrequency (ISO Duration)
	if executionFrequency := remoteResource.GetExecutionFrequency(); executionFrequency != nil {
		data.ExecutionFrequency = types.StringValue(fmt.Sprintf("%v", executionFrequency.String()))
	}

	// ScriptContent (base64 encoded)
	decodedContent, err := base64.StdEncoding.DecodeString(string(remoteResource.GetScriptContent()))
	if err != nil {
		tflog.Warn(ctx, "Failed to decode base64 script content", map[string]interface{}{
			"error": err.Error(),
		})
		data.ScriptContent = types.StringValue(string(remoteResource.GetScriptContent())) // Use original if decode fails
	} else {
		data.ScriptContent = types.StringValue(string(decodedContent))
	}

	tflog.Debug(ctx, "Finished mapping remote resource state to Terraform state", map[string]interface{}{
		"resourceId": data.ID.ValueString(),
	})
}
