package graphBetaNetworkManagedTLSCertificate

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

func (d *NetworkManagedTLSCertificateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var object NetworkManagedTLSCertificateDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	readAdapter := d.client.GetAdapter()
	readRequest := abstractions.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(
		abstractions.GET,
		managedTLSCertificateItemURLTemplate,
		map[string]string{
			"baseurl":                 readAdapter.GetBaseUrl(),
			"managedTLSCertificateId": object.ID.ValueString(),
		},
	)
	readRequest.Headers.TryAdd("Accept", "application/json")
	// Without this header Graph returns status "disabled" as "unknownFutureValue".
	readRequest.Headers.TryAdd("Prefer", "include-unknown-enum-members")
	readResult, err := readAdapter.Send(ctx, readRequest, createManagedTLSCertificateResponseFromDiscriminatorValue, managedTLSCertificateErrorMapping)
	certificate, err := decodeGetManagedTLSCertificateResponse(readResult, err)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
		return
	}

	mapRemoteStateToDataSource(&object, certificate)
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
}

const managedTLSCertificateItemURLTemplate = "{+baseurl}/networkaccess/tls/managedCertificateAuthorityCertificates/{managedTLSCertificateId}"

var managedTLSCertificateErrorMapping = abstractions.ErrorMappings{"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue}
