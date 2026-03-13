package graphBetaCrossTenantAccessPartnerSettings

import (
	"context"
	"fmt"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/sentinels"
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
// This function performs two types of validation:
//  1. Validates that the partner tenant_id exists via the Microsoft Graph API
//  2. When b2b_direct_connect_outbound.users_and_groups is configured with access_type = "blocked"
//     and one or more targets are specific user or group GUIDs (i.e. not the special value "AllUsers"),
//     this function batches the GUIDs by type and verifies each exists in the tenant via paginated
//     Graph API list calls with an OData id-in filter.
func validateRequest(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data *CrossTenantAccessPartnerSettingsResourceModel) error {
	tflog.Debug(ctx, "Starting partner settings request validation")

	// Validate tenant ID
	tenantID := data.TenantID.ValueString()
	if err := validateMicrosoftEntraOrganization(ctx, client, tenantID); err != nil {
		return fmt.Errorf("validation failed for tenant_id: %w", err)
	}

	// Validate b2b_direct_connect_outbound users and groups if configured with blocked access
	if err := validateB2bDirectConnectOutbound(ctx, client, data); err != nil {
		return err
	}

	tflog.Debug(ctx, "Partner settings request validation completed successfully")
	return nil
}

// validateB2bDirectConnectOutbound validates user and group IDs in b2b_direct_connect_outbound when access_type is "blocked"
func validateB2bDirectConnectOutbound(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data *CrossTenantAccessPartnerSettingsResourceModel) error {
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

// validateMicrosoftEntraOrganization validates that the partner tenant ID exists
func validateMicrosoftEntraOrganization(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, tenantID string) error {
	tflog.Debug(ctx, fmt.Sprintf("Validating Microsoft Entra organization with tenant ID: %s", tenantID))

	tenantInfo, err := getTenantInformationByTenantID(ctx, client, tenantID)
	if err != nil {
		return fmt.Errorf("tenant ID '%s' could not be validated: %w", tenantID, err)
	}

	if tenantInfo == nil {
		return fmt.Errorf("%w: '%s' does not exist or is not accessible", sentinels.ErrInvalidTenantID, tenantID)
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully validated tenant ID: %s", tenantID))
	return nil
}

// getTenantInformationByTenantID fetches tenant information for validation
func getTenantInformationByTenantID(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, tenantID string) (graphmodels.TenantInformationable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Fetching tenant information for tenant ID: %s", tenantID))

	tenantInfo, err := client.TenantRelationships().FindTenantInformationByTenantIdWithTenantId(&tenantID).Get(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tenant information: %w", err)
	}

	return tenantInfo, nil
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
		return fmt.Errorf("b2b_direct_connect_outbound validation failed: %w: %s", sentinels.ErrInvalidUserGUIDs, strings.Join(missing, ", "))
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
		return fmt.Errorf("b2b_direct_connect_outbound validation failed: %w: %s", sentinels.ErrInvalidGroupGUIDs, strings.Join(missing, ", "))
	}

	tflog.Debug(ctx, fmt.Sprintf("Validated %d group IDs exist in tenant", len(groupIDs)))
	return nil
}
