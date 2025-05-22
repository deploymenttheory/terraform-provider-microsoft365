package graphBetaManagedDeviceCleanupRule

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs and returns a ManagedDeviceCleanupRule
func constructResource(ctx context.Context, data ManagedDeviceCleanupRuleResourceModel) (graphmodels.ManagedDeviceCleanupRuleable, error) {
	tflog.Debug(ctx, "Starting managed device cleanup rule construction")

	rule := graphmodels.NewManagedDeviceCleanupRule()

	constructors.SetStringProperty(data.DisplayName, rule.SetDisplayName)
	constructors.SetStringProperty(data.Description, rule.SetDescription)
	constructors.SetInt32Property(data.DeviceInactivityBeforeRetirementInDays, rule.SetDeviceInactivityBeforeRetirementInDays)

	// Set the platform type enum property
	if !data.DeviceCleanupRulePlatformType.IsNull() && !data.DeviceCleanupRulePlatformType.IsUnknown() {
		err := constructors.SetEnumProperty(
			data.DeviceCleanupRulePlatformType,
			graphmodels.ParseDeviceCleanupRulePlatformType,
			func(val *graphmodels.DeviceCleanupRulePlatformType) {
				rule.SetDeviceCleanupRulePlatformType(val)
			},
		)
		if err != nil {
			return nil, fmt.Errorf("error setting device cleanup rule platform type: %v", err)
		}
	}

	if err := constructors.DebugLogGraphObject(ctx, "Constructed managed device cleanup rule", rule); err != nil {
		tflog.Error(ctx, "Failed to log managed device cleanup rule", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return rule, nil
}
