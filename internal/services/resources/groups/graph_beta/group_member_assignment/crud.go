package graphBetaGroupMemberAssignment

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Create handles the Create operation.
func (r *GroupMemberAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object GroupMemberAssignmentResourceModel

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
	memberId := object.MemberID.ValueString()
	memberObjectType := object.MemberObjectType.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Adding member of type %s to group %s", memberObjectType, groupId))

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
		Members().
		Ref().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	// Create composite ID since Microsoft Graph API doesn't return a unique assignment ID
	// Group member assignments are just relationships, not objects with their own IDs
	// We construct a composite ID from group_id/member_id to uniquely identify this relationship
	compositeID := fmt.Sprintf("%s/%s", groupId, memberId)
	object.ID = types.StringValue(compositeID)

	// Set initial state
	object.MemberType = types.StringValue("Unknown") // Will be updated in the read operation
	object.MemberDisplayName = types.StringValue("")
	// MemberObjectType should already be set from the plan

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = "Create"
	opts.ResourceTypeName = constants.PROVIDER_NAME + "_" + ResourceName

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
func (r *GroupMemberAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object GroupMemberAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read of resource: %s", ResourceName))

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
	memberId := object.MemberID.ValueString()

	members, err := r.client.
		Groups().
		ByGroupId(groupId).
		Members().
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	var memberObject graphmodels.DirectoryObjectable
	if members != nil && members.GetValue() != nil {
		for _, member := range members.GetValue() {
			if member.GetId() != nil && *member.GetId() == memberId {
				memberObject = member
				break
			}
		}
	}

	MapRemoteStateToTerraform(ctx, &object, memberObject)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation.
func (r *GroupMemberAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state GroupMemberAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update of resource: %s", ResourceName))

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

	// For group member assignments, if either group_id or member_id changes,
	// we need to remove the old assignment and create a new one
	if plan.GroupID.ValueString() != state.GroupID.ValueString() ||
		plan.MemberID.ValueString() != state.MemberID.ValueString() {

		err := r.client.
			Groups().
			ByGroupId(state.GroupID.ValueString()).
			Members().
			ByDirectoryObjectId(state.MemberID.ValueString()).
			Ref().
			Delete(ctx, nil)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Update (Remove old member)", r.WritePermissions)
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
			Members().
			Ref().
			Post(ctx, requestBody, nil)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Update (Add new member)", r.WritePermissions)
			return
		}

		// Update the composite ID to reflect the new group_id/member_id relationship
		// (since Microsoft Graph API doesn't provide unique assignment IDs)
		compositeID := fmt.Sprintf("%s/%s", plan.GroupID.ValueString(), plan.MemberID.ValueString())
		plan.ID = types.StringValue(compositeID)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = "Update"
	opts.ResourceTypeName = constants.PROVIDER_NAME + "_" + ResourceName

	err := crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after update",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s", ResourceName))
}

// Delete handles the Delete operation.
func (r *GroupMemberAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object GroupMemberAssignmentResourceModel

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

	groupId := object.GroupID.ValueString()
	memberId := object.MemberID.ValueString()

	err := r.client.
		Groups().
		ByGroupId(groupId).
		Members().
		ByDirectoryObjectId(memberId).
		Ref().
		Delete(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
