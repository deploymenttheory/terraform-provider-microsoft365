package graphBetaAutopatchGroups

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// MapRemoteStateToTerraform maps the remote state from the API to the Terraform model
func MapRemoteStateToTerraform(ctx context.Context, data *AutopatchGroupsResourceModel, autopatchGroup map[string]any) {
	tflog.Debug(ctx, fmt.Sprintf("Mapping %s resource from API to Terraform state", ResourceName))

	// Map top-level fields using convert helpers
	data.ID = convert.MapToFrameworkString(autopatchGroup, "id")
	data.Name = convert.MapToFrameworkString(autopatchGroup, "name")
	data.Description = convert.MapToFrameworkString(autopatchGroup, "description")
	data.TenantId = convert.MapToFrameworkString(autopatchGroup, "tenantId")
	data.Type = convert.MapToFrameworkString(autopatchGroup, "type")
	data.Status = convert.MapToFrameworkString(autopatchGroup, "status")
	data.DistributionType = convert.MapToFrameworkString(autopatchGroup, "distributionType")
	data.FlowId = convert.MapToFrameworkString(autopatchGroup, "flowId")
	data.FlowType = convert.MapToFrameworkString(autopatchGroup, "flowType")
	data.FlowStatus = convert.MapToFrameworkString(autopatchGroup, "flowStatus")
	data.UmbrellaGroupId = convert.MapToFrameworkString(autopatchGroup, "umbrellaGroupId")
	data.IsLockedByPolicy = convert.MapToFrameworkBool(autopatchGroup, "isLockedByPolicy")
	data.ReadOnly = convert.MapToFrameworkBool(autopatchGroup, "readOnly")
	data.UserHasAllScopeTag = convert.MapToFrameworkBool(autopatchGroup, "userHasAllScopeTag")
	data.EnableDriverUpdate = convert.MapToFrameworkBool(autopatchGroup, "enableDriverUpdate")
	data.NumberOfRegisteredDevices = convert.MapToFrameworkInt32(autopatchGroup, "numberOfRegisteredDevices")
	data.EnabledContentTypes = convert.MapToFrameworkInt32(autopatchGroup, "enabledContentTypes")

	// Map scope tags (numbers to strings) using convert helper
	data.ScopeTags = convert.MapToFrameworkStringSetFromNumbers(ctx, autopatchGroup, "scopeTags")

	// Map global_user_managed_aad_groups using types.SetValueFrom
	if globalGroupsRaw, ok := autopatchGroup["globalUserManagedAadGroups"].([]any); ok {
		globalGroups := make([]GlobalUserManagedAadGroup, 0, len(globalGroupsRaw))
		for _, groupRaw := range globalGroupsRaw {
			if groupMap, ok := groupRaw.(map[string]any); ok {
				globalGroups = append(globalGroups, GlobalUserManagedAadGroup{
					Id:   convert.MapToFrameworkString(groupMap, "id"),
					Type: convert.MapToFrameworkString(groupMap, "type"),
				})
			}
		}
		globalGroupsSet, diags := types.SetValueFrom(ctx, types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":   types.StringType,
				"type": types.StringType,
			},
		}, globalGroups)
		if diags.HasError() {
			tflog.Error(ctx, "Failed to convert global user managed AAD groups to set")
			data.GlobalUserManagedAadGroups = types.SetNull(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"id":   types.StringType,
					"type": types.StringType,
				},
			})
		} else {
			data.GlobalUserManagedAadGroups = globalGroupsSet
		}
	} else {
		emptySet, _ := types.SetValueFrom(ctx, types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":   types.StringType,
				"type": types.StringType,
			},
		}, []GlobalUserManagedAadGroup{})
		data.GlobalUserManagedAadGroups = emptySet
	}

	// Map deployment_groups as attr.Value list for complex nested structures
	if deploymentGroupsRaw, ok := autopatchGroup["deploymentGroups"].([]any); ok {
		deploymentGroupsList := mapDeploymentGroupsList(ctx, deploymentGroupsRaw)
		data.DeploymentGroups = deploymentGroupsList
	} else {
		data.DeploymentGroups = types.ListValueMust(
			getDeploymentGroupObjectType(),
			[]attr.Value{},
		)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping %s resource from API to Terraform state", ResourceName))
}

