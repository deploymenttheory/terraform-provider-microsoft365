package graphBetaPolicySet

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteStateToTerraform(ctx context.Context, data *PolicySetResourceModel, remoteResource graphmodels.PolicySetable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.Status = convert.GraphToFrameworkEnum(remoteResource.GetStatus())
	data.ErrorCode = convert.GraphToFrameworkEnum(remoteResource.GetErrorCode())

	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTags())

	if assignments := remoteResource.GetAssignments(); assignments != nil {
		data.Assignments = mapAssignmentsToTerraform(ctx, assignments)
	} else {
		data.Assignments = types.SetNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"type":     types.StringType,
				"group_id": types.StringType,
			},
		})
	}

	if items := remoteResource.GetItems(); items != nil {
		data.Items = mapItemsToTerraform(ctx, items)
	} else {
		data.Items = types.SetNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"payload_id": types.StringType,
				"type":       types.StringType,
				"intent":     types.StringType,
				"settings": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"odata_type":                  types.StringType,
						"vpn_configuration_id":        types.StringType,
						"uninstall_on_device_removal": types.BoolType,
						"is_removable":                types.BoolType,
						"prevent_managed_app_backup":  types.BoolType,
					},
				},
			},
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping remote state to Terraform for resource %s with id %s", ResourceName, data.ID.ValueString()))
}

func mapAssignmentsToTerraform(_ context.Context, assignments []graphmodels.PolicySetAssignmentable) types.Set {
	if len(assignments) == 0 {
		return types.SetValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"type":     types.StringType,
				"group_id": types.StringType,
			},
		}, []attr.Value{})
	}

	assignmentValues := make([]attr.Value, len(assignments))
	for i, assignment := range assignments {
		assignmentAttrs := map[string]attr.Value{
			"type":     types.StringNull(),
			"group_id": types.StringNull(),
		}

		if target := assignment.GetTarget(); target != nil {
			// Map odata type to user-friendly type
			if odataType := target.GetOdataType(); odataType != nil {
				switch *odataType {
				case "#microsoft.graph.groupAssignmentTarget":
					assignmentAttrs["type"] = types.StringValue("groupAssignmentTarget")
				case "#microsoft.graph.exclusionGroupAssignmentTarget":
					assignmentAttrs["type"] = types.StringValue("exclusionGroupAssignmentTarget")
				}
			}

			// Extract group ID for both group and exclusion group targets
			if groupTarget, ok := target.(graphmodels.GroupAssignmentTargetable); ok {
				if groupId := groupTarget.GetGroupId(); groupId != nil {
					assignmentAttrs["group_id"] = types.StringValue(*groupId)
				}
			} else if exclusionGroupTarget, ok := target.(graphmodels.ExclusionGroupAssignmentTargetable); ok {
				if groupId := exclusionGroupTarget.GetGroupId(); groupId != nil {
					assignmentAttrs["group_id"] = types.StringValue(*groupId)
				}
			}
		}

		assignmentValues[i] = types.ObjectValueMust(
			map[string]attr.Type{
				"type":     types.StringType,
				"group_id": types.StringType,
			},
			assignmentAttrs,
		)
	}

	return types.SetValueMust(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":     types.StringType,
			"group_id": types.StringType,
		},
	}, assignmentValues)
}

