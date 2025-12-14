package graphBetaAuthenticationStrengthPolicy

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

	if err := validateRequest(ctx, r.client, object.DisplayName.ValueString(), ""); err != nil {
		resp.Diagnostics.AddError(
			"Validation Error",
			fmt.Sprintf("Authentication strength policy validation failed: %s", err.Error()),
		)
		return
	}

	requestBody, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Create Method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	createdPolicy, err := r.client.
		Identity().
		ConditionalAccess().
		AuthenticationStrength().
		Policies().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	object.ID = types.StringValue(*createdPolicy.GetId())

	tflog.Debug(ctx, fmt.Sprintf("Successfully created %s with ID: %s", ResourceName, object.ID.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = "Create"
	opts.ResourceTypeName = ResourceName
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

	policyId := object.ID.ValueString()

	policy, err := r.client.
		Identity().
		ConditionalAccess().
		AuthenticationStrength().
		Policies().
		ByAuthenticationStrengthPolicyId(policyId).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &object, policy)

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

	if err := validateRequest(ctx, r.client, plan.DisplayName.ValueString(), state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"Validation Error",
			fmt.Sprintf("Authentication strength policy validation failed: %s", err.Error()),
		)
		return
	}

	// Step 1: Update allowedCombinations if changed
	if !plan.AllowedCombinations.Equal(state.AllowedCombinations) {
		tflog.Debug(ctx, "Allowed combinations changed, updating...")

		requestBody, err := constructAllowedCombinationsUpdateForSDK(ctx, &plan)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing allowed combinations update",
				fmt.Sprintf("Could not construct update request: %s: %s", ResourceName, err.Error()),
			)
			return
		}

		_, err = r.client.
			Identity().
			ConditionalAccess().
			AuthenticationStrength().
			Policies().
			ByAuthenticationStrengthPolicyId(state.ID.ValueString()).
			UpdateAllowedCombinations().
			Post(ctx, requestBody, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, "Update", r.WritePermissions)
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

		for i, planConfig := range planConfigs {
			var stateConfig *CombinationConfigurationModel
			if i < len(stateConfigs) {
				stateConfig = &stateConfigs[i]
			}

			if stateConfig != nil && combinationConfigurationsEqual(ctx, &planConfig, stateConfig) {
				continue
			}

			if stateConfig == nil || stateConfig.ID.IsNull() || stateConfig.ID.IsUnknown() {
				tflog.Warn(ctx, fmt.Sprintf("Combination configuration at index %d has no ID, may require full resource recreation", i))
				continue
			}

			configID := stateConfig.ID.ValueString()

			requestBody, err := constructCombinationConfigurationUpdateForSDK(ctx, &planConfig)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error constructing combination configuration update",
					fmt.Sprintf("Could not construct update request for config %s: %s: %s", configID, ResourceName, err.Error()),
				)
				return
			}

			_, err = r.client.
				Identity().
				ConditionalAccess().
				AuthenticationStrength().
				Policies().
				ByAuthenticationStrengthPolicyId(plan.ID.ValueString()).
				CombinationConfigurations().
				ByAuthenticationCombinationConfigurationId(configID).
				Patch(ctx, requestBody, nil)

			if err != nil {
				errors.HandleKiotaGraphError(ctx, err, resp, "Update", r.WritePermissions)
				return
			}
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = "Update"
	opts.ResourceTypeName = ResourceName
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

	policyId := object.ID.ValueString()

	err := r.client.
		Identity().
		ConditionalAccess().
		AuthenticationStrength().
		Policies().
		ByAuthenticationStrengthPolicyId(policyId).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
