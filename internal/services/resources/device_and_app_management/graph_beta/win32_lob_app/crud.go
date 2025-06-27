package graphBetaWin32LobApp

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
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Create handles the Create operation.
func (r *Win32LobAppResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object Win32LobAppResourceModel

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
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	createdResource, err := r.client.
		DeviceAppManagement().
		MobileApps().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	resourceAsWin32LobApp, ok := createdResource.(models.Win32LobAppable)
	if !ok {
		resp.Diagnostics.AddError(
			"Type Assertion Error",
			fmt.Sprintf("Expected resource of type Win32LobApp for %s, but got %T",
				ResourceName, createdResource),
		)
		return
	}

	object.ID = types.StringValue(*resourceAsWin32LobApp.GetId())

	MapRemoteStateToTerraform(ctx, &object, resourceAsWin32LobApp)

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
func (r *Win32LobAppResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object Win32LobAppResourceModel

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

	fetchedResource, err := r.client.
		DeviceAppManagement().
		MobileApps().
		ByMobileAppId(object.ID.ValueString()).
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	resourceAsWin32LobApp, ok := fetchedResource.(models.Win32LobAppable)
	if !ok {
		resp.Diagnostics.AddError(
			"Type Assertion Error",
			fmt.Sprintf("Expected resource of type Win32LobApp for %s, but got %T",
				ResourceName, fetchedResource),
		)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, resourceAsWin32LobApp)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation.
func (r *Win32LobAppResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan Win32LobAppResourceModel
	var state Win32LobAppResourceModel

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

	requestBody, err := constructResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for update method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	updatedResource, err := r.client.
		DeviceAppManagement().
		MobileApps().
		ByMobileAppId(state.ID.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
	}

	resourceAsWin32LobApp, ok := updatedResource.(models.Win32LobAppable)
	if !ok {
		resp.Diagnostics.AddError(
			"Type Assertion Error",
			fmt.Sprintf("Expected resource of type Win32LobApp for %s, but got %T",
				ResourceName, updatedResource),
		)
		return
	}

	MapRemoteStateToTerraform(ctx, &plan, resourceAsWin32LobApp)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
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

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s", ResourceName))
}

// Delete handles the Delete operation.
func (r *Win32LobAppResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object Win32LobAppResourceModel

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
		DeviceAppManagement().
		MobileApps().
		ByMobileAppId(object.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
