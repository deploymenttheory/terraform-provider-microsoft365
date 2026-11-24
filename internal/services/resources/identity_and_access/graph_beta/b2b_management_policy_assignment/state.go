package graphBetaIdentityAndAccessB2bManagementPolicyAssignment

import (
	"context"
	"fmt"
)

// findDirectoryObjectTypeByID reports whether the directory object is an application or a service principal.
func (r *B2bManagementPolicyAssignmentResource) findDirectoryObjectTypeByID(ctx context.Context, directoryObjectID string) (string, error) {
	directoryObject, err := r.client.DirectoryObjects().ByDirectoryObjectId(directoryObjectID).Get(ctx, nil)
	if err != nil {
		return "", err //nolint:wrapcheck // Callers read the ODataError status code from it.
	}

	odataType := "<nil>"
	if directoryObject.GetOdataType() != nil {
		odataType = *directoryObject.GetOdataType()
	}
	directoryObjectType, ok := directoryObjectTypes[odataType]
	if !ok {
		return "", fmt.Errorf("directory object %s is of type %s: %w", directoryObjectID, odataType, errUnsupportedDirectoryObjectType)
	}

	return directoryObjectType, nil
}