// mapDeploymentGroupsList maps deployment groups from API to types.List
func mapDeploymentGroupsList(ctx context.Context, deploymentGroupsRaw []any) types.List {
	if len(deploymentGroupsRaw) == 0 {
		return types.ListValueMust(getDeploymentGroupObjectType(), []attr.Value{})
	}

	deploymentGroupsValues := make([]attr.Value, 0, len(deploymentGroupsRaw))

	for _, groupRaw := range deploymentGroupsRaw {
		groupMap, ok := groupRaw.(map[string]any)
		if !ok {
			continue
		}

		deploymentGroupAttrs := map[string]attr.Value{
			"aad_id":                           convert.MapToFrameworkString(groupMap, "aadId"),
			"name":                             convert.MapToFrameworkString(groupMap, "name"),
			"distribution":                     convert.MapToFrameworkInt32(groupMap, "distribution"),
			"failed_prerequisite_check_count":  convert.MapToFrameworkInt32(groupMap, "failedPreRequisiteCheckCount"),
			"user_managed_aad_groups":          mapUserManagedAadGroups(ctx, groupMap),
			"deployment_group_policy_settings": mapPolicySettingsToObject(ctx, groupMap),
		}

		deploymentGroupsValues = append(deploymentGroupsValues, types.ObjectValueMust(
			getDeploymentGroupObjectType().AttrTypes,
			deploymentGroupAttrs,
		))
	}

	return types.ListValueMust(getDeploymentGroupObjectType(), deploymentGroupsValues)
}

// mapUserManagedAadGroups maps user managed AAD groups to types.Set
func mapUserManagedAadGroups(ctx context.Context, groupMap map[string]any) types.Set {
	userGroupsRaw, ok := groupMap["userManagedAadGroups"].([]any)
	if !ok || len(userGroupsRaw) == 0 {
		return types.SetNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":   types.StringType,
				"name": types.StringType,
				"type": types.StringType,
			},
		})
	}

	userGroups := make([]UserManagedAadGroup, 0, len(userGroupsRaw))
	for _, ugRaw := range userGroupsRaw {
		if ugMap, ok := ugRaw.(map[string]any); ok {
			userGroups = append(userGroups, UserManagedAadGroup{
				Id:   convert.MapToFrameworkString(ugMap, "id"),
				Name: convert.MapToFrameworkString(ugMap, "name"),
				Type: convert.MapToFrameworkString(ugMap, "type"),
			})
		}
	}

	userGroupsSet, diags := types.SetValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":   types.StringType,
			"name": types.StringType,
			"type": types.StringType,
		},
	}, userGroups)

	if diags.HasError() {
		tflog.Error(ctx, "Failed to convert user managed AAD groups to set")
		return types.SetNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":   types.StringType,
				"name": types.StringType,
				"type": types.StringType,
			},
		})
	}

	return userGroupsSet
}

// mapPolicySettingsToObject maps policy settings to types.Object
func mapPolicySettingsToObject(ctx context.Context, groupMap map[string]any) types.Object {
	policyRaw, ok := groupMap["deploymentGroupPolicySettings"].(map[string]any)
	if !ok {
		return types.ObjectNull(getPolicySettingsAttrTypes())
	}

	policyAttrs := map[string]attr.Value{
		"aad_group_name":                      convert.MapToFrameworkString(policyRaw, "aadGroupName"),
		"is_update_settings_modified":         convert.MapToFrameworkBool(policyRaw, "isUpdateSettingsModified"),
		"device_configuration_setting":        mapDeviceConfigurationSetting(ctx, policyRaw),
		"dnf_update_cloud_setting":            mapDnfUpdateCloudSetting(ctx, policyRaw),
		"office_dcv2_setting":                 mapOfficeDCv2Setting(ctx, policyRaw),
		"edge_dcv2_setting":                   mapEdgeDCv2Setting(ctx, policyRaw),
		"feature_update_anchor_cloud_setting": mapFeatureUpdateAnchorCloudSetting(ctx, policyRaw),
	}

	return types.ObjectValueMust(getPolicySettingsAttrTypes(), policyAttrs)
}

// mapDeviceConfigurationSetting maps device configuration setting
func mapDeviceConfigurationSetting(ctx context.Context, policyRaw map[string]any) types.Object {
	deviceConfigRaw, ok := policyRaw["deviceConfigurationSetting"].(map[string]any)
	if !ok {
		return types.ObjectNull(map[string]attr.Type{
			"policy_id":                   types.StringType,
			"update_behavior":             types.StringType,
			"notification_setting":        types.StringType,
			"quality_deployment_settings": getQualityDeploymentSettingsObjectType(),
			"feature_deployment_settings": getFeatureDeploymentSettingsObjectType(),
		})
	}

	deviceConfigAttrs := map[string]attr.Value{
		"policy_id":                   convert.MapToFrameworkString(deviceConfigRaw, "policyId"),
		"update_behavior":             convert.MapToFrameworkString(deviceConfigRaw, "updateBehavior"),
		"notification_setting":        convert.MapToFrameworkString(deviceConfigRaw, "notificationSetting"),
		"quality_deployment_settings": mapQualityDeploymentSettings(ctx, deviceConfigRaw),
		"feature_deployment_settings": mapFeatureDeploymentSettings(ctx, deviceConfigRaw),
	}

	return types.ObjectValueMust(map[string]attr.Type{
		"policy_id":                   types.StringType,
		"update_behavior":             types.StringType,
		"notification_setting":        types.StringType,
		"quality_deployment_settings": getQualityDeploymentSettingsObjectType(),
		"feature_deployment_settings": getFeatureDeploymentSettingsObjectType(),
	}, deviceConfigAttrs)
}

