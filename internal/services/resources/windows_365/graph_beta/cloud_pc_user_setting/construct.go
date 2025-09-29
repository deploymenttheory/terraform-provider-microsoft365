package graphBetaCloudPcUserSetting

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructResource(ctx context.Context, data *CloudPcUserSettingResourceModel) (*graphmodels.CloudPcUserSetting, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewCloudPcUserSetting()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphBool(data.LocalAdminEnabled, requestBody.SetLocalAdminEnabled)
	convert.FrameworkToGraphBool(data.ResetEnabled, requestBody.SetResetEnabled)
	convert.FrameworkToGraphBool(data.SelfServiceEnabled, requestBody.SetSelfServiceEnabled)

	if data.RestorePointSetting != nil {
		restorePointSetting := graphmodels.NewCloudPcRestorePointSetting()
		convert.FrameworkToGraphInt32(data.RestorePointSetting.FrequencyInHours, restorePointSetting.SetFrequencyInHours)

		err := convert.FrameworkToGraphEnum(data.RestorePointSetting.FrequencyType, graphmodels.ParseCloudPcRestorePointFrequencyType, restorePointSetting.SetFrequencyType)
		if err != nil {
			return nil, fmt.Errorf("error setting CloudPcRestorePointFrequencyType: %v", err)
		}

		convert.FrameworkToGraphBool(data.RestorePointSetting.UserRestoreEnabled, restorePointSetting.SetUserRestoreEnabled)
		requestBody.SetRestorePointSetting(restorePointSetting)
	}

	if data.CrossRegionDisasterRecoverySetting != nil {
		disasterRecoverySetting := graphmodels.NewCloudPcCrossRegionDisasterRecoverySetting()
		convert.FrameworkToGraphBool(data.CrossRegionDisasterRecoverySetting.MaintainCrossRegionRestorePointEnabled, disasterRecoverySetting.SetMaintainCrossRegionRestorePointEnabled)
		convert.FrameworkToGraphBool(data.CrossRegionDisasterRecoverySetting.UserInitiatedDisasterRecoveryAllowed, disasterRecoverySetting.SetUserInitiatedDisasterRecoveryAllowed)

		err := convert.FrameworkToGraphEnum(data.CrossRegionDisasterRecoverySetting.DisasterRecoveryType, graphmodels.ParseCloudPcDisasterRecoveryType, disasterRecoverySetting.SetDisasterRecoveryType)
		if err != nil {
			return nil, fmt.Errorf("error setting DisasterRecoveryType: %v", err)
		}

		if data.CrossRegionDisasterRecoverySetting.DisasterRecoveryNetworkSetting != nil {
			networkType := data.CrossRegionDisasterRecoverySetting.DisasterRecoveryNetworkSetting.NetworkType.ValueString()
			regionName := data.CrossRegionDisasterRecoverySetting.DisasterRecoveryNetworkSetting.RegionName.ValueString()
			regionGroup := data.CrossRegionDisasterRecoverySetting.DisasterRecoveryNetworkSetting.RegionGroup.ValueString()

			switch networkType {
			case "microsoftHosted":
				networkSetting := graphmodels.NewCloudPcDisasterRecoveryMicrosoftHostedNetworkSetting()
				odataType := "#microsoft.graph.cloudPcDisasterRecoveryMicrosoftHostedNetworkSetting"
				networkSetting.SetOdataType(&odataType)
				if regionName != "" {
					networkSetting.SetRegionName(&regionName)
				} else {
					defaultRegion := "automatic"
					networkSetting.SetRegionName(&defaultRegion)
				}
				if regionGroup != "" {
					if val, err := graphmodels.ParseCloudPcRegionGroup(regionGroup); err == nil && val != nil {
						networkSetting.SetRegionGroup(val.(*graphmodels.CloudPcRegionGroup))
					}
				}
				disasterRecoverySetting.SetDisasterRecoveryNetworkSetting(networkSetting)
			case "azureNetworkConnection":
				// handle azureNetworkConnection if needed
			}
		}

		requestBody.SetCrossRegionDisasterRecoverySetting(disasterRecoverySetting)
	}

	if data.NotificationSetting != nil {
		notificationSetting := graphmodels.NewCloudPcNotificationSetting()
		convert.FrameworkToGraphBool(data.NotificationSetting.RestartPromptsDisabled, notificationSetting.SetRestartPromptsDisabled)
		requestBody.SetNotificationSetting(notificationSetting)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
