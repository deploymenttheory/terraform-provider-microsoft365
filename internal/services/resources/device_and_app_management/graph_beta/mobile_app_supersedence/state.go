package graphBetaMobileAppSupersedence

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// mapResourceToState maps the Graph API response to the Terraform state.
func mapResourceToState(ctx context.Context, data *MobileAppSupersedenceResourceModel, graphResponse graphmodels.MobileAppSupersedenceable) diag.Diagnostics {
	var diags diag.Diagnostics

	if graphResponse == nil {
		return diags
	}

	tflog.Debug(ctx, fmt.Sprintf("Mapping %s resource to state", ResourceName))

	data.ID = convert.GraphToFrameworkString(graphResponse.GetId())
	data.TargetID = convert.GraphToFrameworkString(graphResponse.GetTargetId())
	data.TargetDisplayName = convert.GraphToFrameworkString(graphResponse.GetTargetDisplayName())
	data.TargetDisplayVersion = convert.GraphToFrameworkString(graphResponse.GetTargetDisplayVersion())
	data.TargetPublisher = convert.GraphToFrameworkString(graphResponse.GetTargetPublisher())
	data.TargetPublisherDisplayName = convert.GraphToFrameworkString(graphResponse.GetTargetPublisherDisplayName())
	data.SourceID = convert.GraphToFrameworkString(graphResponse.GetSourceId())
	data.SourceDisplayName = convert.GraphToFrameworkString(graphResponse.GetSourceDisplayName())
	data.SourceDisplayVersion = convert.GraphToFrameworkString(graphResponse.GetSourceDisplayVersion())
	data.SourcePublisherDisplayName = convert.GraphToFrameworkString(graphResponse.GetSourcePublisherDisplayName())
	data.TargetType = convert.GraphToFrameworkEnum(graphResponse.GetTargetType())
	data.SupersedenceType = convert.GraphToFrameworkEnum(graphResponse.GetSupersedenceType())
	data.SupersededAppCount = convert.GraphToFrameworkInt32(graphResponse.GetSupersededAppCount())
	data.SupersedingAppCount = convert.GraphToFrameworkInt32(graphResponse.GetSupersedingAppCount())

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping %s resource to state", ResourceName))

	return diags
}
