package graphBetaCloudPcProvisioningPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructResource(ctx context.Context, data *CloudPcProvisioningPolicyResourceModel, client *msgraphbetasdk.GraphServiceClient) (*models.CloudPcProvisioningPolicy, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	if err := validateResource(ctx, client, data); err != nil {
		return nil, err
	}

	requestBody := models.NewCloudPcProvisioningPolicy()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphString(data.CloudPcNamingTemplate, requestBody.SetCloudPcNamingTemplate)
	convert.FrameworkToGraphString(data.ImageId, requestBody.SetImageId)

	// Set imageDisplayName based on imageId
	// these have been extracted from api calls. undocumented.
	if !data.ImageId.IsNull() && !data.ImageId.IsUnknown() {
		switch data.ImageId.ValueString() {
		case "microsoftwindowsdesktop_windows-ent-cpc_win11-25h2-ent-cpc":
			val := "Windows 11 Enterprise 25H2"
			requestBody.SetImageDisplayName(&val)
		case "microsoftwindowsdesktop_windows-ent-cpc_win11-25h2-ent-cpc-m365":
			val := "Windows 11 Enterprise + Microsoft 365 Apps 25H2"
			requestBody.SetImageDisplayName(&val)
		case "microsoftwindowsdesktop_windows-ent-cpc_win11-24H2-ent-cpc":
			val := "Windows 11 Enterprise 24H2"
			requestBody.SetImageDisplayName(&val)
		case "microsoftwindowsdesktop_windows-ent-cpc_win11-24H2-ent-cpc-m365":
			val := "Windows 11 Enterprise + Microsoft 365 Apps 24H2"
			requestBody.SetImageDisplayName(&val)
		case "microsoftwindowsdesktop_windows-ent-cpc_win11-23h2-ent-cpc":
			val := "Windows 11 Enterprise 23H2"
			requestBody.SetImageDisplayName(&val)
		case "microsoftwindowsdesktop_windows-ent-cpc_win11-23h2-ent-cpc-m365":
			val := "Windows 11 Enterprise + Microsoft 365 Apps 23H2"
			requestBody.SetImageDisplayName(&val)
		}
	}

	convert.FrameworkToGraphBool(data.EnableSingleSignOn, requestBody.SetEnableSingleSignOn)
	convert.FrameworkToGraphBool(data.LocalAdminEnabled, requestBody.SetLocalAdminEnabled)

	// Set ManagedBy if present
	if !data.ManagedBy.IsNull() && !data.ManagedBy.IsUnknown() && data.ManagedBy.ValueString() != "" {
		if err := convert.FrameworkToGraphEnum(data.ManagedBy, models.ParseCloudPcManagementService, requestBody.SetManagedBy); err != nil {
			return nil, fmt.Errorf("failed to set managedBy: %v", err)
		}
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.ScopeIds, requestBody.SetScopeIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	if data.Autopatch != nil {
		autopatch := models.NewCloudPcProvisioningPolicyAutopatch()
		convert.FrameworkToGraphString(data.Autopatch.AutopatchGroupId, autopatch.SetAutopatchGroupId)
		requestBody.SetAutopatch(autopatch)
	}

	if data.AutopilotConfiguration != nil {
		autopilotConfig := models.NewCloudPcAutopilotConfiguration()
		convert.FrameworkToGraphString(data.AutopilotConfiguration.DevicePreparationProfileId, autopilotConfig.SetDevicePreparationProfileId)
		convert.FrameworkToGraphInt32(data.AutopilotConfiguration.ApplicationTimeoutInMinutes, autopilotConfig.SetApplicationTimeoutInMinutes)
		convert.FrameworkToGraphBool(data.AutopilotConfiguration.OnFailureDeviceAccessDenied, autopilotConfig.SetOnFailureDeviceAccessDenied)
		requestBody.SetAutopilotConfiguration(autopilotConfig)
	}

	if err := convert.FrameworkToGraphEnum(data.ImageType,
		models.ParseCloudPcProvisioningPolicyImageType,
		requestBody.SetImageType); err != nil {
		return nil, fmt.Errorf("failed to set image type: %v", err)
	}

	if err := convert.FrameworkToGraphEnum(data.ProvisioningType,
		models.ParseCloudPcProvisioningType,
		requestBody.SetProvisioningType); err != nil {
		return nil, fmt.Errorf("failed to set provisioning type: %v", err)
	}

	if data.MicrosoftManagedDesktop != nil &&
		!data.MicrosoftManagedDesktop.ManagedType.IsUnknown() &&
		!data.MicrosoftManagedDesktop.Profile.IsUnknown() {
		mmd := models.NewMicrosoftManagedDesktop()

		if err := convert.FrameworkToGraphEnum(data.MicrosoftManagedDesktop.ManagedType,
			models.ParseMicrosoftManagedDesktopType,
			mmd.SetManagedType); err != nil {
			return nil, fmt.Errorf("failed to set Microsoft Managed Desktop type: %v", err)
		}

		convert.FrameworkToGraphString(data.MicrosoftManagedDesktop.Profile, mmd.SetProfile)
		requestBody.SetMicrosoftManagedDesktop(mmd)
	}

	if len(data.DomainJoinConfigurations) > 0 {
		var domainJoinConfigs []models.CloudPcDomainJoinConfigurationable
		for _, config := range data.DomainJoinConfigurations {
			domainJoinConfig := models.NewCloudPcDomainJoinConfiguration()

			if err := convert.FrameworkToGraphEnum(config.DomainJoinType,
				models.ParseCloudPcDomainJoinType,
				domainJoinConfig.SetDomainJoinType); err != nil {
				return nil, fmt.Errorf("failed to set domain join type: %v", err)
			}

			convert.FrameworkToGraphString(config.OnPremisesConnectionId, domainJoinConfig.SetOnPremisesConnectionId)
			convert.FrameworkToGraphString(config.RegionName, domainJoinConfig.SetRegionName)

			if !config.RegionGroup.IsNull() && !config.RegionGroup.IsUnknown() && config.RegionGroup.ValueString() != "" {
				if val, err := models.ParseCloudPcRegionGroup(config.RegionGroup.ValueString()); err == nil && val != nil {
					domainJoinConfig.SetRegionGroup(val.(*models.CloudPcRegionGroup))
				}
			}

			domainJoinConfigs = append(domainJoinConfigs, domainJoinConfig)
		}
		requestBody.SetDomainJoinConfigurations(domainJoinConfigs)
	}

	if data.WindowsSetting != nil {
		windowsSetting := models.NewCloudPcWindowsSetting()
		convert.FrameworkToGraphString(data.WindowsSetting.Locale, windowsSetting.SetLocale)
		requestBody.SetWindowsSetting(windowsSetting)
	}

	requestBody.SetAdditionalData(map[string]any{
		"userExperienceType": "cloudPc",
	})

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
