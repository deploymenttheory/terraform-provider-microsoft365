package graphBetaLinuxDeviceCompliancePolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the remote state to the Terraform schema
func MapRemoteStateToTerraform(ctx context.Context, data *LinuxDeviceCompliancePolicyResourceModel, remoteResource models.DeviceManagementCompliancePolicyable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Starting MapRemoteStateToTerraform for resource: %s", ResourceName))

	// Map base resource properties
	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.Name = convert.GraphToFrameworkString(remoteResource.GetName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())
	data.IsAssigned = convert.GraphToFrameworkBool(remoteResource.GetIsAssigned())
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())

	// Map platforms (should always be linux)
	if platforms := remoteResource.GetPlatforms(); platforms != nil {
		data.Platforms = types.StringValue("linux")
	}

	// Map technologies (should always be linuxMdm)
	if technologies := remoteResource.GetTechnologies(); technologies != nil {
		data.Technologies = types.StringValue("linuxMdm")
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
		data.Assignments = types.SetNull(LinuxDeviceCompliancePolicyAssignmentType())
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

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}

// LinuxDeviceCompliancePolicyAssignmentType returns the object type for LinuxDeviceCompliancePolicyAssignmentModel
func LinuxDeviceCompliancePolicyAssignmentType() attr.Type {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":        types.StringType,
			"group_id":    types.StringType,
			"filter_id":   types.StringType,
			"filter_type": types.StringType,
		},
	}
}

// StateConfigurationPolicySettings maps the settings from the API response to Terraform state
func StateConfigurationPolicySettings(ctx context.Context, data *LinuxDeviceCompliancePolicyResourceModel, settingsResponse models.DeviceManagementConfigurationSettingCollectionResponseable, plan *LinuxDeviceCompliancePolicyResourceModel) error {
	if settingsResponse == nil {
		tflog.Debug(ctx, "Settings response is nil")
		return nil
	}

	settings := settingsResponse.GetValue()
	if settings == nil {
		tflog.Debug(ctx, "Settings value is nil")
		return nil
	}

	tflog.Debug(ctx, fmt.Sprintf("Processing %d settings for Linux device compliance policy", len(settings)))

	// Initialize all settings to null/default values first
	data.DistributionAllowedDistros = types.ListNull(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":            types.StringType,
			"minimum_version": types.StringType,
			"maximum_version": types.StringType,
		},
	})
	data.CustomComplianceRequired = types.BoolNull()
	data.CustomComplianceDiscoveryScript = types.StringNull()
	data.CustomComplianceRules = types.StringNull()
	data.DeviceEncryptionRequired = types.BoolNull()
	data.PasswordPolicyMinimumDigits = types.Int32Null()
	data.PasswordPolicyMinimumLength = types.Int32Null()
	data.PasswordPolicyMinimumLowercase = types.Int32Null()
	data.PasswordPolicyMinimumSymbols = types.Int32Null()
	data.PasswordPolicyMinimumUppercase = types.Int32Null()

	// Process each setting
	for _, setting := range settings {
		if setting == nil {
			continue
		}

		settingInstance := setting.GetSettingInstance()
		if settingInstance == nil {
			continue
		}

		settingDefinitionId := settingInstance.GetSettingDefinitionId()
		if settingDefinitionId == nil {
			continue
		}

		tflog.Debug(ctx, fmt.Sprintf("Processing setting: %s", *settingDefinitionId))

		switch *settingDefinitionId {
		case "linux_distribution_alloweddistros":
			if err := mapGroupSettingCollectionInstanceToState(ctx, data, settingInstance); err != nil {
				return fmt.Errorf("failed to map distribution allowed distros: %w", err)
			}

		case "linux_customcompliance_required":
			if err := mapChoiceSettingInstanceToState(ctx, data, settingInstance, plan); err != nil {
				return fmt.Errorf("failed to map custom compliance required: %w", err)
			}

		case "linux_deviceencryption_required":
			if err := mapDeviceEncryptionRequiredToState(ctx, data, settingInstance); err != nil {
				return fmt.Errorf("failed to map device encryption required: %w", err)
			}

		case "linux_passwordpolicy_minimumdigits":
			if err := mapSimpleSettingInstanceWithIntegerValueToState(ctx, &data.PasswordPolicyMinimumDigits, settingInstance); err != nil {
				return fmt.Errorf("failed to map password policy minimum digits: %w", err)
			}

		case "linux_passwordpolicy_minimumlength":
			if err := mapSimpleSettingInstanceWithIntegerValueToState(ctx, &data.PasswordPolicyMinimumLength, settingInstance); err != nil {
				return fmt.Errorf("failed to map password policy minimum length: %w", err)
			}

		case "linux_passwordpolicy_minimumlowercase":
			if err := mapSimpleSettingInstanceWithIntegerValueToState(ctx, &data.PasswordPolicyMinimumLowercase, settingInstance); err != nil {
				return fmt.Errorf("failed to map password policy minimum lowercase: %w", err)
			}

		case "linux_passwordpolicy_minimumsymbols":
			if err := mapSimpleSettingInstanceWithIntegerValueToState(ctx, &data.PasswordPolicyMinimumSymbols, settingInstance); err != nil {
				return fmt.Errorf("failed to map password policy minimum symbols: %w", err)
			}

		case "linux_passwordpolicy_minimumuppercase":
			if err := mapSimpleSettingInstanceWithIntegerValueToState(ctx, &data.PasswordPolicyMinimumUppercase, settingInstance); err != nil {
				return fmt.Errorf("failed to map password policy minimum uppercase: %w", err)
			}

		default:
			tflog.Debug(ctx, fmt.Sprintf("Unknown setting definition ID: %s", *settingDefinitionId))
		}
	}

	tflog.Debug(ctx, "Finished mapping Linux device compliance policy settings to state")
	return nil
}

