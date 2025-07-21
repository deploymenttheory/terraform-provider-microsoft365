package graphIOSMobileAppConfiguration

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors"
)

func (r *IOSMobileAppConfigurationResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var plan IOSMobileAppConfigurationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(
		ctx,
		plan.Timeouts.Create,
		CreateTimeout*time.Second,
		&resp.Diagnostics,
	)
	if cancel == nil {
		return
	}
	defer cancel()

	// Construct the resource from the plan
	resource := constructResource(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(
		ctx,
		fmt.Sprintf("Creating iOS Mobile App Configuration: %s", plan.DisplayName.ValueString()),
	)

	// Create the resource
	createdResource, err := r.client.DeviceAppManagement().
		MobileAppConfigurations().
		Post(ctx, resource, nil)
	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	// Map the response to state
	iosMobileAppConfig, ok := createdResource.(models.IosMobileAppConfigurationable)
	if !ok {
		resp.Diagnostics.AddError(
			"Type Assertion Failed",
			"Failed to assert created resource as IosMobileAppConfigurationable",
		)
		return
	}

	mapResourceToState(ctx, iosMobileAppConfig, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Handle assignments if provided
	if len(plan.Assignments) > 0 {
		tflog.Debug(ctx, "Processing assignments for iOS Mobile App Configuration")
		assignments := constructAssignments(ctx, &plan, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		// Update with assignments
		updateResource := models.NewIosMobileAppConfiguration()
		updateResource.SetAssignments(assignments)

		configId := plan.Id.ValueString()
		_, err = r.client.DeviceAppManagement().
			MobileAppConfigurations().
			ByManagedDeviceMobileAppConfigurationId(configId).
			Patch(ctx, updateResource, nil)
		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Update assignments", r.WritePermissions)
			return
		}
	}

	// Read back to get the latest state including assignments
	tflog.Debug(ctx, "Reading back iOS Mobile App Configuration after create")
	refreshedResource, err := r.client.DeviceAppManagement().
		MobileAppConfigurations().
		ByManagedDeviceMobileAppConfigurationId(plan.Id.ValueString()).
		Get(ctx, nil)
	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read after create", r.ReadPermissions)
		return
	}

	iosMobileAppConfig, ok = refreshedResource.(models.IosMobileAppConfigurationable)
	if !ok {
		resp.Diagnostics.AddError(
			"Type Assertion Failed",
			"Failed to assert refreshed resource as IosMobileAppConfigurationable",
		)
		return
	}

	mapResourceToState(ctx, iosMobileAppConfig, &plan, &resp.Diagnostics)

	// Read assignments
	assignmentsResp, err := r.client.DeviceAppManagement().
		MobileAppConfigurations().
		ByManagedDeviceMobileAppConfigurationId(plan.Id.ValueString()).
		Assignments().
		Get(ctx, nil)
	if err != nil {
		tflog.Warn(ctx, "Failed to read assignments", map[string]interface{}{"error": err.Error()})
	} else if assignmentsResp != nil && assignmentsResp.GetValue() != nil {
		plan.Assignments = mapAssignmentsToState(ctx, assignmentsResp.GetValue(), &resp.Diagnostics)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Debug(
		ctx,
		fmt.Sprintf("Finished creating iOS Mobile App Configuration: %s", plan.Id.ValueString()),
	)
}

func (r *IOSMobileAppConfigurationResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var state IOSMobileAppConfigurationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(
		ctx,
		state.Timeouts.Read,
		ReadTimeout*time.Second,
		&resp.Diagnostics,
	)
	if cancel == nil {
		return
	}
	defer cancel()

	tflog.Debug(
		ctx,
		fmt.Sprintf("Reading iOS Mobile App Configuration: %s", state.Id.ValueString()),
	)

	// Read the resource
	resource, err := r.client.DeviceAppManagement().
		MobileAppConfigurations().
		ByManagedDeviceMobileAppConfigurationId(state.Id.ValueString()).
		Get(ctx, nil)
	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	// Map to state
	iosMobileAppConfig, ok := resource.(models.IosMobileAppConfigurationable)
	if !ok {
		resp.Diagnostics.AddError(
			"Type Assertion Failed",
			"Failed to assert resource as IosMobileAppConfigurationable",
		)
		return
	}

	mapResourceToState(ctx, iosMobileAppConfig, &state, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read assignments
	assignmentsResp, err := r.client.DeviceAppManagement().
		MobileAppConfigurations().
		ByManagedDeviceMobileAppConfigurationId(state.Id.ValueString()).
		Assignments().
		Get(ctx, nil)
	if err != nil {
		tflog.Warn(ctx, "Failed to read assignments", map[string]interface{}{"error": err.Error()})
	} else if assignmentsResp != nil && assignmentsResp.GetValue() != nil {
		state.Assignments = mapAssignmentsToState(ctx, assignmentsResp.GetValue(), &resp.Diagnostics)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(
		ctx,
		fmt.Sprintf("Finished reading iOS Mobile App Configuration: %s", state.Id.ValueString()),
	)
}

func (r *IOSMobileAppConfigurationResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan IOSMobileAppConfigurationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state IOSMobileAppConfigurationResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(
		ctx,
		plan.Timeouts.Update,
		UpdateTimeout*time.Second,
		&resp.Diagnostics,
	)
	if cancel == nil {
		return
	}
	defer cancel()

	// Construct the resource update
	updateResource := constructResource(ctx, &plan, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(
		ctx,
		fmt.Sprintf("Updating iOS Mobile App Configuration: %s", state.Id.ValueString()),
	)

	// Update the resource
	configId := state.Id.ValueString()
	_, err := r.client.DeviceAppManagement().
		MobileAppConfigurations().
		ByManagedDeviceMobileAppConfigurationId(configId).
		Patch(ctx, updateResource, nil)
	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
	}

	// Handle assignments update
	if !req.Plan.Raw.IsNull() {
		planAssignments := constructAssignments(ctx, &plan, &resp.Diagnostics)
		if resp.Diagnostics.HasError() {
			return
		}

		// Update assignments by setting them on the resource
		assignmentUpdate := models.NewIosMobileAppConfiguration()
		assignmentUpdate.SetAssignments(planAssignments)

		_, err = r.client.DeviceAppManagement().
			MobileAppConfigurations().
			ByManagedDeviceMobileAppConfigurationId(configId).
			Patch(ctx, assignmentUpdate, nil)
		if err != nil {
			errors.HandleGraphError(ctx, err, resp, "Update assignments", r.WritePermissions)
			return
		}
	}

	// Read back the updated resource
	tflog.Debug(ctx, "Reading back iOS Mobile App Configuration after update")
	refreshedResource, err := r.client.DeviceAppManagement().
		MobileAppConfigurations().
		ByManagedDeviceMobileAppConfigurationId(configId).
		Get(ctx, nil)
	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read after update", r.ReadPermissions)
		return
	}

	iosMobileAppConfig, ok := refreshedResource.(models.IosMobileAppConfigurationable)
	if !ok {
		resp.Diagnostics.AddError(
			"Type Assertion Failed",
			"Failed to assert refreshed resource as IosMobileAppConfigurationable",
		)
		return
	}

	mapResourceToState(ctx, iosMobileAppConfig, &plan, &resp.Diagnostics)

	// Read assignments
	assignmentsResp, err := r.client.DeviceAppManagement().
		MobileAppConfigurations().
		ByManagedDeviceMobileAppConfigurationId(configId).
		Assignments().
		Get(ctx, nil)
	if err != nil {
		tflog.Warn(ctx, "Failed to read assignments", map[string]interface{}{"error": err.Error()})
	} else if assignmentsResp != nil && assignmentsResp.GetValue() != nil {
		plan.Assignments = mapAssignmentsToState(ctx, assignmentsResp.GetValue(), &resp.Diagnostics)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	tflog.Debug(
		ctx,
		fmt.Sprintf("Finished updating iOS Mobile App Configuration: %s", state.Id.ValueString()),
	)
}

func (r *IOSMobileAppConfigurationResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var state IOSMobileAppConfigurationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(
		ctx,
		state.Timeouts.Delete,
		DeleteTimeout*time.Second,
		&resp.Diagnostics,
	)
	if cancel == nil {
		return
	}
	defer cancel()

	tflog.Debug(
		ctx,
		fmt.Sprintf("Deleting iOS Mobile App Configuration: %s", state.Id.ValueString()),
	)

	// Delete the resource
	err := r.client.DeviceAppManagement().
		MobileAppConfigurations().
		ByManagedDeviceMobileAppConfigurationId(state.Id.ValueString()).
		Delete(ctx, nil)
	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	tflog.Debug(
		ctx,
		fmt.Sprintf(
			"Successfully deleted iOS Mobile App Configuration: %s",
			state.Id.ValueString(),
		),
	)
}
