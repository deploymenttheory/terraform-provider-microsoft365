package graphBetaAzureNetworkConnection

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructResource(ctx context.Context, data *CloudPcOnPremisesConnectionResourceModel) (models.CloudPcOnPremisesConnectionable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))
	requestBody := models.NewCloudPcOnPremisesConnection()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	if err := convert.FrameworkToGraphEnum(data.ConnectionType, models.ParseCloudPcOnPremisesConnectionType, requestBody.SetConnectionType); err != nil {
		return nil, fmt.Errorf("failed to set connectionType: %v", err)
	}
	convert.FrameworkToGraphString(data.AdDomainName, requestBody.SetAdDomainName)
	convert.FrameworkToGraphString(data.AdDomainUsername, requestBody.SetAdDomainUsername)
	convert.FrameworkToGraphString(data.AdDomainPassword, requestBody.SetAdDomainPassword)
	convert.FrameworkToGraphString(data.OrganizationalUnit, requestBody.SetOrganizationalUnit)
	convert.FrameworkToGraphString(data.ResourceGroupId, requestBody.SetResourceGroupId)
	convert.FrameworkToGraphString(data.SubnetId, requestBody.SetSubnetId)
	convert.FrameworkToGraphString(data.SubscriptionId, requestBody.SetSubscriptionId)
	convert.FrameworkToGraphString(data.VirtualNetworkId, requestBody.SetVirtualNetworkId)

	return requestBody, nil
}
