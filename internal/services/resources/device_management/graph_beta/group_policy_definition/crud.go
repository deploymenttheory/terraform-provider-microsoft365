package graphBetaGroupPolicyDefinition

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Create handles the Create operation
func (r *GroupPolicyDefinitionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object GroupPolicyDefinitionResourceModel

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

	// Initialize AdditionalData
	object.AdditionalData = make(map[string]any)

	// Resolve policy definition and presentations
	err := resolveGroupPolicyDefinition(ctx, r.client, &object, "create")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error resolving policy definition",
			fmt.Sprintf("Could not resolve policy definition and presentations: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	// Validate user-provided values against catalog
	if err := ValidateValues(ctx, r.client, &object); err != nil {
		resp.Diagnostics.AddError(
			"Validation Error",
			fmt.Sprintf("Invalid configuration: %s", err.Error()),
		)
		return
	}

	// Construct the API request
	requestBody, err := constructResource(ctx, &object, "create")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing request",
			fmt.Sprintf("Could not construct updateDefinitionValues request: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	// Call the API
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

	// Resolve instance IDs after creation
	tflog.Debug(ctx, "[CREATE] About to resolve instance IDs after API creation")
	err = resolveGroupPolicyDefinition(ctx, r.client, &object, "read")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error resolving instance IDs after creation",
			fmt.Sprintf("Could not resolve instance IDs after creation: %s: %s", ResourceName, err.Error()),
		)
		return
	}
	tflog.Debug(ctx, fmt.Sprintf("[CREATE] Resolved presentations in AdditionalData: %+v", object.AdditionalData["resolvedPresentations"]))

	// Save to state - this will be updated by ReadWithRetry
	tflog.Debug(ctx, "[CREATE] Saving initial state before ReadWithRetry")
	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Debug(ctx, "[CREATE] Initial state saved, about to call ReadWithRetry")

	// Read back to ensure consistency
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

// Read handles the Read operation
func (r *GroupPolicyDefinitionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object GroupPolicyDefinitionResourceModel

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

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Initialize AdditionalData if needed
	if object.AdditionalData == nil {
		object.AdditionalData = make(map[string]any)
	}

	// Resolve policy definition and presentations
	err := resolveGroupPolicyDefinition(ctx, r.client, &object, "read")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error resolving policy definition",
			fmt.Sprintf("Could not resolve policy definition: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	// Get the definition value instance ID
	definitionValueInstanceID, ok := object.AdditionalData["definitionValueInstanceID"].(string)
	if !ok {
		resp.Diagnostics.AddError(
			"Error resolving instance ID",
			"Could not find definition value instance ID",
		)
		return
	}

	// Read the definition value from the API
	groupPolicyConfigurationID := object.GroupPolicyConfigurationID.ValueString()

	definitionValue, err := r.client.
		DeviceManagement().
		GroupPolicyConfigurations().
		ByGroupPolicyConfigurationId(groupPolicyConfigurationID).
		DefinitionValues().
		ByGroupPolicyDefinitionValueId(definitionValueInstanceID).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	// Get presentation values using the collection endpoint with $expand=presentation
	// Note: We don't actually need the presentation reference here since state mapping
	// uses the resolvedPresentations from AdditionalData, but we keep it for consistency
	presentationValuesResponse, err := r.client.
		DeviceManagement().
		GroupPolicyConfigurations().
		ByGroupPolicyConfigurationId(groupPolicyConfigurationID).
		DefinitionValues().
		ByGroupPolicyDefinitionValueId(definitionValueInstanceID).
		PresentationValues().
		Get(ctx, nil)

	var presentationValues []graphmodels.GroupPolicyPresentationValueable
	if err != nil {
		tflog.Warn(ctx, fmt.Sprintf("[READ] Failed to get presentation values: %s", err.Error()))
		presentationValues = []graphmodels.GroupPolicyPresentationValueable{}
	} else if presentationValuesResponse != nil {
		presentationValues = presentationValuesResponse.GetValue()
		if presentationValues == nil {
			presentationValues = []graphmodels.GroupPolicyPresentationValueable{}
		}
	} else {
		presentationValues = []graphmodels.GroupPolicyPresentationValueable{}
	}
	tflog.Debug(ctx, fmt.Sprintf("[READ] Got %d presentation values from API", len(presentationValues)))

	// Map remote state to Terraform
	tflog.Debug(ctx, "[READ] About to call MapRemoteStateToTerraform")
	MapRemoteStateToTerraform(ctx, &object, presentationValues, definitionValue)
	tflog.Debug(ctx, fmt.Sprintf("[READ] After mapping, object.Values type: %T, null: %v, unknown: %v", object.Values, object.Values.IsNull(), object.Values.IsUnknown()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	tflog.Debug(ctx, "[READ] State has been set")

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation
func (r *GroupPolicyDefinitionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var object GroupPolicyDefinitionResourceModel

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

	// Initialize AdditionalData
	if object.AdditionalData == nil {
		object.AdditionalData = make(map[string]any)
	}

	// Resolve policy definition and presentations
	err := resolveGroupPolicyDefinition(ctx, r.client, &object, "update")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error resolving policy definition",
			fmt.Sprintf("Could not resolve policy definition: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	// Validate user-provided values against catalog
	if err := ValidateValues(ctx, r.client, &object); err != nil {
		resp.Diagnostics.AddError(
			"Validation Error",
			fmt.Sprintf("Invalid configuration: %s", err.Error()),
		)
		return
	}

	// Construct the update request
	requestBody, err := constructResource(ctx, &object, "update")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing update request",
			fmt.Sprintf("Could not construct updateDefinitionValues request: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	// Call the API
	groupPolicyConfigurationID := object.GroupPolicyConfigurationID.ValueString()

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
	if resp.Diagnostics.HasError() {
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

	tflog.Debug(ctx, fmt.Sprintf("Finished Update Method: %s", ResourceName))
}

// Delete handles the Delete operation
func (r *GroupPolicyDefinitionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object GroupPolicyDefinitionResourceModel

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

	// Initialize AdditionalData
	if object.AdditionalData == nil {
		object.AdditionalData = make(map[string]any)
	}

	// Resolve to get instance IDs
	err := resolveGroupPolicyDefinition(ctx, r.client, &object, "read")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error resolving policy definition for deletion",
			fmt.Sprintf("Could not resolve policy definition: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	// Construct the delete request
	requestBody, err := constructResource(ctx, &object, "delete")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing delete request",
			fmt.Sprintf("Could not construct delete request: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	// Call the API
	groupPolicyConfigurationID := object.GroupPolicyConfigurationID.ValueString()

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

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
