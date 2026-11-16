package graphBetaM365AdminSharepointSettings

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

var (
	errReadIdleSessionSignOut = errors.New("failed to read idle_session_sign_out")
	errInvalidSyncAppDomain   = errors.New("invalid allowed_domain_guids_for_sync_app value")
)

// constructResource converts the Terraform resource model to a Kiota SDK model for the PATCH request.
// The settings need no lookups to build the request, so the client is unused; it is kept so the
// signature matches the other resources in this domain.
func constructResource(
	ctx context.Context,
	_ *msgraphbetasdk.GraphServiceClient,
	data *SharepointSettingsResourceModel,
) (graphmodels.SharepointSettingsable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody, err := constructRequestBody(ctx, data)
	if err != nil {
		return nil, err
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructRequestBody builds the PATCH body from the plan.
//
// PATCH /admin/sharepoint/settings keeps every property that is not in the request body, so only known
// (configured or carried over from state) values are sent. A value that is unknown or null is omitted
// instead of being sent as null, and an unset collection is omitted rather than sent as null (R4).
// A configured empty collection is sent as [] so that a list can be cleared.
func constructRequestBody(ctx context.Context, data *SharepointSettingsResourceModel) (graphmodels.SharepointSettingsable, error) {
	requestBody := graphmodels.NewSharepointSettings()

	if err := constructAllowedDomainGuidsForSyncApp(ctx, data.AllowedDomainGuidsForSyncApp, requestBody); err != nil {
		return nil, err
	}

	convert.FrameworkToGraphInt32(data.DeletedUserPersonalSiteRetentionPeriodInDays, requestBody.SetDeletedUserPersonalSiteRetentionPeriodInDays)

	if err := convert.FrameworkToGraphStringSetOmitUnset(ctx, data.ExcludedFileExtensionsForSyncApp, requestBody.SetExcludedFileExtensionsForSyncApp); err != nil {
		return nil, fmt.Errorf("failed to set excluded_file_extensions_for_sync_app: %w", err)
	}

	if err := constructIdleSessionSignOut(ctx, data.IdleSessionSignOut, requestBody); err != nil {
		return nil, err
	}

	if err := convert.FrameworkToGraphEnum(data.ImageTaggingOption, graphmodels.ParseImageTaggingChoice, requestBody.SetImageTaggingOption); err != nil {
		return nil, fmt.Errorf("failed to set image_tagging_option: %w", err)
	}

	convert.FrameworkToGraphBool(data.IsCommentingOnSitePagesEnabled, requestBody.SetIsCommentingOnSitePagesEnabled)
	convert.FrameworkToGraphBool(data.IsFileActivityNotificationEnabled, requestBody.SetIsFileActivityNotificationEnabled)
	convert.FrameworkToGraphBool(data.IsLegacyAuthProtocolsEnabled, requestBody.SetIsLegacyAuthProtocolsEnabled)
	convert.FrameworkToGraphBool(data.IsLoopEnabled, requestBody.SetIsLoopEnabled)
	convert.FrameworkToGraphBool(data.IsMacSyncAppEnabled, requestBody.SetIsMacSyncAppEnabled)
	convert.FrameworkToGraphBool(data.IsRequireAcceptingUserToMatchInvitedUserEnabled, requestBody.SetIsRequireAcceptingUserToMatchInvitedUserEnabled)
	convert.FrameworkToGraphBool(data.IsResharingByExternalUsersEnabled, requestBody.SetIsResharingByExternalUsersEnabled)
	convert.FrameworkToGraphBool(data.IsSharePointMobileNotificationEnabled, requestBody.SetIsSharePointMobileNotificationEnabled)
	convert.FrameworkToGraphBool(data.IsSharePointNewsfeedEnabled, requestBody.SetIsSharePointNewsfeedEnabled)
	convert.FrameworkToGraphBool(data.IsSiteCreationEnabled, requestBody.SetIsSiteCreationEnabled)
	convert.FrameworkToGraphBool(data.IsSiteCreationUIEnabled, requestBody.SetIsSiteCreationUIEnabled)
	convert.FrameworkToGraphBool(data.IsSitePagesCreationEnabled, requestBody.SetIsSitePagesCreationEnabled)
	convert.FrameworkToGraphBool(data.IsSitesStorageLimitAutomatic, requestBody.SetIsSitesStorageLimitAutomatic)
	convert.FrameworkToGraphBool(data.IsSyncButtonHiddenOnPersonalSite, requestBody.SetIsSyncButtonHiddenOnPersonalSite)
	convert.FrameworkToGraphBool(data.IsUnmanagedSyncAppForTenantRestricted, requestBody.SetIsUnmanagedSyncAppForTenantRestricted)

	convert.FrameworkToGraphInt64(data.PersonalSiteDefaultStorageLimitInMB, requestBody.SetPersonalSiteDefaultStorageLimitInMB)

	if err := convert.FrameworkToGraphStringSetOmitUnset(ctx, data.SharingAllowedDomainList, requestBody.SetSharingAllowedDomainList); err != nil {
		return nil, fmt.Errorf("failed to set sharing_allowed_domain_list: %w", err)
	}

	if err := convert.FrameworkToGraphStringSetOmitUnset(ctx, data.SharingBlockedDomainList, requestBody.SetSharingBlockedDomainList); err != nil {
		return nil, fmt.Errorf("failed to set sharing_blocked_domain_list: %w", err)
	}

	if err := convert.FrameworkToGraphEnum(data.SharingCapability, graphmodels.ParseSharingCapabilities, requestBody.SetSharingCapability); err != nil {
		return nil, fmt.Errorf("failed to set sharing_capability: %w", err)
	}

	err := convert.FrameworkToGraphEnum(data.SharingDomainRestrictionMode, graphmodels.ParseSharingDomainRestrictionMode, requestBody.SetSharingDomainRestrictionMode)
	if err != nil {
		return nil, fmt.Errorf("failed to set sharing_domain_restriction_mode: %w", err)
	}

	convert.FrameworkToGraphString(data.SiteCreationDefaultManagedPath, requestBody.SetSiteCreationDefaultManagedPath)
	convert.FrameworkToGraphInt32(data.SiteCreationDefaultStorageLimitInMB, requestBody.SetSiteCreationDefaultStorageLimitInMB)
	convert.FrameworkToGraphString(data.TenantDefaultTimezone, requestBody.SetTenantDefaultTimezone)

	return requestBody, nil
}

// constructAllowedDomainGuidsForSyncApp sets the trusted sync app domain GUIDs when the set is known.
func constructAllowedDomainGuidsForSyncApp(ctx context.Context, value types.Set, requestBody graphmodels.SharepointSettingsable) error {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	var guids []string
	if diags := value.ElementsAs(ctx, &guids, false); diags.HasError() {
		return fmt.Errorf("%w: %v", errInvalidSyncAppDomain, diags.Errors())
	}

	domains := make([]uuid.UUID, 0, len(guids))
	for _, guid := range guids {
		parsed, err := uuid.Parse(guid)
		if err != nil {
			return fmt.Errorf("%w %q: %w", errInvalidSyncAppDomain, guid, err)
		}
		domains = append(domains, parsed)
	}

	requestBody.SetAllowedDomainGuidsForSyncApp(domains)
	return nil
}

// constructIdleSessionSignOut sets the idle session sign-out policy with only its known values.
// An idle_session_sign_out object with no known values is omitted rather than sent as an empty object.
func constructIdleSessionSignOut(ctx context.Context, value types.Object, requestBody graphmodels.SharepointSettingsable) error {
	if value.IsNull() || value.IsUnknown() {
		return nil
	}

	var idleSessionSignOut IdleSessionSignOutModel
	if diags := value.As(ctx, &idleSessionSignOut, basetypes.ObjectAsOptions{UnhandledUnknownAsEmpty: true}); diags.HasError() {
		return fmt.Errorf("%w: %v", errReadIdleSessionSignOut, diags.Errors())
	}

	policy := graphmodels.NewIdleSessionSignOut()
	convert.FrameworkToGraphBool(idleSessionSignOut.IsEnabled, policy.SetIsEnabled)
	convert.FrameworkToGraphInt64(idleSessionSignOut.SignOutAfterInSeconds, policy.SetSignOutAfterInSeconds)
	convert.FrameworkToGraphInt64(idleSessionSignOut.WarnAfterInSeconds, policy.SetWarnAfterInSeconds)

	if policy.GetIsEnabled() != nil || policy.GetSignOutAfterInSeconds() != nil || policy.GetWarnAfterInSeconds() != nil {
		requestBody.SetIdleSessionSignOut(policy)
	}

	return nil
}
