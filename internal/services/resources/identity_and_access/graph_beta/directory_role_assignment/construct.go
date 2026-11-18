package graphBetaDirectoryRoleAssignment

import graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

func constructResource(data *DirectoryRoleAssignmentResourceModel) graphmodels.UnifiedRoleAssignmentable {
	body := graphmodels.NewUnifiedRoleAssignment()
	body.SetRoleDefinitionId(data.RoleDefinitionID.ValueStringPointer())
	body.SetPrincipalId(data.PrincipalID.ValueStringPointer())
	body.SetDirectoryScopeId(new("/"))
	return body
}
