package graphBetaCloudPC

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/state"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps a cloud PC to a model
func MapRemoteStateToDataSource(ctx context.Context, data graphmodels.CloudPCable) CloudPCItemDataSourceModel {
	model := CloudPCItemDataSourceModel{
		Id:                       state.StringPointerValue(data.GetId()),
		AadDeviceId:              state.StringPointerValue(data.GetAadDeviceId()),
		DisplayName:              state.StringPointerValue(data.GetDisplayName()),
		ImageDisplayName:         state.StringPointerValue(data.GetImageDisplayName()),
		ManagedDeviceId:          state.StringPointerValue(data.GetManagedDeviceId()),
		ManagedDeviceName:        state.StringPointerValue(data.GetManagedDeviceName()),
		ProvisioningPolicyId:     state.StringPointerValue(data.GetProvisioningPolicyId()),
		ProvisioningPolicyName:   state.StringPointerValue(data.GetProvisioningPolicyName()),
		OnPremisesConnectionName: state.StringPointerValue(data.GetOnPremisesConnectionName()),
		ServicePlanId:            state.StringPointerValue(data.GetServicePlanId()),
		ServicePlanName:          state.StringPointerValue(data.GetServicePlanName()),
		ServicePlanType:          state.EnumPtrToTypeString(data.GetServicePlanType()),
		Status:                   state.EnumPtrToTypeString(data.GetStatus()),
		UserPrincipalName:        state.StringPointerValue(data.GetUserPrincipalName()),
		LastModifiedDateTime:     state.TimeToString(data.GetLastModifiedDateTime()),
		//StatusDetails:            state.StringPointerValue(data.GetStatusDetails()),
		GracePeriodEndDateTime: state.TimeToString(data.GetGracePeriodEndDateTime()),
		ProvisioningType:       state.EnumPtrToTypeString(data.GetProvisioningType()),
		DeviceRegionName:       state.StringPointerValue(data.GetDeviceRegionName()),
		DiskEncryptionState:    state.EnumPtrToTypeString(data.GetDiskEncryptionState()),
	}

	return model
}
