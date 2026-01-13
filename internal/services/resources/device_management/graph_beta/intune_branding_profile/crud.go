package graphBetaDeviceManagementIntuneBrandingProfile

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
)

// Create handles the Create operation for the Intune branding profile resource.
func (r *IntuneBrandingProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan IntuneBrandingProfileResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Step 1: Create the base profile with minimal required fields
	createBody, err := constructCreateResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing create resource",
			fmt.Sprintf("Could not construct create resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	created, err := r.client.
		DeviceManagement().
		IntuneBrandingProfiles().
		Post(ctx, createBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	plan.ID = types.StringValue(*created.GetId())

	// Step 2: Update the profile with all the remaining configuration
	updateBody, err := constructUpdateResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing update resource for create method",
			fmt.Sprintf("Could not construct update resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.DeviceManagement().
		IntuneBrandingProfiles().
		ByIntuneBrandingProfileId(plan.ID.ValueString()).
		Patch(ctx, updateBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	// Step 3: Call the patch method again to Upload landing page image if provided
	if plan.LandingPageCustomizedImage != nil {
		landingPageBody, err := constructLandingPageImageResource(ctx, &plan)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing landing page image resource for create method",
				fmt.Sprintf("Could not construct landing page image resource: %s: %s", ResourceName, err.Error()),
			)
			return
		}

		_, err = r.client.DeviceManagement().
			IntuneBrandingProfiles().
			ByIntuneBrandingProfileId(plan.ID.ValueString()).
			Patch(ctx, landingPageBody, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
			return
		}
	}

	requestAssignment, err := constructAssignment(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing assignment for Create Method",
			fmt.Sprintf("Could not construct assignment: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	err = r.client.
		DeviceManagement().
		IntuneBrandingProfiles().
		ByIntuneBrandingProfileId(plan.ID.ValueString()).
		Assign().
		Post(ctx, requestAssignment, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
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

// Read handles the Read operation for the Intune branding profile resource.
func (r *IntuneBrandingProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object IntuneBrandingProfileResourceModel

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

	tflog.Debug(ctx, fmt.Sprintf("Reading %s with ID: %s (operation: %s)", ResourceName, object.ID.ValueString(), operation))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	remote, err := r.client.
		DeviceManagement().
		IntuneBrandingProfiles().
		ByIntuneBrandingProfileId(object.ID.ValueString()).
		Get(ctx, &devicemanagement.IntuneBrandingProfilesIntuneBrandingProfileItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &devicemanagement.IntuneBrandingProfilesIntuneBrandingProfileItemRequestBuilderGetQueryParameters{
				Expand: []string{"assignments"},
			},
		})

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, remote)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for the Intune branding profile resource.
func (r *IntuneBrandingProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan IntuneBrandingProfileResourceModel
	var state IntuneBrandingProfileResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update method for: %s", ResourceName))

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

	requestBody, err := constructUpdateResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for update method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.DeviceManagement().
		IntuneBrandingProfiles().
		ByIntuneBrandingProfileId(state.ID.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	// Handle landing page image update if provided
	if plan.LandingPageCustomizedImage != nil {
		landingPageBody, err := constructLandingPageImageResource(ctx, &plan)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing landing page image resource for update method",
				fmt.Sprintf("Could not construct landing page image resource: %s: %s", ResourceName, err.Error()),
			)
			return
		}

		_, err = r.client.DeviceManagement().
			IntuneBrandingProfiles().
			ByIntuneBrandingProfileId(state.ID.ValueString()).
			Patch(ctx, landingPageBody, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}
	}

	// Always handle assignments - either update with new assignments or remove all assignments if nil
	requestAssignment, err := constructAssignment(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing assignment for update method",
			fmt.Sprintf("Could not construct assignment: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	err = r.client.
		DeviceManagement().
		IntuneBrandingProfiles().
		ByIntuneBrandingProfileId(state.ID.ValueString()).
		Assign().
		Post(ctx, requestAssignment, nil)

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

// Delete handles the Delete operation for the Intune branding profile resource.
func (r *IntuneBrandingProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object IntuneBrandingProfileResourceModel

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

	err := r.client.DeviceManagement().
		IntuneBrandingProfiles().
		ByIntuneBrandingProfileId(object.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfTfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
