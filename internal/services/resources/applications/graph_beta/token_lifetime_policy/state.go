package graphBetaApplicationsTokenLifetimePolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps a remote TokenLifetimePolicy to the Terraform state model
func MapRemoteResourceStateToTerraform(ctx context.Context, data *TokenLifetimePolicyResourceModel, remoteResource graphmodels.TokenLifetimePolicyable) {
	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.Definition = convert.GraphToFrameworkStringList(remoteResource.GetDefinition())
	data.IsOrganizationDefault = convert.GraphToFrameworkBool(remoteResource.GetIsOrganizationDefault())
	data.DeletedDateTime = convert.GraphToFrameworkTime(remoteResource.GetDeletedDateTime())
}
