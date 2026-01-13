package graphBetaAutopatchGroups

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

	jsonBytes, err := json.Marshal(requestBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error marshaling request body",
			fmt.Sprintf("Could not marshal request body: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	url := r.APIEndpoint + r.ResourcePath
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBytes))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating HTTP request",
			fmt.Sprintf("Could not create HTTP request: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	httpReq.Header.Set("Content-Type", "application/json")

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

	if httpResp.StatusCode != http.StatusOK && httpResp.StatusCode != http.StatusCreated {
		errors.HandleHTTPGraphError(ctx, httpResp, resp, constants.TfOperationCreate, r.WritePermissions)
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

// Read handles the Read operation for Autopatch Groups resources.
func (r *AutopatchGroupsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object AutopatchGroupsResourceModel

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

	url := r.APIEndpoint + r.ResourcePath + "/" + object.ID.ValueString()
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating HTTP request",
			fmt.Sprintf("Could not create HTTP request: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	httpReq.Header.Set("Content-Type", "application/json")

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

	if prettyJson, err := json.MarshalIndent(baseResource, "", "  "); err == nil {
		tflog.Debug(ctx, fmt.Sprintf("Raw API Response:\n%s", string(prettyJson)))
	}

	// Check if resource is still in "Creating" state and retry until "Active"
	if status, ok := baseResource["status"].(string); ok && status == "Creating" {
		tflog.Debug(ctx, fmt.Sprintf("Resource %s status is 'Creating', waiting for 'Active' state...", ResourceName))

		maxRetries := 30
		retryDelay := 10 * time.Second

		for i := range maxRetries {
			tflog.Debug(ctx, fmt.Sprintf("Retry %d/%d: Waiting %v before checking status again", i+1, maxRetries, retryDelay))

			select {
			case <-ctx.Done():
				resp.Diagnostics.AddError(
					"Context cancelled while waiting for resource to become active",
					fmt.Sprintf("Resource %s is still in 'Creating' state", ResourceName),
				)
				return
			case <-time.After(retryDelay):
			}

			retryReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error creating retry HTTP request",
					fmt.Sprintf("Could not create HTTP request: %s: %s", ResourceName, err.Error()),
				)
				return
			}
			retryReq.Header.Set("Content-Type", "application/json")

			retryResp, err := client.DoWithRetry(ctx, r.httpClient, retryReq, 10)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error making retry HTTP request",
					fmt.Sprintf("Could not make HTTP request: %s: %s", ResourceName, err.Error()),
				)
				return
			}

			if retryResp.StatusCode != http.StatusOK {
				retryResp.Body.Close()
				errors.HandleHTTPGraphError(ctx, retryResp, resp, operation, r.ReadPermissions)
				return
			}

			if err := json.NewDecoder(retryResp.Body).Decode(&baseResource); err != nil {
				retryResp.Body.Close()
				resp.Diagnostics.AddError(
					"Error unmarshaling retry response",
					fmt.Sprintf("Could not unmarshal response: %s: %s", ResourceName, err.Error()),
				)
				return
			}
			retryResp.Body.Close()

			if status, ok := baseResource["status"].(string); ok {
				switch status {
				case "Active":
					tflog.Debug(ctx, fmt.Sprintf("Resource is now 'Active' after %d retries", i+1))
					goto exitRetryLoop
				case "Creating":
					// Continue waiting
				default:
					tflog.Debug(ctx, fmt.Sprintf("Resource status changed to '%s'", status))
					goto exitRetryLoop
				}
			}

			if i == maxRetries-1 {
				resp.Diagnostics.AddWarning(
					"Resource still in 'Creating' state",
					fmt.Sprintf("Resource %s is still in 'Creating' state after %d retries. Proceeding with incomplete data.", ResourceName, maxRetries),
				)
			}
		}
	exitRetryLoop:

		if prettyJson, err := json.MarshalIndent(baseResource, "", "  "); err == nil {
			tflog.Debug(ctx, fmt.Sprintf("Final API Response after status check:\n%s", string(prettyJson)))
		}
	}

	MapRemoteStateToTerraform(ctx, &object, baseResource)

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
			"Error constructing resource for update method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	requestBody["id"] = state.ID.ValueString()

	jsonBytes, err := json.Marshal(requestBody)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error marshaling request body",
			fmt.Sprintf("Could not marshal request body: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	url := r.APIEndpoint + r.ResourcePath
	httpReq, err := http.NewRequestWithContext(ctx, "PUT", url, bytes.NewReader(jsonBytes))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating HTTP request",
			fmt.Sprintf("Could not create HTTP request: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	httpReq.Header.Set("Content-Type", "application/json")

	tflog.Debug(ctx, fmt.Sprintf("Making PUT request to: %s", url))

	httpResp, err := client.DoWithRetry(ctx, r.httpClient, httpReq, 10)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error making HTTP request",
			fmt.Sprintf("Could not make HTTP request: %s: %s", ResourceName, err.Error()),
		)
		return
	}
	defer httpResp.Body.Close()

	tflog.Debug(ctx, fmt.Sprintf("PUT request response status: %d %s", httpResp.StatusCode, httpResp.Status))

	if httpResp.StatusCode != http.StatusNoContent && httpResp.StatusCode != http.StatusOK {
		errors.HandleHTTPGraphError(ctx, httpResp, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

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

	url := r.APIEndpoint + r.ResourcePath + "/" + object.ID.ValueString()
	httpReq, err := http.NewRequestWithContext(ctx, constants.TfOperationDelete, url, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating HTTP request",
			fmt.Sprintf("Could not create HTTP request: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	httpReq.Header.Set("Content-Type", "application/json")

	tflog.Debug(ctx, fmt.Sprintf("Making DELETE request to: %s", url))

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
		errors.HandleHTTPGraphError(ctx, httpResp, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
