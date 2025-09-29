package graphBetaCloudPcUserSetting

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteStateToTerraform(ctx context.Context, data *CloudPcUserSettingResourceModel, remoteResource models.CloudPcUserSettingable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())
	data.LocalAdminEnabled = convert.GraphToFrameworkBool(remoteResource.GetLocalAdminEnabled())
	data.ResetEnabled = convert.GraphToFrameworkBool(remoteResource.GetResetEnabled())
	data.SelfServiceEnabled = convert.GraphToFrameworkBool(remoteResource.GetSelfServiceEnabled())

	if restorePointSetting := remoteResource.GetRestorePointSetting(); restorePointSetting != nil {
		// Normalize frequency_type: map "twelveHours" to "default" for Terraform state
		frequencyType := ""
		if ft := restorePointSetting.GetFrequencyType(); ft != nil {
			ftStr := ft.String()
			if ftStr == "twelveHours" {
				frequencyType = "default"
			} else {
				frequencyType = ftStr
			}
		}
		data.RestorePointSetting = &RestorePointSettingModel{
			FrequencyInHours:   convert.GraphToFrameworkInt32(restorePointSetting.GetFrequencyInHours()),
			FrequencyType:      types.StringValue(frequencyType),
			UserRestoreEnabled: convert.GraphToFrameworkBool(restorePointSetting.GetUserRestoreEnabled()),
		}
	} else {
		data.RestorePointSetting = nil
	}

	if disasterRecoverySetting := remoteResource.GetCrossRegionDisasterRecoverySetting(); disasterRecoverySetting != nil {
		data.CrossRegionDisasterRecoverySetting = &CrossRegionDisasterRecoverySettingModel{
			MaintainCrossRegionRestorePointEnabled: convert.GraphToFrameworkBool(disasterRecoverySetting.GetMaintainCrossRegionRestorePointEnabled()),
			UserInitiatedDisasterRecoveryAllowed:   convert.GraphToFrameworkBool(disasterRecoverySetting.GetUserInitiatedDisasterRecoveryAllowed()),
			DisasterRecoveryType:                   convert.GraphToFrameworkEnum(disasterRecoverySetting.GetDisasterRecoveryType()),
		}

		// Handle DisasterRecoveryNetworkSetting
		if networkSetting := disasterRecoverySetting.GetDisasterRecoveryNetworkSetting(); networkSetting != nil {
			var regionNamePtr, regionGroupPtr *string
			var networkType string

			// Use OData type string to determine networkType
			odataType := ""
			if networkSetting.GetOdataType() != nil {
				odataType = *networkSetting.GetOdataType()
			}
			switch odataType {
			case "#microsoft.graph.cloudPcDisasterRecoveryMicrosoftHostedNetworkSetting":
				if ns, ok := networkSetting.(models.CloudPcDisasterRecoveryMicrosoftHostedNetworkSettingable); ok {
					if rn := ns.GetRegionName(); rn != nil && *rn != "" {
						regionNamePtr = rn
					}
					if rg := ns.GetRegionGroup(); rg != nil {
						rgStr := rg.String()
						regionGroupPtr = &rgStr
					}
				}
				networkType = "microsoftHosted"
			case "#microsoft.graph.cloudPcDisasterRecoveryAzureNetworkConnectionSetting":
				networkType = "azureNetworkConnection"
			default:
				networkType = ""
			}

			data.CrossRegionDisasterRecoverySetting.DisasterRecoveryNetworkSetting = &DisasterRecoveryNetworkSettingModel{
				NetworkType: types.StringValue(networkType),
				RegionName:  convert.GraphToFrameworkString(regionNamePtr),
				RegionGroup: convert.GraphToFrameworkString(regionGroupPtr),
			}
		} else {
			data.CrossRegionDisasterRecoverySetting.DisasterRecoveryNetworkSetting = &DisasterRecoveryNetworkSettingModel{
				NetworkType: convert.GraphToFrameworkString(nil),
				RegionName:  convert.GraphToFrameworkString(nil),
				RegionGroup: convert.GraphToFrameworkString(nil),
			}
		}
	} else {
		data.CrossRegionDisasterRecoverySetting = nil
	}

	// Handle NotificationSetting
	if notificationSetting := remoteResource.GetNotificationSetting(); notificationSetting != nil {
		data.NotificationSetting = &NotificationSettingModel{
			RestartPromptsDisabled: convert.GraphToFrameworkBool(notificationSetting.GetRestartPromptsDisabled()),
		}
	} else {
		data.NotificationSetting = nil
	}

	assignments := remoteResource.GetAssignments()
	tflog.Debug(ctx, "Retrieved assignments from remote resource", map[string]any{
		"assignmentCount": len(assignments),
		"resourceId":      data.ID.ValueString(),
	})

	if len(assignments) == 0 {
		tflog.Debug(ctx, "No assignments found, setting assignments to null", map[string]any{
			"resourceId": data.ID.ValueString(),
		})
		data.Assignments = types.SetNull(CloudPcUserSettingAssignmentType())
	} else {
		tflog.Debug(ctx, "Starting assignment mapping process", map[string]any{
			"resourceId":      data.ID.ValueString(),
			"assignmentCount": len(assignments),
		})
		MapAssignmentsToTerraform(ctx, data, assignments)
		tflog.Debug(ctx, "Completed assignment mapping process", map[string]any{
			"resourceId": data.ID.ValueString(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))
}
