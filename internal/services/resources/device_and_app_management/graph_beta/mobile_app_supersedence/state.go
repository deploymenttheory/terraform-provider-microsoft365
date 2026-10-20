package graphBetaMobileAppSupersedence

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// mapResourceToState maps the Graph API response to the Terraform state.
func mapResourceToState(
	ctx context.Context,
	data *MobileAppSupersedenceResourceModel,
	graphResponse graphmodels.MobileAppSupersedenceable,
) diag.Diagnostics {
	var diags diag.Diagnostics

	tflog.Debug(ctx, fmt.Sprintf("Mapping %s resource to state", ResourceName))

	if graphResponse == nil || graphResponse.GetId() == nil || *graphResponse.GetId() == "" ||
		graphResponse.GetTargetId() == nil || *graphResponse.GetTargetId() == "" ||
		graphResponse.GetSupersedenceType() == nil {
		diags.AddError(
			"Incomplete app supersedence response",
			"Microsoft Graph did not return the relationship ID, target ID, or supersedence type. The existing state has been retained.",
		)
		return diags
	}
	if graphResponse.GetSourceId() != nil && *graphResponse.GetSourceId() != "" {
		data.SourceID = appIDToState(data.SourceID, graphResponse.GetSourceId())
	}
	if data.SourceID.IsNull() || data.SourceID.IsUnknown() || data.SourceID.ValueString() == "" {
		diags.AddError(
			"Incomplete app supersedence response",
			"Microsoft Graph did not return the source app ID required to import this supersedence.",
		)
		return diags
	}
	data.ID = convert.GraphToFrameworkString(graphResponse.GetId())
	data.TargetID = appIDToState(data.TargetID, graphResponse.GetTargetId())
	data.SupersedenceType = convert.GraphToFrameworkEnum(graphResponse.GetSupersedenceType())
	data.TargetType = convert.GraphToFrameworkEnum(graphResponse.GetTargetType())
	data.AppRelationshipType = data.TargetType
	data.SourceDisplayName = convert.GraphToFrameworkString(graphResponse.GetSourceDisplayName())
	data.SourceDisplayVersion = convert.GraphToFrameworkString(
		graphResponse.GetSourceDisplayVersion(),
	)
	data.SourcePublisherDisplayName = convert.GraphToFrameworkString(
		graphResponse.GetSourcePublisherDisplayName(),
	)
	data.TargetDisplayName = convert.GraphToFrameworkString(graphResponse.GetTargetDisplayName())
	data.TargetDisplayVersion = convert.GraphToFrameworkString(
		graphResponse.GetTargetDisplayVersion(),
	)
	data.TargetPublisher = convert.GraphToFrameworkString(graphResponse.GetTargetPublisher())
	data.TargetPublisherDisplayName = convert.GraphToFrameworkString(
		graphResponse.GetTargetPublisherDisplayName(),
	)
	data.SupersededAppCount = convert.GraphToFrameworkInt32(graphResponse.GetSupersededAppCount())
	data.SupersedingAppCount = convert.GraphToFrameworkInt32(graphResponse.GetSupersedingAppCount())

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping %s resource to state", ResourceName))

	return diags
}

// appIDToState preserves configured GUID casing when Graph returns the same app ID normalized.
func appIDToState(prior types.String, remote *string) types.String {
	if remote != nil && !prior.IsNull() && !prior.IsUnknown() &&
		strings.EqualFold(prior.ValueString(), *remote) {
		return prior
	}
	return convert.GraphToFrameworkString(remote)
}
