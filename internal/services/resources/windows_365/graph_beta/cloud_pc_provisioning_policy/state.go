package graphBetaCloudPcProvisioningPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteStateToTerraform(ctx context.Context, data *CloudPcProvisioningPolicyResourceModel, remoteResource models.CloudPcProvisioningPolicyable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]interface{}{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()).ValueString(),
	})

	// Set basic properties
	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.CloudPcNamingTemplate = convert.GraphToFrameworkString(remoteResource.GetCloudPcNamingTemplate())
	data.AlternateResourceUrl = convert.GraphToFrameworkString(remoteResource.GetAlternateResourceUrl())
	data.CloudPcGroupDisplayName = convert.GraphToFrameworkString(remoteResource.GetCloudPcGroupDisplayName())
	data.EnableSingleSignOn = convert.GraphToFrameworkBool(remoteResource.GetEnableSingleSignOn())
	data.GracePeriodInHours = convert.GraphToFrameworkInt32(remoteResource.GetGracePeriodInHours())
	data.ImageDisplayName = convert.GraphToFrameworkString(remoteResource.GetImageDisplayName())
	data.ImageId = convert.GraphToFrameworkString(remoteResource.GetImageId())
	data.ImageType = convert.GraphToFrameworkEnum(remoteResource.GetImageType())
	data.LocalAdminEnabled = convert.GraphToFrameworkBool(remoteResource.GetLocalAdminEnabled())
	data.ProvisioningType = convert.GraphToFrameworkEnum(remoteResource.GetProvisioningType())
	data.ManagedBy = convert.GraphToFrameworkEnum(remoteResource.GetManagedBy())
	data.ScopeIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetScopeIds())

	// Handle Autopatch
	if autopatch := remoteResource.GetAutopatch(); autopatch != nil {
		data.Autopatch = &AutopatchModel{
			AutopatchGroupId: convert.GraphToFrameworkString(autopatch.GetAutopatchGroupId()),
		}
	} else {
		data.Autopatch = nil
	}

	// Handle AutopilotConfiguration
	if apc := remoteResource.GetAutopilotConfiguration(); apc != nil {
		data.AutopilotConfiguration = &AutopilotConfigurationModel{
			DevicePreparationProfileId:  convert.GraphToFrameworkString(apc.GetDevicePreparationProfileId()),
			ApplicationTimeoutInMinutes: convert.GraphToFrameworkInt32(apc.GetApplicationTimeoutInMinutes()),
			OnFailureDeviceAccessDenied: convert.GraphToFrameworkBool(apc.GetOnFailureDeviceAccessDenied()),
		}
	} else {
		data.AutopilotConfiguration = nil
	}

	// Handle Microsoft Managed Desktop
	if mmd := remoteResource.GetMicrosoftManagedDesktop(); mmd != nil {
		data.MicrosoftManagedDesktop = &MicrosoftManagedDesktopModel{
			ManagedType: convert.GraphToFrameworkEnum(mmd.GetManagedType()),
			Profile:     convert.GraphToFrameworkString(mmd.GetProfile()),
		}
	} else {
		data.MicrosoftManagedDesktop = nil
	}

	// Handle Domain Join Configurations
	if domainJoinConfigs := remoteResource.GetDomainJoinConfigurations(); domainJoinConfigs != nil {
		data.DomainJoinConfigurations = make([]DomainJoinConfigurationModel, len(domainJoinConfigs))
		for i, config := range domainJoinConfigs {
			data.DomainJoinConfigurations[i] = DomainJoinConfigurationModel{
				DomainJoinType:         convert.GraphToFrameworkEnum(config.GetDomainJoinType()),
				OnPremisesConnectionId: convert.GraphToFrameworkString(config.GetOnPremisesConnectionId()),
				RegionName:             convert.GraphToFrameworkString(config.GetRegionName()),
				RegionGroup:            convert.GraphToFrameworkEnum(config.GetRegionGroup()),
			}
		}
	} else {
		data.DomainJoinConfigurations = []DomainJoinConfigurationModel{}
	}

	// Handle Windows Settings
	if windowsSetting := remoteResource.GetWindowsSetting(); windowsSetting != nil {
		data.WindowsSetting = &WindowsSettingModel{
			Locale: convert.GraphToFrameworkString(windowsSetting.GetLocale()),
		}
	} else {
		data.WindowsSetting = nil
	}

	// // Handle Assignments
	// if assignments := remoteResource.GetAssignments(); assignments != nil {
	// 	tflog.Debug(ctx, fmt.Sprintf("Found %d assignments for policy", len(assignments)))
	// 	data.Assignments = MapAssignmentsToTerraform(ctx, assignments)
	// } else {
	// 	tflog.Debug(ctx, "No assignments found for policy")
	// 	data.Assignments = []CloudPcProvisioningPolicyAssignmentModel{}
	// }

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))
}
