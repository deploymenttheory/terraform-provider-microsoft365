package graphIOSMobileAppConfiguration

import (
	"context"
	"encoding/base64"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

func mapResourceToDataSourceState(
	ctx context.Context,
	resource models.IosMobileAppConfigurationable,
	data *IOSMobileAppConfigurationDataSourceModel,
	diagnostics *diag.Diagnostics,
) {
	if resource == nil {
		tflog.Debug(ctx, "Resource is nil, skipping state mapping")
		return
	}

	tflog.Debug(ctx, "Mapping iOS Mobile App Configuration resource to data source state")

	if resource.GetId() != nil {
		data.Id = types.StringValue(*resource.GetId())
	}

	if resource.GetDisplayName() != nil {
		data.DisplayName = types.StringValue(*resource.GetDisplayName())
	}

	if resource.GetDescription() != nil {
		data.Description = types.StringValue(*resource.GetDescription())
	} else {
		data.Description = types.StringNull()
	}

	// Handle targeted mobile apps - treat empty list as null to match Terraform behavior
	targetedApps := resource.GetTargetedMobileApps()
	if len(targetedApps) == 0 {
		data.TargetedMobileApps = types.ListNull(types.StringType)
	} else {
		data.TargetedMobileApps = convert.GraphToFrameworkStringList(targetedApps)
	}

	// Handle encoded setting XML - treat empty byte arrays as null
	encodedXml := resource.GetEncodedSettingXml()
	if len(encodedXml) > 0 {
		encodedString := base64.StdEncoding.EncodeToString(encodedXml)
		data.EncodedSettingXml = types.StringValue(encodedString)
	} else {
		data.EncodedSettingXml = types.StringNull()
	}

	if resource.GetSettings() != nil {
		data.Settings = mapSettingsToDataSourceState(ctx, resource.GetSettings(), diagnostics)
	}

	data.CreatedDateTime = convert.GraphToFrameworkTime(resource.GetCreatedDateTime())

	data.LastModifiedDateTime = convert.GraphToFrameworkTime(resource.GetLastModifiedDateTime())

	data.Version = convert.GraphToFrameworkInt32(resource.GetVersion())
}

func mapSettingsToDataSourceState(
	ctx context.Context,
	settings []models.AppConfigurationSettingItemable,
	diagnostics *diag.Diagnostics,
) []IOSMobileAppConfigurationSetting {
	if len(settings) == 0 {
		return nil
	}

	result := make([]IOSMobileAppConfigurationSetting, 0, len(settings))

	for _, setting := range settings {
		if setting == nil {
			continue
		}

		stateSetting := IOSMobileAppConfigurationSetting{}

		if setting.GetAppConfigKey() != nil {
			stateSetting.AppConfigKey = types.StringValue(*setting.GetAppConfigKey())
		}

		if setting.GetAppConfigKeyType() != nil {
			// Convert enum to string - the enum value needs to be converted to its string representation
			keyType := setting.GetAppConfigKeyType()
			if keyType != nil {
				stateSetting.AppConfigKeyType = types.StringValue(keyType.String())
			}
		}

		if setting.GetAppConfigKeyValue() != nil {
			stateSetting.AppConfigKeyValue = types.StringValue(*setting.GetAppConfigKeyValue())
		}

		result = append(result, stateSetting)
	}

	return result
}

func mapAssignmentsToDataSourceState(
	ctx context.Context,
	assignments []models.ManagedDeviceMobileAppConfigurationAssignmentable,
	diagnostics *diag.Diagnostics,
) []IOSMobileAppConfigurationAssignment {
	if len(assignments) == 0 {
		return nil
	}

	result := make([]IOSMobileAppConfigurationAssignment, 0, len(assignments))

	for _, assignment := range assignments {
		if assignment == nil {
			continue
		}

		stateAssignment := IOSMobileAppConfigurationAssignment{}

		if assignment.GetId() != nil {
			stateAssignment.Id = types.StringValue(*assignment.GetId())
		}

		if assignment.GetTarget() != nil {
			stateAssignment.Target = mapAssignmentTargetToDataSourceState(
				ctx,
				assignment.GetTarget(),
				diagnostics,
			)
		}

		result = append(result, stateAssignment)
	}

	return result
}

func mapAssignmentTargetToDataSourceState(
	ctx context.Context,
	target models.DeviceAndAppManagementAssignmentTargetable,
	diagnostics *diag.Diagnostics,
) *IOSMobileAppConfigurationAssignmentTarget {
	if target == nil {
		return nil
	}

	result := &IOSMobileAppConfigurationAssignmentTarget{}

	if target.GetOdataType() != nil {
		result.ODataType = types.StringValue(*target.GetOdataType())
	}

	// Check if target has a GroupId (both GroupAssignmentTargetable and ExclusionGroupAssignmentTargetable have this)
	if typedTarget, ok := target.(models.GroupAssignmentTargetable); ok {
		if typedTarget.GetGroupId() != nil {
			result.GroupId = types.StringValue(*typedTarget.GetGroupId())
		}
	}

	return result
}
