package graphBetaDirectoryRoleAssignment

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

func MapRemoteStateToTerraform(data *DirectoryRoleAssignmentResourceModel, remote graphmodels.UnifiedRoleAssignmentable) {
	data.ID = types.StringPointerValue(remote.GetId())
	data.RoleDefinitionID = types.StringPointerValue(remote.GetRoleDefinitionId())
	data.PrincipalID = types.StringPointerValue(remote.GetPrincipalId())
	data.DirectoryScopeID = types.StringPointerValue(remote.GetDirectoryScopeId())
}
