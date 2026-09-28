package graphBetaIdentityAndAccessB2bManagementPolicyAssignment

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/hashicorp/terraform-plugin-framework/types"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models/odataerrors"
)

// The documented POST/DELETE /policies/b2bManagementPolicies/{id}/appliesTo/$ref endpoints reject
// every request with 400 "Property 'appliesTo' is read-only and cannot be set.". The assignment is
// instead written from the application or service principal side through the policies navigation
// property, which the SDK does not generate. GET .../appliesTo reflects assignments made this way.
const (
	policyReferenceCollectionURLTemplate = "{+baseurl}/{directoryObjectCollection}/{directoryObjectId}/policies/$ref"
	policyReferenceItemURLTemplate       = "{+baseurl}/{directoryObjectCollection}/{directoryObjectId}/policies/{policyId}/$ref"

	directoryObjectTypeApplication      = "application"
	directoryObjectTypeServicePrincipal = "servicePrincipal"
)

// errAppliesToUnreadable marks an appliesTo response that could not be parsed. When the caller
// lacks Application.Read.All and the policy applies to an application, Microsoft Graph answers
// GET .../appliesTo with HTTP 200 and a body that is cut off mid-object and followed by an
// InternalServerError ("The property 'signInAudienceRestrictions[Nullable=False]' ... has a null
// value"). With $select=id the same caller instead gets an empty collection, which would silently
// drop the assignment from state, so the full representation is requested and the failure is
// surfaced.
var errAppliesToUnreadable = errors.New("the appliesTo response could not be parsed")

var policyReferenceErrorMapping = abstractions.ErrorMappings{
	"XXX": odataerrors.CreateODataErrorFromDiscriminatorValue,
}

// compositeID builds the synthetic assignment ID.
func compositeID(policyID, directoryObjectID string) types.String {
	return types.StringValue(policyID + "/" + directoryObjectID)
}

// directoryObjectTypeFromODataType maps an @odata.type to the supported assignment target type.
func directoryObjectTypeFromODataType(odataType *string) (string, bool) {
	if odataType == nil {
		return "", false
	}
	switch *odataType {
	case "#microsoft.graph.application":
		return directoryObjectTypeApplication, true
	case "#microsoft.graph.servicePrincipal":
		return directoryObjectTypeServicePrincipal, true
	default:
		return "", false
	}
}

func directoryObjectCollection(directoryObjectType string) (string, error) {
	switch directoryObjectType {
	case directoryObjectTypeApplication:
		return "applications", nil
	case directoryObjectTypeServicePrincipal:
		return "servicePrincipals", nil
	default:
		return "", fmt.Errorf("unsupported directory object type %q: B2B management policies can only be applied to applications and service principals", directoryObjectType)
	}
}

// resolveDirectoryObjectType reads the directory object to determine whether it is an application
// or a service principal.
func (r *B2bManagementPolicyAssignmentResource) resolveDirectoryObjectType(ctx context.Context, directoryObjectID string) (string, error) {
	directoryObject, err := r.client.DirectoryObjects().ByDirectoryObjectId(directoryObjectID).Get(ctx, nil)
	if err != nil {
		return "", err
	}

	directoryObjectType, ok := directoryObjectTypeFromODataType(directoryObject.GetOdataType())
	if !ok {
		odataType := "<nil>"
		if directoryObject.GetOdataType() != nil {
			odataType = *directoryObject.GetOdataType()
		}
		return "", fmt.Errorf("directory object %s is of type %s: B2B management policies can only be applied to applications and service principals", directoryObjectID, odataType)
	}

	return directoryObjectType, nil
}

// addPolicyReference sends POST /{applications|servicePrincipals}/{id}/policies/$ref.
func (r *B2bManagementPolicyAssignmentResource) addPolicyReference(ctx context.Context, directoryObjectType, directoryObjectID, policyID string) error {
	collection, err := directoryObjectCollection(directoryObjectType)
	if err != nil {
		return err
	}

	adapter := r.client.GetAdapter()
	requestInfo := abstractions.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(
		abstractions.POST,
		policyReferenceCollectionURLTemplate,
		map[string]string{
			"baseurl":                   adapter.GetBaseUrl(),
			"directoryObjectCollection": collection,
			"directoryObjectId":         directoryObjectID,
		},
	)
	requestInfo.Headers.TryAdd("Accept", "application/json")

	requestBody := graphmodels.NewReferenceCreate()
	odataID := fmt.Sprintf("%s/policies/b2bManagementPolicies/%s", adapter.GetBaseUrl(), policyID)
	requestBody.SetOdataId(&odataID)

	if err := requestInfo.SetContentFromParsable(ctx, adapter, "application/json", requestBody); err != nil {
		return fmt.Errorf("set B2B management policy reference request content: %w", err)
	}

	return adapter.SendNoContent(ctx, requestInfo, policyReferenceErrorMapping)
}

// removePolicyReference sends DELETE /{applications|servicePrincipals}/{id}/policies/{policyId}/$ref.
func (r *B2bManagementPolicyAssignmentResource) removePolicyReference(ctx context.Context, directoryObjectType, directoryObjectID, policyID string) error {
	collection, err := directoryObjectCollection(directoryObjectType)
	if err != nil {
		return err
	}

	adapter := r.client.GetAdapter()
	requestInfo := abstractions.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(
		abstractions.DELETE,
		policyReferenceItemURLTemplate,
		map[string]string{
			"baseurl":                   adapter.GetBaseUrl(),
			"directoryObjectCollection": collection,
			"directoryObjectId":         directoryObjectID,
			"policyId":                  policyID,
		},
	)
	requestInfo.Headers.TryAdd("Accept", "application/json")

	return adapter.SendNoContent(ctx, requestInfo, policyReferenceErrorMapping)
}

// listAppliesTo returns every directory object the policy applies to, following @odata.nextLink.
func (r *B2bManagementPolicyAssignmentResource) listAppliesTo(ctx context.Context, policyID string) ([]graphmodels.DirectoryObjectable, error) {
	builder := r.client.Policies().B2bManagementPolicies().ByB2bManagementPolicyId(policyID).AppliesTo()

	page, err := builder.Get(ctx, nil)
	if err != nil {
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
			return nil, wrapAppliesToError(err)
		}
	}

	return directoryObjects, nil
}

// wrapAppliesToError tags errors that are not Graph API, transport or context errors, i.e. a
// successful response whose body failed to parse, with errAppliesToUnreadable.
func wrapAppliesToError(err error) error {
	var apiError abstractions.ApiErrorable
	var urlError *url.Error
	if errors.As(err, &apiError) || errors.As(err, &urlError) ||
		errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return err
	}
	return fmt.Errorf("%w: %w", errAppliesToUnreadable, err)
}