// mapQualityDeploymentSettings maps quality deployment settings
func mapQualityDeploymentSettings(ctx context.Context, deviceConfigRaw map[string]any) types.Object {
	qualityRaw, ok := deviceConfigRaw["qualityDeploymentSettings"].(map[string]any)
	if !ok {
		return types.ObjectNull(getQualityDeploymentSettingsObjectType().AttrTypes)
	}

	qualityAttrs := map[string]attr.Value{
		"deadline":     convert.MapToFrameworkInt32(qualityRaw, "deadline"),
		"deferral":     convert.MapToFrameworkInt32(qualityRaw, "deferral"),
		"grace_period": convert.MapToFrameworkInt32(qualityRaw, "gracePeriod"),
	}

	return types.ObjectValueMust(getQualityDeploymentSettingsObjectType().AttrTypes, qualityAttrs)
}

// mapFeatureDeploymentSettings maps feature deployment settings
func mapFeatureDeploymentSettings(ctx context.Context, deviceConfigRaw map[string]any) types.Object {
	featureRaw, ok := deviceConfigRaw["featureDeploymentSettings"].(map[string]any)
	if !ok {
		return types.ObjectNull(getFeatureDeploymentSettingsObjectType().AttrTypes)
	}

	featureAttrs := map[string]attr.Value{
		"deadline": convert.MapToFrameworkInt32(featureRaw, "deadline"),
		"deferral": convert.MapToFrameworkInt32(featureRaw, "deferral"),
	}

	return types.ObjectValueMust(getFeatureDeploymentSettingsObjectType().AttrTypes, featureAttrs)
}

// mapDnfUpdateCloudSetting maps DNF update cloud setting
func mapDnfUpdateCloudSetting(ctx context.Context, policyRaw map[string]any) types.Object {
	dnfRaw, ok := policyRaw["dnfUpdateCloudSetting"].(map[string]any)
	if !ok {
		return types.ObjectNull(map[string]attr.Type{
			"policy_id":                   types.StringType,
			"approval_type":               types.StringType,
			"deployment_deferral_in_days": types.Int32Type,
		})
	}

	dnfAttrs := map[string]attr.Value{
		"policy_id":                   convert.MapToFrameworkString(dnfRaw, "policyId"),
		"approval_type":               convert.MapToFrameworkString(dnfRaw, "approvalType"),
		"deployment_deferral_in_days": convert.MapToFrameworkInt32(dnfRaw, "deploymentDeferralInDays"),
	}

	return types.ObjectValueMust(map[string]attr.Type{
		"policy_id":                   types.StringType,
		"approval_type":               types.StringType,
		"deployment_deferral_in_days": types.Int32Type,
	}, dnfAttrs)
}

// mapOfficeDCv2Setting maps Office DCv2 setting
func mapOfficeDCv2Setting(ctx context.Context, policyRaw map[string]any) types.Object {
	officeRaw, ok := policyRaw["officeDCv2Setting"].(map[string]any)
	if !ok {
		return types.ObjectNull(map[string]attr.Type{
			"policy_id":                  types.StringType,
			"deadline":                   types.Int32Type,
			"deferral":                   types.Int32Type,
			"hide_update_notifications":  types.BoolType,
			"target_channel":             types.StringType,
			"enable_automatic_update":    types.BoolType,
			"hide_enable_disable_update": types.BoolType,
			"enable_office_mgmt":         types.BoolType,
			"update_path":                types.StringType,
		})
	}

	officeAttrs := map[string]attr.Value{
		"policy_id":                  convert.MapToFrameworkString(officeRaw, "policyId"),
		"deadline":                   convert.MapToFrameworkInt32(officeRaw, "deadline"),
		"deferral":                   convert.MapToFrameworkInt32(officeRaw, "deferral"),
		"hide_update_notifications":  convert.MapToFrameworkBool(officeRaw, "hideUpdateNotifications"),
		"target_channel":             convert.MapToFrameworkString(officeRaw, "targetChannel"),
		"enable_automatic_update":    convert.MapToFrameworkBool(officeRaw, "enableAutomaticUpdate"),
		"hide_enable_disable_update": convert.MapToFrameworkBool(officeRaw, "hideEnableDisableUpdate"),
		"enable_office_mgmt":         convert.MapToFrameworkBool(officeRaw, "enableOfficeMgmt"),
		"update_path":                convert.MapToFrameworkString(officeRaw, "updatePath"),
	}

	return types.ObjectValueMust(map[string]attr.Type{
		"policy_id":                  types.StringType,
		"deadline":                   types.Int32Type,
		"deferral":                   types.Int32Type,
		"hide_update_notifications":  types.BoolType,
		"target_channel":             types.StringType,
		"enable_automatic_update":    types.BoolType,
		"hide_enable_disable_update": types.BoolType,
		"enable_office_mgmt":         types.BoolType,
		"update_path":                types.StringType,
	}, officeAttrs)
}

