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

// readbackMatch is what the read-back after a write waits for: the written policy settings.
var readbackMatch = crud.AttributeMatchOptions{Paths: crud.RootPaths("is_enabled", "notify_reviewers", "reminders_enabled", "request_duration_in_days", "reviewers")}

// Create handles the Create operation for admin consent request policy resources.
//
// Operation: Applies the complete writable configuration to the existing tenant singleton
// API Calls:
//   - PUT /policies/adminConsentRequestPolicy
//
// Reference: https://learn.microsoft.com/en-us/graph/api/adminconsentrequestpolicy-update?view=graph-rest-beta
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

// Read handles the Read operation for admin consent request policy resources.
//
// Operation: Retrieves the tenant admin consent request policy
// API Calls:
//   - GET /policies/adminConsentRequestPolicy
//
// Reference: https://learn.microsoft.com/en-us/graph/api/adminconsentrequestpolicy-get?view=graph-rest-beta
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

// Update handles the Update operation for admin consent request policy resources.
//
// Operation: Applies the complete writable configuration to the existing tenant singleton
// API Calls:
//   - PUT /policies/adminConsentRequestPolicy
//
// Reference: https://learn.microsoft.com/en-us/graph/api/adminconsentrequestpolicy-update?view=graph-rest-beta
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

// Delete handles the Delete operation for admin consent request policy resources.
//
// Operation: Removes Terraform state while retaining the tenant policy
// API Calls:
//   - None (state only)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/resources/adminconsentrequestpolicy?view=graph-rest-beta
func (r *AdminConsentRequestPolicyResource) Delete(ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.State.RemoveResource(ctx)
}
