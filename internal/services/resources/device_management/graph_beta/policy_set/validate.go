package graphBetaPolicySet

import (
	"context"
	"fmt"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	deviceappmanagement "github.com/microsoftgraph/msgraph-beta-sdk-go/deviceappmanagement"
	devicemanagement "github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

// validateRequest validates the policy set request data before creating/updating
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data *PolicySetResourceModel, resp any, requiredPermissions []string) error {
	if data.Items.IsNull() || data.Items.IsUnknown() {
		return nil
	}

	var itemModels []PolicySetItemModel
	diags := data.Items.ElementsAs(ctx, &itemModels, false)
	if diags.HasError() {
		return fmt.Errorf("failed to convert items for validation: %v", diags)
	}

	for i, item := range itemModels {
		if err := validatePolicySetItem(ctx, client, &item, i, resp, requiredPermissions); err != nil {
			return err
		}
	}

	return nil
}

// validatePolicySetItem validates a single policy set item
func validatePolicySetItem(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, item *PolicySetItemModel, index int, resp any, requiredPermissions []string) error {
	itemType := item.Type.ValueString()
	payloadId := item.PayloadId.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Validating policy set item %d: type=%s, payloadId=%s", index, itemType, payloadId))

	switch itemType {
	case "app":
		return validateMobileAppId(ctx, client, payloadId, index, resp, requiredPermissions)
	case "app_configuration_policy":
		return validateManagedAppPolicyId(ctx, client, payloadId, index, resp, requiredPermissions)
	case "app_protection_policy":
		return validateManagedAppPolicyId(ctx, client, payloadId, index, resp, requiredPermissions)
	case "device_configuration_profile":
		return validateDeviceConfigurationId(ctx, client, payloadId, index, resp, requiredPermissions)
	case "device_management_configuration_policy":
		// TODO: Add validation for device management configuration policies if needed
		return nil
	case "device_compliance_policy":
		return validateDeviceCompliancePolicyId(ctx, client, payloadId, index, resp, requiredPermissions)
	case "windows_autopilot_deployment_profile":
		return validateWindowsAutopilotDeploymentProfileId(ctx, client, payloadId, index, resp, requiredPermissions)
	default:
		return fmt.Errorf("unsupported policy set item type at index %d: %s", index, itemType)
	}
}

// validateMobileAppId validates that a mobile app ID exists in Intune
func validateMobileAppId(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, appId string, index int, resp any, requiredPermissions []string) error {
	tflog.Debug(ctx, fmt.Sprintf("Validating mobile app ID: %s", appId))

	// filter as per graph x ray request
	filter := "(isof('microsoft.graph.iosStoreApp') or isof('microsoft.graph.iosLobApp') or (isof('microsoft.graph.managedIOSStoreApp') and microsoft.graph.managedApp/appAvailability eq microsoft.graph.managedAppAvailability'global') or isof('microsoft.graph.managedIOSLobApp') or isof('microsoft.graph.androidStoreApp') or isof('microsoft.graph.androidLobApp') or (isof('microsoft.graph.managedAndroidStoreApp') and microsoft.graph.managedApp/appAvailability eq microsoft.graph.managedAppAvailability'global') or isof('microsoft.graph.managedAndroidLobApp') or isof('microsoft.graph.officeSuiteApp') or isof('microsoft.graph.webApp') or isof('microsoft.graph.windowsMobileMSI') or isof('microsoft.graph.windowsMicrosoftEdgeApp') or isof('microsoft.graph.macOSOfficeSuiteApp') or isof('microsoft.graph.macOSLobApp') or isof('microsoft.graph.macOSMicrosoftEdgeApp') or isof('microsoft.graph.macOSMicrosoftDefenderApp') or (isof('microsoft.graph.managedAndroidStoreApp') and microsoft.graph.managedApp/appAvailability eq microsoft.graph.managedAppAvailability'lineOfBusiness') or (isof('microsoft.graph.managedIOSStoreApp') and microsoft.graph.managedApp/appAvailability eq microsoft.graph.managedAppAvailability'lineOfBusiness')) and (microsoft.graph.managedApp/appAvailability eq null or microsoft.graph.managedApp/appAvailability eq 'lineOfBusiness' or isAssigned eq true)"
	orderby := "displayName"

	requestConfig := &deviceappmanagement.MobileAppsRequestBuilderGetRequestConfiguration{
		QueryParameters: &deviceappmanagement.MobileAppsRequestBuilderGetQueryParameters{
			Filter:  &filter,
			Orderby: []string{orderby},
		},
	}

	apps, err := client.
		DeviceAppManagement().
		MobileApps().
		Get(ctx, requestConfig)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "validation", requiredPermissions)
		return err
	}

	if apps == nil || apps.GetValue() == nil {
		return fmt.Errorf("mobile app ID %s at index %d is not a valid Intune mobile app ID", appId, index)
	}

	for _, app := range apps.GetValue() {
		if app.GetId() != nil && *app.GetId() == appId {
			tflog.Debug(ctx, fmt.Sprintf("Mobile app ID validated successfully: %s", appId))
			return nil
		}
	}

	return fmt.Errorf("mobile app ID %s at index %d is not a valid Intune mobile app ID", appId, index)
}

