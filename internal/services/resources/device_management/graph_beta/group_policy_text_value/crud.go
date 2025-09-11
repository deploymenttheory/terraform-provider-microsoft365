package graphBetaGroupPolicyTextValue

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Create handles the Create operation for Group Policy Text Value resources.
func (r *GroupPolicyTextValueResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object GroupPolicyTextValueResourceModel

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

	// Resolve definition and presentation IDs for creation
	err := groupPolicyIDResolver(ctx, &object, r.client, "create")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error resolving IDs",
			fmt.Sprintf("Could not resolve definition and presentation IDs: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	requestBody, err := constructResource(ctx, &object, "create")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing updateDefinitionValues request",
			fmt.Sprintf("Could not construct updateDefinitionValues request: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	groupPolicyConfigurationID := object.GroupPolicyConfigurationID.ValueString()

	err = r.client.
		DeviceManagement().
		GroupPolicyConfigurations().
		ByGroupPolicyConfigurationId(groupPolicyConfigurationID).
		UpdateDefinitionValues().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	// After creating via updateDefinitionValues, we need to resolve again and return the instance IDs
	err = groupPolicyIDResolver(ctx, &object, r.client, "read")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error resolving instance IDs after creation",
			fmt.Sprintf("Could not resolve definition value and presentation value instance IDs after creation: %s: %s", ResourceName, err.Error()),
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
	opts.Operation = "Create"
	opts.ResourceTypeName = constants.PROVIDER_NAME + "_" + ResourceName

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

// Read handles the Read operation for Group Policy Text Value resources.
func (r *GroupPolicyTextValueResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object GroupPolicyTextValueResourceModel

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

	tflog.Debug(ctx, fmt.Sprintf("Reading %s with configuration ID: %s", ResourceName, object.GroupPolicyConfigurationID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Resolve definition and presentation IDs for reading
	err := groupPolicyIDResolver(ctx, &object, r.client, "read")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error resolving IDs during read",
			fmt.Sprintf("Could not resolve definition and presentation IDs: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	groupPolicyConfigurationID := object.GroupPolicyConfigurationID.ValueString()
	groupPolicyDefinitionValueID := object.GroupPolicyDefinitionValueID.ValueString()
	presentationValueID := object.ID.ValueString()

	if presentationValueID == "" || groupPolicyDefinitionValueID == "" {
		resp.Diagnostics.AddError(
			"Resource not found",
			fmt.Sprintf("Could not find %s resource in configuration", ResourceName),
		)
		return
	}

	presentationValue, err := r.client.
		DeviceManagement().
		GroupPolicyConfigurations().
		ByGroupPolicyConfigurationId(groupPolicyConfigurationID).
		DefinitionValues().
		ByGroupPolicyDefinitionValueId(groupPolicyDefinitionValueID).
		PresentationValues().
		ByGroupPolicyPresentationValueId(presentationValueID).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	definitionValue, err := r.client.
		DeviceManagement().
		GroupPolicyConfigurations().
		ByGroupPolicyConfigurationId(groupPolicyConfigurationID).
		DefinitionValues().
		ByGroupPolicyDefinitionValueId(groupPolicyDefinitionValueID).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	// Map the presentation value data
	if textValue, ok := presentationValue.(graphmodels.GroupPolicyPresentationValueTextable); ok {
		MapRemoteStateToTerraform(ctx, &object, textValue, definitionValue)
	} else {
		resp.Diagnostics.AddError(
			"Type assertion error",
			fmt.Sprintf("Could not cast response to GroupPolicyPresentationValueText for resource: %s", ResourceName),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for Group Policy Text Value resources.
func (r *GroupPolicyTextValueResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var object GroupPolicyTextValueResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Resolve definition and presentation IDs for update
	err := groupPolicyIDResolver(ctx, &object, r.client, "update")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error resolving IDs for update",
			fmt.Sprintf("Could not resolve definition and presentation IDs: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	// Construct the updateDefinitionValues request for update
	requestBody, err := constructResource(ctx, &object, "update")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for update",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	groupPolicyConfigurationID := object.GroupPolicyConfigurationID.ValueString()

	// Call updateDefinitionValues to update the definition value with presentation value
	err = r.client.
		DeviceManagement().
		GroupPolicyConfigurations().
		ByGroupPolicyConfigurationId(groupPolicyConfigurationID).
		UpdateDefinitionValues().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = "Update"
	opts.ResourceTypeName = constants.PROVIDER_NAME + "_" + ResourceName

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

// Delete handles the Delete operation for Group Policy Text Value resources.
func (r *GroupPolicyTextValueResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object GroupPolicyTextValueResourceModel

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

	groupPolicyConfigurationID := object.GroupPolicyConfigurationID.ValueString()

	requestBody, err := constructResource(ctx, &object, "delete")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing delete request",
			fmt.Sprintf("Could not construct updateDefinitionValues delete request: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	err = r.client.
		DeviceManagement().
		GroupPolicyConfigurations().
		ByGroupPolicyConfigurationId(groupPolicyConfigurationID).
		UpdateDefinitionValues().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))

	resp.State.RemoveResource(ctx)
}
