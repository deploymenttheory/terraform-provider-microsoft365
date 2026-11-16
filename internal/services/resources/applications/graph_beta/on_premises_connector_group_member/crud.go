package graphBetaApplicationsOnPremisesConnectorGroupMember

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
)

// Create handles the Create operation for connector group member resources.
//
// Operation: Assigns a connector to a connector group
// API Calls:
//   - POST /onPremisesPublishingProfiles/applicationProxy/connectors/{id}/memberOf/$ref
//
// Reference: https://learn.microsoft.com/en-us/graph/api/connector-post-memberof?view=graph-rest-beta
func (r *OnPremisesConnectorGroupMemberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object OnPremisesConnectorGroupMemberResourceModel

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

	requestBody, err := constructResource(ctx, object.ConnectorGroupID.ValueString(), r.client.GetAdapter().GetBaseUrl())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	err = r.client.
		OnPremisesPublishingProfiles().
		ByOnPremisesPublishingProfileId(onPremisesPublishingProfileID).
		Connectors().
		ByConnectorId(object.ConnectorID.ValueString()).
		MemberOf().
		Ref().
		Post(ctx, requestBody, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	object.ID = types.StringValue(object.ConnectorID.ValueString())

	tflog.Debug(ctx, fmt.Sprintf("Successfully assigned connector %s to connector group %s", object.ConnectorID.ValueString(), object.ConnectorGroupID.ValueString()))

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

// Read handles the Read operation for connector group member resources.
//
// Operation: Retrieves the connector and the connector group it belongs to
// API Calls:
//   - GET /onPremisesPublishingProfiles/applicationProxy/connectors/{id}
//   - GET /onPremisesPublishingProfiles/applicationProxy/connectors/{id}/memberOf
//
// Reference: https://learn.microsoft.com/en-us/graph/api/connector-get?view=graph-rest-beta
// Reference: https://learn.microsoft.com/en-us/graph/api/connector-list-memberof?view=graph-rest-beta
func (r *OnPremisesConnectorGroupMemberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object OnPremisesConnectorGroupMemberResourceModel
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

	connector, err := r.client.
		OnPremisesPublishingProfiles().
		ByOnPremisesPublishingProfileId(onPremisesPublishingProfileID).
		Connectors().
		ByConnectorId(object.ConnectorID.ValueString()).
		Get(ctx, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	memberships, err := r.listConnectorMemberships(ctx, object.ConnectorID.ValueString())
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	connectorGroup := firstConnectorGroup(memberships)

	if connectorGroup == nil {
		if operation == constants.TfOperationRead {
			tflog.Warn(ctx, fmt.Sprintf("Connector %s is not a member of any connector group, removing %s from state", object.ConnectorID.ValueString(), ResourceName))
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Connector group membership not yet visible",
			fmt.Sprintf("Connector %s does not yet report a connector group, retry may be needed", object.ConnectorID.ValueString()),
		)
		return
	}

	// After create, Graph can briefly report the previous connector group.
	// Returning an error lets ReadWithRetry try again. Update waits for the new
	// group with a consistency predicate instead, because its read starts from
	// the prior state. During a normal refresh a different group is real drift,
	// so it is mapped into state and Terraform plans to move the connector back.
	currentGroupID := connectorGroupIDOf(connectorGroup)
	expectedGroupID := object.ConnectorGroupID.ValueString()
	if operation == constants.TfOperationCreate && !sameID(currentGroupID, expectedGroupID) {
		resp.Diagnostics.AddError(
			"Connector group membership not yet updated",
			fmt.Sprintf("Connector %s still reports connector group %s instead of %s", object.ConnectorID.ValueString(), currentGroupID, expectedGroupID),
		)
		return
	}

	MapRemoteResourceStateToTerraform(ctx, &object, connector, connectorGroup)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for connector group member resources.
//
// Operation: Moves the connector to a different connector group
// API Calls:
//   - POST /onPremisesPublishingProfiles/applicationProxy/connectors/{id}/memberOf/$ref
//
// Reference: https://learn.microsoft.com/en-us/graph/api/connector-post-memberof?view=graph-rest-beta
func (r *OnPremisesConnectorGroupMemberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan OnPremisesConnectorGroupMemberResourceModel
	var state OnPremisesConnectorGroupMemberResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update of resource: %s", ResourceName))

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

	requestBody, err := constructResource(ctx, plan.ConnectorGroupID.ValueString(), r.client.GetAdapter().GetBaseUrl())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	// Posting a new membership moves the connector; a connector belongs to
	// exactly one connector group, so no separate removal is needed.
	err = r.client.
		OnPremisesPublishingProfiles().
		ByOnPremisesPublishingProfileId(onPremisesPublishingProfileID).
		Connectors().
		ByConnectorId(state.ConnectorID.ValueString()).
		MemberOf().
		Ref().
		Post(ctx, requestBody, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationUpdate
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = connectorGroupConsistencyPredicate(plan.ConnectorGroupID.ValueString())

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

// Delete handles the Delete operation for connector group member resources.
//
// Operation: Moves the connector back to the tenant's default connector group
// API Calls:
//   - GET /onPremisesPublishingProfiles/applicationProxy/connectors/{id}/memberOf
//   - GET /onPremisesPublishingProfiles/applicationProxy/connectorGroups
//   - POST /onPremisesPublishingProfiles/applicationProxy/connectors/{id}/memberOf/$ref
//
// Reference: https://learn.microsoft.com/en-us/graph/api/connector-post-memberof?view=graph-rest-beta
// Note: A connector must always belong to a connector group, and Microsoft Learn
// documents no removal call. Delete therefore assigns the connector back to the
// default connector group, matching how Microsoft Entra places new connectors.
// If the connector has already moved to another group outside Terraform, or the
// managed group is the default group, Delete only removes the resource from state.
func (r *OnPremisesConnectorGroupMemberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object OnPremisesConnectorGroupMemberResourceModel

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

	memberships, err := r.listConnectorMemberships(ctx, object.ConnectorID.ValueString())
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.ReadPermissions)
		return
	}

	currentGroup := firstConnectorGroup(memberships)

	if !sameID(connectorGroupIDOf(currentGroup), object.ConnectorGroupID.ValueString()) {
		tflog.Warn(ctx, "Connector is no longer in the managed connector group; removing from state without moving it", map[string]any{
			"connector_id":               object.ConnectorID.ValueString(),
			"state_connector_group_id":   object.ConnectorGroupID.ValueString(),
			"current_connector_group_id": connectorGroupIDOf(currentGroup),
		})
		resp.State.RemoveResource(ctx)
		return
	}

	connectorGroups, err := r.listConnectorGroups(ctx)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.ReadPermissions)
		return
	}

	defaultGroup := selectDefaultConnectorGroup(connectorGroups)

	if defaultGroup == nil {
		resp.Diagnostics.AddError(
			"Default connector group not found",
			fmt.Sprintf("No default connector group to move connector %s back to. Move it manually or remove it from state.", object.ConnectorID.ValueString()),
		)
		return
	}

	if sameID(connectorGroupIDOf(defaultGroup), object.ConnectorGroupID.ValueString()) {
		tflog.Info(ctx, fmt.Sprintf("Connector %s is already in the default connector group; removing %s from state", object.ConnectorID.ValueString(), ResourceName))
		resp.State.RemoveResource(ctx)
		return
	}

	requestBody, err := constructResource(ctx, connectorGroupIDOf(defaultGroup), r.client.GetAdapter().GetBaseUrl())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	err = r.client.
		OnPremisesPublishingProfiles().
		ByOnPremisesPublishingProfileId(onPremisesPublishingProfileID).
		Connectors().
		ByConnectorId(object.ConnectorID.ValueString()).
		MemberOf().
		Ref().
		Post(ctx, requestBody, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Moved connector %s back to default connector group %s", object.ConnectorID.ValueString(), connectorGroupIDOf(defaultGroup)))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
