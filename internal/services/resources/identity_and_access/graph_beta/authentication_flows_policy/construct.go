package graphBetaAuthenticationFlowsPolicy

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

var errInvalidSignUp = errors.New("self_service_sign_up.is_enabled must be a known Boolean")

// constructResource sends only the writable self-service sign-up configuration.
func constructResource(
	ctx context.Context,
	data *AuthenticationFlowsPolicyResourceModel,
) (graphmodels.AuthenticationFlowsPolicyable, error) {
	if data.SelfServiceSignUp.IsNull() || data.SelfServiceSignUp.IsUnknown() {
		return nil, errInvalidSignUp
	}
	var config SelfServiceSignUpModel
	if diags := data.SelfServiceSignUp.As(
		ctx,
		&config,
		basetypes.ObjectAsOptions{},
	); diags.HasError() {
		return nil, fmt.Errorf("%w: %v", errInvalidSignUp, diags)
	}
	if config.IsEnabled.IsNull() || config.IsEnabled.IsUnknown() {
		return nil, errInvalidSignUp
	}
	body := graphmodels.NewAuthenticationFlowsPolicy()
	signUp := graphmodels.NewSelfServiceSignUpAuthenticationFlowConfiguration()
	convert.FrameworkToGraphBool(config.IsEnabled, signUp.SetIsEnabled)
	body.SetSelfServiceSignUp(signUp)
	if err := constructors.DebugLogGraphObject(
		ctx,
		fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName),
		body,
	); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
	}
	return body, nil
}
