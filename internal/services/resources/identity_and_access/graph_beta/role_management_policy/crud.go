package graphBetaRoleManagementPolicy

import (
	"context"
	stderrors "errors"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
)

// Create adopts an existing role management policy and applies the configured rules.
//
// Operation: Role management policies are created by Microsoft Entra for every role and scope and cannot be created
// or deleted. Create finds the policy, validates the configured rules against it and updates them.
// API Calls:
//   - GET /policies/roleManagementPolicyAssignments?$filter=scopeId eq '{scopeId}' and scopeType eq '{scopeType}' and roleDefinitionId eq '{id}'
//   - GET /policies/roleManagementPolicies/{id}?$expand=rules
//   - PATCH /policies/roleManagementPolicies/{id}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/unifiedrolemanagementpolicy-update?view=graph-rest-beta
func (r *RoleManagementPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object RoleManagementPolicyResourceModel

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

	requestBody, err := constructResource(ctx, r.client, &object)
	if err != nil {
		var graphErr *odataerrors.ODataError
		if stderrors.As(err, &graphErr) {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.ReadPermissions)
			return
		}
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.Policies().RoleManagementPolicies().ByUnifiedRoleManagementPolicyId(*requestBody.GetId()).Patch(ctx, requestBody, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	object.ID = types.StringValue(*requestBody.GetId())

	tflog.Debug(ctx, fmt.Sprintf("Successfully adopted %s with ID: %s", ResourceName, object.ID.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationCreate
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = roleManagementPolicyConsistencyPredicate(&object)

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

// Read handles the Read operation for role management policy resources.
//
// Operation: Retrieves a role management policy and its rules by ID
// API Calls:
//   - GET /policies/roleManagementPolicies/{id}?$expand=rules
//   - GET /policies/roleManagementPolicyAssignments?$filter=scopeId eq '{scopeId}' and scopeType eq '{scopeType}'
//     (only while role_definition_id is unknown, for example after import)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/unifiedrolemanagementpolicy-get?view=graph-rest-beta
func (r *RoleManagementPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object RoleManagementPolicyResourceModel
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

	importAll := false
	if req.Private != nil {
		value, diags := req.Private.GetKey(ctx, importAllRulesKey)
		resp.Diagnostics.Append(diags...)
		importAll = len(value) > 0
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

	remoteResource, err := getPolicyWithRules(ctx, r.client, object.ID.ValueString())
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	object = MapRemoteResourceStateToTerraform(ctx, object, remoteResource, importAll)

	// The policy does not expose its role definition. Look it up once, after import or when the policy was
	// identified by policy_id, and keep it in state afterwards.
	if object.RoleDefinitionId.IsNull() || object.RoleDefinitionId.IsUnknown() || object.RoleDefinitionId.ValueString() == "" {
		roleDefinitionId, err := lookupRoleDefinitionId(ctx, r.client, object.ID.ValueString(), object.ScopeId.ValueString(), object.ScopeType.ValueString())
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
			return
		}
		object.RoleDefinitionId = types.StringNull()
		if roleDefinitionId != "" {
			object.RoleDefinitionId = types.StringValue(roleDefinitionId)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if importAll && resp.Private != nil {
		resp.Diagnostics.Append(resp.Private.SetKey(ctx, importAllRulesKey, nil)...)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
}

// Update handles the Update operation for role management policy resources.
//
// Operation: Validates the configured rules against the policy and updates them. Rules that are not configured are
// not sent and stay unchanged.
// API Calls:
//   - GET /policies/roleManagementPolicies/{id}?$expand=rules
//   - PATCH /policies/roleManagementPolicies/{id}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/unifiedrolemanagementpolicy-update?view=graph-rest-beta
func (r *RoleManagementPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan RoleManagementPolicyResourceModel
	var state RoleManagementPolicyResourceModel

	tflog.Debug(ctx, fmt.Sprintf("Updating %s", ResourceName))

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

	requestBody, err := constructResource(ctx, r.client, &plan)
	if err != nil {
		var graphErr *odataerrors.ODataError
		if stderrors.As(err, &graphErr) {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.ReadPermissions)
			return
		}
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	_, err = r.client.Policies().RoleManagementPolicies().ByUnifiedRoleManagementPolicyId(state.ID.ValueString()).Patch(ctx, requestBody, nil)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
		return
	}

	// Read maps only the rules tracked in its prior state and keeps a configured duration spelling when Graph returns an
	// equivalent one. Reading against the plan lets it pick up rules added in this update and the new spellings; the
	// response state is still written only by Read.
	readReq := resource.ReadRequest{State: tfsdk.State{Schema: req.Plan.Schema, Raw: req.Plan.Raw}, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationUpdate
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = roleManagementPolicyConsistencyPredicate(&plan)

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

// Delete releases Terraform management without changing the role management policy.
//
// Operation: Role management policies cannot be deleted. Delete only removes the resource from Terraform state and
// leaves the policy rules as they are. To return a rule to its previous values, apply them before destroying.
// API Calls: none
func (r *RoleManagementPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s", ResourceName))

	tflog.Info(ctx, fmt.Sprintf("Removing %s from Terraform state; the role management policy is left unchanged", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
