package graphBetaMacOSDepEnrollmentProfile

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
)

// Create handles the Create operation for macOS DEP Enrollment Profile resources.
func (r *MacOSDepEnrollmentProfileResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var object MacOSDepEnrollmentProfileResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// admin_account_password is a write-only attribute: its value lives only in
	// the config, never in the plan/state model. Read it from config and place
	// it on the model so constructResource sends the real password (otherwise it
	// would send null and the DEP-created admin account gets no password).
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("admin_account_password"), &object.AdminAccountPassword)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(
		ctx,
		object.Timeouts.Create,
		CreateTimeout*time.Second,
		&resp.Diagnostics,
	)
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

	depId, err := r.resolveDepOnboardingSettingsId(ctx, object.DepOnboardingSettingsID)
	if err != nil {
		errors.HandleKiotaGraphError(
			ctx,
			err,
			resp,
			constants.TfOperationCreate,
			r.WritePermissions,
		)
		resp.Diagnostics.AddError("Failed to resolve dep_onboarding_settings_id", err.Error())
		return
	}
	// Persist resolved dep_onboarding_settings_id to state so subsequent operations reuse it
	object.DepOnboardingSettingsID = types.StringValue(depId)

	created, err := r.client.
		DeviceManagement().
		DepOnboardingSettings().
		ByDepOnboardingSettingId(depId).
		EnrollmentProfiles().
		Post(ctx, requestBody, nil)
	if err != nil {
		errors.HandleKiotaGraphError(
			ctx,
			err,
			resp,
			constants.TfOperationCreate,
			r.WritePermissions,
		)
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

// Read handles the Read operation for macOS DEP Enrollment Profile resources.
func (r *MacOSDepEnrollmentProfileResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var object MacOSDepEnrollmentProfileResourceModel
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

	ctx, cancel := crud.HandleTimeout(
		ctx,
		object.Timeouts.Read,
		ReadTimeout*time.Second,
		&resp.Diagnostics,
	)
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

	depId := object.DepOnboardingSettingsID.ValueString()
	if depId == "" {
		var err error
		depId, err = r.resolveDepOnboardingSettingsId(ctx, object.DepOnboardingSettingsID)
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
			resp.Diagnostics.AddError("Failed to resolve dep_onboarding_settings_id", err.Error())
			return
		}
	}

	remote, err := r.client.
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

	profile, ok := remote.(graphmodels.DepMacOSEnrollmentProfileable)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected enrollment profile type",
			fmt.Sprintf(
				"Resource %s with ID %s is not a depMacOSEnrollmentProfile",
				ResourceName,
				object.ID.ValueString(),
			),
		)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, profile, depId)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for macOS DEP Enrollment Profile resources.
func (r *MacOSDepEnrollmentProfileResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan MacOSDepEnrollmentProfileResourceModel
	var state MacOSDepEnrollmentProfileResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the write-only admin_account_password from config (it is null in the
	// plan model) so the update sends the real password. See Create.
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("admin_account_password"), &plan.AdminAccountPassword)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Updating %s with ID: %s", ResourceName, state.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(
		ctx,
		plan.Timeouts.Update,
		UpdateTimeout*time.Second,
		&resp.Diagnostics,
	)
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

	depId := state.DepOnboardingSettingsID.ValueString()
	if depId == "" {
		var err error
		depId, err = r.resolveDepOnboardingSettingsId(ctx, state.DepOnboardingSettingsID)
		if err != nil {
			errors.HandleKiotaGraphError(
				ctx,
				err,
				resp,
				constants.TfOperationUpdate,
				r.WritePermissions,
			)
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
		errors.HandleKiotaGraphError(
			ctx,
			err,
			resp,
			constants.TfOperationUpdate,
			r.WritePermissions,
		)
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

	tflog.Debug(
		ctx,
		fmt.Sprintf("Finished updating %s with ID: %s", ResourceName, state.ID.ValueString()),
	)
}

// Delete handles the Delete operation for macOS DEP Enrollment Profile resources.
func (r *MacOSDepEnrollmentProfileResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	var object MacOSDepEnrollmentProfileResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(
		ctx,
		object.Timeouts.Delete,
		DeleteTimeout*time.Second,
		&resp.Diagnostics,
	)
	if cancel == nil {
		return
	}
	defer cancel()

	depId := object.DepOnboardingSettingsID.ValueString()
	if depId == "" {
		var err error
		depId, err = r.resolveDepOnboardingSettingsId(ctx, object.DepOnboardingSettingsID)
		if err != nil {
			errors.HandleKiotaGraphError(
				ctx,
				err,
				resp,
				constants.TfOperationDelete,
				r.WritePermissions,
			)
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
		errors.HandleKiotaGraphError(
			ctx,
			err,
			resp,
			constants.TfOperationDelete,
			r.WritePermissions,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
