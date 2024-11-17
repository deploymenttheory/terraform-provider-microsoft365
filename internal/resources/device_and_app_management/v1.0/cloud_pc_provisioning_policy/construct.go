package graphCloudPcProvisioningPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/construct"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model
func constructResource(ctx context.Context, typeName string, data *CloudPcProvisioningPolicyResourceModel) (*models.CloudPcProvisioningPolicy, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", typeName))

	requestBody := models.NewCloudPcProvisioningPolicy()

	displayName := data.DisplayName.ValueString()
	requestBody.SetDisplayName(&displayName)

	if !data.Description.IsNull() {
		description := data.Description.ValueString()
		requestBody.SetDescription(&description)
	}

	if !data.CloudPcNamingTemplate.IsNull() {
		namingTemplate := data.CloudPcNamingTemplate.ValueString()
		requestBody.SetCloudPcNamingTemplate(&namingTemplate)
	}

	if !data.EnableSingleSignOn.IsNull() {
		enableSSO := data.EnableSingleSignOn.ValueBool()
		requestBody.SetEnableSingleSignOn(&enableSSO)
	}

	if !data.ImageId.IsNull() {
		imageId := data.ImageId.ValueString()
		requestBody.SetImageId(&imageId)
	}

	if !data.ImageType.IsNull() {
		imageTypeStr := data.ImageType.ValueString()
		imageType, err := models.ParseCloudPcProvisioningPolicyImageType(imageTypeStr)
		if err != nil {
			return nil, fmt.Errorf("invalid image type: %s", err)
		}
		requestBody.SetImageType(imageType.(*models.CloudPcProvisioningPolicyImageType))
	}

	if !data.LocalAdminEnabled.IsNull() {
		localAdminEnabled := data.LocalAdminEnabled.ValueBool()
		requestBody.SetLocalAdminEnabled(&localAdminEnabled)
	}

	if !data.ProvisioningType.IsNull() {
		provisioningTypeStr := data.ProvisioningType.ValueString()
		provisioningType, err := models.ParseCloudPcProvisioningType(provisioningTypeStr)
		if err != nil {
			return nil, fmt.Errorf("invalid provisioning type: %s", err)
		}
		requestBody.SetProvisioningType(provisioningType.(*models.CloudPcProvisioningType))
	}

	if data.MicrosoftManagedDesktop != nil {
		mmd := models.NewMicrosoftManagedDesktop()

		if managedType := data.MicrosoftManagedDesktop.ManagedType; !managedType.IsNull() {
			managedTypeStr := managedType.ValueString()
			managedTypeAny, err := models.ParseMicrosoftManagedDesktopType(managedTypeStr)
			if err != nil {
				return nil, fmt.Errorf("invalid Microsoft Managed Desktop type: %s", err)
			}
			if managedTypeAny != nil {
				managedTypeEnum, ok := managedTypeAny.(*models.MicrosoftManagedDesktopType)
				if !ok {
					return nil, fmt.Errorf("unexpected type for Microsoft Managed Desktop type")
				}
				mmd.SetManagedType(managedTypeEnum)
			}
		}

		if profile := data.MicrosoftManagedDesktop.Profile; !profile.IsNull() {
			profileValue := profile.ValueString()
			mmd.SetProfile(&profileValue)
		}

		requestBody.SetMicrosoftManagedDesktop(mmd)
	}

	if len(data.DomainJoinConfigurations) > 0 {
		var domainJoinConfigs []models.CloudPcDomainJoinConfigurationable
		for _, config := range data.DomainJoinConfigurations {
			domainJoinConfig := models.NewCloudPcDomainJoinConfiguration()

			if !config.DomainJoinType.IsNull() {
				domainJoinTypeStr := config.DomainJoinType.ValueString()
				domainJoinTypeAny, err := models.ParseCloudPcDomainJoinType(domainJoinTypeStr)
				if err != nil {
					return nil, fmt.Errorf("invalid domain join type: %s", err)
				}
				if domainJoinTypeAny != nil {
					domainJoinTypeEnum, ok := domainJoinTypeAny.(*models.CloudPcDomainJoinType)
					if !ok {
						return nil, fmt.Errorf("unexpected type for domain join type")
					}
					domainJoinConfig.SetDomainJoinType(domainJoinTypeEnum)
				}
			}

			if !config.OnPremisesConnectionId.IsNull() {
				onPremisesConnectionId := config.OnPremisesConnectionId.ValueString()
				domainJoinConfig.SetOnPremisesConnectionId(&onPremisesConnectionId)
			}

			if !config.RegionName.IsNull() {
				regionName := config.RegionName.ValueString()
				domainJoinConfig.SetRegionName(&regionName)
			}

			domainJoinConfigs = append(domainJoinConfigs, domainJoinConfig)
		}
		requestBody.SetDomainJoinConfigurations(domainJoinConfigs)
	}

	if data.WindowsSetting != nil {
		windowsSetting := models.NewCloudPcWindowsSetting()

		if !data.WindowsSetting.Locale.IsNull() {
			locale := data.WindowsSetting.Locale.ValueString()
			windowsSetting.SetLocale(&locale)
		}

		requestBody.SetWindowsSetting(windowsSetting)
	}

	if err := construct.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", typeName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", typeName))

	return requestBody, nil
}
