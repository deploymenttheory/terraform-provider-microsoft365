package graphBetaIdentityAndAccessLifecycleWorkflowSettings

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	identitygovernance "github.com/microsoftgraph/msgraph-beta-sdk-go/models/identitygovernance"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

var errReadEmailSettings = errors.New("failed to read email_settings")

// constructResource converts the Terraform resource model to a Kiota SDK model for the PATCH request.
// The settings need no lookups to build the request, so the client is unused; it is kept so the
// signature matches the other resources in this domain.
func constructResource(
	ctx context.Context,
	_ *msgraphbetasdk.GraphServiceClient,
	data *LifecycleWorkflowSettingsResourceModel,
) (identitygovernance.LifecycleManagementSettingsable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody, err := constructRequestBody(ctx, data)
	if err != nil {
		return nil, err
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructRequestBody builds the PATCH body from the plan.
//
// PATCH /identityGovernance/lifecycleWorkflows/settings merges the request into the current settings,
// including the properties inside emailSettings, so only known (configured or carried over from state)
// values are sent. Graph rejects null for emailSettings and emailSettings.senderDomain, so a value that
// is unknown or null is omitted instead of being sent as null.
func constructRequestBody(ctx context.Context, data *LifecycleWorkflowSettingsResourceModel) (identitygovernance.LifecycleManagementSettingsable, error) {
	requestBody := identitygovernance.NewLifecycleManagementSettings()

	convert.FrameworkToGraphInt32(data.WorkflowScheduleIntervalInHours, requestBody.SetWorkflowScheduleIntervalInHours)

	if !data.EmailSettings.IsNull() && !data.EmailSettings.IsUnknown() {
		var emailSettings EmailSettingsModel
		if diags := data.EmailSettings.As(ctx, &emailSettings, basetypes.ObjectAsOptions{UnhandledUnknownAsEmpty: true}); diags.HasError() {
			return nil, fmt.Errorf("%w: %v", errReadEmailSettings, diags.Errors())
		}

		settings := graphmodels.NewEmailSettings()
		convert.FrameworkToGraphString(emailSettings.SenderDomain, settings.SetSenderDomain)
		convert.FrameworkToGraphBool(emailSettings.UseCompanyBranding, settings.SetUseCompanyBranding)

		if settings.GetSenderDomain() != nil || settings.GetUseCompanyBranding() != nil {
			requestBody.SetEmailSettings(settings)
		}
	}

	return requestBody, nil
}