// mapGroupSettingCollectionInstanceToState maps the distribution allowed distros setting to state
func mapGroupSettingCollectionInstanceToState(ctx context.Context, data *LinuxDeviceCompliancePolicyResourceModel, settingInstance models.DeviceManagementConfigurationSettingInstanceable) error {
	groupCollectionInstance, ok := settingInstance.(models.DeviceManagementConfigurationGroupSettingCollectionInstanceable)
	if !ok {
		return fmt.Errorf("expected DeviceManagementConfigurationGroupSettingCollectionInstance")
	}

	groupValues := groupCollectionInstance.GetGroupSettingCollectionValue()
	if groupValues == nil || len(groupValues) == 0 {
		data.DistributionAllowedDistros = types.ListNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"type":            types.StringType,
				"minimum_version": types.StringType,
				"maximum_version": types.StringType,
			},
		})
		return nil
	}

	var distributionObjects []attr.Value
	objectType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":            types.StringType,
			"minimum_version": types.StringType,
			"maximum_version": types.StringType,
		},
	}

	for _, groupValue := range groupValues {
		if groupValue == nil {
			continue
		}

		children := groupValue.GetChildren()
		if children == nil {
			continue
		}

		distribution := map[string]attr.Value{
			"type":            types.StringNull(),
			"minimum_version": types.StringNull(),
			"maximum_version": types.StringNull(),
		}

		for _, child := range children {
			if child == nil {
				continue
			}

			childDefId := child.GetSettingDefinitionId()
			if childDefId == nil {
				continue
			}

			switch *childDefId {
			case "linux_distribution_alloweddistros_item_$type":
				if choiceInstance, ok := child.(models.DeviceManagementConfigurationChoiceSettingInstanceable); ok {
					if choiceValue := choiceInstance.GetChoiceSettingValue(); choiceValue != nil {
						if value := choiceValue.GetValue(); value != nil {
							// Extract the type from the choice value (e.g., "linux_distribution_alloweddistros_item_$type_ubuntu" -> "ubuntu")
							if len(*value) > len("linux_distribution_alloweddistros_item_$type_") {
								typeValue := (*value)[len("linux_distribution_alloweddistros_item_$type_"):]
								distribution["type"] = types.StringValue(typeValue)
							}
						}
					}
				}

			case "linux_distribution_alloweddistros_item_minimumversion":
				if simpleInstance, ok := child.(models.DeviceManagementConfigurationSimpleSettingInstanceable); ok {
					if simpleValue := simpleInstance.GetSimpleSettingValue(); simpleValue != nil {
						if stringValue, ok := simpleValue.(models.DeviceManagementConfigurationStringSettingValueable); ok {
							if value := stringValue.GetValue(); value != nil {
								distribution["minimum_version"] = types.StringValue(*value)
							}
						}
					}
				}

			case "linux_distribution_alloweddistros_item_maximumversion":
				if simpleInstance, ok := child.(models.DeviceManagementConfigurationSimpleSettingInstanceable); ok {
					if simpleValue := simpleInstance.GetSimpleSettingValue(); simpleValue != nil {
						if stringValue, ok := simpleValue.(models.DeviceManagementConfigurationStringSettingValueable); ok {
							if value := stringValue.GetValue(); value != nil {
								distribution["maximum_version"] = types.StringValue(*value)
							}
						}
					}
				}
			}
		}

		obj, diag := types.ObjectValue(objectType.AttrTypes, distribution)
		if diag.HasError() {
			return fmt.Errorf("failed to create distribution object: %v", diag.Errors())
		}
		distributionObjects = append(distributionObjects, obj)
	}

	list, diag := types.ListValue(objectType, distributionObjects)
	if diag.HasError() {
		return fmt.Errorf("failed to create distribution list: %v", diag.Errors())
	}

	data.DistributionAllowedDistros = list
	return nil
}

