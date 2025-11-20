package graphBetaGroupOwnerAssignment

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Create handles the Create operation.
func (r *GroupOwnerAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object GroupOwnerAssignmentResourceModel

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

	groupId := object.GroupID.ValueString()
	ownerId := object.OwnerID.ValueString()
	ownerObjectType := object.OwnerObjectType.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Adding owner of type %s to group %s", ownerObjectType, groupId))

	requestBody, err := constructResource(ctx, &object, r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Create method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	err = r.client.
		Groups().
		ByGroupId(groupId).
		Owners().
		Ref().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	// Create composite ID since Microsoft Graph API doesn't return a unique assignment ID
	// Group owner assignments are just relationships, not objects with their own IDs
	// We construct a composite ID from group_id/owner_id to uniquely identify this relationship
	compositeID := fmt.Sprintf("%s/%s", groupId, ownerId)
	object.ID = types.StringValue(compositeID)

	object.OwnerType = types.StringValue("Unknown") // Will be updated in the read operation
	object.OwnerDisplayName = types.StringValue("")

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = "Create"
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
func (r *GroupOwnerAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object GroupOwnerAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", ResourceName))

	operation := "Read"
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

	groupId := object.GroupID.ValueString()
	ownerId := object.OwnerID.ValueString()

	owners, err := r.client.
		Groups().
		ByGroupId(groupId).
		Owners().
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	var ownerObject graphmodels.DirectoryObjectable
	if owners != nil && owners.GetValue() != nil {
		for _, owner := range owners.GetValue() {
			if owner.GetId() != nil && *owner.GetId() == ownerId {
				ownerObject = owner
				break
			}
		}
	}

	MapRemoteStateToTerraform(ctx, &object, ownerObject)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation.
func (r *GroupOwnerAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state GroupOwnerAssignmentResourceModel

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

	// For group owner assignments, if either group_id or owner_id changes,
	// we need to remove the old assignment and create a new one
	if plan.GroupID.ValueString() != state.GroupID.ValueString() ||
		plan.OwnerID.ValueString() != state.OwnerID.ValueString() {

		err := r.client.
			Groups().
			ByGroupId(state.GroupID.ValueString()).
			Owners().
			ByDirectoryObjectId(state.OwnerID.ValueString()).
			Ref().
			Delete(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, "Update (Remove old owner)", r.WritePermissions)
			return
		}

		requestBody, err := constructResource(ctx, &plan, r.client)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing resource for Update method",
				fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
			)
			return
		}

		err = r.client.
			Groups().
			ByGroupId(plan.GroupID.ValueString()).
			Owners().
			Ref().
			Post(ctx, requestBody, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, "Update (Add new owner)", r.WritePermissions)
			return
		}

		compositeID := fmt.Sprintf("%s/%s", plan.GroupID.ValueString(), plan.OwnerID.ValueString())
		plan.ID = types.StringValue(compositeID)
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = "Update"
	opts.ResourceTypeName = ResourceName

	err := crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
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
func (r *GroupOwnerAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object GroupOwnerAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Delete of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	groupId := object.GroupID.ValueString()
	ownerId := object.OwnerID.ValueString()

	err := r.client.
		Groups().
		ByGroupId(groupId).
		Owners().
		ByDirectoryObjectId(ownerId).
		Ref().
		Delete(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
