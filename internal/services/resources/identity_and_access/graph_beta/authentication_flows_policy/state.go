package graphBetaAuthenticationFlowsPolicy

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

var errInvalidPolicyResponse = errors.New(
	"response must contain the authenticationFlowsPolicy singleton and selfServiceSignUp.isEnabled",
)

// mapRemoteState validates and maps the singleton object returned by the beta GET endpoint.
func mapRemoteState(
	ctx context.Context,
	data *AuthenticationFlowsPolicyResourceModel,
	remote graphmodels.AuthenticationFlowsPolicyable,
) error {
	if remote == nil || remote.GetId() == nil || *remote.GetId() != singletonID ||
		remote.GetSelfServiceSignUp() == nil ||
		remote.GetSelfServiceSignUp().GetIsEnabled() == nil {
		return errInvalidPolicyResponse
	}
	signUp, diags := types.ObjectValueFrom(
		ctx,
		selfServiceSignUpAttributeTypes(),
		SelfServiceSignUpModel{
			IsEnabled: convert.GraphToFrameworkBool(remote.GetSelfServiceSignUp().GetIsEnabled()),
		},
	)
	if diags.HasError() {
		return fmt.Errorf("%w: %v", errInvalidPolicyResponse, diags)
	}
	data.ID = types.StringValue(singletonID)
	data.DisplayName = convert.GraphToFrameworkString(remote.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remote.GetDescription())
	data.SelfServiceSignUp = signUp
	return nil
}
