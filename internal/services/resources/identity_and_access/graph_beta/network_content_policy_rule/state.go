package graphBetaNetworkContentPolicyRule

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// MapRemoteStateToTerraform maps a file rule onto the Terraform state. Comma-separated activities and session types are
// split into sets, and every empty value maps to null.
func MapRemoteStateToTerraform(ctx context.Context, data *NetworkContentPolicyRuleResourceModel, remote *contentPolicyRuleResponse) {
	if remote == nil {
		return
	}
	splitValues := func(value *string) []string {
		if value == nil {
			return nil
		}
		var result []string
		for part := range strings.SplitSeq(*value, ",") {
			if trimmed := strings.TrimSpace(part); trimmed != "" {
				result = append(result, trimmed)
			}
		}
		return result
	}
	// Kiota skipped null array elements, so they are skipped here too.
	presentValues := func(values []*string) []string {
		var result []string
		for _, value := range values {
			if value != nil {
				result = append(result, *value)
			}
		}
		return result
	}
	var status *string
	if remote.Settings != nil {
		status = remote.Settings.Status
	}
	var activities, sessionTypes *string
	var contentTypes, textContentTypes []*string
	var remoteDestinations []*contentPolicyRuleDestinationResponse
	if conditions := remote.MatchingConditions; conditions != nil {
		if conditions.FileAttributes != nil {
			activities = conditions.FileAttributes.Activities
			contentTypes = conditions.FileAttributes.ContentTypes
			textContentTypes = conditions.FileAttributes.TextContentTypes
		}
		if conditions.Sources != nil {
			sessionTypes = conditions.Sources.SessionType
		}
		remoteDestinations = conditions.Destinations
	}

	data.ID = convert.GraphToFrameworkString(remote.ID)
	data.Name = convert.GraphToFrameworkString(remote.Name)
	data.Description = convert.GraphToFrameworkString(remote.Description)
	data.Action = convert.GraphToFrameworkString(remote.Action)
	data.Priority = convert.GraphToFrameworkInt64(remote.Priority)
	data.Status = convert.GraphToFrameworkString(status)
	data.Activities = convert.GraphToFrameworkStringSet(ctx, splitValues(activities))
	data.ContentTypes = convert.GraphToFrameworkStringSet(ctx, presentValues(contentTypes))
	data.TextContentTypes = convert.GraphToFrameworkStringSet(ctx, presentValues(textContentTypes))
	data.SessionTypes = convert.GraphToFrameworkStringSet(ctx, splitValues(sessionTypes))

	destinations := make([]attr.Value, 0, len(remoteDestinations))
	for _, destination := range remoteDestinations {
		if destination == nil {
			continue
		}
		destinationType := types.StringNull()
		if destination.ODataType != nil {
			switch *destination.ODataType {
			case filePolicyWebCategoryDestinationODataType:
				destinationType = types.StringValue(destinationTypeWebCategory)
			case filePolicyURLDestinationODataType:
				destinationType = types.StringValue(destinationTypeURL)
			default:
				destinationType = types.StringValue(destinationTypeFQDN)
			}
		}
		fields := map[string]attr.Value{"type": destinationType, "values": convert.GraphToFrameworkStringSet(ctx, presentValues(destination.Values))}
		destinations = append(destinations, types.ObjectValueMust(contentPolicyRuleDestinationObjectType.AttrTypes, fields))
	}
	data.Destinations = types.ListValueMust(contentPolicyRuleDestinationObjectType, destinations)
	tflog.Debug(ctx, fmt.Sprintf("Finished mapping %s with id %s", ResourceName, data.ID.ValueString()))
}
