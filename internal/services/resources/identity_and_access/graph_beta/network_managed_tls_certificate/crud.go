package graphBetaNetworkManagedTLSCertificate

import (
	"context"
	stderrors "errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
)

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
	body, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing resource for Create Method", err.Error())
		return
	}
	createAdapter := r.client.GetAdapter()
	createRequest, err := newManagedTLSCertificateRequestInformation(ctx, createAdapter, abstractions.POST, "", body)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}
	createResult, err := createAdapter.Send(ctx, createRequest, createManagedTLSCertificateResponseFromDiscriminatorValue, managedTLSCertificateErrorMapping)
	created, err := decodeCreateManagedTLSCertificateResponse(createResult, err)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}
	if created.id == nil {
		resp.Diagnostics.AddError("Error creating Microsoft-managed TLS certificate", "The API returned an invalid response without an id.")
		return
	}

	// Save the ID and computed values now, so a failure below leaves no unknowns.
	object.ID = types.StringValue(*created.id)
	mapComputed(ctx, &object, created)
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The Entra portal does not send status in the POST payload. A newly created
	// certificate is disabled by default, and Graph rejects an immediate explicit
	// disabled PATCH while it is provisioning. Only enabling requires the
	// distinct status PATCH and lifecycle wait.
	if desiredEnabled {
		statusBody, err := constructStatusUpdate(ctx, true)
		if err != nil {
			resp.Diagnostics.AddError("Error constructing status update after Create Method", err.Error())
			return
		}
		updateAdapter := r.client.GetAdapter()
		updateRequest, err := newManagedTLSCertificateRequestInformation(ctx, updateAdapter, abstractions.PATCH, object.ID.ValueString(), statusBody)
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
			return
		}
		if err := updateAdapter.SendNoContent(ctx, updateRequest, managedTLSCertificateErrorMapping); err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
			return
		}

		var settled *managedTLSCertificateResponse

		err = crud.PollUntil(ctx, managedTLSCertificateStatusPollInterval, func(ctx context.Context) (bool, error) {
			readbackAdapter := r.client.GetAdapter()
			readbackRequest, err := newManagedTLSCertificateRequestInformation(ctx, readbackAdapter, abstractions.GET, object.ID.ValueString(), nil)
			if err != nil {
				info := errors.GraphError(ctx, err)
				if info.StatusCode == 400 || info.StatusCode == 401 || info.StatusCode == 403 {
					return false, &crud.FatalPollError{Err: err}
				}
				return false, err
			}
			readbackResult, err := readbackAdapter.Send(ctx, readbackRequest, createManagedTLSCertificateResponseFromDiscriminatorValue, managedTLSCertificateErrorMapping)
			certificate, err := decodeSendManagedTLSCertificateResponse(readbackResult, err)
			if err != nil {
				info := errors.GraphError(ctx, err)
				if info.StatusCode == 400 || info.StatusCode == 401 || info.StatusCode == 403 {
					return false, &crud.FatalPollError{Err: err}
				}
				return false, err
			}
			settled = certificate
			if certificate.status == nil {
				return false, errMissingCertificateStatus
			}

			status := strings.ToLower(*certificate.status)
			switch status {
			case "expired", "revoked":
				return false, &crud.FatalPollError{Err: fmt.Errorf("%w: %q", errTerminalStatus, status)}
			case "active", "expiring":
				return true, nil
			default:
				return false, fmt.Errorf("%w with status %q", errPendingEnable, status)
			}
		})
		if err != nil {
			// settled is returned with the error so callers can report the last status seen.
			err = fmt.Errorf("wait for certificate status: %w", err)
		}
		if err != nil && keepAfterWaitError(err) {
			keepPlanned(ctx, &object, settled, err, &resp.Diagnostics)
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
	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationCreate
	opts.ResourceTypeName = ResourceName
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

	readAdapter := r.client.GetAdapter()
	readRequest, err := newManagedTLSCertificateRequestInformation(ctx, readAdapter, abstractions.GET, object.ID.ValueString(), nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}
	readResult, err := readAdapter.Send(ctx, readRequest, createManagedTLSCertificateResponseFromDiscriminatorValue, managedTLSCertificateErrorMapping)
	certificate, err := decodeSendManagedTLSCertificateResponse(readResult, err)
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
		body, err := constructStatusUpdate(ctx, plan.Enabled.ValueBool())
		if err != nil {
			resp.Diagnostics.AddError("Error constructing resource for Update Method", err.Error())
			return
		}
		updateAdapter := r.client.GetAdapter()
		updateRequest, err := newManagedTLSCertificateRequestInformation(ctx, updateAdapter, abstractions.PATCH, state.ID.ValueString(), body)
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}
		if err := updateAdapter.SendNoContent(ctx, updateRequest, managedTLSCertificateErrorMapping); err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}

		var settled *managedTLSCertificateResponse

		err = crud.PollUntil(ctx, managedTLSCertificateStatusPollInterval, func(ctx context.Context) (bool, error) {
			readbackAdapter := r.client.GetAdapter()
			readbackRequest, err := newManagedTLSCertificateRequestInformation(ctx, readbackAdapter, abstractions.GET, state.ID.ValueString(), nil)
			if err != nil {
				info := errors.GraphError(ctx, err)
				if info.StatusCode == 400 || info.StatusCode == 401 || info.StatusCode == 403 {
					return false, &crud.FatalPollError{Err: err}
				}
				return false, err
			}
			readbackResult, err := readbackAdapter.Send(ctx, readbackRequest, createManagedTLSCertificateResponseFromDiscriminatorValue, managedTLSCertificateErrorMapping)
			certificate, err := decodeSendManagedTLSCertificateResponse(readbackResult, err)
			if err != nil {
				info := errors.GraphError(ctx, err)
				if info.StatusCode == 400 || info.StatusCode == 401 || info.StatusCode == 403 {
					return false, &crud.FatalPollError{Err: err}
				}
				return false, err
			}
			settled = certificate
			if certificate.status == nil {
				return false, errMissingCertificateStatus
			}

			status := strings.ToLower(*certificate.status)
			switch {
			case status == "expired" || status == "revoked":
				return false, &crud.FatalPollError{Err: fmt.Errorf("%w: %q", errTerminalStatus, status)}
			case plan.Enabled.ValueBool() && (status == "active" || status == "expiring"):
				return true, nil
			case !plan.Enabled.ValueBool() && status == "disabled":
				return true, nil
			case plan.Enabled.ValueBool():
				return false, fmt.Errorf("%w with status %q", errPendingEnable, status)
			default:
				return false, fmt.Errorf("%w with status %q", errPendingDisable, status)
			}
		})
		if err != nil {
			// settled is returned with the error so callers can report the last status seen.
			err = fmt.Errorf("wait for certificate status: %w", err)
		}
		if err != nil && keepAfterWaitError(err) {
			plan.ID = state.ID
			keepPlanned(ctx, &plan, settled, err, &resp.Diagnostics)
			resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
			return
		}
		if err != nil {
			resp.Diagnostics.AddError("Error waiting for Microsoft-managed TLS certificate status update", err.Error())
			return
		}
		MapRemoteStateToTerraform(ctx, &plan, settled)
	}

	plan.ID = state.ID
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}
	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationUpdate
	opts.ResourceTypeName = ResourceName
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

	deleteAdapter := r.client.GetAdapter()
	deleteRequest, err := newManagedTLSCertificateRequestInformation(ctx, deleteAdapter, abstractions.DELETE, object.ID.ValueString(), nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}
	if err := deleteAdapter.SendNoContent(ctx, deleteRequest, managedTLSCertificateErrorMapping); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}
	resp.State.RemoveResource(ctx)
}

