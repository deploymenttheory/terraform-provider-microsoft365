package graphBetaApplicationsTokenLifetimePolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/sentinels"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// validateRequest checks for duplicate display names before create/update operations
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, displayName string, excludeID string) error {
	if displayName == "" {
		return sentinels.ErrEmptyDisplayName
	}

	policies, err := client.Policies().TokenLifetimePolicies().Get(ctx, nil)
	if err != nil {
		tflog.Warn(ctx, "Failed to list token lifetime policies for validation, proceeding without duplicate check", map[string]any{
			"error": err.Error(),
		})
		return nil
	}

	if policies == nil || policies.GetValue() == nil {
		return nil
	}

	for _, policy := range policies.GetValue() {
		if policy.GetId() != nil && *policy.GetId() == excludeID {
			continue
		}
		if policy.GetDisplayName() != nil && *policy.GetDisplayName() == displayName {
			return fmt.Errorf("%w: display name '%s' is already used by policy with ID '%s'. Use terraform import to manage this resource",
				sentinels.ErrDuplicateDisplayName, displayName, *policy.GetId())
		}
	}

	tflog.Debug(ctx, "Token lifetime policy display name validation passed", map[string]any{
		"display_name": displayName,
	})
	return nil
}
