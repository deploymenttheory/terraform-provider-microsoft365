package graphBetaConditionalAccessPolicy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/generic_client"
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

	requestBody, err := constructResource(ctx, r.httpClient, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Create Method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	jsonBytes, err := json.Marshal(requestBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error marshaling request body",
			fmt.Sprintf("Could not marshal request body: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	url := r.httpClient.GetBaseURL() + r.ResourcePath
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBytes))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating HTTP request",
			fmt.Sprintf("Could not create HTTP request: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Making POST request to: %s", url))

	// Use retry logic with exponential backoff for 429 errors (max 20 retries)
	httpResp, err := client.DoWithRetry(ctx, r.httpClient, httpReq, 10)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error making HTTP request",
			fmt.Sprintf("Could not make HTTP request: %s: %s", ResourceName, err.Error()),
		)
		return
	}
	defer httpResp.Body.Close()

	tflog.Debug(ctx, fmt.Sprintf("POST request response status: %d %s", httpResp.StatusCode, httpResp.Status))

	if httpResp.StatusCode != http.StatusCreated {
		errors.HandleHTTPGraphError(ctx, httpResp, resp, "Create", r.WritePermissions)
		return
	}

	var createdResource map[string]any
	if err := json.NewDecoder(httpResp.Body).Decode(&createdResource); err != nil {
		resp.Diagnostics.AddError(
			"Error parsing response",
			fmt.Sprintf("Could not parse response: %s: %s", ResourceName, err.Error()),
		)
		return
	}

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

	// Use retry logic with exponential backoff for 429 errors (max 20 retries)
	httpResp, err := client.DoWithRetry(ctx, r.httpClient, httpReq, 10)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error making HTTP request",
			fmt.Sprintf("Could not make HTTP request: %s: %s", ResourceName, err.Error()),
		)
		return
	}
	defer httpResp.Body.Close()

	tflog.Debug(ctx, fmt.Sprintf("GET request response status: %d %s", httpResp.StatusCode, httpResp.Status))

	if httpResp.StatusCode != http.StatusOK {
		errors.HandleHTTPGraphError(ctx, httpResp, resp, operation, r.ReadPermissions)
		return
	}

	var baseResource map[string]any
	if err := json.NewDecoder(httpResp.Body).Decode(&baseResource); err != nil {
		resp.Diagnostics.AddError(
			"Error unmarshaling response",
			fmt.Sprintf("Could not unmarshal response: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	// Debug: Pretty print the raw API response
	if prettyJson, err := json.MarshalIndent(baseResource, "", "  "); err == nil {
		tflog.Debug(ctx, fmt.Sprintf("Raw API Response:\n%s", string(prettyJson)))
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

	requestBody, err := constructResource(ctx, r.httpClient, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Update Method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	jsonBytes, err := json.Marshal(requestBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error marshaling request body",
			fmt.Sprintf("Could not marshal request body: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	url := r.httpClient.GetBaseURL() + r.ResourcePath + "/" + state.ID.ValueString()
	httpReq, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewReader(jsonBytes))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating HTTP request",
			fmt.Sprintf("Could not create HTTP request: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Making PATCH request to: %s", url))

	// Use retry logic with exponential backoff for 429 errors (max 20 retries)
	httpResp, err := client.DoWithRetry(ctx, r.httpClient, httpReq, 10)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error making HTTP request",
			fmt.Sprintf("Could not make HTTP request: %s: %s", ResourceName, err.Error()),
		)
		return
	}
	defer httpResp.Body.Close()

	tflog.Debug(ctx, fmt.Sprintf("PATCH request response status: %d %s", httpResp.StatusCode, httpResp.Status))

	if httpResp.StatusCode != http.StatusNoContent && httpResp.StatusCode != http.StatusOK {
		errors.HandleHTTPGraphError(ctx, httpResp, resp, "Update", r.WritePermissions)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Add delay before read to allow for eventual consistency
	// it's possible to perform a valid GET request before the update has propagated,
	// and the response will lack the updated values.
	// Testing has shown a range between around 5-10 seconds is sufficient for the update to propagate and reach eventual consistency.
	// The period varies depending on what fields are updated..
	tflog.Debug(ctx, "Waiting 10 seconds for conditional access policy update to propagate and reach eventual consistency")
	time.Sleep(10 * time.Second)

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

	// Use retry logic with exponential backoff for 429 errors (max 20 retries)
	httpResp, err := client.DoWithRetry(ctx, r.httpClient, httpReq, 10)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error making HTTP request",
			fmt.Sprintf("Could not make HTTP request: %s: %s", ResourceName, err.Error()),
		)
		return
	}
	defer httpResp.Body.Close()

	tflog.Debug(ctx, fmt.Sprintf("DELETE request response status: %d %s", httpResp.StatusCode, httpResp.Status))

	if httpResp.StatusCode != http.StatusNoContent && httpResp.StatusCode != http.StatusNotFound {
		errors.HandleHTTPGraphError(ctx, httpResp, resp, "Delete", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
