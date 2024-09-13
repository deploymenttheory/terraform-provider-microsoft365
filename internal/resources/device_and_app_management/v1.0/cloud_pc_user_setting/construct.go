package graphCloudPcUserSetting

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

func constructResource(ctx context.Context, data *CloudPcUserSettingResourceModel) (*models.CloudPcUserSetting, error) {
	tflog.Debug(ctx, "Constructing CloudPcUserSetting resource")

	requestBody := models.NewCloudPcUserSetting()

	if !data.DisplayName.IsNull() && !data.DisplayName.IsUnknown() {
		displayName := data.DisplayName.ValueString()
		requestBody.SetDisplayName(&displayName)
	}

	if !data.LocalAdminEnabled.IsNull() && !data.LocalAdminEnabled.IsUnknown() {
		localAdminEnabled := data.LocalAdminEnabled.ValueBool()
		requestBody.SetLocalAdminEnabled(&localAdminEnabled)
	}

	if !data.ResetEnabled.IsNull() && !data.ResetEnabled.IsUnknown() {
		resetEnabled := data.ResetEnabled.ValueBool()
		requestBody.SetResetEnabled(&resetEnabled)
	}

	if data.RestorePointSetting != nil {
		restorePointSetting, err := constructRestorePointSetting(data.RestorePointSetting)
		if err != nil {
			return nil, fmt.Errorf("error constructing restore point setting: %v", err)
		}
		requestBody.SetRestorePointSetting(restorePointSetting)
	}

	// Debug logging
	debugPrintRequestBody(ctx, requestBody)

	return requestBody, nil
}

func constructRestorePointSetting(data *CloudPcRestorePointSettingModel) (models.CloudPcRestorePointSettingable, error) {
	if data == nil {
		return nil, nil
	}

	restorePointSetting := models.NewCloudPcRestorePointSetting()

	if !data.FrequencyType.IsNull() && !data.FrequencyType.IsUnknown() {
		frequencyTypeStr := data.FrequencyType.ValueString()
		frequencyTypeAny, err := models.ParseCloudPcRestorePointFrequencyType(frequencyTypeStr)
		if err != nil {
			return nil, fmt.Errorf("error parsing frequency type: %v", err)
		}
		if frequencyTypeAny != nil {
			frequencyType, ok := frequencyTypeAny.(*models.CloudPcRestorePointFrequencyType)
			if !ok {
				return nil, fmt.Errorf("unexpected type for frequency type: %T", frequencyTypeAny)
			}
			restorePointSetting.SetFrequencyType(frequencyType)
		}
	}

	if !data.UserRestoreEnabled.IsNull() && !data.UserRestoreEnabled.IsUnknown() {
		userRestoreEnabled := data.UserRestoreEnabled.ValueBool()
		restorePointSetting.SetUserRestoreEnabled(&userRestoreEnabled)
	}

	return restorePointSetting, nil
}

func debugPrintRequestBody(ctx context.Context, requestBody *models.CloudPcUserSetting) {
	requestMap := map[string]interface{}{
		"displayName":         requestBody.GetDisplayName(),
		"localAdminEnabled":   requestBody.GetLocalAdminEnabled(),
		"resetEnabled":        requestBody.GetResetEnabled(),
		"restorePointSetting": debugMapRestorePointSetting(requestBody.GetRestorePointSetting()),
	}

	requestBodyJSON, err := json.MarshalIndent(requestMap, "", "  ")
	if err != nil {
		tflog.Error(ctx, "Error marshalling request body to JSON", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	tflog.Debug(ctx, "Constructed CloudPcUserSetting resource", map[string]interface{}{
		"requestBody": string(requestBodyJSON),
	})
}

func debugMapRestorePointSetting(setting models.CloudPcRestorePointSettingable) map[string]interface{} {
	if setting == nil {
		return nil
	}
	return map[string]interface{}{
		"frequencyType":      setting.GetFrequencyType(),
		"userRestoreEnabled": setting.GetUserRestoreEnabled(),
	}
}
