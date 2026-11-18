package graphBetaOAuth2PermissionGrant

import (
	"context"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

func MapRemoteStateToTerraform(ctx context.Context, data *OAuth2PermissionGrantResourceModel, remote graphmodels.OAuth2PermissionGrantable) {
	data.ID = convert.GraphToFrameworkString(remote.GetId())
	data.ClientID = convert.GraphToFrameworkString(remote.GetClientId())
	data.ResourceID = convert.GraphToFrameworkString(remote.GetResourceId())
	data.ConsentType = convert.GraphToFrameworkString(remote.GetConsentType())
	data.PrincipalID = convert.GraphToFrameworkString(remote.GetPrincipalId())
	scopes := convert.GraphToFrameworkString(remote.GetScope()).ValueString()
	values := strings.Fields(scopes)
	slices.Sort(values)
	data.Scopes = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, slices.Compact(values))
}

func grantConsistencyPredicate(expected *OAuth2PermissionGrantResourceModel) func(context.Context, tfsdk.State) bool {
	matchingReads := 0
	return func(ctx context.Context, state tfsdk.State) bool {
		var actual OAuth2PermissionGrantResourceModel
		if state.Get(ctx, &actual).HasError() {
			matchingReads = 0
			return false
		}
		matches := actual.ID.Equal(expected.ID) && actual.ClientID.Equal(expected.ClientID) &&
			actual.ResourceID.Equal(expected.ResourceID) && actual.ConsentType.Equal(expected.ConsentType) &&
			actual.PrincipalID.Equal(expected.PrincipalID) && actual.Scopes.Equal(expected.Scopes)
		if !matches {
			matchingReads = 0
			return false
		}
		// A single matching replica can be followed by a stale scope response on
		// Terraform's immediate refresh. Require successive matching observations.
		matchingReads++
		return matchingReads >= 5
	}
}
