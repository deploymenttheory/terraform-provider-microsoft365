package graphBetaWindowsDeviceComplianceNotifications

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructLocalizedMessage creates a single localized notification message
// If isUpdate is true, excludes locale from the request body (required for PATCH operations)
func constructLocalizedMessage(ctx context.Context, msg *LocalizedNotificationMessageModel, isUpdate bool) (*graphmodels.LocalizedNotificationMessage, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing localized message for locale: %s (isUpdate: %t)", msg.Locale.ValueString(), isUpdate))

	requestBody := graphmodels.NewLocalizedNotificationMessage()

	// Only include locale for create operations, not for updates
	if !isUpdate {
		convert.FrameworkToGraphString(msg.Locale, requestBody.SetLocale)
	}

	convert.FrameworkToGraphString(msg.Subject, requestBody.SetSubject)
	convert.FrameworkToGraphString(msg.MessageTemplate, requestBody.SetMessageTemplate)
	convert.FrameworkToGraphBool(msg.IsDefault, requestBody.SetIsDefault)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON for localized message %s (isUpdate: %t)", msg.Locale.ValueString(), isUpdate), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log localized message", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing localized message for locale: %s", msg.Locale.ValueString()))

	return requestBody, nil
}
