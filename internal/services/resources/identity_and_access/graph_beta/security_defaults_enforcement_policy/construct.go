package graphBetaIdentityAndAccessSecurityDefaultsEnforcementPolicy

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

var errInvalidIsEnabled = errors.New("is_enabled must be a known Boolean")

// constructResource converts the Terraform resource model to a Kiota SDK model for the PATCH request.
// The policy needs no lookups to build the request, so the client is unused; it is kept so the
// signature matches the other resources in this domain.
//
// Only isEnabled is writable. displayName and description are set by Microsoft Entra ID and are never sent.
func constructResource(
	ctx context.Context,
	_ *msgraphbetasdk.GraphServiceClient,
	data *SecurityDefaultsEnforcementPolicyResourceModel,
) (graphmodels.IdentitySecurityDefaultsEnforcementPolicyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	if data.IsEnabled.IsNull() || data.IsEnabled.IsUnknown() {
		return nil, errInvalidIsEnabled
	}

	requestBody := graphmodels.NewIdentitySecurityDefaultsEnforcementPolicy()
	convert.FrameworkToGraphBool(data.IsEnabled, requestBody.SetIsEnabled)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
