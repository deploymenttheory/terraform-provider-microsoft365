package graphBetaManagedDeviceCleanupRule

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// validateRequest ensures only one cleanup rule exists per platform in the tenant.
// If a rule already exists matching the requested platform (excluding the resource being updated), an error is returned.
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, platform graphmodels.DeviceCleanupRulePlatformType, excludeResourceID *string) error {
	tflog.Debug(ctx, "Validating uniqueness of managed device cleanup rule per platform", map[string]any{
		"platform": platform.String(),
	})

	if client == nil {
		return fmt.Errorf("graph client is not configured")
	}

	list, err := client.
		DeviceManagement().
		ManagedDeviceCleanupRules().
		Get(ctx, nil)

	if err != nil {
		return fmt.Errorf("failed to retrieve existing managed device cleanup rules for validation: %s", err.Error())
	}

	if list == nil || list.GetValue() == nil {
		return nil
	}

	for _, existing := range list.GetValue() {
		if existing == nil || existing.GetDeviceCleanupRulePlatformType() == nil {
			continue
		}

		if *existing.GetDeviceCleanupRulePlatformType() == platform {
			// If updating, allow the same resource ID
			if excludeResourceID != nil && existing.GetId() != nil && *existing.GetId() == *excludeResourceID {
				continue
			}
			return fmt.Errorf("only one of this resource can exist at once in this tenant")
		}
	}

	return nil
}
