package graphConditionalAccessTermsOfUse

import (
	"context"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// validateRequest validates that the proposed agreement configuration is valid
// This is a placeholder for any future validation logic
func validateRequest(ctx context.Context, data *ConditionalAccessTermsOfUseResourceModel) error {
	tflog.Debug(ctx, "Starting validation of terms of use agreement", map[string]any{
		"displayName": data.DisplayName.ValueString(),
	})

	// For now, we don't have specific validation requirements
	// This function can be extended in the future if needed

	tflog.Debug(ctx, "Terms of use agreement validation passed")
	return nil
}
