package graphCloudPcDeviceImage

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/construct"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model
func constructResource(ctx context.Context, data *CloudPcDeviceImageResourceModel) (*models.CloudPcDeviceImage, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := models.NewCloudPcDeviceImage()

	if !data.DisplayName.IsNull() {
		displayName := data.DisplayName.ValueString()
		requestBody.SetDisplayName(&displayName)
	}

	if !data.SourceImageResourceId.IsNull() {
		sourceImageResourceId := data.SourceImageResourceId.ValueString()
		requestBody.SetSourceImageResourceId(&sourceImageResourceId)
	}

	if !data.Version.IsNull() {
		version := data.Version.ValueString()
		requestBody.SetVersion(&version)
	}

	if err := construct.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
