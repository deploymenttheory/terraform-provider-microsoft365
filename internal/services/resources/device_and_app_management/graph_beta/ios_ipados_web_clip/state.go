package graphBetaIOSiPadOSWebClip

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	sharedmodels "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/shared_models/graph_beta/device_and_app_management"
	sharedstater "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/state/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// mapResourceToState maps the Graph API response to the Terraform state.
func mapResourceToState(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, graphResponse graphmodels.MobileAppable, data *IOSiPadOSWebClipResourceModel) error {
	iosWebClip, ok := graphResponse.(graphmodels.IosiPadOSWebClipable)
	if !ok {
		return fmt.Errorf("expected IosiPadOSWebClipable but got %T", graphResponse)
	}

	data.ID = convert.GraphToFrameworkString(iosWebClip.GetId())
	data.DisplayName = convert.GraphToFrameworkString(iosWebClip.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(iosWebClip.GetDescription())
	data.Publisher = convert.GraphToFrameworkString(iosWebClip.GetPublisher())
	data.InformationUrl = convert.GraphToFrameworkString(iosWebClip.GetInformationUrl())
	data.PrivacyInformationUrl = convert.GraphToFrameworkString(iosWebClip.GetPrivacyInformationUrl())
	data.Owner = convert.GraphToFrameworkString(iosWebClip.GetOwner())
	data.Developer = convert.GraphToFrameworkString(iosWebClip.GetDeveloper())
	data.Notes = convert.GraphToFrameworkString(iosWebClip.GetNotes())
	data.IsFeatured = convert.GraphToFrameworkBool(iosWebClip.GetIsFeatured())
	data.CreatedDateTime = convert.GraphToFrameworkTime(iosWebClip.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(iosWebClip.GetLastModifiedDateTime())
	data.PublishingState = convert.GraphToFrameworkEnum(iosWebClip.GetPublishingState())
	data.DependentAppCount = convert.GraphToFrameworkInt32(iosWebClip.GetDependentAppCount())
	data.IsAssigned = convert.GraphToFrameworkBool(iosWebClip.GetIsAssigned())
	data.SupersededAppCount = convert.GraphToFrameworkInt32(iosWebClip.GetSupersededAppCount())
	data.SupersedingAppCount = convert.GraphToFrameworkInt32(iosWebClip.GetSupersedingAppCount())
	data.UploadState = convert.GraphToFrameworkInt32(iosWebClip.GetUploadState())
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, iosWebClip.GetRoleScopeTagIds())
	data.AppUrl = convert.GraphToFrameworkString(iosWebClip.GetAppUrl())
	data.FullScreenEnabled = convert.GraphToFrameworkBool(iosWebClip.GetFullScreenEnabled())
	data.IgnoreManifestScope = convert.GraphToFrameworkBool(iosWebClip.GetIgnoreManifestScope())
	data.PreComposedIconEnabled = convert.GraphToFrameworkBool(iosWebClip.GetPreComposedIconEnabled())
	data.TargetApplicationBundleIdentifier = convert.GraphToFrameworkString(iosWebClip.GetTargetApplicationBundleIdentifier())
	data.UseManagedBrowser = convert.GraphToFrameworkBool(iosWebClip.GetUseManagedBrowser())

	if data.AppIcon != nil {
		tflog.Debug(ctx, "Preserving original app_icon values from configuration")
	} else if largeIcon := iosWebClip.GetLargeIcon(); largeIcon != nil {
		data.AppIcon = &sharedmodels.MobileAppIconResourceModel{
			IconFilePathSource: types.StringNull(),
			IconURLSource:      types.StringNull(),
		}
	} else {
		data.AppIcon = nil
	}

	// Map categories
	data.Categories = sharedstater.MapMobileAppCategoriesStateToTerraform(ctx, iosWebClip.GetCategories())

	return nil
}
