package graphBetaNetworkContentPolicyRule

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func MapRemoteStateToTerraform(ctx context.Context, data *NetworkContentPolicyRuleResourceModel, remote *contentPolicyRuleResponse) {
	if remote == nil {
		return
	}
	data.ID = convert.GraphToFrameworkString(remote.id)
	data.Name = convert.GraphToFrameworkString(remote.name)
	data.Description = convert.GraphToFrameworkString(remote.description)
	data.Action = convert.GraphToFrameworkString(remote.action)
	data.Priority = convert.GraphToFrameworkInt64(remote.priority)
	data.Status = convert.GraphToFrameworkString(remote.status)
	data.Activities = convert.GraphToFrameworkStringSet(ctx, remote.activities)
	data.ContentTypes = convert.GraphToFrameworkStringSet(ctx, remote.contentTypes)
	data.TextContentTypes = convert.GraphToFrameworkStringSet(ctx, remote.textContentTypes)
	data.SessionTypes = convert.GraphToFrameworkStringSet(ctx, remote.sessionTypes)

	destinations := make([]attr.Value, 0, len(remote.destinations))
	for _, destination := range remote.destinations {
		destinationType := types.StringNull()
		if destination.odataType != nil {
			destinationType = types.StringValue(terraformContentPolicyRuleDestinationType(*destination.odataType))
		}
		destinations = append(destinations, types.ObjectValueMust(contentPolicyRuleDestinationObjectType().AttrTypes, map[string]attr.Value{
			"type":   destinationType,
			"values": convert.GraphToFrameworkStringSet(ctx, destination.values),
		}))
	}
	data.Destinations = types.ListValueMust(contentPolicyRuleDestinationObjectType(), destinations)
	tflog.Debug(ctx, fmt.Sprintf("Finished mapping %s with id %s", ResourceName, data.ID.ValueString()))
}

func contentPolicyRuleDestinationObjectType() types.ObjectType {
	return types.ObjectType{AttrTypes: map[string]attr.Type{
		"type":   types.StringType,
		"values": types.SetType{ElemType: types.StringType},
	}}
}
