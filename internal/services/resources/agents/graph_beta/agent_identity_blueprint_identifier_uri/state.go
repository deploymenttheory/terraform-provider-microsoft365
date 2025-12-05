package graphBetaAgentIdentityBlueprintIdentifierUri

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps the application response to the Terraform state
func MapRemoteResourceStateToTerraform(ctx context.Context, data *AgentIdentityBlueprintIdentifierUriResourceModel, application graphmodels.Applicationable) {
	if application == nil {
		tflog.Warn(ctx, "Received nil application in MapRemoteResourceStateToTerraform")
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Starting to map remote state to Terraform state for %s", ResourceName))

	// Map scope from API response
	api := application.GetApi()
	if api != nil {
		scopes := api.GetOauth2PermissionScopes()
		tflog.Debug(ctx, fmt.Sprintf("Found %d oauth2PermissionScopes in API response", len(scopes)))

		if len(scopes) > 0 {
			// Find the scope matching our value or use the first one
			var matchedScope graphmodels.PermissionScopeable

			// If we have a scope value to match against, use it
			if data.Scope != nil && !data.Scope.Value.IsNull() && !data.Scope.Value.IsUnknown() {
				scopeValue := data.Scope.Value.ValueString()
				for _, scope := range scopes {
					if scope.GetValue() != nil && *scope.GetValue() == scopeValue {
						matchedScope = scope
						tflog.Debug(ctx, fmt.Sprintf("Found matching scope with value: %s", scopeValue))
						break
					}
				}
			}

			// If no exact match, use the first scope (for imports where we don't know the value)
			if matchedScope == nil && len(scopes) > 0 {
				matchedScope = scopes[0]
				tflog.Debug(ctx, "No exact scope match found, using first scope")
			}

			if matchedScope != nil {
				// Initialize Scope if nil (e.g., during import)
				if data.Scope == nil {
					data.Scope = &OAuth2PermissionScopeModel{}
				}

				data.Scope.ID = convert.GraphToFrameworkUUID(matchedScope.GetId())
				data.Scope.AdminConsentDescription = convert.GraphToFrameworkString(matchedScope.GetAdminConsentDescription())
				data.Scope.AdminConsentDisplayName = convert.GraphToFrameworkString(matchedScope.GetAdminConsentDisplayName())
				data.Scope.IsEnabled = convert.GraphToFrameworkBool(matchedScope.GetIsEnabled())
				data.Scope.Type = convert.GraphToFrameworkString(matchedScope.GetTypeEscaped())
				data.Scope.Value = convert.GraphToFrameworkString(matchedScope.GetValue())
				tflog.Debug(ctx, fmt.Sprintf("Mapped scope.id: %s", data.Scope.ID.ValueString()))
			}
		} else {
			tflog.Debug(ctx, "No oauth2PermissionScopes found in API response")
			data.Scope = nil
		}
	} else {
		tflog.Debug(ctx, "API object is nil in application response")
		data.Scope = nil
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping %s with identifier_uri: %s", ResourceName, data.IdentifierUri.ValueString()))
}
