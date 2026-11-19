package graphBetaIdentityAndAccessB2bManagementPolicyAssignment

import (
	"context"
	stderrors "errors"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/constants"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta"
)

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

	// Objects created moments earlier may not have reached every Entra replica yet.
	if err := crud.PollUntil(ctx, 2*time.Second, func(ctx context.Context) (bool, error) {
		_, err := r.client.
			Policies().
			B2bManagementPolicies().
			ByB2bManagementPolicyId(policyID).
			Get(ctx, nil)

		if err == nil {
			return true, nil
		}

		if errors.GraphError(ctx, err).StatusCode != 404 {
			return false, &crud.FatalPollError{Err: err}
		}

		tflog.Debug(ctx, fmt.Sprintf("B2B management policy %s not visible yet, awaiting Entra propagation", policyID))
		return false, fmt.Errorf("B2B management policy %s is not visible in the directory yet: %w", policyID, err)
	}); err != nil {
		resp.Diagnostics.AddError(
			"Error verifying B2B management policy before assignment",
			fmt.Sprintf("B2B management policy %s could not be read prior to assigning it to directory object %s: %s", policyID, directoryObjectID, err.Error()),
		)
		return
	}

	var directoryObjectType string

	err := crud.PollUntil(ctx, 2*time.Second, func(ctx context.Context) (bool, error) {
		directoryObject, err := r.client.DirectoryObjects().ByDirectoryObjectId(directoryObjectID).Get(ctx, nil)
		resolvedType, err := decodeDirectoryObjectType(directoryObject, err, directoryObjectID)
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

	// Retry only 404s from replicas the policy has not reached. An existing reference
	// must be imported rather than adopted, since destroy would otherwise remove it.
	attempt := 0
	if err := crud.PollUntil(ctx, 2*time.Second, func(ctx context.Context) (bool, error) {
		attempt++

		adapter := r.client.GetAdapter()
		requestInfo, err := constructPolicyReferenceRequest(ctx, adapter, abstractions.POST, directoryObjectType, directoryObjectID, policyID)
		if err == nil {
			err = adapter.SendNoContent(ctx, requestInfo, policyReferenceErrorMapping)
		}
		if err == nil {
			return true, nil
		}

		if isReferenceAlreadyExistsError(err) {
			return false, &crud.FatalPollError{Err: fmt.Errorf(
				"B2B management policy %s is already applied to %s %s; import it with ID %q instead of creating it: %w",
				policyID, directoryObjectType, directoryObjectID, policyID+"/"+directoryObjectID, err)}
		}

		if errors.GraphError(ctx, err).StatusCode == 404 {
			tflog.Debug(ctx, fmt.Sprintf("Assigning B2B management policy %s returned 404 (attempt %d), awaiting Entra propagation", policyID, attempt))
			return false, err
		}

		return false, &crud.FatalPollError{Err: err}
	}); err != nil {
		if isReferenceAlreadyExistsError(err) {
			resp.Diagnostics.AddError("B2B management policy assignment already exists", err.Error())
			return
		}
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationCreate, r.WritePermissions)
		return
	}

	object.ID = compositeID(policyID, directoryObjectID)
	object.DirectoryObjectType = types.StringValue(directoryObjectType)

	resp.Diagnostics.Append(resp.State.Set(ctx, &object)...)
	if resp.Diagnostics.HasError() {
		return
	}

	readReq := resource.ReadRequest{State: resp.State, ProviderMeta: req.ProviderMeta}
	stateContainer := &crud.CreateResponseContainer{CreateResponse: resp}

	opts := crud.DefaultReadWithRetryOptions()
	opts.Operation = constants.TfOperationCreate
	opts.ResourceTypeName = ResourceName
	opts.MaxRetries = 60
	opts.RetryInterval = 5 * time.Second
	opts.ConsistencyPredicate = b2bManagementPolicyAssignmentConsistencyPredicate(&object)

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

// isReferenceAlreadyExistsError matches Graph's 400 for a reference that is already present.
func isReferenceAlreadyExistsError(err error) bool {
	var odataError *odataerrors.ODataError
	if !stderrors.As(err, &odataError) || odataError.GetStatusCode() != 400 || odataError.GetErrorEscaped() == nil {
		return false
	}
	message := odataError.GetErrorEscaped().GetMessage()
	return message != nil && strings.Contains(*message, "added object references already exist")
}

// Read handles the Read operation for B2B Management Policy Assignment resources.
//
// Operation: Verifies the B2B management policy applies to the directory object
// API Calls:
//   - GET /policies/b2bManagementPolicies/{b2bManagementPolicyId}/appliesTo
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

	find := func(ctx context.Context) (graphmodels.DirectoryObjectable, error) {
		builder := r.client.Policies().B2bManagementPolicies().ByB2bManagementPolicyId(policyID).AppliesTo()
		page, err := builder.Get(ctx, nil)
		if err != nil {
			if errors.GraphError(ctx, err).StatusCode == 404 {
				return nil, nil //nolint:nilnil // No matching assignment is a valid lookup result.
			}
			return nil, wrapAppliesToError(err)
		}
		var directoryObjects []graphmodels.DirectoryObjectable
		for page != nil {
			directoryObjects = append(directoryObjects, page.GetValue()...)
			nextLink := page.GetOdataNextLink()
			if nextLink == nil || *nextLink == "" {
				break
			}
			page, err = builder.WithUrl(*nextLink).Get(ctx, nil)
			if err != nil {
				if errors.GraphError(ctx, err).StatusCode == 404 {
					return nil, nil //nolint:nilnil // No matching assignment is a valid lookup result.
				}
				return nil, wrapAppliesToError(err)
			}
		}
		for _, directoryObject := range directoryObjects {
			if directoryObject.GetId() != nil && *directoryObject.GetId() == directoryObjectID {
				return directoryObject, nil
			}
		}
		return nil, nil //nolint:nilnil // No matching assignment is a valid lookup result.
	}
	directoryObject, err := find(ctx)
	if err == nil && directoryObject == nil && operation == constants.TfOperationRead {
		tflog.Debug(ctx, fmt.Sprintf("Directory object %s not in appliesTo of B2B management policy %s, confirming it is not a stale Entra replica", directoryObjectID, policyID))
		select {
		case <-time.After(notFoundConfirmationInterval):
		case <-ctx.Done():
			err = ctx.Err()
		}
		if err == nil {
			pollCtx, cancelPoll := context.WithTimeout(ctx, notFoundConfirmationWindow)
			defer cancelPoll()
			pollErr := crud.PollUntil(pollCtx, notFoundConfirmationInterval, func(ctx context.Context) (bool, error) {
				found, findErr := find(ctx)
				if findErr != nil {
					return false, &crud.FatalPollError{Err: findErr}
				}
				if found == nil {
					return false, errAssignmentNotVisible
				}
				directoryObject = found
				return true, nil
			})
			switch {
			case pollErr == nil:
				tflog.Debug(ctx, fmt.Sprintf("Directory object %s is in appliesTo again; the miss came from a stale Entra replica", directoryObjectID))
			case pollCtx.Err() != nil && ctx.Err() == nil:
				directoryObject = nil
			default:
				err = pollErr
			}
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
		if directoryObjectType, ok := directoryObjectTypeFromODataType(directoryObject.GetOdataType()); ok {
			object.DirectoryObjectType = types.StringValue(directoryObjectType)
		}
		object.ID = compositeID(policyID, directoryObjectID)

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
		directoryObject, err := r.client.DirectoryObjects().ByDirectoryObjectId(directoryObjectID).Get(ctx, nil)
		resolvedType, err := decodeDirectoryObjectType(directoryObject, err, directoryObjectID)
		if err != nil {
			errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
			return
		}
		directoryObjectType = resolvedType
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing B2B management policy %s from %s %s", policyID, directoryObjectType, directoryObjectID))

	adapter := r.client.GetAdapter()
	requestInfo, err := constructPolicyReferenceRequest(ctx, adapter, abstractions.DELETE, directoryObjectType, directoryObjectID, policyID)
	if err == nil {
		err = adapter.SendNoContent(ctx, requestInfo, policyReferenceErrorMapping)
	}
	if err != nil {
		errors.HandleKiotaGraphError(ctx, err, resp, constants.TfOperationDelete, r.WritePermissions)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Removing %s from Terraform state", ResourceName))
	resp.State.RemoveResource(ctx)

	tflog.Debug(ctx, fmt.Sprintf("Finished Delete Method: %s", ResourceName))
}

// errAppliesToUnreadable marks an unparseable appliesTo response. Without Application.Read.All,
// Graph returns HTTP 200 with the body cut off at an application it cannot serialize. $select=id
// is not a way around it: it returns an empty list, which would drop the assignment from state.
var errAppliesToUnreadable = stderrors.New("the appliesTo response could not be parsed")

var policyReferenceErrorMapping = abstractions.ErrorMappings{
	"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue,
}

// Variables so unit tests can shorten them (export_test.go).
var (
	notFoundConfirmationWindow   = 60 * time.Second
	notFoundConfirmationInterval = 3 * time.Second
)

var errAssignmentNotVisible = stderrors.New("assignment not visible in appliesTo")
