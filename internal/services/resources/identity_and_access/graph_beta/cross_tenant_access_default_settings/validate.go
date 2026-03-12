package graphBetaCrossTenantAccessDefaultSettings

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/groups"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/users"
	graphcore "github.com/microsoftgraph/msgraph-sdk-go-core"
)

// validateRequest validates the plan data before the API request is constructed.
//
// When b2b_direct_connect_outbound.users_and_groups is configured with access_type = "blocked"
// and one or more targets are specific user or group GUIDs (i.e. not the special value "AllUsers"),
// this function batches the GUIDs by type and verifies each exists in the tenant via paginated
// Graph API list calls with an OData id-in filter. An error is returned for any GUID not found.
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data *CrossTenantAccessDefaultSettingsResourceModel) error {
	if data.B2bDirectConnectOutbound.IsNull() || data.B2bDirectConnectOutbound.IsUnknown() {
		return nil
	}

	attrs := data.B2bDirectConnectOutbound.Attributes()

	usersAndGroupsObj, ok := attrs["users_and_groups"].(types.Object)
	if !ok || usersAndGroupsObj.IsNull() || usersAndGroupsObj.IsUnknown() {
		return nil
	}

	uagAttrs := usersAndGroupsObj.Attributes()

	accessType, ok := uagAttrs["access_type"].(types.String)
	if !ok || accessType.IsNull() || accessType.ValueString() != "blocked" {
		return nil
	}

	targetsSet, ok := uagAttrs["targets"].(types.Set)
	if !ok || targetsSet.IsNull() {
		return nil
	}

	var userIDs []string
	var groupIDs []string

	for _, elem := range targetsSet.Elements() {
		targetObj, ok := elem.(types.Object)
		if !ok || targetObj.IsNull() {
			continue
		}

		targetAttrs := targetObj.Attributes()

		targetVal, ok := targetAttrs["target"].(types.String)
		if !ok || targetVal.IsNull() {
			continue
		}

		target := targetVal.ValueString()
		if target == "AllUsers" {
			continue
		}

		targetType, ok := targetAttrs["target_type"].(types.String)
		if !ok || targetType.IsNull() {
			continue
		}

		switch targetType.ValueString() {
		case "user":
			userIDs = append(userIDs, target)
		case "group":
			groupIDs = append(groupIDs, target)
		}
	}

	if len(userIDs) > 0 {
		if err := validateUsersExist(ctx, client, userIDs); err != nil {
			return err
		}
	}

	if len(groupIDs) > 0 {
		if err := validateGroupsExist(ctx, client, groupIDs); err != nil {
			return err
		}
	}

	return nil
}

// validateUsersExist verifies all supplied user GUIDs exist in the tenant using a paginated
// list call with an OData id-in filter. A single error is returned listing any GUIDs not found.
func validateUsersExist(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, userIDs []string) error {
	quotedIDs := make([]string, len(userIDs))
	for i, id := range userIDs {
		quotedIDs[i] = "'" + id + "'"
	}
	filter := "id in (" + strings.Join(quotedIDs, ",") + ")"
	selectFields := []string{"id"}

	tflog.Debug(ctx, fmt.Sprintf("Validating %d user IDs against tenant via paginated users list", len(userIDs)))

	requestParams := &users.UsersRequestBuilderGetRequestConfiguration{
		QueryParameters: &users.UsersRequestBuilderGetQueryParameters{
			Filter: &filter,
			Select: selectFields,
		},
	}

	usersResponse, err := client.Users().Get(ctx, requestParams)
	if err != nil {
		return fmt.Errorf("b2b_direct_connect_outbound validation failed: could not query users: %w", err)
	}

	foundIDs := make(map[string]bool)

	pageIterator, err := graphcore.NewPageIterator[graphmodels.Userable](
		usersResponse,
		client.GetAdapter(),
		graphmodels.CreateUserCollectionResponseFromDiscriminatorValue,
	)
	if err != nil {
		return fmt.Errorf("b2b_direct_connect_outbound validation failed: could not create user page iterator: %w", err)
	}

	err = pageIterator.Iterate(ctx, func(item graphmodels.Userable) bool {
		if item != nil && item.GetId() != nil {
			foundIDs[*item.GetId()] = true
		}
		return true
	})
	if err != nil {
		return fmt.Errorf("b2b_direct_connect_outbound validation failed: error iterating user pages: %w", err)
	}

	var missing []string
	for _, id := range userIDs {
		if !foundIDs[id] {
			missing = append(missing, id)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("b2b_direct_connect_outbound validation failed: user IDs not found in tenant: %s", strings.Join(missing, ", "))
	}

	tflog.Debug(ctx, fmt.Sprintf("Validated %d user IDs exist in tenant", len(userIDs)))
	return nil
}

// validateGroupsExist verifies all supplied group GUIDs exist in the tenant using a paginated
// list call with an OData id-in filter. A single error is returned listing any GUIDs not found.
func validateGroupsExist(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, groupIDs []string) error {
	quotedIDs := make([]string, len(groupIDs))
	for i, id := range groupIDs {
		quotedIDs[i] = "'" + id + "'"
	}
	filter := "id in (" + strings.Join(quotedIDs, ",") + ")"
	selectFields := []string{"id"}

	tflog.Debug(ctx, fmt.Sprintf("Validating %d group IDs against tenant via paginated groups list", len(groupIDs)))

	requestParams := &groups.GroupsRequestBuilderGetRequestConfiguration{
		QueryParameters: &groups.GroupsRequestBuilderGetQueryParameters{
			Filter: &filter,
			Select: selectFields,
		},
	}

	groupsResponse, err := client.Groups().Get(ctx, requestParams)
	if err != nil {
		return fmt.Errorf("b2b_direct_connect_outbound validation failed: could not query groups: %w", err)
	}

	foundIDs := make(map[string]bool)

	pageIterator, err := graphcore.NewPageIterator[graphmodels.Groupable](
		groupsResponse,
		client.GetAdapter(),
		graphmodels.CreateGroupCollectionResponseFromDiscriminatorValue,
	)
	if err != nil {
		return fmt.Errorf("b2b_direct_connect_outbound validation failed: could not create group page iterator: %w", err)
	}

	err = pageIterator.Iterate(ctx, func(item graphmodels.Groupable) bool {
		if item != nil && item.GetId() != nil {
			foundIDs[*item.GetId()] = true
		}
		return true
	})
	if err != nil {
		return fmt.Errorf("b2b_direct_connect_outbound validation failed: error iterating group pages: %w", err)
	}

	var missing []string
	for _, id := range groupIDs {
		if !foundIDs[id] {
			missing = append(missing, id)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("b2b_direct_connect_outbound validation failed: group IDs not found in tenant: %s", strings.Join(missing, ", "))
	}

	tflog.Debug(ctx, fmt.Sprintf("Validated %d group IDs exist in tenant", len(groupIDs)))
	return nil
}
