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
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructResource(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data *PolicySetResourceModel, resp any, requiredPermissions []string) (graphmodels.PolicySetable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	// Validate the request data first
	if err := validateRequest(ctx, client, data, resp, requiredPermissions); err != nil {
		return nil, err
	}

	requestBody := graphmodels.NewPolicySet()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTags); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	if !data.Assignments.IsNull() && !data.Assignments.IsUnknown() {
		assignments, err := constructAssignments(ctx, data.Assignments)
		if err != nil {
			return nil, fmt.Errorf("failed to construct assignments: %s", err)
		}
		requestBody.SetAssignments(assignments)
	}

	if !data.Items.IsNull() && !data.Items.IsUnknown() {
		items, err := constructItems(ctx, data.Items)
		if err != nil {
			return nil, fmt.Errorf("failed to construct items: %s", err)
		}
		requestBody.SetItems(items)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

func constructAssignments(ctx context.Context, assignmentsSet types.Set) ([]graphmodels.PolicySetAssignmentable, error) {
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

		target, err := constructAssignmentTarget(ctx, &assignmentModel)
		if err != nil {
			return nil, fmt.Errorf("failed to construct assignment target: %s", err)
		}
		assignment.SetTarget(target)

		assignments[i] = assignment
	}

	return assignments, nil
}

func constructAssignmentTarget(_ context.Context, assignmentModel *PolicySetAssignmentModel) (graphmodels.DeviceAndAppManagementAssignmentTargetable, error) {
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

func constructItems(ctx context.Context, itemsSet types.Set) ([]graphmodels.PolicySetItemable, error) {
	if itemsSet.IsNull() || itemsSet.IsUnknown() {
		return nil, nil
	}

	var itemModels []PolicySetItemModel
	diags := itemsSet.ElementsAs(ctx, &itemModels, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to convert items set: %v", diags)
	}

	items := make([]graphmodels.PolicySetItemable, len(itemModels))
	for i, itemModel := range itemModels {
		item, err := constructPolicySetItem(ctx, &itemModel)
		if err != nil {
			return nil, fmt.Errorf("failed to construct policy set item at index %d: %s", i, err)
		}
		items[i] = item
	}

	return items, nil
}

func constructPolicySetItem(ctx context.Context, itemModel *PolicySetItemModel) (graphmodels.PolicySetItemable, error) {

	odataType, err := resolveODataTypeForSetItem(itemModel.Type.ValueString())
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
			settings, err := constructPolicySetItemSettings(ctx, itemModel.Settings)
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

	default:
		return nil, fmt.Errorf("unsupported policy set item type: %s", odataType)
	}
}

// resolveODataTypeForSetItem maps user-friendly type names to OData types
func resolveODataTypeForSetItem(userType string) (string, error) {
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
	default:
		return "", fmt.Errorf("unsupported policy set item type: %s. Valid types are: app, app_configuration_policy, app_protection_policy, device_configuration_profile, device_management_configuration_policy, device_compliance_policy, windows_autopilot_deployment_profile", userType)
	}
}

func constructPolicySetItemSettings(ctx context.Context, settingsObj types.Object) (graphmodels.MobileAppAssignmentSettingsable, error) {
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
