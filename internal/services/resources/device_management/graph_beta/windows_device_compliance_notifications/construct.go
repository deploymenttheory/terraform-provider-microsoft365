package graphBetaWindowsDeviceComplianceNotifications

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource builds the base resource request body.
func constructResource(ctx context.Context, data *WindowsDeviceComplianceNotificationsResourceModel) (*graphmodels.NotificationMessageTemplate, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing base %s template from model", ResourceName))

	if err := validateRequest(ctx, data); err != nil {
		return nil, fmt.Errorf("validation failed: %s", err.Error())
	}

	requestBody := graphmodels.NewNotificationMessageTemplate()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)

	// Handle branding options as a set with bitwise flags
	if !data.BrandingOptions.IsNull() && !data.BrandingOptions.IsUnknown() {
		var brandingValues []string
		data.BrandingOptions.ElementsAs(ctx, &brandingValues, false)

		var brandingOptions graphmodels.NotificationTemplateBrandingOptions = 0

		for _, value := range brandingValues {
			switch value {
			case "none":
				brandingOptions |= graphmodels.NONE_NOTIFICATIONTEMPLATEBRANDINGOPTIONS
			case "includeCompanyLogo":
				brandingOptions |= graphmodels.INCLUDECOMPANYLOGO_NOTIFICATIONTEMPLATEBRANDINGOPTIONS
			case "includeCompanyName":
				brandingOptions |= graphmodels.INCLUDECOMPANYNAME_NOTIFICATIONTEMPLATEBRANDINGOPTIONS
			case "includeContactInformation":
				brandingOptions |= graphmodels.INCLUDECONTACTINFORMATION_NOTIFICATIONTEMPLATEBRANDINGOPTIONS
			case "includeCompanyPortalLink":
				brandingOptions |= graphmodels.INCLUDECOMPANYPORTALLINK_NOTIFICATIONTEMPLATEBRANDINGOPTIONS
			case "includeDeviceDetails":
				brandingOptions |= graphmodels.INCLUDEDEVICEDETAILS_NOTIFICATIONTEMPLATEBRANDINGOPTIONS
			}
		}

		if brandingOptions == 0 {
			brandingOptions = graphmodels.NONE_NOTIFICATIONTEMPLATEBRANDINGOPTIONS
		}

		requestBody.SetBrandingOptions(&brandingOptions)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tag ids: %s", err)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for base %s template", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing base %s template", ResourceName))

	return requestBody, nil
}
