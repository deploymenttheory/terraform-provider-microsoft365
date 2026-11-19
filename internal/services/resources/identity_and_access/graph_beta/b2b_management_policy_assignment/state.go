package graphBetaIdentityAndAccessB2bManagementPolicyAssignment

import (
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

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

// wrapAppliesToError tags only JSON parse failures, so auth or transport errors are not blamed on
// a missing Application.Read.All.
func wrapAppliesToError(err error) error {
	var apiError abstractions.ApiErrorable
	if errors.As(err, &apiError) || !strings.Contains(err.Error(), kiotaInvalidJSONMessage) {
		return err
	}
	return fmt.Errorf("%w: %w", errAppliesToUnreadable, err)
}

// decodeDirectoryObjectType validates the supported B2B assignment target type.
func decodeDirectoryObjectType(directoryObject graphmodels.DirectoryObjectable, err error, directoryObjectID string) (string, error) {
	if err != nil {
		return "", err
	}
	directoryObjectType, ok := directoryObjectTypeFromODataType(directoryObject.GetOdataType())
	if !ok {
		odataType := "<nil>"
		if directoryObject.GetOdataType() != nil {
			odataType = *directoryObject.GetOdataType()
		}
		return "", fmt.Errorf("directory object %s is of type %s: %w", directoryObjectID, odataType, errUnsupportedAssignmentTarget)
	}
	return directoryObjectType, nil
}
