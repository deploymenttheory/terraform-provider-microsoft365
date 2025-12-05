package graphBetaAgentIdentityBlueprintFederatedIdentityCredential

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateFromSDK maps the Kiota SDK FederatedIdentityCredentialable to Terraform state
func MapRemoteStateFromSDK(ctx context.Context, data *AgentIdentityBlueprintFederatedIdentityCredentialResourceModel, credential graphmodels.FederatedIdentityCredentialable) {
	tflog.Debug(ctx, fmt.Sprintf("Mapping federated identity credential from SDK response for %s", ResourceName))

	data.ID = convert.GraphToFrameworkString(credential.GetId())
	data.Name = convert.GraphToFrameworkString(credential.GetName())
	data.Issuer = convert.GraphToFrameworkString(credential.GetIssuer())
	data.Subject = convert.GraphToFrameworkString(credential.GetSubject())
	data.Description = convert.GraphToFrameworkString(credential.GetDescription())
	data.Audiences = convert.GraphToFrameworkStringSet(ctx, credential.GetAudiences())

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s from SDK", ResourceName, data.ID.ValueString()))
}
