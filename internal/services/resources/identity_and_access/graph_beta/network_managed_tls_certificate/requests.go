package graphBetaNetworkManagedTLSCertificate

import (
	"context"
	"fmt"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	s "github.com/microsoft/kiota-abstractions-go/serialization"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
)

const (
	managedTLSCertificatesURLTemplate    = "{+baseurl}/networkaccess/tls/managedCertificateAuthorityCertificates"
	managedTLSCertificateItemURLTemplate = managedTLSCertificatesURLTemplate + "/{managedTLSCertificateId}"
)

var managedTLSCertificateErrorMapping = abstractions.ErrorMappings{"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue}

func (r *NetworkManagedTLSCertificateResource) createManagedTLSCertificate(ctx context.Context, requestBody s.Parsable) (*managedTLSCertificateResponse, error) {
	result, err := r.sendManagedTLSCertificate(ctx, abstractions.POST, "", requestBody)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, fmt.Errorf("create Microsoft-managed TLS certificate returned nil response")
	}
	return result, nil
}

func (r *NetworkManagedTLSCertificateResource) getManagedTLSCertificate(ctx context.Context, certificateID string) (*managedTLSCertificateResponse, error) {
	return r.sendManagedTLSCertificate(ctx, abstractions.GET, certificateID, nil)
}

func (r *NetworkManagedTLSCertificateResource) sendManagedTLSCertificate(ctx context.Context, method abstractions.HttpMethod, certificateID string, requestBody s.Parsable) (*managedTLSCertificateResponse, error) {
	adapter := r.client.GetAdapter()
	requestInfo, err := newManagedTLSCertificateRequestInformation(ctx, adapter, method, certificateID, requestBody)
	if err != nil {
		return nil, err
	}
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

func (r *NetworkManagedTLSCertificateResource) updateManagedTLSCertificate(ctx context.Context, certificateID string, requestBody s.Parsable) error {
	adapter := r.client.GetAdapter()
	requestInfo, err := newManagedTLSCertificateRequestInformation(ctx, adapter, abstractions.PATCH, certificateID, requestBody)
	if err != nil {
		return err
	}
	return adapter.SendNoContent(ctx, requestInfo, managedTLSCertificateErrorMapping)
}

func (r *NetworkManagedTLSCertificateResource) deleteManagedTLSCertificate(ctx context.Context, certificateID string) error {
	adapter := r.client.GetAdapter()
	requestInfo, err := newManagedTLSCertificateRequestInformation(ctx, adapter, abstractions.DELETE, certificateID, nil)
	if err != nil {
		return err
	}
	return adapter.SendNoContent(ctx, requestInfo, managedTLSCertificateErrorMapping)
}

func newManagedTLSCertificateRequestInformation(ctx context.Context, adapter abstractions.RequestAdapter, method abstractions.HttpMethod, certificateID string, requestBody s.Parsable) (*abstractions.RequestInformation, error) {
	pathParameters := map[string]string{"baseurl": adapter.GetBaseUrl()}
	urlTemplate := managedTLSCertificatesURLTemplate
	if certificateID != "" {
		urlTemplate = managedTLSCertificateItemURLTemplate
		pathParameters["managedTLSCertificateId"] = certificateID
	}

	requestInfo := abstractions.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(method, urlTemplate, pathParameters)
	requestInfo.Headers.TryAdd("Accept", "application/json")
	if requestBody != nil {
		if err := requestInfo.SetContentFromParsable(ctx, adapter, "application/json", requestBody); err != nil {
			return nil, fmt.Errorf("set Microsoft-managed TLS certificate request content: %w", err)
		}
	}
	return requestInfo, nil
}
