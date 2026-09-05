package graphBetaNetworkTLSInspectionPolicyRule

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

func MapRemoteStateToTerraform(
	ctx context.Context,
	data *NetworkTLSInspectionPolicyRuleResourceModel,
	remote *tlsInspectionPolicyRuleResponse,
) error {
	if remote == nil || remote.id == nil || *remote.id == "" || remote.name == nil ||
		remote.action == nil ||
		remote.priority == nil ||
		remote.status == nil {
		return fmt.Errorf("%w: missing required fields", errInvalidResponse)
	}
	if *remote.priority < 100 || remote.conditions == nil {
		return errSystemRule
	}
	if *remote.priority > math.MaxInt32 {
		return fmt.Errorf(
			"%w: priority %d exceeds the supported API range",
			errInvalidResponse,
			*remote.priority,
		)
	}
	// Preserve the raw status, but never guess enabled for an unknown API status.
	data.Status = types.StringValue(*remote.status)
	if *remote.status != "enabled" && *remote.status != "disabled" {
		return fmt.Errorf(
			"%w: unsupported status %q; enabled cannot be determined",
			errInvalidResponse,
			*remote.status,
		)
	}
	destinations := make([]attr.Value, 0, len(remote.conditions.destinations))
	for _, destination := range remote.conditions.destinations {
		if destination.odataType == nil {
			return fmt.Errorf("%w: destination is missing @odata.type", errInvalidResponse)
		}
		var destinationType string
		switch strings.TrimPrefix(*destination.odataType, "#") {
		case "microsoft.graph.networkaccess.tlsInspectionFqdnDestination":
			destinationType = destinationTypeFQDN
		case "microsoft.graph.networkaccess.tlsInspectionWebCategoryDestination":
			destinationType = destinationTypeWebCategory
		default:
			return fmt.Errorf(
				"%w: unsupported destination type %q",
				errInvalidResponse,
				*destination.odataType,
			)
		}
		remoteValues := destination.values
		if remoteValues == nil {
			remoteValues = []string{}
		}
		values, diags := types.SetValueFrom(ctx, types.StringType, remoteValues)
		if diags.HasError() {
			return fmt.Errorf("%w: convert destination values: %v", errInvalidResponse, diags)
		}
		destinations = append(
			destinations,
			types.ObjectValueMust(
				tlsInspectionPolicyRuleDestinationObjectType().AttrTypes,
				map[string]attr.Value{"type": types.StringValue(destinationType), "values": values},
			),
		)
	}
	data.ID = convert.GraphToFrameworkString(remote.id)
	data.Name = convert.GraphToFrameworkString(remote.name)
	data.Description = convert.GraphToFrameworkString(remote.description)
	data.Action = convert.GraphToFrameworkString(remote.action)
	data.Priority = types.Int32Value(int32(*remote.priority))
	data.Enabled = types.BoolValue(*remote.status == "enabled")
	data.Destinations = types.ListValueMust(
		tlsInspectionPolicyRuleDestinationObjectType(),
		destinations,
	)
	return nil
}

func tlsInspectionPolicyRuleDestinationObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"type":   types.StringType,
			"values": types.SetType{ElemType: types.StringType},
		},
	}
}
