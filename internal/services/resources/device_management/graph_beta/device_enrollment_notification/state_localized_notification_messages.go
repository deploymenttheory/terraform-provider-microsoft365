package graphBetaDeviceEnrollmentNotification

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// StateLocalizedNotificationMessages maps localized notification messages from notification templates to Terraform state
func StateLocalizedNotificationMessages(ctx context.Context, data *DeviceEnrollmentNotificationResourceModel, templates []graphmodels.NotificationMessageTemplateable, templateTypes []string) {
	if len(templates) == 0 {
		tflog.Debug(ctx, "No templates provided for localized message mapping")
		return
	}

	tflog.Debug(ctx, "Starting localized notification messages mapping", map[string]any{
		"templateCount": len(templates),
		"resourceId":    data.ID.ValueString(),
	})

	// Process localized messages from templates
	localizedMessages := []LocalizedNotificationMessageModel{}

	for i, template := range templates {
		// Get localized messages from the template
		messages := template.GetLocalizedNotificationMessages()
		if messages == nil {
			continue
		}

		for _, message := range messages {
			if message == nil {
				continue
			}

			locale := message.GetLocale()
			subject := message.GetSubject()
			messageTemplate := message.GetMessageTemplate()
			isDefault := message.GetIsDefault()

			if locale == nil || subject == nil || messageTemplate == nil {
				continue
			}

			// Normalize locale to lowercase to match schema normalization and api response behaviour
			normalizedLocale := strings.ToLower(*locale)

			localizedMsg := LocalizedNotificationMessageModel{
				Locale:          types.StringValue(normalizedLocale),
				Subject:         types.StringValue(*subject),
				MessageTemplate: types.StringValue(*messageTemplate),
				TemplateType:    types.StringValue(templateTypes[i]),
			}

			if isDefault != nil {
				localizedMsg.IsDefault = types.BoolValue(*isDefault)
			} else {
				localizedMsg.IsDefault = types.BoolNull()
			}

			localizedMessages = append(localizedMessages, localizedMsg)
		}
	}

	// Set the localized messages in the object
	if len(localizedMessages) > 0 {
		attrTypes := map[string]attr.Type{
			"locale":           types.StringType,
			"subject":          types.StringType,
			"message_template": types.StringType,
			"template_type":    types.StringType,
			"is_default":       types.BoolType,
		}

		localizedMessagesValue, diags := types.SetValueFrom(ctx, types.ObjectType{AttrTypes: attrTypes}, localizedMessages)
		if !diags.HasError() {
			data.LocalizedNotificationMessages = localizedMessagesValue
		} else {
			data.LocalizedNotificationMessages = types.SetNull(types.ObjectType{AttrTypes: attrTypes})
		}
	} else {
		attrTypes := map[string]attr.Type{
			"locale":           types.StringType,
			"subject":          types.StringType,
			"message_template": types.StringType,
			"template_type":    types.StringType,
			"is_default":       types.BoolType,
		}
		data.LocalizedNotificationMessages = types.SetNull(types.ObjectType{AttrTypes: attrTypes})
	}

	tflog.Debug(ctx, "Finished mapping localized notification messages to Terraform state", map[string]any{
		"finalMessageCount": len(localizedMessages),
		"resourceId":        data.ID.ValueString(),
	})
}
