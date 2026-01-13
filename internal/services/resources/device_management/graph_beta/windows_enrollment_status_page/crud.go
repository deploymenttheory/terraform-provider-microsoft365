package graphBetaWindowsEnrollmentStatusPage

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
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

// Create handles the Create operation for Windows Enrollment Status Page resources.
func (r *WindowsEnrollmentStatusPageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object WindowsEnrollmentStatusPageResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Create, time.Duration(CreateTimeout)*time.Second, &resp.Diagnostics)
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

	baseResource, err := r.client.
		DeviceManagement().
		DeviceEnrollmentConfigurations().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	object.ID = types.StringValue(*baseResource.GetId())

	tflog.Debug(ctx, fmt.Sprintf("Created %s with ID: %s", ResourceName, object.ID.ValueString()))

	tflog.Debug(ctx, fmt.Sprintf("Assignments detected, constructing assignment request for policy ID: %s", object.ID.ValueString()))

	requestAssignment, err := constructAssignment(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing assignment for Create Method",
			fmt.Sprintf("Could not construct assignment: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	err = r.client.
		DeviceManagement().
		DeviceEnrollmentConfigurations().
		ByDeviceEnrollmentConfigurationId(object.ID.ValueString()).
		Assign().
		Post(ctx, requestAssignment, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully posted assignments for policy ID: %s", object.ID.ValueString()))

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

// Read handles the Read operation for Windows Enrollment Status Page resources.
func (r *WindowsEnrollmentStatusPageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WindowsEnrollmentStatusPageResourceModel

	operation := constants.TfOperationRead
	if ctxOp := ctx.Value("retry_operation"); ctxOp != nil {
		if opStr, ok := ctxOp.(string); ok {
			operation = opStr
		}
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading %s with ID: %s (operation: %s)", ResourceName, state.ID.ValueString(), operation))

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Read, time.Duration(ReadTimeout)*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "resource_id", state.ID.ValueString())
	tflog.Debug(ctx, "Starting read of enrollment status page configuration")

	enrollmentConfig, err := r.client.
		DeviceManagement().
		DeviceEnrollmentConfigurations().
		ByDeviceEnrollmentConfigurationId(state.ID.ValueString()).
		Get(ctx, &devicemanagement.DeviceEnrollmentConfigurationsDeviceEnrollmentConfigurationItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &devicemanagement.DeviceEnrollmentConfigurationsDeviceEnrollmentConfigurationItemRequestBuilderGetQueryParameters{
				Expand: []string{"assignments"},
			},
		})

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, &state, enrollmentConfig)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for Windows Enrollment Status Page resources.
func (r *WindowsEnrollmentStatusPageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan WindowsEnrollmentStatusPageResourceModel
	var state WindowsEnrollmentStatusPageResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting update method for: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, time.Duration(UpdateTimeout)*time.Second, &resp.Diagnostics)
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
		DeviceManagement().
		DeviceEnrollmentConfigurations().
		ByDeviceEnrollmentConfigurationId(plan.ID.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	if !plan.Assignments.IsNull() && !plan.Assignments.IsUnknown() {
		tflog.Debug(ctx, fmt.Sprintf("Assignments detected, constructing assignment request for policy ID: %s", plan.ID.ValueString()))

		assignRequestBody, err := constructAssignment(ctx, &plan)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing assignments",
				fmt.Sprintf("Failed to construct assignments for policy: %s", err.Error()),
			)
			return
		}

		err = r.client.
			DeviceManagement().
			DeviceEnrollmentConfigurations().
			ByDeviceEnrollmentConfigurationId(plan.ID.ValueString()).
			Assign().
			Post(ctx, assignRequestBody, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("Successfully posted assignments for policy ID: %s", state.ID.ValueString()))
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

// Delete handles the Delete operation for Windows Enrollment Status Page resources.
func (r *WindowsEnrollmentStatusPageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state WindowsEnrollmentStatusPageResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Delete, time.Duration(DeleteTimeout)*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	err := r.client.
		DeviceManagement().
		DeviceEnrollmentConfigurations().
		ByDeviceEnrollmentConfigurationId(state.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfTfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
