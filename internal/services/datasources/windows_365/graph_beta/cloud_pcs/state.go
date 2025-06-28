package graphBetaCloudPcs

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteStateToDataSource(ctx context.Context, data graphmodels.CloudPCable) CloudPcItemModel {
	model := CloudPcItemModel{
		ID:                       convert.GraphToFrameworkString(data.GetId()),
		DisplayName:              convert.GraphToFrameworkString(data.GetDisplayName()),
		AADDeviceID:              convert.GraphToFrameworkString(data.GetAadDeviceId()),
		ImageDisplayName:         convert.GraphToFrameworkString(data.GetImageDisplayName()),
		ManagedDeviceID:          convert.GraphToFrameworkString(data.GetManagedDeviceId()),
		ManagedDeviceName:        convert.GraphToFrameworkString(data.GetManagedDeviceName()),
		ProvisioningPolicyID:     convert.GraphToFrameworkString(data.GetProvisioningPolicyId()),
		ProvisioningPolicyName:   convert.GraphToFrameworkString(data.GetProvisioningPolicyName()),
		OnPremisesConnectionName: convert.GraphToFrameworkString(data.GetOnPremisesConnectionName()),
		ServicePlanID:            convert.GraphToFrameworkString(data.GetServicePlanId()),
		ServicePlanName:          convert.GraphToFrameworkString(data.GetServicePlanName()),
		DeviceRegionName:         convert.GraphToFrameworkString(data.GetDeviceRegionName()),
		UserPrincipalName:        convert.GraphToFrameworkString(data.GetUserPrincipalName()),
	}

	// Handle enum types
	model.Status = convert.GraphToFrameworkEnum(data.GetStatus())
	model.ServicePlanType = convert.GraphToFrameworkEnum(data.GetServicePlanType())
	model.ProvisioningType = convert.GraphToFrameworkEnum(data.GetProvisioningType())
	model.DiskEncryptionState = convert.GraphToFrameworkEnum(data.GetDiskEncryptionState())
	model.ProductType = convert.GraphToFrameworkEnum(data.GetProductType())
	model.UserAccountType = convert.GraphToFrameworkEnum(data.GetUserAccountType())

	// Handle datetime fields
	model.LastModifiedDateTime = convert.GraphToFrameworkTime(data.GetLastModifiedDateTime())
	model.GracePeriodEndDateTime = convert.GraphToFrameworkTime(data.GetGracePeriodEndDateTime())

	// Handle status detail
	if statusDetail := data.GetStatusDetail(); statusDetail != nil {
		if code := statusDetail.GetCode(); code != nil {
			model.StatusDetailCode = convert.GraphToFrameworkString(code)
		}
		if message := statusDetail.GetMessage(); message != nil {
			model.StatusDetailMessage = convert.GraphToFrameworkString(message)
		}
	}

	// Handle connection setting
	if connectionSetting := data.GetConnectionSetting(); connectionSetting != nil {
		model.EnableSingleSignOn = convert.GraphToFrameworkBool(connectionSetting.GetEnableSingleSignOn())
	}

	return model
}
