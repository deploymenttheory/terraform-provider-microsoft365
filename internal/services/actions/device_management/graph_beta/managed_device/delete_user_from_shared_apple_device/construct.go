package graphBetaDeleteUserFromSharedAppleDevice

import (
	"context"

	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

// constructManagedDeviceRequest builds the request body for deleting a user from a managed device
func constructManagedDeviceRequest(ctx context.Context, userPrincipalName string) *devicemanagement.ManagedDevicesItemDeleteUserFromSharedAppleDevicePostRequestBody {
	requestBody := devicemanagement.NewManagedDevicesItemDeleteUserFromSharedAppleDevicePostRequestBody()
	requestBody.SetUserPrincipalName(&userPrincipalName)
	return requestBody
}

// constructComanagedDeviceRequest builds the request body for deleting a user from a co-managed device
func constructComanagedDeviceRequest(ctx context.Context, userPrincipalName string) *devicemanagement.ComanagedDevicesItemDeleteUserFromSharedAppleDevicePostRequestBody {
	requestBody := devicemanagement.NewComanagedDevicesItemDeleteUserFromSharedAppleDevicePostRequestBody()
	requestBody.SetUserPrincipalName(&userPrincipalName)
	return requestBody
}