// validateManagedAppPolicyId validates that a managed app policy ID exists in Intune
func validateManagedAppPolicyId(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, policyId string, index int, resp any, requiredPermissions []string) error {
	tflog.Debug(ctx, fmt.Sprintf("Validating managed app policy ID: %s", policyId))

	// Create request configuration for managedAppPolicies endpoint
	requestConfig := &deviceappmanagement.ManagedAppPoliciesRequestBuilderGetRequestConfiguration{
		QueryParameters: &deviceappmanagement.ManagedAppPoliciesRequestBuilderGetQueryParameters{},
	}

	policies, err := client.
		DeviceAppManagement().
		ManagedAppPolicies().
		Get(ctx, requestConfig)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "validation", requiredPermissions)
		return err
	}

	if policies == nil || policies.GetValue() == nil {
		return fmt.Errorf("managed app policy ID %s at index %d is not a valid Intune managed app policy ID", policyId, index)
	}

	// Check if the policy ID exists in the returned list
	for _, policy := range policies.GetValue() {
		if policy.GetId() != nil && *policy.GetId() == policyId {
			tflog.Debug(ctx, fmt.Sprintf("Managed app policy ID validated successfully: %s", policyId))
			return nil
		}
	}

	return fmt.Errorf("managed app policy ID %s at index %d is not a valid Intune managed app policy ID", policyId, index)
}

// validateDeviceConfigurationId validates that a device configuration ID exists in Intune
func validateDeviceConfigurationId(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, configId string, index int, resp any, requiredPermissions []string) error {
	tflog.Debug(ctx, fmt.Sprintf("Validating device configuration ID: %s", configId))

	// Create request configuration for deviceConfigurations endpoint with select and top parameters
	selectFields := "id,displayName,lastModifiedDateTime,roleScopeTagIds,microsoft.graph.unsupportedDeviceConfiguration/originalEntityTypeName"
	top := int32(1000)

	requestConfig := &devicemanagement.DeviceConfigurationsRequestBuilderGetRequestConfiguration{
		QueryParameters: &devicemanagement.DeviceConfigurationsRequestBuilderGetQueryParameters{
			Select: []string{selectFields},
			Top:    &top,
		},
	}

	configurations, err := client.
		DeviceManagement().
		DeviceConfigurations().
		Get(ctx, requestConfig)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "validation", requiredPermissions)
		return err
	}

	if configurations == nil || configurations.GetValue() == nil {
		return fmt.Errorf("device configuration ID %s at index %d is not a valid Intune device configuration ID", configId, index)
	}

	// Check if the configuration ID exists in the returned list
	for _, config := range configurations.GetValue() {
		if config.GetId() != nil && *config.GetId() == configId {
			tflog.Debug(ctx, fmt.Sprintf("Device configuration ID validated successfully: %s", configId))
			return nil
		}
	}

	return fmt.Errorf("device configuration ID %s at index %d is not a valid Intune device configuration ID", configId, index)
}

