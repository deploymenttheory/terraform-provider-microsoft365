package graphBetaGroupPolicyUploadedDefinitionFiles

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Create handles the Create operation for Group Policy Uploaded Definition Files resources.
func (r *GroupPolicyUploadedDefinitionFileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object GroupPolicyUploadedDefinitionFileResourceModel

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

	if !object.ForceDefinitionFileUpload.IsNull() && object.ForceDefinitionFileUpload.ValueBool() {
		tflog.Info(ctx, "force_definition_file_upload is true, checking for existing definition files with same namespace")
		if err := r.deleteExistingDefinitionFileByContent(ctx, object.Content.ValueString()); err != nil {
			resp.Diagnostics.AddError(
				"Error deleting existing definition file",
				fmt.Sprintf("Could not delete existing definition file with same namespace: %s", err.Error()),
			)
			return
		}
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
		GroupPolicyUploadedDefinitionFiles().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	object.ID = types.StringValue(*baseResource.GetId())

	// Check ADMX upload status before proceeding
	status, statusDetails, err := r.checkADMXUploadStatus(ctx, object.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error checking ADMX upload status",
			fmt.Sprintf("Could not check ADMX upload status: %s", err.Error()),
		)
		return
	}

	if status == "uploadFailed" {
		// Resource was already deleted in the checkADMXUploadStatus function
		errorMessage := "The upload of this ADMX file has failed."
		if statusDetails != "" {
			errorMessage += fmt.Sprintf(" Details: %s", statusDetails)
		}

		resp.Diagnostics.AddError(
			"ADMX upload failed",
			errorMessage,
		)
		return
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

// Read handles the Read operation for Group Policy Uploaded Definition Files resources.
func (r *GroupPolicyUploadedDefinitionFileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object GroupPolicyUploadedDefinitionFileResourceModel
	var identity sharedmodels.ResourceIdentity

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
		GroupPolicyUploadedDefinitionFiles().
		ByGroupPolicyUploadedDefinitionFileId(object.ID.ValueString()).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, respResource)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	identity.ID = object.ID.ValueString()

	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, identity)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for Group Policy Uploaded Definition Files resources.
func (r *GroupPolicyUploadedDefinitionFileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan GroupPolicyUploadedDefinitionFileResourceModel
	var state GroupPolicyUploadedDefinitionFileResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update method for: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Updating %s with ID: %s", ResourceName, state.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Group policy uploaded definition files don't support updates
	// We need to delete the existing one and create a new one
	err := r.client.
		DeviceManagement().
		GroupPolicyUploadedDefinitionFiles().
		ByGroupPolicyUploadedDefinitionFileId(state.ID.ValueString()).
		Remove().
		Post(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	// Wait for removal to complete by polling until the resource returns 404
	err = r.checkADMXRemovalStatus(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for ADMX removal to complete during update",
			fmt.Sprintf("Could not verify ADMX removal completion: %s", err.Error()),
		)
		return
	}

	if !plan.ForceDefinitionFileUpload.IsNull() && plan.ForceDefinitionFileUpload.ValueBool() {
		tflog.Info(ctx, "force_definition_file_upload is true, checking for any other existing definition files with same namespace")
		if err := r.deleteExistingDefinitionFileByContent(ctx, plan.Content.ValueString()); err != nil {
			resp.Diagnostics.AddError(
				"Error deleting existing definition file",
				fmt.Sprintf("Could not delete existing definition file with same namespace: %s", err.Error()),
			)
			return
		}
	}

	requestBody, err := constructResource(ctx, &plan, true)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for update method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	baseResource, err := r.client.
		DeviceManagement().
		GroupPolicyUploadedDefinitionFiles().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	plan.ID = types.StringValue(*baseResource.GetId())

	// Check ADMX upload status before proceeding
	status, statusDetails, err := r.checkADMXUploadStatus(ctx, plan.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error checking ADMX upload status",
			fmt.Sprintf("Could not check ADMX upload status: %s", err.Error()),
		)
		return
	}

	if status == "uploadFailed" {
		// Resource was already deleted in the checkADMXUploadStatus function
		errorMessage := "The upload of this ADMX file has failed."
		if statusDetails != "" {
			errorMessage += fmt.Sprintf(" Details: %s", statusDetails)
		}

		resp.Diagnostics.AddError(
			"ADMX upload failed",
			errorMessage,
		)
		return
	}

	// Update state with the new plan (which has the new ID) before reading
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
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

// Delete handles the Delete operation for Group Policy Uploaded Definition Files resources.
func (r *GroupPolicyUploadedDefinitionFileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object GroupPolicyUploadedDefinitionFileResourceModel

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
		GroupPolicyUploadedDefinitionFiles().
		ByGroupPolicyUploadedDefinitionFileId(object.ID.ValueString()).
		Remove().
		Post(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	err = r.checkADMXRemovalStatus(ctx, object.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error waiting for ADMX removal to complete",
			fmt.Sprintf("Could not verify ADMX removal completion: %s", err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
