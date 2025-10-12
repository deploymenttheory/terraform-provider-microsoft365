package graphBetaPolicySet

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	deviceappmanagement "github.com/microsoftgraph/msgraph-beta-sdk-go/deviceappmanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructBasePatchRequest constructs the base PATCH request for the Policy Set resource
func constructBasePatchRequest(ctx context.Context, data *PolicySetResourceModel) (graphmodels.PolicySetable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing base PATCH request for %s resource", ResourceName))

	patchBody := graphmodels.NewPolicySet()
	convert.FrameworkToGraphString(data.DisplayName, patchBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, patchBody.SetDescription)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for %s base PATCH", ResourceName), patchBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	return patchBody, nil
}

// constructItemsUpdateRequest constructs the items update request for the Policy Set resources
func constructItemsUpdateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, currentData, planData *PolicySetResourceModel, resp any, requiredPermissions []string) (deviceappmanagement.PolicySetsItemUpdatePostRequestBodyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing policy set items update request for %s resource", ResourceName))

	if err := validateRequest(ctx, client, planData, resp, requiredPermissions); err != nil {
		return nil, err
	}

	updateRequest := deviceappmanagement.NewPolicySetsItemUpdatePostRequestBody()

	currentItems := make(map[string]PolicySetItemModel)
	if !currentData.Items.IsNull() && !currentData.Items.IsUnknown() {
		var currentItemModels []PolicySetItemModel
		diags := currentData.Items.ElementsAs(ctx, &currentItemModels, false)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to convert current items: %v", diags)
		}
		for _, item := range currentItemModels {
			key := item.PayloadId.ValueString()
			currentItems[key] = item
		}
	}

	plannedItems := make(map[string]PolicySetItemModel)
	if !planData.Items.IsNull() && !planData.Items.IsUnknown() {
		var plannedItemModels []PolicySetItemModel
		diags := planData.Items.ElementsAs(ctx, &plannedItemModels, false)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to convert planned items: %v", diags)
		}
		for _, item := range plannedItemModels {
			key := item.PayloadId.ValueString()
			plannedItems[key] = item
		}
	}

	addedItems := make([]graphmodels.PolicySetItemable, 0)
	for key, item := range plannedItems {
		if _, exists := currentItems[key]; !exists {
			policySetItem, err := constructPolicySetItemForUpdate(ctx, &item)
			if err != nil {
				return nil, fmt.Errorf("failed to construct added item: %s", err)
			}
			addedItems = append(addedItems, policySetItem)
		}
	}

	updatedItems := make([]graphmodels.PolicySetItemable, 0)
	for key, item := range plannedItems {
		if currentItem, exists := currentItems[key]; exists {
			if itemHasChanges(&currentItem, &item) {
				policySetItem, err := constructPolicySetItemForUpdate(ctx, &item)
				if err != nil {
					return nil, fmt.Errorf("failed to construct updated item: %s", err)
				}
				updatedItems = append(updatedItems, policySetItem)
			}
		}
	}

	deletedItems := make([]string, 0)
	for key, currentItem := range currentItems {
		if _, exists := plannedItems[key]; !exists {
			if !currentItem.ID.IsNull() && !currentItem.ID.IsUnknown() {
				deletedItems = append(deletedItems, currentItem.ID.ValueString())
			}
		}
	}

	updateRequest.SetAddedPolicySetItems(addedItems)
	updateRequest.SetUpdatedPolicySetItems(updatedItems)
	updateRequest.SetDeletedPolicySetItems(deletedItems)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for %s items update", ResourceName), updateRequest); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s items update request", ResourceName))

	return updateRequest, nil
}

func constructAssignmentsUpdateRequest(ctx context.Context, data *PolicySetResourceModel) (deviceappmanagement.PolicySetsItemUpdatePostRequestBodyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing assignments update request for %s resource", ResourceName))

	updateRequest := deviceappmanagement.NewPolicySetsItemUpdatePostRequestBody()

	if !data.Assignments.IsNull() && !data.Assignments.IsUnknown() {
		assignments, err := constructAssignmentsForUpdate(ctx, data.Assignments)
		if err != nil {
			return nil, fmt.Errorf("failed to construct assignments for update: %s", err)
		}
		updateRequest.SetAssignments(assignments)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for %s assignments update", ResourceName), updateRequest); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s assignments update request", ResourceName))

	return updateRequest, nil
}

