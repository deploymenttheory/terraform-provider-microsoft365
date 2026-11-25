package graphBetaAdminConsentRequestPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// MapRemoteStateToTerraform maps the remote Microsoft Graph response to Terraform state.
func MapRemoteStateToTerraform(ctx context.Context, data *AdminConsentRequestPolicyResourceModel, remote graphmodels.AdminConsentRequestPolicyable) diag.Diagnostics {
	data.ID = types.StringValue(singletonID)
	data.IsEnabled = convert.GraphToFrameworkBool(remote.GetIsEnabled())
	data.NotifyReviewers = convert.GraphToFrameworkBool(remote.GetNotifyReviewers())
	data.RemindersEnabled = convert.GraphToFrameworkBool(remote.GetRemindersEnabled())
	data.RequestDurationInDays = convert.GraphToFrameworkInt32(remote.GetRequestDurationInDays())
	data.Version = convert.GraphToFrameworkInt32(remote.GetVersion())
	selections := make([]ReviewerModel, 0, len(remote.GetReviewers()))
	for _, reviewer := range remote.GetReviewers() {
		selections = append(
			selections,
			ReviewerModel{
				Query:     convert.GraphToFrameworkString(reviewer.GetQuery()),
				QueryType: convert.GraphToFrameworkString(reviewer.GetQueryType()),
				QueryRoot: convert.GraphToFrameworkString(reviewer.GetQueryRoot()),
			},
		)
	}
	var diags diag.Diagnostics
	data.Reviewers, diags = types.SetValueFrom(ctx, types.ObjectType{AttrTypes: reviewerAttributeTypes}, selections)
	return diags
}
