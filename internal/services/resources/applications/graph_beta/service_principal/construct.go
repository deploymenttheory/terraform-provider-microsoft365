package graphBetaServicePrincipal

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs a ServicePrincipal resource from the Terraform model
func constructResource(ctx context.Context, data *ServicePrincipalResourceModel) (graphmodels.ServicePrincipalable, error) {
	requestBody := graphmodels.NewServicePrincipal()

	// Required field: appId
	appId := data.AppID.ValueString()
	requestBody.SetAppId(&appId)

	// Optional boolean fields using helpers
	convert.FrameworkToGraphBool(data.AccountEnabled, requestBody.SetAccountEnabled)
	convert.FrameworkToGraphBool(data.AppRoleAssignmentRequired, requestBody.SetAppRoleAssignmentRequired)

	// Optional string fields using helpers
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphString(data.LoginURL, requestBody.SetLoginUrl)
	convert.FrameworkToGraphString(data.Notes, requestBody.SetNotes)
	convert.FrameworkToGraphString(data.PreferredSingleSignOnMode, requestBody.SetPreferredSingleSignOnMode)

	// Optional collection fields using helpers
	if err := convert.FrameworkToGraphStringSet(ctx, data.Tags, requestBody.SetTags); err != nil {
		return nil, fmt.Errorf("failed to set tags: %w", err)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.NotificationEmailAddresses, requestBody.SetNotificationEmailAddresses); err != nil {
		return nil, fmt.Errorf("failed to set notification_email_addresses: %w", err)
	}

	return requestBody, nil
}
