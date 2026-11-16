package graphBetaIdentityAndAccessLifecycleWorkflow

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/identitygovernance"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
)

// purgeTimeout bounds how long Delete waits for a deleted workflow to appear in deletedItems
// before purging it.
const purgeTimeout = 60 * time.Second

// Create handles the Create operation for lifecycle workflow resources.
//
// Operation: Creates a new lifecycle workflow with its execution conditions and tasks
// API Calls:
//   - POST /identityGovernance/lifecycleWorkflows/workflows
//
// Reference: https://learn.microsoft.com/en-us/graph/api/identitygovernance-lifecycleworkflowscontainer-post-workflows?view=graph-rest-beta
func (r *LifecycleWorkflowResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object LifecycleWorkflowResourceModel

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
		IdentityGovernance().
		LifecycleWorkflows().
		Workflows().
		Post(ctx, requestBody, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	if createdResource == nil || createdResource.GetId() == nil {
		resp.Diagnostics.AddError(
			"Error extracting resource ID",
			fmt.Sprintf("Could not extract ID from created resource: %s. The API may not return the full resource on creation.", ResourceName),
		)
		return
	}

	object.ID = types.StringValue(*createdResource.GetId())

	tflog.Debug(ctx, fmt.Sprintf("Successfully created %s with ID: %s", ResourceName, object.ID.ValueString()))

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

// Read handles the Read operation for lifecycle workflow resources.
//
// Operation: Retrieves a lifecycle workflow by ID. Tasks are expanded by default.
// API Calls:
//   - GET /identityGovernance/lifecycleWorkflows/workflows/{id}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/identitygovernance-workflow-get?view=graph-rest-beta
func (r *LifecycleWorkflowResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object LifecycleWorkflowResourceModel
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

	remoteResource, err := r.client.
		IdentityGovernance().
		LifecycleWorkflows().
		Workflows().
		ByWorkflowId(object.ID.ValueString()).
		Get(ctx, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	object = MapRemoteResourceStateToTerraform(ctx, object, remoteResource)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for lifecycle workflow resources.
//
// Operation: Updates an existing lifecycle workflow in place
// API Calls:
//   - PATCH /identityGovernance/lifecycleWorkflows/workflows/{id} (only display_name, description, is_enabled, is_scheduling_enabled changed)
//   - POST /identityGovernance/lifecycleWorkflows/workflows/{id}/createNewVersion (category, execution_conditions or tasks changed)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/identitygovernance-workflow-update?view=graph-rest-beta
// Reference: https://learn.microsoft.com/en-us/graph/api/identitygovernance-workflow-createnewversion?view=graph-rest-beta
// Note: Graph rejects PATCH bodies containing anything other than the four basic properties. createNewVersion
// keeps the workflow ID, increments version and replaces all properties, so a new version also carries basic property changes.
func (r *LifecycleWorkflowResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan LifecycleWorkflowResourceModel
	var state LifecycleWorkflowResourceModel

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

	kind, diags := planUpdateKind(ctx, req.Plan, req.State)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	prevModified := types.StringNull()
	prevVersion := types.Int64Null()

	switch kind {
	case updateNewVersion:
		workflow, err := constructResource(ctx, &plan)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing resource",
				fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
			)
			return
		}

		requestBody := identitygovernance.NewLifecycleWorkflowsWorkflowsItemMicrosoftGraphIdentityGovernanceCreateNewVersionCreateNewVersionPostRequestBody()
		requestBody.SetWorkflow(workflow)

		tflog.Debug(ctx, fmt.Sprintf("Creating new version of %s with ID: %s", ResourceName, state.ID.ValueString()))

		_, err = r.client.
			IdentityGovernance().
			LifecycleWorkflows().
			Workflows().
			ByWorkflowId(state.ID.ValueString()).
			MicrosoftGraphIdentityGovernanceCreateNewVersion().
			Post(ctx, requestBody, nil)
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}

		prevModified = state.LastModifiedDateTime
		prevVersion = state.Version
	case updateBasic:
		requestBody, err := constructBasicUpdate(ctx, &plan)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error constructing resource",
				fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
			)
			return
		}

		_, err = r.client.
			IdentityGovernance().
			LifecycleWorkflows().
			Workflows().
			ByWorkflowId(state.ID.ValueString()).
			Patch(ctx, requestBody, nil)
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}

		prevModified = state.LastModifiedDateTime
	default:
		tflog.Debug(ctx, "No Graph properties changed, skipping write")
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationUpdate
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = lifecycleWorkflowConsistencyPredicate(prevModified, prevVersion)

	err := crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after update",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished updating %s with ID: %s", ResourceName, state.ID.ValueString()))
}

// Delete handles the Delete operation for lifecycle workflow resources.
//
// Operation: Deletes a lifecycle workflow and permanently removes it from deleted items
// API Calls:
//   - DELETE /identityGovernance/lifecycleWorkflows/workflows/{id}
//   - DELETE /identityGovernance/lifecycleWorkflows/deletedItems/workflows/{id}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/identitygovernance-workflow-delete?view=graph-rest-beta
// Reference: https://learn.microsoft.com/en-us/graph/api/identitygovernance-deleteditemcontainer-delete?view=graph-rest-beta
// Note: A deleted workflow stays in deletedItems for 30 days and counts toward the tenant's workflow limit,
// so destroy also purges it. A workflow removed by Terraform cannot be restored.
func (r *LifecycleWorkflowResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object LifecycleWorkflowResourceModel

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

	workflowId := object.ID.ValueString()

	tflog.Info(ctx, fmt.Sprintf("Deleting %s with ID: %s", ResourceName, workflowId))

	err := r.client.
		IdentityGovernance().
		LifecycleWorkflows().
		Workflows().
		ByWorkflowId(workflowId).
		Delete(ctx, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Soft deleted %s with ID: %s, purging from deleted items", ResourceName, workflowId))

	purgeCtx, purgeCancel := context.WithTimeout(ctx, purgeTimeout)
	defer purgeCancel()

	err = crud.PollUntil(purgeCtx, 2*time.Second, func(ctx context.Context) (bool, error) {
		purgeErr := r.client.
			IdentityGovernance().
			LifecycleWorkflows().
			DeletedItems().
			Workflows().
			ByWorkflowId(workflowId).
			Delete(ctx, nil)
		if purgeErr == nil {
			return true, nil
		}
		if errors.GraphError(ctx, purgeErr).StatusCode == http.StatusNotFound {
			return false, fmt.Errorf("workflow %s is not yet in deleted items: %w", workflowId, purgeErr)
		}
		return false, &crud.FatalPollError{Err: purgeErr}
	})

	if err != nil {
		resp.Diagnostics.AddWarning(
			"Workflow deleted but not purged",
			fmt.Sprintf("%s %s was deleted but could not be permanently removed from /identityGovernance/lifecycleWorkflows/deletedItems/workflows: %s. "+
				"It is removed automatically after 30 days and counts toward the tenant workflow limit until then.", ResourceName, workflowId, err.Error()),
		)
	} else {
		tflog.Debug(ctx, fmt.Sprintf("Permanently deleted %s with ID: %s", ResourceName, workflowId))
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
