package graphBetaVisionOSDeviceEnrollmentPolicy

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	customrequest "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
)

// Create handles the Create operation for the visionOS ADE enrollment policy.
func (r *VisionOSDeviceEnrollmentPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object VisionOSDeviceEnrollmentPolicyResourceModel

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

	depOnboardingSettingsId, err := r.resolveDepOnboardingSettingsId(ctx, object.DepOnboardingSettingsId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error resolving dep_onboarding_settings_id",
			err.Error(),
		)
		return
	}
	object.DepOnboardingSettingsId = types.StringValue(depOnboardingSettingsId)

	requestBody, err := constructResource(ctx, &object, depOnboardingSettingsId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Create Method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	baseResource, err := r.client.
		DeviceManagement().
		ConfigurationPolicies().
		Post(ctx, requestBody, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	object.ID = types.StringValue(*baseResource.GetId())

	if !object.IsDefaultPolicyAssignment.IsNull() && !object.IsDefaultPolicyAssignment.IsUnknown() && object.IsDefaultPolicyAssignment.ValueBool() {
		tflog.Debug(ctx, fmt.Sprintf("Setting policy ID %s as default visionOS enrollment profile for dep_onboarding_settings_id %s", object.ID.ValueString(), depOnboardingSettingsId))

		if err := r.setDefaultVisionOSProfileWithRetry(ctx, depOnboardingSettingsId, object.ID.ValueString()); err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
			return
		}
	}

	if !object.DeviceSecurityGroup.IsNull() && !object.DeviceSecurityGroup.IsUnknown() {
		deviceSecurityGroupID := object.DeviceSecurityGroup.ValueString()

		if diagnostics := validateSecurityGroupOwnership(ctx, r.client, deviceSecurityGroupID); diagnostics.HasError() {
			resp.Diagnostics.Append(diagnostics...)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("Calling setEnrollmentTimeDeviceMembershipTarget for policy ID: %s with group ID: %s",
			object.ID.ValueString(), deviceSecurityGroupID))

		// Raw request: the generated SDK method resolves this action to a type-cast URL
		// (".../microsoft.management.services.api.setEnrollmentTimeDeviceMembershipTarget") that
		// the Intune backend rejects with a 500. The Intune admin center itself calls the plain
		// action name posted here, so this bypasses the SDK to match its behavior.
		if err := customrequest.PostRequestNoContent(ctx, r.client.GetAdapter(), customrequest.PostRequestConfig{
			APIVersion:  customrequest.GraphAPIBeta,
			Endpoint:    fmt.Sprintf("deviceManagement/configurationPolicies('%s')/setEnrollmentTimeDeviceMembershipTarget", object.ID.ValueString()),
			RequestBody: constructEnrollmentTimeDeviceMembershipTargetBody(deviceSecurityGroupID),
		}); err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
			return
		}
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

	if err := crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts); err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after create",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s", ResourceName))
}

