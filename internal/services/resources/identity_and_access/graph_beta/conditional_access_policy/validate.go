package graphBetaConditionalAccessPolicy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/client"
	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/generic_client"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// RoleDefinition represents a role definition from Microsoft Graph API
type RoleDefinition struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	Description string `json:"description"`
}

// RoleDefinitionsResponse represents the response from the role definitions API
type RoleDefinitionsResponse struct {
	Value []RoleDefinition `json:"value"`
}

// TenantInformation represents a Microsoft Entra organization information
type TenantInformation struct {
	TenantID            string  `json:"tenantId"`
	FederationBrandName *string `json:"federationBrandName"`
	DisplayName         *string `json:"displayName"`
	DefaultDomainName   *string `json:"defaultDomainName"`
}

// UserDefinition represents a user from Microsoft Graph API
type UserDefinition struct {
	ID                string  `json:"id"`
	UserPrincipalName string  `json:"userPrincipalName"`
	DisplayName       *string `json:"displayName"`
}

// UsersResponse represents the response from the users API
type UsersResponse struct {
	Value []UserDefinition `json:"value"`
}

// NamedLocation represents a named location from Microsoft Graph API
type NamedLocation struct {
	ID          string `json:"id"`
	DisplayName string `json:"displayName"`
	ODataType   string `json:"@odata.type"`
	IsTrusted   *bool  `json:"isTrusted,omitempty"` // Only present for IP locations
}

// NamedLocationsResponse represents the response from the named locations API
type NamedLocationsResponse struct {
	Value []NamedLocation `json:"value"`
}

// validateRequest validates the entire conditional access policy request
func validateRequest(ctx context.Context, httpClient *client.AuthenticatedHTTPClient, data *ConditionalAccessPolicyResourceModel) error {
	tflog.Debug(ctx, "Starting conditional access policy request validation")

	if data.Conditions != nil && data.Conditions.Users != nil && !data.Conditions.Users.ExcludeRoles.IsNull() {
		if err := validateExcludeRoles(ctx, httpClient, data.Conditions.Users.ExcludeRoles); err != nil {
			return fmt.Errorf("validation failed for exclude_roles: %w", err)
		}
	}

	if data.Conditions != nil && data.Conditions.Users != nil && !data.Conditions.Users.IncludeRoles.IsNull() {
		if err := validateIncludeRoles(ctx, httpClient, data.Conditions.Users.IncludeRoles); err != nil {
			return fmt.Errorf("validation failed for include_roles: %w", err)
		}
	}

	if data.Conditions != nil && data.Conditions.Users != nil && !data.Conditions.Users.IncludeUsers.IsNull() {
		if err := validateIncludeUsers(ctx, httpClient, data.Conditions.Users.IncludeUsers); err != nil {
			return fmt.Errorf("validation failed for include_users: %w", err)
		}
	}

	if data.Conditions != nil && data.Conditions.Users != nil && !data.Conditions.Users.ExcludeUsers.IsNull() {
		if err := validateExcludeUsers(ctx, httpClient, data.Conditions.Users.ExcludeUsers); err != nil {
			return fmt.Errorf("validation failed for exclude_users: %w", err)
		}
	}

	// Validate external tenant IDs in include_guests_or_external_users if present
	if data.Conditions != nil && data.Conditions.Users != nil && !data.Conditions.Users.IncludeGuestsOrExternalUsers.IsNull() {
		guestsAttrs := data.Conditions.Users.IncludeGuestsOrExternalUsers.Attributes()
		if externalTenantsAttr, ok := guestsAttrs["external_tenants"]; ok {
			if externalTenantsObj, ok := externalTenantsAttr.(types.Object); ok && !externalTenantsObj.IsNull() {
				tenantsAttrs := externalTenantsObj.Attributes()
				if membersAttr, ok := tenantsAttrs["members"]; ok {
					if membersSet, ok := membersAttr.(types.Set); ok && !membersSet.IsNull() {
						if err := validateMicrosoftEntraOrganization(ctx, httpClient, membersSet); err != nil {
							return fmt.Errorf("validation failed for include_guests_or_external_users.external_tenants.members: %w", err)
						}
					}
				}
			}
		}
	}

	// Validate external tenant IDs in exclude_guests_or_external_users if present
	if data.Conditions != nil && data.Conditions.Users != nil && !data.Conditions.Users.ExcludeGuestsOrExternalUsers.IsNull() {
		guestsAttrs := data.Conditions.Users.ExcludeGuestsOrExternalUsers.Attributes()
		if externalTenantsAttr, ok := guestsAttrs["external_tenants"]; ok {
			if externalTenantsObj, ok := externalTenantsAttr.(types.Object); ok && !externalTenantsObj.IsNull() {
				tenantsAttrs := externalTenantsObj.Attributes()
				if membersAttr, ok := tenantsAttrs["members"]; ok {
					if membersSet, ok := membersAttr.(types.Set); ok && !membersSet.IsNull() {
						if err := validateMicrosoftEntraOrganization(ctx, httpClient, membersSet); err != nil {
							return fmt.Errorf("validation failed for exclude_guests_or_external_users.external_tenants.members: %w", err)
						}
					}
				}
			}
		}
	}

	// Validate user inclusion assignment requirements
	if data.Conditions != nil && data.Conditions.Users != nil {
		if err := validateUserInclusionAssignments(ctx, data.Conditions.Users); err != nil {
			return fmt.Errorf("validation failed for 'include_users': %w", err)
		}
	}

	// Validate application inclusion assignment requirements
	if data.Conditions != nil && data.Conditions.Applications != nil {
		if err := validateApplicationInclusionAssignments(ctx, data.Conditions.Applications); err != nil {
			return fmt.Errorf("validation failed for 'include_applications': %w", err)
		}
	}

	// Validate trusted locations
	if data.Conditions != nil && data.Conditions.Locations != nil {
		if err := validateTrustedLocations(ctx, httpClient, data.Conditions.Locations); err != nil {
			return fmt.Errorf("validation failed for 'include_locations' or 'exclude_locations': %w", err)
		}
	}

	tflog.Debug(ctx, "Conditional access policy request validation completed successfully")
	return nil
}

