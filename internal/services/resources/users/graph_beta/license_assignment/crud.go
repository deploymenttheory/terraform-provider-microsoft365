package graphBetaUserLicenseAssignment

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/users"
)

// Create handles the Create operation for user license assignment resources.
//
// Operation: Assigns a license to a user
// API Calls:
//   - POST /users/{id}/assignLicense
//
// Reference: https://learn.microsoft.com/en-us/graph/api/user-assignlicense?view=graph-rest-beta
// Note: Composite ID (userId_skuId) is constructed as Graph API does not return unique assignment IDs
func (r *UserLicenseAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object UserLicenseAssignmentResourceModel

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

	// Create composite ID: user_id_sku_id
	object.ID = types.StringValue(fmt.Sprintf("%s_%s", object.UserId.ValueString(), object.SkuId.ValueString()))

	// Ensure disabled_plans is set to empty set if not provided (can't be unknown)
	if object.DisabledPlans.IsNull() || object.DisabledPlans.IsUnknown() {
		object.DisabledPlans = types.SetValueMust(types.StringType, []attr.Value{})
	}

	requestBody, err := constructAddLicensesRequest(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing license assignment request",
			fmt.Sprintf("Could not construct license assignment request: %s", err.Error()),
		)
		return
	}

	_, err = r.client.
		Users().
		ByUserId(object.UserId.ValueString()).
		AssignLicense().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully assigned licenses to user: %s", object.UserId.ValueString()))

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

// Read handles the Read operation for user license assignment resources.
//
// Operation: Retrieves user's assigned licenses to verify assignment exists
// API Calls:
//   - GET /users/{id}?$select=id,userPrincipalName,assignedLicenses
//
// Reference: https://learn.microsoft.com/en-us/graph/api/user-get?view=graph-rest-beta
func (r *UserLicenseAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object UserLicenseAssignmentResourceModel

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

	tflog.Debug(ctx, fmt.Sprintf("Reading user license assignments for user ID: %s", object.UserId.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestParameters := &users.UserItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &users.UserItemRequestBuilderGetQueryParameters{
			Select: []string{"id", "userPrincipalName", "assignedLicenses"},
		},
	}

	user, err := r.client.
		Users().
		ByUserId(object.UserId.ValueString()).
		Get(ctx, requestParameters)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &object, user)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for user license assignment resources.
//
// Operation: Updates disabled plans for an existing license assignment
// API Calls:
//   - POST /users/{id}/assignLicense
//
// Reference: https://learn.microsoft.com/en-us/graph/api/user-assignlicense?view=graph-rest-beta
// Note: Only disabled_plans can be updated; sku_id changes trigger resource replacement
func (r *UserLicenseAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan UserLicenseAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update method for: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Update the license (mainly disabled_plans since sku_id has RequiresReplace)
	requestBody, err := constructUpdateLicenseRequest(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing license assignment request for update",
			fmt.Sprintf("Could not construct license assignment request: %s", err.Error()),
		)
		return
	}

	_, err = r.client.
		Users().
		ByUserId(plan.UserId.ValueString()).
		AssignLicense().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully updated license for user: %s", plan.UserId.ValueString()))

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

	tflog.Debug(ctx, fmt.Sprintf("Finished updating %s with ID: %s", ResourceName, plan.ID.ValueString()))
}

// Delete handles the Delete operation for user license assignment resources.
//
// Operation: Removes a license from a user
// API Calls:
//   - POST /users/{id}/assignLicense
//
// Reference: https://learn.microsoft.com/en-us/graph/api/user-assignlicense?view=graph-rest-beta
// Note: License removal is performed by passing the skuId in removeLicenses array
func (r *UserLicenseAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object UserLicenseAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Delete method for: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Remove the single license managed by this resource
	requestBody, err := constructRemoveLicenseRequest(ctx, object.SkuId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing license removal request",
			fmt.Sprintf("Could not construct license removal request: %s", err.Error()),
		)
		return
	}

	_, err = r.client.
		Users().
		ByUserId(object.UserId.ValueString()).
		AssignLicense().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully removed license %s from user: %s", object.SkuId.ValueString(), object.UserId.ValueString()))

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