// mapChoiceSettingInstanceToState maps the custom compliance required setting to state
func mapChoiceSettingInstanceToState(ctx context.Context, data *LinuxDeviceCompliancePolicyResourceModel, settingInstance models.DeviceManagementConfigurationSettingInstanceable, plan *LinuxDeviceCompliancePolicyResourceModel) error {
	choiceInstance, ok := settingInstance.(models.DeviceManagementConfigurationChoiceSettingInstanceable)
	if !ok {
		return fmt.Errorf("expected DeviceManagementConfigurationChoiceSettingInstance")
	}

	choiceValue := choiceInstance.GetChoiceSettingValue()
	if choiceValue == nil {
		return fmt.Errorf("choice setting value is nil")
	}

	value := choiceValue.GetValue()
	if value == nil {
		return fmt.Errorf("choice value is nil")
	}

	// Map the choice value to boolean
	required := *value == "linux_customcompliance_required_true"
	data.CustomComplianceRequired = types.BoolValue(required)

	// Process children if custom compliance is enabled
	if required {
		children := choiceValue.GetChildren()
		if children != nil {
			for _, child := range children {
				if child == nil {
					continue
				}

				childDefId := child.GetSettingDefinitionId()
				if childDefId == nil {
					continue
				}

				switch *childDefId {
				case "linux_customcompliance_discoveryscript":
					if simpleInstance, ok := child.(models.DeviceManagementConfigurationSimpleSettingInstanceable); ok {
						if simpleValue := simpleInstance.GetSimpleSettingValue(); simpleValue != nil {
							if referenceValue, ok := simpleValue.(models.DeviceManagementConfigurationReferenceSettingValueable); ok {
								if value := referenceValue.GetValue(); value != nil {
									data.CustomComplianceDiscoveryScript = types.StringValue(*value)
								}
							}
						}
					}

				case "linux_customcompliance_rules":
					// Use planned value for sensitive data like custom compliance rules
					if plan != nil && !plan.CustomComplianceRules.IsNull() && !plan.CustomComplianceRules.IsUnknown() {
						data.CustomComplianceRules = plan.CustomComplianceRules
					} else {
						if simpleInstance, ok := child.(models.DeviceManagementConfigurationSimpleSettingInstanceable); ok {
							if simpleValue := simpleInstance.GetSimpleSettingValue(); simpleValue != nil {
								if stringValue, ok := simpleValue.(models.DeviceManagementConfigurationStringSettingValueable); ok {
									if value := stringValue.GetValue(); value != nil {
										data.CustomComplianceRules = types.StringValue(*value)
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return nil
}

// mapDeviceEncryptionRequiredToState maps the device encryption required setting to state
func mapDeviceEncryptionRequiredToState(ctx context.Context, data *LinuxDeviceCompliancePolicyResourceModel, settingInstance models.DeviceManagementConfigurationSettingInstanceable) error {
	choiceInstance, ok := settingInstance.(models.DeviceManagementConfigurationChoiceSettingInstanceable)
	if !ok {
		return fmt.Errorf("expected DeviceManagementConfigurationChoiceSettingInstance")
	}

	choiceValue := choiceInstance.GetChoiceSettingValue()
	if choiceValue == nil {
		return fmt.Errorf("choice setting value is nil")
	}

	value := choiceValue.GetValue()
	if value == nil {
		return fmt.Errorf("choice value is nil")
	}

	// Map the choice value to boolean
	required := *value == "linux_deviceencryption_required_true"
	data.DeviceEncryptionRequired = types.BoolValue(required)

	return nil
}

// mapSimpleSettingInstanceWithIntegerValueToState maps a password policy integer setting to state
func mapSimpleSettingInstanceWithIntegerValueToState(ctx context.Context, target *types.Int32, settingInstance models.DeviceManagementConfigurationSettingInstanceable) error {
	simpleInstance, ok := settingInstance.(models.DeviceManagementConfigurationSimpleSettingInstanceable)
	if !ok {
		return fmt.Errorf("expected DeviceManagementConfigurationSimpleSettingInstance")
	}

	simpleValue := simpleInstance.GetSimpleSettingValue()
	if simpleValue == nil {
		return fmt.Errorf("simple setting value is nil")
	}

	integerValue, ok := simpleValue.(models.DeviceManagementConfigurationIntegerSettingValueable)
	if !ok {
		return fmt.Errorf("expected DeviceManagementConfigurationIntegerSettingValue")
	}

	value := integerValue.GetValue()
	if value == nil {
		return fmt.Errorf("integer value is nil")
	}

	*target = types.Int32Value(*value)
	return nil
}
