package graphBetaDeviceManagementIntuneBrandingProfile

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the remote IntuneBrandingProfile resource state to Terraform state
func MapRemoteStateToTerraform(ctx context.Context, data *IntuneBrandingProfileResourceModel, remoteResource graphmodels.IntuneBrandingProfileable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceName": remoteResource.GetDisplayName(),
		"resourceId":   remoteResource.GetId(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.ProfileName = convert.GraphToFrameworkString(remoteResource.GetProfileName())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.ContactITName = convert.GraphToFrameworkString(remoteResource.GetContactITName())
	data.ContactITPhoneNumber = convert.GraphToFrameworkString(remoteResource.GetContactITPhoneNumber())
	data.ContactITEmailAddress = convert.GraphToFrameworkString(remoteResource.GetContactITEmailAddress())
	data.ContactITNotes = convert.GraphToFrameworkString(remoteResource.GetContactITNotes())
	data.OnlineSupportSiteUrl = convert.GraphToFrameworkString(remoteResource.GetOnlineSupportSiteUrl())
	data.OnlineSupportSiteName = convert.GraphToFrameworkString(remoteResource.GetOnlineSupportSiteName())
	data.PrivacyUrl = convert.GraphToFrameworkString(remoteResource.GetPrivacyUrl())
	data.CustomPrivacyMessage = convert.GraphToFrameworkString(remoteResource.GetCustomPrivacyMessage())
	data.CustomCanSeePrivacyMessage = convert.GraphToFrameworkString(remoteResource.GetCustomCanSeePrivacyMessage())
	data.CustomCantSeePrivacyMessage = convert.GraphToFrameworkString(remoteResource.GetCustomCantSeePrivacyMessage())
	data.EnrollmentAvailability = convert.GraphToFrameworkEnum(remoteResource.GetEnrollmentAvailability())
	data.ShowLogo = convert.GraphToFrameworkBool(remoteResource.GetShowLogo())
	data.ShowDisplayNameNextToLogo = convert.GraphToFrameworkBool(remoteResource.GetShowDisplayNameNextToLogo())
	data.IsRemoveDeviceDisabled = convert.GraphToFrameworkBool(remoteResource.GetIsRemoveDeviceDisabled())
	data.IsFactoryResetDisabled = convert.GraphToFrameworkBool(remoteResource.GetIsFactoryResetDisabled())
	data.ShowAzureADEnterpriseApps = convert.GraphToFrameworkBool(remoteResource.GetShowAzureADEnterpriseApps())
	data.ShowOfficeWebApps = convert.GraphToFrameworkBool(remoteResource.GetShowOfficeWebApps())
	data.SendDeviceOwnershipChangePushNotification = convert.GraphToFrameworkBool(remoteResource.GetSendDeviceOwnershipChangePushNotification())
	data.DisableClientTelemetry = convert.GraphToFrameworkBool(remoteResource.GetDisableClientTelemetry())
	data.IsDefaultProfile = convert.GraphToFrameworkBool(remoteResource.GetIsDefaultProfile())
	data.CreatedDateTime = convert.GraphToFrameworkTime(remoteResource.GetCreatedDateTime())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, remoteResource.GetRoleScopeTagIds())

	if themeColor := remoteResource.GetThemeColor(); themeColor != nil {
		data.ThemeColor = &RgbColorResourceModel{}

		if r := themeColor.GetR(); r != nil {
			data.ThemeColor.R = types.Int32Value(int32(*r))
		}

		if g := themeColor.GetG(); g != nil {
			data.ThemeColor.G = types.Int32Value(int32(*g))
		}

		if b := themeColor.GetB(); b != nil {
			data.ThemeColor.B = types.Int32Value(int32(*b))
		}
	}

	if blockedActions := remoteResource.GetCompanyPortalBlockedActions(); len(blockedActions) > 0 {
		var actions []*CompanyPortalBlockedActionResourceModel
		for _, action := range blockedActions {
			actionModel := &CompanyPortalBlockedActionResourceModel{
				Platform:  convert.GraphToFrameworkEnum(action.GetPlatform()),
				OwnerType: convert.GraphToFrameworkEnum(action.GetOwnerType()),
				Action:    convert.GraphToFrameworkEnum(action.GetAction()),
			}
			actions = append(actions, actionModel)
		}
		data.CompanyPortalBlockedActions = actions
	} else {
		data.CompanyPortalBlockedActions = []*CompanyPortalBlockedActionResourceModel{}
	}

	// Note: We don't map the image properties back as they are input-only
	// The API doesn't return the actual image data, so we preserve the user's input

	tflog.Debug(ctx, fmt.Sprintf("Finished mapping resource %s with id %s", ResourceName, data.ID.ValueString()))
}
