package graphBetaAdministrativeUnit

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource converts the Terraform resource model to a Kiota SDK model
// Returns an AdministrativeUnit that can be serialized by Kiota
func constructResource(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data *AdministrativeUnitResourceModel) (graphmodels.AdministrativeUnitable, error) {

	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	if err := validateRequest(ctx, client, data); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	requestBody := graphmodels.NewAdministrativeUnit()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphBool(data.IsMemberManagementRestricted, requestBody.SetIsMemberManagementRestricted)
	convert.FrameworkToGraphString(data.MembershipRule, requestBody.SetMembershipRule)
	convert.FrameworkToGraphString(data.MembershipRuleProcessingState, requestBody.SetMembershipRuleProcessingState)
	convert.FrameworkToGraphString(data.MembershipType, requestBody.SetMembershipType)
	convert.FrameworkToGraphString(data.Visibility, requestBody.SetVisibility)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
