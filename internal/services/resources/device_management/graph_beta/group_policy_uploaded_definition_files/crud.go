package graphBetaGroupPolicyUploadedDefinitionFiles

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
		errors.HandleKiotaGraphError(ctx, err, resp, "Create", r.WritePermissions)
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

// Read handles the Read operation for Group Policy Uploaded Definition Files resources.
func (r *GroupPolicyUploadedDefinitionFileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object GroupPolicyUploadedDefinitionFileResourceModel

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
		Delete(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Update - Delete", r.WritePermissions)
		return
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
		errors.HandleKiotaGraphError(ctx, err, resp, "Update - Create", r.WritePermissions)
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

	// For group policy uploaded definition files, we need to call the remove endpoint
	err := r.client.
		DeviceManagement().
		GroupPolicyUploadedDefinitionFiles().
		ByGroupPolicyUploadedDefinitionFileId(object.ID.ValueString()).
		Remove().
		Post(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}

// checkADMXUploadStatus polls and monitors the status of an ADMX file upload operation until completion.
//
// This function handles the asynchronous nature of ADMX file uploads in Intune, which can take time to process.
// It periodically checks the status of the upload operation and waits for it to complete, with three possible outcomes:
//
// 1. "available" - Upload succeeded and the ADMX file is ready for use
// 2. "uploadFailed" - Upload failed, in which case it automatically cleans up by deleting the failed resource
// 3. "uploadInProgress" - Upload is still processing, function will continue polling until timeout
//
// The function also retrieves detailed error information when an upload fails by examining
// the groupPolicyOperations collection, which contains specific error messages that can help
// diagnose issues such as missing dependency files or format problems.
//
// Parameters:
//   - ctx: The context for controlling cancellation and timeout
//   - id: The ID of the group policy uploaded definition file to check
//
// Returns:
//   - string: The final status of the upload operation ("available", "uploadFailed", etc.)
//   - string: Detailed error message if the upload failed, empty string otherwise
//   - error: Any error that occurred during the status check process itself
func (r *GroupPolicyUploadedDefinitionFileResource) checkADMXUploadStatus(ctx context.Context, id string) (string, string, error) {
	maxRetries := 30
	retryInterval := 5 * time.Second

	for i := 0; i < maxRetries; i++ {

		respResource, err := r.client.
			DeviceManagement().
			GroupPolicyUploadedDefinitionFiles().
			ByGroupPolicyUploadedDefinitionFileId(id).
			Get(ctx, &devicemanagement.GroupPolicyUploadedDefinitionFilesGroupPolicyUploadedDefinitionFileItemRequestBuilderGetRequestConfiguration{
				QueryParameters: &devicemanagement.GroupPolicyUploadedDefinitionFilesGroupPolicyUploadedDefinitionFileItemRequestBuilderGetQueryParameters{
					Expand: []string{"groupPolicyOperations"},
				},
			})

		if err != nil {
			return "", "", fmt.Errorf("failed to get upload status: %v", err)
		}

		if status := respResource.GetStatus(); status != nil {
			statusStr := status.String()

			statusDetails := ""
			if operations := respResource.GetGroupPolicyOperations(); len(operations) > 0 {
				// Look for the most recent upload operation
				for _, op := range operations {
					if op.GetOperationType() != nil && op.GetOperationType().String() == "upload" {
						if op.GetStatusDetails() != nil {
							statusDetails = *op.GetStatusDetails()
							break
						}
					}
				}
			}

			switch statusStr {
			case "available":
				// Upload successful
				return statusStr, statusDetails, nil
			case "uploadFailed":
				// Upload failed - delete the resource
				tflog.Debug(ctx, fmt.Sprintf("ADMX upload failed with details: %s", statusDetails))

				deleteErr := r.client.
					DeviceManagement().
					GroupPolicyUploadedDefinitionFiles().
					ByGroupPolicyUploadedDefinitionFileId(id).
					Delete(ctx, nil)

				if deleteErr != nil {
					tflog.Warn(ctx, fmt.Sprintf("Failed to delete failed upload: %s", deleteErr.Error()))
				}

				return statusStr, statusDetails, nil
			case "uploadInProgress":
				// Still in progress, wait and retry
				tflog.Debug(ctx, fmt.Sprintf("ADMX upload in progress (attempt %d/%d), waiting %v before checking again",
					i+1, maxRetries, retryInterval))
				time.Sleep(retryInterval)
				continue
			default:
				// Unknown status
				return statusStr, statusDetails, fmt.Errorf("unknown ADMX upload status: %s", statusStr)
			}
		}

		time.Sleep(retryInterval)
	}

	return "", "", fmt.Errorf("timed out waiting for ADMX upload to complete after %d attempts", maxRetries)
}
