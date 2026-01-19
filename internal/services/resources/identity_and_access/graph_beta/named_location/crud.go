package graphBetaNamedLocation

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
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

	createdResource, err := r.client.
		Identity().
		ConditionalAccess().
		NamedLocations().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	if createdResource == nil || createdResource.GetId() == nil {
		resp.Diagnostics.AddError(
			"Error extracting resource ID",
			"Created resource ID is missing from response",
		)
		return
	}

	object.ID = types.StringValue(*createdResource.GetId())
	tflog.Debug(ctx, fmt.Sprintf("Successfully created %s with ID: %s", ResourceName, object.ID.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationCreate
	opts.ResourceTypeName = ResourceName
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

	tflog.Debug(ctx, "Making GET request to retrieve named location")

	remoteResource, err := r.client.
		Identity().
		ConditionalAccess().
		NamedLocations().
		ByNamedLocationId(object.ID.ValueString()).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &object, remoteResource)

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

	tflog.Debug(ctx, "Making PATCH request to update named location")

	_, err = r.client.
		Identity().
		ConditionalAccess().
		NamedLocations().
		ByNamedLocationId(state.ID.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
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

	// Step 1: Get current resource to check if it's a trusted IP location
	tflog.Debug(ctx, "Making initial GET request to check resource before deletion")

	currentResource, err := r.client.
		Identity().
		ConditionalAccess().
		NamedLocations().
		ByNamedLocationId(object.ID.ValueString()).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.ReadPermissions)
		return
	}

	// Check if this is a trusted IP named location
	var needsPatch bool
	if ipLocation, ok := currentResource.(*graphmodels.IpNamedLocation); ok {
		isTrusted := ipLocation.GetIsTrusted()
		needsPatch = isTrusted != nil && *isTrusted
	}

	// Step 2: If it's a trusted IP location, patch it to set isTrusted=false
	if needsPatch {
		tflog.Debug(ctx, "Resource is a trusted IP location, patching to set isTrusted=false")

		patchBody, err := constructResourceForDeletion(ctx)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing deletion patch body",
				fmt.Sprintf("Could not construct deletion patch body: %s: %s", ResourceName, err.Error()),
			)
			return
		}

		tflog.Debug(ctx, "Making PATCH request to set isTrusted=false")

		_, err = r.client.
			Identity().
			ConditionalAccess().
			NamedLocations().
			ByNamedLocationId(object.ID.ValueString()).
			Patch(ctx, patchBody, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
			return
		}

		tflog.Debug(ctx, "Successfully patched isTrusted=false")

		// Step 3: Poll until isTrusted=false is confirmed (eventual consistency)
		maxRetries := 10
		retryDelay := 2 * time.Second

		for i := range maxRetries {
			tflog.Debug(ctx, fmt.Sprintf("Verification attempt %d/%d: checking if isTrusted=false", i+1, maxRetries))

			select {
			case <-time.After(retryDelay):
			case <-ctx.Done():
				resp.Diagnostics.AddError(
					"Context cancelled during isTrusted verification",
					fmt.Sprintf("Context was cancelled while waiting for isTrusted=false: %s", ctx.Err()),
				)
				return
			}

			verifyResource, err := r.client.
				Identity().
				ConditionalAccess().
				NamedLocations().
				ByNamedLocationId(object.ID.ValueString()).
				Get(ctx, nil)

			if err != nil {
				errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.ReadPermissions)
				return
			}

			if ipLocation, ok := verifyResource.(*graphmodels.IpNamedLocation); ok {
				verifyIsTrusted := ipLocation.GetIsTrusted()
				if verifyIsTrusted == nil || !*verifyIsTrusted {
					tflog.Debug(ctx, "Confirmed isTrusted=false, proceeding to delete")
					break
				}
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

		// Wait for eventual consistency of patch operation to complete
		// before attempting deletion
		consistencyDelay := 10 * time.Second
		tflog.Debug(ctx, fmt.Sprintf("Waiting %v for eventual consistency before deletion", consistencyDelay))

		select {
		case <-time.After(consistencyDelay):
		case <-ctx.Done():
			resp.Diagnostics.AddError(
				"Context cancelled during eventual consistency wait",
				fmt.Sprintf("Context was cancelled while waiting for eventual consistency: %s", ctx.Err()),
			)
			return
		}
	}

	// Step 4: Wait for conditional access policies to stop referencing this named location
	// Named locations cannot be deleted while still referenced by CA policies
	maxChecks := 12
	checkInterval := 5 * time.Second
	namedLocationID := object.ID.ValueString()

	for i := range maxChecks {
		tflog.Debug(ctx, fmt.Sprintf("Check %d/%d: Verifying no conditional access policies reference named location %s", i+1, maxChecks, namedLocationID))

		// Query all conditional access policies
		policies, err := r.client.
			Identity().
			ConditionalAccess().
			Policies().
			Get(ctx, nil)

		if err != nil {
			// If we can't check, log but don't fail - attempt deletion anyway
			tflog.Warn(ctx, fmt.Sprintf("Unable to verify CA policy references (check %d/%d), proceeding with deletion", i+1, maxChecks))
			break
		}

		// Check if any policy references this named location
		referencingPolicies := []string{}
		if policies != nil && policies.GetValue() != nil {
			for _, policy := range policies.GetValue() {
				if policy.GetConditions() != nil && policy.GetConditions().GetLocations() != nil {
					locations := policy.GetConditions().GetLocations()

					// Check includeLocations
					if includeLocations := locations.GetIncludeLocations(); includeLocations != nil {
						if slices.Contains(includeLocations, namedLocationID) {
							policyName := "unknown"
							if policy.GetDisplayName() != nil {
								policyName = *policy.GetDisplayName()
							}
							referencingPolicies = append(referencingPolicies, policyName)
						}
					}

					// Check excludeLocations
					if excludeLocations := locations.GetExcludeLocations(); excludeLocations != nil {
						if slices.Contains(excludeLocations, namedLocationID) {
							policyName := "unknown"
							if policy.GetDisplayName() != nil {
								policyName = *policy.GetDisplayName()
							}
							referencingPolicies = append(referencingPolicies, policyName)
						}
					}
				}
			}
		}

		if len(referencingPolicies) == 0 {
			tflog.Debug(ctx, fmt.Sprintf("No conditional access policies reference named location %s, waiting before deletion", namedLocationID))

			// Brief delay after list check passes to prevent executing DELETE too quickly
			// Graph API needs time to settle after CA policy deletions before accepting named location deletion
			validationSyncDelay := 2 * time.Second
			tflog.Debug(ctx, fmt.Sprintf("Waiting %v before DELETE to allow Graph API to settle", validationSyncDelay))

			select {
			case <-time.After(validationSyncDelay):
				tflog.Debug(ctx, "Validation sync delay complete, proceeding with deletion")
			case <-ctx.Done():
				resp.Diagnostics.AddError(
					"Context cancelled during validation sync wait",
					fmt.Sprintf("Context was cancelled while waiting for validation sync: %s", ctx.Err()),
				)
				return
			}
			break
		}

		if i == maxChecks-1 {
			resp.Diagnostics.AddError(
				"Named location still referenced by conditional access policies",
				fmt.Sprintf("After %d checks over %v, the following conditional access policies still reference this named location: %v. "+
					"This may indicate the policies were not properly deleted or eventual consistency is taking longer than expected.",
					maxChecks, time.Duration(maxChecks)*checkInterval, referencingPolicies),
			)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("Named location still referenced by %d CA policies: %v. Waiting %v before retry (check %d/%d)",
			len(referencingPolicies), referencingPolicies, checkInterval, i+1, maxChecks))

		select {
		case <-time.After(checkInterval):
		case <-ctx.Done():
			resp.Diagnostics.AddError(
				"Context cancelled during named location deletion wait",
				fmt.Sprintf("Context was cancelled while waiting for conditional access policies to stop referencing named location: %s", ctx.Err()),
			)
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Making DELETE request for %s with ID: %s", ResourceName, object.ID.ValueString()))
	tflog.Debug(ctx, fmt.Sprintf("DELETE URL will be: /identity/conditionalAccess/namedLocations/%s", object.ID.ValueString()))

	deleteOptions := crud.DefaultDeleteWithRetryOptions()
	deleteOptions.ResourceTypeName = ResourceName
	deleteOptions.ResourceID = object.ID.ValueString()
	deleteOptions.RetryInterval = 10 * time.Second
	deleteOptions.MaxRetries = 6 // 6 * 10s = 60s max retry duration

	tflog.Debug(ctx, fmt.Sprintf("Starting DeleteWithRetry: maxRetries=%d, interval=%v", deleteOptions.MaxRetries, deleteOptions.RetryInterval))

	err = crud.DeleteWithRetry(ctx, func(ctx context.Context) error {
		tflog.Debug(ctx, fmt.Sprintf("Executing DELETE call for named location %s", object.ID.ValueString()))
		deleteErr := r.client.
			Identity().
			ConditionalAccess().
			NamedLocations().
			ByNamedLocationId(object.ID.ValueString()).
			Delete(ctx, nil)

		if deleteErr != nil {
			errorInfo := errors.GraphError(ctx, deleteErr)
			tflog.Debug(ctx, fmt.Sprintf("DELETE call returned error: status=%d, category=%s, message=%s",
				errorInfo.StatusCode, errorInfo.Category, errorInfo.ErrorMessage))
		} else {
			tflog.Debug(ctx, fmt.Sprintf("DELETE call succeeded for named location %s", object.ID.ValueString()))
		}
		return deleteErr
	}, deleteOptions)

	if err != nil {
		errorInfo := errors.GraphError(ctx, err)
		tflog.Error(ctx, fmt.Sprintf("DeleteWithRetry failed: status=%d, category=%s, code=%s, message=%s",
			errorInfo.StatusCode, errorInfo.Category, errorInfo.ErrorCode, errorInfo.ErrorMessage))
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully deleted %s with ID: %s", ResourceName, object.ID.ValueString()))

	// Step 5: Remove from Terraform state
	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
