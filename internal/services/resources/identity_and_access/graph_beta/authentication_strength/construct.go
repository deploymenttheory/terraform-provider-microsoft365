package graphBetaAuthenticationStrength

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// constructResource converts the Terraform resource model to a plain map for JSON marshaling
// Returns a map[string]any that can be directly JSON marshaled by the HTTP client
func constructResource(ctx context.Context, data *AuthenticationStrengthResourceModel) (map[string]any, error) {

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

	// Always include an empty combinationConfigurations array
	requestBody["combinationConfigurations"] = []interface{}{}

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
