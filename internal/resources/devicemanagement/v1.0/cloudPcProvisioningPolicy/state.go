package graphCloudPcProvisioningPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

func mapRemoteStateToTerraform(ctx context.Context, data *CloudPcProvisioningPolicyResourceModel, remoteResource models.CloudPcProvisioningPolicyable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": state.StringPtrToString(remoteResource.GetId()),
	})

	data.ID = types.StringValue(state.StringPtrToString(remoteResource.GetId()))
	data.DisplayName = types.StringValue(state.StringPtrToString(remoteResource.GetDisplayName()))
	data.Description = types.StringValue(state.StringPtrToString(remoteResource.GetDescription()))
	data.CloudPcNamingTemplate = types.StringValue(state.StringPtrToString(remoteResource.GetCloudPcNamingTemplate()))
	data.AlternateResourceUrl = types.StringValue(state.StringPtrToString(remoteResource.GetAlternateResourceUrl()))
	data.CloudPcGroupDisplayName = types.StringValue(state.StringPtrToString(remoteResource.GetCloudPcGroupDisplayName()))
	data.EnableSingleSignOn = state.BoolPtrToTypeBool(remoteResource.GetEnableSingleSignOn())
	data.GracePeriodInHours = state.Int32PtrToTypeInt64(remoteResource.GetGracePeriodInHours())
	data.ImageDisplayName = types.StringValue(state.StringPtrToString(remoteResource.GetImageDisplayName()))
	data.ImageId = types.StringValue(state.StringPtrToString(remoteResource.GetImageId()))
	data.ImageType = state.EnumPtrToTypeString(remoteResource.GetImageType())
	data.LocalAdminEnabled = state.BoolPtrToTypeBool(remoteResource.GetLocalAdminEnabled())
	data.ProvisioningType = state.EnumPtrToTypeString(remoteResource.GetProvisioningType())

	if mmd := remoteResource.GetMicrosoftManagedDesktop(); mmd != nil {
		data.MicrosoftManagedDesktop = &MicrosoftManagedDesktopModel{
			ManagedType: state.EnumPtrToTypeString(mmd.GetManagedType()),
			Profile:     types.StringValue(state.StringPtrToString(mmd.GetProfile())),
		}
	} else {
		data.MicrosoftManagedDesktop = nil
	}

	if domainJoinConfigs := remoteResource.GetDomainJoinConfigurations(); domainJoinConfigs != nil {
		data.DomainJoinConfigurations = make([]DomainJoinConfigurationModel, len(domainJoinConfigs))
		for i, config := range domainJoinConfigs {
			data.DomainJoinConfigurations[i] = DomainJoinConfigurationModel{
				DomainJoinType:         state.EnumPtrToTypeString(config.GetDomainJoinType()),
				OnPremisesConnectionId: types.StringValue(state.StringPtrToString(config.GetOnPremisesConnectionId())),
				RegionName:             types.StringValue(state.StringPtrToString(config.GetRegionName())),
			}
		}
	} else {
		data.DomainJoinConfigurations = []DomainJoinConfigurationModel{}
	}

	if windowsSetting := remoteResource.GetWindowsSetting(); windowsSetting != nil {
		data.WindowsSetting = &WindowsSettingModel{
			Locale: types.StringValue(state.StringPtrToString(windowsSetting.GetLocale())),
		}
	} else {
		data.WindowsSetting = nil
	}

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state", map[string]interface{}{
		"resourceId": data.ID.ValueString(),
	})
}