// validateExcludeRoles validates that all role GUIDs in exclude_roles exist in Microsoft Graph
func validateExcludeRoles(ctx context.Context, httpClient *client.AuthenticatedHTTPClient, excludeRoles types.Set) error {
	return validateRoles(ctx, httpClient, excludeRoles, "exclude_roles")
}

// validateIncludeRoles validates that all role GUIDs in include_roles exist in Microsoft Graph
func validateIncludeRoles(ctx context.Context, httpClient *client.AuthenticatedHTTPClient, includeRoles types.Set) error {
	return validateRoles(ctx, httpClient, includeRoles, "include_roles")
}

// validateIncludeUsers validates that all user GUIDs in include_users exist in Microsoft Graph
func validateIncludeUsers(ctx context.Context, httpClient *client.AuthenticatedHTTPClient, includeUsers types.Set) error {
	return validateUsers(ctx, httpClient, includeUsers, "include_users")
}

// validateExcludeUsers validates that all user GUIDs in exclude_users exist in Microsoft Graph
func validateExcludeUsers(ctx context.Context, httpClient *client.AuthenticatedHTTPClient, excludeUsers types.Set) error {
	return validateUsers(ctx, httpClient, excludeUsers, "exclude_users")
}

// validateRoles performs the actual validation logic for role GUIDs
func validateRoles(ctx context.Context, httpClient *client.AuthenticatedHTTPClient, roleSet types.Set, fieldName string) error {
	if roleSet.IsNull() || roleSet.IsUnknown() {
		tflog.Debug(ctx, fmt.Sprintf("Skipping validation for %s: field is null or unknown", fieldName))
		return nil
	}

	var roleGUIDs []string
	elements := roleSet.Elements()

	tflog.Debug(ctx, fmt.Sprintf("Validating %d roles in %s", len(elements), fieldName))

	for _, element := range elements {
		if stringVal, ok := element.(types.String); ok && !stringVal.IsNull() {
			roleGUIDs = append(roleGUIDs, stringVal.ValueString())
		}
	}

	if len(roleGUIDs) == 0 {
		tflog.Debug(ctx, fmt.Sprintf("No role GUIDs to validate in %s", fieldName))
		return nil
	}

	// Fetch role definitions from Microsoft Graph API
	roleDefinitions, err := getRoleDefinitions(ctx, httpClient)
	if err != nil {
		return fmt.Errorf("failed to fetch role definitions from Microsoft Graph: %w", err)
	}

	// Create a map of valid role GUIDs for quick lookup
	validRoleGUIDs := make(map[string]RoleDefinition)
	for _, role := range roleDefinitions {
		validRoleGUIDs[role.ID] = role
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetched %d valid role definitions from Microsoft Graph", len(validRoleGUIDs)))

	// Validate each role GUID
	var invalidRoles []string
	for _, roleGUID := range roleGUIDs {
		if _, exists := validRoleGUIDs[roleGUID]; !exists {
			invalidRoles = append(invalidRoles, roleGUID)
			tflog.Warn(ctx, fmt.Sprintf("Invalid builtin role GUID found in %s: %s", fieldName, roleGUID))
		}
	}

	// Return error if any invalid roles found
	if len(invalidRoles) > 0 {
		return fmt.Errorf("invalid role GUIDs found in %s: %v. Please verify these built in roles exist in your Microsoft Entra ID tenant", fieldName, invalidRoles)
	}

	tflog.Debug(ctx, fmt.Sprintf("All %d role GUIDs in %s are valid", len(roleGUIDs), fieldName))
	return nil
}

// validateUsers performs the actual validation logic for user GUIDs
func validateUsers(ctx context.Context, httpClient *client.AuthenticatedHTTPClient, userSet types.Set, fieldName string) error {
	if userSet.IsNull() || userSet.IsUnknown() {
		tflog.Debug(ctx, fmt.Sprintf("Skipping validation for %s: field is null or unknown", fieldName))
		return nil
	}

	var userGUIDs []string
	var specialValues []string
	elements := userSet.Elements()

	tflog.Debug(ctx, fmt.Sprintf("Validating %d users in %s", len(elements), fieldName))

	for _, element := range elements {
		if stringVal, ok := element.(types.String); ok && !stringVal.IsNull() {
			value := stringVal.ValueString()
			// Check if it's a special value that doesn't need validation
			if value == "All" || value == "None" || value == "GuestsOrExternalUsers" {
				specialValues = append(specialValues, value)
				tflog.Debug(ctx, fmt.Sprintf("Found special value in %s: %s", fieldName, value))
			} else {
				userGUIDs = append(userGUIDs, value)
			}
		}
	}

	if len(userGUIDs) == 0 {
		tflog.Debug(ctx, fmt.Sprintf("No user GUIDs to validate in %s (found %d special values)", fieldName, len(specialValues)))
		return nil
	}

	// Check each user GUID individually since batch lookup might not be efficient
	var invalidUsers []string
	for _, userGUID := range userGUIDs {
		if err := validateUserExists(ctx, httpClient, userGUID); err != nil {
			tflog.Warn(ctx, fmt.Sprintf("Invalid user GUID found in %s: %s", fieldName, userGUID))
			invalidUsers = append(invalidUsers, userGUID)
		}
	}

	// Return error if any invalid users found
	if len(invalidUsers) > 0 {
		return fmt.Errorf("invalid user GUIDs found in %s: %v. Please verify these users exist in your Microsoft Entra ID tenant", fieldName, invalidUsers)
	}

	tflog.Debug(ctx, fmt.Sprintf("All %d user GUIDs in %s are valid", len(userGUIDs), fieldName))
	return nil
}

// getRoleDefinitions retrieves role definitions from Microsoft Graph API
func getRoleDefinitions(ctx context.Context, httpClient *client.AuthenticatedHTTPClient) ([]RoleDefinition, error) {
	tflog.Debug(ctx, "Fetching role definitions from Microsoft Graph API")

	url := httpClient.GetBaseURL() + "/roleManagement/directory/roleDefinitions"
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	tflog.Debug(ctx, fmt.Sprintf("Making GET request to: %s", url))

	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %w", err)
	}
	defer httpResp.Body.Close()

	tflog.Debug(ctx, fmt.Sprintf("GET request response status: %d %s", httpResp.StatusCode, httpResp.Status))

	if httpResp.StatusCode != http.StatusOK {
		errorInfo := errors.ExtractHTTPGraphError(ctx, httpResp)
		return nil, fmt.Errorf("unexpected response from role definitions API: %d %s (RequestID: %s)",
			httpResp.StatusCode,
			httpResp.Status,
			errorInfo.RequestID)
	}

	var response RoleDefinitionsResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error parsing role definitions response: %w", err)
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d role definitions", len(response.Value)))
	return response.Value, nil
}

