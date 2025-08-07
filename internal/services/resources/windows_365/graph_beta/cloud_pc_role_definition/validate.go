package graphBetaRoleDefinition

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// validateRequest validates that all provided role permissions are valid Cloud PC operations
func validateRequest(ctx context.Context, permissions []string, client *msgraphbetasdk.GraphServiceClient, displayName string) ([]string, error) {
	tflog.Debug(ctx, "Validating Cloud PC role permissions and name uniqueness")

	// First validate permissions using static lookup
	validatedPermissions, err := validateCloudPCPermissionsStatic(ctx, permissions)
	if err != nil {
		return nil, err
	}

	// Then validate role name uniqueness using API call
	if displayName != "" {
		err = checkRoleNameUniqueness(ctx, client, displayName)
		if err != nil {
			return nil, err
		}
	}

	return validatedPermissions, nil
}

// validateCloudPCPermissionsStatic provides static validation for Cloud PC permissions
// using a list of known valid Cloud PC permissions extracted from API calls.
// Unlike the intune resource equivalent, there is no api to call for a dynamic list look up
// else we'd use that instead.
func validateCloudPCPermissionsStatic(ctx context.Context, permissions []string) ([]string, error) {
	tflog.Debug(ctx, "Performing static validation for Cloud PC permissions")

	// Known valid Cloud PC permissions from the API documentation
	validCloudPCOperations := map[string]bool{
		"Microsoft.CloudPC/CloudPCs/Read":                                true,
		"Microsoft.CloudPC/CloudPCs/Reprovision":                         true,
		"Microsoft.CloudPC/CloudPCs/Resize":                              true,
		"Microsoft.CloudPC/CloudPCs/EndGracePeriod":                      true,
		"Microsoft.CloudPC/CloudPCs/Restore":                             true,
		"Microsoft.CloudPC/CloudPCs/Reboot":                              true,
		"Microsoft.CloudPC/CloudPCs/Rename":                              true,
		"Microsoft.CloudPC/CloudPCs/Troubleshoot":                        true,
		"Microsoft.CloudPC/CloudPCs/ModifyDiskEncryptionType":            true,
		"Microsoft.CloudPC/CloudPCs/ChangeUserAccountType":               true,
		"Microsoft.CloudPC/CloudPCs/PlaceUnderReview":                    true,
		"Microsoft.CloudPC/CloudPCs/RetryPartnerAgentInstallation":       true,
		"Microsoft.CloudPC/CloudPCs/ApplyCurrentProvisioningPolicy":      true,
		"Microsoft.CloudPC/CloudPCs/CreateSnapshot":                      true,
		"Microsoft.CloudPC/CloudPCs/PowerOn":                             true,
		"Microsoft.CloudPC/CloudPCs/PowerOff":                            true,
		"Microsoft.CloudPC/CloudPCs/DisasterRecoveryFailover":            true,
		"Microsoft.CloudPC/CloudPCs/DisasterRecoveryFailback":            true,
		"Microsoft.CloudPC/CloudPCs/Start":                               true,
		"Microsoft.CloudPC/CloudPCs/Stop":                                true,
		"Microsoft.CloudPC/CloudPCs/GetCloudPcLaunchInfo":                true,
		"Microsoft.CloudPC/CloudPCs/ReinstallAgent":                      true,
		"Microsoft.CloudPC/CloudPCs/CheckAgentStatus":                    true,
		"Microsoft.CloudPC/CloudPCs/RetrieveAgentStatus":                 true,
		"Microsoft.CloudPC/CloudPCs/Provision":                           true,
		"Microsoft.CloudPC/CloudPCs/Deprovision":                         true,
		"Microsoft.CloudPC/DeviceImages/Create":                          true,
		"Microsoft.CloudPC/DeviceImages/Delete":                          true,
		"Microsoft.CloudPC/DeviceImages/Read":                            true,
		"Microsoft.CloudPC/OnPremisesConnections/Create":                 true,
		"Microsoft.CloudPC/OnPremisesConnections/Delete":                 true,
		"Microsoft.CloudPC/OnPremisesConnections/Read":                   true,
		"Microsoft.CloudPC/OnPremisesConnections/Update":                 true,
		"Microsoft.CloudPC/OnPremisesConnections/RunHealthChecks":        true,
		"Microsoft.CloudPC/OnPremisesConnections/UpdateAdDomainPassword": true,
		"Microsoft.CloudPC/ProvisioningPolicies/Assign":                  true,
		"Microsoft.CloudPC/ProvisioningPolicies/Apply":                   true,
		"Microsoft.CloudPC/ProvisioningPolicies/Create":                  true,
		"Microsoft.CloudPC/ProvisioningPolicies/Delete":                  true,
		"Microsoft.CloudPC/ProvisioningPolicies/Read":                    true,
		"Microsoft.CloudPC/ProvisioningPolicies/Update":                  true,
		"Microsoft.CloudPC/UserSettings/Assign":                          true,
		"Microsoft.CloudPC/UserSettings/Create":                          true,
		"Microsoft.CloudPC/UserSettings/Delete":                          true,
		"Microsoft.CloudPC/UserSettings/Read":                            true,
		"Microsoft.CloudPC/UserSettings/Update":                          true,
		"Microsoft.CloudPC/Roles/Read":                                   true,
		"Microsoft.CloudPC/Roles/Create":                                 true,
		"Microsoft.CloudPC/Roles/Update":                                 true,
		"Microsoft.CloudPC/Roles/Delete":                                 true,
		"Microsoft.CloudPC/RoleAssignments/Create":                       true,
		"Microsoft.CloudPC/RoleAssignments/Update":                       true,
		"Microsoft.CloudPC/RoleAssignments/Delete":                       true,
		"Microsoft.CloudPC/AuditData/Read":                               true,
		"Microsoft.CloudPC/SupportedRegion/Read":                         true,
		"Microsoft.CloudPC/ServicePlan/Read":                             true,
		"Microsoft.CloudPC/Snapshot/Read":                                true,
		"Microsoft.CloudPC/Snapshot/Share":                               true,
		"Microsoft.CloudPC/Snapshot/Import":                              true,
		"Microsoft.CloudPC/Snapshot/PurgeImportedSnapshot":               true,
		"Microsoft.CloudPC/OrganizationSettings/Read":                    true,
		"Microsoft.CloudPC/OrganizationSettings/Update":                  true,
		"Microsoft.CloudPC/ExternalPartnerSettings/Read":                 true,
		"Microsoft.CloudPC/ExternalPartnerSettings/Create":               true,
		"Microsoft.CloudPC/ExternalPartnerSettings/Update":               true,
		"Microsoft.CloudPC/PerformanceReports/Read":                      true,
		"Microsoft.CloudPC/SharedUseServicePlans/Read":                   true,
		"Microsoft.CloudPC/FrontLineServicePlans/Read":                   true,
		"Microsoft.CloudPC/SharedUseLicenseUsageReports/Read":            true,
		"Microsoft.CloudPC/FrontlineReports/Read":                        true,
		"Microsoft.CloudPC/CrossRegionDisasterRecovery/Read":             true,
		"Microsoft.CloudPC/BulkActions/Read":                             true,
		"Microsoft.CloudPC/BulkActions/Write":                            true,
		"Microsoft.CloudPC/ActionStatus/Read":                            true,
		"Microsoft.CloudPC/InaccessibleReports/Read":                     true,
		"Microsoft.CloudPC/MaintenanceWindows/Assign":                    true,
		"Microsoft.CloudPC/MaintenanceWindows/Create":                    true,
		"Microsoft.CloudPC/MaintenanceWindows/Delete":                    true,
		"Microsoft.CloudPC/MaintenanceWindows/Read":                      true,
		"Microsoft.CloudPC/MaintenanceWindows/Update":                    true,
		"Microsoft.CloudPC/DeviceRecommendation/Read":                    true,
		"Microsoft.CloudPC/CloudApps/Read":                               true,
		"Microsoft.CloudPC/CloudApps/Publish":                            true,
		"Microsoft.CloudPC/CloudApps/Update":                             true,
		"Microsoft.CloudPC/CloudApps/Reset":                              true,
		"Microsoft.CloudPC/CloudApps/Unpublish":                          true,
		"Microsoft.CloudPC/Settings/Assign":                              true,
		"Microsoft.CloudPC/Settings/Create":                              true,
		"Microsoft.CloudPC/Settings/Read":                                true,
		"Microsoft.CloudPC/Settings/Update":                              true,
		"Microsoft.CloudPC/Settings/Delete":                              true,
		"Microsoft.CloudPC/AdminHighlights/Operate":                      true,
	}

	var invalidPermissions []string
	for _, permission := range permissions {
		if !validCloudPCOperations[permission] {
			invalidPermissions = append(invalidPermissions, permission)
		}
	}

	if len(invalidPermissions) > 0 {
		validOperationsList := make([]string, 0, len(validCloudPCOperations))
		for operation := range validCloudPCOperations {
			validOperationsList = append(validOperationsList, operation)
		}
		return nil, fmt.Errorf("invalid Cloud PC resource operation(s) %v. Valid operations include: %v", invalidPermissions, validOperationsList[:10]) // Show first 10 for readability
	}

	tflog.Debug(ctx, fmt.Sprintf("Static validation passed for %d Cloud PC permissions", len(permissions)))
	return permissions, nil
}

// checkRoleNameUniqueness verifies that the role name is unique among existing roles
func checkRoleNameUniqueness(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, displayName string) error {
	tflog.Debug(ctx, fmt.Sprintf("Checking if role definition with display name '%s' already exists", displayName))

	existingRoles, err := client.
		RoleManagement().
		CloudPC().
		RoleDefinitions().
		Get(ctx, nil)

	if err != nil {
		return fmt.Errorf("failed to retrieve existing role definitions: %v", err)
	}

	roles := existingRoles.GetValue()
	for _, role := range roles {
		if role.GetDisplayName() != nil && *role.GetDisplayName() == displayName {
			return fmt.Errorf("a role definition with the display name '%s' already exists - role names must be unique", displayName)
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Role definition name '%s' is unique", displayName))
	return nil
}
