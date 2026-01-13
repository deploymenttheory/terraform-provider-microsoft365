package graphBetaGroup

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	graphBetaDirectory "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud/graph_beta/directory"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/groups"
)

// Create handles the Create operation.
func (r *GroupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object GroupResourceModel

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

	requestBody, err := constructResource(ctx, &object, false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Create method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	groupObject, err := r.client.
		Groups().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, groupObject)

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

// Read handles the Read operation.
func (r *GroupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object GroupResourceModel

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

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	groupId := object.ID.ValueString()

	// Configure request to include non-default properties via $select
	requestConfig := &groups.GroupItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &groups.GroupItemRequestBuilderGetQueryParameters{
			Select: []string{
				"id",
				"displayName",
				"description",
				"mailNickname",
				"mailEnabled",
				"securityEnabled",
				"groupTypes",
				"visibility",
				"isAssignableToRole",
				"membershipRule",
				"membershipRuleProcessingState",
				"createdDateTime",
			},
		},
	}

	groupObject, err := r.client.
		Groups().
		ByGroupId(groupId).
		Get(ctx, requestConfig)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, groupObject)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation.
func (r *GroupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan GroupResourceModel
	var state GroupResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Updating %s with ID: %s", ResourceName, state.ID.ValueString()))

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

	requestBody, err := constructResource(ctx, &plan, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Update method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	groupId := state.ID.ValueString()

	updatedResource, err := r.client.
		Groups().
		ByGroupId(groupId).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, &plan, updatedResource)

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

// Delete handles the Delete operation.
func (r *GroupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object GroupResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	groupId := object.ID.ValueString()

	// Wait for all licenses to be removed before attempting deletion
	// Groups with active license assignments cannot be deleted
	tflog.Debug(ctx, "Checking for active license assignments before group deletion")
	maxRetries := 30
	retryDelay := 5 * time.Second

	for i := range maxRetries {
		group, err := r.client.Groups().ByGroupId(groupId).Get(ctx, &groups.GroupItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &groups.GroupItemRequestBuilderGetQueryParameters{
				Select: []string{"id", "assignedLicenses"},
			},
		})

		if err != nil {
			tflog.Debug(ctx, "Group not found during license check, proceeding with cleanup")
			break
		}

		assignedLicenses := group.GetAssignedLicenses()
		if len(assignedLicenses) == 0 {
			tflog.Debug(ctx, "No active licenses found, proceeding with group deletion")
			break
		}

		if i == maxRetries-1 {
			resp.Diagnostics.AddError(
				"Delete Failed",
				fmt.Sprintf("Cannot delete group %s: group still has %d active license(s) assigned. Please remove all licenses before deleting the group.", groupId, len(assignedLicenses)),
			)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("Group has %d active license(s), waiting %v before retry (attempt %d/%d)", len(assignedLicenses), retryDelay, i+1, maxRetries))
		time.Sleep(retryDelay)
	}

	// Define soft delete function
	softDeleteFunc := func(ctx context.Context) error {
		return r.client.
			Groups().
			ByGroupId(groupId).
			Delete(ctx, nil)
	}

	deleteOpts := graphBetaDirectory.DeleteOptions{
		ResourceType: graphBetaDirectory.ResourceTypeGroup,
		ResourceID:   groupId,
		ResourceName: object.DisplayName.ValueString(),
	}

	// Execute delete with verification (soft delete + optional hard delete)
	err := graphBetaDirectory.ExecuteDeleteWithVerification(
		ctx,
		r.client,
		softDeleteFunc,
		object.HardDelete.ValueBool(),
		deleteOpts,
	)

	if err != nil {
		resp.Diagnostics.AddError(
			"Delete Failed",
			fmt.Sprintf("Failed to delete %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
