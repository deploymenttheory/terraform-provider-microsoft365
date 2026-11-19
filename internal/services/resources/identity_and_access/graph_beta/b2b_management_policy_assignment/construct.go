package graphBetaIdentityAndAccessB2bManagementPolicyAssignment

import (
	"context"
	stderrors "errors"
	"fmt"

	abstractions "github.com/microsoft/kiota-abstractions-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// The documented appliesTo/$ref endpoints reject writes with 400 "Property 'appliesTo' is
// read-only", so assignments are written through the application or service principal policies
// navigation, which the SDK does not generate. GET .../appliesTo reflects them.
const (
	policyReferenceCollectionURLTemplate = "{+baseurl}/{directoryObjectCollection}/{directoryObjectId}/policies/$ref"
	policyReferenceItemURLTemplate       = "{+baseurl}/{directoryObjectCollection}/{directoryObjectId}/policies/{policyId}/$ref"

	directoryObjectTypeApplication      = "application"
	directoryObjectTypeServicePrincipal = "servicePrincipal"
)

// kiotaInvalidJSONMessage is kiota-serialization-json's (unexported) error for a body failing json.Valid.
const kiotaInvalidJSONMessage = "invalid json type"

func directoryObjectCollection(directoryObjectType string) (string, error) {
	switch directoryObjectType {
	case directoryObjectTypeApplication:
		return "applications", nil
	case directoryObjectTypeServicePrincipal:
		return "servicePrincipals", nil
	default:
		return "", fmt.Errorf("unsupported directory object type %q: %w", directoryObjectType, errUnsupportedAssignmentTarget)
	}
}

// constructPolicyReferenceRequest builds the application or service-principal policy reference request.
func constructPolicyReferenceRequest(
	ctx context.Context,
	adapter abstractions.RequestAdapter,
	method abstractions.HttpMethod,
	directoryObjectType, directoryObjectID, policyID string,
) (*abstractions.RequestInformation, error) {
	collection, err := directoryObjectCollection(directoryObjectType)
	if err != nil {
		return nil, err
	}
	template := policyReferenceCollectionURLTemplate
	if method == abstractions.DELETE {
		template = policyReferenceItemURLTemplate
	}
	requestInfo := abstractions.NewRequestInformationWithMethodAndUrlTemplateAndPathParameters(method, template, map[string]string{
		"baseurl": adapter.GetBaseUrl(), "directoryObjectCollection": collection, "directoryObjectId": directoryObjectID, "policyId": policyID,
	})
	requestInfo.Headers.TryAdd("Accept", "application/json")
	if method == abstractions.POST {
		body := graphmodels.NewReferenceCreate()
		odataID := fmt.Sprintf("%s/policies/b2bManagementPolicies/%s", adapter.GetBaseUrl(), policyID)
		body.SetOdataId(&odataID)
		if err := requestInfo.SetContentFromParsable(ctx, adapter, "application/json", body); err != nil {
			return nil, fmt.Errorf("set B2B management policy reference request content: %w", err)
		}
	}
	return requestInfo, nil
}

var errUnsupportedAssignmentTarget = stderrors.New("B2B management policies can only be applied to applications and service principals")
