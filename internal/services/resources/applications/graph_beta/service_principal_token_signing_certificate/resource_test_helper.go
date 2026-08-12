package graphBetaServicePrincipalTokenSigningCertificate

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// ServicePrincipalTokenSigningCertificateTestResource implements the types.TestResource interface
type ServicePrincipalTokenSigningCertificateTestResource struct{}

// Exists checks whether the signing key credential exists on the service principal.
func (r ServicePrincipalTokenSigningCertificateTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
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

	keyId := state.Attributes["key_id"]
	for _, credential := range servicePrincipal.GetKeyCredentials() {
		if credentialKeyId := credential.GetKeyId(); credentialKeyId != nil && credentialKeyId.String() == keyId {
			exists := true
			return &exists, nil
		}
	}

	exists := false
	return &exists, nil
}