// validateMicrosoftEntraOrganization validates tenant IDs by checking if they are valid Microsoft Entra organizations
func validateMicrosoftEntraOrganization(ctx context.Context, httpClient *client.AuthenticatedHTTPClient, tenantIDs types.Set) error {
	if tenantIDs.IsNull() || tenantIDs.IsUnknown() {
		tflog.Debug(ctx, "Skipping validation for tenant IDs: field is null or unknown")
		return nil
	}

	var tenantIDsList []string
	elements := tenantIDs.Elements()

	tflog.Debug(ctx, fmt.Sprintf("Validating %d tenant IDs", len(elements)))

	for _, element := range elements {
		if stringVal, ok := element.(types.String); ok && !stringVal.IsNull() {
			tenantIDsList = append(tenantIDsList, stringVal.ValueString())
		}
	}

	if len(tenantIDsList) == 0 {
		tflog.Debug(ctx, "No tenant IDs to validate")
		return nil
	}

	// Validate each tenant ID
	var invalidTenantIDs []string
	var tenantDetails []string

	for _, tenantID := range tenantIDsList {
		tenantInfo, err := getTenantInformationByTenantID(ctx, httpClient, tenantID)
		if err != nil {
			tflog.Warn(ctx, fmt.Sprintf("Error validating tenant ID %s: %v", tenantID, err))
			invalidTenantIDs = append(invalidTenantIDs, tenantID)
			continue
		}

		displayName := "Unknown"
		domainName := "Unknown"

		if tenantInfo.DisplayName != nil {
			displayName = *tenantInfo.DisplayName
		}

		if tenantInfo.DefaultDomainName != nil {
			domainName = *tenantInfo.DefaultDomainName
		}

		tenantDetails = append(tenantDetails, fmt.Sprintf("Tenant ID: %s, Name: %s, Domain: %s",
			tenantID, displayName, domainName))
	}

	// Log validated tenant details
	if len(tenantDetails) > 0 {
		tflog.Debug(ctx, "Validated tenant information:", map[string]any{
			"tenants": tenantDetails,
		})
	}

	// Return error if any invalid tenant IDs found
	if len(invalidTenantIDs) > 0 {
		return fmt.Errorf("invalid Microsoft Entra organization tenant ID found: %v. Please verify these tenant IDs are valid.", invalidTenantIDs)
	}

	tflog.Debug(ctx, fmt.Sprintf("All %d tenant IDs are valid", len(tenantIDsList)))
	return nil
}

