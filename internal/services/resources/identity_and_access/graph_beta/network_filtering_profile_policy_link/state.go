package graphBetaNetworkFilteringProfilePolicyLink

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"
)

func MapRemoteStateToTerraform(ctx context.Context, data *NetworkFilteringProfilePolicyLinkResourceModel, filteringProfileID string, link models.PolicyLinkable) {
	if link == nil {
		tflog.Debug(ctx, "Remote policy link is nil")
		return
	}

	policyLinkID := convert.GraphToFrameworkString(link.GetId())
	data.FilteringProfileID = types.StringValue(filteringProfileID)
	data.PolicyLinkID = policyLinkID
	data.ID = types.StringNull()
	if filteringProfileID != "" && policyLinkID.ValueString() != "" {
		data.ID = types.StringValue(filteringProfileID + "/" + policyLinkID.ValueString())
	}
	data.State = convert.GraphToFrameworkEnum(link.GetState())
	data.Version = convert.GraphToFrameworkString(link.GetVersion())

	var policyODataType string
	if policy := link.GetPolicy(); policy != nil {
		data.PolicyID = convert.GraphToFrameworkString(policy.GetId())
		if policyType := policy.GetOdataType(); policyType != nil {
			policyODataType = *policyType
		}
	}

	var linkODataType string
	if odataType := link.GetOdataType(); odataType != nil {
		linkODataType = *odataType
	}
	matched := false
	for policyType, odataTypes := range policyTypeToODataTypes {
		if (linkODataType != "" && odataTypes.link == linkODataType) || (policyODataType != "" && odataTypes.policy == policyODataType) {
			data.PolicyType = types.StringValue(policyType)
			matched = true
			break
		}
	}
	if !matched {
		tflog.Warn(ctx, "Remote policy link type is not supported by this Terraform resource version", map[string]any{
			"link_odata_type":   linkODataType,
			"policy_odata_type": policyODataType,
		})
	}

	if filteringLink, ok := link.(models.FilteringPolicyLinkable); ok {
		data.Priority = convert.GraphToFrameworkInt64(filteringLink.GetPriority())
		data.LoggingState = convert.GraphToFrameworkEnum(filteringLink.GetLoggingState())
	} else {
		data.Priority = types.Int64Null()
		data.LoggingState = types.StringNull()
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished stating %s with id %s", ResourceName, data.ID.ValueString()))
}
