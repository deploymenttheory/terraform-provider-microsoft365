package graphBetaMacOSCustomAttributeScript

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the remote DeviceCustomAttributeShellScript to the Terraform resource model
func MapRemoteStateToTerraform(ctx context.Context, data *DeviceCustomAttributeShellScriptResourceModel, remoteResource graphmodels.DeviceCustomAttributeShellScriptable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId":   remoteResource.GetId(),
		"resourceName": remoteResource.GetDisplayName(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.CustomAttributeType = convert.GraphToFrameworkEnum(remoteResource.GetCustomAttributeType())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.ScriptContent = convert.GraphToFrameworkBytes(remoteResource.GetScriptContent())
	data.RunAsAccount = convert.GraphToFrameworkEnum(remoteResource.GetRunAsAccount())
	data.FileName = convert.GraphToFrameworkString(remoteResource.GetFileName())
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())

	assignments := remoteResource.GetAssignments()
	tflog.Debug(ctx, "Retrieved assignments from remote resource", map[string]any{
		"assignmentCount": len(assignments),
		"resourceId":      data.ID.ValueString(),
	})

	if len(assignments) == 0 {
		data.Assignments = types.SetNull(MacOSCustomAttributeScriptAssignmentType())
	} else {
		MapAssignmentsToTerraform(ctx, data, assignments)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}
