package graphBetaNetworkManagedTLSCertificate

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
	created, err := r.createManagedTLSCertificate(ctx, body)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}
	if created.id == nil {
		resp.Diagnostics.AddError("Error creating Microsoft-managed TLS certificate", "The API returned an invalid response without an id.")
		return
	}

	object.ID = types.StringValue(*created.id)
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
		if err := r.updateManagedTLSCertificate(ctx, object.ID.ValueString(), statusBody); err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
			return
		}

		settled, err := r.waitForManagedTLSCertificateStatus(ctx, object.ID.ValueString(), true)
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

	certificate, err := r.getManagedTLSCertificate(ctx, object.ID.ValueString())
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
		if err := r.updateManagedTLSCertificate(ctx, state.ID.ValueString(), body); err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}

		settled, err := r.waitForManagedTLSCertificateStatus(ctx, state.ID.ValueString(), plan.Enabled.ValueBool())
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

	if err := r.deleteManagedTLSCertificate(ctx, object.ID.ValueString()); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}
	resp.State.RemoveResource(ctx)
}
