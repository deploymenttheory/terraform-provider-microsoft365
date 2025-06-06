package graphBetaDeviceEnrollmentLimitConfiguration

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform resource model to the Graph API request model
func constructResource(ctx context.Context, data *DeviceEnrollmentLimitConfigurationResourceModel) (models.DeviceEnrollmentLimitConfigurationable, error) {
	requestBody := models.NewDeviceEnrollmentLimitConfiguration()

	constructors.SetStringProperty(data.DisplayName, requestBody.SetDisplayName)
	constructors.SetStringProperty(data.Description, requestBody.SetDescription)
	constructors.SetInt32Property(data.Priority, requestBody.SetPriority)
	constructors.SetInt32Property(data.Limit, requestBody.SetLimit)

	if err := constructors.SetStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
