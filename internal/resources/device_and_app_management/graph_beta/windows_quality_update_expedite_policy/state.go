package graphBetaWindowsQualityUpdateExpeditePolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

	data.ID = types.StringPointerValue(remoteResource.GetId())
	data.DisplayName = types.StringPointerValue(remoteResource.GetDisplayName())
	data.Description = types.StringPointerValue(remoteResource.GetDescription())
	data.CreatedDateTime = state.TimeToString(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = state.TimeToString(remoteResource.GetLastModifiedDateTime())
	data.RoleScopeTagIds = state.StringSliceToSet(ctx, remoteResource.GetRoleScopeTagIds())
	data.ReleaseDateDisplayName = types.StringPointerValue(remoteResource.GetReleaseDateDisplayName())
	data.DeployableContentDisplayName = types.StringPointerValue(remoteResource.GetDeployableContentDisplayName())

	if expeditedSettings := remoteResource.GetExpeditedUpdateSettings(); expeditedSettings != nil {
		data.ExpeditedUpdateSettings = &ExpeditedWindowsQualityUpdateSettings{
			QualityUpdateRelease:  types.StringPointerValue(expeditedSettings.GetQualityUpdateRelease()),
			DaysUntilForcedReboot: state.Int32PtrToTypeInt32(expeditedSettings.GetDaysUntilForcedReboot()),
		}
	} else {
		data.ExpeditedUpdateSettings = nil
	}

	tflog.Debug(ctx, "Finished mapping remote state to Terraform", map[string]interface{}{"resourceId": data.ID.ValueString()})
}
