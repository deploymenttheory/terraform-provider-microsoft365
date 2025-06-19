package graphBetaWindowsQualityUpdateExpeditePolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the Graph API model into the Terraform state model
func MapRemoteResourceStateToTerraform(ctx context.Context, data *WindowsQualityUpdateExpeditePolicyResourceModel, remoteResource graphmodels.WindowsQualityUpdateProfileable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Mapping remote state to Terraform", map[string]interface{}{"resourceId": remoteResource.GetId()})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())
	data.ReleaseDateDisplayName = convert.GraphToFrameworkString(remoteResource.GetReleaseDateDisplayName())
	data.DeployableContentDisplayName = convert.GraphToFrameworkString(remoteResource.GetDeployableContentDisplayName())

	if expeditedSettings := remoteResource.GetExpeditedUpdateSettings(); expeditedSettings != nil {
		data.ExpeditedUpdateSettings = &ExpeditedWindowsQualityUpdateSettings{
			QualityUpdateRelease:  convert.GraphToFrameworkString(expeditedSettings.GetQualityUpdateRelease()),
			DaysUntilForcedReboot: convert.GraphToFrameworkInt32(expeditedSettings.GetDaysUntilForcedReboot()),
		}
	} else {
		data.ExpeditedUpdateSettings = nil
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))
}
