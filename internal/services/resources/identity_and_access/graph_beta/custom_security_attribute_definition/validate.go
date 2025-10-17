package graphBetaCustomSecurityAttributeDefinition

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

const (
	// MaxActiveCustomSecurityAttributes is the maximum number of active custom security attributes allowed in a tenant
	// Reference: https://learn.microsoft.com/en-us/entra/fundamentals/custom-security-attributes-limits
	MaxActiveCustomSecurityAttributes = 500
)

// validateRequest validates the custom security attribute definition request
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data *CustomSecurityAttributeDefinitionResourceModel) error {
	tflog.Debug(ctx, "Starting custom security attribute definition request validation")

	if err := validateActiveAttributeLimit(ctx, client); err != nil {
		return fmt.Errorf("validation failed for active attribute limit: %w", err)
	}

	tflog.Debug(ctx, "Custom security attribute definition request validation completed successfully")
	return nil
}

// validateActiveAttributeLimit checks that the tenant doesn't exceed the limit of 500 active custom security attributes
func validateActiveAttributeLimit(ctx context.Context, client *msgraphbetasdk.GraphServiceClient) error {
	tflog.Debug(ctx, "Validating active custom security attribute count limit")

	activeCount, err := getActiveCustomSecurityAttributeCount(ctx, client)
	if err != nil {
		return fmt.Errorf("failed to retrieve active custom security attribute count: %w", err)
	}

	tflog.Debug(ctx, fmt.Sprintf("Current active custom security attribute count: %d", activeCount))

	if activeCount >= MaxActiveCustomSecurityAttributes {
		return fmt.Errorf("cannot create new custom security attribute: tenant has reached the maximum limit of %d active custom security attributes (current count: %d). You can define up to 500 active objects in a tenant",
			MaxActiveCustomSecurityAttributes, activeCount)
	}

	tflog.Debug(ctx, fmt.Sprintf("Active custom security attribute count is within limit: %d/%d", activeCount, MaxActiveCustomSecurityAttributes))
	return nil
}

// getActiveCustomSecurityAttributeCount retrieves the count of active custom security attributes from Microsoft Graph API
func getActiveCustomSecurityAttributeCount(ctx context.Context, client *msgraphbetasdk.GraphServiceClient) (int, error) {
	tflog.Debug(ctx, "Fetching custom security attribute definitions from Microsoft Graph API")

	result, err := client.
		Directory().
		CustomSecurityAttributeDefinitions().
		Get(ctx, nil)

	if err != nil {
		return 0, fmt.Errorf("error fetching custom security attribute definitions: %w", err)
	}

	if result == nil || result.GetValue() == nil {
		tflog.Debug(ctx, "No custom security attribute definitions found")
		return 0, nil
	}

	activeCount := 0
	for _, attr := range result.GetValue() {
		if attr.GetStatus() != nil && *attr.GetStatus() == "Available" {
			activeCount++
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully fetched custom security attribute definitions: %d total, %d active", len(result.GetValue()), activeCount))
	return activeCount, nil
}
