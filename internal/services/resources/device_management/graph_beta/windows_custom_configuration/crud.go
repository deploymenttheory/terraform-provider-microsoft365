package graphBetaWindowsCustomConfiguration

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Create handles the Create operation.
func (r *WindowsCustomConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object WindowsCustomConfigurationResourceModel

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
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	requestResource, err := r.client.
		DeviceManagement().
		DeviceConfigurations().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	object.ID = types.StringValue(*requestResource.GetId())

	requestAssignment, err := constructAssignment(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing assignment for create method",
			fmt.Sprintf("Could not construct assignment: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.
		DeviceManagement().
		DeviceConfigurations().
		ByDeviceConfigurationId(object.ID.ValueString()).
		Assign().
		Post(ctx, requestAssignment, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State}
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

// Read handles the Read operation for Windows custom configuration profiles.
//
//   - Retrieves the current state from the read request
//   - Gets the resource details with expanded assignments from the API
//   - Resolves encrypted OMA setting values back to their plain text values
//   - Maps the resource details and assignments to Terraform state
//
// The Graph API masks encrypted OMA setting values as "****" in GET responses, so the
// plain text values are retrieved via the getOmaSettingPlainTextValue function to allow
// drift detection against the configured values.
func (r *WindowsCustomConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object WindowsCustomConfigurationResourceModel
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

	respResource, err := r.client.
		DeviceManagement().
		DeviceConfigurations().
		ByDeviceConfigurationId(object.ID.ValueString()).
		Get(ctx, &devicemanagement.DeviceConfigurationsDeviceConfigurationItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &devicemanagement.DeviceConfigurationsDeviceConfigurationItemRequestBuilderGetQueryParameters{
				Expand: []string{"assignments"},
			},
		})

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	if winConfig, ok := respResource.(graphmodels.Windows10CustomConfigurationable); ok {
		r.resolveEncryptedOmaSettingValues(ctx, &object, winConfig)
	}

	MapRemoteResourceStateToTerraform(ctx, &object, respResource)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// resolveEncryptedOmaSettingValues replaces masked ("****") values of encrypted OMA settings
// with their plain text values retrieved via the getOmaSettingPlainTextValue function.
// If the plain text value cannot be retrieved, the value from the current Terraform state is
// preserved (matched by OMA-URI) to avoid perpetual drift.
func (r *WindowsCustomConfigurationResource) resolveEncryptedOmaSettingValues(ctx context.Context, object *WindowsCustomConfigurationResourceModel, config graphmodels.Windows10CustomConfigurationable) {
	stateValues := make(map[string]string)
	if !object.OmaSettings.IsNull() && !object.OmaSettings.IsUnknown() {
		var stateSettings []OmaSettingResourceModel
		if diags := object.OmaSettings.ElementsAs(ctx, &stateSettings, false); !diags.HasError() {
			for _, stateSetting := range stateSettings {
				if !stateSetting.Value.IsNull() && !stateSetting.Value.IsUnknown() {
					stateValues[stateSetting.OmaUri.ValueString()] = stateSetting.Value.ValueString()
				}
			}
		}
	}

	for _, setting := range config.GetOmaSettings() {
		if setting == nil || setting.GetIsEncrypted() == nil || !*setting.GetIsEncrypted() {
			continue
		}

		var plainTextValue *string
		if secretReferenceValueId := setting.GetSecretReferenceValueId(); secretReferenceValueId != nil {
			plainTextResponse, err := r.client.
				DeviceManagement().
				DeviceConfigurations().
				ByDeviceConfigurationId(object.ID.ValueString()).
				GetOmaSettingPlainTextValueWithSecretReferenceValueId(secretReferenceValueId).
				GetAsGetOmaSettingPlainTextValueWithSecretReferenceValueIdGetResponse(ctx, nil)
			if err != nil {
				tflog.Warn(ctx, "Failed to retrieve plain text value for encrypted oma setting, falling back to state value", map[string]any{
					"omaUri": setting.GetOmaUri(),
					"error":  err.Error(),
				})
			} else {
				plainTextValue = plainTextResponse.GetValue()
			}
		}

		if plainTextValue == nil {
			if omaUri := setting.GetOmaUri(); omaUri != nil {
				if stateValue, ok := stateValues[*omaUri]; ok {
					plainTextValue = &stateValue
				}
			}
		}

		if plainTextValue == nil {
			continue
		}

		switch typedSetting := setting.(type) {
		case *graphmodels.OmaSettingString:
			typedSetting.SetValue(plainTextValue)
		case *graphmodels.OmaSettingBase64:
			typedSetting.SetValue(plainTextValue)
		case *graphmodels.OmaSettingStringXml:
			typedSetting.SetValue([]byte(*plainTextValue))
		default:
			tflog.Warn(ctx, "Encrypted oma setting has an unexpected type, leaving masked value", map[string]any{
				"omaUri": setting.GetOmaUri(),
				"type":   fmt.Sprintf("%T", setting),
			})
		}
	}
}

// Update handles the Update operation for Windows custom configuration profiles.
//
// The function performs the following operations:
//   - Patches the existing configuration resource with updated settings using PATCH
//   - Updates assignments using POST if they are defined or removes all assignments if nil
//   - Retrieves the updated resource with expanded assignments
//   - Maps the remote state back to Terraform
func (r *WindowsCustomConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan WindowsCustomConfigurationResourceModel
	var state WindowsCustomConfigurationResourceModel

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
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.
		DeviceManagement().
		DeviceConfigurations().
		ByDeviceConfigurationId(state.ID.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
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

	_, err = r.client.
		DeviceManagement().
		DeviceConfigurations().
		ByDeviceConfigurationId(plan.ID.ValueString()).
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

	readReq := resource.ReadRequest{State: resp.State}
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

// Delete handles the Delete operation for Windows custom configuration profiles.
//
//   - Retrieves the current state from the delete request
//   - Validates the state data and timeout configuration
//   - Sends DELETE request to remove the resource from the API
//   - Cleans up by removing the resource from Terraform state
//
// All assignments and settings associated with the resource are automatically removed as part of the deletion.
func (r *WindowsCustomConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object WindowsCustomConfigurationResourceModel
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

	err := r.client.
		DeviceManagement().
		DeviceConfigurations().
		ByDeviceConfigurationId(object.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
