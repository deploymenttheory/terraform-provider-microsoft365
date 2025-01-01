package graphCloudPcUserSetting

import (
	"context"
	"encoding/json"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

func MapRemoteStateToTerraform(ctx context.Context, data *CloudPcUserSettingResourceModel, remoteState models.CloudPcUserSettingable) {
	if remoteState == nil {
		tflog.Debug(ctx, "Remote state is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": state.StringPtrToString(remoteState.GetId()),
	})

	data.ID = types.StringPointerValue(remoteState.GetId())
	data.DisplayName = types.StringPointerValue(remoteState.GetDisplayName())
	data.CreatedDateTime = state.TimeToString(remoteState.GetCreatedDateTime())
	data.LastModifiedDateTime = state.TimeToString(remoteState.GetLastModifiedDateTime())
	data.LocalAdminEnabled = types.BoolPointerValue(remoteState.GetLocalAdminEnabled())
	data.ResetEnabled = types.BoolPointerValue(remoteState.GetResetEnabled())
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
		FrequencyType:      state.EnumPtrToTypeString(restorePointSetting.GetFrequencyType()),
		UserRestoreEnabled: types.BoolPointerValue(restorePointSetting.GetUserRestoreEnabled()),
	}
}
