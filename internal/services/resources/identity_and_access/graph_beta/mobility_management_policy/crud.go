package graphBetaIdentityAndAccessMobilityManagementPolicy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	graphpolicies "github.com/microsoftgraph/msgraph-beta-sdk-go/policies"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
)

// Create adopts the service-created mobility management policy and applies the configuration.
//
// Operation: Mobility management policies are created by the service when a mobility management application is
// added to the tenant and have no POST endpoint. Create PATCHes the configured settings, then applies the scope.
// API Calls:
//   - GET /policies/{mobileDeviceManagementPolicies|mobileAppManagementPolicies}/{id} (confirms the policy exists)
//   - PATCH /policies/{mobileDeviceManagementPolicies|mobileAppManagementPolicies}/{id} (settings, only when configured)
//   - PATCH /policies/{mobileDeviceManagementPolicies|mobileAppManagementPolicies}/{id} (appliesTo none or all, sent on its own)
//   - POST /policies/{mobileDeviceManagementPolicies|mobileAppManagementPolicies}/{id}/includedGroups/$ref (once per added group)
//   - DELETE /policies/{mobileDeviceManagementPolicies|mobileAppManagementPolicies}/{id}/includedGroups/{groupId}/$ref (once per removed group)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/mobiledevicemanagementpolicies-update?view=graph-rest-beta
func (r *MobilityManagementPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object MobilityManagementPolicyResourceModel

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

	// The policy is service-created, so confirm it exists before adopting it. A missing policy fails here
	// instead of being retried as eventual consistency by the read-back.
	_, err := r.getPolicy(ctx, object.PolicyType.ValueString(), object.PolicyID.ValueString())
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.ReadPermissions)
		return
	}

	err = r.applyConfiguration(ctx, &object, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	object.ID = types.StringValue(buildResourceID(object.PolicyType.ValueString(), object.PolicyID.ValueString()))

	tflog.Debug(ctx, fmt.Sprintf("Successfully adopted %s with ID: %s", ResourceName, object.ID.ValueString()))

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

// Read handles the Read operation for the mobility management policy resource.
//
// Operation: Retrieves the mobility management policy and its included group IDs
// API Calls:
//   - GET /policies/{mobileDeviceManagementPolicies|mobileAppManagementPolicies}/{id}
//   - GET /policies/{mobileDeviceManagementPolicies|mobileAppManagementPolicies}/{id}/includedGroups?$select=id (paged)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/mobiledevicemanagementpolicies-get?view=graph-rest-beta
func (r *MobilityManagementPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object MobilityManagementPolicyResourceModel
	var identity sharedmodels.ResourceIdentity

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

	identity.ID = object.ID.ValueString()

	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, identity)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	policyType := object.PolicyType.ValueString()
	policyID := object.PolicyID.ValueString()

	remoteResource, err := r.getPolicy(ctx, policyType, policyID)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	includedGroupIDs, err := r.listIncludedGroupIDs(ctx, policyType, policyID)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	object = MapRemoteResourceStateToTerraform(ctx, object, remoteResource, includedGroupIDs)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for the mobility management policy resource.
//
// Operation: PATCHes the settings that changed, then applies the scope and the included group differences
// API Calls:
//   - PATCH /policies/{mobileDeviceManagementPolicies|mobileAppManagementPolicies}/{id} (settings, only when changed)
//   - PATCH /policies/{mobileDeviceManagementPolicies|mobileAppManagementPolicies}/{id} (appliesTo none or all, sent on its own)
//   - POST /policies/{mobileDeviceManagementPolicies|mobileAppManagementPolicies}/{id}/includedGroups/$ref (once per added group)
//   - DELETE /policies/{mobileDeviceManagementPolicies|mobileAppManagementPolicies}/{id}/includedGroups/{groupId}/$ref (once per removed group)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/mobiledevicemanagementpolicies-update?view=graph-rest-beta
func (r *MobilityManagementPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan MobilityManagementPolicyResourceModel
	var state MobilityManagementPolicyResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Updating %s", ResourceName))

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

	err := r.applyConfiguration(ctx, &plan, &state)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
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

// Delete releases Terraform management without changing the tenant mobility management policy.
//
// Operation: Mobility management policies are owned by the service and cannot be deleted. Resetting appliesTo to
// none would silently stop automatic enrollment for every user in scope, so Delete only removes the resource
// from Terraform state and leaves the policy, its scope and its included groups as they are.
// API Calls: none
func (r *MobilityManagementPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s", ResourceName))

	tflog.Info(ctx, fmt.Sprintf("Removing %s from Terraform state; the tenant mobility management policy is left unchanged", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}

// applyConfiguration writes the configured settings, then the scope. appliesTo cannot be PATCHed together with
// other properties, so it is sent on its own. selected is applied by adding group references; none and all are
// PATCHed and any group references left afterwards are removed. Only differing group references are changed.
func (r *MobilityManagementPolicyResource) applyConfiguration(
	ctx context.Context,
	plan *MobilityManagementPolicyResourceModel,
	state *MobilityManagementPolicyResourceModel,
) error {
	policyType := plan.PolicyType.ValueString()
	policyID := plan.PolicyID.ValueString()

	if settings := constructSettingsRequest(ctx, plan, state); settings != nil {
		if err := r.patchPolicy(ctx, policyType, policyID, settings); err != nil {
			return err
		}
	}

	if !isConfigured(plan.AppliesTo) {
		tflog.Debug(ctx, fmt.Sprintf("applies_to is not configured for %s; leaving the scope and included groups unchanged", ResourceName))
		return nil
	}

	appliesTo := plan.AppliesTo.ValueString()
	desiredGroupIDs := make([]string, 0)

	if appliesTo == appliesToSelected {
		if diags := plan.IncludedGroups.ElementsAs(ctx, &desiredGroupIDs, false); diags.HasError() {
			return fmt.Errorf("%w: %s", errReadIncludedGroups, diags.Errors()[0].Detail())
		}
	} else {
		remoteResource, err := r.getPolicy(ctx, policyType, policyID)
		if err != nil {
			return err
		}
		if remoteResource.GetAppliesTo() == nil || remoteResource.GetAppliesTo().String() != appliesTo {
			requestBody, err := constructAppliesToRequest(ctx, policyType, appliesTo)
			if err != nil {
				return err
			}
			if err := r.patchPolicy(ctx, policyType, policyID, requestBody); err != nil {
				return err
			}
		}
	}

	currentGroupIDs, err := r.listIncludedGroupIDs(ctx, policyType, policyID)
	if err != nil {
		return err
	}

	toAdd, toRemove := diffGroupIDs(currentGroupIDs, desiredGroupIDs)

	// Add before removing so a selected scope never passes through an empty group list.
	for _, groupID := range toAdd {
		tflog.Debug(ctx, fmt.Sprintf("Adding group %s to %s %s", groupID, ResourceName, policyID))
		if err := r.addIncludedGroup(ctx, policyType, policyID, groupID); err != nil {
			return err
		}
	}
	for _, groupID := range toRemove {
		tflog.Debug(ctx, fmt.Sprintf("Removing group %s from %s %s", groupID, ResourceName, policyID))
		if err := r.removeIncludedGroup(ctx, policyType, policyID, groupID); err != nil {
			return err
		}
	}

	return nil
}

// getPolicy reads the policy from the collection that matches the policy type
func (r *MobilityManagementPolicyResource) getPolicy(ctx context.Context, policyType, policyID string) (graphmodels.MobilityManagementPolicyable, error) {
	if policyType == policyTypeApp {
		//nolint:wrapcheck // Keep the Kiota error type for shared Graph HTTP status handling.
		return r.client.Policies().MobileAppManagementPolicies().ByMobileAppManagementPolicyId(policyID).Get(ctx, nil)
	}
	//nolint:wrapcheck // Keep the Kiota error type for shared Graph HTTP status handling.
	return r.client.Policies().MobileDeviceManagementPolicies().ByMobileDeviceManagementPolicyId(policyID).Get(ctx, nil)
}

// patchPolicy PATCHes the policy in the collection that matches the policy type
func (r *MobilityManagementPolicyResource) patchPolicy(ctx context.Context, policyType, policyID string, requestBody graphmodels.MobilityManagementPolicyable) error {
	if appBody, ok := requestBody.(graphmodels.MobileAppManagementPolicyable); ok && policyType == policyTypeApp {
		_, err := r.client.Policies().MobileAppManagementPolicies().ByMobileAppManagementPolicyId(policyID).Patch(ctx, appBody, nil)
		//nolint:wrapcheck // Keep the Kiota error type for shared Graph HTTP status handling.
		return err
	}
	if deviceBody, ok := requestBody.(graphmodels.MobileDeviceManagementPolicyable); ok && policyType == policyTypeDevice {
		_, err := r.client.Policies().MobileDeviceManagementPolicies().ByMobileDeviceManagementPolicyId(policyID).Patch(ctx, deviceBody, nil)
		//nolint:wrapcheck // Keep the Kiota error type for shared Graph HTTP status handling.
		return err
	}
	return fmt.Errorf("%w: body %T, policy type %q", errPolicyBodyMismatch, requestBody, policyType)
}

// listIncludedGroupIDs returns the IDs of every group in includedGroups, following @odata.nextLink
func (r *MobilityManagementPolicyResource) listIncludedGroupIDs(ctx context.Context, policyType, policyID string) ([]string, error) {
	groupIDs := make([]string, 0)
	nextLink := ""

	for {
		var page graphmodels.GroupCollectionResponseable
		var err error

		if policyType == policyTypeApp {
			builder := r.client.Policies().MobileAppManagementPolicies().ByMobileAppManagementPolicyId(policyID).IncludedGroups()
			if nextLink != "" {
				builder = builder.WithUrl(nextLink)
			}
			query := &graphpolicies.MobileAppManagementPoliciesItemIncludedGroupsRequestBuilderGetQueryParameters{Select: []string{"id"}}
			page, err = builder.Get(ctx, &graphpolicies.MobileAppManagementPoliciesItemIncludedGroupsRequestBuilderGetRequestConfiguration{QueryParameters: query})
		} else {
			builder := r.client.Policies().MobileDeviceManagementPolicies().ByMobileDeviceManagementPolicyId(policyID).IncludedGroups()
			if nextLink != "" {
				builder = builder.WithUrl(nextLink)
			}
			query := &graphpolicies.MobileDeviceManagementPoliciesItemIncludedGroupsRequestBuilderGetQueryParameters{Select: []string{"id"}}
			page, err = builder.Get(ctx, &graphpolicies.MobileDeviceManagementPoliciesItemIncludedGroupsRequestBuilderGetRequestConfiguration{QueryParameters: query})
		}
		if err != nil {
			//nolint:wrapcheck // Keep the Kiota error type for shared Graph HTTP status handling.
			return nil, err
		}
		if page == nil {
			break
		}

		for _, group := range page.GetValue() {
			if group != nil && group.GetId() != nil {
				groupIDs = append(groupIDs, *group.GetId())
			}
		}

		if page.GetOdataNextLink() == nil || *page.GetOdataNextLink() == "" {
			break
		}
		nextLink = *page.GetOdataNextLink()
	}

	return groupIDs, nil
}

// addIncludedGroup adds one group reference to includedGroups
func (r *MobilityManagementPolicyResource) addIncludedGroup(ctx context.Context, policyType, policyID, groupID string) error {
	requestBody := constructGroupReference(groupID)
	if policyType == policyTypeApp {
		//nolint:wrapcheck // Keep the Kiota error type for shared Graph HTTP status handling.
		return r.client.Policies().MobileAppManagementPolicies().ByMobileAppManagementPolicyId(policyID).IncludedGroups().Ref().Post(ctx, requestBody, nil)
	}
	//nolint:wrapcheck // Keep the Kiota error type for shared Graph HTTP status handling.
	return r.client.Policies().MobileDeviceManagementPolicies().ByMobileDeviceManagementPolicyId(policyID).IncludedGroups().Ref().Post(ctx, requestBody, nil)
}

// removeIncludedGroup removes one group reference from includedGroups
func (r *MobilityManagementPolicyResource) removeIncludedGroup(ctx context.Context, policyType, policyID, groupID string) error {
	if policyType == policyTypeApp {
		//nolint:wrapcheck // Keep the Kiota error type for shared Graph HTTP status handling.
		return r.client.Policies().MobileAppManagementPolicies().ByMobileAppManagementPolicyId(policyID).IncludedGroups().ByGroupId(groupID).Ref().Delete(ctx, nil)
	}
	//nolint:wrapcheck // Keep the Kiota error type for shared Graph HTTP status handling.
	return r.client.Policies().MobileDeviceManagementPolicies().ByMobileDeviceManagementPolicyId(policyID).IncludedGroups().ByGroupId(groupID).Ref().Delete(ctx, nil)
}
