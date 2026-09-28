package graphBetaIdentityAndAccessB2bManagementPolicy

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteResourceStateToTerraform maps a remote B2bManagementPolicy to the Terraform state model
func MapRemoteResourceStateToTerraform(ctx context.Context, data *B2bManagementPolicyResourceModel, remoteResource graphmodels.B2bManagementPolicyable) {
	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Definition = convert.GraphToFrameworkStringList(remoteResource.GetDefinition())
	data.IsOrganizationDefault = convert.GraphToFrameworkBool(remoteResource.GetIsOrganizationDefault())
}
