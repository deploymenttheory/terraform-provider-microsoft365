package graphBetaGroup

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model
func constructResource(ctx context.Context, data *GroupResourceModel) (graphmodels.Groupable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewGroup()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphString(data.MailNickname, requestBody.SetMailNickname)
	convert.FrameworkToGraphBool(data.MailEnabled, requestBody.SetMailEnabled)
	convert.FrameworkToGraphBool(data.SecurityEnabled, requestBody.SetSecurityEnabled)

	if err := convert.FrameworkToGraphStringSet(ctx, data.GroupTypes, requestBody.SetGroupTypes); err != nil {
		return nil, fmt.Errorf("failed to set group types: %s", err)
	}

	convert.FrameworkToGraphString(data.Visibility, requestBody.SetVisibility)
	convert.FrameworkToGraphBool(data.IsAssignableToRole, requestBody.SetIsAssignableToRole)
	convert.FrameworkToGraphString(data.MembershipRule, requestBody.SetMembershipRule)
	convert.FrameworkToGraphString(data.MembershipRuleProcessingState, requestBody.SetMembershipRuleProcessingState)
	convert.FrameworkToGraphString(data.PreferredDataLocation, requestBody.SetPreferredDataLocation)
	convert.FrameworkToGraphString(data.PreferredLanguage, requestBody.SetPreferredLanguage)
	convert.FrameworkToGraphString(data.Theme, requestBody.SetTheme)
	convert.FrameworkToGraphString(data.Classification, requestBody.SetClassification)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
