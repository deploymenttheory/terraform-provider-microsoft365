package graphBetaAutopatchGroups

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// mapTypeStringToInt converts type strings to integers for API requests
// "Device" = 0, "None" = 1
func mapTypeStringToInt(typeStr string) int {
	switch typeStr {
	case "None":
		return 1
	case "Device":
		return 0
	default:
		return 0 // Default to Device
	}
}

// constructResource constructs a JSON request body for the Autopatch Groups API
func constructResource(ctx context.Context, data *AutopatchGroupsResourceModel) (map[string]any, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from Terraform configuration", ResourceName))

	requestBody := make(map[string]any)

	// Required fields
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		requestBody["name"] = data.Name.ValueString()
	}

	// Optional fields
	if !data.Description.IsNull() && !data.Description.IsUnknown() {
		requestBody["description"] = data.Description.ValueString()
	}

	// Global User Managed AAD Groups
	if !data.GlobalUserManagedAadGroups.IsNull() && !data.GlobalUserManagedAadGroups.IsUnknown() {
		var globalGroups []GlobalUserManagedAadGroup
		data.GlobalUserManagedAadGroups.ElementsAs(ctx, &globalGroups, false)

		globalGroupsAPI := make([]map[string]any, 0, len(globalGroups))
		for _, group := range globalGroups {
			groupMap := make(map[string]any)
			if !group.Id.IsNull() && !group.Id.IsUnknown() {
				groupMap["id"] = group.Id.ValueString()
			}
			if !group.Type.IsNull() && !group.Type.IsUnknown() {
				// Convert string to integer for API (POST uses integers)
				groupMap["type"] = mapTypeStringToInt(group.Type.ValueString())
			}
			globalGroupsAPI = append(globalGroupsAPI, groupMap)
		}
		requestBody["globalUserManagedAadGroups"] = globalGroupsAPI
	} else {
		requestBody["globalUserManagedAadGroups"] = []any{}
	}

	// Deployment Groups
	if !data.DeploymentGroups.IsNull() && !data.DeploymentGroups.IsUnknown() {
		var deploymentGroups []DeploymentGroup
		data.DeploymentGroups.ElementsAs(ctx, &deploymentGroups, false)

		deploymentGroupsAPI := make([]map[string]any, 0, len(deploymentGroups))
		for _, group := range deploymentGroups {
			groupMap := make(map[string]any)

			if !group.AadId.IsNull() && !group.AadId.IsUnknown() {
				groupMap["aadId"] = group.AadId.ValueString()
			}
			if !group.Name.IsNull() && !group.Name.IsUnknown() {
				groupMap["name"] = group.Name.ValueString()
			}
			if !group.Distribution.IsNull() && !group.Distribution.IsUnknown() {
				groupMap["distribution"] = group.Distribution.ValueInt32()
			}
			if !group.FailedPrerequisiteCheckCount.IsNull() && !group.FailedPrerequisiteCheckCount.IsUnknown() {
				groupMap["failedPreRequisiteCheckCount"] = group.FailedPrerequisiteCheckCount.ValueInt32()
			}

			// User Managed AAD Groups within deployment group
			if !group.UserManagedAadGroups.IsNull() && !group.UserManagedAadGroups.IsUnknown() {
				var userGroups []UserManagedAadGroup
				group.UserManagedAadGroups.ElementsAs(ctx, &userGroups, false)

				userGroupsAPI := make([]map[string]any, 0, len(userGroups))
				for _, userGroup := range userGroups {
					userGroupMap := make(map[string]any)
					if !userGroup.Id.IsNull() && !userGroup.Id.IsUnknown() {
						userGroupMap["id"] = userGroup.Id.ValueString()
					}
					if !userGroup.Name.IsNull() && !userGroup.Name.IsUnknown() {
						userGroupMap["name"] = userGroup.Name.ValueString()
					}
					if !userGroup.Type.IsNull() && !userGroup.Type.IsUnknown() {
						// Convert string to integer for API (POST uses integers)
						userGroupMap["type"] = mapTypeStringToInt(userGroup.Type.ValueString())
					}
					userGroupsAPI = append(userGroupsAPI, userGroupMap)
				}
				groupMap["userManagedAadGroups"] = userGroupsAPI
			} else {
				groupMap["userManagedAadGroups"] = []any{}
			}

			// Deployment Group Policy Settings
			if group.DeploymentGroupPolicySettings != nil {
				policyMap, err := constructDeploymentGroupPolicySettings(ctx, group.DeploymentGroupPolicySettings)
				if err != nil {
					return nil, fmt.Errorf("error constructing deployment group policy settings: %w", err)
				}
				groupMap["deploymentGroupPolicySettings"] = policyMap
			}

			deploymentGroupsAPI = append(deploymentGroupsAPI, groupMap)
		}
		requestBody["deploymentGroups"] = deploymentGroupsAPI
	} else {
		requestBody["deploymentGroups"] = []any{}
	}

	// Default fields based on API example
	requestBody["windowsUpdateSettings"] = []any{}
	requestBody["status"] = "Unknown"
	requestBody["type"] = "Unknown"
	requestBody["distributionType"] = "Unknown"
	requestBody["driverUpdateSettings"] = []any{}

	if !data.EnableDriverUpdate.IsNull() && !data.EnableDriverUpdate.IsUnknown() {
		requestBody["enableDriverUpdate"] = data.EnableDriverUpdate.ValueBool()
	} else {
		requestBody["enableDriverUpdate"] = true // Default from API example
	}

	// Scope Tags - convert from string set to int array for API
	if !data.ScopeTags.IsNull() && !data.ScopeTags.IsUnknown() {
		var scopeTagStrings []string
		convert.FrameworkToGraphStringSet(ctx, data.ScopeTags, func(tags []string) {
			scopeTagStrings = tags
		})

		scopeTagsAPI := make([]int, 0, len(scopeTagStrings))
		for _, tagStr := range scopeTagStrings {
			if tagInt, err := strconv.Atoi(tagStr); err == nil {
				scopeTagsAPI = append(scopeTagsAPI, tagInt)
			}
		}
		requestBody["scopeTags"] = scopeTagsAPI
	} else {
		requestBody["scopeTags"] = []int{0} // Default
	}

	// Enabled Content Types
	if !data.EnabledContentTypes.IsNull() && !data.EnabledContentTypes.IsUnknown() {
		requestBody["enabledContentTypes"] = data.EnabledContentTypes.ValueInt32()
	} else {
		requestBody["enabledContentTypes"] = 31 // Default from API example
	}

	jsonData, err := json.MarshalIndent(requestBody, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	tflog.Debug(ctx, fmt.Sprintf("Constructed %s JSON payload", ResourceName), map[string]any{
		"json_payload": string(jsonData),
	})

	return requestBody, nil
}

