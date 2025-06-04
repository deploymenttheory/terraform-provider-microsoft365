package graphBetaGroupOwnerAssignment

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/errors"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Create handles the Create operation.
func (r *GroupOwnerAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object GroupOwnerAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s_%s", r.ProviderTypeName, r.TypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	deadline, _ := ctx.Deadline()
	retryTimeout := time.Until(deadline) - time.Second

	groupId := object.GroupID.ValueString()
	ownerId := object.OwnerID.ValueString()
	ownerObjectType := object.OwnerObjectType.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Adding owner of type %s to group %s", ownerObjectType, groupId))

	requestBody, err := constructResource(ctx, &object, r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Create method",
			fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
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
		errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	// Create composite ID since Microsoft Graph API doesn't return a unique assignment ID
	// Group owner assignments are just relationships, not objects with their own IDs
	// We construct a composite ID from group_id/owner_id to uniquely identify this relationship
	compositeID := fmt.Sprintf("%s/%s", groupId, ownerId)
	object.ID = types.StringValue(compositeID)

	// Set initial state
	object.OwnerType = types.StringValue("Unknown") // Will be updated in the read operation
	object.OwnerDisplayName = types.StringValue("")
	// OwnerObjectType should already be set from the plan

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err = retry.RetryContext(ctx, retryTimeout, func() *retry.RetryError {
		readResp := &resource.ReadResponse{State: resp.State}
		r.Read(ctx, resource.ReadRequest{
			State:        resp.State,
			ProviderMeta: req.ProviderMeta,
		}, readResp)

		if readResp.Diagnostics.HasError() {
			return retry.NonRetryableError(fmt.Errorf("error reading resource state after Create Method: %s", readResp.Diagnostics.Errors()))
		}

		resp.State = readResp.State
		return nil
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for resource creation",
			fmt.Sprintf("Failed to verify resource creation: %s", err),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Read handles the Read operation.
func (r *GroupOwnerAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object GroupOwnerAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read of resource: %s_%s", r.ProviderTypeName, r.TypeName))

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
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
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

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Update handles the Update operation.
func (r *GroupOwnerAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state GroupOwnerAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update of resource: %s_%s", r.ProviderTypeName, r.TypeName))

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
			errors.HandleGraphError(ctx, err, resp, "Update (Remove old owner)", r.WritePermissions)
			return
		}

		requestBody, err := constructResource(ctx, &plan, r.client)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing resource for Update method",
				fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
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
			errors.HandleGraphError(ctx, err, resp, "Update (Add new owner)", r.WritePermissions)
			return
		}

		// Update composite ID
		compositeID := fmt.Sprintf("%s/%s", plan.GroupID.ValueString(), plan.OwnerID.ValueString())
		plan.ID = types.StringValue(compositeID)
	}

	// Read to get current state
	readResp := &resource.ReadResponse{State: resp.State}
	r.Read(ctx, resource.ReadRequest{
		State:        resp.State,
		ProviderMeta: req.ProviderMeta,
	}, readResp)

	if readResp.Diagnostics.HasError() {
		resp.Diagnostics.Append(readResp.Diagnostics...)
		return
	}

	resp.State = readResp.State

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Delete handles the Delete operation.
func (r *GroupOwnerAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object GroupOwnerAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Delete of resource: %s_%s", r.ProviderTypeName, r.TypeName))

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
		errors.HandleGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s_%s", r.ProviderTypeName, r.TypeName))
}
