package graphBetaNetworkCrossTenantAccessSettings

import (
	"context"
	"fmt"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
)

func constructResource(ctx context.Context, data *NetworkCrossTenantAccessSettingsResourceModel) (graphmodels.CrossTenantAccessSettingsable, error) {
	body := graphmodels.NewCrossTenantAccessSettings()
	if err := convert.FrameworkToGraphEnum(data.NetworkPacketTaggingStatus, graphmodels.ParseStatus, body.SetNetworkPacketTaggingStatus); err != nil {
		return nil, fmt.Errorf("construct network packet tagging status: %w", err)
	}
	return body, nil
}
