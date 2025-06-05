package graphBetaDeviceEnrollmentNotificationConfiguration

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	msgraphbetamodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform resource model to the Graph API request model
func constructResource(ctx context.Context, data *DeviceEnrollmentNotificationConfigurationResourceModel) (models.DeviceEnrollmentNotificationConfigurationable, error) {
	requestBody := models.NewDeviceEnrollmentNotificationConfiguration()

	constructors.SetStringProperty(data.DisplayName, requestBody.SetDisplayName)
	constructors.SetStringProperty(data.Description, requestBody.SetDescription)
	constructors.SetInt32Property(data.Priority, requestBody.SetPriority)

	// Platform type is always Windows for enrollment notifications
	platformType := models.WINDOWS_ENROLLMENTRESTRICTIONPLATFORMTYPE
	requestBody.SetPlatformType(&platformType)

	brandingOptions := models.EnrollmentNotificationBrandingOptions(models.NONE_ENROLLMENTNOTIFICATIONBRANDINGOPTIONS)
	requestBody.SetBrandingOptions(&brandingOptions)

	var notificationTemplates []string
	if !data.TemplateTypes.IsNull() && !data.TemplateTypes.IsUnknown() {
		for _, element := range data.TemplateTypes.Elements() {
			if stringVal, ok := element.(types.String); ok && !stringVal.IsNull() {
				templateType := stringVal.ValueString()

				notificationTemplates = append(notificationTemplates, fmt.Sprintf("%s_00000000-0000-0000-0000-000000000000", templateType))
			}
		}
	}
	if len(notificationTemplates) > 0 {
		requestBody.SetNotificationTemplates(notificationTemplates)
	}

	if !data.NotificationMessageTemplateId.IsNull() && !data.NotificationMessageTemplateId.IsUnknown() {
		if templateId, err := uuid.Parse(data.NotificationMessageTemplateId.ValueString()); err == nil {
			requestBody.SetNotificationMessageTemplateId(&templateId)
		}
	}

	if err := constructors.SetStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructLocalizedMessage constructs a localized message request body
func constructLocalizedMessage(ctx context.Context, message *LocalizedNotificationMessageModel) *msgraphbetamodels.LocalizedNotificationMessage {
	if message == nil {
		return nil
	}

	requestBody := msgraphbetamodels.NewLocalizedNotificationMessage()

	// Set default locale if not specified
	locale := "en-us"
	if !message.Locale.IsNull() && !message.Locale.IsUnknown() {
		locale = message.Locale.ValueString()
	}
	requestBody.SetLocale(&locale)

	if !message.Subject.IsNull() && !message.Subject.IsUnknown() {
		subject := message.Subject.ValueString()
		requestBody.SetSubject(&subject)
	}

	if !message.MessageTemplate.IsNull() && !message.MessageTemplate.IsUnknown() {
		messageTemplate := message.MessageTemplate.ValueString()
		requestBody.SetMessageTemplate(&messageTemplate)
	}

	// Set default to true if not specified
	isDefault := true
	if !message.IsDefault.IsNull() && !message.IsDefault.IsUnknown() {
		isDefault = message.IsDefault.ValueBool()
	}
	requestBody.SetIsDefault(&isDefault)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return requestBody
}
