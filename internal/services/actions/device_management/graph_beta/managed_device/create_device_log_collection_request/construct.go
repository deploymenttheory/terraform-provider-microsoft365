package graphBetaCreateDeviceLogCollectionRequestManagedDevice

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	msgraphbetamodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructManagedDeviceRequest(ctx context.Context, device ManagedDeviceLogCollection) *devicemanagement.ManagedDevicesItemCreateDeviceLogCollectionRequestPostRequestBody {
	requestBody := devicemanagement.NewManagedDevicesItemCreateDeviceLogCollectionRequestPostRequestBody()

	// Create the deviceLogCollectionRequest object
	logCollectionRequest := msgraphbetamodels.NewDeviceLogCollectionRequest()

	// Set template type if provided, otherwise use default "predefined"
	templateType := "predefined"
	if !device.TemplateType.IsNull() && !device.TemplateType.IsUnknown() {
		templateType = device.TemplateType.ValueString()
	}

	// Convert string to enum
	var templateTypeEnum msgraphbetamodels.DeviceLogCollectionTemplateType
	switch templateType {
	case "predefined":
		templateTypeEnum = msgraphbetamodels.PREDEFINED_DEVICELOGCOLLECTIONTEMPLATETYPE
	case "unknownFutureValue":
		templateTypeEnum = msgraphbetamodels.UNKNOWNFUTUREVALUE_DEVICELOGCOLLECTIONTEMPLATETYPE
	default:
		// Default to predefined
		templateTypeEnum = msgraphbetamodels.PREDEFINED_DEVICELOGCOLLECTIONTEMPLATETYPE
	}

	logCollectionRequest.SetTemplateType(&templateTypeEnum)
	requestBody.SetTemplateType(logCollectionRequest)

	if err := constructors.DebugLogGraphObject(ctx, "Final managed device log collection request", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}

	return requestBody
}

func constructComanagedDeviceRequest(ctx context.Context, device ComanagedDeviceLogCollection) *devicemanagement.ComanagedDevicesItemCreateDeviceLogCollectionRequestPostRequestBody {
	requestBody := devicemanagement.NewComanagedDevicesItemCreateDeviceLogCollectionRequestPostRequestBody()

	// Create the deviceLogCollectionRequest object
	logCollectionRequest := msgraphbetamodels.NewDeviceLogCollectionRequest()

	// Set template type if provided, otherwise use default "predefined"
	templateType := "predefined"
	if !device.TemplateType.IsNull() && !device.TemplateType.IsUnknown() {
		templateType = device.TemplateType.ValueString()
	}

	// Convert string to enum
	var templateTypeEnum msgraphbetamodels.DeviceLogCollectionTemplateType
	switch templateType {
	case "predefined":
		templateTypeEnum = msgraphbetamodels.PREDEFINED_DEVICELOGCOLLECTIONTEMPLATETYPE
	case "unknownFutureValue":
		templateTypeEnum = msgraphbetamodels.UNKNOWNFUTUREVALUE_DEVICELOGCOLLECTIONTEMPLATETYPE
	default:
		// Default to predefined
		templateTypeEnum = msgraphbetamodels.PREDEFINED_DEVICELOGCOLLECTIONTEMPLATETYPE
	}

	logCollectionRequest.SetTemplateType(&templateTypeEnum)
	requestBody.SetTemplateType(logCollectionRequest)

	if err := constructors.DebugLogGraphObject(ctx, "Final co-managed device log collection request", requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}

	return requestBody
}
