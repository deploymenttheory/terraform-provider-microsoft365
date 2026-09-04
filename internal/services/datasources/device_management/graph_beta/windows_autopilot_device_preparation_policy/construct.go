package graphBetaWindowsAutopilotDevicePreparationPolicy

import (
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// ConstructPolicyItems constructs a list of item models from a list of configuration policies
func ConstructPolicyItems(policies []graphmodels.DeviceManagementConfigurationPolicyable) []WindowsAutopilotDevicePreparationPolicyItemModel {
	if policies == nil {
		return []WindowsAutopilotDevicePreparationPolicyItemModel{}
	}

	items := make([]WindowsAutopilotDevicePreparationPolicyItemModel, 0, len(policies))
	for _, policy := range policies {
		if policy != nil {
			items = append(items, ConstructPolicyItem(policy))
		}
	}

	return items
}

// ConstructPolicyItem constructs an item model from a configuration policy
func ConstructPolicyItem(policy graphmodels.DeviceManagementConfigurationPolicyable) WindowsAutopilotDevicePreparationPolicyItemModel {
	item := WindowsAutopilotDevicePreparationPolicyItemModel{
		ID:                                types.StringPointerValue(policy.GetId()),
		Name:                              types.StringPointerValue(policy.GetName()),
		Description:                       types.StringPointerValue(policy.GetDescription()),
		CreatedDateTime:                   types.StringNull(),
		LastModifiedDateTime:              types.StringNull(),
		CreationSource:                    types.StringPointerValue(policy.GetCreationSource()),
		Platforms:                         types.StringNull(),
		Technologies:                      types.StringNull(),
		SettingCount:                      types.Int64Null(),
		IsAssigned:                        types.BoolPointerValue(policy.GetIsAssigned()),
		DisableEntraGroupPolicyAssignment: types.BoolPointerValue(policy.GetDisableEntraGroupPolicyAssignment()),
		Priority:                          types.Int64Null(),
		RoleScopeTagIds:                   []types.String{},
	}

	if created := policy.GetCreatedDateTime(); created != nil {
		item.CreatedDateTime = types.StringValue(created.Format("2006-01-02T15:04:05Z"))
	}

	if lastModified := policy.GetLastModifiedDateTime(); lastModified != nil {
		item.LastModifiedDateTime = types.StringValue(lastModified.Format("2006-01-02T15:04:05Z"))
	}

	if platforms := policy.GetPlatforms(); platforms != nil {
		item.Platforms = types.StringValue(platforms.String())
	}

	if technologies := policy.GetTechnologies(); technologies != nil {
		item.Technologies = types.StringValue(technologies.String())
	}

	if settingCount := policy.GetSettingCount(); settingCount != nil {
		item.SettingCount = types.Int64Value(int64(*settingCount))
	}

	if priorityMetaData := policy.GetPriorityMetaData(); priorityMetaData != nil {
		if priority := priorityMetaData.GetPriority(); priority != nil {
			item.Priority = types.Int64Value(int64(*priority))
		}
	}

	if roleScopeTagIds := policy.GetRoleScopeTagIds(); len(roleScopeTagIds) > 0 {
		item.RoleScopeTagIds = make([]types.String, 0, len(roleScopeTagIds))
		for _, roleScopeTagId := range roleScopeTagIds {
			item.RoleScopeTagIds = append(item.RoleScopeTagIds, types.StringValue(roleScopeTagId))
		}
	}

	if templateRef := policy.GetTemplateReference(); templateRef != nil {
		item.TemplateReference = &TemplateReferenceModel{
			TemplateId:             types.StringPointerValue(templateRef.GetTemplateId()),
			TemplateFamily:         types.StringNull(),
			TemplateDisplayName:    types.StringPointerValue(templateRef.GetTemplateDisplayName()),
			TemplateDisplayVersion: types.StringPointerValue(templateRef.GetTemplateDisplayVersion()),
			DeploymentMode:         types.StringNull(),
		}

		if templateFamily := templateRef.GetTemplateFamily(); templateFamily != nil {
			item.TemplateReference.TemplateFamily = types.StringValue(templateFamily.String())
		}

		if mode := deploymentModeFromTemplateId(templateRef.GetTemplateId()); mode != "" {
			item.TemplateReference.DeploymentMode = types.StringValue(mode)
		}
	}

	return item
}

// ConstructPolicyAssignments constructs assignment models from a list of policy assignments
func ConstructPolicyAssignments(assignments []graphmodels.DeviceManagementConfigurationPolicyAssignmentable) []PolicyAssignmentModel {
	if assignments == nil {
		return []PolicyAssignmentModel{}
	}

	items := make([]PolicyAssignmentModel, 0, len(assignments))
	for _, assignment := range assignments {
		if assignment == nil {
			continue
		}

		item := PolicyAssignmentModel{
			ID:         types.StringPointerValue(assignment.GetId()),
			Type:       types.StringNull(),
			GroupId:    types.StringNull(),
			FilterId:   types.StringNull(),
			FilterType: types.StringValue("none"),
		}

		target := assignment.GetTarget()
		if target == nil {
			items = append(items, item)
			continue
		}

		if odataType := target.GetOdataType(); odataType != nil {
			item.Type = types.StringValue(strings.TrimPrefix(*odataType, "#microsoft.graph."))
		}

		if groupTarget, ok := target.(graphmodels.GroupAssignmentTargetable); ok {
			if groupId := groupTarget.GetGroupId(); groupId != nil && *groupId != "" {
				item.GroupId = types.StringValue(*groupId)
			}
		}

		if filterId := target.GetDeviceAndAppManagementAssignmentFilterId(); filterId != nil &&
			*filterId != "" && *filterId != "00000000-0000-0000-0000-000000000000" {
			item.FilterId = types.StringValue(*filterId)
		}

		if filterType := target.GetDeviceAndAppManagementAssignmentFilterType(); filterType != nil {
			switch *filterType {
			case graphmodels.INCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE:
				item.FilterType = types.StringValue("include")
			case graphmodels.EXCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE:
				item.FilterType = types.StringValue("exclude")
			default:
				item.FilterType = types.StringValue("none")
			}
		}

		items = append(items, item)
	}

	return items
}
