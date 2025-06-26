package graphBetaConditionalAccessPolicy

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	customrequest "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Create handles the Create operation for Conditional Access Policy resources.
//
//   - Retrieves the planned configuration from the create request
//   - Constructs the resource request body from the plan
//   - Sends POST request to create the conditional access policy
//   - Captures the new resource ID from the response
//   - Sets initial state with planned values
//   - Calls Read operation to fetch the latest state from the API
//   - Updates the final state with the fresh data from the API
//
// The function ensures the conditional access policy is created with all specified
// conditions, grant controls, and session controls properly configured.
func (r *ConditionalAccessPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object ConditionalAccessPolicyResourceModel

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

	// Revert to custom request approach - SDK doesn't have Policies() method
	postConfig := customrequest.PostRequestConfig{
		APIVersion:  customrequest.GraphAPIBeta,
		Endpoint:    r.ResourcePath,
		RequestBody: requestBody,
	}

	createdResource, err := customrequest.PostRequest(
		ctx,
		r.client.GetAdapter(),
		postConfig,
		models.CreateConditionalAccessPolicyDetailFromDiscriminatorValue,
		nil,
	)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	// TODO this model approach is failing as there appears to be no correct model in the sdk. so this all falls apart.
	if createdResource != nil {
		if castedResource, ok := createdResource.(*models.AppliedConditionalAccessPolicy); ok {
			if castedResource.GetId() != nil {
				object.ID = types.StringValue(*castedResource.GetId())
				tflog.Debug(ctx, fmt.Sprintf("Successfully created %s with ID: %s", ResourceName, *castedResource.GetId()))
			} else {
				resp.Diagnostics.AddError(
					"Error extracting resource ID",
					"Created resource ID is nil",
				)
				return
			}
		} else {
			resp.Diagnostics.AddError(
				"Error casting created resource",
				"Could not cast created resource to ConditionalAccessPolicyDetail",
			)
			return
		}
	}

	if object.ID.IsNull() || object.ID.IsUnknown() {
		resp.Diagnostics.AddError(
			"Error extracting resource ID",
			fmt.Sprintf("Could not extract ID from created resource: %s. The API may not return the full resource on creation.", ResourceName),
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

// Read handles the Read operation for Conditional Access Policy resources.
//
//   - Retrieves the current state from the read request
//   - Gets the conditional access policy details from the API
//   - Maps the policy details to Terraform state
//
// The function ensures that all components (conditions, grant controls, session controls)
// are properly read and mapped into the Terraform state, providing a complete view
// of the resource's current configuration on the server.
func (r *ConditionalAccessPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object ConditionalAccessPolicyResourceModel

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

	getConfig := customrequest.GetRequestConfig{
		APIVersion:        customrequest.GraphAPIBeta,
		Endpoint:          r.ResourcePath,
		ResourceIDPattern: "/{id}",
		ResourceID:        object.ID.ValueString(),
	}

	var responseBytes json.RawMessage
	responseBytes, err := customrequest.GetRequestByResourceId(
		ctx,
		r.client.GetAdapter(),
		getConfig,
	)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	var baseResource map[string]interface{}
	err = json.Unmarshal(responseBytes, &baseResource)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error unmarshaling response",
			fmt.Sprintf("Could not unmarshal response: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &object, baseResource)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for Conditional Access Policy resources.
//
//   - Retrieves the planned changes from the update request
//   - Constructs the resource request body from the plan
//   - Sends PATCH request to update the conditional access policy
//   - Sets initial state with planned values
//   - Calls Read operation to fetch the latest state from the API
//   - Updates the final state with the fresh data from the API
//
// The function ensures that the policy is updated with the new configuration
// and the final state reflects the actual state of the resource on the server.
func (r *ConditionalAccessPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var object ConditionalAccessPolicyResourceModel

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

	requestBody, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Update Method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	patchConfig := customrequest.PatchRequestConfig{
		APIVersion:        customrequest.GraphAPIBeta,
		Endpoint:          r.ResourcePath,
		ResourceIDPattern: "/{id}",
		ResourceID:        object.ID.ValueString(),
		RequestBody:       requestBody,
	}

	err = customrequest.PatchRequestByResourceId(
		ctx,
		r.client.GetAdapter(),
		patchConfig,
	)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
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

// Delete handles the Delete operation for Conditional Access Policy resources.
//
//   - Retrieves the current state from the delete request
//   - Validates the state data and timeout configuration
//   - Sends DELETE request to remove the conditional access policy from the API
//   - Cleans up by removing the resource from Terraform state
//
// The function ensures that the policy is completely removed from the
// Microsoft Graph API and cleans up the Terraform state accordingly.
func (r *ConditionalAccessPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object ConditionalAccessPolicyResourceModel

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

	deleteConfig := customrequest.DeleteRequestConfig{
		APIVersion:        customrequest.GraphAPIBeta,
		Endpoint:          r.ResourcePath,
		ResourceIDPattern: "/{id}",
		ResourceID:        object.ID.ValueString(),
	}

	err := customrequest.DeleteRequestByResourceId(
		ctx,
		r.client.GetAdapter(),
		deleteConfig,
	)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
