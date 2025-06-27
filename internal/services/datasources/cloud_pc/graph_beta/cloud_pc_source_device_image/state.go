package cloudPcSourceDeviceImage

import (
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// This file is reserved for future state mapping helpers if needed.

// MapRemoteStateToDataSource maps a CloudPcSourceDeviceImageable to the Terraform item model.
func MapRemoteStateToDataSource(img models.CloudPcSourceDeviceImageable) CloudPcSourceDeviceImageItemModel {
	return CloudPcSourceDeviceImageItemModel{
		ID:                      convert.GraphToFrameworkString(img.GetId()),
		ResourceId:              convert.GraphToFrameworkString(img.GetResourceId()),
		DisplayName:             convert.GraphToFrameworkString(img.GetDisplayName()),
		SubscriptionId:          convert.GraphToFrameworkString(img.GetSubscriptionId()),
		SubscriptionDisplayName: convert.GraphToFrameworkString(img.GetSubscriptionDisplayName()),
	}
}
