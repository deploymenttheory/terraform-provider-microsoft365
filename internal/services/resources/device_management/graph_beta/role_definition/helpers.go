package graphBetaRoleDefinition

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// checkRoleNameUniqueness verifies that the role name is unique among existing roles
func checkRoleNameUniqueness(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, displayName string) error {
	tflog.Debug(ctx, fmt.Sprintf("Checking if role definition with display name '%s' already exists", displayName))

	existingRoles, err := client.
		DeviceManagement().
		RoleDefinitions().
		Get(ctx, nil)

	if err != nil {
		return fmt.Errorf("failed to retrieve existing role definitions: %v", err)
	}

	roles := existingRoles.GetValue()
	for _, role := range roles {
		if role.GetDisplayName() != nil && *role.GetDisplayName() == displayName {
			return fmt.Errorf("a role definition with the display name '%s' already exists - role names must be unique", displayName)
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Role definition name '%s' is unique", displayName))
	return nil
}
