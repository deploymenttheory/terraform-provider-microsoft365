package graphBetaDeviceEnrollmentNotificationConfiguration

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/state"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func mapRemoteStateToTerraform(ctx context.Context, data *DeviceEnrollmentNotificationConfigurationResourceModel, remoteResource models.DeviceEnrollmentConfigurationable) {
	if remoteResource == nil {
		return
	}

	if notificationConfig, ok := remoteResource.(models.DeviceEnrollmentNotificationConfigurationable); ok {
		data.ID = state.StringPointerValue(notificationConfig.GetId())
		data.DisplayName = state.StringPointerValue(notificationConfig.GetDisplayName())
		data.Description = state.StringPointerValue(notificationConfig.GetDescription())
		data.Priority = state.Int32PointerValue(notificationConfig.GetPriority())
		data.CreatedDateTime = state.TimeToString(notificationConfig.GetCreatedDateTime())
		data.LastModifiedDateTime = state.TimeToString(notificationConfig.GetLastModifiedDateTime())
		data.Version = state.Int32PointerValue(notificationConfig.GetVersion())

		if configType := notificationConfig.GetDeviceEnrollmentConfigurationType(); configType != nil {
			data.DeviceEnrollmentConfigurationType = types.StringValue(configType.String())
		}

		if platformType := notificationConfig.GetPlatformType(); platformType != nil {
			data.PlatformType = types.StringValue(platformType.String())
		}

		if data.TemplateTypes.IsNull() || data.TemplateTypes.IsUnknown() {
			if templateType := notificationConfig.GetTemplateType(); templateType != nil {
				data.TemplateTypes = state.StringSliceToSet(ctx, []string{templateType.String()})
			}
		}

		if brandingOptions := notificationConfig.GetBrandingOptions(); brandingOptions != nil {
			data.BrandingOptions = types.StringValue(brandingOptions.String())
		}

		if data.NotificationMessageTemplateId.IsNull() || data.NotificationMessageTemplateId.IsUnknown() {
			if templateId := notificationConfig.GetNotificationMessageTemplateId(); templateId != nil {
				data.NotificationMessageTemplateId = types.StringValue(templateId.String())
			}
		}

		if data.NotificationTemplates.IsNull() || data.NotificationTemplates.IsUnknown() {
			if notificationTemplates := notificationConfig.GetNotificationTemplates(); notificationTemplates != nil {
				data.NotificationTemplates = state.StringSliceToSet(ctx, notificationTemplates)
			}
		}

		if roleScopeTagIds := notificationConfig.GetRoleScopeTagIds(); roleScopeTagIds != nil {
			data.RoleScopeTagIds = state.StringSliceToSet(ctx, roleScopeTagIds)
		}
	}
}
