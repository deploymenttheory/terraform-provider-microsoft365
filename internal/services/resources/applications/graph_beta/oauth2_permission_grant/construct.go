package graphBetaOAuth2PermissionGrant

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"
	"time"

	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

var errScopes = errors.New("failed to extract delegated permission scopes")

// constructResource sends only the mutable scope field for PATCH. Beta requires
// non-null legacy timestamps on POST while explicitly ignoring their values.
func constructResource(ctx context.Context, data *OAuth2PermissionGrantResourceModel, create bool) (graphmodels.OAuth2PermissionGrantable, error) {
	var scopes []string
	if diags := data.Scopes.ElementsAs(ctx, &scopes, false); diags.HasError() {
		return nil, fmt.Errorf("%w: %v", errScopes, diags.Errors())
	}
	slices.Sort(scopes)
	body := graphmodels.NewOAuth2PermissionGrant()
	body.SetScope(new(strings.Join(scopes, " ")))
	if create {
		convert.FrameworkToGraphString(data.ClientID, body.SetClientId)
		convert.FrameworkToGraphString(data.ResourceID, body.SetResourceId)
		convert.FrameworkToGraphString(data.ConsentType, body.SetConsentType)
		convert.FrameworkToGraphString(data.PrincipalID, body.SetPrincipalId)
		start := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
		expiry := time.Date(9999, time.December, 31, 23, 59, 59, 0, time.UTC)
		body.SetStartTime(&start)
		body.SetExpiryTime(&expiry)
	}
	return body, nil
}
