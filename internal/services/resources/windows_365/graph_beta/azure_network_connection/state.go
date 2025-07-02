package graphBetaAzureNetworkConnection

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetamodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// priorPassword should be the value from the prior state or plan
func MapRemoteStateToTerraform(ctx context.Context, data *CloudPcOnPremisesConnectionResourceModel, remote msgraphbetamodels.CloudPcOnPremisesConnectionable, priorPassword types.String) {
	if remote == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}
	data.ID = convert.GraphToFrameworkString(remote.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remote.GetDisplayName())
	data.ConnectionType = convert.GraphToFrameworkEnum(remote.GetConnectionType())
	data.AdDomainName = convert.GraphToFrameworkString(remote.GetAdDomainName())
	data.AdDomainUsername = convert.GraphToFrameworkString(remote.GetAdDomainUsername())
	// Preserve sensitive password from prior state/plan if not returned by API
	if remote.GetAdDomainPassword() != nil {
		data.AdDomainPassword = convert.GraphToFrameworkString(remote.GetAdDomainPassword())
	} else {
		data.AdDomainPassword = priorPassword
	}
	data.OrganizationalUnit = convert.GraphToFrameworkString(remote.GetOrganizationalUnit())
	data.ResourceGroupId = convert.GraphToFrameworkString(remote.GetResourceGroupId())
	data.SubnetId = convert.GraphToFrameworkString(remote.GetSubnetId())
	data.SubscriptionId = convert.GraphToFrameworkString(remote.GetSubscriptionId())
	data.VirtualNetworkId = convert.GraphToFrameworkString(remote.GetVirtualNetworkId())
	data.HealthCheckStatus = convert.GraphToFrameworkEnum(remote.GetHealthCheckStatus())
	data.ManagedBy = convert.GraphToFrameworkEnum(remote.GetManagedBy())
	data.InUse = convert.GraphToFrameworkBool(remote.GetInUse())
}
