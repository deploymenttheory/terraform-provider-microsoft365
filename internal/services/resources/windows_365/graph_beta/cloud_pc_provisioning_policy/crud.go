// Package graphBetaCloudPcProvisioningPolicy manages Windows 365 Cloud PC provisioning policies.
//
// SDK deviation: unknown enum value
//
// Evidence:
//   - SDK: github.com/microsoftgraph/msgraph-beta-sdk-go@v0.168.0, models.CloudPcRegionGroup.
//     ParseCloudPcRegionGroup knows default through kenya and returns nil for uae, so the SDK cannot
//     set regionGroup = uae and its deserialiser drops uae from a GET response.
//   - Learn: https://learn.microsoft.com/en-us/graph/api/resources/cloudpcdomainjoinconfiguration?view=graph-rest-beta
//     links the cloudPcRegionGroup enum without uae; https://learn.microsoft.com/en-us/windows-365/enterprise/enhanced-resiliency-mhn
//     lists UAE as a region group of the Middle East geography.
//   - Observed: Intune portal create request captured in #4904, POST /beta/deviceManagement/virtualEndpoint/provisioningPolicies
//     with domainJoinConfigurations [{"type": "azureADJoin", "geographicLocationType": "middleEast", "regionName": "automatic",
//     "regionGroup": "uae"}]. No live GET was available; the read path was checked against SDK source.
//
// Handling: construct.go writes a region group the SDK enum lacks into the domain join configuration's additionalData,
// which the SDK serialiser writes as "regionGroup". Read fetches the policy as raw JSON through custom_requests,
// decodes it with the SDK model and copies each raw regionGroup the SDK dropped into additionalData, where state.go reads it.
//
// Removal condition: models.CloudPcRegionGroup includes uae.
package graphBetaCloudPcProvisioningPolicy

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
)

