package graphCloudPcUserSetting

import (
	"context"
	"encoding/json"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

func mapRemoteStateToTerraform(ctx context.Context, data *CloudPcUserSettingResourceModel, remoteState models.CloudPcUserSettingable) {
	if remoteState == nil {
		tflog.Debug(ctx, "Remote state is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": helpers.StringPtrToString(remoteState.GetId()),
	})

	data.ID = types.StringValue(helpers.StringPtrToString(remoteState.GetId()))
	data.DisplayName = types.StringValue(helpers.StringPtrToString(remoteState.GetDisplayName()))
	data.CreatedDateTime = helpers.TimeToString(remoteState.GetCreatedDateTime())
	data.LastModifiedDateTime = helpers.TimeToString(remoteState.GetLastModifiedDateTime())
	data.LocalAdminEnabled = helpers.BoolPtrToTypeBool(remoteState.GetLocalAdminEnabled())
	data.ResetEnabled = helpers.BoolPtrToTypeBool(remoteState.GetResetEnabled())
	data.RestorePointSetting = mapRestorePointSetting(remoteState.GetRestorePointSetting())

	finalState, _ := json.MarshalIndent(data, "", "  ")
	tflog.Debug(ctx, "Finished mapping remote state to Terraform state", map[string]interface{}{
		"resourceId": data.ID.ValueString(),
		"finalState": string(finalState),
	})
}

func mapRestorePointSetting(restorePointSetting models.CloudPcRestorePointSettingable) *CloudPcRestorePointSettingModel {
	if restorePointSetting == nil {
		return nil
	}

	return &CloudPcRestorePointSettingModel{
		FrequencyType:      helpers.EnumPtrToTypeString(restorePointSetting.GetFrequencyType()),
		UserRestoreEnabled: helpers.BoolPtrToTypeBool(restorePointSetting.GetUserRestoreEnabled()),
	}
}
