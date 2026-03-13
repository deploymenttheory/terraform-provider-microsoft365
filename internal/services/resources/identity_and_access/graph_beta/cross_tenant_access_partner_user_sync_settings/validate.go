package graphBetaCrossTenantAccessPartnerUserSyncSettings

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// validateRequest validates the partner user sync settings request
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data *CrossTenantAccessPartnerUserSyncSettingsResourceModel) error {
	tflog.Debug(ctx, "Starting partner user sync settings request validation")

	tenantID := data.TenantID.ValueString()
	if err := validateMicrosoftEntraOrganization(ctx, client, tenantID); err != nil {
		return fmt.Errorf("validation failed for tenant_id: %w", err)
	}

	tflog.Debug(ctx, "Partner user sync settings request validation completed successfully")
	return nil
}

// validateMicrosoftEntraOrganization validates a single tenant ID by checking if it is a valid Microsoft Entra organization
func validateMicrosoftEntraOrganization(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, tenantID string) error {
	tflog.Debug(ctx, fmt.Sprintf("Validating tenant ID: %s", tenantID))

	tenantInfo, err := getTenantInformationByTenantID(ctx, client, tenantID)
	if err != nil {
		tflog.Warn(ctx, fmt.Sprintf("Error validating tenant ID %s: %v", tenantID, err))
		return fmt.Errorf("invalid Microsoft Entra organization tenant ID: %s. Please verify this tenant ID is valid", tenantID)
	}

	displayName := "Unknown"
	domainName := "Unknown"

	if tenantInfo.GetDisplayName() != nil {
		displayName = *tenantInfo.GetDisplayName()
	}

	if tenantInfo.GetDefaultDomainName() != nil {
		domainName = *tenantInfo.GetDefaultDomainName()
	}

	tflog.Debug(ctx, "Validated tenant information", map[string]any{
		"tenant_id":   tenantID,
		"name":        displayName,
		"domain":      domainName,
	})

	tflog.Debug(ctx, "Tenant ID validation completed successfully")
	return nil
}

// getTenantInformationByTenantID retrieves tenant information from Microsoft Graph
func getTenantInformationByTenantID(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, tenantID string) (graphmodels.TenantInformationable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Fetching tenant information for tenant ID: %s", tenantID))

	tenantInfo, err := client.
		TenantRelationships().
		FindTenantInformationByTenantIdWithTenantId(&tenantID).
		Get(ctx, nil)

	if err != nil {
		return nil, fmt.Errorf("tenant ID validation failed: %w", err)
	}

	return tenantInfo, nil
}
