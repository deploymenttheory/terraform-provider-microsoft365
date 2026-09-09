package graphBetaNetworkManagedTLSCertificate

import (
	"context"
	"fmt"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
)

const managedTLSCertificateItemURLTemplate = "{+baseurl}/networkaccess/tls/managedCertificateAuthorityCertificates/{managedTLSCertificateId}"

var managedTLSCertificateErrorMapping = abstractions.ErrorMappings{"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue}

func (d *NetworkManagedTLSCertificateDataSource) getManagedTLSCertificate(ctx context.Context, certificateID string) (*managedTLSCertificateResponse, error) {
	adapter := d.client.GetAdapter()
	requestInfo := abstractions.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(
		abstractions.GET,
		managedTLSCertificateItemURLTemplate,
		map[string]string{
			"baseurl":                 adapter.GetBaseUrl(),
			"managedTLSCertificateId": certificateID,
		},
	)
	requestInfo.Headers.TryAdd("Accept", "application/json")

	result, err := adapter.Send(ctx, requestInfo, createManagedTLSCertificateResponseFromDiscriminatorValue, managedTLSCertificateErrorMapping)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, fmt.Errorf("Microsoft-managed TLS certificate request returned nil response")
	}
	certificate, ok := result.(*managedTLSCertificateResponse)
	if !ok {
		return nil, fmt.Errorf("Microsoft-managed TLS certificate request returned %T, expected managedTLSCertificateResponse", result)
	}
	return certificate, nil
}
