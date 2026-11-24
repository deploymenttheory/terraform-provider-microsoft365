// Package graphBetaIdentityAndAccessB2bManagementPolicyAssignment manages the assignment of a B2B management policy to an
// application or service principal.
//
// SDK deviation: missing type
//
// Evidence:
//   - SDK: github.com/microsoftgraph/msgraph-beta-sdk-go@v0.168.0 has no request builder for
//     /applications/{id}/policies/$ref, /servicePrincipals/{id}/policies/$ref or their /policies/{policyId}/$ref items.
//     The applications and serviceprincipals packages build only the claimsMappingPolicies, homeRealmDiscoveryPolicies,
//     tokenIssuancePolicies and tokenLifetimePolicies references.
//   - Learn: https://learn.microsoft.com/en-us/graph/api/b2bmanagementpolicy-post-appliesto?view=graph-rest-beta
//     documents POST /policies/b2bManagementPolicies/{id}/appliesTo/$ref instead.
//   - Observed: the documented appliesTo/$ref endpoints reject writes with 400 "Property 'appliesTo' is read-only", so
//     assignments are written through the application or service principal policies navigation, and
//     GET /policies/b2bManagementPolicies/{id}/appliesTo reflects them.
//
// Removal condition: the beta SDK ships request builders for the applications and servicePrincipals policies/$ref
// navigation, or the documented appliesTo/$ref endpoints accept writes.
//
// The appliesTo GET also goes through custom_requests, so that every page is collected by its @odata.nextLink paging
// and a body Graph cut short is reported as a JSON syntax error rather than matched on the SDK parser's message.
package graphBetaIdentityAndAccessB2bManagementPolicyAssignment

import (
	"context"
	"encoding/json"
	stderrors "errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	customrequests "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/custom_requests"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
)

// Variables so unit tests can shorten them (export_test.go).
var (
	notFoundConfirmationWindow   = 60 * time.Second
	notFoundConfirmationInterval = 3 * time.Second
)

var (
	// errAppliesToUnreadable marks an unparseable appliesTo response. Without Application.Read.All,
	// Graph returns HTTP 200 with the body cut off at an application it cannot serialize. $select=id
	// is not a way around it: it returns an empty list, which would drop the assignment from state.
	errAppliesToUnreadable = stderrors.New("the appliesTo response could not be parsed")

	errAssignmentNotVisible = stderrors.New("assignment not visible in appliesTo")

	errUnsupportedDirectoryObjectType = stderrors.New("B2B management policies can only be applied to applications and service principals")
)

// directoryObjectTypes maps the @odata.type of a supported assignment target to its directory_object_type.
var directoryObjectTypes = map[string]string{
	"#microsoft.graph.application":      "application",
	"#microsoft.graph.servicePrincipal": "servicePrincipal",
}

// directoryObjectCollections maps a directory_object_type to the Graph collection its policies navigation lives in.
var directoryObjectCollections = map[string]string{
	"application":      "applications",
	"servicePrincipal": "servicePrincipals",
}