// validateWindowsAutopilotDeploymentProfileId validates that a Windows Autopilot deployment profile ID exists in Intune
func validateWindowsAutopilotDeploymentProfileId(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, profileId string, index int, resp any, requiredPermissions []string) error {
	tflog.Debug(ctx, fmt.Sprintf("Validating Windows Autopilot deployment profile ID: %s", profileId))

	// Create request configuration for windowsAutopilotDeploymentProfiles endpoint
	requestConfig := &devicemanagement.WindowsAutopilotDeploymentProfilesRequestBuilderGetRequestConfiguration{
		QueryParameters: &devicemanagement.WindowsAutopilotDeploymentProfilesRequestBuilderGetQueryParameters{},
	}

	profiles, err := client.
		DeviceManagement().
		WindowsAutopilotDeploymentProfiles().
		Get(ctx, requestConfig)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "validation", requiredPermissions)
		return err
	}

	if profiles == nil || profiles.GetValue() == nil {
		return fmt.Errorf("Windows Autopilot deployment profile ID %s at index %d is not a valid Intune deployment profile ID", profileId, index)
	}

	// Check if the profile ID exists in the returned list
	for _, profile := range profiles.GetValue() {
		if profile.GetId() != nil && *profile.GetId() == profileId {
			tflog.Debug(ctx, fmt.Sprintf("Windows Autopilot deployment profile ID validated successfully: %s", profileId))
			return nil
		}
	}

	return fmt.Errorf("Windows Autopilot deployment profile ID %s at index %d is not a valid Intune deployment profile ID", profileId, index)
}

// validateDeviceCompliancePolicyId validates that a device compliance policy ID exists in Intune
func validateDeviceCompliancePolicyId(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, policyId string, index int, resp any, requiredPermissions []string) error {
	tflog.Debug(ctx, fmt.Sprintf("Validating device compliance policy ID: %s", policyId))

	// Create request configuration for deviceCompliancePolicies endpoint with complex select, expand, and top parameters
	selectFields := []string{
		"id",
		"displayName",
		"lastModifiedDateTime",
		"roleScopeTagIds",
		"microsoft.graph.androidCompliancePolicy/deviceThreatProtectionRequiredSecurityLevel",
		"microsoft.graph.androidWorkProfileCompliancePolicy/deviceThreatProtectionRequiredSecurityLevel",
		"microsoft.graph.iosCompliancePolicy/deviceThreatProtectionRequiredSecurityLevel",
		"microsoft.graph.windows10CompliancePolicy/deviceThreatProtectionRequiredSecurityLevel",
		"microsoft.graph.iosCompliancePolicy/advancedThreatProtectionRequiredSecurityLevel",
		"microsoft.graph.androidWorkProfileCompliancePolicy/advancedThreatProtectionRequiredSecurityLevel",
		"microsoft.graph.androidDeviceOwnerCompliancePolicy/advancedThreatProtectionRequiredSecurityLevel",
		"microsoft.graph.androidDeviceOwnerCompliancePolicy/deviceThreatProtectionRequiredSecurityLevel",
		"microsoft.graph.androidCompliancePolicy/advancedThreatProtectionRequiredSecurityLevel",
	}
	expand := []string{"assignments"}
	top := int32(1000)

	requestConfig := &devicemanagement.DeviceCompliancePoliciesRequestBuilderGetRequestConfiguration{
		QueryParameters: &devicemanagement.DeviceCompliancePoliciesRequestBuilderGetQueryParameters{
			Select: selectFields,
			Expand: expand,
			Top:    &top,
		},
	}

	policies, err := client.
		DeviceManagement().
		DeviceCompliancePolicies().
		Get(ctx, requestConfig)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "validation", requiredPermissions)
		return err
	}

	if policies == nil || policies.GetValue() == nil {
		return fmt.Errorf("device compliance policy ID %s at index %d is not a valid Intune compliance policy ID", policyId, index)
	}

	// Check if the policy ID exists in the returned list
	for _, policy := range policies.GetValue() {
		if policy.GetId() != nil && *policy.GetId() == policyId {
			tflog.Debug(ctx, fmt.Sprintf("Device compliance policy ID validated successfully: %s", policyId))
			return nil
		}
	}

	return fmt.Errorf("device compliance policy ID %s at index %d is not a valid Intune compliance policy ID", policyId, index)
}
