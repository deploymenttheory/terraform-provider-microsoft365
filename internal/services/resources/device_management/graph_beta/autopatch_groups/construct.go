package graphBetaAutopatchGroups

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// constructResource constructs a JSON request body for the Autopatch Groups API
func constructResource(ctx context.Context, data *AutopatchGroupsResourceModel) (map[string]interface{}, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from Terraform configuration", ResourceName))

	requestBody := make(map[string]interface{})

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
		
		globalGroupsAPI := make([]map[string]interface{}, 0, len(globalGroups))
		for _, group := range globalGroups {
			groupMap := make(map[string]interface{})
			if !group.Id.IsNull() && !group.Id.IsUnknown() {
				groupMap["id"] = group.Id.ValueString()
			}
			if !group.Type.IsNull() && !group.Type.IsUnknown() {
				groupMap["type"] = group.Type.ValueString()
			}
			globalGroupsAPI = append(globalGroupsAPI, groupMap)
		}
		requestBody["globalUserManagedAadGroups"] = globalGroupsAPI
	} else {
		requestBody["globalUserManagedAadGroups"] = []interface{}{}
	}

	// Deployment Groups
	if !data.DeploymentGroups.IsNull() && !data.DeploymentGroups.IsUnknown() {
		var deploymentGroups []DeploymentGroup
		data.DeploymentGroups.ElementsAs(ctx, &deploymentGroups, false)
		
		deploymentGroupsAPI := make([]map[string]interface{}, 0, len(deploymentGroups))
		for _, group := range deploymentGroups {
			groupMap := make(map[string]interface{})
			
			if !group.AadId.IsNull() && !group.AadId.IsUnknown() {
				groupMap["aadId"] = group.AadId.ValueString()
			}
			if !group.Name.IsNull() && !group.Name.IsUnknown() {
				groupMap["name"] = group.Name.ValueString()
			}
			if !group.Distribution.IsNull() && !group.Distribution.IsUnknown() {
				groupMap["distribution"] = group.Distribution.ValueInt64()
			}
			if !group.FailedPrerequisiteCheckCount.IsNull() && !group.FailedPrerequisiteCheckCount.IsUnknown() {
				groupMap["failedPreRequisiteCheckCount"] = group.FailedPrerequisiteCheckCount.ValueInt64()
			}

			// User Managed AAD Groups within deployment group
			if !group.UserManagedAadGroups.IsNull() && !group.UserManagedAadGroups.IsUnknown() {
				var userGroups []UserManagedAadGroup
				group.UserManagedAadGroups.ElementsAs(ctx, &userGroups, false)
				
				userGroupsAPI := make([]map[string]interface{}, 0, len(userGroups))
				for _, userGroup := range userGroups {
					userGroupMap := make(map[string]interface{})
					if !userGroup.Id.IsNull() && !userGroup.Id.IsUnknown() {
						userGroupMap["id"] = userGroup.Id.ValueString()
					}
					if !userGroup.Name.IsNull() && !userGroup.Name.IsUnknown() {
						userGroupMap["name"] = userGroup.Name.ValueString()
					}
					if !userGroup.Type.IsNull() && !userGroup.Type.IsUnknown() {
						userGroupMap["type"] = userGroup.Type.ValueInt64()
					}
					userGroupsAPI = append(userGroupsAPI, userGroupMap)
				}
				groupMap["userManagedAadGroups"] = userGroupsAPI
			} else {
				groupMap["userManagedAadGroups"] = []interface{}{}
			}

			// Deployment Group Policy Settings
			if group.DeploymentGroupPolicySettings != nil {
				policyMap := make(map[string]interface{})
				
				if !group.DeploymentGroupPolicySettings.AadGroupName.IsNull() && !group.DeploymentGroupPolicySettings.AadGroupName.IsUnknown() {
					policyMap["aadGroupName"] = group.DeploymentGroupPolicySettings.AadGroupName.ValueString()
				}
				if !group.DeploymentGroupPolicySettings.IsUpdateSettingsModified.IsNull() && !group.DeploymentGroupPolicySettings.IsUpdateSettingsModified.IsUnknown() {
					policyMap["isUpdateSettingsModified"] = group.DeploymentGroupPolicySettings.IsUpdateSettingsModified.ValueBool()
				}

				// Device Configuration Settings
				if group.DeploymentGroupPolicySettings.DeviceConfigurationSetting != nil {
					deviceConfigMap := make(map[string]interface{})
					
					if !group.DeploymentGroupPolicySettings.DeviceConfigurationSetting.PolicyId.IsNull() && !group.DeploymentGroupPolicySettings.DeviceConfigurationSetting.PolicyId.IsUnknown() {
						deviceConfigMap["policyId"] = group.DeploymentGroupPolicySettings.DeviceConfigurationSetting.PolicyId.ValueString()
					}
					if !group.DeploymentGroupPolicySettings.DeviceConfigurationSetting.UpdateBehavior.IsNull() && !group.DeploymentGroupPolicySettings.DeviceConfigurationSetting.UpdateBehavior.IsUnknown() {
						deviceConfigMap["updateBehavior"] = group.DeploymentGroupPolicySettings.DeviceConfigurationSetting.UpdateBehavior.ValueString()
					}
					if !group.DeploymentGroupPolicySettings.DeviceConfigurationSetting.NotificationSetting.IsNull() && !group.DeploymentGroupPolicySettings.DeviceConfigurationSetting.NotificationSetting.IsUnknown() {
						deviceConfigMap["notificationSetting"] = group.DeploymentGroupPolicySettings.DeviceConfigurationSetting.NotificationSetting.ValueString()
					}

					// Quality Deployment Settings
					if group.DeploymentGroupPolicySettings.DeviceConfigurationSetting.QualityDeploymentSettings != nil {
						qualityMap := make(map[string]interface{})
						if !group.DeploymentGroupPolicySettings.DeviceConfigurationSetting.QualityDeploymentSettings.Deadline.IsNull() && !group.DeploymentGroupPolicySettings.DeviceConfigurationSetting.QualityDeploymentSettings.Deadline.IsUnknown() {
							qualityMap["deadline"] = group.DeploymentGroupPolicySettings.DeviceConfigurationSetting.QualityDeploymentSettings.Deadline.ValueInt64()
						}
						if !group.DeploymentGroupPolicySettings.DeviceConfigurationSetting.QualityDeploymentSettings.Deferral.IsNull() && !group.DeploymentGroupPolicySettings.DeviceConfigurationSetting.QualityDeploymentSettings.Deferral.IsUnknown() {
							qualityMap["deferral"] = group.DeploymentGroupPolicySettings.DeviceConfigurationSetting.QualityDeploymentSettings.Deferral.ValueInt64()
						}
						if !group.DeploymentGroupPolicySettings.DeviceConfigurationSetting.QualityDeploymentSettings.GracePeriod.IsNull() && !group.DeploymentGroupPolicySettings.DeviceConfigurationSetting.QualityDeploymentSettings.GracePeriod.IsUnknown() {
							qualityMap["gracePeriod"] = group.DeploymentGroupPolicySettings.DeviceConfigurationSetting.QualityDeploymentSettings.GracePeriod.ValueInt64()
						}
						deviceConfigMap["qualityDeploymentSettings"] = qualityMap
					}

					// Feature Deployment Settings
					if group.DeploymentGroupPolicySettings.DeviceConfigurationSetting.FeatureDeploymentSettings != nil {
						featureMap := make(map[string]interface{})
						if !group.DeploymentGroupPolicySettings.DeviceConfigurationSetting.FeatureDeploymentSettings.Deadline.IsNull() && !group.DeploymentGroupPolicySettings.DeviceConfigurationSetting.FeatureDeploymentSettings.Deadline.IsUnknown() {
							featureMap["deadline"] = group.DeploymentGroupPolicySettings.DeviceConfigurationSetting.FeatureDeploymentSettings.Deadline.ValueInt64()
						}
						if !group.DeploymentGroupPolicySettings.DeviceConfigurationSetting.FeatureDeploymentSettings.Deferral.IsNull() && !group.DeploymentGroupPolicySettings.DeviceConfigurationSetting.FeatureDeploymentSettings.Deferral.IsUnknown() {
							featureMap["deferral"] = group.DeploymentGroupPolicySettings.DeviceConfigurationSetting.FeatureDeploymentSettings.Deferral.ValueInt64()
						}
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
				
				groupMap["deploymentGroupPolicySettings"] = policyMap
			}
			
			deploymentGroupsAPI = append(deploymentGroupsAPI, groupMap)
		}
		requestBody["deploymentGroups"] = deploymentGroupsAPI
	} else {
		requestBody["deploymentGroups"] = []interface{}{}
	}

	// Default fields based on API example
	requestBody["windowsUpdateSettings"] = []interface{}{}
	requestBody["status"] = "Unknown"
	requestBody["type"] = "Unknown"
	requestBody["distributionType"] = "Unknown"
	requestBody["driverUpdateSettings"] = []interface{}{}

	// Optional boolean fields
	if !data.EnableDriverUpdate.IsNull() && !data.EnableDriverUpdate.IsUnknown() {
		requestBody["enableDriverUpdate"] = data.EnableDriverUpdate.ValueBool()
	} else {
		requestBody["enableDriverUpdate"] = true // Default from API example
	}

	// Scope Tags
	if !data.ScopeTags.IsNull() && !data.ScopeTags.IsUnknown() {
		var scopeTags []types.Int64
		data.ScopeTags.ElementsAs(ctx, &scopeTags, false)
		
		scopeTagsAPI := make([]int64, 0, len(scopeTags))
		for _, tag := range scopeTags {
			if !tag.IsNull() && !tag.IsUnknown() {
				scopeTagsAPI = append(scopeTagsAPI, tag.ValueInt64())
			}
		}
		requestBody["scopeTags"] = scopeTagsAPI
	} else {
		requestBody["scopeTags"] = []int64{0} // Default
	}

	// Enabled Content Types
	if !data.EnabledContentTypes.IsNull() && !data.EnabledContentTypes.IsUnknown() {
		requestBody["enabledContentTypes"] = data.EnabledContentTypes.ValueInt64()
	} else {
		requestBody["enabledContentTypes"] = 31 // Default from API example
	}

	jsonData, err := json.MarshalIndent(requestBody, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	tflog.Debug(ctx, fmt.Sprintf("Constructed %s JSON payload", ResourceName), map[string]interface{}{
		"json_payload": string(jsonData),
	})

	return requestBody, nil
}