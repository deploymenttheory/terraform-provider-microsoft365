package graphBetaGroup

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/resources/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource maps the Terraform schema to the SDK model
func constructResource(ctx context.Context, data *GroupResourceModel) (graphmodels.Groupable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewGroup()

	constructors.SetStringProperty(data.DisplayName, requestBody.SetDisplayName)
	constructors.SetStringProperty(data.Description, requestBody.SetDescription)
	constructors.SetStringProperty(data.MailNickname, requestBody.SetMailNickname)

	if !data.MailEnabled.IsNull() && !data.MailEnabled.IsUnknown() {
		mailEnabled := data.MailEnabled.ValueBool()
		requestBody.SetMailEnabled(&mailEnabled)
	}

	if !data.SecurityEnabled.IsNull() && !data.SecurityEnabled.IsUnknown() {
		securityEnabled := data.SecurityEnabled.ValueBool()
		requestBody.SetSecurityEnabled(&securityEnabled)
	}

	if err := constructors.SetStringSet(ctx, data.GroupTypes, requestBody.SetGroupTypes); err != nil {
		return nil, fmt.Errorf("failed to set group types: %s", err)
	}

	constructors.SetStringProperty(data.Visibility, requestBody.SetVisibility)

	if !data.IsAssignableToRole.IsNull() && !data.IsAssignableToRole.IsUnknown() {
		isAssignableToRole := data.IsAssignableToRole.ValueBool()
		requestBody.SetIsAssignableToRole(&isAssignableToRole)
	}

	constructors.SetStringProperty(data.MembershipRule, requestBody.SetMembershipRule)
	constructors.SetStringProperty(data.MembershipRuleProcessingState, requestBody.SetMembershipRuleProcessingState)
	constructors.SetStringProperty(data.PreferredDataLocation, requestBody.SetPreferredDataLocation)
	constructors.SetStringProperty(data.PreferredLanguage, requestBody.SetPreferredLanguage)
	constructors.SetStringProperty(data.Theme, requestBody.SetTheme)
	constructors.SetStringProperty(data.Classification, requestBody.SetClassification)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
} 