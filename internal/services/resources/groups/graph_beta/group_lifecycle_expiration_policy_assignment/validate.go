package graphBetaGroupLifecycleExpirationPolicyAssignment

import (
	"context"
	"fmt"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// validateRequest performs validation before creating or deleting a policy assignment.
func (r *GroupLifecycleExpirationPolicyAssignmentResource) validateRequest(
	ctx context.Context,
	groupID string,
	diagnostics *diag.Diagnostics,
) (policyID string, err error) {
	tflog.Debug(ctx, "Starting request validation", map[string]any{
		"groupId": groupID,
	})

	// Validate tenant policy exists
	policyID, err = r.validateTenantPolicyExists(ctx, diagnostics)
	if err != nil || policyID == "" {
		return "", err
	}

	// Validate group is Microsoft 365 type
	err = r.validateGroupIsOfTypeM365(ctx, groupID, diagnostics)
	if err != nil {
		return "", err
	}

	tflog.Debug(ctx, "Request validation completed", map[string]any{
		"groupId":  groupID,
		"policyId": policyID,
	})

	return policyID, nil
}

// validateTenantPolicyExists validates that a lifecycle policy exists in the tenant.
func (r *GroupLifecycleExpirationPolicyAssignmentResource) validateTenantPolicyExists(
	ctx context.Context,
	diagnostics *diag.Diagnostics,
) (string, error) {
	tflog.Debug(ctx, "Validating tenant lifecycle policy exists")

	policies, err := r.client.GroupLifecyclePolicies().Get(ctx, nil)
	if err != nil {
		tflog.Error(ctx, "Failed to retrieve lifecycle policies", map[string]any{
			"error": err.Error(),
		})
		diagnostics.AddError(
			"Failed to get lifecycle policy",
			fmt.Sprintf("Error retrieving lifecycle policy: %s", err.Error()),
		)
		return "", err
	}

	if policies == nil || policies.GetValue() == nil || len(policies.GetValue()) == 0 {
		tflog.Error(ctx, "No lifecycle policy found in tenant")
		diagnostics.AddError(
			"No lifecycle policy found",
			"No lifecycle policy exists in the tenant. You must create a group lifecycle policy before assigning groups to it.",
		)
		return "", fmt.Errorf("no lifecycle policy found")
	}

	policy := policies.GetValue()[0]
	policyID := policy.GetId()
	if policyID == nil {
		tflog.Error(ctx, "Lifecycle policy ID is null")
		diagnostics.AddError(
			"Invalid policy ID",
			"The lifecycle policy ID is null",
		)
		return "", fmt.Errorf("policy ID is null")
	}

	managedGroupTypes := policy.GetManagedGroupTypes()
	if managedGroupTypes == nil || *managedGroupTypes != "Selected" {
		tflog.Error(ctx, "Lifecycle policy managedGroupTypes is not set to Selected", map[string]any{
			"policyId":          *policyID,
			"managedGroupTypes": managedGroupTypes,
		})
		diagnostics.AddError(
			"Invalid lifecycle policy configuration",
			fmt.Sprintf("The lifecycle policy must have managedGroupTypes set to 'Selected' to assign individual groups. Current value: %v. "+
				"This resource is only applicable when the policy is configured to manage selected groups, not all groups.", managedGroupTypes),
		)
		return "", fmt.Errorf("lifecycle policy managedGroupTypes is not 'Selected'")
	}

	tflog.Debug(ctx, "Tenant lifecycle policy validated", map[string]any{
		"policyId":          *policyID,
		"managedGroupTypes": *managedGroupTypes,
	})

	return *policyID, nil
}

// validateGroupIsOfTypeM365 validates that the group is a Microsoft 365 (Unified) group.
func (r *GroupLifecycleExpirationPolicyAssignmentResource) validateGroupIsOfTypeM365(
	ctx context.Context,
	groupID string,
	diagnostics *diag.Diagnostics,
) error {
	tflog.Debug(ctx, "Validating group is Microsoft 365 type", map[string]any{
		"groupId": groupID,
	})

	group, err := r.client.
		Groups().
		ByGroupId(groupID).
		Get(ctx, nil)

	if err != nil {
		errorInfo := errors.GraphError(ctx, err)
		tflog.Error(ctx, "Failed to retrieve group", map[string]any{
			"groupId":    groupID,
			"statusCode": errorInfo.StatusCode,
			"errorCode":  errorInfo.ErrorCode,
		})
		diagnostics.AddError(
			"Failed to retrieve group",
			fmt.Sprintf("Error retrieving group %s: %s", groupID, err.Error()),
		)
		return err
	}

	if group == nil {
		tflog.Error(ctx, "Group not found", map[string]any{
			"groupId": groupID,
		})
		diagnostics.AddError(
			"Group not found",
			fmt.Sprintf("Group %s does not exist", groupID),
		)
		return fmt.Errorf("group %s not found", groupID)
	}

	groupTypes := group.GetGroupTypes()
	isM365Group := false
	for _, groupType := range groupTypes {
		if groupType == "Unified" {
			isM365Group = true
			break
		}
	}

	if !isM365Group {
		tflog.Error(ctx, "Group is not a Microsoft 365 group", map[string]any{
			"groupId":    groupID,
			"groupTypes": groupTypes,
		})
		diagnostics.AddError(
			"Invalid group type",
			fmt.Sprintf("Group %s is not a Microsoft 365 group. Only Microsoft 365 (Unified) groups can be assigned to lifecycle policies. Group types: %v", groupID, groupTypes),
		)
		return fmt.Errorf("group %s is not a Microsoft 365 group", groupID)
	}

	tflog.Debug(ctx, "Group validated as Microsoft 365 type", map[string]any{
		"groupId":    groupID,
		"groupTypes": groupTypes,
	})

	return nil
}
