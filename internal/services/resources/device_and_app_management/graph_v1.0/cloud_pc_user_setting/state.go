package graphCloudPcUserSetting

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

func MapRemoteStateToTerraform(ctx context.Context, data *CloudPcUserSettingResourceModel, remoteResource models.CloudPcUserSettingable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote state is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())
	data.LocalAdminEnabled = convert.GraphToFrameworkBool(remoteResource.GetLocalAdminEnabled())
	data.ResetEnabled = convert.GraphToFrameworkBool(remoteResource.GetResetEnabled())
	data.RestorePointSetting = mapRestorePointSetting(remoteResource.GetRestorePointSetting())

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))

}

func mapRestorePointSetting(restorePointSetting models.CloudPcRestorePointSettingable) *CloudPcRestorePointSettingModel {
	if restorePointSetting == nil {
		return nil
	}

	return &CloudPcRestorePointSettingModel{
		FrequencyType:      convert.GraphToFrameworkEnum(restorePointSetting.GetFrequencyType()),
		UserRestoreEnabled: convert.GraphToFrameworkBool(restorePointSetting.GetUserRestoreEnabled()),
	}
}
