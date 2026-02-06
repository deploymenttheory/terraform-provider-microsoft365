package graphBetaApplicationsOnPremisesPublishing

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs a Graph SDK Application object from the Terraform model
// This is used for PATCH operations to update the onPremisesPublishing property
func constructResource(ctx context.Context, data *OnPremisesPublishingResourceModel) (graphmodels.Applicationable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from plan", ResourceName))

	requestBody := graphmodels.NewApplication()
	onPremisesPublishing := graphmodels.NewOnPremisesPublishing()

	convert.FrameworkToGraphString(data.ApplicationType, onPremisesPublishing.SetApplicationType)
	convert.FrameworkToGraphString(data.InternalUrl, onPremisesPublishing.SetInternalUrl)
	convert.FrameworkToGraphString(data.ExternalUrl, onPremisesPublishing.SetExternalUrl)

	err := convert.FrameworkToGraphEnum(
		data.ExternalAuthenticationType,
		graphmodels.ParseExternalAuthenticationType,
		func(val *graphmodels.ExternalAuthenticationType) {
			onPremisesPublishing.SetExternalAuthenticationType(val)
		},
	)
	if err != nil {
		return nil, fmt.Errorf("error setting external authentication type: %w", err)
	}

	convert.FrameworkToGraphBool(data.IsAccessibleViaZTNAClient, onPremisesPublishing.SetIsAccessibleViaZTNAClient)
	convert.FrameworkToGraphBool(data.IsHttpOnlyCookieEnabled, onPremisesPublishing.SetIsHttpOnlyCookieEnabled)
	convert.FrameworkToGraphBool(data.IsOnPremPublishingEnabled, onPremisesPublishing.SetIsOnPremPublishingEnabled)
	convert.FrameworkToGraphBool(data.IsPersistentCookieEnabled, onPremisesPublishing.SetIsPersistentCookieEnabled)
	convert.FrameworkToGraphBool(data.IsSecureCookieEnabled, onPremisesPublishing.SetIsSecureCookieEnabled)
	convert.FrameworkToGraphBool(data.IsStateSessionEnabled, onPremisesPublishing.SetIsStateSessionEnabled)
	convert.FrameworkToGraphBool(data.IsTranslateHostHeaderEnabled, onPremisesPublishing.SetIsTranslateHostHeaderEnabled)
	convert.FrameworkToGraphBool(data.IsTranslateLinksInBodyEnabled, onPremisesPublishing.SetIsTranslateLinksInBodyEnabled)

	requestBody.SetOnPremisesPublishing(onPremisesPublishing)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
