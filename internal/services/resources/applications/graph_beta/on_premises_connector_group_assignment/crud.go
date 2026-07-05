package graphBetaApplicationsOnPremisesConnectorGroupAssignment

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
)

func (r *OnPremisesConnectorGroupAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object OnPremisesConnectorGroupAssignmentResourceModel

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

	requestBody, err := constructResource(ctx, &object, r.client.GetAdapter().GetBaseUrl())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing connector group assignment",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	// Assignment is only valid for applications that already have Application
	// Proxy on-premises publishing enabled. Learn documents the $ref endpoint at
	// https://learn.microsoft.com/en-us/graph/api/connectorgroup-post-applications?view=graph-rest-beta,
	// but the prerequisite is easy to miss from this resource alone. Direct API
	// verification on 2026-07-05 returned Application_NotFound with
	// "Application '{id}' not found or OnPremisesPublishing is not enabled for
	// your tenant." for an existing application without onPremisesPublishing.
	err = r.client.
		Applications().
		ByApplicationId(object.ApplicationID.ValueString()).
		ConnectorGroup().
		Ref().
		Put(ctx, requestBody, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	object.ID = types.StringValue(compositeID(object.ApplicationID.ValueString(), object.ConnectorGroupID.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// The assignment endpoint is eventually consistent enough that an immediate
	// GET can briefly return the previous connector group. Reuse the provider's
	// normal create read-retry flow so the resource only persists once Graph
	// reports the requested connector group.
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

func (r *OnPremisesConnectorGroupAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object OnPremisesConnectorGroupAssignmentResourceModel
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

	connectorGroup, err := r.client.
		Applications().
		ByApplicationId(object.ApplicationID.ValueString()).
		ConnectorGroup().
		Get(ctx, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	if connectorGroup == nil || connectorGroup.GetId() == nil || *connectorGroup.GetId() != object.ConnectorGroupID.ValueString() {
		// Deleting /applications/{id}/connectorGroup/$ref does not make the
		// navigation property disappear in this tenant. Direct API verification
		// on 2026-07-05 showed GET /applications/{id}/connectorGroup returning
		// the tenant default connector group after DELETE. Treat any current
		// connector group ID different from connector_group_id as this managed
		// assignment being absent.
		if operation == constants.TfOperationRead {
			tflog.Warn(ctx, "Connector group assignment not found on application, removing from state", map[string]any{
				"application_id":     object.ApplicationID.ValueString(),
				"connector_group_id": object.ConnectorGroupID.ValueString(),
			})
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Connector group assignment not found",
			fmt.Sprintf("Connector group %s is not yet assigned to application %s, retry may be needed", object.ConnectorGroupID.ValueString(), object.ApplicationID.ValueString()),
		)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &object, connectorGroup)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

func (r *OnPremisesConnectorGroupAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan OnPremisesConnectorGroupAssignmentResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *OnPremisesConnectorGroupAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object OnPremisesConnectorGroupAssignmentResourceModel

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

	// Learn documents the assignment as a $ref relationship:
	// https://learn.microsoft.com/en-us/graph/api/connectorgroup-post-applications?view=graph-rest-beta
	// The generated SDK exposes Delete on that same /connectorGroup/$ref path.
	// Direct API verification on 2026-07-05 confirmed DELETE works, but Graph may
	// then expose the tenant default connector group on the application navigation
	// property. Read therefore uses ID comparison rather than requiring
	// GET /connectorGroup to return 404.
	err := r.client.
		Applications().
		ByApplicationId(object.ApplicationID.ValueString()).
		ConnectorGroup().
		Ref().
		Delete(ctx, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
