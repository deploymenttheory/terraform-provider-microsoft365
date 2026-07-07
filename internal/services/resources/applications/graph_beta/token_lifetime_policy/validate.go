package graphBetaApplicationsTokenLifetimePolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/sentinels"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// validateRequest validates the token lifetime policy request before creation or update.
// It checks for duplicate display names as the API only allows one policy per display name.
// The excludeID parameter allows skipping validation for a specific resource ID (used during updates).
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, displayName string, excludeID string) error {
	tflog.Debug(ctx, fmt.Sprintf("Validating %s request for display_name: %s (excludeID: %s)", ResourceName, displayName, excludeID))

	if displayName == "" {
		return sentinels.ErrEmptyDisplayName
	}

	policies, err := client.Policies().TokenLifetimePolicies().Get(ctx, nil)
	if err != nil {
		tflog.Warn(ctx, "Failed to list token lifetime policies for validation, skipping duplicate name check", map[string]any{
			"error": err.Error(),
		})
		// Don't fail validation if we can't fetch policies - allow operation to proceed
		// The API will return an error if there's an actual conflict
		return nil
	}

	if policies == nil || policies.GetValue() == nil {
		tflog.Debug(ctx, "No existing token lifetime policies found")
		return nil
	}

	for _, policy := range policies.GetValue() {
		if policy.GetId() == nil || policy.GetDisplayName() == nil {
			continue
		}

		policyID := *policy.GetId()
		existingDisplayName := *policy.GetDisplayName()

		// Skip validation for the resource being updated
		if excludeID != "" && policyID == excludeID {
			tflog.Debug(ctx, fmt.Sprintf("Skipping validation for current resource ID: %s", excludeID))
			continue
		}

		if existingDisplayName == displayName {
			return fmt.Errorf(
				"%w: a token lifetime policy with display_name '%s' already exists (ID: %s). "+
					"Only one policy per display_name is allowed. Please use a different display_name or import the existing policy",
				sentinels.ErrDuplicateDisplayName,
				displayName,
				policyID,
			)
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Validation passed: no duplicate display_name found for '%s'", displayName))
	return nil
}
