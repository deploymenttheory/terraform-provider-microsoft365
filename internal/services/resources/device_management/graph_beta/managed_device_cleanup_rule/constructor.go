package graphBetaManagedDeviceCleanupRule

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs and returns a ManagedDeviceCleanupRule
func constructResource(ctx context.Context, data ManagedDeviceCleanupRuleResourceModel) (graphmodels.ManagedDeviceCleanupRuleable, error) {
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
	}

	if err := constructors.DebugLogGraphObject(ctx, "Constructed managed device cleanup rule", rule); err != nil {
		tflog.Error(ctx, "Failed to log managed device cleanup rule", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return rule, nil
}