func itemHasChanges(current, planned *PolicySetItemModel) bool {
	if current.Intent.ValueString() != planned.Intent.ValueString() {
		return true
	}
	if !current.Settings.Equal(planned.Settings) {
		return true
	}
	return false
}

func constructAssignmentsForUpdate(ctx context.Context, assignmentsSet types.Set) ([]graphmodels.PolicySetAssignmentable, error) {
	if assignmentsSet.IsNull() || assignmentsSet.IsUnknown() {
		return nil, nil
	}

	var assignmentModels []PolicySetAssignmentModel
	diags := assignmentsSet.ElementsAs(ctx, &assignmentModels, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to convert assignments set: %v", diags)
	}

	assignments := make([]graphmodels.PolicySetAssignmentable, len(assignmentModels))
	for i, assignmentModel := range assignmentModels {
		assignment := graphmodels.NewPolicySetAssignment()

		target, err := constructAssignmentTargetForUpdate(ctx, &assignmentModel)
		if err != nil {
			return nil, fmt.Errorf("failed to construct assignment target: %s", err)
		}
		assignment.SetTarget(target)

		assignments[i] = assignment
	}

	return assignments, nil
}

func constructAssignmentTargetForUpdate(_ context.Context, assignmentModel *PolicySetAssignmentModel) (graphmodels.DeviceAndAppManagementAssignmentTargetable, error) {
	assignmentType := assignmentModel.Type.ValueString()

	switch assignmentType {
	case "groupAssignmentTarget":
		target := graphmodels.NewGroupAssignmentTarget()
		convert.FrameworkToGraphString(assignmentModel.GroupId, target.SetGroupId)
		return target, nil
	case "exclusionGroupAssignmentTarget":
		target := graphmodels.NewExclusionGroupAssignmentTarget()
		convert.FrameworkToGraphString(assignmentModel.GroupId, target.SetGroupId)
		return target, nil
	default:
		return nil, fmt.Errorf("unsupported assignment target type: %s. Valid types are: 'groupAssignmentTarget', 'exclusionGroupAssignmentTarget'", assignmentType)
	}
}

func constructPolicySetItemForUpdate(ctx context.Context, itemModel *PolicySetItemModel) (graphmodels.PolicySetItemable, error) {
	odataType, err := resolveODataTypeForSetItemUpdate(itemModel.Type.ValueString())
	if err != nil {
		return nil, err
	}

	switch odataType {
	case "#microsoft.graph.mobileAppPolicySetItem":
		item := graphmodels.NewMobileAppPolicySetItem()
		convert.FrameworkToGraphString(itemModel.PayloadId, item.SetPayloadId)
		if !itemModel.Intent.IsNull() && !itemModel.Intent.IsUnknown() {
			intentValue := itemModel.Intent.ValueString()
			switch intentValue {
			case "required":
				intent := graphmodels.REQUIRED_INSTALLINTENT
				item.SetIntent(&intent)
			case "available":
				intent := graphmodels.AVAILABLE_INSTALLINTENT
				item.SetIntent(&intent)
			case "uninstall":
				intent := graphmodels.UNINSTALL_INSTALLINTENT
				item.SetIntent(&intent)
			case "availableWithoutEnrollment":
				intent := graphmodels.AVAILABLEWITHOUTENROLLMENT_INSTALLINTENT
				item.SetIntent(&intent)
			}
		}

		if !itemModel.Settings.IsNull() && !itemModel.Settings.IsUnknown() {
			settings, err := constructPolicySetItemSettingsForUpdate(ctx, itemModel.Settings)
			if err != nil {
				return nil, fmt.Errorf("failed to construct mobile app settings: %s", err)
			}
			item.SetSettings(settings)
		}

		return item, nil

	case "#microsoft.graph.targetedManagedAppConfigurationPolicySetItem":
		item := graphmodels.NewTargetedManagedAppConfigurationPolicySetItem()
		convert.FrameworkToGraphString(itemModel.PayloadId, item.SetPayloadId)
		return item, nil

	case "#microsoft.graph.managedAppProtectionPolicySetItem":
		item := graphmodels.NewManagedAppProtectionPolicySetItem()
		convert.FrameworkToGraphString(itemModel.PayloadId, item.SetPayloadId)
		return item, nil

	case "#microsoft.graph.deviceConfigurationPolicySetItem":
		item := graphmodels.NewDeviceConfigurationPolicySetItem()
		convert.FrameworkToGraphString(itemModel.PayloadId, item.SetPayloadId)
		return item, nil

	case "#microsoft.graph.deviceManagementConfigurationPolicyPolicySetItem":
		item := graphmodels.NewDeviceManagementConfigurationPolicyPolicySetItem()
		convert.FrameworkToGraphString(itemModel.PayloadId, item.SetPayloadId)
		return item, nil

	case "#microsoft.graph.deviceCompliancePolicyPolicySetItem":
		item := graphmodels.NewDeviceCompliancePolicyPolicySetItem()
		convert.FrameworkToGraphString(itemModel.PayloadId, item.SetPayloadId)
		return item, nil

	case "#microsoft.graph.windowsAutopilotDeploymentProfilePolicySetItem":
		item := graphmodels.NewWindowsAutopilotDeploymentProfilePolicySetItem()
		convert.FrameworkToGraphString(itemModel.PayloadId, item.SetPayloadId)
		return item, nil

	case "#microsoft.graph.windows10EnrollmentCompletionPageConfigurationPolicySetItem":
		item := graphmodels.NewWindows10EnrollmentCompletionPageConfigurationPolicySetItem()
		convert.FrameworkToGraphString(itemModel.PayloadId, item.SetPayloadId)
		return item, nil

	default:
		return nil, fmt.Errorf("unsupported policy set item type: %s", odataType)
	}
}

