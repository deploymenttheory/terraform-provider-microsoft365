package graphBetaServicePrincipal

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs a ServicePrincipal resource from the Terraform model.
// isUpdate distinguishes PATCH construction from POST construction: on update, properties
// removed from the configuration must be cleared remotely with an explicit JSON null,
// which the generated setters cannot produce (they omit nil fields).
func constructResource(ctx context.Context, data *ServicePrincipalResourceModel, isUpdate bool) (graphmodels.ServicePrincipalable, error) {
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

	if err := convert.FrameworkToGraphStringSet(ctx, data.AlternativeNames, requestBody.SetAlternativeNames); err != nil {
		return nil, fmt.Errorf("failed to set alternative_names: %w", err)
	}

	// Properties removed from the configuration must be cleared remotely with an explicit
	// JSON null, which the generated setters cannot produce (they omit nil fields).
	// SetAdditionalData replaces the whole map, so all explicit nulls accumulate here and
	// are applied once at the end.
	explicitNulls := map[string]any{}

	// Optional nested object: samlSingleSignOnSettings
	if !data.SamlSingleSignOnSettings.IsNull() && !data.SamlSingleSignOnSettings.IsUnknown() {
		var samlSettings SamlSingleSignOnSettingsModel
		diags := data.SamlSingleSignOnSettings.As(ctx, &samlSettings, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			return nil, fmt.Errorf("failed to parse saml_single_sign_on_settings: %v", diags.Errors())
		}

		graphSamlSettings := graphmodels.NewSamlSingleSignOnSettings()
		convert.FrameworkToGraphString(samlSettings.RelayState, graphSamlSettings.SetRelayState)
		requestBody.SetSamlSingleSignOnSettings(graphSamlSettings)
	} else {
		explicitNulls["samlSingleSignOnSettings"] = nil
	}

	// Optional UUID field: tokenEncryptionKeyId. A stale key reference must be clearable,
	// so removal from the configuration is also sent as an explicit null.
	if !data.TokenEncryptionKeyID.IsNull() && !data.TokenEncryptionKeyID.IsUnknown() {
		if err := convert.FrameworkToGraphUUID(data.TokenEncryptionKeyID, requestBody.SetTokenEncryptionKeyId); err != nil {
			return nil, fmt.Errorf("failed to set token_encryption_key_id: %w", err)
		}
	} else {
		explicitNulls["tokenEncryptionKeyId"] = nil
	}

	if isUpdate && len(explicitNulls) > 0 {
		requestBody.SetAdditionalData(explicitNulls)
	}

	return requestBody, nil
}
