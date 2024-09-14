package graphCloudPcDeviceImage

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/construct"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model
func constructResource(ctx context.Context, data *CloudPcDeviceImageResourceModel) (*models.CloudPcDeviceImage, error) {
	tflog.Debug(ctx, "Constructing CloudPcDeviceImage Resource")
	construct.DebugPrintStruct(ctx, "Constructed CloudPcDeviceImage Resource from model", data)

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

	return requestBody, nil
}
