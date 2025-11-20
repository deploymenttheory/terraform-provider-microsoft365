package graphBetaWindowsDeviceCompliancePolicy

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

// Create handles the Create operation for Device Compliance Policy resources.
func (r *WindowsDeviceCompliancePolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object DeviceCompliancePolicyResourceModel

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

	// Validate request data, specifically notification_message_cc_list group IDs
	resp.Diagnostics.Append(validateRequest(ctx, r.client, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	requestBody, err := constructResource(ctx, &object, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	baseResource, err := r.client.
		DeviceManagement().
		DeviceCompliancePolicies().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	object.ID = types.StringValue(*baseResource.GetId())

	requestAssignment, err := constructAssignment(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing assignment for Create Method",
			fmt.Sprintf("Could not construct assignment: %s: %s", ResourceName, err.Error()),
		)
		return
	}
	_, err = r.client.
		DeviceManagement().
		DeviceCompliancePolicies().
		ByDeviceCompliancePolicyId(object.ID.ValueString()).
		Assign().
		Post(ctx, requestAssignment, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

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

// Read handles the Read operation for Device Compliance Policy resources.
func (r *WindowsDeviceCompliancePolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object DeviceCompliancePolicyResourceModel

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

	tflog.Debug(ctx, fmt.Sprintf("Reading %s with ID: %s", ResourceName, object.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	respResource, err := r.client.
		DeviceManagement().
		DeviceCompliancePolicies().
		ByDeviceCompliancePolicyId(object.ID.ValueString()).
		Get(ctx, &devicemanagement.DeviceCompliancePoliciesDeviceCompliancePolicyItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &devicemanagement.DeviceCompliancePoliciesDeviceCompliancePolicyItemRequestBuilderGetQueryParameters{
				Expand: []string{"assignments", "scheduledActionsForRule($expand=scheduledActionConfigurations)"},
			},
		})

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, respResource)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for Device Compliance Policy resources.
func (r *WindowsDeviceCompliancePolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan DeviceCompliancePolicyResourceModel
	var state DeviceCompliancePolicyResourceModel

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

	// Validate request data, specifically notification_message_cc_list group IDs
	resp.Diagnostics.Append(validateRequest(ctx, r.client, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	requestBody, err := constructResource(ctx, &plan, false)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for update method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.
		DeviceManagement().
		DeviceCompliancePolicies().
		ByDeviceCompliancePolicyId(state.ID.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
	}

	// Handle scheduleActionsForRules if present
	if !plan.ScheduledActionsForRule.IsNull() && !plan.ScheduledActionsForRule.IsUnknown() {
		var scheduledActionsModels []ScheduledActionForRuleModel
		diags := plan.ScheduledActionsForRule.ElementsAs(ctx, &scheduledActionsModels, false)
		if diags.HasError() {
			resp.Diagnostics.AddError(
				"Error parsing scheduled actions for rules",
				fmt.Sprintf("Could not parse scheduled actions list: %v", diags.Errors()),
			)
			return
		}

		for _, scheduledAction := range scheduledActionsModels {
			scheduleRequestBody, err := constructDeviceComplianceScheduledActionForRulesWithPatchMethod(ctx, scheduledAction)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error constructing scheduled actions for rules",
					fmt.Sprintf("Could not construct scheduled actions request: %s", err.Error()),
				)
				return
			}

			err = r.client.
				DeviceManagement().
				DeviceCompliancePolicies().
				ByDeviceCompliancePolicyId(state.ID.ValueString()).
				ScheduleActionsForRules().
				Post(ctx, scheduleRequestBody, nil)

			if err != nil {
				errors.HandleKiotaGraphError(ctx, err, resp, "Update", r.WritePermissions)
				return
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully scheduled actions for policy ID: %s",
				state.ID.ValueString()))
		}
	}

	// Always handle assignments - either update with new assignments or remove all assignments if nil
	requestAssignment, err := constructAssignment(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing assignment for update method",
			fmt.Sprintf("Could not construct assignment: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.
		DeviceManagement().
		DeviceCompliancePolicies().
		ByDeviceCompliancePolicyId(state.ID.ValueString()).
		Assign().
		Post(ctx, requestAssignment, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Update - Assignments", r.WritePermissions)
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = "Update"
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

// Delete handles the Delete operation for Device Compliance Policy resources.
func (r *WindowsDeviceCompliancePolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object DeviceCompliancePolicyResourceModel

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

	err := r.client.
		DeviceManagement().
		DeviceCompliancePolicies().
		ByDeviceCompliancePolicyId(object.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
