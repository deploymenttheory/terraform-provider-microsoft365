package graphBetaCloudPcOrganizationSettings

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
)

// Create handles the Create operation for Cloud PC organization settings resources.
//
// Operation: Configures organization-wide Cloud PC settings
// API Calls:
//   - PATCH /deviceManagement/virtualEndpoint/organizationSettings
//
// Reference: https://learn.microsoft.com/en-us/graph/api/cloudpcorganizationsettings-update?view=graph-rest-beta
// Note: This is a singleton resource; settings always exist and are configured rather than created
func (r *CloudPcOrganizationSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CloudPcOrganizationSettingsResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting creation of singleton resource: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.
		DeviceManagement().
		VirtualEndpoint().
		OrganizationSettings().
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	plan.ID = types.StringValue(SingletonID)

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

// Read handles the Read operation for Cloud PC organization settings resources.
//
// Operation: Retrieves organization-wide Cloud PC settings
// API Calls:
//   - GET /deviceManagement/virtualEndpoint/organizationSettings
//
// Reference: https://learn.microsoft.com/en-us/graph/api/cloudpcorganizationsettings-get?view=graph-rest-beta
func (r *CloudPcOrganizationSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CloudPcOrganizationSettingsResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read method for singleton: %s", ResourceName))

	operation := constants.TfOperationRead
	if ctxOp := ctx.Value("retry_operation"); ctxOp != nil {
		if opStr, ok := ctxOp.(string); ok {
			operation = opStr
		}
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	remote, err := r.client.
		DeviceManagement().
		VirtualEndpoint().
		OrganizationSettings().
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, &state, remote)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for Cloud PC organization settings resources.
//
// Operation: Updates organization-wide Cloud PC settings
// API Calls:
//   - PATCH /deviceManagement/virtualEndpoint/organizationSettings
//
// Reference: https://learn.microsoft.com/en-us/graph/api/cloudpcorganizationsettings-update?view=graph-rest-beta
func (r *CloudPcOrganizationSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CloudPcOrganizationSettingsResourceModel
	var state CloudPcOrganizationSettingsResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Updating singleton resource: %s", ResourceName))

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
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.
		DeviceManagement().
		VirtualEndpoint().
		OrganizationSettings().
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	plan.ID = types.StringValue(SingletonID)

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
	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s", ResourceName))
}

// Delete handles the Delete operation for Cloud PC organization settings resources.
//
// Operation: Resets organization settings to default values
// API Calls:
//   - PATCH /deviceManagement/virtualEndpoint/organizationSettings
//
// Reference: https://learn.microsoft.com/en-us/graph/api/cloudpcorganizationsettings-update?view=graph-rest-beta
// Note: Singleton resource cannot be deleted; this operation resets settings to defaults and removes from Terraform state
func (r *CloudPcOrganizationSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CloudPcOrganizationSettingsResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting deletion (reset) of singleton resource: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	reset := CloudPcOrganizationSettingsResourceModel{
		ID:                  types.StringValue(SingletonID),
		EnableMEMAutoEnroll: types.BoolValue(false),
		EnableSingleSignOn:  types.BoolValue(false),
		OsVersion:           types.StringValue("windows11"),    // Default as per docs
		UserAccountType:     types.StringValue("standardUser"), // Default as per docs
		WindowsSettings:     nil,                               // Default
	}

	requestBody, err := constructResource(ctx, &reset)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.
		DeviceManagement().
		VirtualEndpoint().
		OrganizationSettings().
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}
	resp.State.RemoveResource(ctx)
	tflog.Debug(ctx, fmt.Sprintf("Finished Delete (reset) Method: %s", ResourceName))
}
