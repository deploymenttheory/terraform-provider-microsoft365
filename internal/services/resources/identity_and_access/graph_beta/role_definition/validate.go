package graphBetaIdentityAndAccessRoleDefinition

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
)

// directoryResourceNamespace is the RBAC namespace whose resource actions custom directory roles can grant
const directoryResourceNamespace = "microsoft.directory"

// listResourceActionsHint tells the user where to find the valid actions
const listResourceActionsHint = "List valid actions with GET /roleManagement/directory/resourceNamespaces/microsoft.directory/resourceActions"

// errUnknownResourceActions is returned when configured actions are missing from the microsoft.directory namespace
var errUnknownResourceActions = errors.New("resource actions not found in the microsoft.directory namespace")

// validateRequest validates the role definition request before it is sent to Microsoft Graph
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data *RoleDefinitionResourceModel) error {
	tflog.Debug(ctx, "Starting role definition request validation")

	if err := validateResourceActions(ctx, client, data.RolePermissions); err != nil {
		return fmt.Errorf("validation failed for role_permissions.allowed_resource_actions: %w", err)
	}

	tflog.Debug(ctx, "Role definition request validation completed successfully")
	return nil
}

// validateResourceActions checks every configured action exists in the microsoft.directory namespace.
// Microsoft Entra accepts unknown actions on create and silently rewrites them (for example
// 'applications/basic/read' becomes 'applications/standard/read'), which would surface as
// permanent drift, so unknown actions are rejected before the write.
func validateResourceActions(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, permissions []RoleDefinitionPermission) error {
	configured := make(map[string]bool)
	for _, permission := range permissions {
		if permission.AllowedResourceActions.IsNull() || permission.AllowedResourceActions.IsUnknown() {
			continue
		}
		for _, element := range permission.AllowedResourceActions.Elements() {
			if action, ok := element.(interface{ ValueString() string }); ok {
				configured[action.ValueString()] = true
			}
		}
	}

	if len(configured) == 0 {
		return nil
	}

	known, err := listDirectoryResourceActions(ctx, client)
	if err != nil {
		return err
	}

	var unknown []string
	for action := range configured {
		if !known[strings.ToLower(action)] {
			unknown = append(unknown, action)
		}
	}

	if len(unknown) > 0 {
		sort.Strings(unknown)
		return fmt.Errorf("%w: %s. %s", errUnknownResourceActions, strings.Join(unknown, ", "), listResourceActionsHint)
	}

	return nil
}

// listDirectoryResourceActions returns the lower-cased names of all microsoft.directory resource actions
func listDirectoryResourceActions(ctx context.Context, client *msgraphbetasdk.GraphServiceClient) (map[string]bool, error) {
	actions := make(map[string]bool)

	result, err := client.
		RoleManagement().
		Directory().
		ResourceNamespaces().
		ByUnifiedRbacResourceNamespaceId(directoryResourceNamespace).
		ResourceActions().
		Get(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to list %s resource actions: %w", directoryResourceNamespace, err)
	}

	for {
		for _, action := range result.GetValue() {
			if name := action.GetName(); name != nil {
				actions[strings.ToLower(*name)] = true
			}
		}

		nextLink := result.GetOdataNextLink()
		if nextLink == nil || *nextLink == "" {
			break
		}

		result, err = client.
			RoleManagement().
			Directory().
			ResourceNamespaces().
			ByUnifiedRbacResourceNamespaceId(directoryResourceNamespace).
			ResourceActions().
			WithUrl(*nextLink).
			Get(ctx, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to list %s resource actions: %w", directoryResourceNamespace, err)
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Loaded %d %s resource actions for validation", len(actions), directoryResourceNamespace))

	return actions, nil
}
