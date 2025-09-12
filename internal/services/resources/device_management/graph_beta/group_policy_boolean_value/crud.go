package graphBetaGroupPolicyBooleanValue

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

// Create handles the Create operation for Group Policy Boolean Value resources.
func (r *GroupPolicyBooleanValueResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object GroupPolicyBooleanValueResourceModel

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
	err, _ := GroupPolicyIDResolver(ctx, &object, r.client, "create")
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
	err, _ = GroupPolicyIDResolver(ctx, &object, r.client, "read")
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

// Read handles the Read operation for Group Policy Boolean Value resources.
func (r *GroupPolicyBooleanValueResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object GroupPolicyBooleanValueResourceModel

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

	err, statusCode := GroupPolicyIDResolver(ctx, &object, r.client, "read")
	if err != nil {
		if statusCode == 500 {
			warningTitle := fmt.Sprintf("Reading of resource %s with group policy configuration ID %s failed. Policy setting no longer exists and will be removed from state.", ResourceName, object.GroupPolicyConfigurationID.ValueString())
			warningDetail := fmt.Sprintf("The group policy setting '%s' has been removed from the group policy configuration. %s. Removing from state.", object.PolicyName.ValueString(), err.Error())
			resp.Diagnostics.AddWarning(warningTitle, warningDetail)
			resp.State.RemoveResource(ctx)
			return
		} else {
			resp.Diagnostics.AddError(
				"Error resolving IDs during read",
				fmt.Sprintf("Could not resolve definition and presentation IDs for %s (HTTP %d): %s", ResourceName, statusCode, err.Error()),
			)
			return
		}
	}

	groupPolicyConfigurationID := object.GroupPolicyConfigurationID.ValueString()
	groupPolicyDefinitionValueID := object.GroupPolicyDefinitionValueID.ValueString()

	if groupPolicyDefinitionValueID == "" {
		resp.Diagnostics.AddError(
			"Resource not found",
			fmt.Sprintf("Could not find %s resource in configuration", ResourceName),
		)
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

	// Get all presentation values for this definition value
	allPresentationValues, err := r.client.
		DeviceManagement().
		GroupPolicyConfigurations().
		ByGroupPolicyConfigurationId(groupPolicyConfigurationID).
		DefinitionValues().
		ByGroupPolicyDefinitionValueId(groupPolicyDefinitionValueID).
		PresentationValues().
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	if allPresentationValues == nil || allPresentationValues.GetValue() == nil {
		resp.Diagnostics.AddError(
			"No presentation values found",
			fmt.Sprintf("No presentation values found for %s resource in configuration", ResourceName),
		)
		return
	}

	// Filter for boolean presentation values only
	var booleanPresentationValues []graphmodels.GroupPolicyPresentationValueable
	for _, presValue := range allPresentationValues.GetValue() {
		if presValue == nil {
			continue
		}

		odataType := presValue.GetOdataType()
		if odataType != nil && *odataType == "#microsoft.graph.groupPolicyPresentationValueBoolean" {
			booleanPresentationValues = append(booleanPresentationValues, presValue)
		}
	}

	if len(booleanPresentationValues) == 0 {
		resp.Diagnostics.AddError(
			"No boolean presentation values found",
			fmt.Sprintf("No boolean presentation values found for %s resource in configuration", ResourceName),
		)
		return
	}

	// Map the presentation values data
	MapRemoteStateToTerraform(ctx, &object, booleanPresentationValues, definitionValue)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for Group Policy Boolean Value resources.
func (r *GroupPolicyBooleanValueResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var object GroupPolicyBooleanValueResourceModel

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
	err, _ := GroupPolicyIDResolver(ctx, &object, r.client, "update")
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

// Delete handles the Delete operation for Group Policy Boolean Value resources.
func (r *GroupPolicyBooleanValueResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object GroupPolicyBooleanValueResourceModel

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
