package graphBetaDirectoryRole

import (
	"context"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToDataSource maps a single DirectoryRoleable to a DirectoryRoleModel.
func MapRemoteStateToDataSource(ctx context.Context, role graphmodels.DirectoryRoleable) DirectoryRoleModel {
	if role == nil {
		tflog.Debug(ctx, "DirectoryRole is nil, returning empty model")
		return DirectoryRoleModel{}
	}

	return DirectoryRoleModel{
		ID:             convert.GraphToFrameworkString(role.GetId()),
		DisplayName:    convert.GraphToFrameworkString(role.GetDisplayName()),
		Description:    convert.GraphToFrameworkString(role.GetDescription()),
		RoleTemplateID: convert.GraphToFrameworkString(role.GetRoleTemplateId()),
	}
}
