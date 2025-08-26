package graphBetaCloudPcProvisioningPolicy

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

// Create handles the Create operation.
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
		errors.HandleKiotaGraphError(ctx, err, resp, "Create", r.WritePermissions)
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
			errors.HandleKiotaGraphError(ctx, err, resp, "CreateAssignments", r.WritePermissions)
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

// Read handles the Read operation.
func (r *CloudPcProvisioningPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object CloudPcProvisioningPolicyResourceModel

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

	tflog.Debug(ctx, fmt.Sprintf("Reading %s with ID: %s (operation: %s)", ResourceName, object.ID.ValueString(), operation))

	ctx, cancel := crud.HandleTimeout(ctx, object.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	tflog.Debug(ctx, "Fetching provisioning policy with expanded assignments")

	provisioningPolicy, err := r.client.
		DeviceManagement().
		VirtualEndpoint().
		ProvisioningPolicies().
		ByCloudPcProvisioningPolicyId(object.ID.ValueString()).
		Get(ctx, &devicemanagement.VirtualEndpointProvisioningPoliciesCloudPcProvisioningPolicyItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &devicemanagement.VirtualEndpointProvisioningPoliciesCloudPcProvisioningPolicyItemRequestBuilderGetQueryParameters{
				Expand: []string{"assignments"},
				Select: []string{"*"},
			},
		})

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, provisioningPolicy)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation.
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
		errors.HandleKiotaGraphError(ctx, err, resp, "Update", r.WritePermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, &plan, updated)

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
		errors.HandleKiotaGraphError(ctx, err, resp, "UpdateAssignments", r.WritePermissions)
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
				errors.HandleKiotaGraphError(ctx, err, resp, "ApplySingleSignOn", r.WritePermissions)
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
				errors.HandleKiotaGraphError(ctx, err, resp, "ApplyRegion", r.WritePermissions)
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
				errors.HandleKiotaGraphError(ctx, err, resp, "ApplyRegionToSelected", r.WritePermissions)
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

	tflog.Debug(ctx, fmt.Sprintf("Finished updating %s with ID: %s", ResourceName, state.ID.ValueString()))
}

// Delete handles the Delete operation.
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
		errors.HandleKiotaGraphError(ctx, err, resp, "Delete", r.WritePermissions)
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
		errors.HandleKiotaGraphError(ctx, err, resp, "Delete", r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
