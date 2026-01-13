package graphBetaWindowsAutopilotDeploymentProfile

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Create handles the Create operation.
func (r *WindowsAutopilotDeploymentProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object WindowsAutopilotDeploymentProfileResourceModel

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

	requestBody, err := constructResource(ctx, &object, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	createdResource, err := r.client.
		DeviceManagement().
		WindowsAutopilotDeploymentProfiles().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	object.ID = types.StringValue(*createdResource.GetId())

	// Handle assignments if provided
	if !object.Assignments.IsNull() && !object.Assignments.IsUnknown() {
		assignments, err := constructAssignments(ctx, object.Assignments)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing assignments",
				fmt.Sprintf("Could not construct assignments: %s: %s", ResourceName, err.Error()),
			)
			return
		}

		// Create assignments
		for _, assignment := range assignments {
			_, err := r.client.
				DeviceManagement().
				WindowsAutopilotDeploymentProfiles().
				ByWindowsAutopilotDeploymentProfileId(object.ID.ValueString()).
				Assignments().
				Post(ctx, assignment, nil)

			if err != nil {
				errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
				return
			}
		}
	}

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
func (r *WindowsAutopilotDeploymentProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object WindowsAutopilotDeploymentProfileResourceModel

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

	respResource, err := r.client.
		DeviceManagement().
		WindowsAutopilotDeploymentProfiles().
		ByWindowsAutopilotDeploymentProfileId(object.ID.ValueString()).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &object, respResource)

	// Fetch assignments
	assignmentsResp, err := r.client.
		DeviceManagement().
		WindowsAutopilotDeploymentProfiles().
		ByWindowsAutopilotDeploymentProfileId(object.ID.ValueString()).
		Assignments().
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation+" assignments", r.ReadPermissions)
		return
	}

	// Map assignments to Terraform state
	if assignmentsResp != nil && assignmentsResp.GetValue() != nil {
		assignmentsSet, err := mapRemoteAssignmentsToTerraform(ctx, assignmentsResp.GetValue())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error mapping assignments to Terraform state",
				fmt.Sprintf("Could not map assignments: %s: %s", ResourceName, err.Error()),
			)
			return
		}
		object.Assignments = assignmentsSet
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation.
func (r *WindowsAutopilotDeploymentProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan WindowsAutopilotDeploymentProfileResourceModel
	var state WindowsAutopilotDeploymentProfileResourceModel

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

	requestBody, err := constructResource(ctx, &plan, false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.
		DeviceManagement().
		WindowsAutopilotDeploymentProfiles().
		ByWindowsAutopilotDeploymentProfileId(state.ID.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	// Handle assignment updates
	// First, delete existing assignments
	existingAssignments, err := r.client.
		DeviceManagement().
		WindowsAutopilotDeploymentProfiles().
		ByWindowsAutopilotDeploymentProfileId(state.ID.ValueString()).
		Assignments().
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.ReadPermissions)
		return
	}

	// Delete existing assignments
	if existingAssignments != nil && existingAssignments.GetValue() != nil {
		for _, assignment := range existingAssignments.GetValue() {
			if assignment.GetId() != nil {
				err := r.client.
					DeviceManagement().
					WindowsAutopilotDeploymentProfiles().
					ByWindowsAutopilotDeploymentProfileId(state.ID.ValueString()).
					Assignments().
					ByWindowsAutopilotDeploymentProfileAssignmentId(*assignment.GetId()).
					Delete(ctx, nil)

				if err != nil {
					errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
					return
				}
			}
		}
	}

	// Create new assignments if provided
	if !plan.Assignments.IsNull() && !plan.Assignments.IsUnknown() {
		assignments, err := constructAssignments(ctx, plan.Assignments)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing assignments during update",
				fmt.Sprintf("Could not construct assignments: %s: %s", ResourceName, err.Error()),
			)
			return
		}

		// Create assignments
		for _, assignment := range assignments {
			_, err := r.client.
				DeviceManagement().
				WindowsAutopilotDeploymentProfiles().
				ByWindowsAutopilotDeploymentProfileId(state.ID.ValueString()).
				Assignments().
				Post(ctx, assignment, nil)

			if err != nil {
				errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
				return
			}
		}
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
func (r *WindowsAutopilotDeploymentProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object WindowsAutopilotDeploymentProfileResourceModel

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

	// First, get all assignments and delete them
	assignments, err := r.client.
		DeviceManagement().
		WindowsAutopilotDeploymentProfiles().
		ByWindowsAutopilotDeploymentProfileId(object.ID.ValueString()).
		Assignments().
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfTfOperationDelete, r.ReadPermissions)
		return
	}

	// Delete existing assignments if any exist
	if assignments != nil && assignments.GetValue() != nil && len(assignments.GetValue()) > 0 {
		tflog.Debug(ctx, fmt.Sprintf("Deleting %d assignments for Windows Autopilot Deployment Profile: %s", len(assignments.GetValue()), object.ID.ValueString()))

		for _, assignment := range assignments.GetValue() {
			if assignment.GetId() != nil {
				err := r.client.
					DeviceManagement().
					WindowsAutopilotDeploymentProfiles().
					ByWindowsAutopilotDeploymentProfileId(object.ID.ValueString()).
					Assignments().
					ByWindowsAutopilotDeploymentProfileAssignmentId(*assignment.GetId()).
					Delete(ctx, nil)

				if err != nil {
					errors.HandleKiotaGraphError(ctx, err, resp, constants.TfTfOperationDelete, r.WritePermissions)
					return
				}
			}
		}

		tflog.Debug(ctx, fmt.Sprintf("Successfully deleted all assignments for Windows Autopilot Deployment Profile: %s", object.ID.ValueString()))
	}

	err = r.client.
		DeviceManagement().
		WindowsAutopilotDeploymentProfiles().
		ByWindowsAutopilotDeploymentProfileId(object.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfTfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