// getTenantInformationByTenantID retrieves tenant information from Microsoft Graph API
func getTenantInformationByTenantID(ctx context.Context, httpClient *client.AuthenticatedHTTPClient, tenantID string) (*TenantInformation, error) {
	tflog.Debug(ctx, fmt.Sprintf("Validating tenant ID: %s", tenantID))

	// First attempt: Try direct lookup by tenant ID
	url := httpClient.GetBaseURL() + "/tenantRelationships/findTenantInformationByTenantId(tenantId='" + tenantID + "')"
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	tflog.Debug(ctx, fmt.Sprintf("Making GET request to: %s", url))

	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %w", err)
	}
	defer httpResp.Body.Close()

	tflog.Debug(ctx, fmt.Sprintf("GET request response status: %d %s", httpResp.StatusCode, httpResp.Status))

	if httpResp.StatusCode == http.StatusOK {
		var tenantInfo TenantInformation
		if err := json.NewDecoder(httpResp.Body).Decode(&tenantInfo); err != nil {
			return nil, fmt.Errorf("error parsing tenant information response: %w", err)
		}

		return &tenantInfo, nil
	}

	errorInfo := errors.ExtractHTTPGraphError(ctx, httpResp)
	return nil, fmt.Errorf("tenant ID validation failed: %d %s (RequestID: %s)",
		httpResp.StatusCode,
		httpResp.Status,
		errorInfo.RequestID)
}

