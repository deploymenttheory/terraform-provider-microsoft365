package graphBetaAuthenticationStrengthPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// constructAllowedCombinationsUpdate constructs the request body for updating allowedCombinations
// POST /policies/authenticationStrengthPolicies/{id}/updateAllowedCombinations
func constructAllowedCombinationsUpdate(ctx context.Context, plan *AuthenticationStrengthPolicyResourceModel) (map[string]any, error) {
	var allowedCombinations []string
	if err := convert.FrameworkToGraphStringSet(ctx, plan.AllowedCombinations, func(values []string) {
		allowedCombinations = values
	}); err != nil {
		return nil, fmt.Errorf("failed to convert allowedCombinations: %w", err)
	}

	requestBody := map[string]any{
		"allowedCombinations": allowedCombinations,
	}

	return requestBody, nil
}

// constructCombinationConfigurationUpdate constructs the request body for updating a single combination configuration
// PATCH /identity/conditionalAccess/authenticationStrength/policies/{id}/combinationConfigurations/{configId}
func constructCombinationConfigurationUpdate(ctx context.Context, config *CombinationConfigurationModel) (map[string]any, error) {
	requestBody := make(map[string]any)

	// Add appliesToCombinations (required)
	var appliesToCombinations []string
	if err := convert.FrameworkToGraphStringSet(ctx, config.AppliesToCombinations, func(values []string) {
		appliesToCombinations = values
	}); err != nil {
		return nil, fmt.Errorf("failed to convert appliesToCombinations: %w", err)
	}
	requestBody["appliesToCombinations"] = appliesToCombinations

	// Add type-specific fields if present
	if !config.AllowedAAGUIDs.IsNull() && !config.AllowedAAGUIDs.IsUnknown() {
		var allowedAAGUIDs []string
		if err := convert.FrameworkToGraphStringSet(ctx, config.AllowedAAGUIDs, func(values []string) {
			allowedAAGUIDs = values
		}); err != nil {
			return nil, fmt.Errorf("failed to convert allowedAAGUIDs: %w", err)
		}
		requestBody["allowedAAGUIDs"] = allowedAAGUIDs
	}

	if !config.AllowedIssuerSkis.IsNull() && !config.AllowedIssuerSkis.IsUnknown() {
		var allowedIssuerSkis []string
		if err := convert.FrameworkToGraphStringSet(ctx, config.AllowedIssuerSkis, func(values []string) {
			allowedIssuerSkis = values
		}); err != nil {
			return nil, fmt.Errorf("failed to convert allowedIssuerSkis: %w", err)
		}
		requestBody["allowedIssuerSkis"] = allowedIssuerSkis
	}

	if !config.AllowedIssuers.IsNull() && !config.AllowedIssuers.IsUnknown() {
		var allowedIssuers []string
		if err := convert.FrameworkToGraphStringSet(ctx, config.AllowedIssuers, func(values []string) {
			allowedIssuers = values
		}); err != nil {
			return nil, fmt.Errorf("failed to convert allowedIssuers: %w", err)
		}
		requestBody["allowedIssuers"] = allowedIssuers
	}

	if !config.AllowedPolicyOIDs.IsNull() && !config.AllowedPolicyOIDs.IsUnknown() {
		var allowedPolicyOIDs []string
		if err := convert.FrameworkToGraphStringSet(ctx, config.AllowedPolicyOIDs, func(values []string) {
			allowedPolicyOIDs = values
		}); err != nil {
			return nil, fmt.Errorf("failed to convert allowedPolicyOIDs: %w", err)
		}
		requestBody["allowedPolicyOIDs"] = allowedPolicyOIDs
	}

	return requestBody, nil
}

// combinationConfigurationsEqual checks if two combination configurations are equal
func combinationConfigurationsEqual(ctx context.Context, a, b *CombinationConfigurationModel) bool {
	return a.ODataType.Equal(b.ODataType) &&
		a.AppliesToCombinations.Equal(b.AppliesToCombinations) &&
		a.AllowedAAGUIDs.Equal(b.AllowedAAGUIDs) &&
		a.AllowedIssuerSkis.Equal(b.AllowedIssuerSkis) &&
		a.AllowedIssuers.Equal(b.AllowedIssuers) &&
		a.AllowedPolicyOIDs.Equal(b.AllowedPolicyOIDs)
}
