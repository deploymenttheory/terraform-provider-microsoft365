package graphBetaSettingsCatalogConfigurationPolicy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// listResourceAssignments checks if a policy has any assignments by querying the assignments endpoint.
// This is more reliable than using the isAssigned field which can be stale/incorrect in the API.
func (r *SettingsCatalogListResource) listResourceAssignments(ctx context.Context, policyId string) (bool, error) {
	tflog.Debug(ctx, fmt.Sprintf("Checking assignments for policy: %s", policyId))

	assignmentsResponse, err := r.client.
		DeviceManagement().
		ConfigurationPolicies().
		ByDeviceManagementConfigurationPolicyId(policyId).
		Assignments().
		Get(ctx, nil)

	if err != nil {
		tflog.Warn(ctx, fmt.Sprintf("Error checking assignments for policy %s: %v", policyId, err))
		return false, err
	}

	if assignmentsResponse == nil {
		tflog.Debug(ctx, fmt.Sprintf("Policy %s has nil assignments response", policyId))
		return false, nil
	}

	assignments := assignmentsResponse.GetValue()
	hasAssignments := len(assignments) > 0
	tflog.Debug(ctx, fmt.Sprintf("Policy %s has %d assignments (hasAssignments=%t)", policyId, len(assignments), hasAssignments))
	return hasAssignments, nil
}