// validateUserExists checks if a user with the given GUID exists in Microsoft Graph
func validateUserExists(ctx context.Context, httpClient *client.AuthenticatedHTTPClient, userGUID string) error {
	tflog.Debug(ctx, fmt.Sprintf("Validating user GUID: %s", userGUID))

	url := httpClient.GetBaseURL() + "/users/" + userGUID
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %w", err)
	}

	// Add select parameter to minimize data transfer - only fetch id and userPrincipalName
	q := httpReq.URL.Query()
	q.Add("$select", "id,userPrincipalName,displayName")
	httpReq.URL.RawQuery = q.Encode()

	tflog.Debug(ctx, fmt.Sprintf("Making GET request to: %s", httpReq.URL.String()))

	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("error making HTTP request: %w", err)
	}
	defer httpResp.Body.Close()

	tflog.Debug(ctx, fmt.Sprintf("GET request response status: %d %s", httpResp.StatusCode, httpResp.Status))

	if httpResp.StatusCode == http.StatusOK {
		var user UserDefinition
		if err := json.NewDecoder(httpResp.Body).Decode(&user); err != nil {
			return fmt.Errorf("error parsing user response: %w", err)
		}

		displayName := "Unknown"
		if user.DisplayName != nil {
			displayName = *user.DisplayName
		}

		tflog.Debug(ctx, fmt.Sprintf("Successfully validated user: ID=%s, UPN=%s, DisplayName=%s",
			user.ID, user.UserPrincipalName, displayName))
		return nil
	}

	if httpResp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("user not found")
	}

	errorInfo := errors.ExtractHTTPGraphError(ctx, httpResp)
	return fmt.Errorf("user validation failed: %d %s (RequestID: %s)",
		httpResp.StatusCode,
		httpResp.Status,
		errorInfo.RequestID)
}

