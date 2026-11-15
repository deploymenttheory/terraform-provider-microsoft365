package graphBetaIdentityAndAccessAuthenticationMethodsPolicy

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

// Create adopts the tenant authentication methods policy and applies the configuration.
//
// Operation: The authenticationMethodsPolicy is a singleton that always exists in the tenant and has no POST
// endpoint. Create PATCHes the configured policy-level settings and each configured method configuration.
// API Calls:
//   - PATCH /policies/authenticationMethodsPolicy (only when a policy-level setting is configured)
//   - PATCH /policies/authenticationMethodsPolicy/authenticationMethodConfigurations/{id} (once per configured method)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/authenticationmethodspolicy-update?view=graph-rest-beta
func (r *AuthenticationMethodsPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object AuthenticationMethodsPolicyResourceModel

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

	requestBodies, err := constructResource(ctx, &object, nil)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	err = r.patchAuthenticationMethodsPolicy(ctx, requestBodies)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	// This is a singleton resource: set the static ID so the subsequent Read can locate the state entry.
	object.ID = types.StringValue(singletonID)

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

// Read handles the Read operation for the authentication methods policy resource.
//
// Operation: Retrieves the tenant authentication methods policy, which includes its method configurations
// API Calls:
//   - GET /policies/authenticationMethodsPolicy
//
// Reference: https://learn.microsoft.com/en-us/graph/api/authenticationmethodspolicy-get?view=graph-rest-beta
func (r *AuthenticationMethodsPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object AuthenticationMethodsPolicyResourceModel
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

	identity.ID = singletonID

	if resp.Identity != nil {
		resp.Diagnostics.Append(resp.Identity.Set(ctx, identity)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	remoteResource, err := r.client.
		Policies().
		AuthenticationMethodsPolicy().
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

// Update handles the Update operation for the authentication methods policy resource.
//
// Operation: PATCHes the policy-level settings and each method configuration that changed
// API Calls:
//   - PATCH /policies/authenticationMethodsPolicy (only when a policy-level setting changed)
//   - PATCH /policies/authenticationMethodsPolicy/authenticationMethodConfigurations/{id} (once per changed method)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/authenticationmethodconfiguration-update?view=graph-rest-beta
func (r *AuthenticationMethodsPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan AuthenticationMethodsPolicyResourceModel
	var state AuthenticationMethodsPolicyResourceModel

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

	requestBodies, err := constructResource(ctx, &plan, &state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error constructing resource",
			fmt.Sprintf("Could not construct resource: %s: %s", ResourceName, err.Error()),
		)
		return
	}

	err = r.patchAuthenticationMethodsPolicy(ctx, requestBodies)
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationUpdate, r.WritePermissions)
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

	tflog.Debug(ctx, fmt.Sprintf("Finished updating %s with ID: %s", ResourceName, singletonID))
}

// Delete releases Terraform management without changing the tenant-wide authentication methods policy.
//
// Operation: The singleton policy cannot be deleted. Delete only removes the resource from Terraform
// state and leaves the tenant policy and its method configurations as they are.
// API Calls: none
func (r *AuthenticationMethodsPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, fmt.Sprintf("Starting deletion of resource: %s", ResourceName))

	tflog.Info(ctx, fmt.Sprintf("Removing %s from Terraform state; the tenant authentication methods policy is left unchanged", ResourceName))

	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}

// patchAuthenticationMethodsPolicy sends the policy-level PATCH, when there is one, and then PATCHes each
// changed method configuration at its own URI.
func (r *AuthenticationMethodsPolicyResource) patchAuthenticationMethodsPolicy(ctx context.Context, requestBodies *requestBodies) error {
	if requestBodies.Policy != nil {
		_, err := r.client.
			Policies().
			AuthenticationMethodsPolicy().
			Patch(ctx, requestBodies.Policy, nil)
		if err != nil {
			return err //nolint:wrapcheck // Keep the Kiota error type for shared Graph HTTP status handling.
		}
	}

	for _, method := range requestBodies.Methods {
		tflog.Debug(ctx, fmt.Sprintf("Updating %s method configuration %s", ResourceName, method.ID))

		_, err := r.client.
			Policies().
			AuthenticationMethodsPolicy().
			AuthenticationMethodConfigurations().
			ByAuthenticationMethodConfigurationId(method.ID).
			Patch(ctx, method.Body, nil)
		if err != nil {
			return err //nolint:wrapcheck // Keep the Kiota error type for shared Graph HTTP status handling.
		}
	}

	return nil
}
