// Main entry point to construct the intune windows driver update inventory resource for the Terraform provider.
package graphBetaWindowsDriverUpdateInventory

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Main entry point to construct the intune windows driver update inventory resource for the Terraform provider.
func constructResource(ctx context.Context, data *WindowsDriverUpdateInventoryResourceModel) (graphmodels.WindowsDriverUpdateInventoryable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewWindowsDriverUpdateInventory()

	constructors.SetStringProperty(data.Name, requestBody.SetName)
	constructors.SetStringProperty(data.Version, requestBody.SetVersion)
	constructors.SetStringProperty(data.Manufacturer, requestBody.SetManufacturer)
	constructors.SetStringProperty(data.DriverClass, requestBody.SetDriverClass)

	if err := constructors.SetEnumProperty(data.ApprovalStatus, graphmodels.ParseDriverApprovalStatus, requestBody.SetApprovalStatus); err != nil {
		return nil, fmt.Errorf("invalid approval status: %s", err)
	}

	if err := constructors.SetEnumProperty(data.Category, graphmodels.ParseDriverCategory, requestBody.SetCategory); err != nil {
		return nil, fmt.Errorf("invalid category: %s", err)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