// validateUserInclusionAssignments validates that conditional access policies have proper user inclusion assignments
// If all user inclusion assignment fields are null or empty, then include_users must contain either "All" or "None"
func validateUserInclusionAssignments(ctx context.Context, users *ConditionalAccessUsers) error {
	tflog.Debug(ctx, "Starting conditional access policy user inclusion assignment validation")

	// Helper function to check if a set has non-empty values
	hasNonEmptyValues := func(set types.Set) bool {
		if set.IsNull() || set.IsUnknown() {
			return false
		}
		for _, element := range set.Elements() {
			if stringVal, ok := element.(types.String); ok && !stringVal.IsNull() && stringVal.ValueString() != "" {
				return true
			}
		}
		return false
	}

	// Helper function to check if include_users has specific user GUIDs (not just "All"/"None")
	hasSpecificUsers := func() bool {
		if users.IncludeUsers.IsNull() || users.IncludeUsers.IsUnknown() {
			return false
		}
		for _, element := range users.IncludeUsers.Elements() {
			if stringVal, ok := element.(types.String); ok && !stringVal.IsNull() {
				value := stringVal.ValueString()
				if value != "All" && value != "None" && value != "" {
					return true
				}
			}
		}
		return false
	}

	// Helper function to check if guests/external users are configured
	hasGuestTargeting := func() bool {
		if users.IncludeGuestsOrExternalUsers.IsNull() || users.IncludeGuestsOrExternalUsers.IsUnknown() {
			return false
		}
		guestsAttrs := users.IncludeGuestsOrExternalUsers.Attributes()
		if guestTypesAttr, ok := guestsAttrs["guest_or_external_user_types"]; ok {
			if guestTypesSet, ok := guestTypesAttr.(types.Set); ok {
				return hasNonEmptyValues(guestTypesSet)
			}
		}
		return false
	}

	// Check if any user targeting is configured
	hasUserTargeting := hasSpecificUsers() ||
		hasNonEmptyValues(users.IncludeGroups) ||
		hasNonEmptyValues(users.IncludeRoles) ||
		hasGuestTargeting()

	// If no user targeting, validate include_users has "All" or "None"
	if !hasUserTargeting {
		if users.IncludeUsers.IsNull() || users.IncludeUsers.IsUnknown() || len(users.IncludeUsers.Elements()) == 0 {
			return fmt.Errorf("when conditional access policy user inclusion assignments are empty for 'include_users', 'include_groups', 'include_roles', and 'include_guests_or_external_users', then 'include_users' must be either 'All' or 'None'")
		}

		// Verify include_users contains only "All" or "None"
		for _, element := range users.IncludeUsers.Elements() {
			if stringVal, ok := element.(types.String); ok && !stringVal.IsNull() {
				value := stringVal.ValueString()
				if value != "All" && value != "None" {
					return fmt.Errorf("when conditional access policy user inclusion assignments are empty for 'include_users', 'include_groups', 'include_roles', and 'include_guests_or_external_users', then 'include_users' must contain only 'All' or 'None'")
				}
			}
		}
	} else {
		// If specific user targeting exists, validate include_users cannot be "All" or "None"
		if !users.IncludeUsers.IsNull() && !users.IncludeUsers.IsUnknown() {
			for _, element := range users.IncludeUsers.Elements() {
				if stringVal, ok := element.(types.String); ok && !stringVal.IsNull() {
					value := stringVal.ValueString()
					if value == "All" || value == "None" {
						return fmt.Errorf("when conditional access policy has specific user inclusion assignments configured for 'include_groups', 'include_roles', or 'include_guests_or_external_users', then 'include_users' cannot contain 'All' or 'None'")
					}
				}
			}
		}
	}

	tflog.Debug(ctx, "conditional access policy user inclusion assignment validation completed successfully")
	return nil
}

// validateTrustedLocations validates that location GUIDs in include_locations and exclude_locations exist in Microsoft Graph
func validateTrustedLocations(ctx context.Context, httpClient *client.AuthenticatedHTTPClient, locations *ConditionalAccessLocations) error {
	tflog.Debug(ctx, "Starting trusted locations validation")

	// Validate include_locations
	if !locations.IncludeLocations.IsNull() && !locations.IncludeLocations.IsUnknown() {
		if err := validateLocationSet(ctx, httpClient, locations.IncludeLocations, "include_locations"); err != nil {
			return fmt.Errorf("validation failed for include_locations: %w", err)
		}
	}

	// Validate exclude_locations
	if !locations.ExcludeLocations.IsNull() && !locations.ExcludeLocations.IsUnknown() {
		if err := validateLocationSet(ctx, httpClient, locations.ExcludeLocations, "exclude_locations"); err != nil {
			return fmt.Errorf("validation failed for exclude_locations: %w", err)
		}
	}

	tflog.Debug(ctx, "Trusted locations validation completed successfully")
	return nil
}

