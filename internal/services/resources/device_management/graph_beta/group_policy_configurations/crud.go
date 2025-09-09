package graphBetaGroupPolicyConfigurations

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
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// Helper functions for logging
func stringPtrToString(s *string) string {
	if s == nil {
		return "nil"
	}
	return *s
}

func boolPtrToBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// Create handles the creation of a Group Policy Configuration resource.
func (r *GroupPolicyConfigurationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object GroupPolicyConfigurationResourceModel

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

	// Step 1: Create the basic group policy configuration
	requestBody, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	// Convert the request body to the SDK model
	groupPolicyConfig := models.NewGroupPolicyConfiguration()
	if displayName, ok := requestBody["displayName"].(string); ok {
		groupPolicyConfig.SetDisplayName(&displayName)
	}
	if description, ok := requestBody["description"].(string); ok {
		groupPolicyConfig.SetDescription(&description)
	}
	if roleScopeTagIds, ok := requestBody["roleScopeTagIds"].([]string); ok {
		groupPolicyConfig.SetRoleScopeTagIds(roleScopeTagIds)
	}

	baseResource, err := r.client.
		DeviceManagement().
		GroupPolicyConfigurations().
		Post(ctx, groupPolicyConfig, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Create", r.WritePermissions)
		return
	}

	object.ID = types.StringValue(*baseResource.GetId())

	// Step 2: Update definition values if provided
	if !object.DefinitionValues.IsNull() && !object.DefinitionValues.IsUnknown() {
		tflog.Debug(ctx, fmt.Sprintf("Updating definition values for %s with ID: %s", ResourceName, *baseResource.GetId()))

		// Create lookup service for definition ID resolution
		lookupService := NewDefinitionLookupService(r.client)

		// Use the constructor to build the SDK request body
		requestBody, err := constructUpdateDefinitionValuesRequest(ctx, &object, lookupService)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing definition values request",
				fmt.Sprintf("Could not construct definition values request for %s: %s", ResourceName, err.Error()),
			)
			return
		}

		// Call the updateDefinitionValues action using the SDK
		tflog.Debug(ctx, fmt.Sprintf("About to call updateDefinitionValues API for %s with ID: %s", ResourceName, *baseResource.GetId()))

		err = r.client.
			DeviceManagement().
			GroupPolicyConfigurations().
			ByGroupPolicyConfigurationId(*baseResource.GetId()).
			UpdateDefinitionValues().
			Post(ctx, requestBody, nil)

		if err != nil {
			tflog.Error(ctx, fmt.Sprintf("updateDefinitionValues API call failed for %s with ID %s: %s", ResourceName, *baseResource.GetId(), err.Error()))
			resp.Diagnostics.AddError(
				"Error updating definition values",
				fmt.Sprintf("Could not update definition values for %s with ID %s: %s", ResourceName, *baseResource.GetId(), err.Error()),
			)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("Successfully called updateDefinitionValues API for %s with ID: %s", ResourceName, *baseResource.GetId()))

		tflog.Debug(ctx, fmt.Sprintf("Successfully updated definition values for %s with ID: %s", ResourceName, *baseResource.GetId()))
	}

	object.ID = types.StringValue(*baseResource.GetId())

	requestAssignment, err := constructAssignment(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing assignment for Create Method",
			fmt.Sprintf("Could not construct assignment: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.
		DeviceManagement().
		GroupPolicyConfigurations().
		ByGroupPolicyConfigurationId(object.ID.ValueString()).
		Assign().
		Post(ctx, requestAssignment, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Create", r.WritePermissions)
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

// Read handles the reading of a Group Policy Configuration resource.
func (r *GroupPolicyConfigurationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state GroupPolicyConfigurationResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Read Method: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Step 1: Read the basic group policy configuration
	groupPolicyConfig, err := r.client.
		DeviceManagement().
		GroupPolicyConfigurations().
		ByGroupPolicyConfigurationId(state.ID.ValueString()).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	// Use state mapping function
	err = MapRemoteResourceStateToTerraform(ctx, &state, groupPolicyConfig)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error mapping resource state",
			fmt.Sprintf("Could not map resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	// Step 2: Read definition values with expanded definition details
	tflog.Debug(ctx, fmt.Sprintf("[CRUD] Making API call: GET /deviceManagement/groupPolicyConfigurations/%s/definitionValues", state.ID.ValueString()))
	definitionValuesResponse, err := r.client.
		DeviceManagement().
		GroupPolicyConfigurations().
		ByGroupPolicyConfigurationId(state.ID.ValueString()).
		DefinitionValues().
		Get(ctx, &devicemanagement.GroupPolicyConfigurationsItemDefinitionValuesRequestBuilderGetRequestConfiguration{
			QueryParameters: &devicemanagement.GroupPolicyConfigurationsItemDefinitionValuesRequestBuilderGetQueryParameters{
				Expand: []string{"definition($select=id,classType,displayName,policyType,hasRelatedDefinitions,version,minUserCspVersion,minDeviceCspVersion)"},
			},
		})

	if err != nil {
		tflog.Error(ctx, fmt.Sprintf("[CRUD] Failed to get definition values: %s", err.Error()))
		errors.HandleKiotaGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	// Log raw response
	if definitionValuesResponse != nil {
		tflog.Debug(ctx, fmt.Sprintf("[CRUD] Definition values response count: %d", len(definitionValuesResponse.GetValue())))
		for i, defValue := range definitionValuesResponse.GetValue() {
			if defValue != nil {
				tflog.Debug(ctx, fmt.Sprintf("[CRUD] DefValue[%d]: ID=%s, Enabled=%v, PresentationValues=%d",
					i,
					stringPtrToString(defValue.GetId()),
					boolPtrToBool(defValue.GetEnabled()),
					len(defValue.GetPresentationValues())))

				for j, presValue := range defValue.GetPresentationValues() {
					if presValue != nil {
						tflog.Debug(ctx, fmt.Sprintf("[CRUD] DefValue[%d].PresValue[%d]: ID=%s, Type=%T",
							i, j,
							stringPtrToString(presValue.GetId()),
							presValue))
					}
				}
			}
		}
	}

	// Step 2a: Get detailed presentation values for each definition value
	if definitionValuesResponse != nil && definitionValuesResponse.GetValue() != nil {
		for _, defValue := range definitionValuesResponse.GetValue() {
			if defValue != nil && defValue.GetId() != nil && defValue.GetPresentationValues() != nil {
				for _, presValue := range defValue.GetPresentationValues() {
					if presValue != nil && presValue.GetId() != nil {
						// GET /deviceManagement/groupPolicyConfigurations/{groupPolicyConfigurationId}/definitionValues/{groupPolicyDefinitionValueId}/presentationValues/{groupPolicyPresentationValueId}/definitionValue
						apiURL := fmt.Sprintf("/deviceManagement/groupPolicyConfigurations/%s/definitionValues/%s/presentationValues/%s/definitionValue",
							state.ID.ValueString(), *defValue.GetId(), *presValue.GetId())
						tflog.Debug(ctx, fmt.Sprintf("[CRUD] Making detailed API call: GET %s", apiURL))

						detailResponse, err := r.client.
							DeviceManagement().
							GroupPolicyConfigurations().
							ByGroupPolicyConfigurationId(state.ID.ValueString()).
							DefinitionValues().
							ByGroupPolicyDefinitionValueId(*defValue.GetId()).
							PresentationValues().
							ByGroupPolicyPresentationValueId(*presValue.GetId()).
							DefinitionValue().
							Get(ctx, nil)

						if err != nil {
							tflog.Debug(ctx, fmt.Sprintf("[CRUD] Failed detailed API call for %s: %s", apiURL, err.Error()))
						} else {
							tflog.Debug(ctx, fmt.Sprintf("[CRUD] Detailed API call successful for %s: Response=%v", apiURL, detailResponse != nil))
						}
					}
				}
			}
		}
	}

	// Map definition values to Terraform state
	if definitionValuesResponse != nil && definitionValuesResponse.GetValue() != nil {
		// Create lookup service for reverse lookup
		lookupService := NewDefinitionLookupService(r.client)

		err = MapRemoteDefinitionValuesToTerraform(ctx, &state, definitionValuesResponse.GetValue(), lookupService)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error mapping definition values state",
				fmt.Sprintf("Could not map definition values state: %s: %s", ResourceName, err.Error()),
			)
			return
		}
	} else {
		// Set to null if no definition values
		state.DefinitionValues = types.SetNull(types.ObjectType{
			AttrTypes: getDefinitionValueAttrTypes(),
		})
	}

	// Step 3: Read assignments
	assignmentsResponse, err := r.client.
		DeviceManagement().
		GroupPolicyConfigurations().
		ByGroupPolicyConfigurationId(state.ID.ValueString()).
		Assignments().
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Read", r.ReadPermissions)
		return
	}

	if assignmentsResponse != nil {
		MapAssignmentsToTerraform(ctx, &state, assignmentsResponse.GetValue())
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the updating of a Group Policy Configuration resource.
func (r *GroupPolicyConfigurationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan GroupPolicyConfigurationResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update Method: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	// Step 1: Update the basic group policy configuration properties
	requestBody, err := constructResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource for update: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	// Convert the request body to the SDK model
	groupPolicyConfig := models.NewGroupPolicyConfiguration()
	if displayName, ok := requestBody["displayName"].(string); ok {
		groupPolicyConfig.SetDisplayName(&displayName)
	}
	if description, ok := requestBody["description"].(string); ok {
		groupPolicyConfig.SetDescription(&description)
	}
	if roleScopeTagIds, ok := requestBody["roleScopeTagIds"].([]string); ok {
		groupPolicyConfig.SetRoleScopeTagIds(roleScopeTagIds)
	}

	_, err = r.client.
		DeviceManagement().
		GroupPolicyConfigurations().
		ByGroupPolicyConfigurationId(plan.ID.ValueString()).
		Patch(ctx, groupPolicyConfig, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
	}

	// Step 2: Update definition values if provided
	if !plan.DefinitionValues.IsNull() && !plan.DefinitionValues.IsUnknown() {
		tflog.Debug(ctx, fmt.Sprintf("Updating definition values for %s with ID: %s", ResourceName, plan.ID.ValueString()))

		// Create lookup service for definition ID resolution
		lookupService := NewDefinitionLookupService(r.client)

		// Use the constructor to build the SDK request body
		requestBody, err := constructUpdateDefinitionValuesRequest(ctx, &plan, lookupService)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing definition values request",
				fmt.Sprintf("Could not construct definition values request for %s: %s", ResourceName, err.Error()),
			)
			return
		}

		// Call the updateDefinitionValues action using the SDK
		err = r.client.DeviceManagement().
			GroupPolicyConfigurations().
			ByGroupPolicyConfigurationId(plan.ID.ValueString()).
			UpdateDefinitionValues().
			Post(ctx, requestBody, nil)

		if err != nil {
			resp.Diagnostics.AddError(
				"Error updating definition values",
				fmt.Sprintf("Could not update definition values for %s with ID %s: %s", ResourceName, plan.ID.ValueString(), err.Error()),
			)
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("Successfully updated definition values for %s with ID: %s", ResourceName, plan.ID.ValueString()))
	}

	// Step 3: Update assignments
	requestAssignment, err := constructAssignment(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing assignment for Update Method",
			fmt.Sprintf("Could not construct assignment: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.
		DeviceManagement().
		GroupPolicyConfigurations().
		ByGroupPolicyConfigurationId(plan.ID.ValueString()).
		Assign().
		Post(ctx, requestAssignment, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
	}

	// Step 4: Read the updated resource state
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

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

// Delete handles the deletion of a Group Policy Configuration resource.
func (r *GroupPolicyConfigurationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state GroupPolicyConfigurationResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Delete Method: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(ctx, state.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	err := r.client.
		DeviceManagement().
		GroupPolicyConfigurations().
		ByGroupPolicyConfigurationId(state.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
