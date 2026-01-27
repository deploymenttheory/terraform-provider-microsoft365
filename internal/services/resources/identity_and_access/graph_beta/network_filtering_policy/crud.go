package graphBetaNetworkFilteringPolicy

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/license"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Create handles the Create operation for Filtering Policy resources.
//
// Operation: Creates a new network filtering policy for Global Secure Access
// API Calls:
//   - POST /networkAccess/filteringPolicies
//
// Reference: https://learn.microsoft.com/en-us/graph/api/networkaccess-filteringprofile-post-policies?view=graph-rest-beta
// Note: Requires specific Microsoft Entra licensing for Global Secure Access
func (r *NetworkFilteringPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object NetworkFilteringPolicyResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate tenant has required license before attempting API call
	if !license.HasRequiredLicense(ctx, r.client, "NetworkFilteringPolicy") {
		resp.Diagnostics.AddError(
			"Missing Required License",
			fmt.Sprintf(
				"This resource requires a tenant license that was not found.\n\n%s",
				license.FormatRequiredLicensesMessage("NetworkFilteringPolicy"),
			),
		)
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Create Method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	baseResource, err := r.client.
		NetworkAccess().
		FilteringPolicies().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	object.ID = types.StringValue(*baseResource.GetId())
	tflog.Debug(ctx, fmt.Sprintf("Successfully created %s with ID: %s", ResourceName, *baseResource.GetId()))

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

// Read handles the Read operation for Filtering Policy resources.
//
// Operation: Retrieves a network filtering policy by ID
// API Calls:
//   - GET /networkAccess/filteringPolicies/{filteringPolicyId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/networkaccess-filteringpolicy-get?view=graph-rest-beta
func (r *NetworkFilteringPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object NetworkFilteringPolicyResourceModel

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

	policy, err := r.client.
		NetworkAccess().
		FilteringPolicies().
		ByFilteringPolicyId(object.ID.ValueString()).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, policy)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for Filtering Policy resources.
//
// Operation: Updates an existing network filtering policy
// API Calls:
//   - PATCH /networkAccess/filteringPolicies/{filteringPolicyId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/networkaccess-filteringpolicy-update?view=graph-rest-beta
func (r *NetworkFilteringPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan NetworkFilteringPolicyResourceModel
	var state NetworkFilteringPolicyResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Updating %s with ID: %s", ResourceName, state.ID.ValueString()))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)   // desired state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...) // current state (for ID)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate tenant has required license before attempting API call
	if !license.HasRequiredLicense(ctx, r.client, "NetworkFilteringPolicy") {
		resp.Diagnostics.AddError(
			"Missing Required License",
			fmt.Sprintf(
				"This resource requires a tenant license that was not found.\n\n%s",
				license.FormatRequiredLicensesMessage("NetworkFilteringPolicy"),
			),
		)
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Update Method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.
		NetworkAccess().
		FilteringPolicies().
		ByFilteringPolicyId(state.ID.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
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

// Delete handles the Delete operation for Filtering Policy resources.
//
// Operation: Deletes a network filtering policy
// API Calls:
//   - DELETE /networkAccess/filteringPolicies/{filteringPolicyId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/networkaccess-filteringprofile-delete-policies?view=graph-rest-beta
func (r *NetworkFilteringPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object NetworkFilteringPolicyResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate tenant has required license before attempting API call
	if !license.HasRequiredLicense(ctx, r.client, "NetworkFilteringPolicy") {
		resp.Diagnostics.AddError(
			"Missing Required License",
			fmt.Sprintf(
				"This resource requires a tenant license that was not found.\n\n%s",
				license.FormatRequiredLicensesMessage("NetworkFilteringPolicy"),
			),
		)
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	err := r.client.
		NetworkAccess().
		FilteringPolicies().
		ByFilteringPolicyId(object.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
