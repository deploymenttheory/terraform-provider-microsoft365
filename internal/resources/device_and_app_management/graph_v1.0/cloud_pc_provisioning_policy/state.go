package graphCloudPcProvisioningPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

func MapRemoteStateToTerraform(ctx context.Context, data *CloudPcProvisioningPolicyResourceModel, remoteResource models.CloudPcProvisioningPolicyable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": state.StringPtrToString(remoteResource.GetId()),
	})

	// Set basic properties
	data.ID = types.StringPointerValue(remoteResource.GetId())
	data.DisplayName = types.StringPointerValue(remoteResource.GetDisplayName())
	data.Description = types.StringPointerValue(remoteResource.GetDescription())
	data.CloudPcNamingTemplate = types.StringPointerValue(remoteResource.GetCloudPcNamingTemplate())
	data.AlternateResourceUrl = types.StringPointerValue(remoteResource.GetAlternateResourceUrl())
	data.CloudPcGroupDisplayName = types.StringPointerValue(remoteResource.GetCloudPcGroupDisplayName())
	data.EnableSingleSignOn = types.BoolPointerValue(remoteResource.GetEnableSingleSignOn())
	data.GracePeriodInHours = state.Int32PtrToTypeInt64(remoteResource.GetGracePeriodInHours())
	data.ImageDisplayName = types.StringPointerValue(remoteResource.GetImageDisplayName())
	data.ImageId = types.StringPointerValue(remoteResource.GetImageId())
	data.ImageType = state.EnumPtrToTypeString(remoteResource.GetImageType())
	data.LocalAdminEnabled = types.BoolPointerValue(remoteResource.GetLocalAdminEnabled())
	data.ProvisioningType = state.EnumPtrToTypeString(remoteResource.GetProvisioningType())

	// Handle Microsoft Managed Desktop
	if mmd := remoteResource.GetMicrosoftManagedDesktop(); mmd != nil {
		data.MicrosoftManagedDesktop = &MicrosoftManagedDesktopModel{
			ManagedType: state.EnumPtrToTypeString(mmd.GetManagedType()),
			Profile:     types.StringPointerValue(mmd.GetProfile()),
		}
	} else {
		data.MicrosoftManagedDesktop = nil
	}

	// Handle Domain Join Configurations
	if domainJoinConfigs := remoteResource.GetDomainJoinConfigurations(); domainJoinConfigs != nil {
		data.DomainJoinConfigurations = make([]DomainJoinConfigurationModel, len(domainJoinConfigs))
		for i, config := range domainJoinConfigs {
			data.DomainJoinConfigurations[i] = DomainJoinConfigurationModel{
				DomainJoinType:         state.EnumPtrToTypeString(config.GetDomainJoinType()),
				OnPremisesConnectionId: types.StringPointerValue(config.GetOnPremisesConnectionId()),
				RegionName:             types.StringPointerValue(config.GetRegionName()),
			}
		}
	} else {
		data.DomainJoinConfigurations = []DomainJoinConfigurationModel{}
	}

	// Handle Windows Settings
	if windowsSetting := remoteResource.GetWindowsSetting(); windowsSetting != nil {
		data.WindowsSetting = &WindowsSettingModel{
			Locale: types.StringPointerValue(windowsSetting.GetLocale()),
		}
	} else {
		data.WindowsSetting = nil
	}

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state", map[string]interface{}{
		"resourceId": data.ID.ValueString(),
	})
}
