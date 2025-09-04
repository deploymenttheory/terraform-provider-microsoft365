package graphBetaAutopatchGroups

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// MapRemoteStateToTerraform maps the remote state from the API to the Terraform model
func MapRemoteStateToTerraform(ctx context.Context, data *AutopatchGroupsResourceModel, autopatchGroup map[string]interface{}) {
	tflog.Debug(ctx, fmt.Sprintf("Mapping %s resource from API to Terraform state", ResourceName))

	// Basic string fields
	if id, ok := autopatchGroup["id"].(string); ok {
		data.ID = types.StringValue(id)
	}

	if name, ok := autopatchGroup["name"].(string); ok {
		data.Name = types.StringValue(name)
	}

	if description, ok := autopatchGroup["description"].(string); ok {
		data.Description = types.StringValue(description)
	} else {
		data.Description = types.StringNull()
	}

	if tenantId, ok := autopatchGroup["tenantId"].(string); ok {
		data.TenantId = types.StringValue(tenantId)
	} else {
		data.TenantId = types.StringNull()
	}

	if groupType, ok := autopatchGroup["type"].(string); ok {
		data.Type = types.StringValue(groupType)
	} else {
		data.Type = types.StringNull()
	}

	if status, ok := autopatchGroup["status"].(string); ok {
		data.Status = types.StringValue(status)
	} else {
		data.Status = types.StringNull()
	}

	if distributionType, ok := autopatchGroup["distributionType"].(string); ok {
		data.DistributionType = types.StringValue(distributionType)
	} else {
		data.DistributionType = types.StringNull()
	}

	if flowId, ok := autopatchGroup["flowId"].(string); ok {
		data.FlowId = types.StringValue(flowId)
	} else {
		data.FlowId = types.StringNull()
	}

	if flowType, ok := autopatchGroup["flowType"].(string); ok {
		data.FlowType = types.StringValue(flowType)
	} else {
		data.FlowType = types.StringNull()
	}

	if flowStatus, ok := autopatchGroup["flowStatus"].(string); ok {
		data.FlowStatus = types.StringValue(flowStatus)
	} else {
		data.FlowStatus = types.StringNull()
	}

	if umbrellaGroupId, ok := autopatchGroup["umbrellaGroupId"].(string); ok {
		data.UmbrellaGroupId = types.StringValue(umbrellaGroupId)
	} else {
		data.UmbrellaGroupId = types.StringNull()
	}

	// Boolean fields
	if isLockedByPolicy, ok := autopatchGroup["isLockedByPolicy"].(bool); ok {
		data.IsLockedByPolicy = types.BoolValue(isLockedByPolicy)
	} else {
		data.IsLockedByPolicy = types.BoolNull()
	}

	if readOnly, ok := autopatchGroup["readOnly"].(bool); ok {
		data.ReadOnly = types.BoolValue(readOnly)
	} else {
		data.ReadOnly = types.BoolNull()
	}

	if userHasAllScopeTag, ok := autopatchGroup["userHasAllScopeTag"].(bool); ok {
		data.UserHasAllScopeTag = types.BoolValue(userHasAllScopeTag)
	} else {
		data.UserHasAllScopeTag = types.BoolNull()
	}

	if enableDriverUpdate, ok := autopatchGroup["enableDriverUpdate"].(bool); ok {
		data.EnableDriverUpdate = types.BoolValue(enableDriverUpdate)
	} else {
		data.EnableDriverUpdate = types.BoolNull()
	}

	// Numeric fields
	if numberOfRegisteredDevices, ok := autopatchGroup["numberOfRegisteredDevices"].(float32); ok {
		data.NumberOfRegisteredDevices = types.Int32Value(int32(numberOfRegisteredDevices))
	} else {
		data.NumberOfRegisteredDevices = types.Int32Null()
	}

	if enabledContentTypes, ok := autopatchGroup["enabledContentTypes"].(float32); ok {
		data.EnabledContentTypes = types.Int32Value(int32(enabledContentTypes))
	} else {
		data.EnabledContentTypes = types.Int32Null()
	}

	// Scope Tags
	if scopeTagsRaw, ok := autopatchGroup["scopeTags"].([]interface{}); ok {
		scopeTagsValues := make([]attr.Value, 0, len(scopeTagsRaw))
		for _, tagRaw := range scopeTagsRaw {
			if tagFloat, ok := tagRaw.(float32); ok {
				scopeTagsValues = append(scopeTagsValues, types.StringValue(fmt.Sprintf("%.0f", tagFloat)))
			}
		}
		if len(scopeTagsValues) > 0 {
			data.ScopeTags = types.SetValueMust(types.StringType, scopeTagsValues)
		} else {
			data.ScopeTags = types.SetValueMust(types.StringType, []attr.Value{types.StringValue("0")})
		}
	} else {
		data.ScopeTags = types.SetValueMust(types.StringType, []attr.Value{types.StringValue("0")})
	}

	// Global User Managed AAD Groups
	if globalGroupsRaw, ok := autopatchGroup["globalUserManagedAadGroups"].([]interface{}); ok {
		globalGroupsValues := make([]attr.Value, 0, len(globalGroupsRaw))
		for _, groupRaw := range globalGroupsRaw {
			if groupMap, ok := groupRaw.(map[string]interface{}); ok {
				globalGroup := map[string]attr.Value{
					"id":   types.StringNull(),
					"type": types.StringNull(),
				}

				if id, ok := groupMap["id"].(string); ok {
					globalGroup["id"] = types.StringValue(id)
				}
				if groupType, ok := groupMap["type"].(string); ok {
					globalGroup["type"] = types.StringValue(groupType)
				}

				globalGroupsValues = append(globalGroupsValues, types.ObjectValueMust(
					map[string]attr.Type{
						"id":   types.StringType,
						"type": types.StringType,
					},
					globalGroup,
				))
			}
		}
		if len(globalGroupsValues) > 0 {
			data.GlobalUserManagedAadGroups = types.SetValueMust(
				types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"id":   types.StringType,
						"type": types.StringType,
					},
				},
				globalGroupsValues,
			)
		} else {
			data.GlobalUserManagedAadGroups = types.SetValueMust(
				types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"id":   types.StringType,
						"type": types.StringType,
					},
				},
				[]attr.Value{},
			)
		}
	} else {
		data.GlobalUserManagedAadGroups = types.SetValueMust(
			types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"id":   types.StringType,
					"type": types.StringType,
				},
			},
			[]attr.Value{},
		)
	}

	// Deployment Groups - this is complex, so for now set to empty
	// TODO: Implement full deployment groups mapping
	data.DeploymentGroups = types.SetValueMust(
		types.ObjectType{
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
							"type": types.Int32Type,
						},
					},
				},
				"deployment_group_policy_settings": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"aad_group_name":              types.StringType,
						"is_update_settings_modified": types.BoolType,
						"device_configuration_setting": types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"policy_id":            types.StringType,
								"update_behavior":      types.StringType,
								"notification_setting": types.StringType,
								"quality_deployment_settings": types.ObjectType{
									AttrTypes: map[string]attr.Type{
										"deadline":     types.Int32Type,
										"deferral":     types.Int32Type,
										"grace_period": types.Int32Type,
									},
								},
								"feature_deployment_settings": types.ObjectType{
									AttrTypes: map[string]attr.Type{
										"deadline": types.Int32Type,
										"deferral": types.Int32Type,
									},
								},
							},
						},
					},
				},
			},
		},
		[]attr.Value{},
	)

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping %s resource from API to Terraform state", ResourceName))
}
