package graphBetaNetworkManagedTLSCertificate

import (
	"context"

	"github.com/hashicorp/terraform-plugin-testing/terraform"
	abstractions "github.com/microsoft/kiota-abstractions-go"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/acceptance"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

// NetworkManagedTLSCertificateTestResource implements the acceptance TestResource interface.
type NetworkManagedTLSCertificateTestResource struct{}

// Exists checks whether the Microsoft-managed TLS certificate authority exists in Microsoft Graph.
func (NetworkManagedTLSCertificateTestResource) Exists(ctx context.Context, _ any, state *terraform.InstanceState) (*bool, error) {
	graphClient, err := acceptance.TestGraphClient()
	if err != nil {
		return nil, err
	}

	resource := NetworkManagedTLSCertificateResource{client: graphClient}
	readAdapter := resource.client.GetAdapter()
	readRequest, err := newManagedTLSCertificateRequestInformation(ctx, readAdapter, abstractions.GET, state.ID, nil)
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
	readResult, err := readAdapter.Send(ctx, readRequest, createManagedTLSCertificateResponseFromDiscriminatorValue, managedTLSCertificateErrorMapping)
	_, err = decodeSendManagedTLSCertificateResponse(readResult, err)
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
