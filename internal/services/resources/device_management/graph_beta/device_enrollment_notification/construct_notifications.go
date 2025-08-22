package graphBetaDeviceEnrollmentNotification

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructLocalizedNotificationMessage maps the Terraform localized notification message data to the SDK model.
func constructLocalizedNotificationMessage(ctx context.Context, message LocalizedNotificationMessageModel, isForUpdate ...bool) (graphmodels.LocalizedNotificationMessageable, error) {
	tflog.Debug(ctx, "Constructing LocalizedNotificationMessage from Terraform state")

	requestBody := graphmodels.NewLocalizedNotificationMessage()

	// Always include subject and messageTemplate
	convert.FrameworkToGraphString(message.Subject, requestBody.SetSubject)
	convert.FrameworkToGraphString(message.MessageTemplate, requestBody.SetMessageTemplate)

	// Only include locale and isDefault for create operations (not updates)
	if len(isForUpdate) == 0 || !isForUpdate[0] {
		convert.FrameworkToGraphString(message.Locale, requestBody.SetLocale)
		convert.FrameworkToGraphBool(message.IsDefault, requestBody.SetIsDefault)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
