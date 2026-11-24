// Package graphBetaNetworkManagedTLSCertificate manages the Microsoft-managed certificate authority used by Global
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
//   - Observed: the Entra portal creates the certificate with name, commonName, organizationName and
//     validityMonths, and changes status with a status-only PATCH. Without
//     Prefer: include-unknown-enum-members Graph returns status "disabled" as "unknownFutureValue", so the header is
//     sent on every request. Requests go through custom_requests in raw JSON mode with the plain structs in model.go.
//
// Removal condition: the beta SDK ships a request builder and model for
// networkaccess/tls/managedCertificateAuthorityCertificates whose status enum includes the lifecycle values this
// resource reads (creating, enrolling, active, expiring, expired, revoked, disabled).
package graphBetaNetworkManagedTLSCertificate

import (
	"context"
	stderrors "errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
)

const managedTLSCertificateStatusPollInterval = 5 * time.Second

// Warning for a status change that apply stopped waiting for. It records the planned enabled value with the last
// status seen; the next refresh reads the actual status, and the next apply retries if it differs.
const (
	statusChangePendingSummary = "Microsoft-managed TLS certificate status change has not finished"
	statusChangePendingDetail  = "Certificate %s was set to enabled = %t, but apply stopped waiting before Graph finished (last status %q): %s. " +
		"The next refresh shows the actual status, and the next apply retries if needed. " +
		"Raise timeouts.create or timeouts.update to wait longer."
)

var (
	// errTerminalStatus marks a certificate that can no longer reach the requested status. Only a terminal status
	// fails a status wait: after a timeout, cancel or transient error Graph may still finish, and failing Create
	// would taint the resource, whose replacement issues a new root CA.
	errTerminalStatus = stderrors.New("certificate entered a terminal status")
	errPendingEnable  = stderrors.New("certificate activation is still pending")
	errPendingDisable = stderrors.New("certificate deactivation is still pending")
	errNoStatus       = stderrors.New("the certificate Graph returned has no lifecycle status")
	errEmptyResponse  = stderrors.New("the Microsoft-managed TLS certificate request returned no response body")
)

// managedTLSCertificateHeaders are sent with every request. Status values such as "disabled" are evolvable enum
// members; without the Prefer header Graph returns them as "unknownFutureValue".
var managedTLSCertificateHeaders = map[string]string{"Prefer": "include-unknown-enum-members"}

func (r *NetworkManagedTLSCertificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object NetworkManagedTLSCertificateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	desiredEnabled := object.Enabled.ValueBool()
	requestBody, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing resource for Create Method", err.Error())
		return
	}

	// A response with no body leaves created nil.
	var created *managedTLSCertificateResponse
	createConfig := customrequests.PostRequestConfig{
		APIVersion:  customrequests.GraphAPIBeta,
		Endpoint:    r.ResourcePath,
		RawBody:     requestBody,
		Headers:     managedTLSCertificateHeaders,
		RawResponse: &created,
	}
	_, err = customrequests.PostRequest(ctx, r.client.GetAdapter(), createConfig, nil, nil)
	if err == nil && created == nil {
		err = errEmptyResponse
	}
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}
	if created.ID == nil {
		resp.Diagnostics.AddError("Error creating Microsoft-managed TLS certificate", "The API returned an invalid response without an id.")
		return
	}

	// Save the ID and computed values now, so a failure below leaves no unknowns.
	object.ID = types.StringValue(*created.ID)
	mapManagedTLSCertificateComputedToState(ctx, &object, created)
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The Entra portal does not send status in the POST payload. A newly created
	// certificate is disabled by default, and Graph rejects an immediate explicit
	// disabled PATCH while it is provisioning. Only enabling requires the
	// distinct status PATCH and lifecycle wait.
	if desiredEnabled {
		status := "enabled"
		statusBody := managedTLSCertificateRequest{Status: &status}
		tflog.Debug(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), map[string]any{"request_body": statusBody})

		patchConfig := customrequests.PatchRequestConfig{
			APIVersion:        customrequests.GraphAPIBeta,
			Endpoint:          r.ResourcePath,
			ResourceIDPattern: customrequests.ResourceIDPatternSegment,
			ResourceID:        object.ID.ValueString(),
			RawBody:           statusBody,
			Headers:           managedTLSCertificateHeaders,
		}
		if err := customrequests.PatchRequestByResourceId(ctx, r.client.GetAdapter(), patchConfig); err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
			return
		}

		settled, err := r.waitForManagedTLSCertificateStatus(ctx, object.ID.ValueString(), true)
		if err != nil && !stderrors.Is(err, errTerminalStatus) {
			mapManagedTLSCertificateComputedToState(ctx, &object, settled)
			detail := fmt.Sprintf(statusChangePendingDetail, object.ID.ValueString(), object.Enabled.ValueBool(), object.Status.ValueString(), err)
			resp.Diagnostics.AddWarning(statusChangePendingSummary, detail)
			resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
			if resp.Identity != nil {
				resp.Diagnostics.Append(resp.Identity.Set(ctx, sharedmodels.ResourceIdentity{ID: object.ID.ValueString()})...)
			}
			return
		}
		if err != nil {
			resp.Diagnostics.AddError("Error waiting for Microsoft-managed TLS certificate activation", err.Error())
			return
		}
		MapRemoteStateToTerraform(ctx, &object, settled)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	opts := crud.NewReadWithRetryOptions(constants.TfOperationCreate, ResourceName)
	if err := crud.ReadWithRetry(ctx, r.Read, readReq, &crud.CreateResponseContainer{CreateResponse: resp}, opts); err != nil {
		resp.Diagnostics.AddError("Error reading resource state after create", fmt.Sprintf("Could not read resource state: %s", err))
	}
}

func (r *NetworkManagedTLSCertificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object NetworkManagedTLSCertificateResourceModel
	var identity sharedmodels.ResourceIdentity
	operation := constants.TfOperationRead
	if op, ok := ctx.Value("retry_operation").(string); ok {
		operation = op
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	identity.ID = object.ID.ValueString()
	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, identity)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	var certificate *managedTLSCertificateResponse
	getConfig := customrequests.GetRequestConfig{
		APIVersion:        customrequests.GraphAPIBeta,
		Endpoint:          r.ResourcePath,
		ResourceIDPattern: customrequests.ResourceIDPatternSegment,
		ResourceID:        object.ID.ValueString(),
		Headers:           managedTLSCertificateHeaders,
		RawResponse:       &certificate,
	}
	_, err := customrequests.GetRequestByResourceId(ctx, r.client.GetAdapter(), getConfig)
	if err == nil && certificate == nil {
		err = errEmptyResponse
	}
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}
	MapRemoteStateToTerraform(ctx, &object, certificate)
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
}

func (r *NetworkManagedTLSCertificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state NetworkManagedTLSCertificateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	if !plan.Enabled.Equal(state.Enabled) {
		status := "disabled"
		if plan.Enabled.ValueBool() {
			status = "enabled"
		}
		statusBody := managedTLSCertificateRequest{Status: &status}
		tflog.Debug(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), map[string]any{"request_body": statusBody})

		patchConfig := customrequests.PatchRequestConfig{
			APIVersion:        customrequests.GraphAPIBeta,
			Endpoint:          r.ResourcePath,
			ResourceIDPattern: customrequests.ResourceIDPatternSegment,
			ResourceID:        state.ID.ValueString(),
			RawBody:           statusBody,
			Headers:           managedTLSCertificateHeaders,
		}
		if err := customrequests.PatchRequestByResourceId(ctx, r.client.GetAdapter(), patchConfig); err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}

		settled, err := r.waitForManagedTLSCertificateStatus(ctx, state.ID.ValueString(), plan.Enabled.ValueBool())
		if err != nil && !stderrors.Is(err, errTerminalStatus) {
			// The planned enabled value is kept with the last status seen, so the next apply retries the change.
			plan.ID = state.ID
			mapManagedTLSCertificateComputedToState(ctx, &plan, settled)
			detail := fmt.Sprintf(statusChangePendingDetail, plan.ID.ValueString(), plan.Enabled.ValueBool(), plan.Status.ValueString(), err)
			resp.Diagnostics.AddWarning(statusChangePendingSummary, detail)
			resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
			return
		}
		if err != nil {
			resp.Diagnostics.AddError("Error waiting for Microsoft-managed TLS certificate status update", err.Error())
			return
		}
		MapRemoteStateToTerraform(ctx, &state, settled)
	}

	// Keep the prior known state if the read-back fails after a successful update.
	object := state
	object.Timeouts = plan.Timeouts

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	opts := crud.NewReadWithRetryOptions(constants.TfOperationUpdate, ResourceName)
	if err := crud.ReadWithRetry(ctx, r.Read, readReq, &crud.UpdateResponseContainer{UpdateResponse: resp}, opts); err != nil {
		resp.Diagnostics.AddError("Error reading resource state after update", fmt.Sprintf("Could not read resource state: %s", err))
	}
}

func (r *NetworkManagedTLSCertificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object NetworkManagedTLSCertificateResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	deleteConfig := customrequests.DeleteRequestConfig{
		APIVersion:        customrequests.GraphAPIBeta,
		Endpoint:          r.ResourcePath,
		ResourceIDPattern: customrequests.ResourceIDPatternSegment,
		ResourceID:        object.ID.ValueString(),
		Headers:           managedTLSCertificateHeaders,
	}
	if err := customrequests.DeleteRequestByResourceId(ctx, r.client.GetAdapter(), deleteConfig); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}
	resp.State.RemoveResource(ctx)
}

// waitForManagedTLSCertificateStatus polls the certificate until its lifecycle status reflects enabled. The last
// certificate read is returned with any error, so callers can report the last status seen.
func (r *NetworkManagedTLSCertificateResource) waitForManagedTLSCertificateStatus(ctx context.Context, certificateID string, enabled bool) (*managedTLSCertificateResponse, error) {
	var latest *managedTLSCertificateResponse

	getConfig := customrequests.GetRequestConfig{
		APIVersion:        customrequests.GraphAPIBeta,
		Endpoint:          r.ResourcePath,
		ResourceIDPattern: customrequests.ResourceIDPatternSegment,
		ResourceID:        certificateID,
		Headers:           managedTLSCertificateHeaders,
	}
	err := crud.PollUntil(ctx, managedTLSCertificateStatusPollInterval, func(ctx context.Context) (bool, error) {
		var certificate *managedTLSCertificateResponse
		getConfig.RawResponse = &certificate
		_, err := customrequests.GetRequestByResourceId(ctx, r.client.GetAdapter(), getConfig)
		if err == nil && certificate == nil {
			err = errEmptyResponse
		}
		if err != nil {
			info := errors.GraphError(ctx, err)
			if info.StatusCode == 400 || info.StatusCode == 401 || info.StatusCode == 403 {
				return false, &crud.FatalPollError{Err: err}
			}
			return false, err
		}
		latest = certificate
		if certificate.Status == nil {
			return false, errNoStatus
		}

		status := strings.ToLower(*certificate.Status)
		switch {
		case status == "expired" || status == "revoked":
			return false, &crud.FatalPollError{Err: fmt.Errorf("%w: %q", errTerminalStatus, status)}
		case enabled && (status == "active" || status == "expiring"):
			return true, nil
		case !enabled && status == "disabled":
			return true, nil
		case enabled:
			return false, fmt.Errorf("%w with status %q", errPendingEnable, status)
		default:
			return false, fmt.Errorf("%w with status %q", errPendingDisable, status)
		}
	})
	if err != nil {
		return latest, fmt.Errorf("wait for certificate status: %w", err)
	}
	return latest, nil
}
