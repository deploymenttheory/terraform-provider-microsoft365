package graphBetaAdminConsentRequestPolicy

import (
	"context"
	"errors"
	"fmt"

	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

var errReviewers = errors.New("invalid admin consent reviewer configuration")

// constructResource converts the Terraform resource model to a Microsoft Graph request body.
func constructResource(ctx context.Context, data *AdminConsentRequestPolicyResourceModel) (graphmodels.AdminConsentRequestPolicyable, error) {
	body := graphmodels.NewAdminConsentRequestPolicy()
	convert.FrameworkToGraphBool(data.IsEnabled, body.SetIsEnabled)
	convert.FrameworkToGraphBool(data.NotifyReviewers, body.SetNotifyReviewers)
	convert.FrameworkToGraphBool(data.RemindersEnabled, body.SetRemindersEnabled)
	convert.FrameworkToGraphInt32(data.RequestDurationInDays, body.SetRequestDurationInDays)
	var selections []ReviewerModel
	if diags := data.Reviewers.ElementsAs(ctx, &selections, false); diags.HasError() {
		return nil, fmt.Errorf("%w: %s", errReviewers, diags.Errors()[0].Detail())
	}
	reviewers := make([]graphmodels.AccessReviewReviewerScopeable, 0, len(selections))
	for _, selection := range selections {
		reviewer := graphmodels.NewAccessReviewReviewerScope()
		convert.FrameworkToGraphString(selection.Query, reviewer.SetQuery)
		convert.FrameworkToGraphString(selection.QueryType, reviewer.SetQueryType)
		convert.FrameworkToGraphString(selection.QueryRoot, reviewer.SetQueryRoot)
		reviewers = append(reviewers, reviewer)
	}
	body.SetReviewers(reviewers)
	return body, nil
}
