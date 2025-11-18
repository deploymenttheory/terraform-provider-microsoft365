package graphBetaAuthenticationStrengthPolicy

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

// Create handles the Create operation for Authentication Strength Policy resources.
func (r *AuthenticationStrengthPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object AuthenticationStrengthPolicyResourceModel

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

	// Add required resource headers
	httpReq.Header.Add("x-ms-command-name", "AuthenticationStrengths - AddCustomAuthStrength")

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
	opts.MaxRetries = 60
	opts.RetryInterval = 5 * time.Second

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

// Read handles the Read operation for Authentication Strength Policy resources.
func (r *AuthenticationStrengthPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object AuthenticationStrengthPolicyResourceModel

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

	url := r.httpClient.GetBaseURL() + r.ResourcePath + "/" + object.ID.ValueString() + "?$expand=combinationConfigurations"
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

// Update handles the Update operation for Authentication Strength Policy resources.
// The Graph API requires two separate operations for updating different parts of the resource:
// 1. POST /policies/authenticationStrengthPolicies/{id}/updateAllowedCombinations - for allowedCombinations
// 2. PATCH /identity/conditionalAccess/authenticationStrength/policies/{id}/combinationConfigurations/{configId} - for each configuration
func (r *AuthenticationStrengthPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AuthenticationStrengthPolicyResourceModel
	var state AuthenticationStrengthPolicyResourceModel

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

	// Step 1: Update allowedCombinations if changed
	if !plan.AllowedCombinations.Equal(state.AllowedCombinations) {
		tflog.Debug(ctx, "Allowed combinations changed, updating...")

		// Construct request body
		requestBody, err := constructAllowedCombinationsUpdate(ctx, &plan)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing allowed combinations update",
				fmt.Sprintf("Could not construct update request: %s: %s", ResourceName, err.Error()),
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

		// Make API call
		url := r.httpClient.GetBaseURL() + "/policies/authenticationStrengthPolicies/" + state.ID.ValueString() + "/updateAllowedCombinations"
		httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(jsonBytes))
		if err != nil {
			resp.Diagnostics.AddError(
				"Error creating HTTP request",
				fmt.Sprintf("Could not create HTTP request: %s: %s", ResourceName, err.Error()),
			)
			return
		}

		httpReq.Header.Add("x-ms-command-name", "AuthenticationStrengths - UpdateAllowedCombinations")

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

		tflog.Debug(ctx, fmt.Sprintf("POST updateAllowedCombinations response status: %d %s", httpResp.StatusCode, httpResp.Status))

		if httpResp.StatusCode != http.StatusOK {
			errors.HandleHTTPGraphError(ctx, httpResp, resp, "Update", r.WritePermissions)
			return
		}
	}

	// Step 2: Update combinationConfigurations if changed
	if !plan.CombinationConfigurations.Equal(state.CombinationConfigurations) {
		tflog.Debug(ctx, "Combination configurations changed, updating...")

		var planConfigs []CombinationConfigurationModel
		var stateConfigs []CombinationConfigurationModel

		if !plan.CombinationConfigurations.IsNull() && !plan.CombinationConfigurations.IsUnknown() {
			diags := plan.CombinationConfigurations.ElementsAs(ctx, &planConfigs, false)
			if diags.HasError() {
				resp.Diagnostics.Append(diags...)
				return
			}
		}

		if !state.CombinationConfigurations.IsNull() && !state.CombinationConfigurations.IsUnknown() {
			diags := state.CombinationConfigurations.ElementsAs(ctx, &stateConfigs, false)
			if diags.HasError() {
				resp.Diagnostics.Append(diags...)
				return
			}
		}

		// Match plan configs to state configs by index
		for i, planConfig := range planConfigs {
			// Find the corresponding state config
			var stateConfig *CombinationConfigurationModel
			if i < len(stateConfigs) {
				stateConfig = &stateConfigs[i]
			}

			// Skip if nothing changed for this config
			if stateConfig != nil && combinationConfigurationsEqual(ctx, &planConfig, stateConfig) {
				continue
			}

			// If the state config has an ID, we update it
			if stateConfig == nil || stateConfig.ID.IsNull() || stateConfig.ID.IsUnknown() {
				tflog.Warn(ctx, fmt.Sprintf("Combination configuration at index %d has no ID, may require full resource recreation", i))
				continue
			}

			configID := stateConfig.ID.ValueString()

			// Construct request body
			requestBody, err := constructCombinationConfigurationUpdate(ctx, &planConfig)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error constructing combination configuration update",
					fmt.Sprintf("Could not construct update request for config %s: %s: %s", configID, ResourceName, err.Error()),
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

			// Make API call
			url := r.httpClient.GetBaseURL() + r.ResourcePath + "/" + plan.ID.ValueString() + "/combinationConfigurations/" + configID
			httpReq, err := http.NewRequestWithContext(ctx, "PATCH", url, bytes.NewReader(jsonBytes))
			if err != nil {
				resp.Diagnostics.AddError(
					"Error creating HTTP request",
					fmt.Sprintf("Could not create HTTP request: %s: %s", ResourceName, err.Error()),
				)
				return
			}

			httpReq.Header.Add("x-ms-command-name", "AuthenticationStrengths - UpdateCombinationConfiguration")

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

			tflog.Debug(ctx, fmt.Sprintf("PATCH combinationConfiguration response status: %d %s", httpResp.StatusCode, httpResp.Status))

			if httpResp.StatusCode != http.StatusNoContent && httpResp.StatusCode != http.StatusOK {
				errors.HandleHTTPGraphError(ctx, httpResp, resp, "Update", r.WritePermissions)
				return
			}
		}
	}

	// Step 3: Read back the updated resource
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = "Update"
	opts.ResourceTypeName = constants.PROVIDER_NAME + "_" + ResourceName
	opts.MaxRetries = 60
	opts.RetryInterval = 5 * time.Second

	err := crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after update",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished updating %s with ID: %s", ResourceName, state.ID.ValueString()))
}

// Delete handles the Delete operation for Authentication Strength Policy resources.
func (r *AuthenticationStrengthPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object AuthenticationStrengthPolicyResourceModel

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
