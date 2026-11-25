package graphBetaOAuth2PermissionGrant

import (
	"context"
	"slices"
	"strings"

	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// MapRemoteStateToTerraform maps the remote Microsoft Graph response to Terraform state.
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
