package graphBetaAuthorizationPolicy

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

// readbackAttributes are the attributes readbackMatch compares.
var readbackAttributes = []string{
	"id", "allowed_to_sign_up_email_based_subscriptions", "allowed_to_use_sspr", "allow_email_verified_users_to_join_organization", "allow_invites_from",
	"allow_user_consent_for_risky_apps", "block_msol_powershell", "description", "display_name", "enabled_preview_features", "guest_user_role_id",
	"permission_grant_policy_ids_assigned_to_default_user_role", "default_user_role_permissions",
}

// readbackMatch is what the read-back after a write waits for: the singleton id and the written settings.
var readbackMatch = crud.AttributeMatchOptions{
	Paths:              crud.RootPaths(readbackAttributes...),
	IgnoreNullExpected: true,
}

// Create configures the authorization policy.
//
// The authorizationPolicy resource is a singleton that always exists in the tenant; there is no documented
// POST endpoint to create it. This operation therefore uses a PATCH request to apply the desired
// configuration, mirroring the Update behaviour.
//
// API Calls:
//   - PATCH /policies/authorizationPolicy/authorizationPolicy
//
// Reference: https://learn.microsoft.com/en-us/graph/api/authorizationpolicy-update?view=graph-rest-beta
func (r *AuthorizationPolicyResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var object AuthorizationPolicyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(
		ctx,
		object.Timeouts.Create,
		CreateTimeout*time.Second,
		&resp.Diagnostics,
	)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, &object)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Create Method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	// The permanent singleton can briefly return 404 while a preceding write replicates.
	writeOpts := crud.PreconditionRetryOptions{
		MaxRetries: 15, RetryInterval: 2 * time.Second, StatusCodes: []int{404},
		ResourceTypeName: ResourceName, Operation: constants.TfOperationCreate,
	}
	err = crud.WriteWithPreconditionRetry(ctx, func(ctx context.Context) error {
		_, patchErr := r.client.Policies().AuthorizationPolicy().ByAuthorizationPolicyId(singletonID).Patch(ctx, requestBody, nil)
		if patchErr != nil {
			return fmt.Errorf("patch authorization policy: %w", patchErr)
		}
		return nil
	}, writeOpts)
	if err != nil {
		errors.HandleKiotaGraphError(
			ctx,
			err,
			resp,
			constants.TfOperationCreate,
			r.WritePermissions,
		)
		return
	}

	// This is a singleton resource — set the static ID immediately so the subsequent Read
	// can locate the correct state entry.
	object.ID = types.StringValue(singletonID)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Graph can queue role-permission changes and overwrite invitations afterwards.
	// Confirm those settings first, then apply the configured invitations separately.
	if requestBody.GetAllowInvitesFrom() != nil && requestBody.GetDefaultUserRolePermissions() != nil {
		names := make([]string, 0, len(readbackAttributes))
		for _, name := range readbackAttributes {
			if name != "allow_invites_from" {
				names = append(names, name)
			}
		}
		settingsMatch := readbackMatch
		settingsMatch.Paths = crud.RootPaths(names...)
		settingsOpts := crud.NewReadWithRetryOptions(constants.TfOperationCreate, ResourceName,
			crud.ReadWithConsistencyPredicate(crud.RequireStableReads(crud.AttributesMatchPredicate(resp.State, settingsMatch))))
		if deadline, ok := ctx.Deadline(); ok {
			settingsOpts.MaxRetries = int(time.Until(deadline) / settingsOpts.RetryInterval)
		}
		// Keep intermediate invitation values out of the desired final state.
		pending := &resource.CreateResponse{State: resp.State}
		err = crud.ReadWithRetry(ctx, r.Read,
			resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta},
			&crud.CreateResponseContainer{CreateResponse: pending}, settingsOpts)
		if err != nil {
			resp.Diagnostics.AddError("Error confirming authorization settings before invitations", err.Error())
			return
		}
		invitations := graphmodels.NewAuthorizationPolicy()
		invitations.SetAllowInvitesFrom(requestBody.GetAllowInvitesFrom())
		err = crud.WriteWithPreconditionRetry(ctx, func(ctx context.Context) error {
			_, patchErr := r.client.Policies().AuthorizationPolicy().ByAuthorizationPolicyId(singletonID).Patch(ctx, invitations, nil)
			if patchErr != nil {
				return fmt.Errorf("patch authorization invitations: %w", patchErr)
			}
			return nil
		}, writeOpts)
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
			return
		}
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationCreate
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = crud.RequireStableReads(crud.AttributesMatchPredicate(resp.State, readbackMatch))
	if deadline, ok := ctx.Deadline(); ok {
		opts.MaxRetries = int(time.Until(deadline) / opts.RetryInterval)
	}

	err = crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after create",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}
}

