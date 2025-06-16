package graphBetaTermsAndConditions

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs and returns a TermsAndConditions
func constructResource(ctx context.Context, data TermsAndConditionsResourceModel) (graphmodels.TermsAndConditionsable, error) {
	tflog.Debug(ctx, "Starting terms and conditions construction")

	requestBody := graphmodels.NewTermsAndConditions()

	constructors.SetStringProperty(data.DisplayName, requestBody.SetDisplayName)
	constructors.SetStringProperty(data.Title, requestBody.SetTitle)
	constructors.SetStringProperty(data.BodyText, requestBody.SetBodyText)
	constructors.SetStringProperty(data.Description, requestBody.SetDescription)
	constructors.SetStringProperty(data.AcceptanceStatement, requestBody.SetAcceptanceStatement)
	constructors.SetInt32Property(data.Version, requestBody.SetVersion)

	if err := constructors.SetStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	if err := constructors.DebugLogGraphObject(ctx, "Constructed terms and conditions", requestBody); err != nil {
		tflog.Error(ctx, "Failed to log terms and conditions", map[string]interface{}{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}
