package graphBetaWindowsUpdatesAutopatchDeviceRegistration

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/admin"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/windowsupdates"
)

// Create handles the Create operation for Windows Updates autopatch device registration resources.
//
// Operation: Enrolls devices into Windows Update for Business deployment service
// API Calls:
//   - POST /admin/windows/updates/updatableAssets/microsoft.graph.windowsUpdates.enrollAssetsById
//
// Reference: https://learn.microsoft.com/en-us/graph/api/windowsupdates-updatableasset-enrollassetsbyid?view=graph-rest-beta
func (r *WindowsUpdatesAutopatchDeviceRegistrationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object WindowsUpdatesAutopatchDeviceRegistrationResourceModel

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

	object.ID = types.StringValue(object.UpdateCategory.ValueString())

	if object.EntraDeviceObjectIds.IsNull() || object.EntraDeviceObjectIds.IsUnknown() {
		object.EntraDeviceObjectIds = types.SetValueMust(types.StringType, []attr.Value{})
	}

	if err := r.validateRequest(ctx, &object, &resp.Diagnostics); err != nil {
		return
	}

	requestBody, err := constructEnrollRequest(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	err = r.client.
		Admin().
		Windows().
		Updates().
		UpdatableAssets().
		MicrosoftGraphWindowsUpdatesEnrollAssetsById().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, "Enroll POST completed, verifying enrollment propagation")

	if err := r.verifyEnrollmentComplete(ctx, &object); err != nil {
		resp.Diagnostics.AddError(
			"Enrollment Verification Failed",
			fmt.Sprintf("Devices were enrolled but verification failed: %s", err.Error()),
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

// Read handles the Read operation for Windows Updates autopatch device registration resources.
//
// Operation: Retrieves all enrolled azureADDevice objects and filters by update category
// API Calls:
//   - GET /admin/windows/updates/updatableAssets?$filter=isof('microsoft.graph.windowsUpdates.azureADDevice')
//
// Reference: https://learn.microsoft.com/en-us/graph/api/adminwindowsupdates-list-updatableassets-azureaddevice?view=graph-rest-beta&tabs=go
//
// Important Notes:
//   - The $filter parameter is required to query only azureADDevice objects (derived type)
//   - Tested $select parameter - it does not work with derived type properties like 'enrollment'
//   - The API returns full objects including enrollment information by default
//   - State mapping behavior:
//   - Normal Read: Filters by planned device IDs from state, returns only managed devices
//   - Import: No planned IDs available, returns ALL devices enrolled for the category
//   - Import behavior: Returns all devices currently enrolled for the specified update category,
//     which may differ from the configuration if other devices are enrolled externally
func (r *WindowsUpdatesAutopatchDeviceRegistrationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object WindowsUpdatesAutopatchDeviceRegistrationResourceModel
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

	// During import, update_category may not be set but id contains the update category value
	if object.UpdateCategory.IsNull() || object.UpdateCategory.IsUnknown() || object.UpdateCategory.ValueString() == "" {
		if !object.ID.IsNull() && !object.ID.IsUnknown() && object.ID.ValueString() != "" {
			object.UpdateCategory = object.ID
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Reading %s for update category: %s", ResourceName, object.UpdateCategory.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Fetch all azureADDevice assets using a type filter
	// This filters on microsoft.graph.windowsUpdates.azureADDevice
	// and returns devices with their enrollment information
	filter := "isof('microsoft.graph.windowsUpdates.azureADDevice')"
	result, err := r.client.
		Admin().
		Windows().
		Updates().
		UpdatableAssets().
		Get(ctx, &admin.WindowsUpdatesUpdatableAssetsRequestBuilderGetRequestConfiguration{
			QueryParameters: &admin.WindowsUpdatesUpdatableAssetsRequestBuilderGetQueryParameters{
				Filter: &filter,
			},
		})

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	var devices []windowsupdates.UpdatableAssetable
	if result != nil {
		devices = result.GetValue()
	}

	MapRemoteStateToTerraform(ctx, &object, devices)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	identity.ID = object.ID.ValueString()

	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, identity)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for Windows Updates autopatch device registration resources.
//
// Operation: Enrolls new devices and unenrolls removed devices from Windows Update for Business
// API Calls:
//   - POST /admin/windows/updates/updatableAssets/microsoft.graph.windowsUpdates.enrollAssetsById
//   - POST /admin/windows/updates/updatableAssets/microsoft.graph.windowsUpdates.unenrollAssetsById
//
// Reference: https://learn.microsoft.com/en-us/graph/api/windowsupdates-updatableasset-unenrollassetsbyid?view=graph-rest-beta&tabs=http
// https://learn.microsoft.com/en-us/graph/api/windowsupdates-updatableasset-enrollassetsbyid?view=graph-rest-beta
func (r *WindowsUpdatesAutopatchDeviceRegistrationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan WindowsUpdatesAutopatchDeviceRegistrationResourceModel
	var state WindowsUpdatesAutopatchDeviceRegistrationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Updating %s for update category: %s", ResourceName, state.UpdateCategory.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	planDeviceIDs := make(map[string]bool)
	if !plan.EntraDeviceObjectIds.IsNull() && !plan.EntraDeviceObjectIds.IsUnknown() {
		elements := plan.EntraDeviceObjectIds.Elements()
		for _, elem := range elements {
			if strVal, ok := elem.(types.String); ok {
				planDeviceIDs[strVal.ValueString()] = true
			}
		}
	}

	stateDeviceIDs := make(map[string]bool)
	if !state.EntraDeviceObjectIds.IsNull() && !state.EntraDeviceObjectIds.IsUnknown() {
		elements := state.EntraDeviceObjectIds.Elements()
		for _, elem := range elements {
			if strVal, ok := elem.(types.String); ok {
				stateDeviceIDs[strVal.ValueString()] = true
			}
		}
	}

	var devicesToEnroll []string
	for id := range planDeviceIDs {
		if !stateDeviceIDs[id] {
			devicesToEnroll = append(devicesToEnroll, id)
		}
	}

	var devicesToUnenroll []string
	for id := range stateDeviceIDs {
		if !planDeviceIDs[id] {
			devicesToUnenroll = append(devicesToUnenroll, id)
		}
	}

	if len(devicesToEnroll) > 0 {
		enrollModel := WindowsUpdatesAutopatchDeviceRegistrationResourceModel{
			UpdateCategory:       plan.UpdateCategory,
			EntraDeviceObjectIds: types.SetValueMust(types.StringType, convertStringsToAttrValues(devicesToEnroll)),
		}

		if err := r.validateRequest(ctx, &enrollModel, &resp.Diagnostics); err != nil {
			return
		}

		enrollRequest, err := constructEnrollRequest(ctx, &enrollModel)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing enroll request",
				fmt.Sprintf("Could not construct enroll request: %s: %s", ResourceName, err.Error()),
			)
			return
		}

		err = r.client.
			Admin().
			Windows().
			Updates().
			UpdatableAssets().
			MicrosoftGraphWindowsUpdatesEnrollAssetsById().
			Post(ctx, enrollRequest, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}

		tflog.Debug(ctx, "Enroll POST completed in Update, verifying enrollment propagation")

		if err := r.verifyEnrollmentComplete(ctx, &enrollModel); err != nil {
			resp.Diagnostics.AddError(
				"Enrollment Verification Failed",
				fmt.Sprintf("Devices were enrolled but verification failed: %s", err.Error()),
			)
			return
		}
	}

	if len(devicesToUnenroll) > 0 {
		unenrollModel := WindowsUpdatesAutopatchDeviceRegistrationResourceModel{
			UpdateCategory:       plan.UpdateCategory,
			EntraDeviceObjectIds: types.SetValueMust(types.StringType, convertStringsToAttrValues(devicesToUnenroll)),
		}

		unenrollRequest, err := constructUnenrollRequest(ctx, &unenrollModel)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing unenroll request",
				fmt.Sprintf("Could not construct unenroll request: %s: %s", ResourceName, err.Error()),
			)
			return
		}

		err = r.client.
			Admin().
			Windows().
			Updates().
			UpdatableAssets().
			MicrosoftGraphWindowsUpdatesUnenrollAssetsById().
			Post(ctx, unenrollRequest, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}

		tflog.Debug(ctx, "Unenroll POST completed in Update, verifying unenrollment propagation")

		if err := r.verifyUnenrollmentComplete(ctx, &unenrollModel); err != nil {
			resp.Diagnostics.AddError(
				"Unenrollment Verification Failed",
				fmt.Sprintf("Devices were unenrolled but verification failed: %s", err.Error()),
			)
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationUpdate
	opts.ResourceTypeName = ResourceName

	err := crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after update",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished updating %s for update category: %s", ResourceName, state.UpdateCategory.ValueString()))
}

// Delete handles the Delete operation for Windows Updates autopatch device registration resources.
//
// Operation: Unenrolls all devices from Windows Update for Business deployment service
// API Calls:
//   - POST /admin/windows/updates/updatableAssets/microsoft.graph.windowsUpdates.unenrollAssetsById
//
// Reference: https://learn.microsoft.com/en-us/graph/api/windowsupdates-updatableasset-unenrollassetsbyid?view=graph-rest-beta
func (r *WindowsUpdatesAutopatchDeviceRegistrationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object WindowsUpdatesAutopatchDeviceRegistrationResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Delete of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructUnenrollRequest(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing unenroll request",
			fmt.Sprintf("Could not construct unenroll request: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	err = r.client.
		Admin().
		Windows().
		Updates().
		UpdatableAssets().
		MicrosoftGraphWindowsUpdatesUnenrollAssetsById().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, "Unenroll POST completed, verifying unenrollment propagation")

	if err := r.verifyUnenrollmentComplete(ctx, &object); err != nil {
		resp.Diagnostics.AddError(
			"Unenrollment Verification Failed",
			fmt.Sprintf("Devices were unenrolled but verification failed: %s", err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}

func convertStringsToAttrValues(strings []string) []attr.Value {
	values := make([]attr.Value, len(strings))
	for i, s := range strings {
		values[i] = types.StringValue(s)
	}
	return values
}
