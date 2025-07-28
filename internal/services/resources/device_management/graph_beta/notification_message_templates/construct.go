package graphBetaNotificationMessageTemplates

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model
func constructResource(ctx context.Context, data *NotificationMessageTemplateResourceModel) (*graphmodels.NotificationMessageTemplate, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewNotificationMessageTemplate()

	// String fields using helper functions
	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphString(data.DefaultLocale, requestBody.SetDefaultLocale)

	// Handle branding options enum conversion
	if !data.BrandingOptions.IsNull() && !data.BrandingOptions.IsUnknown() {
		brandingValue := data.BrandingOptions.ValueString()
		var brandingOptions graphmodels.NotificationTemplateBrandingOptions
		
		switch brandingValue {
		case "none":
			brandingOptions = graphmodels.NONE_NOTIFICATIONTEMPLATEBRANDINGOPTIONS
		case "includeCompanyLogo":
			brandingOptions = graphmodels.INCLUDECOMPANYLOGO_NOTIFICATIONTEMPLATEBRANDINGOPTIONS
		case "includeCompanyName":
			brandingOptions = graphmodels.INCLUDECOMPANYNAME_NOTIFICATIONTEMPLATEBRANDINGOPTIONS
		case "includeContactInformation":
			brandingOptions = graphmodels.INCLUDECONTACTINFORMATION_NOTIFICATIONTEMPLATEBRANDINGOPTIONS
		case "includeCompanyPortalLink":
			brandingOptions = graphmodels.INCLUDECOMPANYPORTALLINK_NOTIFICATIONTEMPLATEBRANDINGOPTIONS
		case "includeDeviceDetails":
			brandingOptions = graphmodels.INCLUDEDEVICEDETAILS_NOTIFICATIONTEMPLATEBRANDINGOPTIONS
		default:
			brandingOptions = graphmodels.NONE_NOTIFICATIONTEMPLATEBRANDINGOPTIONS
		}
		requestBody.SetBrandingOptions(&brandingOptions)
	}

	// Role scope tag IDs - using string slice helper
	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tag ids: %s", err)
	}

	// Handle localized notification messages
	if !data.LocalizedNotificationMessages.IsNull() && !data.LocalizedNotificationMessages.IsUnknown() {
		var localizedMessages []LocalizedNotificationMessageModel
		data.LocalizedNotificationMessages.ElementsAs(ctx, &localizedMessages, false)
		
		var graphLocalizedMessages []graphmodels.LocalizedNotificationMessageable
		for _, msg := range localizedMessages {
			graphMsg := graphmodels.NewLocalizedNotificationMessage()
			
			convert.FrameworkToGraphString(msg.Locale, graphMsg.SetLocale)
			convert.FrameworkToGraphString(msg.Subject, graphMsg.SetSubject)
			convert.FrameworkToGraphString(msg.MessageTemplate, graphMsg.SetMessageTemplate)
			convert.FrameworkToGraphBool(msg.IsDefault, graphMsg.SetIsDefault)
			
			graphLocalizedMessages = append(graphLocalizedMessages, graphMsg)
		}
		requestBody.SetLocalizedNotificationMessages(graphLocalizedMessages)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}