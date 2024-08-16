package graphCloudPcProvisioningPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
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
		"resourceId": helpers.StringPtrToString(remoteResource.GetId()),
	})

	data.ID = types.StringValue(helpers.StringPtrToString(remoteResource.GetId()))
	data.DisplayName = types.StringValue(helpers.StringPtrToString(remoteResource.GetDisplayName()))
	data.Description = types.StringValue(helpers.StringPtrToString(remoteResource.GetDescription()))
	data.CloudPcNamingTemplate = types.StringValue(helpers.StringPtrToString(remoteResource.GetCloudPcNamingTemplate()))
	data.AlternateResourceUrl = types.StringValue(helpers.StringPtrToString(remoteResource.GetAlternateResourceUrl()))
	data.CloudPcGroupDisplayName = types.StringValue(helpers.StringPtrToString(remoteResource.GetCloudPcGroupDisplayName()))
	data.EnableSingleSignOn = helpers.BoolPtrToTypeBool(remoteResource.GetEnableSingleSignOn())
	data.GracePeriodInHours = helpers.Int32PtrToTypeInt64(remoteResource.GetGracePeriodInHours())
	data.ImageDisplayName = types.StringValue(helpers.StringPtrToString(remoteResource.GetImageDisplayName()))
	data.ImageId = types.StringValue(helpers.StringPtrToString(remoteResource.GetImageId()))
	data.ImageType = helpers.EnumPtrToTypeString(remoteResource.GetImageType())
	data.LocalAdminEnabled = helpers.BoolPtrToTypeBool(remoteResource.GetLocalAdminEnabled())
	data.ProvisioningType = helpers.EnumPtrToTypeString(remoteResource.GetProvisioningType())

	if mmd := remoteResource.GetMicrosoftManagedDesktop(); mmd != nil {
		data.MicrosoftManagedDesktop = &MicrosoftManagedDesktopModel{
			ManagedType: helpers.EnumPtrToTypeString(mmd.GetManagedType()),
			Profile:     types.StringValue(helpers.StringPtrToString(mmd.GetProfile())),
		}
	} else {
		data.MicrosoftManagedDesktop = nil
	}

	if domainJoinConfigs := remoteResource.GetDomainJoinConfigurations(); domainJoinConfigs != nil {
		data.DomainJoinConfigurations = make([]DomainJoinConfigurationModel, len(domainJoinConfigs))
		for i, config := range domainJoinConfigs {
			data.DomainJoinConfigurations[i] = DomainJoinConfigurationModel{
				DomainJoinType:         helpers.EnumPtrToTypeString(config.GetDomainJoinType()),
				OnPremisesConnectionId: types.StringValue(helpers.StringPtrToString(config.GetOnPremisesConnectionId())),
				RegionName:             types.StringValue(helpers.StringPtrToString(config.GetRegionName())),
			}
		}
	} else {
		data.DomainJoinConfigurations = []DomainJoinConfigurationModel{}
	}

	if windowsSetting := remoteResource.GetWindowsSetting(); windowsSetting != nil {
		data.WindowsSetting = &WindowsSettingModel{
			Locale: types.StringValue(helpers.StringPtrToString(windowsSetting.GetLocale())),
		}
	} else {
		data.WindowsSetting = nil
	}

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state", map[string]interface{}{
		"resourceId": data.ID.ValueString(),
	})
}
