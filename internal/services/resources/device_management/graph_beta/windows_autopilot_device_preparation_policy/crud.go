package graphBetaWindowsAutopilotDevicePreparationPolicy

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
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Create handles the Create operation for Windows Autopilot Device Preparation Policy.
func (r *WindowsAutopilotDevicePreparationPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object WindowsAutopilotDevicePreparationPolicyResourceModel

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

	requestBody, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Create Method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	baseResource, err := r.client.
		DeviceManagement().
		ConfigurationPolicies().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	object.ID = types.StringValue(*baseResource.GetId())

	// set specified device security group as the enrollment time device membership target (Just-In-Time configuration)
	if !object.DeviceSecurityGroup.IsNull() && !object.DeviceSecurityGroup.IsUnknown() {
		deviceSecurityGroupID := object.DeviceSecurityGroup.ValueString()

		// Validate that the security group has the required ownership
		diagnostics := validateSecurityGroupOwnership(ctx, r.client, deviceSecurityGroupID)
		if diagnostics.HasError() {
			resp.Diagnostics.Append(diagnostics...)
			return
		}

		// Create the request body for setting enrollment time device membership target
		requestBody, err := constructJustInTimeAssignmentBody(ctx, deviceSecurityGroupID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing enrollment time device membership target",
				fmt.Sprintf("Could not construct enrollment time device membership target: %s", err.Error()),
			)
			return
		}

		_, err = r.client.
			DeviceManagement().
			ConfigurationPolicies().
			ByDeviceManagementConfigurationPolicyId(object.ID.ValueString()).
			SetEnrollmentTimeDeviceMembershipTarget().
			Post(ctx, requestBody, nil)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
			return
		}

		tflog.Info(ctx, fmt.Sprintf("Successfully assigned device security group %s as enrollment time device membership target", deviceSecurityGroupID))
	}

	if object.Assignments != nil {
		requestAssignment, err := constructAssignment(ctx, object.Assignments)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing assignment for Create Method",
				fmt.Sprintf("Could not construct assignment: %s: %s", ResourceName, err.Error()),
			)
			return
		}

		ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
		if cancel == nil {
			return
		}
		defer cancel()

		_, err = r.client.
			DeviceManagement().
			ConfigurationPolicies().
			ByDeviceManagementConfigurationPolicyId(object.ID.ValueString()).
			Assign().
			Post(ctx, requestAssignment, nil)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
			return
		}
	}

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

// Read handles the Read operation for Windows Autopilot Device Preparation Policy.
func (r *WindowsAutopilotDevicePreparationPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object WindowsAutopilotDevicePreparationPolicyResourceModel
	var baseResource models.DeviceManagementConfigurationPolicyable
	var settingsResponse models.DeviceManagementConfigurationSettingCollectionResponseable
	var assignmentsResponse models.DeviceManagementConfigurationPolicyAssignmentCollectionResponseable

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", ResourceName))

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

	baseResource, err := r.client.
		DeviceManagement().
		ConfigurationPolicies().
		ByDeviceManagementConfigurationPolicyId(object.ID.ValueString()).
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	mapResourceToState(ctx, &object, baseResource)

	settingsResponse, err = r.client.
		DeviceManagement().
		ConfigurationPolicies().
		ByDeviceManagementConfigurationPolicyId(object.ID.ValueString()).
		Settings().
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	err = mapSettingsToState(ctx, &object, settingsResponse)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error mapping settings state",
			fmt.Sprintf("Could not map settings to Terraform state: %s", err.Error()),
		)
		return
	}

	assignmentsResponse, err = r.client.
		DeviceManagement().
		ConfigurationPolicies().
		ByDeviceManagementConfigurationPolicyId(object.ID.ValueString()).
		Assignments().
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	mapAssignmentsToState(ctx, &object, assignmentsResponse)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for Windows Autopilot Device Preparation Policy.
func (r *WindowsAutopilotDevicePreparationPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan WindowsAutopilotDevicePreparationPolicyResourceModel
	var state WindowsAutopilotDevicePreparationPolicyResourceModel

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

	requestBody, err := constructResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Update Method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.
		DeviceManagement().
		ConfigurationPolicies().
		ByDeviceManagementConfigurationPolicyId(state.ID.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
	}

	// If device security group is specified, set it as the enrollment time device membership target (Just-In-Time configuration)
	if !plan.DeviceSecurityGroup.IsNull() && !plan.DeviceSecurityGroup.IsUnknown() {
		deviceSecurityGroupID := plan.DeviceSecurityGroup.ValueString()

		diagnostics := validateSecurityGroupOwnership(ctx, r.client, deviceSecurityGroupID)
		if diagnostics.HasError() {
			resp.Diagnostics.Append(diagnostics...)
			return
		}

		requestBody, err := constructJustInTimeAssignmentBody(ctx, deviceSecurityGroupID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing enrollment time device membership target",
				fmt.Sprintf("Could not construct enrollment time device membership target: %s", err.Error()),
			)
			return
		}

		_, err = r.client.
			DeviceManagement().
			ConfigurationPolicies().
			ByDeviceManagementConfigurationPolicyId(state.ID.ValueString()).
			SetEnrollmentTimeDeviceMembershipTarget().
			Post(ctx, requestBody, nil)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
			return
		}

		tflog.Info(ctx, fmt.Sprintf("Successfully assigned device security group %s as enrollment time device membership target", deviceSecurityGroupID))
	}

	if plan.Assignments != nil {
		requestAssignment, err := constructAssignment(ctx, plan.Assignments)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing assignment for Update Method",
				fmt.Sprintf("Could not construct assignment: %s: %s", ResourceName, err.Error()),
			)
			return
		}

		ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
		if cancel == nil {
			return
		}
		defer cancel()

		_, err = r.client.
			DeviceManagement().
			ConfigurationPolicies().
			ByDeviceManagementConfigurationPolicyId(state.ID.ValueString()).
			Assign().
			Post(ctx, requestAssignment, nil)

		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
			return
		}
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = "Update"
	opts.ResourceTypeName = constants.PROVIDER_NAME + "_" + ResourceName

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

// Delete handles the Delete operation for Windows Autopilot Device Preparation Policy.
func (r *WindowsAutopilotDevicePreparationPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object WindowsAutopilotDevicePreparationPolicyResourceModel

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
		ConfigurationPolicies().
		ByDeviceManagementConfigurationPolicyId(object.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
