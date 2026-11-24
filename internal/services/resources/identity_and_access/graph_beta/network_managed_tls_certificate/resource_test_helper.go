package graphBetaNetworkManagedTLSCertificate

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// NetworkManagedTLSCertificateTestResource implements the acceptance TestResource interface.
type NetworkManagedTLSCertificateTestResource struct{}

// Exists checks whether the Microsoft-managed TLS certificate authority exists in Microsoft Graph.
func (NetworkManagedTLSCertificateTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return nil, err
	}

	getConfig := customrequests.GetRequestConfig{
		APIVersion:        customrequests.GraphAPIBeta,
		Endpoint:          "networkaccess/tls/managedCertificateAuthorityCertificates",
		ResourceIDPattern: customrequests.ResourceIDPatternSegment,
		ResourceID:        state.ID,
		Headers:           managedTLSCertificateHeaders,
	}
	_, err = customrequests.GetRequestByResourceId(ctx, graphClient.GetAdapter(), getConfig)
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

	exists := true
	return &exists, nil
}
