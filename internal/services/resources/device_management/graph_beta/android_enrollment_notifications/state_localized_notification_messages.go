package graphBetaAndroidEnrollmentNotifications

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// StateLocalizedNotificationMessages maps localized notification messages from notification templates to Terraform state
func StateLocalizedNotificationMessages(ctx context.Context, data *AndroidEnrollmentNotificationsResourceModel, templates []graphmodels.NotificationMessageTemplateable, templateTypes []string) {
	if len(templates) == 0 {
		tflog.Debug(ctx, "No templates provided for localized message mapping")
		data.LocalizedNotificationMessages = types.SetNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"locale":           types.StringType,
				"subject":          types.StringType,
				"message_template": types.StringType,
				"template_type":    types.StringType,
				"is_default":       types.BoolType,
			},
		})
		return
	}

	tflog.Debug(ctx, "Starting localized notification messages mapping", map[string]interface{}{
		"templateCount": len(templates),
		"resourceId":    data.ID.ValueString(),
	})

	// Process localized messages from templates
	localizedMessages := []LocalizedNotificationMessageModel{}

	for i, template := range templates {
		// Get localized messages from the template
		messages := template.GetLocalizedNotificationMessages()
		if messages == nil {
			tflog.Debug(ctx, "No localized messages found for template", map[string]interface{}{
				"templateIndex": i,
				"templateType":  templateTypes[i],
				"resourceId":    data.ID.ValueString(),
			})
			continue
		}

		tflog.Debug(ctx, "Processing localized messages for template", map[string]interface{}{
			"templateIndex": i,
			"templateType":  templateTypes[i],
			"messageCount":  len(messages),
			"resourceId":    data.ID.ValueString(),
		})

		for _, message := range messages {
			if message == nil {
				continue
			}

			locale := message.GetLocale()
			subject := message.GetSubject()
			messageTemplate := message.GetMessageTemplate()
			isDefault := message.GetIsDefault()

			if locale == nil || subject == nil || messageTemplate == nil {
				tflog.Debug(ctx, "Skipping incomplete localized message", map[string]interface{}{
					"templateIndex": i,
					"templateType":  templateTypes[i],
					"hasLocale":     locale != nil,
					"hasSubject":    subject != nil,
					"hasTemplate":   messageTemplate != nil,
					"resourceId":    data.ID.ValueString(),
				})
				continue
			}

			localizedMsg := LocalizedNotificationMessageModel{
				Locale:          types.StringValue(*locale),
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
			tflog.Debug(ctx, "Added localized message to collection", map[string]interface{}{
				"locale":          *locale,
				"subject":         *subject,
				"messageTemplate": *messageTemplate,
				"templateType":    templateTypes[i],
				"isDefault":       isDefault,
				"resourceId":      data.ID.ValueString(),
			})
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
			tflog.Debug(ctx, "Successfully set localized notification messages in state", map[string]interface{}{
				"messageCount":      len(localizedMessages),
				"resourceId":        data.ID.ValueString(),
				"finalStateValues":  localizedMessages,
			})
		} else {
			tflog.Warn(ctx, "Failed to set localized notification messages in state", map[string]interface{}{
				"errors":     diags.Errors(),
				"resourceId": data.ID.ValueString(),
			})
			data.LocalizedNotificationMessages = types.SetNull(types.ObjectType{AttrTypes: attrTypes})
		}
	} else {
		tflog.Debug(ctx, "No valid localized messages found, setting to null", map[string]interface{}{
			"resourceId": data.ID.ValueString(),
		})
		attrTypes := map[string]attr.Type{
			"locale":           types.StringType,
			"subject":          types.StringType,
			"message_template": types.StringType,
			"template_type":    types.StringType,
			"is_default":       types.BoolType,
		}
		data.LocalizedNotificationMessages = types.SetNull(types.ObjectType{AttrTypes: attrTypes})
	}

	tflog.Debug(ctx, "Finished mapping localized notification messages to Terraform state", map[string]interface{}{
		"finalMessageCount": len(localizedMessages),
		"resourceId":        data.ID.ValueString(),
	})
}
