package graphBetaRoleDefinitionAssignment

import (
	"context"
	"fmt"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// validateRoleDefinitionExists checks if the role definition exists
func validateRoleDefinitionExists(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, roleDefinitionId string) error {
	_, err := client.
		DeviceManagement().
		RoleDefinitions().
		ByRoleDefinitionId(roleDefinitionId).
		Get(ctx, nil)

	if err != nil {
		return fmt.Errorf("role definition with ID '%s' not found: %s", roleDefinitionId, err.Error())
	}

	return nil
}

// isBuiltInRole checks if the role definition ID corresponds to a built-in role
func isBuiltInRole(roleDefinitionId string) (string, bool) {
	for roleName, roleId := range BuiltInIntuneRoleDefinitions {
		if roleId == roleDefinitionId {
			return roleName, true
		}
	}
	return "", false
}

// getRoleDefinitionIdByName gets the role definition ID by built-in role name
func getRoleDefinitionIdByName(roleName string) (string, bool) {
	roleId, exists := BuiltInIntuneRoleDefinitions[roleName]
	return roleId, exists
}
