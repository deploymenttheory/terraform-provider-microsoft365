package graphBetaManagedDeviceCleanupRule

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs and returns a ManagedDeviceCleanupRule
// Performs uniqueness validation per platform before returning the object
func constructResource(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data ManagedDeviceCleanupRuleResourceModel, isUpdate bool) (graphmodels.ManagedDeviceCleanupRuleable, error) {
	tflog.Debug(ctx, "Starting managed device cleanup rule construction")

	rule := graphmodels.NewManagedDeviceCleanupRule()

	convert.FrameworkToGraphString(data.DisplayName, rule.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, rule.SetDescription)
	convert.FrameworkToGraphInt32(data.DeviceInactivityBeforeRetirementInDays, rule.SetDeviceInactivityBeforeRetirementInDays)

	// Set the platform type enum property
	if !data.DeviceCleanupRulePlatformType.IsNull() && !data.DeviceCleanupRulePlatformType.IsUnknown() {
		err := convert.FrameworkToGraphEnum(
			data.DeviceCleanupRulePlatformType,
			graphmodels.ParseDeviceCleanupRulePlatformType,
			func(val *graphmodels.DeviceCleanupRulePlatformType) {
				rule.SetDeviceCleanupRulePlatformType(val)
			},
		)
		if err != nil {
			return nil, fmt.Errorf("error setting device cleanup rule platform type: %v", err)
		}
		// Perform uniqueness validation against tenant for this platform
		if client != nil && rule.GetDeviceCleanupRulePlatformType() != nil {
			var excludeResourceID *string
			if isUpdate && !data.ID.IsNull() && !data.ID.IsUnknown() {
				id := data.ID.ValueString()
				excludeResourceID = &id
			}
			if err := validateRequest(ctx, client, *rule.GetDeviceCleanupRulePlatformType(), excludeResourceID); err != nil {
				return nil, err
			}
		}
	}

	if err := constructors.DebugLogGraphObject(ctx, "Constructed managed device cleanup rule", rule); err != nil {
		tflog.Error(ctx, "Failed to log managed device cleanup rule", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return rule, nil
}
