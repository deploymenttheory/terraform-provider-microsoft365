package graphCloudPcProvisioningPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

func constructResource(ctx context.Context, data *CloudPcProvisioningPolicyResourceModel) (*models.CloudPcProvisioningPolicy, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := models.NewCloudPcProvisioningPolicy()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphString(data.CloudPcNamingTemplate, requestBody.SetCloudPcNamingTemplate)
	convert.FrameworkToGraphString(data.ImageId, requestBody.SetImageId)
	convert.FrameworkToGraphBool(data.EnableSingleSignOn, requestBody.SetEnableSingleSignOn)
	convert.FrameworkToGraphBool(data.LocalAdminEnabled, requestBody.SetLocalAdminEnabled)

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

	if data.MicrosoftManagedDesktop != nil {
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

			domainJoinConfigs = append(domainJoinConfigs, domainJoinConfig)
		}
		requestBody.SetDomainJoinConfigurations(domainJoinConfigs)
	}

	// Handle Windows Settings
	if data.WindowsSetting != nil {
		windowsSetting := models.NewCloudPcWindowsSetting()
		convert.FrameworkToGraphString(data.WindowsSetting.Locale, windowsSetting.SetLocale)
		requestBody.SetWindowsSetting(windowsSetting)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
