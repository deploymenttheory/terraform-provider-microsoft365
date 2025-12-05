package graphBetaAgentIdentityBlueprintIdentifierUri

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs the Application object for PATCH request
// Note: This function may modify data.Scope.ID if a new UUID is generated
func constructResource(ctx context.Context, data *AgentIdentityBlueprintIdentifierUriResourceModel, existingUris []string) (graphmodels.Applicationable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	application := graphmodels.NewApplication()

	// Build the identifier URIs list - add the new URI to existing ones
	allUris := make([]string, 0, len(existingUris)+1)

	// Add existing URIs first (excluding the one we're adding/updating)
	targetUri := data.IdentifierUri.ValueString()
	for _, uri := range existingUris {
		if uri != targetUri {
			allUris = append(allUris, uri)
		}
	}

	// Add the new/updated URI
	allUris = append(allUris, targetUri)
	application.SetIdentifierUris(allUris)

	// Configure OAuth2 permission scope if provided
	if data.Scope != nil {
		api := graphmodels.NewApiApplication()
		scope := graphmodels.NewPermissionScope()

		// Set scope ID - generate if not provided or unknown
		if data.Scope.ID.IsNull() || data.Scope.ID.IsUnknown() {
			newID := uuid.New()
			scope.SetId(&newID)
			// Update the data model with the generated ID so it's known after apply
			data.Scope.ID = types.StringValue(newID.String())
			tflog.Debug(ctx, fmt.Sprintf("Generated new scope ID: %s", newID.String()))
		} else {
			if err := convert.FrameworkToGraphUUID(data.Scope.ID, scope.SetId); err != nil {
				return nil, fmt.Errorf("failed to parse scope id: %w", err)
			}
		}

		convert.FrameworkToGraphString(data.Scope.AdminConsentDescription, scope.SetAdminConsentDescription)
		convert.FrameworkToGraphString(data.Scope.AdminConsentDisplayName, scope.SetAdminConsentDisplayName)
		convert.FrameworkToGraphBool(data.Scope.IsEnabled, scope.SetIsEnabled)
		convert.FrameworkToGraphString(data.Scope.Type, scope.SetTypeEscaped)
		convert.FrameworkToGraphString(data.Scope.Value, scope.SetValue)

		api.SetOauth2PermissionScopes([]graphmodels.PermissionScopeable{scope})
		application.SetApi(api)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), application); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return application, nil
}

// constructDeleteResource constructs the Application object for removing a URI
func constructDeleteResource(ctx context.Context, identifierUri string, existingUris []string) (graphmodels.Applicationable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing delete request for %s", ResourceName))

	application := graphmodels.NewApplication()

	// Remove the URI from the list
	newUris := make([]string, 0, len(existingUris))
	for _, uri := range existingUris {
		if uri != identifierUri {
			newUris = append(newUris, uri)
		}
	}
	application.SetIdentifierUris(newUris)

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing delete request for %s", ResourceName))

	return application, nil
}
