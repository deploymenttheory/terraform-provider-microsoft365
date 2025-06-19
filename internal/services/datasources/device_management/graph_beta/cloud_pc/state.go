package graphBetaCloudPC

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps a cloud PC to a model
func MapRemoteStateToDataSource(ctx context.Context, data graphmodels.CloudPCable) CloudPCItemDataSourceModel {
	model := CloudPCItemDataSourceModel{
		Id:                       convert.GraphToFrameworkString(data.GetId()),
		AadDeviceId:              convert.GraphToFrameworkString(data.GetAadDeviceId()),
		DisplayName:              convert.GraphToFrameworkString(data.GetDisplayName()),
		ImageDisplayName:         convert.GraphToFrameworkString(data.GetImageDisplayName()),
		ManagedDeviceId:          convert.GraphToFrameworkString(data.GetManagedDeviceId()),
		ManagedDeviceName:        convert.GraphToFrameworkString(data.GetManagedDeviceName()),
		ProvisioningPolicyId:     convert.GraphToFrameworkString(data.GetProvisioningPolicyId()),
		ProvisioningPolicyName:   convert.GraphToFrameworkString(data.GetProvisioningPolicyName()),
		OnPremisesConnectionName: convert.GraphToFrameworkString(data.GetOnPremisesConnectionName()),
		ServicePlanId:            convert.GraphToFrameworkString(data.GetServicePlanId()),
		ServicePlanName:          convert.GraphToFrameworkString(data.GetServicePlanName()),
		ServicePlanType:          convert.GraphToFrameworkEnum(data.GetServicePlanType()),
		Status:                   convert.GraphToFrameworkEnum(data.GetStatus()),
		UserPrincipalName:        convert.GraphToFrameworkString(data.GetUserPrincipalName()),
		LastModifiedDateTime:     convert.GraphToFrameworkTime(data.GetLastModifiedDateTime()),
		//StatusDetails:            convert.GraphToFrameworkString(data.GetStatusDetails()),
		GracePeriodEndDateTime: convert.GraphToFrameworkTime(data.GetGracePeriodEndDateTime()),
		ProvisioningType:       convert.GraphToFrameworkEnum(data.GetProvisioningType()),
		DeviceRegionName:       convert.GraphToFrameworkString(data.GetDeviceRegionName()),
		DiskEncryptionState:    convert.GraphToFrameworkEnum(data.GetDiskEncryptionState()),
	}

	return model
}
