package graphBetaConditionalAccessPolicy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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

	// Convert the request body to JSON bytes
	jsonBytes, err := json.Marshal(requestBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error marshaling request body",
			fmt.Sprintf("Could not marshal request body: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	// Create the HTTP request
	url := r.httpClient.GetBaseURL() + r.ResourcePath
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBytes))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating HTTP request",
			fmt.Sprintf("Could not create HTTP request: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	// Note: Content-Type is automatically set by the AuthenticatedHTTPClient
	tflog.Debug(ctx, fmt.Sprintf("Making POST request to: %s", url))

	// Send the request
	httpResp, err := r.httpClient.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error making HTTP request",
			fmt.Sprintf("Could not make HTTP request: %s: %s", ResourceName, err.Error()),
		)
		return
	}
	defer httpResp.Body.Close()

	// Log the response status for debugging
	tflog.Debug(ctx, fmt.Sprintf("POST request response status: %d %s", httpResp.StatusCode, httpResp.Status))

	// Check for successful creation (201 Created)
	if httpResp.StatusCode != http.StatusCreated {
		// Try to read error response body for better error reporting
		var errorBody map[string]interface{}
		if decodeErr := json.NewDecoder(httpResp.Body).Decode(&errorBody); decodeErr == nil {
			tflog.Debug(ctx, fmt.Sprintf("Error response body: %+v", errorBody))
		}

		resp.Diagnostics.AddError(
			"Error creating resource",
			fmt.Sprintf("Error creating resource: %s: %s", ResourceName, httpResp.Status),
		)
		return
	}

	// Parse the response to get the resource ID
	var createdResource map[string]interface{}
	if err := json.NewDecoder(httpResp.Body).Decode(&createdResource); err != nil {
		resp.Diagnostics.AddError(
			"Error parsing response",
			fmt.Sprintf("Could not parse response: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	// Extract the ID from the response
	id, ok := createdResource["id"].(string)
	if !ok {
		resp.Diagnostics.AddError(
			"Error extracting resource ID",
			"Created resource ID is missing or not a string",
		)
		return
	}

	object.ID = types.StringValue(id)
	tflog.Debug(ctx, fmt.Sprintf("Successfully created %s with ID: %s", ResourceName, id))

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

	// Create the HTTP request
	url := r.httpClient.GetBaseURL() + r.ResourcePath + "/" + object.ID.ValueString()
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating HTTP request",
			fmt.Sprintf("Could not create HTTP request: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Making GET request to: %s", url))

	// Send the request
	httpResp, err := r.httpClient.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error making HTTP request",
			fmt.Sprintf("Could not make HTTP request: %s: %s", ResourceName, err.Error()),
		)
		return
	}
	defer httpResp.Body.Close()

	// Log the response status for debugging
	tflog.Debug(ctx, fmt.Sprintf("GET request response status: %d %s", httpResp.StatusCode, httpResp.Status))

	// Check for successful read (200 OK)
	if httpResp.StatusCode != http.StatusOK {
		// Try to read error response body for better error reporting
		var errorBody map[string]interface{}
		if decodeErr := json.NewDecoder(httpResp.Body).Decode(&errorBody); decodeErr == nil {
			tflog.Debug(ctx, fmt.Sprintf("Error response body: %+v", errorBody))
		}

		// Handle 404 specifically for eventual consistency during retries
		if httpResp.StatusCode == http.StatusNotFound {
			// Check if context is about to expire - if so, remove from state
			if deadline, hasDeadline := ctx.Deadline(); hasDeadline {
				if time.Until(deadline) < 5*time.Second {
					tflog.Warn(ctx, "Resource not found after retries exhausted, removing from state")
					resp.State.RemoveResource(ctx)
					return
				}
			}

			// For 404 during retry attempts, return an error to trigger retry
			resp.Diagnostics.AddError(
				"Resource Not Found",
				fmt.Sprintf("Resource not found (404), this may be due to eventual consistency: %s", httpResp.Status),
			)
			return
		}

		// For other HTTP errors, use the standard error handling
		resp.Diagnostics.AddError(
			"Error reading resource",
			fmt.Sprintf("Error reading resource: %s: %s", ResourceName, httpResp.Status),
		)
		return
	}

	// Parse the response
	var baseResource map[string]interface{}
	if err := json.NewDecoder(httpResp.Body).Decode(&baseResource); err != nil {
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
	var plan ConditionalAccessPolicyResourceModel
	var state ConditionalAccessPolicyResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Updating %s with ID: %s", ResourceName, state.ID.ValueString()))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)   // desired state
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...) // current state (for ID)
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

	// Convert the request body to JSON bytes
	jsonBytes, err := json.Marshal(requestBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error marshaling request body",
			fmt.Sprintf("Could not marshal request body: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	// Create the HTTP request
	url := r.httpClient.GetBaseURL() + r.ResourcePath + "/" + state.ID.ValueString()
	httpReq, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewReader(jsonBytes))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating HTTP request",
			fmt.Sprintf("Could not create HTTP request: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	// Note: Content-Type is automatically set by the AuthenticatedHTTPClient
	tflog.Debug(ctx, fmt.Sprintf("Making PATCH request to: %s", url))

	// Send the request
	httpResp, err := r.httpClient.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error making HTTP request",
			fmt.Sprintf("Could not make HTTP request: %s: %s", ResourceName, err.Error()),
		)
		return
	}
	defer httpResp.Body.Close()

	// Log the response status for debugging
	tflog.Debug(ctx, fmt.Sprintf("PATCH request response status: %d %s", httpResp.StatusCode, httpResp.Status))

	// Check for successful update (204 No Content or 200 OK)
	if httpResp.StatusCode != http.StatusNoContent && httpResp.StatusCode != http.StatusOK {
		// Try to read error response body for better error reporting
		var errorBody map[string]interface{}
		if decodeErr := json.NewDecoder(httpResp.Body).Decode(&errorBody); decodeErr == nil {
			tflog.Debug(ctx, fmt.Sprintf("Error response body: %+v", errorBody))
		}

		resp.Diagnostics.AddError(
			"Error updating resource",
			fmt.Sprintf("Error updating resource: %s: %s", ResourceName, httpResp.Status),
		)
		return
	}

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

	tflog.Debug(ctx, fmt.Sprintf("Finished updating %s with ID: %s", ResourceName, state.ID.ValueString()))
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

	// Create the HTTP request
	url := r.httpClient.GetBaseURL() + r.ResourcePath + "/" + object.ID.ValueString()
	httpReq, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating HTTP request",
			fmt.Sprintf("Could not create HTTP request: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Making DELETE request to: %s", url))

	// Send the request
	httpResp, err := r.httpClient.Do(httpReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error making HTTP request",
			fmt.Sprintf("Could not make HTTP request: %s: %s", ResourceName, err.Error()),
		)
		return
	}
	defer httpResp.Body.Close()

	// Log the response status for debugging
	tflog.Debug(ctx, fmt.Sprintf("DELETE request response status: %d %s", httpResp.StatusCode, httpResp.Status))

	// Check for successful deletion (204 No Content or 404 Not Found)
	if httpResp.StatusCode != http.StatusNoContent && httpResp.StatusCode != http.StatusNotFound {
		// Try to read error response body for better error reporting
		var errorBody map[string]interface{}
		if decodeErr := json.NewDecoder(httpResp.Body).Decode(&errorBody); decodeErr == nil {
			tflog.Debug(ctx, fmt.Sprintf("Error response body: %+v", errorBody))
		}

		resp.Diagnostics.AddError(
			"Error deleting resource",
			fmt.Sprintf("Error deleting resource: %s: %s", ResourceName, httpResp.Status),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