// keepAfterWaitError reports whether a failed status wait should keep the
// certificate. Only a terminal status is fatal. After a timeout, cancel or
// transient error Graph may still finish, and failing Create would taint the
// resource, whose replacement issues a new root CA.
func keepAfterWaitError(err error) bool {
	return !stderrors.Is(err, errTerminalStatus)
}

// keepPlanned records the planned enabled value with the last status seen and
// warns that the change has not finished. The next refresh reads the actual
// status, and the next apply retries if it differs.
func keepPlanned(ctx context.Context, object *NetworkManagedTLSCertificateResourceModel, latest *managedTLSCertificateResponse, err error, diags *diag.Diagnostics) {
	mapComputed(ctx, object, latest)
	diags.AddWarning(
		"Microsoft-managed TLS certificate status change has not finished",
		fmt.Sprintf("Certificate %s was set to enabled = %t, but apply stopped waiting before Graph finished (last status %q): %s. "+
			"The next refresh shows the actual status, and the next apply retries if needed. "+
			"Raise timeouts.create or timeouts.update to wait longer.",
			object.ID.ValueString(), object.Enabled.ValueBool(), object.Status.ValueString(), err),
	)
}

// mapComputed copies the computed values from a response, leaving configured
// attributes as planned, and replaces any value Graph did not return with null.
func mapComputed(ctx context.Context, object *NetworkManagedTLSCertificateResourceModel, remote *managedTLSCertificateResponse) {
	if remote != nil {
		remoteState := *object
		MapRemoteStateToTerraform(ctx, &remoteState, remote)
		object.Status = remoteState.Status
		object.CreatedDateTime = remoteState.CreatedDateTime
		object.ValidityStartDateTime = remoteState.ValidityStartDateTime
		object.ValidityEndDateTime = remoteState.ValidityEndDateTime
		object.Certificate = remoteState.Certificate
	}
	for _, v := range []*types.String{
		&object.Status, &object.CreatedDateTime, &object.ValidityStartDateTime,
		&object.ValidityEndDateTime, &object.Certificate,
	} {
		if v.IsUnknown() {
			*v = types.StringNull()
		}
	}
}

var managedTLSCertificateErrorMapping = abstractions.ErrorMappings{"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue}

const managedTLSCertificateStatusPollInterval = 5 * time.Second

var (
	// errTerminalStatus marks a certificate that can no longer reach the requested status.
	errTerminalStatus = stderrors.New("certificate entered a terminal status")
	errPendingEnable  = stderrors.New("certificate activation is still pending")
	errPendingDisable = stderrors.New("certificate deactivation is still pending")
)

var errMissingCertificateStatus = stderrors.New("certificate response from Microsoft Graph has no lifecycle status")