func mapItemsToTerraform(ctx context.Context, items []graphmodels.PolicySetItemable) types.Set {
	if len(items) == 0 {
		return types.SetValueMust(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"payload_id": types.StringType,
				"type":       types.StringType,
				"intent":     types.StringType,
				"settings": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"odata_type":                  types.StringType,
						"vpn_configuration_id":        types.StringType,
						"uninstall_on_device_removal": types.BoolType,
						"is_removable":                types.BoolType,
						"prevent_managed_app_backup":  types.BoolType,
					},
				},
			},
		}, []attr.Value{})
	}

	itemValues := make([]attr.Value, len(items))
	for i, item := range items {
		itemAttrs := map[string]attr.Value{
			"payload_id": convert.GraphToFrameworkString(item.GetPayloadId()),
			"type":       types.StringNull(),
			"intent":     types.StringNull(),
			"settings": types.ObjectNull(map[string]attr.Type{
				"odata_type":                  types.StringType,
				"vpn_configuration_id":        types.StringType,
				"uninstall_on_device_removal": types.BoolType,
				"is_removable":                types.BoolType,
				"prevent_managed_app_backup":  types.BoolType,
			}),
		}

		if odataType := item.GetOdataType(); odataType != nil {
			itemAttrs["type"] = types.StringValue(resolveSetItemForOdataType(*odataType))
		}

		if mobileAppItem, ok := item.(graphmodels.MobileAppPolicySetItemable); ok {
			if intent := mobileAppItem.GetIntent(); intent != nil {
				itemAttrs["intent"] = types.StringValue(intent.String())
			}

			if settings := mobileAppItem.GetSettings(); settings != nil {
				itemAttrs["settings"] = mapMobileAppSettingsToTerraform(ctx, settings)
			}
		}

		itemValues[i] = types.ObjectValueMust(
			map[string]attr.Type{
				"payload_id": types.StringType,
				"type":       types.StringType,
				"intent":     types.StringType,
				"settings": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"odata_type":                  types.StringType,
						"vpn_configuration_id":        types.StringType,
						"uninstall_on_device_removal": types.BoolType,
						"is_removable":                types.BoolType,
						"prevent_managed_app_backup":  types.BoolType,
					},
				},
			},
			itemAttrs,
		)
	}

	return types.SetValueMust(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"payload_id": types.StringType,
			"type":       types.StringType,
			"intent":     types.StringType,
			//"guided_deployment_tags": types.SetType{ElemType: types.StringType},
			"settings": types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"odata_type":                  types.StringType,
					"vpn_configuration_id":        types.StringType,
					"uninstall_on_device_removal": types.BoolType,
					"is_removable":                types.BoolType,
					"prevent_managed_app_backup":  types.BoolType,
				},
			},
		},
	}, itemValues)
}

func mapMobileAppSettingsToTerraform(_ context.Context, settings graphmodels.MobileAppAssignmentSettingsable) types.Object {
	settingsAttrs := map[string]attr.Value{
		"odata_type":                  types.StringNull(),
		"vpn_configuration_id":        types.StringNull(),
		"uninstall_on_device_removal": types.BoolNull(),
		"is_removable":                types.BoolNull(),
		"prevent_managed_app_backup":  types.BoolNull(),
	}

	if odataType := settings.GetOdataType(); odataType != nil {
		settingsAttrs["odata_type"] = types.StringValue(*odataType)
	}

	if iosSettings, ok := settings.(graphmodels.IosStoreAppAssignmentSettingsable); ok {
		if vpnId := iosSettings.GetVpnConfigurationId(); vpnId != nil {
			settingsAttrs["vpn_configuration_id"] = types.StringValue(*vpnId)
		}
		if uninstall := iosSettings.GetUninstallOnDeviceRemoval(); uninstall != nil {
			settingsAttrs["uninstall_on_device_removal"] = types.BoolValue(*uninstall)
		}
		if removable := iosSettings.GetIsRemovable(); removable != nil {
			settingsAttrs["is_removable"] = types.BoolValue(*removable)
		}
		if preventBackup := iosSettings.GetPreventManagedAppBackup(); preventBackup != nil {
			settingsAttrs["prevent_managed_app_backup"] = types.BoolValue(*preventBackup)
		}
	}

	return types.ObjectValueMust(
		map[string]attr.Type{
			"odata_type":                  types.StringType,
			"vpn_configuration_id":        types.StringType,
			"uninstall_on_device_removal": types.BoolType,
			"is_removable":                types.BoolType,
			"prevent_managed_app_backup":  types.BoolType,
		},
		settingsAttrs,
	)
}

// resolveSetItemForOdataType maps OData types back to user-friendly type names
func resolveSetItemForOdataType(odataType string) string {
	switch odataType {
	case "#microsoft.graph.mobileAppPolicySetItem":
		return "app"
	case "#microsoft.graph.targetedManagedAppConfigurationPolicySetItem":
		return "app_configuration_policy"
	case "#microsoft.graph.managedAppProtectionPolicySetItem":
		return "app_protection_policy"
	case "#microsoft.graph.deviceConfigurationPolicySetItem":
		return "device_configuration_profile"
	case "#microsoft.graph.deviceManagementConfigurationPolicyPolicySetItem":
		return "device_management_configuration_policy"
	case "#microsoft.graph.deviceCompliancePolicyPolicySetItem":
		return "device_compliance_policy"
	case "#microsoft.graph.windowsAutopilotDeploymentProfilePolicySetItem":
		return "windows_autopilot_deployment_profile"
	default:
		// Return the original OData type if we don't have a mapping
		return odataType
	}
}
