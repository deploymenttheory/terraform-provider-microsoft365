package graphBetaApplicationsIpApplicationSegment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructResource(ctx context.Context, data *IpApplicationSegmentResourceModel) (graphmodels.IpApplicationSegmentable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewIpApplicationSegment()

	convert.FrameworkToGraphString(data.DestinationHost, requestBody.SetDestinationHost)

	if err := convert.FrameworkToGraphEnum(
		data.DestinationType,
		graphmodels.ParsePrivateNetworkDestinationType,
		requestBody.SetDestinationType); err != nil {
		return nil, fmt.Errorf("failed to set destination type: %v", err)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.Ports, requestBody.SetPorts); err != nil {
		return nil, fmt.Errorf("failed to set ports: %w", err)
	}

	if err := convert.FrameworkToGraphEnum(
		data.Protocol,
		graphmodels.ParsePrivateNetworkProtocol,
		requestBody.SetProtocol); err != nil {
		return nil, fmt.Errorf("failed to set protocol: %v", err)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