// Read handles the Read operation for the visionOS ADE enrollment policy.
func (r *VisionOSDeviceEnrollmentPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object VisionOSDeviceEnrollmentPolicyResourceModel
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

	baseResource, err := r.client.
		DeviceManagement().
		ConfigurationPolicies().
		ByDeviceManagementConfigurationPolicyId(object.ID.ValueString()).
		Get(ctx, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	mapResourceToState(ctx, &object, baseResource)

	isDefault, err := r.resolveIsDefaultPolicyAssignment(ctx, object.DepOnboardingSettingsId.ValueString(), object.ID.ValueString())
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}
	object.IsDefaultPolicyAssignment = types.BoolValue(isDefault)

	allSettings, err := r.listAllPolicySettingsWithPageIterator(ctx, object.ID.ValueString())
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	combinedSettingsResponse := models.NewDeviceManagementConfigurationSettingCollectionResponse()
	combinedSettingsResponse.SetValue(allSettings)

	if err := mapSettingsToState(ctx, &object, combinedSettingsResponse); err != nil {
		resp.Diagnostics.AddError(
			"Error mapping settings state",
			fmt.Sprintf("Could not map settings to Terraform state: %s", err.Error()),
		)
		return
	}

	// Raw request: the generated SDK method sends a POST, but the Graph backend only accepts a
	// GET for this action (confirmed against live traffic) and rejects POST with a 400 "no OData
	// route" error - a stale Kiota metadata issue, not specific to this policy or provider.
	membershipTargetBody, err := customrequest.GetRequestByResourceId(ctx, r.client.GetAdapter(), customrequest.GetRequestConfig{
		APIVersion:        customrequest.GraphAPIBeta,
		Endpoint:          "deviceManagement/configurationPolicies",
		ResourceID:        object.ID.ValueString(),
		ResourceIDPattern: "('id')",
		EndpointSuffix:    "/retrieveEnrollmentTimeDeviceMembershipTarget",
	})
	if err != nil {
		tflog.Warn(ctx, fmt.Sprintf("Could not retrieve enrollment time device membership target for policy ID %s; preserving existing device_security_group state: %s", object.ID.ValueString(), err.Error()))
	} else {
		var membershipTargetResult struct {
			EnrollmentTimeDeviceMembershipTargetValidationStatuses []struct {
				TargetId string `json:"targetId"`
			} `json:"enrollmentTimeDeviceMembershipTargetValidationStatuses"`
		}
		if err := json.Unmarshal(membershipTargetBody, &membershipTargetResult); err != nil {
			resp.Diagnostics.AddError(
				"Error parsing enrollment time device membership target response",
				fmt.Sprintf("Could not parse response for policy ID %s: %s", object.ID.ValueString(), err.Error()),
			)
			return
		}

		object.DeviceSecurityGroup = types.StringNull()
		if statuses := membershipTargetResult.EnrollmentTimeDeviceMembershipTargetValidationStatuses; len(statuses) > 0 {
			object.DeviceSecurityGroup = types.StringValue(statuses[0].TargetId)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// listAllPolicySettingsWithPageIterator retrieves every setting for the policy, since the default
// response truncates to the first 25 settings.
func (r *VisionOSDeviceEnrollmentPolicyResource) listAllPolicySettingsWithPageIterator(ctx context.Context, policyId string) ([]models.DeviceManagementConfigurationSettingable, error) {
	var allSettings []models.DeviceManagementConfigurationSettingable

	settingsResponse, err := r.client.
		DeviceManagement().
		ConfigurationPolicies().
		ByDeviceManagementConfigurationPolicyId(policyId).
		Settings().
		Get(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get settings for policy %s: %w", policyId, err)
	}

	pageIterator, err := graphcore.NewPageIterator[models.DeviceManagementConfigurationSettingable](
		settingsResponse,
		r.client.GetAdapter(),
		models.CreateDeviceManagementConfigurationSettingCollectionResponseFromDiscriminatorValue,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create page iterator: %w", err)
	}

	err = pageIterator.Iterate(ctx, func(item models.DeviceManagementConfigurationSettingable) bool {
		if item != nil {
			allSettings = append(allSettings, item)
		}
		return true
	})
	if err != nil {
		return nil, fmt.Errorf("failed to iterate pages: %w", err)
	}

	return allSettings, nil
}

// Update handles the Update operation for the visionOS ADE enrollment policy.
func (r *VisionOSDeviceEnrollmentPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan VisionOSDeviceEnrollmentPolicyResourceModel
	var state VisionOSDeviceEnrollmentPolicyResourceModel

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

	if err := r.validateRequest(ctx, &plan, &state); err != nil {
		resp.Diagnostics.AddError(
			"Validation Error",
			fmt.Sprintf("Pre-request validation failed for resource %s: %s", ResourceName, err.Error()),
		)
		return
	}

	// creationSource must be resent unchanged on every PUT: omitting it causes Graph to reject the
	// request with a 400, and live Intune admin center traffic always includes it on Update too.
	requestBody, err := constructResource(ctx, &plan, plan.DepOnboardingSettingsId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Update Method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	// Use PUT instead of PATCH because the Graph API does not allow PATCH on the 'settings'
	// navigation property of deviceManagementConfigurationPolicy.
	putRequest := customrequest.PutRequestConfig{
		APIVersion:  customrequest.GraphAPIBeta,
		Endpoint:    r.ResourcePath,
		ResourceID:  state.ID.ValueString(),
		RequestBody: requestBody,
	}

	if err := customrequest.PutRequestByResourceId(ctx, r.client.GetAdapter(), putRequest); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	planIsDefault := !plan.IsDefaultPolicyAssignment.IsNull() && !plan.IsDefaultPolicyAssignment.IsUnknown() && plan.IsDefaultPolicyAssignment.ValueBool()
	stateIsDefault := !state.IsDefaultPolicyAssignment.IsNull() && state.IsDefaultPolicyAssignment.ValueBool()

	if planIsDefault && !stateIsDefault {
		tflog.Debug(ctx, fmt.Sprintf("Setting policy ID %s as default visionOS enrollment profile for dep_onboarding_settings_id %s", state.ID.ValueString(), plan.DepOnboardingSettingsId.ValueString()))

		if err := r.setDefaultVisionOSProfileWithRetry(ctx, plan.DepOnboardingSettingsId.ValueString(), state.ID.ValueString()); err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}
	} else if !planIsDefault && stateIsDefault {
		// validateRequest has already confirmed against the live token default that this policy is
		// no longer the default (Graph has no unset action), so there is nothing to call here - the
		// post-update read derives false on its own.
		tflog.Debug(ctx, fmt.Sprintf("Policy ID %s is no longer the default visionOS enrollment profile for dep_onboarding_settings_id %s; no Graph action required for is_default_policy_assignment false",
			state.ID.ValueString(), plan.DepOnboardingSettingsId.ValueString()))
	}

	planDeviceSecurityGroupID := ""
	if !plan.DeviceSecurityGroup.IsNull() && !plan.DeviceSecurityGroup.IsUnknown() {
		planDeviceSecurityGroupID = plan.DeviceSecurityGroup.ValueString()
	}
	stateDeviceSecurityGroupID := ""
	if !state.DeviceSecurityGroup.IsNull() {
		stateDeviceSecurityGroupID = state.DeviceSecurityGroup.ValueString()
	}

	if planDeviceSecurityGroupID != stateDeviceSecurityGroupID {
		// There is no "update" action for the enrollment time device membership target: changing
		// it requires clearing the existing target (DELETE) before setting the new one (POST).
		if stateDeviceSecurityGroupID != "" {
			tflog.Debug(ctx, fmt.Sprintf("Clearing enrollment time device membership target for policy ID: %s", state.ID.ValueString()))

			if err := customrequest.DeleteRequestByResourceId(ctx, r.client.GetAdapter(), customrequest.DeleteRequestConfig{
				APIVersion:        customrequest.GraphAPIBeta,
				Endpoint:          "deviceManagement/configurationPolicies",
				ResourceID:        state.ID.ValueString(),
				ResourceIDPattern: "('id')",
				EndpointSuffix:    "/clearEnrollmentTimeDeviceMembershipTarget",
			}); err != nil {
				errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
				return
			}
		}

		if planDeviceSecurityGroupID != "" {
			if diagnostics := validateSecurityGroupOwnership(ctx, r.client, planDeviceSecurityGroupID); diagnostics.HasError() {
				resp.Diagnostics.Append(diagnostics...)
				return
			}

			tflog.Debug(ctx, fmt.Sprintf("Calling setEnrollmentTimeDeviceMembershipTarget for policy ID: %s with group ID: %s",
				state.ID.ValueString(), planDeviceSecurityGroupID))

			// Raw request: see the matching comment in Create for why the SDK method is bypassed.
			if err := customrequest.PostRequestNoContent(ctx, r.client.GetAdapter(), customrequest.PostRequestConfig{
				APIVersion:  customrequest.GraphAPIBeta,
				Endpoint:    fmt.Sprintf("deviceManagement/configurationPolicies('%s')/setEnrollmentTimeDeviceMembershipTarget", state.ID.ValueString()),
				RequestBody: constructEnrollmentTimeDeviceMembershipTargetBody(planDeviceSecurityGroupID),
			}); err != nil {
				errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
				return
			}
		}
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationUpdate
	opts.ResourceTypeName = ResourceName

	if err := crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts); err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after update",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished updating %s with ID: %s", ResourceName, state.ID.ValueString()))
}

// Delete handles the Delete operation for the visionOS ADE enrollment policy.
func (r *VisionOSDeviceEnrollmentPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object VisionOSDeviceEnrollmentPolicyResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s with ID: %s", ResourceName, object.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	err := r.client.
		DeviceManagement().
		ConfigurationPolicies().
		ByDeviceManagementConfigurationPolicyId(object.ID.ValueString()).
		Delete(ctx, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
