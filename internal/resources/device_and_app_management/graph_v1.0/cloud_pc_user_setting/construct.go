package graphCloudPcUserSetting

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

func constructResource(ctx context.Context, data *CloudPcUserSettingResourceModel) (*models.CloudPcUserSetting, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := models.NewCloudPcUserSetting()

	// Set basic properties
	constructors.SetStringProperty(data.DisplayName, requestBody.SetDisplayName)
	constructors.SetBoolProperty(data.LocalAdminEnabled, requestBody.SetLocalAdminEnabled)
	constructors.SetBoolProperty(data.ResetEnabled, requestBody.SetResetEnabled)

	// Handle restore point settings
	if data.RestorePointSetting != nil {
		restorePointSetting, err := constructRestorePointSetting(data.RestorePointSetting)
		if err != nil {
			return nil, fmt.Errorf("error constructing restore point setting: %v", err)
		}
		requestBody.SetRestorePointSetting(restorePointSetting)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

func constructRestorePointSetting(data *CloudPcRestorePointSettingModel) (models.CloudPcRestorePointSettingable, error) {
	if data == nil {
		return nil, nil
	}

	restorePointSetting := models.NewCloudPcRestorePointSetting()

	if err := constructors.SetEnumProperty(data.FrequencyType,
		models.ParseCloudPcRestorePointFrequencyType,
		restorePointSetting.SetFrequencyType); err != nil {
		return nil, fmt.Errorf("failed to set frequency type: %v", err)
	}

	constructors.SetBoolProperty(data.UserRestoreEnabled, restorePointSetting.SetUserRestoreEnabled)

	return restorePointSetting, nil
}
