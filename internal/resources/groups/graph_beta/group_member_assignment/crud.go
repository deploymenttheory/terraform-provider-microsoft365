package graphBetaGroupMemberAssignment

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
func (r *GroupMemberAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object GroupMemberAssignmentResourceModel

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
	memberId := object.MemberID.ValueString()
	memberObjectType := object.MemberObjectType.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Adding member of type %s to group %s", memberObjectType, groupId))

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
func (r *GroupMemberAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object GroupMemberAssignmentResourceModel

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

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Update handles the Update operation.
func (r *GroupMemberAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state GroupMemberAssignmentResourceModel

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
				fmt.Sprintf("Could not construct resource: %s_%s: %s", r.ProviderTypeName, r.TypeName, err.Error()),
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

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s_%s", r.ProviderTypeName, r.TypeName))
}

// Delete handles the Delete operation.
func (r *GroupMemberAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object GroupMemberAssignmentResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s_%s", r.ProviderTypeName, r.TypeName))

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

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s_%s", r.ProviderTypeName, r.TypeName))
}
