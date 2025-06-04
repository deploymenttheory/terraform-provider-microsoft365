package graphCloudPcProvisioningPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

func constructResource(ctx context.Context, data *CloudPcProvisioningPolicyResourceModel) (*models.CloudPcProvisioningPolicy, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := models.NewCloudPcProvisioningPolicy()

	constructors.SetStringProperty(data.DisplayName, requestBody.SetDisplayName)
	constructors.SetStringProperty(data.Description, requestBody.SetDescription)
	constructors.SetStringProperty(data.CloudPcNamingTemplate, requestBody.SetCloudPcNamingTemplate)
	constructors.SetStringProperty(data.ImageId, requestBody.SetImageId)
	constructors.SetBoolProperty(data.EnableSingleSignOn, requestBody.SetEnableSingleSignOn)
	constructors.SetBoolProperty(data.LocalAdminEnabled, requestBody.SetLocalAdminEnabled)

	if err := constructors.SetEnumProperty(data.ImageType,
		models.ParseCloudPcProvisioningPolicyImageType,
		requestBody.SetImageType); err != nil {
		return nil, fmt.Errorf("failed to set image type: %v", err)
	}

	if err := constructors.SetEnumProperty(data.ProvisioningType,
		models.ParseCloudPcProvisioningType,
		requestBody.SetProvisioningType); err != nil {
		return nil, fmt.Errorf("failed to set provisioning type: %v", err)
	}

	if data.MicrosoftManagedDesktop != nil {
		mmd := models.NewMicrosoftManagedDesktop()

		if err := constructors.SetEnumProperty(data.MicrosoftManagedDesktop.ManagedType,
			models.ParseMicrosoftManagedDesktopType,
			mmd.SetManagedType); err != nil {
			return nil, fmt.Errorf("failed to set Microsoft Managed Desktop type: %v", err)
		}

		constructors.SetStringProperty(data.MicrosoftManagedDesktop.Profile, mmd.SetProfile)
		requestBody.SetMicrosoftManagedDesktop(mmd)
	}

	if len(data.DomainJoinConfigurations) > 0 {
		var domainJoinConfigs []models.CloudPcDomainJoinConfigurationable
		for _, config := range data.DomainJoinConfigurations {
			domainJoinConfig := models.NewCloudPcDomainJoinConfiguration()

			if err := constructors.SetEnumProperty(config.DomainJoinType,
				models.ParseCloudPcDomainJoinType,
				domainJoinConfig.SetDomainJoinType); err != nil {
				return nil, fmt.Errorf("failed to set domain join type: %v", err)
			}

			constructors.SetStringProperty(config.OnPremisesConnectionId, domainJoinConfig.SetOnPremisesConnectionId)
			constructors.SetStringProperty(config.RegionName, domainJoinConfig.SetRegionName)

			domainJoinConfigs = append(domainJoinConfigs, domainJoinConfig)
		}
		requestBody.SetDomainJoinConfigurations(domainJoinConfigs)
	}

	// Handle Windows Settings
	if data.WindowsSetting != nil {
		windowsSetting := models.NewCloudPcWindowsSetting()
		constructors.SetStringProperty(data.WindowsSetting.Locale, windowsSetting.SetLocale)
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
