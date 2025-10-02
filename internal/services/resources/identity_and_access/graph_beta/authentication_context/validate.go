package graphBetaAuthenticationContext

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// validateRequest validates that the proposed authentication context ID is not already in use
// Only runs validation during create operations when isCreate is true
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data *AuthenticationContextResourceModel, isCreate bool) error {
	if !isCreate {
		tflog.Debug(ctx, "Skipping validation for update operation", map[string]any{
			"proposedId": data.ID.ValueString(),
		})
		return nil
	}

	tflog.Debug(ctx, "Starting validation of authentication context ID for create operation", map[string]any{
		"proposedId": data.ID.ValueString(),
	})

	authContexts, err := client.
		Identity().
		ConditionalAccess().
		AuthenticationContextClassReferences().
		Get(ctx, nil)

	if err != nil {
		return fmt.Errorf("could not retrieve authentication context class references for validation: %w", err)
	}

	if authContexts == nil || authContexts.GetValue() == nil {
		tflog.Debug(ctx, "No existing authentication contexts found")
		return nil
	}

	proposedId := data.ID.ValueString()
	existingContexts := authContexts.GetValue()

	for _, context := range existingContexts {
		if context.GetId() != nil && *context.GetId() == proposedId {
			tflog.Error(ctx, "Authentication context ID already exists", map[string]any{
				"proposedId": proposedId,
				"existingId": *context.GetId(),
			})
			return fmt.Errorf("authentication context class reference with ID '%s' already exists", proposedId)
		}
	}

	tflog.Debug(ctx, "Authentication context ID validation passed", map[string]any{
		"proposedId":    proposedId,
		"existingCount": len(existingContexts),
	})

	return nil
}
