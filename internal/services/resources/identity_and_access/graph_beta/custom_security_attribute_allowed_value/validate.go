package graphBetaCustomSecurityAttributeAllowedValue

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	// MaxAllowedValuesPerDefinition is the maximum number of allowed values per custom security attribute definition
	// Reference: https://learn.microsoft.com/en-us/entra/fundamentals/custom-security-attributes-limits
	MaxAllowedValuesPerDefinition = 100
)

// validateRequest validates the allowed value request
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data *CustomSecurityAttributeAllowedValueResourceModel) error {
	tflog.Debug(ctx, "Starting allowed value request validation")

	if err := validateAllowedValueLimit(ctx, client, data.CustomSecurityAttributeDefinitionId.ValueString()); err != nil {
		return fmt.Errorf("validation failed for allowed value limit: %w", err)
	}

	tflog.Debug(ctx, "Allowed value request validation completed successfully")
	return nil
}

// validateAllowedValueLimit checks that the definition doesn't exceed the limit of 100 allowed values
func validateAllowedValueLimit(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, definitionId string) error {
	tflog.Debug(ctx, fmt.Sprintf("Validating allowed value count limit for definition: %s", definitionId))

	allowedValueCount, err := getAllowedValueCount(ctx, client, definitionId)
	if err != nil {
		return fmt.Errorf("failed to retrieve allowed value count: %w", err)
	}

	tflog.Debug(ctx, fmt.Sprintf("Current allowed value count for definition %s: %d", definitionId, allowedValueCount))

	if allowedValueCount >= MaxAllowedValuesPerDefinition {
		return fmt.Errorf("cannot create new allowed value: custom security attribute definition '%s' has reached the maximum limit of %d allowed values (current count: %d). You can define up to 100 allowedValue objects per customSecurityAttributeDefinition",
			definitionId, MaxAllowedValuesPerDefinition, allowedValueCount)
	}

	tflog.Debug(ctx, fmt.Sprintf("Allowed value count is within limit: %d/%d", allowedValueCount, MaxAllowedValuesPerDefinition))
	return nil
}

// getAllowedValueCount retrieves the count of allowed values for a specific custom security attribute definition
func getAllowedValueCount(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, definitionId string) (int, error) {
	tflog.Debug(ctx, fmt.Sprintf("Fetching allowed values for definition: %s", definitionId))

	result, err := client.
		Directory().
		CustomSecurityAttributeDefinitions().
		ByCustomSecurityAttributeDefinitionId(definitionId).
		AllowedValues().
		Get(ctx, nil)

	if err != nil {
		return 0, fmt.Errorf("error fetching allowed values for definition %s: %w", definitionId, err)
	}

	if result == nil || result.GetValue() == nil {
		tflog.Debug(ctx, fmt.Sprintf("No allowed values found for definition: %s", definitionId))
		return 0, nil
	}

	allowedValueCount := len(result.GetValue())
	tflog.Debug(ctx, fmt.Sprintf("Successfully fetched allowed values for definition %s: %d total", definitionId, allowedValueCount))

	return allowedValueCount, nil
}
