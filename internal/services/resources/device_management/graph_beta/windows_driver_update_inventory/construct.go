// Main entry point to construct the intune windows driver update inventory resource for the Terraform provider.
package graphBetaWindowsDriverUpdateInventory

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Main entry point to construct the intune windows driver update inventory resource for the Terraform provider.
func constructResource(ctx context.Context, data *WindowsDriverUpdateInventoryResourceModel) (graphmodels.WindowsDriverUpdateInventoryable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewWindowsDriverUpdateInventory()

	convert.FrameworkToGraphString(data.Name, requestBody.SetName)
	convert.FrameworkToGraphString(data.Version, requestBody.SetVersion)
	convert.FrameworkToGraphString(data.Manufacturer, requestBody.SetManufacturer)
	convert.FrameworkToGraphString(data.DriverClass, requestBody.SetDriverClass)

	if err := convert.FrameworkToGraphEnum(data.ApprovalStatus, graphmodels.ParseDriverApprovalStatus, requestBody.SetApprovalStatus); err != nil {
		return nil, fmt.Errorf("invalid approval status: %s", err)
	}

	if err := convert.FrameworkToGraphEnum(data.Category, graphmodels.ParseDriverCategory, requestBody.SetCategory); err != nil {
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
