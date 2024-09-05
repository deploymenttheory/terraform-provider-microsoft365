package graphCloudPcDeviceImage

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model
func constructResource(ctx context.Context, data *CloudPcDeviceImageResourceModel) (*models.CloudPcDeviceImage, error) {
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

	// Convert the request body to JSON for logging
	requestBodyJSON, err := json.MarshalIndent(map[string]interface{}{
		"displayName":           requestBody.GetDisplayName(),
		"sourceImageResourceId": requestBody.GetSourceImageResourceId(),
		"version":               requestBody.GetVersion(),
	}, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("error marshalling request body to JSON: %s", err)
	}

	tflog.Debug(ctx, "Constructed Cloud PC Device Image resource:\n"+string(requestBodyJSON))

	return requestBody, nil
}
