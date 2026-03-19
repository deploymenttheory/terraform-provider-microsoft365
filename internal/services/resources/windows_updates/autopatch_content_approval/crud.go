package graphBetaWindowsUpdatesAutopatchContentApproval

import (
	"context"
	"fmt"
	"time"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	admin "github.com/microsoftgraph/msgraph-beta-sdk-go/admin"
)

// Create handles the Create operation for Windows Update content approval resources.
//
// Operation: Creates a new compliance change (content approval) for an update policy
// API Calls:
//   - POST /admin/windows/updates/updatePolicies/{updatePolicyId}/complianceChanges
//
// Reference: https://learn.microsoft.com/en-us/graph/api/windowsupdates-updatepolicy-post-compliancechanges-contentapproval?view=graph-rest-beta&tabs=http
func (r *WindowsUpdatesAutopatchContentApprovalResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object WindowsUpdatesAutopatchContentApprovalResourceModel

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

	createdResource, err := r.client.
		Admin().
		Windows().
		Updates().
		UpdatePolicies().
		ByUpdatePolicyId(object.UpdatePolicyId.ValueString()).
		ComplianceChanges().
		Post(ctx, requestBody, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	object.ID = convert.GraphToFrameworkString(createdResource.GetId())

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

// Read handles the Read operation for Windows Update content approval resources.
//
// Operation: Retrieves a compliance change (content approval) by ID
// API Calls:
//   - GET /admin/windows/updates/updatePolicies/{updatePolicyId}/complianceChanges/{complianceChangeId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/windowsupdates-contentapproval-get?view=graph-rest-beta&tabs=http
func (r *WindowsUpdatesAutopatchContentApprovalResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object WindowsUpdatesAutopatchContentApprovalResourceModel
	var identity sharedmodels.ResourceIdentity

	tflog.Debug(ctx, fmt.Sprintf("Starting Read of resource: %s", ResourceName))

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

	remoteResource, err := r.client.
		Admin().
		Windows().
		Updates().
		UpdatePolicies().
		ByUpdatePolicyId(object.UpdatePolicyId.ValueString()).
		ComplianceChanges().
		ByComplianceChangeId(object.ID.ValueString()).
		Get(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	MapRemoteStateToTerraform(ctx, &object, remoteResource)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	identity.ID = object.ID.ValueString()

	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, identity)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for Windows Update content approval resources.
//
// Operation: Updates an existing compliance change (content approval)
// API Calls:
//   - PATCH /admin/windows/updates/updatePolicies/{updatePolicyId}/complianceChanges/{complianceChangeId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/windowsupdates-contentapproval-update?view=graph-rest-beta&tabs=http
//
// Note: The API requires that PATCH requests contain exactly 1 changed property for compliance changes.
func (r *WindowsUpdatesAutopatchContentApprovalResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan WindowsUpdatesAutopatchContentApprovalResourceModel
	var state WindowsUpdatesAutopatchContentApprovalResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Starting Update of resource: %s", ResourceName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Updating %s with ID: %s", ResourceName, state.ID.ValueString()))

	ctx, cancel := crud.HandleTimeout(ctx, plan.Timeouts.Update, UpdateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructUpdateResource(ctx, &plan, &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing update request",
			fmt.Sprintf("Could not construct update request: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.
		Admin().
		Windows().
		Updates().
		UpdatePolicies().
		ByUpdatePolicyId(state.UpdatePolicyId.ValueString()).
		ComplianceChanges().
		ByComplianceChangeId(state.ID.ValueString()).
		Patch(ctx, requestBody, nil)

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

// Delete handles the Delete operation for Windows Update content approval resources.
//
// Creating a contentApproval causes the Graph API to implicitly create a deployment
// bound to the update policy's audience. Deleting the compliance change does not
// remove that deployment, which prevents the audience from being deleted later.
// This method therefore:
//  1. Deletes the compliance change itself.
//  2. Resolves the audience ID from the update policy.
//  3. Lists all deployments whose audience matches and deletes each one.
//
// API Calls:
//   - DELETE /admin/windows/updates/updatePolicies/{updatePolicyId}/complianceChanges/{complianceChangeId}
//   - GET    /admin/windows/updates/updatePolicies/{updatePolicyId}
//   - GET    /admin/windows/updates/deployments?$filter=audience/id eq '{audienceId}'
//   - DELETE /admin/windows/updates/deployments/{deploymentId}  (repeated per deployment)
// Delete handles the Delete operation for Windows Update content approval resources.
//
// Operation: Deletes a compliance change (content approval)
// API Calls:
//   - DELETE /admin/windows/updates/updatePolicies/{updatePolicyId}/complianceChanges/{complianceChangeId}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/windowsupdates-compliancechange-delete?view=graph-rest-beta
//
// Note: Includes retry logic to handle cases where the audience is still being used by deployments
// during cleanup race conditions.
func (r *WindowsUpdatesAutopatchContentApprovalResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object WindowsUpdatesAutopatchContentApprovalResourceModel

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

	// Step 1: delete the compliance change.
	err := r.client.
		Admin().
		Windows().
		Updates().
		UpdatePolicies().
		ByUpdatePolicyId(object.UpdatePolicyId.ValueString()).
		ComplianceChanges().
		ByComplianceChangeId(object.ID.ValueString()).
		Delete(ctx, nil)

	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	// Step 2: resolve the audience ID from the update policy so we can find
	// the implicit deployment that was created when this compliance change was made.
	policy, err := r.client.
		Admin().
		Windows().
		Updates().
		UpdatePolicies().
		ByUpdatePolicyId(object.UpdatePolicyId.ValueString()).
		Get(ctx, nil)

	if err != nil {
		// Log but do not fail — the compliance change is already gone.
		tflog.Warn(ctx, fmt.Sprintf("Could not retrieve update policy to resolve audience ID for deployment cleanup: %s", err.Error()))
		resp.State.RemoveResource(ctx)
		return
	}

	audience := policy.GetAudience()
	if audience == nil || audience.GetId() == nil {
		tflog.Warn(ctx, "Update policy returned no audience; skipping implicit deployment cleanup")
		resp.State.RemoveResource(ctx)
		return
	}

	audienceId := *audience.GetId()

	// Step 3: list all deployments and delete any bound to this audience.
	// The Graph API does not support OData filtering on audience/id for this endpoint,
	// so we retrieve all deployments and filter in memory.
	deploymentsResp, err := r.client.
		Admin().
		Windows().
		Updates().
		Deployments().
		Get(ctx, &admin.WindowsUpdatesDeploymentsRequestBuilderGetRequestConfiguration{
			QueryParameters: &admin.WindowsUpdatesDeploymentsRequestBuilderGetQueryParameters{
				Select: []string{"id", "audience"},
			},
		})

	if err != nil {
		tflog.Warn(ctx, fmt.Sprintf("Could not list deployments for audience %s during cleanup: %s", audienceId, err.Error()))
		resp.State.RemoveResource(ctx)
		return
	}

	for _, deployment := range deploymentsResp.GetValue() {
		if deployment.GetId() == nil {
			continue
		}
		dAudience := deployment.GetAudience()
		if dAudience == nil || dAudience.GetId() == nil || *dAudience.GetId() != audienceId {
			continue
		}
		deploymentId := *deployment.GetId()
		tflog.Debug(ctx, fmt.Sprintf("Deleting implicit deployment %s created by content approval", deploymentId))
		delErr := r.client.
			Admin().
			Windows().
			Updates().
			Deployments().
			ByDeploymentId(deploymentId).
			Delete(ctx, nil)
		if delErr != nil {
			tflog.Warn(ctx, fmt.Sprintf("Failed to delete implicit deployment %s: %s", deploymentId, delErr.Error()))
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
