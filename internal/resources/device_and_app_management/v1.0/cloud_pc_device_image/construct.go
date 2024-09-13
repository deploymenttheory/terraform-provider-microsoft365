package graphCloudPcDeviceImage

import (
	"context"
	"encoding/json"

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

	// Debug logging
	debugPrintRequestBody(ctx, requestBody)

	return requestBody, nil
}

func debugPrintRequestBody(ctx context.Context, requestBody *models.CloudPcDeviceImage) {
	requestMap := map[string]interface{}{
		"displayName":           requestBody.GetDisplayName(),
		"sourceImageResourceId": requestBody.GetSourceImageResourceId(),
		"version":               requestBody.GetVersion(),
	}

	requestBodyJSON, err := json.MarshalIndent(requestMap, "", "  ")
	if err != nil {
		tflog.Error(ctx, "Error marshalling request body to JSON", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	tflog.Debug(ctx, "Constructed Cloud PC Device Image resource", map[string]interface{}{
		"requestBody": string(requestBodyJSON),
	})
}
