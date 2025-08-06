package graphBetaRoleDefinition

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// validateRequest validates that all provided role permissions exist in the list of available operations
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, permissions []string, resp any, readPermissions []string) ([]string, error) {
	tflog.Debug(ctx, "Validating Intune role permissions against available Intune resource operations")

	operations, err := client.
		DeviceManagement().
		ResourceOperations().
		Get(ctx, nil)

	if err != nil {
		errors.HandleGraphError(ctx, err, resp, "Read", readPermissions)
		return nil, err
	}

	validOperations := make(map[string]bool)
	var validOperationsList []string
	for _, op := range operations.GetValue() {
		if op.GetId() != nil {
			validOperations[*op.GetId()] = true
			validOperationsList = append(validOperationsList, *op.GetId())
		}
	}

	var invalidPermissions []string
	for _, permission := range permissions {
		if !validOperations[permission] {
			invalidPermissions = append(invalidPermissions, permission)
		}
	}

	if len(invalidPermissions) > 0 {
		return nil, fmt.Errorf("invalid resource operation ID(s) %v. Valid operations are: %v", invalidPermissions, validOperationsList)
	}

	tflog.Debug(ctx, fmt.Sprintf("Validated %d permissions successfully", len(permissions)))
	return permissions, nil
}
