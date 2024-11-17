package graphBetaAssignmentFilter

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/construct"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model
func constructResource(ctx context.Context, data *AssignmentFilterResourceModel) (*models.DeviceAndAppManagementAssignmentFilter, error) {
	tflog.Debug(ctx, "Constructing 'resource_name' resource from model")

	requestBody := models.NewDeviceAndAppManagementAssignmentFilter()

	displayName := data.DisplayName.ValueString()
	requestBody.SetDisplayName(&displayName)

	// add more fields here as needed

	if err := construct.DebugLogGraphObject(ctx, "Final JSON to be sent to Graph API", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}
