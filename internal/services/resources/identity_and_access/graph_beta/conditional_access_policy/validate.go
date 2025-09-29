package graphBetaConditionalAccessPolicy

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

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
	roleDefinitions, err := fetchRoleDefinitions(ctx, httpClient)
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

// fetchRoleDefinitions retrieves role definitions from Microsoft Graph API
func fetchRoleDefinitions(ctx context.Context, httpClient *client.AuthenticatedHTTPClient) ([]RoleDefinition, error) {
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
