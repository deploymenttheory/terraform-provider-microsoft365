package graphBetaNetworkThreatIntelligencePolicyRule

import (
	"context"
	"fmt"
	"math"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	models "github.com/microsoftgraph/msgraph-beta-sdk-go/models/networkaccess"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

var destinationObjectType = types.ObjectType{AttrTypes: map[string]attr.Type{
	"type": types.StringType, "values": types.ListType{ElemType: types.StringType},
}}

type unknownRuleStatusError struct{ status string }

func (e *unknownRuleStatusError) Error() string {
	return fmt.Sprintf(
		"unsupported observed rule status %q; enabled cannot be inferred; retaining its previous value",
		e.status,
	)
}

// MapRemoteStateToTerraform preserves the API representation, including repeated FQDNs.
func MapRemoteStateToTerraform(
	ctx context.Context,
	data *NetworkThreatIntelligencePolicyRuleResourceModel,
	response models.PolicyRuleable,
) error {
	remote, ok := response.(models.ThreatIntelligenceRuleable)
	if !ok || remote == nil || remote.GetId() == nil || *remote.GetId() == "" ||
		remote.GetName() == nil ||
		remote.GetPriority() == nil ||
		remote.GetAction() == nil ||
		remote.GetSettings() == nil ||
		remote.GetSettings().GetStatus() == nil ||
		remote.GetMatchingConditions() == nil ||
		remote.GetMatchingConditions().GetSeverity() == nil ||
		remote.GetMatchingConditions().GetDestinations() == nil {
		return fmt.Errorf(
			"%w: invalid threat intelligence rule response: missing required fields or discriminator",
			errInvalidResponse,
		)
	}
	priority := *remote.GetPriority()
	if priority == 65000 {
		return fmt.Errorf(
			"%w: the service-managed default rule cannot be managed or imported independently; manage its parent policy instead",
			errInvalidResponse,
		)
	}
	if priority < 100 || priority > math.MaxInt32 {
		return fmt.Errorf("%w: unsupported rule priority %d", errInvalidResponse, priority)
	}
	status := remote.GetSettings().GetStatus().String()
	if status != "enabled" && status != "disabled" {
		return &unknownRuleStatusError{status: status}
	}
	if len(remote.GetMatchingConditions().GetDestinations()) > 1 {
		return fmt.Errorf(
			"%w: multiple FQDN destination groups are unsupported",
			errInvalidResponse,
		)
	}
	groups := make(
		[]ThreatIntelligencePolicyRuleDestinationModel,
		0,
		len(remote.GetMatchingConditions().GetDestinations()),
	)
	for _, destination := range remote.GetMatchingConditions().GetDestinations() {
		fqdn, ok := destination.(models.ThreatIntelligenceFqdnDestinationable)
		if !ok {
			return fmt.Errorf(
				"%w: unsupported destination type in API response",
				errInvalidResponse,
			)
		}
		if len(fqdn.GetValues()) == 0 {
			return fmt.Errorf("%w: missing FQDN destination values", errInvalidResponse)
		}
		values := convert.GraphToFrameworkStringList(fqdn.GetValues())
		groups = append(
			groups,
			ThreatIntelligencePolicyRuleDestinationModel{
				Type:   types.StringValue(destinationTypeFQDN),
				Values: values,
			},
		)
	}
	destinations, diags := types.ListValueFrom(ctx, destinationObjectType, groups)
	if diags.HasError() {
		return fmt.Errorf("%w: map destinations: %v", errInvalidResponse, diags)
	}
	data.ID = convert.GraphToFrameworkString(remote.GetId())
	data.Name = convert.GraphToFrameworkString(remote.GetName())
	data.Description = convert.GraphToFrameworkString(remote.GetDescription())
	data.Action = convert.GraphToFrameworkEnum(remote.GetAction())
	data.Priority = types.Int32Value(int32(priority))
	data.Enabled = types.BoolValue(status == "enabled")
	data.Status = types.StringValue(status)
	data.Severity = convert.GraphToFrameworkEnum(remote.GetMatchingConditions().GetSeverity())
	data.Destinations = destinations
	return nil
}
