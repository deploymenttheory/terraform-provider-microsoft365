package graphBetaNetworkForwardingProfilePolicyLink

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
)

func MapRemoteStateToTerraform(ctx context.Context, data *NetworkForwardingProfilePolicyLinkResourceModel, forwardingProfileID string, link models.PolicyLinkable) {
	if link == nil {
		tflog.Debug(ctx, "Remote forwarding policy link is nil")
		return
	}

	policyLinkID := convert.GraphToFrameworkString(link.GetId())
	data.ForwardingProfileID = types.StringValue(forwardingProfileID)
	data.PolicyLinkID = policyLinkID
	data.ID = compositeID(forwardingProfileID, policyLinkID.ValueString())
	data.State = convert.GraphToFrameworkEnum(link.GetState())
	data.Version = convert.GraphToFrameworkString(link.GetVersion())

	if forwardingLink, ok := link.(models.ForwardingPolicyLinkable); ok {
		data.Priority = convert.GraphToFrameworkInt64(forwardingLink.GetPriority())
	} else {
		data.Priority = types.Int64Null()
	}

	if policy := link.GetPolicy(); policy != nil {
		data.PolicyID = convert.GraphToFrameworkString(policy.GetId())
		data.PolicyName = convert.GraphToFrameworkString(policy.GetName())
		data.PolicyDescription = convert.GraphToFrameworkString(policy.GetDescription())
		if forwardingPolicy, ok := policy.(models.ForwardingPolicyable); ok {
			data.TrafficForwardingType = convert.GraphToFrameworkEnum(forwardingPolicy.GetTrafficForwardingType())
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping %s with id %s", ResourceName, data.ID.ValueString()))
}
