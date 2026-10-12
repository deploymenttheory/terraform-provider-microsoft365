package graphBetaMobileAppDependency

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
	data *MobileAppDependencyResourceModel,
	graphResponse graphmodels.MobileAppDependencyable,
) diag.Diagnostics {
	var diags diag.Diagnostics

	tflog.Debug(ctx, fmt.Sprintf("Mapping %s resource to state", ResourceName))

	if graphResponse == nil || graphResponse.GetId() == nil || *graphResponse.GetId() == "" ||
		graphResponse.GetTargetId() == nil || *graphResponse.GetTargetId() == "" ||
		graphResponse.GetDependencyType() == nil {
		diags.AddError(
			"Incomplete app dependency response",
			"Microsoft Graph did not return the relationship ID, target ID, or dependency type. The existing state has been retained.",
		)
		return diags
	}
	if graphResponse.GetSourceId() != nil && *graphResponse.GetSourceId() != "" {
		data.SourceID = appIDToState(data.SourceID, graphResponse.GetSourceId())
	}
	if data.SourceID.IsNull() || data.SourceID.IsUnknown() || data.SourceID.ValueString() == "" {
		diags.AddError(
			"Incomplete app dependency response",
			"Microsoft Graph did not return the source app ID required to import this dependency.",
		)
		return diags
	}
	data.ID = convert.GraphToFrameworkString(graphResponse.GetId())
	data.TargetID = appIDToState(data.TargetID, graphResponse.GetTargetId())
	data.DependencyType = convert.GraphToFrameworkEnum(graphResponse.GetDependencyType())
	data.TargetType = convert.GraphToFrameworkEnum(graphResponse.GetTargetType())
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
	data.DependentAppCount = convert.GraphToFrameworkInt32(graphResponse.GetDependentAppCount())
	data.DependsOnAppCount = convert.GraphToFrameworkInt32(graphResponse.GetDependsOnAppCount())

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
