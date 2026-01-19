package graphBetaNamedLocation

import (
	"context"
	"fmt"
	"slices"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// handleTrustedIPLocation checks if the resource is a trusted IP location and handles the
// untrusting workflow if needed. Returns false if an error occurred.
func (r *NamedLocationResource) handleTrustedIPLocation(ctx context.Context, resourceID string, resp *resource.DeleteResponse) bool {
	tflog.Debug(ctx, "Making initial GET request to check resource before deletion")

	currentResource, err := r.client.
		Identity().
		ConditionalAccess().
		NamedLocations().
		ByNamedLocationId(resourceID).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.ReadPermissions)
		return false
	}

	needsPatch := r.checkIfTrustedIPLocation(currentResource)
	if !needsPatch {
		return true
	}

	if !r.patchToUntrustIPLocation(ctx, resourceID, resp) {
		return false
	}

	if !r.verifyIPLocationUntrusted(ctx, resourceID, resp) {
		return false
	}

	return r.waitForEventualConsistency(ctx, resp)
}

// checkIfTrustedIPLocation determines if the resource is a trusted IP named location.
func (r *NamedLocationResource) checkIfTrustedIPLocation(resource graphmodels.NamedLocationable) bool {
	if ipLocation, ok := resource.(*graphmodels.IpNamedLocation); ok {
		isTrusted := ipLocation.GetIsTrusted()
		return isTrusted != nil && *isTrusted
	}
	return false
}

// patchToUntrustIPLocation patches a trusted IP location to set isTrusted=false.
func (r *NamedLocationResource) patchToUntrustIPLocation(ctx context.Context, resourceID string, resp *resource.DeleteResponse) bool {
	tflog.Debug(ctx, "Resource is a trusted IP location, patching to set isTrusted=false")

	patchBody, err := constructResourceForDeletion(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing deletion patch body",
			fmt.Sprintf("Could not construct deletion patch body: %s: %s", ResourceName, err.Error()),
		)
		return false
	}

	tflog.Debug(ctx, "Making PATCH request to set isTrusted=false")

	_, err = r.client.
		Identity().
		ConditionalAccess().
		NamedLocations().
		ByNamedLocationId(resourceID).
		Patch(ctx, patchBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return false
	}

	tflog.Debug(ctx, "Successfully patched isTrusted=false")
	return true
}

// verifyIPLocationUntrusted polls until isTrusted=false is confirmed due to eventual consistency.
func (r *NamedLocationResource) verifyIPLocationUntrusted(ctx context.Context, resourceID string, resp *resource.DeleteResponse) bool {
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
			return false
		}

		verifyResource, err := r.client.
			Identity().
			ConditionalAccess().
			NamedLocations().
			ByNamedLocationId(resourceID).
			Get(ctx, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.ReadPermissions)
			return false
		}

		if ipLocation, ok := verifyResource.(*graphmodels.IpNamedLocation); ok {
			verifyIsTrusted := ipLocation.GetIsTrusted()
			if verifyIsTrusted == nil || !*verifyIsTrusted {
				tflog.Debug(ctx, "Confirmed isTrusted=false, proceeding to delete")
				return true
			}
		}

		if i == maxRetries-1 {
			resp.Diagnostics.AddError(
				"Timeout waiting for isTrusted=false",
				fmt.Sprintf("Timed out waiting for isTrusted to become false after %d attempts", maxRetries),
			)
			return false
		}

		tflog.Debug(ctx, fmt.Sprintf("isTrusted still true, retrying in %v", retryDelay))
	}

	return true
}

// waitForEventualConsistency waits for the patch operation to fully propagate before deletion.
func (r *NamedLocationResource) waitForEventualConsistency(ctx context.Context, resp *resource.DeleteResponse) bool {
	consistencyDelay := 10 * time.Second
	tflog.Debug(ctx, fmt.Sprintf("Waiting %v for eventual consistency before deletion", consistencyDelay))

	select {
	case <-time.After(consistencyDelay):
		return true
	case <-ctx.Done():
		resp.Diagnostics.AddError(
			"Context cancelled during eventual consistency wait",
			fmt.Sprintf("Context was cancelled while waiting for eventual consistency: %s", ctx.Err()),
		)
		return false
	}
}