// Create handles the Create operation for B2B Management Policy Assignment resources.
//
// Operation: Applies a B2B management policy to an application or service principal
// API Calls:
//   - GET /policies/b2bManagementPolicies/{b2bManagementPolicyId}
//   - GET /directoryObjects/{directoryObjectId}
//   - POST /{applications|servicePrincipals}/{directoryObjectId}/policies/$ref
//
// Reference: https://learn.microsoft.com/en-us/graph/api/b2bmanagementpolicy-post-appliesto?view=graph-rest-beta
func (r *B2bManagementPolicyAssignmentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var object B2bManagementPolicyAssignmentResourceModel

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

	policyID := object.B2bManagementPolicyID.ValueString()
	directoryObjectID := object.DirectoryObjectID.ValueString()

	// Objects created moments earlier may not have reached every Entra replica yet, so a 404 is polled until the
	// policy and the directory object are readable.
	err := crud.PollUntil(ctx, 2*time.Second, func(ctx context.Context) (bool, error) {
		_, err := r.client.Policies().B2bManagementPolicies().ByB2bManagementPolicyId(policyID).Get(ctx, nil)
		if err == nil {
			return true, nil
		}
		if errors.GraphError(ctx, err).StatusCode != 404 {
			return false, &crud.FatalPollError{Err: err}
		}
		tflog.Debug(ctx, fmt.Sprintf("B2B management policy %s not visible yet, awaiting Entra propagation", policyID))
		return false, fmt.Errorf("B2B management policy %s is not visible in the directory yet: %w", policyID, err)
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error verifying B2B management policy before assignment",
			fmt.Sprintf("B2B management policy %s could not be read prior to assigning it to directory object %s: %s", policyID, directoryObjectID, err.Error()),
		)
		return
	}

	var directoryObjectType string
	err = crud.PollUntil(ctx, 2*time.Second, func(ctx context.Context) (bool, error) {
		resolvedType, err := r.findDirectoryObjectTypeByID(ctx, directoryObjectID)
		if err == nil {
			directoryObjectType = resolvedType
			return true, nil
		}
		if errors.GraphError(ctx, err).StatusCode != 404 {
			return false, &crud.FatalPollError{Err: err}
		}
		tflog.Debug(ctx, fmt.Sprintf("Directory object %s not visible yet, awaiting Entra propagation", directoryObjectID))
		return false, fmt.Errorf("directory object %s is not visible in the directory yet: %w", directoryObjectID, err)
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error resolving directory object before assignment",
			fmt.Sprintf("Directory object %s could not be resolved to an application or service principal: %s", directoryObjectID, err.Error()),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Applying B2B management policy %s to %s %s", policyID, directoryObjectType, directoryObjectID))

	// The $ref POST is retried on 404: a replica the policy has not reached yet rejects it even after a GET
	// succeeded, and nothing is written. Because only 404s are retried, 400 "already exist" means the assignment
	// predates this Create; it is reported with an import hint rather than adopted, since a later destroy would
	// otherwise remove it.
	adapter := r.client.GetAdapter()
	createConfig := customrequests.PostRequestConfig{
		APIVersion: customrequests.GraphAPIBeta,
		Endpoint:   directoryObjectCollections[directoryObjectType] + "/" + directoryObjectID + "/policies/$ref",
		RawBody:    policyReference{ODataID: adapter.GetBaseUrl() + "/policies/b2bManagementPolicies/" + policyID},
	}
	// referenceAlreadyExists matches Graph's 400 for a policy reference that is already present.
	referenceAlreadyExists := func(err error) bool {
		errorInfo := errors.GraphError(ctx, err)
		return errorInfo.StatusCode == 400 && strings.Contains(errorInfo.ErrorMessage, "added object references already exist")
	}
	attempt := 0
	err = crud.PollUntil(ctx, 2*time.Second, func(ctx context.Context) (bool, error) {
		attempt++
		err := customrequests.PostRequestNoContent(ctx, adapter, createConfig)
		if err == nil {
			return true, nil
		}
		if referenceAlreadyExists(err) {
			return false, &crud.FatalPollError{Err: fmt.Errorf(
				"B2B management policy %s is already applied to %s %s; import it with ID %q instead of creating it: %w",
				policyID, directoryObjectType, directoryObjectID, policyID+"/"+directoryObjectID, err)}
		}
		if errors.GraphError(ctx, err).StatusCode == 404 {
			tflog.Debug(ctx, fmt.Sprintf("Assigning B2B management policy %s returned 404 (attempt %d), awaiting Entra propagation", policyID, attempt))
			return false, err
		}
		return false, &crud.FatalPollError{Err: err}
	})
	if err != nil && referenceAlreadyExists(err) {
		resp.Diagnostics.AddError("B2B management policy assignment already exists", err.Error())
		return
	}
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	object.ID = types.StringValue(policyID + "/" + directoryObjectID)
	object.DirectoryObjectType = types.StringValue(directoryObjectType)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Deliberate exception to "no consistency predicate on Create": stale Entra replicas can omit a new assignment
	// from appliesTo, and main already had this predicate. Reads from replicas that do not list the new assignment
	// yet are rejected; accepting them would leave Create without state.
	matchOptions := crud.AttributeMatchOptions{
		Paths:        crud.RootPaths("b2b_management_policy_id", "directory_object_id"),
		RequireKnown: crud.RootPaths("id"),
	}
	assignmentVisible := crud.ReadWithConsistencyPredicate(crud.AttributesMatchPredicate(req.Plan, matchOptions))
	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	opts := crud.NewReadWithRetryOptions(constants.TfOperationCreate, ResourceName, crud.ReadWithMaxRetries(60), crud.ReadWithRetryInterval(5*time.Second), assignmentVisible)
	if err := crud.ReadWithRetry(ctx, r.Read, readReq, &crud.CreateResponseContainer{CreateResponse: resp}, opts); err != nil {
		resp.Diagnostics.AddError("Error reading resource state after create", fmt.Sprintf("Could not read resource state: %s: %s", ResourceName, err.Error()))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished Create Method: %s", ResourceName))
}

// Read handles the Read operation for B2B Management Policy Assignment resources.
//
// Operation: Verifies the B2B management policy applies to the directory object
// API Calls:
//   - GET /policies/b2bManagementPolicies/{b2bManagementPolicyId}/appliesTo (every @odata.nextLink page)
//
// Reference: https://learn.microsoft.com/en-us/graph/api/b2bmanagementpolicy-list-appliesto?view=graph-rest-beta
func (r *B2bManagementPolicyAssignmentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var object B2bManagementPolicyAssignmentResourceModel
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

	policyID := object.B2bManagementPolicyID.ValueString()
	directoryObjectID := object.DirectoryObjectID.ValueString()

	// findAssignment returns the directory object from appliesTo, or errAssignmentNotVisible if it is absent
	// (including a 404 for the policy).
	getConfig := customrequests.GetRequestConfig{
		APIVersion:        customrequests.GraphAPIBeta,
		Endpoint:          "policies/b2bManagementPolicies",
		ResourceIDPattern: customrequests.ResourceIDPatternSegment,
		ResourceID:        policyID,
		EndpointSuffix:    "/appliesTo",
	}
	findAssignment := func(ctx context.Context) (*directoryObjectResponse, error) {
		var appliesTo directoryObjectCollectionResponse
		getConfig.RawResponse = &appliesTo
		if _, err := customrequests.GetRequestByResourceId(ctx, r.client.GetAdapter(), getConfig); err != nil {
			if syntaxError := (*json.SyntaxError)(nil); stderrors.As(err, &syntaxError) {
				return nil, fmt.Errorf("%w: %w", errAppliesToUnreadable, err)
			}
			if errors.GraphError(ctx, err).StatusCode == 404 {
				return nil, errAssignmentNotVisible
			}
			return nil, err //nolint:wrapcheck // HandleKiotaGraphError reads the ODataError from it.
		}
		for _, directoryObject := range appliesTo.Value {
			if directoryObject.ID != nil && *directoryObject.ID == directoryObjectID {
				return &directoryObject, nil
			}
		}
		return nil, errAssignmentNotVisible
	}

	directoryObject, err := findAssignment(ctx)

	// On a refresh, absence is reported only after it persists for notFoundConfirmationWindow: stale replicas
	// briefly return an empty appliesTo, and dropping the resource from state on one would make the next plan
	// re-create an existing assignment.
	assignmentMissing := stderrors.Is(err, errAssignmentNotVisible)
	if assignmentMissing {
		err = nil
	}
	if assignmentMissing && operation == constants.TfOperationRead {
		tflog.Debug(ctx, fmt.Sprintf("Directory object %s not in appliesTo of B2B management policy %s, confirming it is not a stale Entra replica", directoryObjectID, policyID))

		select {
		case <-time.After(notFoundConfirmationInterval):
		case <-ctx.Done():
			err = ctx.Err()
		}

		if err == nil {
			pollCtx, pollCancel := context.WithTimeout(ctx, notFoundConfirmationWindow)
			pollErr := crud.PollUntil(pollCtx, notFoundConfirmationInterval, func(ctx context.Context) (bool, error) {
				found, findErr := findAssignment(ctx)
				if stderrors.Is(findErr, errAssignmentNotVisible) {
					return false, findErr
				}
				if findErr != nil {
					return false, &crud.FatalPollError{Err: findErr}
				}
				directoryObject = found
				return true, nil
			})

			switch {
			case pollErr == nil:
				tflog.Debug(ctx, fmt.Sprintf("Directory object %s is in appliesTo again; the miss came from a stale replica", directoryObjectID))
			case pollCtx.Err() != nil && ctx.Err() == nil:
				directoryObject = nil
			default:
				directoryObject, err = nil, pollErr
			}
			pollCancel()
		}
	}

	if stderrors.Is(err, errAppliesToUnreadable) {
		resp.Diagnostics.AddError(
			"Error reading B2B management policy appliesTo",
			fmt.Sprintf("Microsoft Graph returned an unparseable appliesTo response for B2B management policy %s (%s). "+
				"This happens when the policy applies to an application the caller cannot read: grant Application.Read.All "+
				"(or Application.ReadWrite.All) in addition to Policy.Read.B2BManagementPolicy.", policyID, err.Error()),
		)
		return
	}
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, operation, r.ReadPermissions)
		return
	}

	if directoryObject != nil {
		if directoryObject.ODataType != nil {
			if directoryObjectType, ok := directoryObjectTypes[*directoryObject.ODataType]; ok {
				object.DirectoryObjectType = types.StringValue(directoryObjectType)
			}
		}
		object.ID = types.StringValue(policyID + "/" + directoryObjectID)

		resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
		if resp.Diagnostics.HasError() {
			return
		}

		tflog.Debug(ctx, fmt.Sprintf("Finished Read Method: %s", ResourceName))
		return
	}

	tflog.Debug(ctx, "B2B management policy assignment not found, removing from state", map[string]any{
		"b2b_management_policy_id": policyID,
		"directory_object_id":      directoryObjectID,
	})
	resp.State.RemoveResource(ctx)
}

// Update handles the Update operation for B2B Management Policy Assignment resources.
//
// Operation: Since both identifying fields have RequiresReplace, this is effectively a no-op
// (terraform will destroy and recreate). This Update implementation handles the edge case where
// only timeout changes occur.
func (r *B2bManagementPolicyAssignmentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan B2bManagementPolicyAssignmentResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete handles the Delete operation for B2B Management Policy Assignment resources.
//
// Operation: Removes a B2B management policy from an application or service principal
// API Calls:
//   - DELETE /{applications|servicePrincipals}/{directoryObjectId}/policies/{b2bManagementPolicyId}/$ref
//
// Reference: https://learn.microsoft.com/en-us/graph/api/b2bmanagementpolicy-delete-appliesto?view=graph-rest-beta
func (r *B2bManagementPolicyAssignmentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var object B2bManagementPolicyAssignmentResourceModel

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

	policyID := object.B2bManagementPolicyID.ValueString()
	directoryObjectID := object.DirectoryObjectID.ValueString()

	directoryObjectType := object.DirectoryObjectType.ValueString()
	if directoryObjectType == "" {
		resolvedType, err := r.findDirectoryObjectTypeByID(ctx, directoryObjectID)
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
			return
		}
		directoryObjectType = resolvedType
	}

	collection, ok := directoryObjectCollections[directoryObjectType]
	if !ok {
		err := fmt.Errorf("unsupported directory object type %q: %w", directoryObjectType, errUnsupportedDirectoryObjectType)
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing B2B management policy %s from %s %s", policyID, directoryObjectType, directoryObjectID))

	deleteConfig := customrequests.DeleteRequestConfig{
		APIVersion:        customrequests.GraphAPIBeta,
		Endpoint:          collection + "/" + directoryObjectID + "/policies",
		ResourceIDPattern: customrequests.ResourceIDPatternSegment,
		ResourceID:        policyID,
		EndpointSuffix:    "/$ref",
	}
	if err := customrequests.DeleteRequestByResourceId(ctx, r.client.GetAdapter(), deleteConfig); err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))
	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}