// constructDeploymentGroupPolicySettings constructs the deployment group policy settings
func constructDeploymentGroupPolicySettings(ctx context.Context, settings *DeploymentGroupPolicySettings) (map[string]any, error) {
	policyMap := make(map[string]any)

	convert.FrameworkToGraphString(settings.AadGroupName, func(val *string) {
		if val != nil {
			policyMap["aadGroupName"] = *val
		}
	})

	convert.FrameworkToGraphBool(settings.IsUpdateSettingsModified, func(val *bool) {
		if val != nil {
			policyMap["isUpdateSettingsModified"] = *val
		}
	})

	// Device Configuration Settings
	if settings.DeviceConfigurationSetting != nil {
		deviceConfigMap := make(map[string]any)

		convert.FrameworkToGraphString(settings.DeviceConfigurationSetting.PolicyId, func(val *string) {
			if val != nil {
				deviceConfigMap["policyId"] = *val
			}
		})

		convert.FrameworkToGraphString(settings.DeviceConfigurationSetting.UpdateBehavior, func(val *string) {
			if val != nil {
				deviceConfigMap["updateBehavior"] = *val
			}
		})

		convert.FrameworkToGraphString(settings.DeviceConfigurationSetting.NotificationSetting, func(val *string) {
			if val != nil {
				deviceConfigMap["notificationSetting"] = *val
			}
		})

		// Quality Deployment Settings
		if settings.DeviceConfigurationSetting.QualityDeploymentSettings != nil {
			qualityMap := make(map[string]any)
			convert.FrameworkToGraphInt32(settings.DeviceConfigurationSetting.QualityDeploymentSettings.Deadline, func(val *int32) {
				if val != nil {
					qualityMap["deadline"] = *val
				}
			})
			convert.FrameworkToGraphInt32(settings.DeviceConfigurationSetting.QualityDeploymentSettings.Deferral, func(val *int32) {
				if val != nil {
					qualityMap["deferral"] = *val
				}
			})
			convert.FrameworkToGraphInt32(settings.DeviceConfigurationSetting.QualityDeploymentSettings.GracePeriod, func(val *int32) {
				if val != nil {
					qualityMap["gracePeriod"] = *val
				}
			})
			deviceConfigMap["qualityDeploymentSettings"] = qualityMap
		}

		// Feature Deployment Settings
		if settings.DeviceConfigurationSetting.FeatureDeploymentSettings != nil {
			featureMap := make(map[string]any)
			convert.FrameworkToGraphInt32(settings.DeviceConfigurationSetting.FeatureDeploymentSettings.Deadline, func(val *int32) {
				if val != nil {
					featureMap["deadline"] = *val
				}
			})
			convert.FrameworkToGraphInt32(settings.DeviceConfigurationSetting.FeatureDeploymentSettings.Deferral, func(val *int32) {
				if val != nil {
					featureMap["deferral"] = *val
				}
			})
			deviceConfigMap["featureDeploymentSettings"] = featureMap
		}

		// Add required null fields from API example
		deviceConfigMap["updateFrequencyUI"] = nil
		deviceConfigMap["installDays"] = nil
		deviceConfigMap["installTime"] = nil
		deviceConfigMap["activeHourEndTime"] = nil
		deviceConfigMap["activeHourStartTime"] = nil

		policyMap["deviceConfigurationSetting"] = deviceConfigMap
	}

	// DNF Update Cloud Setting
	if settings.DnfUpdateCloudSetting != nil {
		dnfMap := make(map[string]any)
		convert.FrameworkToGraphString(settings.DnfUpdateCloudSetting.PolicyId, func(val *string) {
			if val != nil {
				dnfMap["policyId"] = *val
			}
		})
		convert.FrameworkToGraphString(settings.DnfUpdateCloudSetting.ApprovalType, func(val *string) {
			if val != nil {
				dnfMap["approvalType"] = *val
			}
		})
		convert.FrameworkToGraphInt32(settings.DnfUpdateCloudSetting.DeploymentDeferralInDays, func(val *int32) {
			if val != nil {
				dnfMap["deploymentDeferralInDays"] = *val
			} else {
				dnfMap["deploymentDeferralInDays"] = nil
			}
		})
		policyMap["dnfUpdateCloudSetting"] = dnfMap
	}

	// Office DCv2 Setting
	if settings.OfficeDCv2Setting != nil {
		officeMap := make(map[string]any)
		convert.FrameworkToGraphString(settings.OfficeDCv2Setting.PolicyId, func(val *string) {
			if val != nil {
				officeMap["policyId"] = *val
			}
		})
		convert.FrameworkToGraphInt32(settings.OfficeDCv2Setting.Deadline, func(val *int32) {
			if val != nil {
				officeMap["deadline"] = *val
			}
		})
		convert.FrameworkToGraphInt32(settings.OfficeDCv2Setting.Deferral, func(val *int32) {
			if val != nil {
				officeMap["deferral"] = *val
			}
		})
		convert.FrameworkToGraphBool(settings.OfficeDCv2Setting.HideUpdateNotifications, func(val *bool) {
			if val != nil {
				officeMap["hideUpdateNotifications"] = *val
			}
		})
		convert.FrameworkToGraphString(settings.OfficeDCv2Setting.TargetChannel, func(val *string) {
			if val != nil {
				officeMap["targetChannel"] = *val
			}
		})
		convert.FrameworkToGraphBool(settings.OfficeDCv2Setting.EnableAutomaticUpdate, func(val *bool) {
			if val != nil {
				officeMap["enableAutomaticUpdate"] = *val
			}
		})
		convert.FrameworkToGraphBool(settings.OfficeDCv2Setting.HideEnableDisableUpdate, func(val *bool) {
			if val != nil {
				officeMap["hideEnableDisableUpdate"] = *val
			}
		})
		convert.FrameworkToGraphBool(settings.OfficeDCv2Setting.EnableOfficeMgmt, func(val *bool) {
			if val != nil {
				officeMap["enableOfficeMgmt"] = *val
			}
		})
		convert.FrameworkToGraphString(settings.OfficeDCv2Setting.UpdatePath, func(val *string) {
			if val != nil {
				officeMap["updatePath"] = *val
			}
		})
		policyMap["officeDCv2Setting"] = officeMap
	}

	// Edge DCv2 Setting
	if settings.EdgeDCv2Setting != nil {
		edgeMap := make(map[string]any)
		convert.FrameworkToGraphString(settings.EdgeDCv2Setting.PolicyId, func(val *string) {
			if val != nil {
				edgeMap["policyId"] = *val
			}
		})
		convert.FrameworkToGraphString(settings.EdgeDCv2Setting.TargetChannel, func(val *string) {
			if val != nil {
				edgeMap["targetChannel"] = *val
			}
		})
		policyMap["edgeDCv2Setting"] = edgeMap
	}

	// Feature Update Anchor Cloud Setting
	if settings.FeatureUpdateAnchorCloudSetting != nil {
		featureAnchorMap := make(map[string]any)
		convert.FrameworkToGraphString(settings.FeatureUpdateAnchorCloudSetting.TargetOSVersion, func(val *string) {
			if val != nil {
				featureAnchorMap["targetOSVersion"] = *val
			}
		})
		convert.FrameworkToGraphBool(settings.FeatureUpdateAnchorCloudSetting.InstallLatestWindows10OnWindows11IneligibleDevice, func(val *bool) {
			if val != nil {
				featureAnchorMap["installLatestWindows10OnWindows11IneligibleDevice"] = *val
			}
		})
		convert.FrameworkToGraphString(settings.FeatureUpdateAnchorCloudSetting.PolicyId, func(val *string) {
			if val != nil {
				featureAnchorMap["policyId"] = *val
			}
		})
		policyMap["featureUpdateAnchorCloudSetting"] = featureAnchorMap
	}

	return policyMap, nil
}