// waitForConditionalAccessPolicyReferences waits until no CA policies reference this named location.
func (r *NamedLocationResource) waitForConditionalAccessPolicyReferences(ctx context.Context, resourceID string, resp *resource.DeleteResponse) bool {
	maxChecks := 12
	checkInterval := 5 * time.Second

	for i := range maxChecks {
		tflog.Debug(ctx, fmt.Sprintf("Check %d/%d: Verifying no conditional access policies reference named location %s", i+1, maxChecks, resourceID))

		policies, err := r.client.
			Identity().
			ConditionalAccess().
			Policies().
			Get(ctx, nil)

		if err != nil {
			tflog.Warn(ctx, fmt.Sprintf("Unable to verify CA policy references (check %d/%d), proceeding with deletion", i+1, maxChecks))
			break
		}

		referencingPolicies := r.findReferencingCAPolicies(policies, resourceID)

		if len(referencingPolicies) == 0 {
			return r.waitForValidationSync(ctx, resourceID, resp)
		}

		if i == maxChecks-1 {
			resp.Diagnostics.AddError(
				"Named location still referenced by conditional access policies",
				fmt.Sprintf("After %d checks over %v, the following conditional access policies still reference this named location: %v. "+
					"This may indicate the policies were not properly deleted or eventual consistency is taking longer than expected.",
					maxChecks, time.Duration(maxChecks)*checkInterval, referencingPolicies),
			)
			return false
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
			return false
		}
	}

	return true
}

// findReferencingCAPolicies returns a list of CA policy names that reference the given named location.
func (r *NamedLocationResource) findReferencingCAPolicies(policies graphmodels.ConditionalAccessPolicyCollectionResponseable, namedLocationID string) []string {
	referencingPolicies := []string{}

	if policies == nil || policies.GetValue() == nil {
		return referencingPolicies
	}

	for _, policy := range policies.GetValue() {
		if policy.GetConditions() == nil || policy.GetConditions().GetLocations() == nil {
			continue
		}

		locations := policy.GetConditions().GetLocations()
		policyName := "unknown"
		if policy.GetDisplayName() != nil {
			policyName = *policy.GetDisplayName()
		}

		if includeLocations := locations.GetIncludeLocations(); includeLocations != nil {
			if slices.Contains(includeLocations, namedLocationID) {
				referencingPolicies = append(referencingPolicies, policyName)
				continue
			}
		}

		if excludeLocations := locations.GetExcludeLocations(); excludeLocations != nil {
			if slices.Contains(excludeLocations, namedLocationID) {
				referencingPolicies = append(referencingPolicies, policyName)
			}
		}
	}

	return referencingPolicies
}

// waitForValidationSync adds a brief delay after CA policy checks pass to allow Graph API to settle.
func (r *NamedLocationResource) waitForValidationSync(ctx context.Context, resourceID string, resp *resource.DeleteResponse) bool {
	tflog.Debug(ctx, fmt.Sprintf("No conditional access policies reference named location %s, waiting before deletion", resourceID))

	validationSyncDelay := 2 * time.Second
	tflog.Debug(ctx, fmt.Sprintf("Waiting %v before DELETE to allow Graph API to settle", validationSyncDelay))

	select {
	case <-time.After(validationSyncDelay):
		tflog.Debug(ctx, "Eventual consistency delay complete, proceeding with deletion")
		return true
	case <-ctx.Done():
		resp.Diagnostics.AddError(
			"Context cancelled during eventual consistency wait",
			fmt.Sprintf("Context was cancelled while waiting for eventual consistency: %s", ctx.Err()),
		)
		return false
	}
}
