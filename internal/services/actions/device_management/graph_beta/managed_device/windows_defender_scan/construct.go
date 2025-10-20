package graphBetaWindowsDefenderScan

import (
	"context"

	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

// constructManagedDeviceRequest builds the request body for managed device Windows Defender scan
func constructManagedDeviceRequest(ctx context.Context, quickScan bool) *devicemanagement.ManagedDevicesItemWindowsDefenderScanPostRequestBody {
	requestBody := devicemanagement.NewManagedDevicesItemWindowsDefenderScanPostRequestBody()
	requestBody.SetQuickScan(&quickScan)
	return requestBody
}

// constructComanagedDeviceRequest builds the request body for co-managed device Windows Defender scan
func constructComanagedDeviceRequest(ctx context.Context, quickScan bool) *devicemanagement.ComanagedDevicesItemWindowsDefenderScanPostRequestBody {
	requestBody := devicemanagement.NewComanagedDevicesItemWindowsDefenderScanPostRequestBody()
	requestBody.SetQuickScan(&quickScan)
	return requestBody
}
