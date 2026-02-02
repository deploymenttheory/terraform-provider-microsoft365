package graphBetaApplicationIdentifierUri

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs the Application object for PATCH request
func constructResource(ctx context.Context, data *ApplicationIdentifierUriResourceModel, existingUris []string) (graphmodels.Applicationable, error) {
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
