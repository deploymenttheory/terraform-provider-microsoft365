package graphBetaCleanWindowsManagedDevice

import (
	"context"

	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

// constructRequest builds the request body for the clean Windows device action
func constructManagedDeviceRequest(ctx context.Context, device ManagedDeviceCleanWindows) *devicemanagement.ManagedDevicesItemCleanWindowsDevicePostRequestBody {
	requestBody := devicemanagement.NewManagedDevicesItemCleanWindowsDevicePostRequestBody()

	keepUserData := device.KeepUserData.ValueBool()
	requestBody.SetKeepUserData(&keepUserData)

	return requestBody
}

// constructComanagedDeviceRequest builds the request body for the clean Windows co-managed device action
func constructComanagedDeviceRequest(ctx context.Context, device ComanagedDeviceCleanWindows) *devicemanagement.ComanagedDevicesItemCleanWindowsDevicePostRequestBody {
	requestBody := devicemanagement.NewComanagedDevicesItemCleanWindowsDevicePostRequestBody()

	keepUserData := device.KeepUserData.ValueBool()
	requestBody.SetKeepUserData(&keepUserData)

	return requestBody
}