// mapEdgeDCv2Setting maps Edge DCv2 setting
func mapEdgeDCv2Setting(ctx context.Context, policyRaw map[string]any) types.Object {
	edgeRaw, ok := policyRaw["edgeDCv2Setting"].(map[string]any)
	if !ok {
		return types.ObjectNull(map[string]attr.Type{
			"policy_id":      types.StringType,
			"target_channel": types.StringType,
		})
	}

	edgeAttrs := map[string]attr.Value{
		"policy_id":      convert.MapToFrameworkString(edgeRaw, "policyId"),
		"target_channel": convert.MapToFrameworkString(edgeRaw, "targetChannel"),
	}

	return types.ObjectValueMust(map[string]attr.Type{
		"policy_id":      types.StringType,
		"target_channel": types.StringType,
	}, edgeAttrs)
}

// mapFeatureUpdateAnchorCloudSetting maps feature update anchor cloud setting
func mapFeatureUpdateAnchorCloudSetting(ctx context.Context, policyRaw map[string]any) types.Object {
	featureUpdateRaw, ok := policyRaw["featureUpdateAnchorCloudSetting"].(map[string]any)
	if !ok {
		return types.ObjectNull(map[string]attr.Type{
			"target_os_version": types.StringType,
			"install_latest_windows10_on_windows11_ineligible_device": types.BoolType,
			"policy_id": types.StringType,
		})
	}

	featureUpdateAttrs := map[string]attr.Value{
		"target_os_version": convert.MapToFrameworkString(featureUpdateRaw, "targetOSVersion"),
		"install_latest_windows10_on_windows11_ineligible_device": convert.MapToFrameworkBool(featureUpdateRaw, "installLatestWindows10OnWindows11IneligibleDevice"),
		"policy_id": convert.MapToFrameworkString(featureUpdateRaw, "policyId"),
	}

	return types.ObjectValueMust(map[string]attr.Type{
		"target_os_version": types.StringType,
		"install_latest_windows10_on_windows11_ineligible_device": types.BoolType,
		"policy_id": types.StringType,
	}, featureUpdateAttrs)
}

// Helper functions for object types
func getDeploymentGroupObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"aad_id":                          types.StringType,
			"name":                            types.StringType,
			"distribution":                    types.Int32Type,
			"failed_prerequisite_check_count": types.Int32Type,
			"user_managed_aad_groups": types.SetType{
				ElemType: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"id":   types.StringType,
						"name": types.StringType,
						"type": types.StringType,
					},
				},
			},
			"deployment_group_policy_settings": types.ObjectType{
				AttrTypes: getPolicySettingsAttrTypes(),
			},
		},
	}
}

func getPolicySettingsAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"aad_group_name":              types.StringType,
		"is_update_settings_modified": types.BoolType,
		"device_configuration_setting": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"policy_id":                   types.StringType,
				"update_behavior":             types.StringType,
				"notification_setting":        types.StringType,
				"quality_deployment_settings": getQualityDeploymentSettingsObjectType(),
				"feature_deployment_settings": getFeatureDeploymentSettingsObjectType(),
			},
		},
		"dnf_update_cloud_setting": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"policy_id":                   types.StringType,
				"approval_type":               types.StringType,
				"deployment_deferral_in_days": types.Int32Type,
			},
		},
		"office_dcv2_setting": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"policy_id":                  types.StringType,
				"deadline":                   types.Int32Type,
				"deferral":                   types.Int32Type,
				"hide_update_notifications":  types.BoolType,
				"target_channel":             types.StringType,
				"enable_automatic_update":    types.BoolType,
				"hide_enable_disable_update": types.BoolType,
				"enable_office_mgmt":         types.BoolType,
				"update_path":                types.StringType,
			},
		},
		"edge_dcv2_setting": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"policy_id":      types.StringType,
				"target_channel": types.StringType,
			},
		},
		"feature_update_anchor_cloud_setting": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"target_os_version": types.StringType,
				"install_latest_windows10_on_windows11_ineligible_device": types.BoolType,
				"policy_id": types.StringType,
			},
		},
	}
}

func getQualityDeploymentSettingsObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"deadline":     types.Int32Type,
			"deferral":     types.Int32Type,
			"grace_period": types.Int32Type,
		},
	}
}

func getFeatureDeploymentSettingsObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"deadline": types.Int32Type,
			"deferral": types.Int32Type,
		},
	}
}
