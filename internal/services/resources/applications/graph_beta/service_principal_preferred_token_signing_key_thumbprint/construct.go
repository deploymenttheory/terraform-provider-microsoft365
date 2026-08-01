package graphBetaServicePrincipalPreferredTokenSigningKeyThumbprint

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs the ServicePrincipal object for the PATCH request
func constructResource(ctx context.Context, data *ServicePrincipalPreferredTokenSigningKeyThumbprintResourceModel) (graphmodels.ServicePrincipalable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	servicePrincipal := graphmodels.NewServicePrincipal()
	thumbprint := data.Thumbprint.ValueString()
	servicePrincipal.SetPreferredTokenSigningKeyThumbprint(&thumbprint)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), servicePrincipal); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return servicePrincipal, nil
}

// constructDeleteResource constructs the ServicePrincipal object that clears
// preferredTokenSigningKeyThumbprint. The generated setter omits nil fields during
// serialization, so the explicit JSON null must be sent via additional data.
func constructDeleteResource(ctx context.Context) (graphmodels.ServicePrincipalable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing delete request for %s", ResourceName))

	servicePrincipal := graphmodels.NewServicePrincipal()
	servicePrincipal.SetAdditionalData(map[string]any{
		"preferredTokenSigningKeyThumbprint": nil,
	})

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing delete request for %s", ResourceName))

	return servicePrincipal, nil
}