func resolveODataTypeForSetItemUpdate(userType string) (string, error) {
	switch userType {
	case "app":
		return "#microsoft.graph.mobileAppPolicySetItem", nil
	case "app_configuration_policy":
		return "#microsoft.graph.targetedManagedAppConfigurationPolicySetItem", nil
	case "app_protection_policy":
		return "#microsoft.graph.managedAppProtectionPolicySetItem", nil
	case "device_configuration_profile":
		return "#microsoft.graph.deviceConfigurationPolicySetItem", nil
	case "device_management_configuration_policy":
		return "#microsoft.graph.deviceManagementConfigurationPolicyPolicySetItem", nil
	case "device_compliance_policy":
		return "#microsoft.graph.deviceCompliancePolicyPolicySetItem", nil
	case "windows_autopilot_deployment_profile":
		return "#microsoft.graph.windowsAutopilotDeploymentProfilePolicySetItem", nil
	case "enrollment_status_page":
		return "#microsoft.graph.windows10EnrollmentCompletionPageConfigurationPolicySetItem", nil
	default:
		return "", fmt.Errorf("unsupported policy set item type: %s. Valid types are: app, app_configuration_policy, app_protection_policy, device_configuration_profile, device_management_configuration_policy, device_compliance_policy, windows_autopilot_deployment_profile, enrollment_status_page", userType)
	}
}

func constructPolicySetItemSettingsForUpdate(ctx context.Context, settingsObj types.Object) (graphmodels.MobileAppAssignmentSettingsable, error) {
	if settingsObj.IsNull() || settingsObj.IsUnknown() {
		return nil, nil
	}

	var settingsModel PolicySetItemSettingsModel
	diags := settingsObj.As(ctx, &settingsModel, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, fmt.Errorf("failed to convert settings object: %v", diags)
	}

	odataType := settingsModel.ODataType.ValueString()

	switch odataType {
	case "#microsoft.graph.iosStoreAppAssignmentSettings":
		settings := graphmodels.NewIosStoreAppAssignmentSettings()
		convert.FrameworkToGraphString(settingsModel.VpnConfigurationId, settings.SetVpnConfigurationId)
		convert.FrameworkToGraphBool(settingsModel.UninstallOnDeviceRemoval, settings.SetUninstallOnDeviceRemoval)
		convert.FrameworkToGraphBool(settingsModel.IsRemovable, settings.SetIsRemovable)
		return settings, nil

	default:
		return nil, fmt.Errorf("unsupported mobile app assignment settings type: %s", odataType)
	}
}
