package graphBetaPolicySet

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
	deviceappmanagement "github.com/microsoftgraph/msgraph-beta-sdk-go/deviceappmanagement"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func (r *PolicySetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object PolicySetResourceModel

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

	requestBody, err := constructResource(ctx, r.client, &object, resp, r.ReadPermissions)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	baseResource, err := r.client.
		DeviceAppManagement().
		PolicySets().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	object.ID = types.StringValue(*baseResource.GetId())

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

func (r *PolicySetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object PolicySetResourceModel

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

	resource, err := r.client.
		DeviceAppManagement().
		PolicySets().
		ByPolicySetId(object.ID.ValueString()).
		Get(ctx, &deviceappmanagement.PolicySetsPolicySetItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &deviceappmanagement.PolicySetsPolicySetItemRequestBuilderGetQueryParameters{
				Expand: []string{"assignments", "items"},
			},
		})

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, resource)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

func (r *PolicySetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan PolicySetResourceModel
	var state PolicySetResourceModel

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

	var err error

	// Step 1: PATCH base properties only if displayName or description changed
	if !plan.DisplayName.Equal(state.DisplayName) || !plan.Description.Equal(state.Description) {
		tflog.Debug(ctx, "Base properties changed, updating via PATCH")

		var patchBody graphmodels.PolicySetable
		patchBody, err = constructBasePatchRequest(ctx, &plan)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing base patch request",
				fmt.Sprintf("Could not construct base patch request: %s: %s", ResourceName, err.Error()),
			)
			return
		}

		_, err = r.client.
			DeviceAppManagement().
			PolicySets().
			ByPolicySetId(state.ID.ValueString()).
			Patch(ctx, patchBody, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}
	} else {
		tflog.Debug(ctx, "Base properties unchanged, skipping PATCH")
	}

	// Step 2: POST to /update for policy set items only if items changed
	if !plan.Items.Equal(state.Items) {
		tflog.Debug(ctx, "Items changed, updating via /update endpoint")

		var itemsUpdateRequest deviceappmanagement.PolicySetsItemUpdatePostRequestBodyable
		itemsUpdateRequest, err = constructItemsUpdateRequest(ctx, r.client, &state, &plan, resp, r.ReadPermissions)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing policy set items update request",
				fmt.Sprintf("Could not construct policy set items update request: %s: %s", ResourceName, err.Error()),
			)
			return
		}

		err = r.client.
			DeviceAppManagement().
			PolicySets().
			ByPolicySetId(state.ID.ValueString()).
			Update().
			Post(ctx, itemsUpdateRequest, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}
	} else {
		tflog.Debug(ctx, "Items unchanged, skipping items update")
	}

	// Step 3: POST to /update for assignments only if assignments changed
	if !plan.Assignments.Equal(state.Assignments) {
		tflog.Debug(ctx, "Assignments changed, updating via /update endpoint")

		var assignmentsUpdateRequest deviceappmanagement.PolicySetsItemUpdatePostRequestBodyable
		assignmentsUpdateRequest, err = constructAssignmentsUpdateRequest(ctx, &plan)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing assignments update request",
				fmt.Sprintf("Could not construct assignments update request: %s: %s", ResourceName, err.Error()),
			)
			return
		}

		err = r.client.
			DeviceAppManagement().
			PolicySets().
			ByPolicySetId(state.ID.ValueString()).
			Update().
			Post(ctx, assignmentsUpdateRequest, nil)

		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}
	} else {
		tflog.Debug(ctx, "Assignments unchanged, skipping assignments update")
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

func (r *PolicySetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object PolicySetResourceModel

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

	err := r.client.
		DeviceAppManagement().
		PolicySets().
		ByPolicySetId(object.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfTfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
