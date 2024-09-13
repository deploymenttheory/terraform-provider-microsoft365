package graphCloudPcProvisioningPolicy

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model
func constructResource(ctx context.Context, data *CloudPcProvisioningPolicyResourceModel) (*models.CloudPcProvisioningPolicy, error) {
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

	// Debug logging
	debugPrintRequestBody(ctx, requestBody)

	return requestBody, nil
}

func debugPrintRequestBody(ctx context.Context, requestBody *models.CloudPcProvisioningPolicy) {
	requestMap := map[string]interface{}{
		"displayName":              requestBody.GetDisplayName(),
		"description":              requestBody.GetDescription(),
		"cloudPcNamingTemplate":    requestBody.GetCloudPcNamingTemplate(),
		"enableSingleSignOn":       requestBody.GetEnableSingleSignOn(),
		"imageId":                  requestBody.GetImageId(),
		"imageType":                requestBody.GetImageType(),
		"localAdminEnabled":        requestBody.GetLocalAdminEnabled(),
		"provisioningType":         requestBody.GetProvisioningType(),
		"microsoftManagedDesktop":  debugMapMicrosoftManagedDesktop(requestBody.GetMicrosoftManagedDesktop()),
		"domainJoinConfigurations": debugMapDomainJoinConfigurations(requestBody.GetDomainJoinConfigurations()),
		"windowsSetting":           debugMapWindowsSetting(requestBody.GetWindowsSetting()),
	}

	requestBodyJSON, err := json.MarshalIndent(requestMap, "", "  ")
	if err != nil {
		tflog.Error(ctx, "Error marshalling request body to JSON", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	tflog.Debug(ctx, "Constructed Cloud PC Provisioning Policy resource", map[string]interface{}{
		"requestBody": string(requestBodyJSON),
	})
}

func debugMapMicrosoftManagedDesktop(mmd models.MicrosoftManagedDesktopable) map[string]interface{} {
	if mmd == nil {
		return nil
	}
	return map[string]interface{}{
		"managedType": mmd.GetManagedType(),
		"profile":     mmd.GetProfile(),
	}
}

func debugMapDomainJoinConfigurations(configs []models.CloudPcDomainJoinConfigurationable) []map[string]interface{} {
	result := make([]map[string]interface{}, len(configs))
	for i, config := range configs {
		result[i] = map[string]interface{}{
			"domainJoinType":         config.GetDomainJoinType(),
			"onPremisesConnectionId": config.GetOnPremisesConnectionId(),
			"regionName":             config.GetRegionName(),
		}
	}
	return result
}

func debugMapWindowsSetting(setting models.CloudPcWindowsSettingable) map[string]interface{} {
	if setting == nil {
		return nil
	}
	return map[string]interface{}{
		"locale": setting.GetLocale(),
	}
}