// Create handles the Create operation for Cloud PC provisioning policy resources.
//
// Operation: Creates a new Cloud PC provisioning policy with optional assignments
// API Calls:
//   - POST /deviceManagement/virtualEndpoint/provisioningPolicies
//   - POST /deviceManagement/virtualEndpoint/provisioningPolicies/{id}/assign (if assignments are configured)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/virtualendpoint-post-provisioningpolicies?view=graph-rest-beta
func (r *CloudPcProvisioningPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan CloudPcProvisioningPolicyResourceModel

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

	requestBody, err := constructResource(ctx, &plan, r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	baseResource, err := r.client.
		DeviceManagement().
		VirtualEndpoint().
		ProvisioningPolicies().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	plan.ID = types.StringValue(*baseResource.GetId())

	if !plan.Assignments.IsNull() && !plan.Assignments.IsUnknown() {
		// Get the number of assignments for logging
		var terraformAssignments []CloudPcProvisioningPolicyAssignmentModel
		diags := plan.Assignments.ElementsAs(ctx, &terraformAssignments, false)
		assignmentCount := 0
		if !diags.HasError() {
			assignmentCount = len(terraformAssignments)
		}

		tflog.Debug(ctx, fmt.Sprintf("Creating %d assignments for policy ID: %s", assignmentCount, plan.ID.ValueString()))

		assignBody, err := constructAssignmentsRequestBody(ctx, plan.Assignments)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing assignments request body",
				fmt.Sprintf("Could not construct assignments request body: %s", err.Error()),
			)
			return
		}

		err = r.client.
			DeviceManagement().
			VirtualEndpoint().
			ProvisioningPolicies().
			ByCloudPcProvisioningPolicyId(plan.ID.ValueString()).
			Assign().
			Post(ctx, assignBody, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
			return
		}

		tflog.Debug(ctx, "Successfully created assignments")
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

// Read handles the Read operation for Cloud PC provisioning policy resources.
//
// Operation: Retrieves a Cloud PC provisioning policy including assignments
// API Calls:
//   - GET /deviceManagement/virtualEndpoint/provisioningPolicies/{id}?$expand=assignments&$select=* (raw JSON, see the package doc)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/cloudpcprovisioningpolicy-get?view=graph-rest-beta
func (r *CloudPcProvisioningPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object CloudPcProvisioningPolicyResourceModel
	var identity sharedmodels.ResourceIdentity

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

	identity.ID = object.ID.ValueString()

	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, identity)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	tflog.Debug(ctx, "Fetching provisioning policy with expanded assignments")

	// SDK deviation: unknown enum value. The policy is read as raw JSON so a region group the SDK enum lacks (uae) survives.
	getConfig := customrequests.GetRequestConfig{
		APIVersion:        customrequests.GraphAPIBeta,
		Endpoint:          "deviceManagement/virtualEndpoint/provisioningPolicies",
		ResourceIDPattern: customrequests.ResourceIDPatternSegment,
		ResourceID:        object.ID.ValueString(),
		EndpointSuffix:    "?$expand=assignments&$select=*",
	}
	body, err := customrequests.GetRequestByResourceId(ctx, r.client.GetAdapter(), getConfig)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	node, err := jsonserialization.NewJsonParseNodeFactory().GetRootParseNode("application/json", body)
	if err != nil {
		resp.Diagnostics.AddError("Error parsing resource", fmt.Sprintf("Could not parse %s response: %s", ResourceName, err.Error()))
		return
	}
	parsed, err := node.GetObjectValue(models.CreateCloudPcProvisioningPolicyFromDiscriminatorValue)
	if err != nil {
		resp.Diagnostics.AddError("Error parsing resource", fmt.Sprintf("Could not decode %s response: %s", ResourceName, err.Error()))
		return
	}
	provisioningPolicy, ok := parsed.(models.CloudPcProvisioningPolicyable)
	if !ok || provisioningPolicy == nil {
		resp.Diagnostics.AddError("Error parsing resource", fmt.Sprintf("Unexpected %s response type %T", ResourceName, parsed))
		return
	}

	var rawRegionGroups struct {
		DomainJoinConfigurations []struct {
			RegionGroup string `json:"regionGroup"`
		} `json:"domainJoinConfigurations"`
	}
	if err := json.Unmarshal(body, &rawRegionGroups); err != nil {
		resp.Diagnostics.AddError("Error parsing resource", fmt.Sprintf("Could not decode %s region groups: %s", ResourceName, err.Error()))
		return
	}
	for i, config := range provisioningPolicy.GetDomainJoinConfigurations() {
		if config.GetRegionGroup() == nil && i < len(rawRegionGroups.DomainJoinConfigurations) && rawRegionGroups.DomainJoinConfigurations[i].RegionGroup != "" {
			config.GetAdditionalData()["regionGroup"] = rawRegionGroups.DomainJoinConfigurations[i].RegionGroup
		}
	}

	MapRemoteStateToTerraform(ctx, &object, provisioningPolicy)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for Cloud PC provisioning policy resources.
//
// Operation: Updates a Cloud PC provisioning policy with assignments and optional apply actions
// API Calls:
//   - PATCH /deviceManagement/virtualEndpoint/provisioningPolicies/{id}
//   - POST /deviceManagement/virtualEndpoint/provisioningPolicies/{id}/assign
//   - POST /deviceManagement/virtualEndpoint/provisioningPolicies/{id}/apply (if apply_to_existing_cloud_pcs is configured)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/cloudpcprovisioningpolicy-update?view=graph-rest-beta
func (r *CloudPcProvisioningPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan CloudPcProvisioningPolicyResourceModel
	var state CloudPcProvisioningPolicyResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update method for: %s", ResourceName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Updating %s with ID: %s", ResourceName, state.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, &plan, r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for update method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	updated, err := r.client.
		DeviceManagement().
		VirtualEndpoint().
		ProvisioningPolicies().
		ByCloudPcProvisioningPolicyId(state.ID.ValueString()).
		Patch(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	plannedAssignments := plan.Assignments
	MapRemoteStateToTerraform(ctx, &plan, updated)
	plan.Assignments = plannedAssignments

	tflog.Debug(ctx, fmt.Sprintf("Updating assignments for policy ID: %s", state.ID.ValueString()))

	assignBody, err := constructAssignmentsRequestBody(ctx, plan.Assignments)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing assignments request body",
			fmt.Sprintf("Could not construct assignments request body: %s", err.Error()),
		)
		return
	}

	err = r.client.
		DeviceManagement().
		VirtualEndpoint().
		ProvisioningPolicies().
		ByCloudPcProvisioningPolicyId(state.ID.ValueString()).
		Assign().
		Post(ctx, assignBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, "Successfully updated assignments")

	if plan.ApplyToExistingCloudPcs != nil {
		if !plan.ApplyToExistingCloudPcs.MicrosoftEntraSingleSignOnForAllDevices.IsNull() &&
			plan.ApplyToExistingCloudPcs.MicrosoftEntraSingleSignOnForAllDevices.ValueBool() {

			tflog.Debug(ctx, fmt.Sprintf("Applying Microsoft Entra Single Sign-On to all devices for policy ID: %s", state.ID.ValueString()))

			applyBody := devicemanagement.NewVirtualEndpointProvisioningPoliciesItemApplyPostRequestBody()
			applyBody.GetAdditionalData()["policySettings"] = "singleSignOn"

			err = r.client.
				DeviceManagement().
				VirtualEndpoint().
				ProvisioningPolicies().
				ByCloudPcProvisioningPolicyId(state.ID.ValueString()).
				Apply().
				Post(ctx, applyBody, nil)

			if err != nil {
				errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
				return
			}

			tflog.Debug(ctx, "Successfully applied Microsoft Entra Single Sign-On to all devices")
		}

		if !plan.ApplyToExistingCloudPcs.RegionOrAzureNetworkConnectionForAllDevices.IsNull() &&
			plan.ApplyToExistingCloudPcs.RegionOrAzureNetworkConnectionForAllDevices.ValueBool() {

			tflog.Debug(ctx, fmt.Sprintf("Applying Region or Azure Network Connection to all devices for policy ID: %s", state.ID.ValueString()))

			applyBody := devicemanagement.NewVirtualEndpointProvisioningPoliciesItemApplyPostRequestBody()
			applyBody.GetAdditionalData()["policySettings"] = "region"

			err = r.client.
				DeviceManagement().
				VirtualEndpoint().
				ProvisioningPolicies().
				ByCloudPcProvisioningPolicyId(state.ID.ValueString()).
				Apply().
				Post(ctx, applyBody, nil)

			if err != nil {
				errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
				return
			}

			tflog.Debug(ctx, "Successfully applied Region or Azure Network Connection to all devices")
		}

		if !plan.ApplyToExistingCloudPcs.RegionOrAzureNetworkConnectionForSelectDevices.IsNull() &&
			plan.ApplyToExistingCloudPcs.RegionOrAzureNetworkConnectionForSelectDevices.ValueBool() {

			tflog.Debug(ctx, fmt.Sprintf("Applying Region or Azure Network Connection to selected devices for policy ID: %s", state.ID.ValueString()))

			applyBody := devicemanagement.NewVirtualEndpointProvisioningPoliciesItemApplyPostRequestBody()
			applyBody.GetAdditionalData()["policySettings"] = "region"

			// Set reservePercentage to 0 to indicate applying to selected devices only
			// This is based on the API behavior where 0 means "selected devices"
			reservePercentage := int32(0)
			applyBody.SetReservePercentage(&reservePercentage)

			err = r.client.
				DeviceManagement().
				VirtualEndpoint().
				ProvisioningPolicies().
				ByCloudPcProvisioningPolicyId(state.ID.ValueString()).
				Apply().
				Post(ctx, applyBody, nil)

			if err != nil {
				errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
				return
			}

			tflog.Debug(ctx, "Successfully applied Region or Azure Network Connection to selected devices")
		}
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

// Delete handles the Delete operation for Cloud PC provisioning policy resources.
//
// Operation: Deletes a Cloud PC provisioning policy after removing assignments
// API Calls:
//   - POST /deviceManagement/virtualEndpoint/provisioningPolicies/{id}/assign (removes assignments)
//   - DELETE /deviceManagement/virtualEndpoint/provisioningPolicies/{id}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/cloudpcprovisioningpolicy-delete?view=graph-rest-beta
// Note: All assignments must be removed before policy deletion
func (r *CloudPcProvisioningPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object CloudPcProvisioningPolicyResourceModel

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

	tflog.Debug(ctx, fmt.Sprintf("Removing all assignments for policy ID: %s before deletion", object.ID.ValueString()))

	// Call the assign endpoint with an empty assignments array to remove all assignments
	// else policy deletion fails with a 400 if there are assignments
	emptyAssignments := types.SetNull(CloudPcProvisioningPolicyAssignmentType())
	assignBody, err := constructAssignmentsRequestBody(ctx, emptyAssignments)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing empty assignments request body",
			fmt.Sprintf("Could not construct assignments request body: %s", err.Error()),
		)
		return
	}

	err = r.client.
		DeviceManagement().
		VirtualEndpoint().
		ProvisioningPolicies().
		ByCloudPcProvisioningPolicyId(object.ID.ValueString()).
		Assign().
		Post(ctx, assignBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, "Successfully removed all assignments")

	time.Sleep(2 * time.Second)

	err = r.client.
		DeviceManagement().
		VirtualEndpoint().
		ProvisioningPolicies().
		ByCloudPcProvisioningPolicyId(object.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
