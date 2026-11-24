package graphBetaNetworkTLSInspectionPolicyRule

import (
	"context"
	"fmt"
	"math"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

func MapRemoteStateToTerraform(
	ctx context.Context,
	data *NetworkTLSInspectionPolicyRuleResourceModel,
	remote models.TlsInspectionRuleable,
) error {
	if remote == nil || remote.GetId() == nil || *remote.GetId() == "" || remote.GetName() == nil ||
		remote.GetPriority() == nil || remote.GetSettings() == nil ||
		remote.GetSettings().GetStatus() == nil {
		return fmt.Errorf("%w: missing required fields", errInvalidResponse)
	}
	// action is read from additional data because the SDK model lacks it (see the SDK deviation in crud.go).
	var action *string
	switch value := remote.GetAdditionalData()["action"].(type) {
	case string:
		action = &value
	case *string:
		action = value
	}
	if action == nil {
		return fmt.Errorf("%w: action is missing or invalid", errInvalidResponse)
	}
	priority := *remote.GetPriority()
	if priority < 100 || remote.GetMatchingConditions() == nil {
		return errSystemRule
	}
	if priority > math.MaxInt32 {
		return fmt.Errorf(
			"%w: priority %d exceeds the supported API range",
			errInvalidResponse,
			priority,
		)
	}
	status := remote.GetSettings().GetStatus().String()
	data.Status = types.StringValue(status)
	if status != models.ENABLED_SECURITYRULESTATUS.String() &&
		status != models.DISABLED_SECURITYRULESTATUS.String() {
		return fmt.Errorf(
			"%w: unsupported status %q; enabled cannot be determined",
			errInvalidResponse,
			status,
		)
	}
	destinations := make([]attr.Value, 0, len(remote.GetMatchingConditions().GetDestinations()))
	for _, remoteDestination := range remote.GetMatchingConditions().GetDestinations() {
		var destinationType string
		var remoteValues []string
		switch destination := remoteDestination.(type) {
		case *models.TlsInspectionFqdnDestination:
			destinationType = destinationTypeFQDN
			remoteValues = destination.GetValues()
		case *models.TlsInspectionWebCategoryDestination:
			destinationType = destinationTypeWebCategory
			remoteValues = destination.GetValues()
		default:
			return fmt.Errorf(
				"%w: unsupported destination type %T",
				errInvalidResponse,
				remoteDestination,
			)
		}
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
				tlsInspectionPolicyRuleDestinationObjectType.AttrTypes,
				map[string]attr.Value{"type": types.StringValue(destinationType), "values": values},
			),
		)
	}
	data.ID = convert.GraphToFrameworkString(remote.GetId())
	data.Name = convert.GraphToFrameworkString(remote.GetName())
	data.Description = convert.GraphToFrameworkString(remote.GetDescription())
	data.Action = convert.GraphToFrameworkString(action)
	// #nosec G115 -- priority was validated against Terraform's supported range above.
	data.Priority = types.Int32Value(int32(priority))
	data.Enabled = types.BoolValue(status == models.ENABLED_SECURITYRULESTATUS.String())
	data.Destinations = types.ListValueMust(
		tlsInspectionPolicyRuleDestinationObjectType,
		destinations,
	)
	return nil
}