// validateLocationSet validates that all location GUIDs in a set exist in Microsoft Graph
func validateLocationSet(ctx context.Context, httpClient *client.AuthenticatedHTTPClient, locationSet types.Set, fieldName string) error {
	if locationSet.IsNull() || locationSet.IsUnknown() {
		tflog.Debug(ctx, fmt.Sprintf("Skipping validation for %s: field is null or unknown", fieldName))
		return nil
	}

	var locationGUIDs []string
	var specialValues []string
	elements := locationSet.Elements()

	tflog.Debug(ctx, fmt.Sprintf("Validating %d locations in %s", len(elements), fieldName))

	for _, element := range elements {
		if stringVal, ok := element.(types.String); ok && !stringVal.IsNull() {
			value := stringVal.ValueString()
			// Check if it's a special value that doesn't need validation
			if value == "All" || value == "AllTrusted" {
				specialValues = append(specialValues, value)
				tflog.Debug(ctx, fmt.Sprintf("Found special value in %s: %s", fieldName, value))
			} else {
				locationGUIDs = append(locationGUIDs, value)
			}
		}
	}

	if len(locationGUIDs) == 0 {
		tflog.Debug(ctx, fmt.Sprintf("No location GUIDs to validate in %s (found %d special values)", fieldName, len(specialValues)))
		return nil
	}

	// Fetch named locations from Microsoft Graph API
	namedLocations, err := fetchNamedLocations(ctx, httpClient)
	if err != nil {
		return fmt.Errorf("failed to fetch named locations from Microsoft Graph: %w", err)
	}

	// Create a map of valid location GUIDs for quick lookup
	validLocationGUIDs := make(map[string]NamedLocation)
	for _, location := range namedLocations {
		validLocationGUIDs[location.ID] = location
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetched %d valid named locations from Microsoft Graph", len(validLocationGUIDs)))

	// Validate each location GUID
	var invalidLocations []string
	for _, locationGUID := range locationGUIDs {
		if location, exists := validLocationGUIDs[locationGUID]; !exists {
			invalidLocations = append(invalidLocations, locationGUID)
			tflog.Warn(ctx, fmt.Sprintf("Invalid location GUID found in %s: %s", fieldName, locationGUID))
		} else {
			tflog.Debug(ctx, fmt.Sprintf("Valid location found in %s: %s (%s)", fieldName, locationGUID, location.DisplayName))
		}
	}

	// Return error if any invalid locations found
	if len(invalidLocations) > 0 {
		return fmt.Errorf("invalid location GUIDs found in %s: %v. Please verify these named locations exist in your Microsoft Entra ID tenant", fieldName, invalidLocations)
	}

	tflog.Debug(ctx, fmt.Sprintf("All %d location GUIDs in %s are valid", len(locationGUIDs), fieldName))
	return nil
}

// fetchNamedLocations retrieves named locations from Microsoft Graph API
func fetchNamedLocations(ctx context.Context, httpClient *client.AuthenticatedHTTPClient) ([]NamedLocation, error) {
	tflog.Debug(ctx, "Fetching named locations from Microsoft Graph API")

	url := httpClient.GetBaseURL() + "/conditionalAccess/namedLocations?$select=id,displayName,microsoft.graph.ipNamedLocation/isTrusted,microsoft.graph.compliantNetworkNamedLocation/compliantNetworkType"
	httpReq, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	tflog.Debug(ctx, fmt.Sprintf("Making GET request to: %s", url))

	httpResp, err := httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %w", err)
	}
	defer httpResp.Body.Close()

	tflog.Debug(ctx, fmt.Sprintf("GET request response status: %d %s", httpResp.StatusCode, httpResp.Status))

	if httpResp.StatusCode != http.StatusOK {
		errorInfo := errors.ExtractHTTPGraphError(ctx, httpResp)
		return nil, fmt.Errorf("unexpected response from named locations API: %d %s (RequestID: %s)",
			httpResp.StatusCode,
			httpResp.Status,
			errorInfo.RequestID)
	}

	var response NamedLocationsResponse
	if err := json.NewDecoder(httpResp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("error parsing named locations response: %w", err)
	}

	tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d named locations", len(response.Value)))
	return response.Value, nil
}

