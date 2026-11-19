package graphBetaDirectoryRoleAssignment

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	grapherrors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
)

var errEmptyAssignment = errors.New("empty directory role assignment response during deletion verification")

// Create handles the Create operation for directory role assignment resources.
//
// Operation: Creates a permanent active tenant-scoped directory role assignment
// API Calls:
//   - POST /roleManagement/directory/roleAssignments
//
// Reference: https://learn.microsoft.com/en-us/graph/api/resources/unifiedroleassignment?view=graph-rest-beta
func (r *DirectoryRoleAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DirectoryRoleAssignmentResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Create, CreateTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	body := constructResource(&data)
	var remote graphmodels.UnifiedRoleAssignmentable
	err := crud.PollUntil(ctx, 2*time.Second, func(ctx context.Context) (bool, error) {
		var postErr error
		remote, postErr = r.client.RoleManagement().Directory().RoleAssignments().Post(ctx, body, nil)
		if postErr == nil {
			return true, nil
		}
		info := grapherrors.GraphError(ctx, postErr)
		if info.StatusCode == 404 ||
			(info.StatusCode == 400 && strings.Contains(strings.ToLower(info.ErrorMessage), "principal") && strings.Contains(strings.ToLower(info.ErrorMessage), "not found")) {
			return false, fmt.Errorf("waiting for directory role assignment references: %w", postErr)
		}
		return false, &crud.FatalPollError{Err: postErr}
	})
	if err != nil {
		grapherrors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}
	if remote == nil || remote.GetId() == nil || *remote.GetId() == "" {
		resp.Diagnostics.AddError("Invalid directory role assignment response", "Microsoft Graph did not return an assignment ID.")
		return
	}
	data.ID = types.StringValue(*remote.GetId())
	data.DirectoryScopeID = types.StringValue("/")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationCreate
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = func(ctx context.Context, state tfsdk.State) bool {
		var actual DirectoryRoleAssignmentResourceModel
		if state.Get(ctx, &actual).HasError() {
			return false
		}
		return actual.PrincipalID.Equal(data.PrincipalID) && actual.RoleDefinitionID.Equal(data.RoleDefinitionID) && actual.DirectoryScopeID.ValueString() == "/"
	}
	if err = crud.ReadWithRetry(
		ctx,
		r.Read,
		resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta},
		&crud.CreateResponseContainer{CreateResponse: resp},
		opts,
	); err != nil {
		resp.Diagnostics.AddError("Error reading directory role assignment after create", err.Error())
	}
}

// Read handles the Read operation for directory role assignment resources.
//
// Operation: Retrieves the directory role assignment by ID
// API Calls:
//   - GET /roleManagement/directory/roleAssignments/{id}
//
// Reference: https://learn.microsoft.com/en-us/graph/api/resources/unifiedroleassignment?view=graph-rest-beta
func (r *DirectoryRoleAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DirectoryRoleAssignmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Read, ReadTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	remote, err := r.client.RoleManagement().Directory().RoleAssignments().ByUnifiedRoleAssignmentId(data.ID.ValueString()).Get(ctx, nil)
	if err != nil {
		operation := constants.TfOperationRead
		if op, ok := ctx.Value("retry_operation").(string); ok {
			operation = op
		}
		grapherrors.HandleKiotaGraphErrorWithOptions(ctx, err, resp, operation, r.ReadPermissions, grapherrors.GraphErrorOptions{PreserveStateOnReadBadRequest: true})
		return
	}
	if remote == nil || remote.GetId() == nil {
		resp.Diagnostics.AddError("Invalid directory role assignment response", "Microsoft Graph did not return an assignment.")
		return
	}
	if remote.GetDirectoryScopeId() == nil || *remote.GetDirectoryScopeId() != "/" || (remote.GetAppScopeId() != nil && *remote.GetAppScopeId() != "") {
		resp.Diagnostics.AddError(
			"Unsupported directory role assignment scope",
			"This resource manages tenant-wide directory assignments only (directoryScopeId '/'). Import a tenant-scoped assignment.",
		)
		return
	}
	MapRemoteStateToTerraform(&data, remote)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, sharedmodels.ResourceIdentity{ID: data.ID.ValueString()})...)
	}
}

// Update handles the Update operation for directory role assignment resources.
//
// Operation: Rejects in-place updates because role and principal changes require replacement
// API Calls:
//   - None (replacement required)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/resources/unifiedroleassignment?view=graph-rest-beta
func (r *DirectoryRoleAssignmentResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Unexpected directory role assignment update", "Assignment identity changes require replacement.")
}

// Delete handles the Delete operation for directory role assignment resources.
//
// Operation: Deletes the directory role assignment
// API Calls:
//   - DELETE /roleManagement/directory/roleAssignments/{id}
//   - GET /roleManagement/directory/roleAssignments/{id} (verifies deletion)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/resources/unifiedroleassignment?view=graph-rest-beta
func (r *DirectoryRoleAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data DirectoryRoleAssignmentResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}
	ctx, cancel := crud.HandleTimeout(ctx, data.Timeouts.Delete, DeleteTimeout*time.Second, &resp.Diagnostics)
	if cancel == nil {
		return
	}
	defer cancel()
	missing := 0
	err := crud.PollUntil(ctx, 2*time.Second, func(ctx context.Context) (bool, error) {
		deleteErr := r.client.RoleManagement().Directory().RoleAssignments().ByUnifiedRoleAssignmentId(data.ID.ValueString()).Delete(ctx, nil)
		if deleteErr != nil && grapherrors.GraphError(ctx, deleteErr).StatusCode != 404 {
			return false, &crud.FatalPollError{Err: deleteErr}
		}
		remote, getErr := r.client.RoleManagement().Directory().RoleAssignments().ByUnifiedRoleAssignmentId(data.ID.ValueString()).Get(ctx, nil)
		if getErr != nil {
			if grapherrors.GraphError(ctx, getErr).StatusCode == 404 {
				missing++
				return missing >= 3, nil
			}
			return false, &crud.FatalPollError{Err: getErr}
		}
		missing = 0
		if remote == nil {
			return false, &crud.FatalPollError{Err: errEmptyAssignment}
		}
		return false, nil
	})
	if ctx.Err() != nil {
		err = ctx.Err()
	}
	if err != nil {
		grapherrors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}
	resp.State.RemoveResource(ctx)
}
