// Package graphBetaNetworkManagedTLSCertificate reads the Microsoft-managed certificate authority used by Global
// Secure Access TLS inspection.
//
// SDK deviation: missing type
//
// Evidence:
//   - SDK: github.com/microsoftgraph/msgraph-beta-sdk-go@v0.168.0 has no request builder for
//     /networkaccess/tls/managedCertificateAuthorityCertificates and no model for its items. The networkaccess
//     package builds only tls/externalCertificateAuthorityCertificates and tls/tlsInspectionPolicies.
//   - Learn: the endpoint is not documented on Microsoft Learn; it is the endpoint the Entra Global Secure Access
//     Microsoft-managed TLS certificate blade calls.
//   - Observed: without Prefer: include-unknown-enum-members Graph returns status "disabled" as
//     "unknownFutureValue", so the header is sent with the GET. The request goes through custom_requests in raw JSON
//     mode with the plain structs in model.go.
//
// Removal condition: the beta SDK ships a request builder and model for
// networkaccess/tls/managedCertificateAuthorityCertificates whose status enum includes the lifecycle values Graph
// returns (creating, enrolling, active, expiring, expired, revoked, disabled).
package graphBetaNetworkManagedTLSCertificate

import (
	"context"
	stderrors "errors"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

var errEmptyResponse = stderrors.New("the Microsoft-managed TLS certificate request returned no response body")

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

	// A response with no body leaves certificate nil.
	var certificate *managedTLSCertificateResponse
	getConfig := customrequests.GetRequestConfig{
		APIVersion:        customrequests.GraphAPIBeta,
		Endpoint:          "networkaccess/tls/managedCertificateAuthorityCertificates",
		ResourceIDPattern: customrequests.ResourceIDPatternSegment,
		ResourceID:        object.ID.ValueString(),
		Headers:           map[string]string{"Prefer": "include-unknown-enum-members"},
		RawResponse:       &certificate,
	}
	_, err := customrequests.GetRequestByResourceId(ctx, d.client.GetAdapter(), getConfig)
	if err == nil && certificate == nil {
		err = errEmptyResponse
	}
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationRead, d.ReadPermissions)
		return
	}

	MapRemoteStateToDataSource(&object, certificate)
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
}