// Read retrieves the current state of the authorization policy from the Graph API.
//
// API Calls:
//   - GET /policies/authorizationPolicy
//
// Reference: https://learn.microsoft.com/en-us/graph/api/authorizationpolicy-get?view=graph-rest-beta
func (r *AuthorizationPolicyResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var object AuthorizationPolicyResourceModel

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

	ctx, cancel := crud.HandleTimeout(
		ctx,
		object.Timeouts.Read,
		ReadTimeout*time.Second,
		&resp.Diagnostics,
	)
	if cancel == nil {
		return
	}
	defer cancel()

	remoteResource, err := r.client.
		Policies().
		AuthorizationPolicy().
		Get(ctx, nil)
	if err != nil {
		// A singleton cannot be deleted by this resource. Keep state on unavailable or malformed reads.
		diagnosticResponse := &resource.UpdateResponse{}
		errors.HandleKiotaGraphError(ctx, err, diagnosticResponse, operation, r.ReadPermissions)
		resp.Diagnostics.Append(diagnosticResponse.Diagnostics...)
		return
	}

	if err := mapRemoteState(ctx, &object, remoteResource); err != nil {
		resp.Diagnostics.AddError("Invalid authorization policy response", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update applies changes to the authorization policy.
//
// API Calls:
//   - PATCH /policies/authorizationPolicy/authorizationPolicy
//
// Reference: https://learn.microsoft.com/en-us/graph/api/authorizationpolicy-update?view=graph-rest-beta
func (r *AuthorizationPolicyResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var plan AuthorizationPolicyResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ctx, cancel := crud.HandleTimeout(
		ctx,
		plan.Timeouts.Update,
		UpdateTimeout*time.Second,
		&resp.Diagnostics,
	)
	if cancel == nil {
		return
	}
	defer cancel()

	requestBody, err := constructResource(ctx, &plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource for Update Method",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	// The permanent singleton can briefly return 404 while a preceding write replicates.
	writeOpts := crud.PreconditionRetryOptions{
		MaxRetries: 15, RetryInterval: 2 * time.Second, StatusCodes: []int{404},
		ResourceTypeName: ResourceName, Operation: constants.TfOperationUpdate,
	}
	err = crud.WriteWithPreconditionRetry(ctx, func(ctx context.Context) error {
		_, patchErr := r.client.Policies().AuthorizationPolicy().ByAuthorizationPolicyId(singletonID).Patch(ctx, requestBody, nil)
		if patchErr != nil {
			return fmt.Errorf("patch authorization policy: %w", patchErr)
		}
		return nil
	}, writeOpts)
	if err != nil {
		errors.HandleKiotaGraphError(
			ctx,
			err,
			resp,
			constants.TfOperationUpdate,
			r.WritePermissions,
		)
		return
	}

	plan.ID = types.StringValue(singletonID)
	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Graph can queue role-permission changes and overwrite invitations afterwards.
	// Confirm those settings first, then apply the configured invitations separately.
	if requestBody.GetAllowInvitesFrom() != nil && requestBody.GetDefaultUserRolePermissions() != nil {
		names := make([]string, 0, len(readbackAttributes))
		for _, name := range readbackAttributes {
			if name != "allow_invites_from" {
				names = append(names, name)
			}
		}
		settingsMatch := readbackMatch
		settingsMatch.Paths = crud.RootPaths(names...)
		settingsOpts := crud.NewReadWithRetryOptions(constants.TfOperationUpdate, ResourceName,
			crud.ReadWithConsistencyPredicate(crud.RequireStableReads(crud.AttributesMatchPredicate(resp.State, settingsMatch))))
		if deadline, ok := ctx.Deadline(); ok {
			settingsOpts.MaxRetries = int(time.Until(deadline) / settingsOpts.RetryInterval)
		}
		// Keep intermediate invitation values out of the desired final state.
		pending := &resource.UpdateResponse{State: resp.State}
		err = crud.ReadWithRetry(ctx, r.Read,
			resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta},
			&crud.UpdateResponseContainer{UpdateResponse: pending}, settingsOpts)
		if err != nil {
			resp.Diagnostics.AddError("Error confirming authorization settings before invitations", err.Error())
			return
		}
		invitations := graphmodels.NewAuthorizationPolicy()
		invitations.SetAllowInvitesFrom(requestBody.GetAllowInvitesFrom())
		err = crud.WriteWithPreconditionRetry(ctx, func(ctx context.Context) error {
			_, patchErr := r.client.Policies().AuthorizationPolicy().ByAuthorizationPolicyId(singletonID).Patch(ctx, invitations, nil)
			if patchErr != nil {
				return fmt.Errorf("patch authorization invitations: %w", patchErr)
			}
			return nil
		}, writeOpts)
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
			return
		}
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.UpdateResponseContainer{UpdateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationUpdate
	opts.ResourceTypeName = ResourceName
	opts.ConsistencyPredicate = crud.RequireStableReads(crud.AttributesMatchPredicate(resp.State, readbackMatch))
	if deadline, ok := ctx.Deadline(); ok {
		opts.MaxRetries = int(time.Until(deadline) / opts.RetryInterval)
	}

	err = crud.ReadWithRetry(ctx, r.Read, readReq, stateContainer, opts)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading resource state after update",
			fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()),
		)
		return
	}
}

// Delete releases Terraform management without changing the tenant-wide settings.
func (r *AuthorizationPolicyResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	resp.State.RemoveResource(ctx)
}
