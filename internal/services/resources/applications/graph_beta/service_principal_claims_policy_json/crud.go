package graphBetaServicePrincipalClaimsPolicyJSON

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	grapherrors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
)

var graphErrorMapping = abstractions.ErrorMappings{"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue}

func (r *ServicePrincipalClaimsPolicyJSONResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ServicePrincipalClaimsPolicyJSONResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	content, err := constructResource(data.PolicyJSON)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing custom claims policy", err.Error())
		return
	}
	request, err := r.client.ServicePrincipals().
		ByServicePrincipalId(data.ServicePrincipalID.ValueString()).
		ClaimsPolicy().
		ToPutRequestInformation(ctx, graphmodels.NewCustomClaimsPolicy(), nil)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing claims policy request", err.Error())
		return
	}
	// Keep raw JSON so beta fields such as groupFilter and all discriminators survive.
	request.Content = content
	err = r.client.GetAdapter().SendNoContent(ctx, request, graphErrorMapping)
	if err != nil {
		grapherrors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}
	data.ID = types.StringValue(data.ServicePrincipalID.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationCreate
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = policyConsistencyPredicate(&data)
	if deadline, ok := ctx.Deadline(); ok {
		opts.MaxRetries = int(time.Until(deadline)/opts.RetryInterval) + 1
	}
	if err = crud.ReadWithRetry(
		ctx,
		r.Read,
		resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta},
		&crud.CreateResponseContainer{CreateResponse: resp},
		opts,
	); err != nil {
		resp.Diagnostics.AddError("Error reading custom claims policy after create", err.Error())
	}
}

func (r *ServicePrincipalClaimsPolicyJSONResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ServicePrincipalClaimsPolicyJSONResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	request, err := r.client.ServicePrincipals().ByServicePrincipalId(data.ServicePrincipalID.ValueString()).ClaimsPolicy().ToGetRequestInformation(ctx, nil)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing claims policy read", err.Error())
		return
	}
	operation := constants.TfOperationRead
	if op, ok := ctx.Value("retry_operation").(string); ok {
		operation = op
	}
	var remote any
	missing := 0
	err = crud.PollUntil(ctx, 2*time.Second, func(ctx context.Context) (bool, error) {
		var getErr error
		remote, getErr = r.client.GetAdapter().SendPrimitive(ctx, request, "[]byte", graphErrorMapping)
		if getErr == nil {
			return true, nil
		}
		if grapherrors.GraphError(ctx, getErr).StatusCode == 404 && operation == constants.TfOperationRead {
			missing++
			return missing >= 3, nil
		}
		return false, &crud.FatalPollError{Err: getErr}
	})
	if err != nil {
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
	if missing >= 3 {
		resp.State.RemoveResource(ctx)
		return
	}
	content, ok := remote.([]byte)
	if !ok || len(content) == 0 {
		resp.Diagnostics.AddError("Invalid custom claims policy response", errPolicyResponse.Error())
		return
	}
	if err = MapRemoteStateToTerraform(&data, content); err != nil {
		resp.Diagnostics.AddError("Invalid custom claims policy response", err.Error())
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, sharedmodels.ResourceIdentity{ID: data.ID.ValueString()})...)
	}
}

func (r *ServicePrincipalClaimsPolicyJSONResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ServicePrincipalClaimsPolicyJSONResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	content, err := constructResource(data.PolicyJSON)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing custom claims policy", err.Error())
		return
	}
	request, err := r.client.ServicePrincipals().
		ByServicePrincipalId(data.ServicePrincipalID.ValueString()).
		ClaimsPolicy().
		ToPutRequestInformation(ctx, graphmodels.NewCustomClaimsPolicy(), nil)
	if err != nil {
		resp.Diagnostics.AddError("Error constructing claims policy request", err.Error())
		return
	}
	request.Content = content
	err = r.client.GetAdapter().SendNoContent(ctx, request, graphErrorMapping)
	if err != nil {
		grapherrors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}
	data.ID = types.StringValue(data.ServicePrincipalID.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationUpdate
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = policyConsistencyPredicate(&data)
	if deadline, ok := ctx.Deadline(); ok {
		opts.MaxRetries = int(time.Until(deadline)/opts.RetryInterval) + 1
	}
	if err = crud.ReadWithRetry(
		ctx,
		r.Read,
		resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta},
		&crud.UpdateResponseContainer{UpdateResponse: resp},
		opts,
	); err != nil {
		resp.Diagnostics.AddError("Error reading custom claims policy after update", err.Error())
	}
}

// Beta DELETE returns 405. Releasing management must retain the SAML policy;
// PUT {} changes token settings and is not equivalent to removing the policy.
func (r *ServicePrincipalClaimsPolicyJSONResource) Delete(ctx context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.State.RemoveResource(ctx)
}
