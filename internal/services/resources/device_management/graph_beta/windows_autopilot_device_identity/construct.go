package graphBetaWindowsAutopilotDeviceIdentity

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs a Windows Autopilot Device Identity object for API requests
func constructResource(ctx context.Context, data *WindowsAutopilotDeviceIdentityResourceModel, forUpdate bool) (graphmodels.WindowsAutopilotDeviceIdentityable, error) {
	resource := graphmodels.NewWindowsAutopilotDeviceIdentity()

	convert.FrameworkToGraphString(data.SerialNumber, resource.SetSerialNumber)
	convert.FrameworkToGraphString(data.GroupTag, resource.SetGroupTag)
	convert.FrameworkToGraphString(data.PurchaseOrderIdentifier, resource.SetPurchaseOrderIdentifier)
	convert.FrameworkToGraphString(data.ProductKey, resource.SetProductKey)
	convert.FrameworkToGraphString(data.Manufacturer, resource.SetManufacturer)
	convert.FrameworkToGraphString(data.Model, resource.SetModel)
	convert.FrameworkToGraphString(data.DisplayName, resource.SetDisplayName)
	convert.FrameworkToGraphString(data.UserPrincipalName, resource.SetUserPrincipalName)

	if err := constructors.DebugLogGraphObject(ctx, "Constructed Windows Autopilot Device Identity Resource", resource); err != nil {
		tflog.Error(ctx, "Failed to log Windows Autopilot Device Identity", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return resource, nil
}
