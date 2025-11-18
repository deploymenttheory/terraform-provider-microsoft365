package graphBetaAuthenticationStrengthPolicy

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// constructResource converts the Terraform resource model to a plain map for JSON marshaling
// Returns a map[string]any that can be directly JSON marshaled by the HTTP client
func constructResource(ctx context.Context, data *AuthenticationStrengthPolicyResourceModel) (map[string]any, error) {

	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := make(map[string]any)

	// Basic properties using convert helpers
	convert.FrameworkToGraphString(data.DisplayName, func(val *string) {
		if val != nil {
			requestBody["displayName"] = *val
		}
	})

	convert.FrameworkToGraphString(data.Description, func(val *string) {
		if val != nil {
			requestBody["description"] = *val
		}
	})

	// Convert allowed combinations to array
	if err := convert.FrameworkToGraphStringSet(ctx, data.AllowedCombinations, func(values []string) {
		if len(values) > 0 {
			requestBody["allowedCombinations"] = values
		}
	}); err != nil {
		return nil, fmt.Errorf("failed to convert allowed combinations: %w", err)
	}

	// Construct combination configurations
	combinationConfigs, err := constructCombinationConfigurations(ctx, data)
	if err != nil {
		return nil, fmt.Errorf("failed to construct combination configurations: %w", err)
	}
	requestBody["combinationConfigurations"] = combinationConfigs

	// Debug logging using plain JSON marshal
	if debugJSON, err := json.MarshalIndent(requestBody, "", "    "); err == nil {
		tflog.Debug(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), map[string]any{
			"json": "\n" + string(debugJSON),
		})
	} else {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructCombinationConfigurations builds the combinationConfigurations array from the Terraform model
func constructCombinationConfigurations(ctx context.Context, data *AuthenticationStrengthPolicyResourceModel) ([]any, error) {
	if data.CombinationConfigurations.IsNull() || data.CombinationConfigurations.IsUnknown() {
		return []any{}, nil
	}

	var configurations []CombinationConfigurationModel
	diags := data.CombinationConfigurations.ElementsAs(ctx, &configurations, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to parse combination configurations: %v", diags)
	}

	result := make([]any, 0, len(configurations))

	for _, config := range configurations {
		// Use ordered map to ensure correct field order for API
		// Microsoft Graph API requires: @odata.type, appliesToCombinations, then type-specific fields
		configMap := make(map[string]any)

		// Collect all fields first
		var odataType string
		var id string
		var appliesToCombinations []string
		var allowedAAGUIDs []string
		var allowedIssuerSkis []string
		var allowedIssuers []string
		var allowedPolicyOIDs []string

		if !config.ODataType.IsNull() && !config.ODataType.IsUnknown() {
			odataType = config.ODataType.ValueString()
		}

		if !config.ID.IsNull() && !config.ID.IsUnknown() {
			id = config.ID.ValueString()
		}

		if err := convert.FrameworkToGraphStringSet(ctx, config.AppliesToCombinations, func(values []string) {
			appliesToCombinations = values
		}); err != nil {
			return nil, fmt.Errorf("failed to convert appliesToCombinations: %w", err)
		}

		if err := convert.FrameworkToGraphStringSet(ctx, config.AllowedAAGUIDs, func(values []string) {
			allowedAAGUIDs = values
		}); err != nil {
			return nil, fmt.Errorf("failed to convert allowedAAGUIDs: %w", err)
		}

		if err := convert.FrameworkToGraphStringSet(ctx, config.AllowedIssuerSkis, func(values []string) {
			allowedIssuerSkis = values
		}); err != nil {
			return nil, fmt.Errorf("failed to convert allowedIssuerSkis: %w", err)
		}

		if err := convert.FrameworkToGraphStringSet(ctx, config.AllowedIssuers, func(values []string) {
			allowedIssuers = values
		}); err != nil {
			return nil, fmt.Errorf("failed to convert allowedIssuers: %w", err)
		}

		if err := convert.FrameworkToGraphStringSet(ctx, config.AllowedPolicyOIDs, func(values []string) {
			allowedPolicyOIDs = values
		}); err != nil {
			return nil, fmt.Errorf("failed to convert allowedPolicyOIDs: %w", err)
		}

		// Build map in correct order (though Go maps don't preserve order, this helps readability)
		// The actual ordering will depend on JSON marshaling
		if odataType != "" {
			configMap["@odata.type"] = odataType
		}

		if id != "" {
			configMap["id"] = id
		}

		// IMPORTANT: appliesToCombinations must come before type-specific fields in the API request
		if len(appliesToCombinations) > 0 {
			configMap["appliesToCombinations"] = appliesToCombinations
		}

		// Type-specific fields
		if len(allowedAAGUIDs) > 0 {
			configMap["allowedAAGUIDs"] = allowedAAGUIDs
		}

		if len(allowedIssuerSkis) > 0 {
			configMap["allowedIssuerSkis"] = allowedIssuerSkis
		}

		if len(allowedIssuers) > 0 {
			configMap["allowedIssuers"] = allowedIssuers
		}

		if len(allowedPolicyOIDs) > 0 {
			configMap["allowedPolicyOIDs"] = allowedPolicyOIDs
		}

		result = append(result, configMap)
	}

	return result, nil
}
