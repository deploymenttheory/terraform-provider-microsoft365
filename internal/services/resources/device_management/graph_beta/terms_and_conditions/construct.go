package graphBetaTermsAndConditions

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs and returns a TermsAndConditions
func constructResource(ctx context.Context, data TermsAndConditionsResourceModel) (graphmodels.TermsAndConditionsable, error) {
	tflog.Debug(ctx, "Starting terms and conditions construction")

	requestBody := graphmodels.NewTermsAndConditions()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Title, requestBody.SetTitle)
	convert.FrameworkToGraphString(data.BodyText, requestBody.SetBodyText)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)
	convert.FrameworkToGraphString(data.AcceptanceStatement, requestBody.SetAcceptanceStatement)
	convert.FrameworkToGraphInt32(data.Version, requestBody.SetVersion)

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	if err := constructors.DebugLogGraphObject(ctx, "Constructed terms and conditions", requestBody); err != nil {
		tflog.Error(ctx, "Failed to log terms and conditions", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}
