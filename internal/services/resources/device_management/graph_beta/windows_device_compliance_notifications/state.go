package graphBetaWindowsDeviceComplianceNotifications

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the base properties of a WindowsDeviceComplianceNotificationsResourceModel to a Terraform state.
func MapRemoteStateToTerraform(ctx context.Context, data *WindowsDeviceComplianceNotificationsResourceModel, remoteResource graphmodels.NotificationMessageTemplateable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(remoteResource.GetId()),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.DefaultLocale = convert.GraphToFrameworkString(remoteResource.GetDefaultLocale())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())

	// Handle localized notification messages
	if remoteResource.GetLocalizedNotificationMessages() != nil {
		localizedMessages := remoteResource.GetLocalizedNotificationMessages()
		var terraformMessages []LocalizedNotificationMessageModel

		for _, msg := range localizedMessages {
			terraformMsg := LocalizedNotificationMessageModel{
				ID:              convert.GraphToFrameworkString(msg.GetId()),
				Locale:          convert.GraphToFrameworkString(msg.GetLocale()),
				Subject:         convert.GraphToFrameworkString(msg.GetSubject()),
				MessageTemplate: convert.GraphToFrameworkString(msg.GetMessageTemplate()),
				IsDefault:       convert.GraphToFrameworkBool(msg.GetIsDefault()),
			}
			terraformMessages = append(terraformMessages, terraformMsg)
		}

		// Convert to framework set using the object type
		if len(terraformMessages) > 0 {
			data.LocalizedNotificationMessages, _ = types.SetValueFrom(ctx, types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"id":               types.StringType,
					"locale":           types.StringType,
					"subject":          types.StringType,
					"message_template": types.StringType,
					"is_default":       types.BoolType,
				},
			}, terraformMessages)
		} else {
			data.LocalizedNotificationMessages = types.SetNull(types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"id":               types.StringType,
					"locale":           types.StringType,
					"subject":          types.StringType,
					"message_template": types.StringType,
					"is_default":       types.BoolType,
				},
			})
		}
	} else {
		data.LocalizedNotificationMessages = types.SetNull(types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"id":               types.StringType,
				"locale":           types.StringType,
				"subject":          types.StringType,
				"message_template": types.StringType,
				"is_default":       types.BoolType,
			},
		})
	}

	// Handle branding options as a set with bitwise flags
	if remoteResource.GetBrandingOptions() != nil {
		brandingOptions := remoteResource.GetBrandingOptions()
		var brandingValues []string

		// Extract individual flags from the bitwise combination
		if *brandingOptions&graphmodels.NONE_NOTIFICATIONTEMPLATEBRANDINGOPTIONS != 0 {
			brandingValues = append(brandingValues, "none")
		}
		if *brandingOptions&graphmodels.INCLUDECOMPANYLOGO_NOTIFICATIONTEMPLATEBRANDINGOPTIONS != 0 {
			brandingValues = append(brandingValues, "includeCompanyLogo")
		}
		if *brandingOptions&graphmodels.INCLUDECOMPANYNAME_NOTIFICATIONTEMPLATEBRANDINGOPTIONS != 0 {
			brandingValues = append(brandingValues, "includeCompanyName")
		}
		if *brandingOptions&graphmodels.INCLUDECONTACTINFORMATION_NOTIFICATIONTEMPLATEBRANDINGOPTIONS != 0 {
			brandingValues = append(brandingValues, "includeContactInformation")
		}
		if *brandingOptions&graphmodels.INCLUDECOMPANYPORTALLINK_NOTIFICATIONTEMPLATEBRANDINGOPTIONS != 0 {
			brandingValues = append(brandingValues, "includeCompanyPortalLink")
		}
		if *brandingOptions&graphmodels.INCLUDEDEVICEDETAILS_NOTIFICATIONTEMPLATEBRANDINGOPTIONS != 0 {
			brandingValues = append(brandingValues, "includeDeviceDetails")
		}

		// If no flags were extracted, default to none
		if len(brandingValues) == 0 {
			brandingValues = []string{"none"}
		}

		// Convert to Terraform set
		brandingSet, diags := types.SetValueFrom(ctx, types.StringType, brandingValues)
		if diags.HasError() {
			tflog.Warn(ctx, "Failed to convert branding options to set, defaulting to none")
			brandingSet, _ = types.SetValueFrom(ctx, types.StringType, []string{"none"})
		}
		data.BrandingOptions = brandingSet
	} else {
		// Default to none if no branding options are present
		brandingSet, _ := types.SetValueFrom(ctx, types.StringType, []string{"none"})
		data.BrandingOptions = brandingSet
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}