// validateApplicationInclusionAssignments validates that conditional access policies have proper application inclusion assignments
// Rules:
// 1. If all fields are empty, include_applications must be set to "None"
// 2. application_filter can only be set if include_applications has GUID-based values
// 3. include_applications and exclude_applications can be set at the same time
// 4. If include_applications is set, then include_user_actions and include_authentication_context_class_references cannot be set
func validateApplicationInclusionAssignments(ctx context.Context, applications *ConditionalAccessApplications) error {
	tflog.Debug(ctx, "Starting conditional access policy application inclusion assignment validation")

	// Helper function to check if a set has non-empty values
	hasNonEmptyValues := func(set types.Set) bool {
		if set.IsNull() || set.IsUnknown() {
			return false
		}
		for _, element := range set.Elements() {
			if stringVal, ok := element.(types.String); ok && !stringVal.IsNull() && stringVal.ValueString() != "" {
				return true
			}
		}
		return false
	}

	// Helper function to check if include_applications allows application_filter
	// application_filter is allowed with GUID values and "Office365", but not with "All" or "None"
	allowsApplicationFilter := func() bool {
		if applications.IncludeApplications.IsNull() || applications.IncludeApplications.IsUnknown() {
			return false
		}
		for _, element := range applications.IncludeApplications.Elements() {
			if stringVal, ok := element.(types.String); ok && !stringVal.IsNull() {
				value := stringVal.ValueString()
				// Allow application_filter for GUIDs and "Office365", but not for "All" or "None"
				if value == "All" || value == "None" {
					return false
				}
				// If it's "Office365" or a GUID-like value, allow application_filter
				if value == "Office365" || (len(value) == 36 && strings.Count(value, "-") == 4) {
					return true
				}
			}
		}
		return false
	}

	// Check if all application targeting fields are empty
	includeApplicationsEmpty := !hasNonEmptyValues(applications.IncludeApplications)
	excludeApplicationsEmpty := !hasNonEmptyValues(applications.ExcludeApplications)
	includeUserActionsEmpty := !hasNonEmptyValues(applications.IncludeUserActions)
	includeAuthContextEmpty := !hasNonEmptyValues(applications.IncludeAuthenticationContextClassReferences)

	// Rule 1: If all fields are empty, include_applications must be "None"
	if includeApplicationsEmpty && excludeApplicationsEmpty && includeUserActionsEmpty && includeAuthContextEmpty {
		return fmt.Errorf("when conditional access policy application fields 'include_applications', 'exclude_applications', 'include_user_actions', and 'include_authentication_context_class_references' are all empty, then 'include_applications' must be set to 'None'")
	}

	// Rule 2: application_filter can only be set if include_applications has GUID or "Office365" values (not "All" or "None")
	if applications.ApplicationFilter != nil && !applications.ApplicationFilter.Mode.IsNull() && !applications.ApplicationFilter.Rule.IsNull() {
		if !allowsApplicationFilter() {
			return fmt.Errorf("conditional access policy 'application_filter' cannot be used when 'include_applications' contains 'All' or 'None' values. It can be used with GUID values or 'Office365'")
		}
	}

	// Rule 4: If include_applications is set, then include_user_actions and include_authentication_context_class_references cannot be set
	if !includeApplicationsEmpty {
		if !includeUserActionsEmpty {
			return fmt.Errorf("conditional access policy cannot have both 'include_applications' and 'include_user_actions' configured at the same time")
		}
		if !includeAuthContextEmpty {
			return fmt.Errorf("conditional access policy cannot have both 'include_applications' and 'include_authentication_context_class_references' configured at the same time")
		}
	}

	tflog.Debug(ctx, "Conditional access policy application inclusion assignment validation completed successfully")
	return nil
}
