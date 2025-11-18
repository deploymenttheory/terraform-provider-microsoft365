package graphBetaNamedLocation

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

// Create handles the Create operation for Named Location resources.
func (r *NamedLocationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object NamedLocationResourceModel

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
	opts.MaxRetries = 60                 // Up from default 30
	opts.RetryInterval = 5 * time.Second // Up from default 2 seconds

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

// Read handles the Read operation for Named Location resources.
func (r *NamedLocationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object NamedLocationResourceModel

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

	MapRemoteResourceStateToTerraform(ctx, &object, baseResource)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for Named Location resources.
func (r *NamedLocationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan NamedLocationResourceModel
	var state NamedLocationResourceModel

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

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = "Update"
	opts.ResourceTypeName = constants.PROVIDER_NAME + "_" + ResourceName
	opts.MaxRetries = 60                 // Up from default 30
	opts.RetryInterval = 5 * time.Second // Up from default 2 seconds

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

// Delete handles the Delete operation for Named Location resources.
//
// This function implements a specialized deletion workflow required by Microsoft Graph's
// Named Location API constraints. The complexity exists because:
//
// 1. TRUSTED IP NAMED LOCATIONS CANNOT BE DELETED DIRECTLY
//   - Microsoft Graph API will reject DELETE requests for IP Named Locations with isTrusted=true
//   - This is a security feature to prevent accidental deletion of trusted network locations
//   - The API requires isTrusted to be explicitly set to false before deletion is allowed
//
// 2. EVENTUAL CONSISTENCY CHALLENGES
//
//   - Microsoft Graph API exhibits eventual consistency behavior
//
//   - A PATCH request to set isTrusted=false may not immediately take effect
//
//   - Subsequent GET requests may still show isTrusted=true for a period of time
//
//   - Attempting DELETE before the change propagates will fail
//
//     3. DELETION WORKFLOW FOR TRUSTED IP LOCATIONS:
//     Step 1: GET resource and check if it's an ipNamedLocation with isTrusted=true
//     Step 2: If conditions met, PATCH to set isTrusted=false
//     Step 3: Poll with GET requests until isTrusted=false is confirmed (eventual consistency)
//     Step 4: Execute DELETE operation
//     Step 5: Remove from Terraform state
//
// 4. DELETION WORKFLOW FOR OTHER NAMED LOCATIONS:
//   - Country Named Locations and non-trusted IP locations can be deleted directly
//   - Skip steps 2-3 and proceed directly to DELETE operation
//
// This approach ensures reliable deletion across all Named Location types while handling
// the API's security constraints and eventual consistency behavior.
func (r *NamedLocationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object NamedLocationResourceModel

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

	getURL := r.httpClient.GetBaseURL() + r.ResourcePath + "/" + object.ID.ValueString()

	var currentResource map[string]any
	var needsPatch bool

	getReq, err := http.NewRequestWithContext(ctx, "GET", getURL, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating GET HTTP request",
			fmt.Sprintf("Could not create GET HTTP request: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Making initial GET request to check resource before deletion: %s", getURL))

	getResp, err := client.DoWithRetry(ctx, r.httpClient, getReq, 10)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error making GET HTTP request",
			fmt.Sprintf("Could not make GET HTTP request: %s: %s", ResourceName, err.Error()),
		)
		return
	}
	defer getResp.Body.Close()

	tflog.Debug(ctx, fmt.Sprintf("GET request response status: %d %s", getResp.StatusCode, getResp.Status))

	if getResp.StatusCode != http.StatusOK {
		if getResp.StatusCode == http.StatusNotFound {
			tflog.Debug(ctx, "Resource not found during pre-deletion check, considering it already deleted")
			resp.State.RemoveResource(ctx)
			return
		}

		errors.HandleHTTPGraphError(ctx, getResp, resp, "Delete", r.ReadPermissions)
		return
	}

	if err := json.NewDecoder(getResp.Body).Decode(&currentResource); err != nil {
		resp.Diagnostics.AddError(
			"Error parsing current resource",
			fmt.Sprintf("Could not parse current resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	// Step 1: Check if this is an IP named location with isTrusted=true
	odataType, _ := currentResource["@odata.type"].(string)
	isTrusted, _ := currentResource["isTrusted"].(bool)
	needsPatch = odataType == "#microsoft.graph.ipNamedLocation" && isTrusted

	// Step 2: If conditions are met, patch with minimum required fields before deletion.
	if needsPatch {
		patchBody, err := constructResourceForDeletion(ctx)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing deletion patch body",
				fmt.Sprintf("Could not construct deletion patch body: %s: %s", ResourceName, err.Error()),
			)
			return
		}

		jsonBytes, err := json.Marshal(patchBody)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error marshaling patch request body",
				fmt.Sprintf("Could not marshal patch request body: %s: %s", ResourceName, err.Error()),
			)
			return
		}

		patchURL := r.httpClient.GetBaseURL() + r.ResourcePath + "/" + object.ID.ValueString()
		patchReq, err := http.NewRequestWithContext(ctx, "PATCH", patchURL, bytes.NewReader(jsonBytes))
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating PATCH HTTP request",
				fmt.Sprintf("Could not create PATCH HTTP request: %s: %s", ResourceName, err.Error()),
			)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("Making PATCH request to set isTrusted=false: %s", patchURL))

		patchResp, err := client.DoWithRetry(ctx, r.httpClient, patchReq, 10)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error making PATCH HTTP request",
				fmt.Sprintf("Could not make PATCH HTTP request: %s: %s", ResourceName, err.Error()),
			)
			return
		}
		defer patchResp.Body.Close()

		tflog.Debug(ctx, fmt.Sprintf("PATCH request response status: %d %s", patchResp.StatusCode, patchResp.Status))

		if patchResp.StatusCode != http.StatusNoContent && patchResp.StatusCode != http.StatusOK {
			errors.HandleHTTPGraphError(ctx, patchResp, resp, "Delete", r.WritePermissions)
			return
		}

		tflog.Debug(ctx, "Successfully patched isTrusted=false")

		maxRetries := 10
		retryDelay := 2 * time.Second

		for i := range maxRetries {
			tflog.Debug(ctx, fmt.Sprintf("Verification attempt %d/%d: checking if isTrusted=false", i+1, maxRetries))

			time.Sleep(retryDelay)

			verifyReq, err := http.NewRequestWithContext(ctx, "GET", getURL, nil)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error creating verification GET HTTP request",
					fmt.Sprintf("Could not create verification GET HTTP request: %s: %s", ResourceName, err.Error()),
				)
				return
			}

			verifyResp, err := client.DoWithRetry(ctx, r.httpClient, verifyReq, 10)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error making verification GET HTTP request",
					fmt.Sprintf("Could not make verification GET HTTP request: %s: %s", ResourceName, err.Error()),
				)
				return
			}

			if verifyResp.StatusCode != http.StatusOK {
				verifyResp.Body.Close()
				if verifyResp.StatusCode == http.StatusNotFound {
					tflog.Debug(ctx, "Resource not found during verification, considering it already deleted")
					resp.State.RemoveResource(ctx)
					return
				}
				errors.HandleHTTPGraphError(ctx, verifyResp, resp, "Delete", r.ReadPermissions)
				return
			}

			var verifyResource map[string]any
			if err := json.NewDecoder(verifyResp.Body).Decode(&verifyResource); err != nil {
				verifyResp.Body.Close()
				resp.Diagnostics.AddError(
					"Error parsing verification resource",
					fmt.Sprintf("Could not parse verification resource: %s: %s", ResourceName, err.Error()),
				)
				return
			}
			verifyResp.Body.Close()

			verifyIsTrusted, _ := verifyResource["isTrusted"].(bool)
			if !verifyIsTrusted {
				tflog.Debug(ctx, "Confirmed isTrusted=false, proceeding to delete")
				break
			}

			if i == maxRetries-1 {
				resp.Diagnostics.AddError(
					"Timeout waiting for isTrusted=false",
					fmt.Sprintf("Timed out waiting for isTrusted to become false after %d attempts", maxRetries),
				)
				return
			}

			tflog.Debug(ctx, fmt.Sprintf("isTrusted still true, retrying in %v", retryDelay))
		}

		// Step 3: Wait for eventual consistency before deletion
		consistencyDelay := 10 * time.Second
		tflog.Debug(ctx, fmt.Sprintf("Waiting %v for eventual consistency before deletion", consistencyDelay))
		time.Sleep(consistencyDelay)
	}

	// Step 4: Execute DELETE operation
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

	// Use retry logic with exponential backoff for 429 errors (max 10 retries)
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

	// Step 5: Remove from Terraform state
	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
