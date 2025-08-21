package graphBetaAutopatchGroups

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

// Create handles the Create operation for Autopatch Groups resources.
func (r *AutopatchGroupsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object AutopatchGroupsResourceModel

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

	// Convert the request body to JSON bytes
	jsonBytes, err := json.Marshal(requestBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error marshaling request body",
			fmt.Sprintf("Could not marshal request body: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	// Create the HTTP request - use custom endpoint
	url := r.APIEndpoint + r.ResourcePath
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBytes))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating HTTP request",
			fmt.Sprintf("Could not create HTTP request: %s: %s", ResourceName, err.Error()),
		)
		return
	}

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

	tflog.Debug(ctx, fmt.Sprintf("POST request response status: %d %s", httpResp.StatusCode, httpResp.Status))

	// Check for successful creation (200 OK or 201 Created)
	if httpResp.StatusCode != http.StatusOK && httpResp.StatusCode != http.StatusCreated {
		var errorBody map[string]interface{}
		if decodeErr := json.NewDecoder(httpResp.Body).Decode(&errorBody); decodeErr == nil {
			tflog.Debug(ctx, fmt.Sprintf("Error response body: %+v", errorBody))
		}

		resp.Diagnostics.AddError(
			"HTTP Error on Create",
			fmt.Sprintf("Received unexpected response status %d when creating %s", httpResp.StatusCode, ResourceName),
		)
		return
	}

	// Parse the response to get the resource data
	var responseData map[string]interface{}
	if err := json.NewDecoder(httpResp.Body).Decode(&responseData); err != nil {
		resp.Diagnostics.AddError(
			"Error parsing response",
			fmt.Sprintf("Could not parse response: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	// Extract ID from response
	if id, ok := responseData["id"].(string); ok {
		object.ID = types.StringValue(id)
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

// Read handles the Read operation for Autopatch Groups resources.
func (r *AutopatchGroupsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object AutopatchGroupsResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for: %s", ResourceName))

	// Get operation type for logging
	_ = "Read"
	if ctxOp := ctx.Value("retry_operation"); ctxOp != nil {
		if _, ok := ctxOp.(string); ok {
			// operation context available if needed
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

	// Create the HTTP request for GET
	url := r.APIEndpoint + r.ResourcePath + "/" + object.ID.ValueString()
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

	tflog.Debug(ctx, fmt.Sprintf("GET request response status: %d %s", httpResp.StatusCode, httpResp.Status))

	// Handle 404 - resource not found
	if httpResp.StatusCode == http.StatusNotFound {
		tflog.Debug(ctx, fmt.Sprintf("Resource %s not found, removing from state", ResourceName))
		resp.State.RemoveResource(ctx)
		return
	}

	// Check for successful read
	if httpResp.StatusCode != http.StatusOK {
		var errorBody map[string]interface{}
		if decodeErr := json.NewDecoder(httpResp.Body).Decode(&errorBody); decodeErr == nil {
			tflog.Debug(ctx, fmt.Sprintf("Error response body: %+v", errorBody))
		}

		resp.Diagnostics.AddError(
			"HTTP Error on Read",
			fmt.Sprintf("Received unexpected response status %d when reading %s", httpResp.StatusCode, ResourceName),
		)
		return
	}

	// Parse the response
	var responseData map[string]interface{}
	if err := json.NewDecoder(httpResp.Body).Decode(&responseData); err != nil {
		resp.Diagnostics.AddError(
			"Error parsing response",
			fmt.Sprintf("Could not parse response: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, responseData)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for Autopatch Groups resources.
func (r *AutopatchGroupsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AutopatchGroupsResourceModel
	var state AutopatchGroupsResourceModel

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
			"Error constructing resource for update method",
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

	// Create the HTTP request for PATCH
	url := r.APIEndpoint + r.ResourcePath + "/" + state.ID.ValueString()
	httpReq, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewReader(jsonBytes))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating HTTP request",
			fmt.Sprintf("Could not create HTTP request: %s: %s", ResourceName, err.Error()),
		)
		return
	}

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

	tflog.Debug(ctx, fmt.Sprintf("PATCH request response status: %d %s", httpResp.StatusCode, httpResp.Status))

	// Check for successful update (200 OK or 204 No Content)
	if httpResp.StatusCode != http.StatusOK && httpResp.StatusCode != http.StatusNoContent {
		var errorBody map[string]interface{}
		if decodeErr := json.NewDecoder(httpResp.Body).Decode(&errorBody); decodeErr == nil {
			tflog.Debug(ctx, fmt.Sprintf("Error response body: %+v", errorBody))
		}

		resp.Diagnostics.AddError(
			"HTTP Error on Update",
			fmt.Sprintf("Received unexpected response status %d when updating %s", httpResp.StatusCode, ResourceName),
		)
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

// Delete handles the Delete operation for Autopatch Groups resources.
func (r *AutopatchGroupsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object AutopatchGroupsResourceModel

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

	// Create the HTTP request for DELETE
	url := r.APIEndpoint + r.ResourcePath + "/" + object.ID.ValueString()
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

	tflog.Debug(ctx, fmt.Sprintf("DELETE request response status: %d %s", httpResp.StatusCode, httpResp.Status))

	// Check for successful deletion (200 OK, 204 No Content, or 404 Not Found)
	if httpResp.StatusCode != http.StatusOK && 
	   httpResp.StatusCode != http.StatusNoContent && 
	   httpResp.StatusCode != http.StatusNotFound {
		var errorBody map[string]interface{}
		if decodeErr := json.NewDecoder(httpResp.Body).Decode(&errorBody); decodeErr == nil {
			tflog.Debug(ctx, fmt.Sprintf("Error response body: %+v", errorBody))
		}

		resp.Diagnostics.AddError(
			"HTTP Error on Delete",
			fmt.Sprintf("Received unexpected response status %d when deleting %s", httpResp.StatusCode, ResourceName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}