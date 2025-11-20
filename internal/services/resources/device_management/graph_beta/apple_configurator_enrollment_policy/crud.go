package graphBetaAppleConfiguratorEnrollmentPolicy

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

// Create handles the Create operation for Apple Configurator Enrollment Policy resources.
func (r *AppleConfiguratorEnrollmentPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object AppleConfiguratorEnrollmentPolicyResourceModel

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

	requestBody, err := constructEnrollmentProfile(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	depId, err := r.resolveDepOnboardingSettingsId(ctx, object.DepOnboardingSettingsId)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Create", r.WritePermissions)
		resp.Diagnostics.AddError("Failed to resolve dep_onboarding_settings_id", err.Error())
		return
	}
	// Persist resolved dep_onboarding_settings_id to state so subsequent operations reuse it
	object.DepOnboardingSettingsId = types.StringValue(depId)
	created, err := r.client.
		DeviceManagement().
		DepOnboardingSettings().
		ByDepOnboardingSettingId(depId).
		EnrollmentProfiles().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	object.ID = types.StringValue(*created.GetId())

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = "Create"
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

// Read handles the Read operation for Apple Configurator Enrollment Policy resources.
func (r *AppleConfiguratorEnrollmentPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object AppleConfiguratorEnrollmentPolicyResourceModel

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

	depId := object.DepOnboardingSettingsId.ValueString()
	if depId == "" {
		var err error
		depId, err = r.resolveDepOnboardingSettingsId(ctx, object.DepOnboardingSettingsId)
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
			resp.Diagnostics.AddError("Failed to resolve dep_onboarding_settings_id", err.Error())
			return
		}
	}

	resource, err := r.client.
		DeviceManagement().
		DepOnboardingSettings().
		ByDepOnboardingSettingId(depId).
		EnrollmentProfiles().
		ByEnrollmentProfileId(object.ID.ValueString()).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, resource, depId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for Apple Configurator Enrollment Policy resources.
func (r *AppleConfiguratorEnrollmentPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AppleConfiguratorEnrollmentPolicyResourceModel
	var state AppleConfiguratorEnrollmentPolicyResourceModel

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

	requestBody, err := constructEnrollmentProfile(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for update method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	depId := state.DepOnboardingSettingsId.ValueString()
	if depId == "" {
		var err error
		depId, err = r.resolveDepOnboardingSettingsId(ctx, state.DepOnboardingSettingsId)
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, "Update", r.WritePermissions)
			resp.Diagnostics.AddError("Failed to resolve dep_onboarding_settings_id", err.Error())
			return
		}
	}

	_, err = r.client.
		DeviceManagement().
		DepOnboardingSettings().
		ByDepOnboardingSettingId(depId).
		EnrollmentProfiles().
		ByEnrollmentProfileId(state.ID.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = "Update"
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

// Delete handles the Delete operation for Apple Configurator Enrollment Policy resources.
func (r *AppleConfiguratorEnrollmentPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object AppleConfiguratorEnrollmentPolicyResourceModel

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

	depId := object.DepOnboardingSettingsId.ValueString()
	if depId == "" {
		var err error
		depId, err = r.resolveDepOnboardingSettingsId(ctx, object.DepOnboardingSettingsId)
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, "Delete", r.WritePermissions)
			resp.Diagnostics.AddError("Failed to resolve dep_onboarding_settings_id", err.Error())
			return
		}
	}

	err := r.client.
		DeviceManagement().
		DepOnboardingSettings().
		ByDepOnboardingSettingId(depId).
		EnrollmentProfiles().
		ByEnrollmentProfileId(object.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
