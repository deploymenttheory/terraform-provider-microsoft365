package graphBetaAzureNetworkConnection

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetamodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructResource(ctx context.Context, data *CloudPcOnPremisesConnectionResourceModel) (msgraphbetamodels.CloudPcOnPremisesConnectionable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))
	requestBody := msgraphbetamodels.NewCloudPcOnPremisesConnection()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphEnum(data.ConnectionType, msgraphbetamodels.ParseCloudPcOnPremisesConnectionType, requestBody.SetConnectionType)
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
