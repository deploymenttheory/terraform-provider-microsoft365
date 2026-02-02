package graphBetaApplicationFederatedIdentityCredential

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func constructResource(ctx context.Context, data *ApplicationFederatedIdentityCredentialResourceModel) (graphmodels.FederatedIdentityCredentialable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewFederatedIdentityCredential()

	convert.FrameworkToGraphString(data.Name, requestBody.SetName)
	convert.FrameworkToGraphString(data.Issuer, requestBody.SetIssuer)
	convert.FrameworkToGraphString(data.Subject, requestBody.SetSubject)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)

	if !data.Audiences.IsNull() && !data.Audiences.IsUnknown() {
		var audiences []string
		diags := data.Audiences.ElementsAs(ctx, &audiences, false)
		if diags.HasError() {
			return nil, fmt.Errorf("failed to extract audiences: %v", diags.Errors())
		}
		requestBody.SetAudiences(audiences)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
