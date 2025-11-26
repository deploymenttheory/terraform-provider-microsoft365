package graphBetaDirectorySettings

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Create handles the Create operation.
// Only 1 Group.Unified settings object can exist per tenant. The behavior depends on the overwrite_existing_settings flag:
// - If overwrite_existing_settings = true: Finds the existing tenant settings and overwrites them (PATCH)
// - If overwrite_existing_settings = false (default): Attempts to create new settings (POST)
func (r *DirectorySettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object DirectorySettingsResourceModel

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
			"Error constructing resource for Create method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	templateType := object.TemplateType.ValueString()
	templateID := getTemplateID(templateType)

	if templateID == "" {
		resp.Diagnostics.AddError(
			"Invalid template type",
			fmt.Sprintf("Unknown template type: %s", templateType),
		)
		return
	}

	var shouldPatch bool
	var existingSettingID string

	if !object.OverwriteExistingSettings.IsNull() && object.OverwriteExistingSettings.ValueBool() {
		tflog.Info(ctx, fmt.Sprintf("overwrite_existing_settings is true, checking for existing %s settings", templateType))

		// Check if an instantiated settings object for this template already exists
		existingSettingID, err = r.resolveInstantiatedDirectorySettingsID(ctx, templateID)
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, "Create - Get Existing Settings", r.ReadPermissions)
			return
		}

		if existingSettingID != "" {
			shouldPatch = true
			tflog.Info(ctx, fmt.Sprintf("Found existing %s settings with ID: %s. Will update with Terraform configuration.", templateType, existingSettingID))
		} else {
			tflog.Info(ctx, fmt.Sprintf("No existing %s settings found. Will create new settings object.", templateType))
		}
	} else {
		tflog.Info(ctx, fmt.Sprintf("Creating new %s settings object (overwrite_existing_settings is false)", templateType))
	}

	if shouldPatch {
		// Update existing settings (PATCH)
		_, err := r.client.
			Settings().
			ByDirectorySettingId(existingSettingID).
			Patch(ctx, requestBody, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, "Create - Update Existing Settings", r.WritePermissions)
			return
		}

		object.ID = types.StringValue(existingSettingID)
	} else {
		// Create new settings (POST)
		requestBody.SetTemplateId(&templateID)

		settingObject, err := r.client.
			Settings().
			Post(ctx, requestBody, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, "Create", r.WritePermissions)
			return
		}

		if settingObject.GetId() != nil {
			object.ID = types.StringValue(*settingObject.GetId())
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Wait for eventual consistency before reading back the resource
	time.Sleep(20 * time.Second)

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

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s with ID: %s", ResourceName, object.ID.ValueString()))
}

// Read handles the Read operation.
func (r *DirectorySettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object DirectorySettingsResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	settingId := object.ID.ValueString()

	settingObject, err := r.client.
		Settings().
		ByDirectorySettingId(settingId).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, settingObject)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation.
func (r *DirectorySettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var object DirectorySettingsResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	settingId := object.ID.ValueString()

	requestBody, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Update method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.
		Settings().
		ByDirectorySettingId(settingId).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Update", r.WritePermissions)
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

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s", ResourceName))
}

// Delete handles the Delete operation.
// For directory settings, we delete the settings object from the tenant.
func (r *DirectorySettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object DirectorySettingsResourceModel

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

	settingId := object.ID.ValueString()

	err := r.client.
		Settings().
		ByDirectorySettingId(settingId).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
