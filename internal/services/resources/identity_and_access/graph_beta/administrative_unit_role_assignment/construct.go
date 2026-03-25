package graphBetaAdministrativeUnitRoleAssignment

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource converts the Terraform resource model to a Kiota SDK model
// Returns a ScopedRoleMembership that can be serialized by Kiota
func constructResource(ctx context.Context, data *AdministrativeUnitRoleAssignmentResourceModel) (graphmodels.ScopedRoleMembershipable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewScopedRoleMembership()

	convert.FrameworkToGraphString(data.RoleID, requestBody.SetRoleId)

	roleMemberInfo := graphmodels.NewIdentity()
	convert.FrameworkToGraphString(data.RoleMemberID, roleMemberInfo.SetId)
	requestBody.SetRoleMemberInfo(roleMemberInfo)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}
