package graphBetaIdentityAndAccessSecurityDefaultsEnforcementPolicy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
)

// Create adopts the tenant security defaults enforcement policy and applies the configuration.
//
// Operation: The security defaults enforcement policy is a singleton that always exists in the tenant and has no
// POST endpoint. Create applies is_enabled with PATCH.
// API Calls:
//   - PATCH /policies/identitySecurityDefaultsEnforcementPolicy
//
// Reference: https://learn.microsoft.com/en-us/graph/api/identitysecuritydefaultsenforcementpolicy-update?view=graph-rest-beta
func (r *SecurityDefaultsEnforcementPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object SecurityDefaultsEnforcementPolicyResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, r.client, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.
		Policies().
		IdentitySecurityDefaultsEnforcementPolicy().
		Patch(ctx, requestBody, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		addEnableConflictDiagnostic(ctx, err, object.IsEnabled, &resp.Diagnostics)
		return
	}

	// This is a singleton resource: set the static ID so the subsequent Read can locate the state entry.
	object.ID = types.StringValue(singletonID)

	tflog.Debug(ctx, fmt.Sprintf("Successfully adopted %s with ID: %s", ResourceName, object.ID.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationCreate
	opts.ResourceTypeName = ResourceName

	err = crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after create",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s", ResourceName))
}

// Read handles the Read operation for the security defaults enforcement policy resource.
//
// Operation: Retrieves the tenant security defaults enforcement policy
// API Calls:
//   - GET /policies/identitySecurityDefaultsEnforcementPolicy
//
// Reference: https://learn.microsoft.com/en-us/graph/api/identitysecuritydefaultsenforcementpolicy-get?view=graph-rest-beta
func (r *SecurityDefaultsEnforcementPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object SecurityDefaultsEnforcementPolicyResourceModel
	var identity sharedmodels.ResourceIdentity

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", ResourceName))

	operation := constants.TfOperationRead
	if ctxOp := ctx.Value("retry_operation"); ctxOp != nil {
		if opStr, ok := ctxOp.(string); ok {
			operation = opStr
		}
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading %s with ID: %s", ResourceName, object.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	identity.ID = singletonID

	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, identity)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	remoteResource, err := r.client.
		Policies().
		IdentitySecurityDefaultsEnforcementPolicy().
		Get(ctx, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	object = MapRemoteResourceStateToTerraform(ctx, object, remoteResource)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for the security defaults enforcement policy resource.
//
// Operation: Applies is_enabled to the tenant security defaults enforcement policy
// API Calls:
//   - PATCH /policies/identitySecurityDefaultsEnforcementPolicy
//
// Reference: https://learn.microsoft.com/en-us/graph/api/identitysecuritydefaultsenforcementpolicy-update?view=graph-rest-beta
func (r *SecurityDefaultsEnforcementPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan SecurityDefaultsEnforcementPolicyResourceModel
	var state SecurityDefaultsEnforcementPolicyResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Updating %s", ResourceName))

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

	requestBody, err := constructResource(ctx, r.client, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.
		Policies().
		IdentitySecurityDefaultsEnforcementPolicy().
		Patch(ctx, requestBody, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		addEnableConflictDiagnostic(ctx, err, plan.IsEnabled, &resp.Diagnostics)
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationUpdate
	opts.ResourceTypeName = ResourceName

	err = crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after update",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished updating %s with ID: %s", ResourceName, state.ID.ValueString()))
}

// Delete releases Terraform management without changing the tenant-wide security defaults enforcement policy.
//
// Operation: The singleton policy cannot be deleted. Delete only removes the resource from Terraform state and
// leaves security defaults as they are, so destroying the resource never enables or disables security defaults.
// API Calls: none
func (r *SecurityDefaultsEnforcementPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s", ResourceName))

	tflog.Info(ctx, fmt.Sprintf("Removing %s from Terraform state; the tenant security defaults setting is left unchanged", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
