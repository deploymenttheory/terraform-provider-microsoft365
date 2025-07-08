package graphBetaWindowsWebApp

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_and_app_management"
	sharedstater "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/state/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// mapResourceToState maps the Graph API response to the Terraform state.
func mapResourceToState(ctx context.Context, data *WindowsWebAppResourceModel, graphResponse graphmodels.WindowsWebAppable) diag.Diagnostics {
	var diags diag.Diagnostics

	if graphResponse == nil {
		return diags
	}

	tflog.Debug(ctx, fmt.Sprintf("Mapping %s resource to state", ResourceName))

	data.ID = convert.GraphToFrameworkString(graphResponse.GetId())
	data.DisplayName = convert.GraphToFrameworkString(graphResponse.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(graphResponse.GetDescription())
	data.Publisher = convert.GraphToFrameworkString(graphResponse.GetPublisher())
	data.InformationUrl = convert.GraphToFrameworkString(graphResponse.GetInformationUrl())
	data.PrivacyInformationUrl = convert.GraphToFrameworkString(graphResponse.GetPrivacyInformationUrl())
	data.Owner = convert.GraphToFrameworkString(graphResponse.GetOwner())
	data.Developer = convert.GraphToFrameworkString(graphResponse.GetDeveloper())
	data.Notes = convert.GraphToFrameworkString(graphResponse.GetNotes())
	data.IsFeatured = convert.GraphToFrameworkBool(graphResponse.GetIsFeatured())
	data.CreatedDateTime = convert.GraphToFrameworkTime(graphResponse.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(graphResponse.GetLastModifiedDateTime())
	data.DependentAppCount = convert.GraphToFrameworkInt32(graphResponse.GetDependentAppCount())
	data.IsAssigned = convert.GraphToFrameworkBool(graphResponse.GetIsAssigned())
	data.SupersededAppCount = convert.GraphToFrameworkInt32(graphResponse.GetSupersededAppCount())
	data.SupersedingAppCount = convert.GraphToFrameworkInt32(graphResponse.GetSupersedingAppCount())
	data.UploadState = convert.GraphToFrameworkInt32(graphResponse.GetUploadState())
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, graphResponse.GetRoleScopeTagIds())
	data.AppUrl = convert.GraphToFrameworkString(graphResponse.GetAppUrl())
	data.PublishingState = convert.GraphToFrameworkEnum(graphResponse.GetPublishingState())

	if data.AppIcon != nil {
		tflog.Debug(ctx, "Preserving original app_icon values from configuration")
	} else if largeIcon := graphResponse.GetLargeIcon(); largeIcon != nil {
		data.AppIcon = &sharedmodels.MobileAppIconResourceModel{
			IconFilePathSource: types.StringNull(),
			IconURLSource:      types.StringNull(),
		}
	} else {
		data.AppIcon = nil
	}

	// Map categories
	data.Categories = sharedstater.MapMobileAppCategoriesStateToTerraform(ctx, graphResponse.GetCategories())

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping %s resource to state", ResourceName))

	return diags
}
