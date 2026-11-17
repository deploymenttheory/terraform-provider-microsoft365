package graphBetaM365AdminSharepointSettings

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// MapRemoteResourceStateToTerraform maps the remote SharePoint settings from Kiota SDK to Terraform state.
// Graph does not return an id for this singleton, so the fixed singleton ID is always used.
func MapRemoteResourceStateToTerraform(
	ctx context.Context,
	data SharepointSettingsResourceModel,
	remoteResource graphmodels.SharepointSettingsable,
) SharepointSettingsResourceModel {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return data
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": singletonID,
	})

	data.ID = types.StringValue(singletonID)
	data.AllowedDomainGuidsForSyncApp = mapAllowedDomainGuidsForSyncApp(ctx, remoteResource)
	data.AvailableManagedPathsForSiteCreation = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, remoteResource.GetAvailableManagedPathsForSiteCreation())
	data.DeletedUserPersonalSiteRetentionPeriodInDays = convert.GraphToFrameworkInt32(remoteResource.GetDeletedUserPersonalSiteRetentionPeriodInDays())
	data.ExcludedFileExtensionsForSyncApp = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, remoteResource.GetExcludedFileExtensionsForSyncApp())
	data.IdleSessionSignOut = mapIdleSessionSignOut(remoteResource.GetIdleSessionSignOut())
	data.ImageTaggingOption = convert.GraphToFrameworkEnum(remoteResource.GetImageTaggingOption())
	data.IsCommentingOnSitePagesEnabled = convert.GraphToFrameworkBool(remoteResource.GetIsCommentingOnSitePagesEnabled())
	data.IsFileActivityNotificationEnabled = convert.GraphToFrameworkBool(remoteResource.GetIsFileActivityNotificationEnabled())
	data.IsLegacyAuthProtocolsEnabled = convert.GraphToFrameworkBool(remoteResource.GetIsLegacyAuthProtocolsEnabled())
	data.IsLoopEnabled = convert.GraphToFrameworkBool(remoteResource.GetIsLoopEnabled())
	data.IsMacSyncAppEnabled = convert.GraphToFrameworkBool(remoteResource.GetIsMacSyncAppEnabled())
	data.IsRequireAcceptingUserToMatchInvitedUserEnabled = convert.GraphToFrameworkBool(remoteResource.GetIsRequireAcceptingUserToMatchInvitedUserEnabled())
	data.IsResharingByExternalUsersEnabled = convert.GraphToFrameworkBool(remoteResource.GetIsResharingByExternalUsersEnabled())
	data.IsSharePointMobileNotificationEnabled = convert.GraphToFrameworkBool(remoteResource.GetIsSharePointMobileNotificationEnabled())
	data.IsSharePointNewsfeedEnabled = convert.GraphToFrameworkBool(remoteResource.GetIsSharePointNewsfeedEnabled())
	data.IsSiteCreationEnabled = convert.GraphToFrameworkBool(remoteResource.GetIsSiteCreationEnabled())
	data.IsSiteCreationUIEnabled = convert.GraphToFrameworkBool(remoteResource.GetIsSiteCreationUIEnabled())
	data.IsSitePagesCreationEnabled = convert.GraphToFrameworkBool(remoteResource.GetIsSitePagesCreationEnabled())
	data.IsSitesStorageLimitAutomatic = convert.GraphToFrameworkBool(remoteResource.GetIsSitesStorageLimitAutomatic())
	data.IsSyncButtonHiddenOnPersonalSite = convert.GraphToFrameworkBool(remoteResource.GetIsSyncButtonHiddenOnPersonalSite())
	data.IsUnmanagedSyncAppForTenantRestricted = convert.GraphToFrameworkBool(remoteResource.GetIsUnmanagedSyncAppForTenantRestricted())
	data.PersonalSiteDefaultStorageLimitInMB = convert.GraphToFrameworkInt64(remoteResource.GetPersonalSiteDefaultStorageLimitInMB())
	data.SharingAllowedDomainList = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, remoteResource.GetSharingAllowedDomainList())
	data.SharingBlockedDomainList = convert.GraphToFrameworkStringSetPreserveEmpty(ctx, remoteResource.GetSharingBlockedDomainList())
	data.SharingCapability = convert.GraphToFrameworkEnum(remoteResource.GetSharingCapability())
	data.SharingDomainRestrictionMode = convert.GraphToFrameworkEnum(remoteResource.GetSharingDomainRestrictionMode())
	data.SiteCreationDefaultManagedPath = convert.GraphToFrameworkString(remoteResource.GetSiteCreationDefaultManagedPath())
	data.SiteCreationDefaultStorageLimitInMB = convert.GraphToFrameworkInt32(remoteResource.GetSiteCreationDefaultStorageLimitInMB())
	data.TenantDefaultTimezone = convert.GraphToFrameworkString(remoteResource.GetTenantDefaultTimezone())

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
	return data
}

// mapAllowedDomainGuidsForSyncApp maps the trusted sync app domain GUIDs as lower-case GUID strings.
func mapAllowedDomainGuidsForSyncApp(ctx context.Context, remoteResource graphmodels.SharepointSettingsable) types.Set {
	domains := remoteResource.GetAllowedDomainGuidsForSyncApp()
	guids := make([]string, 0, len(domains))
	for _, domain := range domains {
		guids = append(guids, domain.String())
	}
	return convert.GraphToFrameworkStringSetPreserveEmpty(ctx, guids)
}

// mapIdleSessionSignOut maps the idle session sign-out policy.
func mapIdleSessionSignOut(remote graphmodels.IdleSessionSignOutable) types.Object {
	if remote == nil {
		return types.ObjectNull(idleSessionSignOutAttrTypes)
	}
	return types.ObjectValueMust(idleSessionSignOutAttrTypes, map[string]attr.Value{
		"is_enabled":                convert.GraphToFrameworkBool(remote.GetIsEnabled()),
		"sign_out_after_in_seconds": convert.GraphToFrameworkInt64(remote.GetSignOutAfterInSeconds()),
		"warn_after_in_seconds":     convert.GraphToFrameworkInt64(remote.GetWarnAfterInSeconds()),
	})
}
