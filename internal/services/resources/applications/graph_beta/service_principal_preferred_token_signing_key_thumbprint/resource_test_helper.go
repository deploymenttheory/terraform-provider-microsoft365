package graphBetaServicePrincipalPreferredTokenSigningKeyThumbprint

import (
	"context"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// ServicePrincipalPreferredTokenSigningKeyThumbprintTestResource implements the types.TestResource interface
type ServicePrincipalPreferredTokenSigningKeyThumbprintTestResource struct{}

// Exists checks whether the service principal has the expected preferred token signing key
// thumbprint set. A cleared property (or a missing service principal) reports non-existence,
// which is what destroy verification expects after Delete sends the explicit null.
func (r ServicePrincipalPreferredTokenSigningKeyThumbprintTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return nil, err
	}

	servicePrincipal, err := graphClient.
		ServicePrincipals().
		ByServicePrincipalId(state.Attributes["service_principal_id"]).
		Get(ctx, nil)

	if err != nil {
		errorInfo := errors.GraphError(ctx, err)
		if errorInfo.StatusCode == 404 ||
			errorInfo.ErrorCode == "ResourceNotFound" ||
			errorInfo.ErrorCode == "Request_ResourceNotFound" {
			exists := false
			return &exists, nil
		}
		return nil, err
	}

	remote := servicePrincipal.GetPreferredTokenSigningKeyThumbprint()
	exists := remote != nil && strings.EqualFold(*remote, state.Attributes["thumbprint"])
	return &exists, nil
}
