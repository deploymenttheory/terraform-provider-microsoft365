package graphBetaNotificationMessageTemplates

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the base properties of a NotificationMessageTemplateResourceModel to a Terraform state.
func MapRemoteStateToTerraform(ctx context.Context, data *NotificationMessageTemplateResourceModel, remoteResource graphmodels.NotificationMessageTemplateable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.DefaultLocale = convert.GraphToFrameworkString(remoteResource.GetDefaultLocale())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())

	// Handle localized notification messages
	if remoteResource.GetLocalizedNotificationMessages() != nil {
		localizedMessages := remoteResource.GetLocalizedNotificationMessages()
		var terraformMessages []LocalizedNotificationMessageModel
		
		for _, msg := range localizedMessages {
			terraformMsg := LocalizedNotificationMessageModel{
				ID:                   convert.GraphToFrameworkString(msg.GetId()),
				Locale:               convert.GraphToFrameworkString(msg.GetLocale()),
				Subject:              convert.GraphToFrameworkString(msg.GetSubject()),
				MessageTemplate:      convert.GraphToFrameworkString(msg.GetMessageTemplate()),
				IsDefault:            convert.GraphToFrameworkBool(msg.GetIsDefault()),
				LastModifiedDateTime: convert.GraphToFrameworkTime(msg.GetLastModifiedDateTime()),
			}
			terraformMessages = append(terraformMessages, terraformMsg)
		}
		
		// Convert to framework set using the object type
		if len(terraformMessages) > 0 {
			data.LocalizedNotificationMessages, _ = types.SetValueFrom(ctx, types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"id":                        types.StringType,
					"locale":                    types.StringType,
					"subject":                   types.StringType,
					"message_template":          types.StringType,
					"is_default":               types.BoolType,
					"last_modified_date_time":  types.StringType,
				},
			}, terraformMessages)
		} else {
			data.LocalizedNotificationMessages = types.SetNull(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"id":                        types.StringType,
					"locale":                    types.StringType,
					"subject":                   types.StringType,
					"message_template":          types.StringType,
					"is_default":               types.BoolType,
					"last_modified_date_time":  types.StringType,
				},
			})
		}
	} else {
		data.LocalizedNotificationMessages = types.SetNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":                        types.StringType,
				"locale":                    types.StringType,
				"subject":                   types.StringType,
				"message_template":          types.StringType,
				"is_default":               types.BoolType,
				"last_modified_date_time":  types.StringType,
			},
		})
	}

	// Handle branding options enum conversion
	if remoteResource.GetBrandingOptions() != nil {
		brandingOptions := remoteResource.GetBrandingOptions()
		switch *brandingOptions {
		case graphmodels.NONE_NOTIFICATIONTEMPLATEBRANDINGOPTIONS:
			data.BrandingOptions = convert.GraphToFrameworkString(&[]string{"none"}[0])
		case graphmodels.INCLUDECOMPANYLOGO_NOTIFICATIONTEMPLATEBRANDINGOPTIONS:
			data.BrandingOptions = convert.GraphToFrameworkString(&[]string{"includeCompanyLogo"}[0])
		case graphmodels.INCLUDECOMPANYNAME_NOTIFICATIONTEMPLATEBRANDINGOPTIONS:
			data.BrandingOptions = convert.GraphToFrameworkString(&[]string{"includeCompanyName"}[0])
		case graphmodels.INCLUDECONTACTINFORMATION_NOTIFICATIONTEMPLATEBRANDINGOPTIONS:
			data.BrandingOptions = convert.GraphToFrameworkString(&[]string{"includeContactInformation"}[0])
		case graphmodels.INCLUDECOMPANYPORTALLINK_NOTIFICATIONTEMPLATEBRANDINGOPTIONS:
			data.BrandingOptions = convert.GraphToFrameworkString(&[]string{"includeCompanyPortalLink"}[0])
		case graphmodels.INCLUDEDEVICEDETAILS_NOTIFICATIONTEMPLATEBRANDINGOPTIONS:
			data.BrandingOptions = convert.GraphToFrameworkString(&[]string{"includeDeviceDetails"}[0])
		default:
			data.BrandingOptions = convert.GraphToFrameworkString(&[]string{"none"}[0])
		}
	} else {
		data.BrandingOptions = convert.GraphToFrameworkString(&[]string{"none"}[0])
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}