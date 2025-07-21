package graphIOSMobileAppConfiguration

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

func constructResource(
	ctx context.Context,
	data *IOSMobileAppConfigurationResourceModel,
	diagnostics *diag.Diagnostics,
) models.IosMobileAppConfigurationable {
	tflog.Debug(ctx, "Constructing iOS Mobile App Configuration resource")

	resource := models.NewIosMobileAppConfiguration()

	// Set display name
	convert.FrameworkToGraphString(data.DisplayName, resource.SetDisplayName)

	// Set description
	convert.FrameworkToGraphString(data.Description, resource.SetDescription)

	// Set targeted mobile apps
	if err := convert.FrameworkToGraphStringList(ctx, data.TargetedMobileApps, resource.SetTargetedMobileApps); err != nil {
		diagnostics.AddError(
			"Failed to convert targeted mobile apps",
			fmt.Sprintf("Failed to set targeted mobile apps: %s", err),
		)
		return nil
	}

	// Set encoded setting XML
	if !data.EncodedSettingXml.IsNull() && !data.EncodedSettingXml.IsUnknown() {
		xmlString := data.EncodedSettingXml.ValueString()
		encodedBytes, err := base64.StdEncoding.DecodeString(xmlString)
		if err != nil {
			diagnostics.AddError(
				"Failed to decode encoded_setting_xml",
				err.Error(),
			)
			return nil
		}
		resource.SetEncodedSettingXml(encodedBytes)
	}

	// Set settings
	if len(data.Settings) > 0 {
		settings := make([]models.AppConfigurationSettingItemable, 0, len(data.Settings))
		for _, setting := range data.Settings {
			settingItem := models.NewAppConfigurationSettingItem()

			// Use helper function for simple string pointers
			if !setting.AppConfigKey.IsNull() && !setting.AppConfigKey.IsUnknown() {
				settingItem.SetAppConfigKey(helpers.StringPtr(setting.AppConfigKey.ValueString()))
			}

			// Use the convert helper for enum conversion
			err := convert.FrameworkToGraphBitmaskEnum(
				setting.AppConfigKeyType,
				models.ParseMdmAppConfigKeyType,
				settingItem.SetAppConfigKeyType,
			)
			if err != nil {
				diagnostics.AddError(
					"Failed to convert app config key type",
					fmt.Sprintf("Error converting app config key type: %s", err),
				)
				return nil
			}

			if !setting.AppConfigKeyValue.IsNull() && !setting.AppConfigKeyValue.IsUnknown() {
				settingItem.SetAppConfigKeyValue(
					helpers.StringPtr(setting.AppConfigKeyValue.ValueString()),
				)
			}

			settings = append(settings, settingItem)
		}
		resource.SetSettings(settings)
	}

	// RoleScopeTagIds is not available in v1.0 API for iOS Mobile App Configuration

	return resource
}

func constructAssignments(
	ctx context.Context,
	data *IOSMobileAppConfigurationResourceModel,
	diagnostics *diag.Diagnostics,
) []models.ManagedDeviceMobileAppConfigurationAssignmentable {
	if len(data.Assignments) == 0 {
		return nil
	}

	assignments := make(
		[]models.ManagedDeviceMobileAppConfigurationAssignmentable,
		0,
		len(data.Assignments),
	)

	for _, assignment := range data.Assignments {
		if assignment.Target == nil {
			continue
		}

		assignmentItem := models.NewManagedDeviceMobileAppConfigurationAssignment()

		target := constructAssignmentTarget(ctx, assignment.Target, diagnostics)
		if target != nil {
			assignmentItem.SetTarget(target)
			assignments = append(assignments, assignmentItem)
		}
	}

	return assignments
}

func constructAssignmentTarget(
	ctx context.Context,
	target *IOSMobileAppConfigurationAssignmentTarget,
	diagnostics *diag.Diagnostics,
) models.DeviceAndAppManagementAssignmentTargetable {
	if target == nil {
		return nil
	}

	odataType := target.ODataType.ValueString()

	switch odataType {
	case "#microsoft.graph.allLicensedUsersAssignmentTarget":
		return models.NewAllLicensedUsersAssignmentTarget()

	case "#microsoft.graph.allDevicesAssignmentTarget":
		return models.NewAllDevicesAssignmentTarget()

	case "#microsoft.graph.groupAssignmentTarget":
		groupTarget := models.NewGroupAssignmentTarget()
		if !target.GroupId.IsNull() && !target.GroupId.IsUnknown() {
			groupTarget.SetGroupId(helpers.StringPtr(target.GroupId.ValueString()))
		}
		return groupTarget

	case "#microsoft.graph.exclusionGroupAssignmentTarget":
		exclusionTarget := models.NewExclusionGroupAssignmentTarget()
		if !target.GroupId.IsNull() && !target.GroupId.IsUnknown() {
			exclusionTarget.SetGroupId(helpers.StringPtr(target.GroupId.ValueString()))
		}
		return exclusionTarget

	default:
		diagnostics.AddError(
			"Invalid assignment target type",
			"Unknown OData type: "+odataType,
		)
		return nil
	}
}
