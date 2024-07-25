package graphCloudPcProvisioningPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

// mapRemoteStateToTerraform maps the remote state from the Graph API to the Terraform resource model.
// It populates the CloudPcProvisioningPolicyResourceModel with data from the CloudPcProvisioningPolicy.
func mapRemoteStateToTerraform(ctx context.Context, data *CloudPcProvisioningPolicyResourceModel, remoteResource models.CloudPcProvisioningPolicyable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	data.ID = types.StringValue(helpers.StringPtrToString(remoteResource.GetId()))
	data.DisplayName = types.StringValue(helpers.StringPtrToString(remoteResource.GetDisplayName()))
	data.Description = types.StringValue(helpers.StringPtrToString(remoteResource.GetDescription()))
	data.CloudPcNamingTemplate = types.StringValue(helpers.StringPtrToString(remoteResource.GetCloudPcNamingTemplate()))
	data.AlternateResourceUrl = types.StringValue(helpers.StringPtrToString(remoteResource.GetAlternateResourceUrl()))
	data.CloudPcGroupDisplayName = types.StringValue(helpers.StringPtrToString(remoteResource.GetCloudPcGroupDisplayName()))

	if enableSSO := remoteResource.GetEnableSingleSignOn(); enableSSO != nil {
		data.EnableSingleSignOn = types.BoolValue(*enableSSO)
	} else {
		data.EnableSingleSignOn = types.BoolNull()
	}

	if gracePeriod := remoteResource.GetGracePeriodInHours(); gracePeriod != nil {
		data.GracePeriodInHours = types.Int64Value(int64(*gracePeriod))
	} else {
		data.GracePeriodInHours = types.Int64Null()
	}

	data.ImageDisplayName = types.StringValue(helpers.StringPtrToString(remoteResource.GetImageDisplayName()))
	data.ImageId = types.StringValue(helpers.StringPtrToString(remoteResource.GetImageId()))

	if imageType := remoteResource.GetImageType(); imageType != nil {
		data.ImageType = types.StringValue(imageType.String())
	} else {
		data.ImageType = types.StringNull()
	}

	if localAdmin := remoteResource.GetLocalAdminEnabled(); localAdmin != nil {
		data.LocalAdminEnabled = types.BoolValue(*localAdmin)
	} else {
		data.LocalAdminEnabled = types.BoolNull()
	}

	if provisioningType := remoteResource.GetProvisioningType(); provisioningType != nil {
		data.ProvisioningType = types.StringValue(provisioningType.String())
	} else {
		data.ProvisioningType = types.StringNull()
	}

	if mmd := remoteResource.GetMicrosoftManagedDesktop(); mmd != nil {
		managedType := ""
		if mt := mmd.GetManagedType(); mt != nil {
			managedType = mt.String()
		}

		data.MicrosoftManagedDesktop = &MicrosoftManagedDesktopModel{
			ManagedType: types.StringValue(managedType),
			Profile:     types.StringValue(helpers.StringPtrToString(mmd.GetProfile())),
		}
	} else {
		data.MicrosoftManagedDesktop = nil
	}

	if domainJoinConfigs := remoteResource.GetDomainJoinConfigurations(); domainJoinConfigs != nil {
		data.DomainJoinConfigurations = make([]DomainJoinConfigurationModel, len(domainJoinConfigs))
		for i, config := range domainJoinConfigs {
			domainJoinType := ""
			if djt := config.GetDomainJoinType(); djt != nil {
				domainJoinType = djt.String()
			}

			data.DomainJoinConfigurations[i] = DomainJoinConfigurationModel{
				DomainJoinType:         types.StringValue(domainJoinType),
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

	tflog.Debug(ctx, "Finished mapping remote state to Terraform")
}
