package graphBetaAospDeviceOwnerCompliancePolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the remote GraphServiceClient object to a Terraform state.
func MapRemoteStateToTerraform(ctx context.Context, data *DeviceCompliancePolicyResourceModel, remoteResource graphmodels.DeviceCompliancePolicyable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())

	// This resource only handles AOSP device owner compliance policies
	if aospPolicy, ok := remoteResource.(*graphmodels.AospDeviceOwnerCompliancePolicy); ok {
		mapAospCompliancePolicyToState(ctx, data, aospPolicy)
	} else {
		tflog.Error(ctx, "Remote resource is not an AOSP device owner compliance policy")
		return
	}

	// Map scheduled actions using SDK getters
	if scheduledActions := remoteResource.GetScheduledActionsForRule(); scheduledActions != nil {
		mappedScheduledActions, err := mapScheduledActionsForRuleToState(ctx, scheduledActions)
		if err != nil {
			tflog.Error(ctx, "Failed to map scheduled actions for rule", map[string]interface{}{
				"error": err.Error(),
			})
		} else {
			data.ScheduledActionsForRule = mappedScheduledActions
		}
	}

	assignments := remoteResource.GetAssignments()
	tflog.Debug(ctx, "Retrieved assignments from remote resource", map[string]interface{}{
		"assignmentCount": len(assignments),
		"resourceId":      data.ID.ValueString(),
	})

	if len(assignments) == 0 {
		tflog.Debug(ctx, "No assignments found, setting assignments to null", map[string]interface{}{
			"resourceId": data.ID.ValueString(),
		})
		data.Assignments = types.SetNull(AospDeviceOwnerCompliancePolicyAssignmentType())
	} else {
		tflog.Debug(ctx, "Starting assignment mapping process", map[string]interface{}{
			"resourceId":      data.ID.ValueString(),
			"assignmentCount": len(assignments),
		})
		MapAssignmentsToTerraform(ctx, data, assignments)
		tflog.Debug(ctx, "Completed assignment mapping process", map[string]interface{}{
			"resourceId": data.ID.ValueString(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}

// mapAospCompliancePolicyToState is a responder function that maps AOSP device owner compliance policy properties.
func mapAospCompliancePolicyToState(ctx context.Context, data *DeviceCompliancePolicyResourceModel, policy *graphmodels.AospDeviceOwnerCompliancePolicy) {
	// Password settings - AOSP uses "password" in SDK but "passcode" in Terraform schema
	data.PasscodeRequired = convert.GraphToFrameworkBool(policy.GetPasswordRequired())
	data.PasscodeMinimumLength = convert.GraphToFrameworkInt32(policy.GetPasswordMinimumLength())
	data.PasscodeMinutesOfInactivityBeforeLock = convert.GraphToFrameworkInt32(policy.GetPasswordMinutesOfInactivityBeforeLock())
	data.PasscodeRequiredType = convert.GraphToFrameworkEnum(policy.GetPasswordRequiredType())

	// OS version settings
	data.OsMinimumVersion = convert.GraphToFrameworkString(policy.GetOsMinimumVersion())
	data.OsMaximumVersion = convert.GraphToFrameworkString(policy.GetOsMaximumVersion())

	// Security settings
	data.SecurityBlockJailbrokenDevices = convert.GraphToFrameworkBool(policy.GetSecurityBlockJailbrokenDevices())
	data.StorageRequireEncryption = convert.GraphToFrameworkBool(policy.GetStorageRequireEncryption())

	// Android-specific settings
	data.MinAndroidSecurityPatchLevel = convert.GraphToFrameworkString(policy.GetMinAndroidSecurityPatchLevel())
}

// mapScheduledActionsForRuleToState maps scheduled actions for rule from SDK to state.
func mapScheduledActionsForRuleToState(ctx context.Context, scheduledActions []graphmodels.DeviceComplianceScheduledActionForRuleable) (types.List, error) {
	scheduledActionType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"rule_name": types.StringType,
			"scheduled_action_configurations": types.SetType{ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"action_type":                  types.StringType,
					"grace_period_hours":           types.Int32Type,
					"notification_template_id":     types.StringType,
					"notification_message_cc_list": types.ListType{ElemType: types.StringType},
				},
			}},
		},
	}

	if len(scheduledActions) == 0 {
		return types.ListNull(scheduledActionType), nil
	}

	actionValues := make([]attr.Value, 0, len(scheduledActions))

	for _, action := range scheduledActions {
		var mappedConfigs types.Set
		if configs := action.GetScheduledActionConfigurations(); configs != nil {
			var err error
			mappedConfigs, err = mapScheduledActionConfigurationsToState(ctx, configs)
			if err != nil {
				return types.ListNull(scheduledActionType), err
			}
		} else {
			mappedConfigs = types.SetNull(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"action_type":                  types.StringType,
					"grace_period_hours":           types.Int32Type,
					"notification_template_id":     types.StringType,
					"notification_message_cc_list": types.ListType{ElemType: types.StringType},
				},
			})
		}

		actionAttrs := map[string]attr.Value{
			"rule_name":                       convert.GraphToFrameworkString(action.GetRuleName()),
			"scheduled_action_configurations": mappedConfigs,
		}

		actionValue, _ := types.ObjectValue(scheduledActionType.AttrTypes, actionAttrs)
		actionValues = append(actionValues, actionValue)
	}

	list, diags := types.ListValue(scheduledActionType, actionValues)
	if diags.HasError() {
		return types.ListNull(scheduledActionType), fmt.Errorf("failed to create scheduled actions list")
	}
	return list, nil
}

// mapScheduledActionConfigurationsToState maps scheduled action configurations from SDK to state.
func mapScheduledActionConfigurationsToState(ctx context.Context, configurations []graphmodels.DeviceComplianceActionItemable) (types.Set, error) {
	configurationType := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"action_type":                  types.StringType,
			"grace_period_hours":           types.Int32Type,
			"notification_template_id":     types.StringType,
			"notification_message_cc_list": types.ListType{ElemType: types.StringType},
		},
	}

	configValues := make([]attr.Value, 0, len(configurations))

	for _, config := range configurations {
		configAttrs := map[string]attr.Value{
			"action_type":                  convert.GraphToFrameworkEnum(config.GetActionType()),
			"grace_period_hours":           convert.GraphToFrameworkInt32(config.GetGracePeriodHours()),
			"notification_template_id":     convert.GraphToFrameworkString(config.GetNotificationTemplateId()),
			"notification_message_cc_list": convert.GraphToFrameworkStringList(config.GetNotificationMessageCCList()),
		}

		configValue, _ := types.ObjectValue(configurationType.AttrTypes, configAttrs)
		configValues = append(configValues, configValue)
	}

	set, diags := types.SetValue(configurationType, configValues)
	if diags.HasError() {
		return types.SetNull(configurationType), fmt.Errorf("failed to create scheduled action configurations set")
	}
	return set, nil
}
