package graphBetaAdminConsentRequestPolicy

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	grapherrors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
)

func (r *AdminConsentRequestPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data AdminConsentRequestPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	body, err := constructResource(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing admin consent request policy", err.Error())
		return
	}
	request, err := r.client.Policies().AdminConsentRequestPolicy().ToPatchRequestInformation(ctx, body, nil)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing admin consent request", err.Error())
		return
	}
	// The service supports PUT with a bodyless 204; generated PATCH is unsupported.
	request.Method = abstractions.PUT
	err = r.client.GetAdapter().SendNoContent(ctx, request, abstractions.ErrorMappings{"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue})
	if err != nil {
		grapherrors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}
	data.ID = types.StringValue(singletonID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationCreate
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = policyConsistencyPredicate(&data)
	if err = crud.ReadWithRetry(
		ctx,
		r.Read,
		resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta},
		&crud.CreateResponseContainer{CreateResponse: resp},
		opts,
	); err != nil {
		resp.Diagnostics.AddError("Error reading admin consent request policy after create", err.Error())
	}
}

func (r *AdminConsentRequestPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data AdminConsentRequestPolicyResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	remote, err := r.client.Policies().AdminConsentRequestPolicy().Get(ctx, nil)
	if err != nil {
		operation := constants.TfOperationRead
		if op, ok := ctx.Value("retry_operation").(string); ok {
			operation = op
		}
		grapherrors.HandleKiotaGraphErrorWithOptions(
			ctx,
			err,
			resp,
			operation,
			r.ReadPermissions,
			grapherrors.GraphErrorOptions{PreserveStateOnReadBadRequest: true, PreserveStateOnReadNotFound: true},
		)
		return
	}
	if remote == nil || remote.GetIsEnabled() == nil || remote.GetNotifyReviewers() == nil || remote.GetRemindersEnabled() == nil || remote.GetRequestDurationInDays() == nil {
		resp.Diagnostics.AddError("Invalid admin consent request policy response", "Microsoft Graph returned an incomplete tenant policy.")
		return
	}
	resp.Diagnostics.Append(MapRemoteStateToTerraform(ctx, &data, remote)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, sharedmodels.ResourceIdentity{ID: singletonID})...)
	}
}

func (r *AdminConsentRequestPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data AdminConsentRequestPolicyResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	body, err := constructResource(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing admin consent request policy", err.Error())
		return
	}
	request, err := r.client.Policies().AdminConsentRequestPolicy().ToPatchRequestInformation(ctx, body, nil)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing admin consent request", err.Error())
		return
	}
	request.Method = abstractions.PUT
	err = r.client.GetAdapter().SendNoContent(ctx, request, abstractions.ErrorMappings{"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue})
	if err != nil {
		grapherrors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}
	data.ID = types.StringValue(singletonID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationUpdate
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = policyConsistencyPredicate(&data)
	if err = crud.ReadWithRetry(
		ctx,
		r.Read,
		resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta},
		&crud.UpdateResponseContainer{UpdateResponse: resp},
		opts,
	); err != nil {
		resp.Diagnostics.AddError("Error reading admin consent request policy after update", err.Error())
	}
}

// Delete releases management of the tenant singleton without changing the policy.
func (r *AdminConsentRequestPolicyResource) Delete(ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.State.RemoveResource(ctx)
}
