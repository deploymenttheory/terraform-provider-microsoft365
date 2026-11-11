package graphBetaAzureNetworkConnection

import (
	"context"
	"strings"

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
	data.ResourceGroupId = azureIDFromRemote(data.ResourceGroupId, remote.GetResourceGroupId())
	data.SubnetId = azureIDFromRemote(data.SubnetId, remote.GetSubnetId())
	data.SubscriptionId = azureIDFromRemote(data.SubscriptionId, remote.GetSubscriptionId())
	data.VirtualNetworkId = azureIDFromRemote(data.VirtualNetworkId, remote.GetVirtualNetworkId())
	data.HealthCheckStatus = convert.GraphToFrameworkEnum(remote.GetHealthCheckStatus())
	data.ManagedBy = convert.GraphToFrameworkEnum(remote.GetManagedBy())
	data.InUse = convert.GraphToFrameworkBool(remote.GetInUse())
}

func azureIDFromRemote(prior types.String, remote *string) types.String {
	// Azure IDs are case-insensitive. Keep the plan/prior state spelling when
	// Graph changes only casing, without hiding a different or missing ID.
	if !prior.IsNull() && !prior.IsUnknown() && remote != nil && strings.EqualFold(prior.ValueString(), *remote) {
		return prior
	}
	return convert.GraphToFrameworkString(remote)
}
