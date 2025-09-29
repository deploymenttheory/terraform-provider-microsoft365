package graphBetaDeviceManagementIntuneBrandingProfile

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	sharedConstructors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors/graph_beta/device_management"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	helpers "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud/graph_beta/device_and_app_management"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructCreateResource maps minimal required fields for initial POST request
func constructCreateResource(ctx context.Context, data *IntuneBrandingProfileResourceModel) (graphmodels.IntuneBrandingProfileable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s create resource from model", ResourceName))

	requestBody := graphmodels.NewIntuneBrandingProfile()

	convert.FrameworkToGraphString(data.ProfileName, requestBody.SetProfileName)
	convert.FrameworkToGraphString(data.ProfileDescription, requestBody.SetProfileDescription)

	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for create resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s create resource", ResourceName))

	return requestBody, nil
}

// constructUpdateResource maps the full Terraform schema to the SDK model for PATCH request
func constructUpdateResource(ctx context.Context, data *IntuneBrandingProfileResourceModel) (graphmodels.IntuneBrandingProfileable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s update resource from model", ResourceName))

	requestBody := graphmodels.NewIntuneBrandingProfile()

	// Include all fields in PATCH based on working browser example
	convert.FrameworkToGraphString(data.ProfileName, requestBody.SetProfileName)
	convert.FrameworkToGraphString(data.ProfileDescription, requestBody.SetProfileDescription)
	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.ContactITName, requestBody.SetContactITName)
	convert.FrameworkToGraphString(data.ContactITPhoneNumber, requestBody.SetContactITPhoneNumber)
	convert.FrameworkToGraphString(data.ContactITEmailAddress, requestBody.SetContactITEmailAddress)
	convert.FrameworkToGraphString(data.ContactITNotes, requestBody.SetContactITNotes)
	convert.FrameworkToGraphString(data.OnlineSupportSiteUrl, requestBody.SetOnlineSupportSiteUrl)
	convert.FrameworkToGraphString(data.OnlineSupportSiteName, requestBody.SetOnlineSupportSiteName)
	convert.FrameworkToGraphString(data.PrivacyUrl, requestBody.SetPrivacyUrl)
	convert.FrameworkToGraphString(data.CustomPrivacyMessage, requestBody.SetCustomPrivacyMessage)
	convert.FrameworkToGraphString(data.CustomCanSeePrivacyMessage, requestBody.SetCustomCanSeePrivacyMessage)
	convert.FrameworkToGraphString(data.CustomCantSeePrivacyMessage, requestBody.SetCustomCantSeePrivacyMessage)

	if err := convert.FrameworkToGraphEnum(data.EnrollmentAvailability, graphmodels.ParseEnrollmentAvailabilityOptions, requestBody.SetEnrollmentAvailability); err != nil {
		return nil, err
	}

	convert.FrameworkToGraphBool(data.ShowLogo, requestBody.SetShowLogo)
	convert.FrameworkToGraphBool(data.ShowDisplayNameNextToLogo, requestBody.SetShowDisplayNameNextToLogo)
	convert.FrameworkToGraphBool(data.IsRemoveDeviceDisabled, requestBody.SetIsRemoveDeviceDisabled)
	convert.FrameworkToGraphBool(data.IsFactoryResetDisabled, requestBody.SetIsFactoryResetDisabled)
	convert.FrameworkToGraphBool(data.ShowAzureADEnterpriseApps, requestBody.SetShowAzureADEnterpriseApps)
	convert.FrameworkToGraphBool(data.ShowOfficeWebApps, requestBody.SetShowOfficeWebApps)
	convert.FrameworkToGraphBool(data.ShowConfigurationManagerApps, requestBody.SetShowConfigurationManagerApps)
	convert.FrameworkToGraphBool(data.DisableDeviceCategorySelection, requestBody.SetDisableDeviceCategorySelection)
	convert.FrameworkToGraphBool(data.SendDeviceOwnershipChangePushNotification, requestBody.SetSendDeviceOwnershipChangePushNotification)
	convert.FrameworkToGraphBool(data.DisableClientTelemetry, requestBody.SetDisableClientTelemetry)
	convert.FrameworkToGraphBool(data.IsDefaultProfile, requestBody.SetIsDefaultProfile)

	// Include roleScopeTagIds based on working browser example
	if err := convert.FrameworkToGraphStringSet(ctx, data.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to set role scope tags: %s", err)
	}

	// Convert int32 to byte
	if data.ThemeColor != nil {
		themeColor := graphmodels.NewRgbColor()

		if !data.ThemeColor.R.IsNull() {
			r := byte(data.ThemeColor.R.ValueInt32())
			themeColor.SetR(&r)
		}

		if !data.ThemeColor.G.IsNull() {
			g := byte(data.ThemeColor.G.ValueInt32())
			themeColor.SetG(&g)
		}

		if !data.ThemeColor.B.IsNull() {
			b := byte(data.ThemeColor.B.ValueInt32())
			themeColor.SetB(&b)
		}

		requestBody.SetThemeColor(themeColor)
	}

	if len(data.CompanyPortalBlockedActions) > 0 {
		var blockedActions []graphmodels.CompanyPortalBlockedActionable
		for _, action := range data.CompanyPortalBlockedActions {
			blockedAction := graphmodels.NewCompanyPortalBlockedAction()
			if err := convert.FrameworkToGraphEnum(action.Platform, graphmodels.ParseDevicePlatformType, blockedAction.SetPlatform); err != nil {
				return nil, err
			}
			if err := convert.FrameworkToGraphEnum(action.OwnerType, graphmodels.ParseOwnerType, blockedAction.SetOwnerType); err != nil {
				return nil, err
			}
			if err := convert.FrameworkToGraphEnum(action.Action, graphmodels.ParseCompanyPortalAction, blockedAction.SetAction); err != nil {
				return nil, err
			}
			blockedActions = append(blockedActions, blockedAction)
		}
		requestBody.SetCompanyPortalBlockedActions(blockedActions)
	}

	if data.ThemeColorLogo != nil {
		logo, tempFiles, err := sharedConstructors.ConstructImage(ctx, data.ThemeColorLogo)
		if err != nil {
			return nil, fmt.Errorf("failed to construct theme color logo: %w", err)
		}
		defer func() {
			for _, tempFile := range tempFiles {
				helpers.CleanupTempFile(ctx, tempFile)
			}
		}()
		requestBody.SetThemeColorLogo(logo)
	}

	if data.LightBackgroundLogo != nil {
		logo, tempFiles, err := sharedConstructors.ConstructImage(ctx, data.LightBackgroundLogo)
		if err != nil {
			return nil, fmt.Errorf("failed to construct light background logo: %w", err)
		}
		defer func() {
			for _, tempFile := range tempFiles {
				helpers.CleanupTempFile(ctx, tempFile)
			}
		}()
		requestBody.SetLightBackgroundLogo(logo)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for update resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s update resource", ResourceName))

	return requestBody, nil
}

// constructLandingPageImageResource creates a separate request body for landing page image upload
func constructLandingPageImageResource(ctx context.Context, data *IntuneBrandingProfileResourceModel) (graphmodels.IntuneBrandingProfileable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s landing page image resource from model", ResourceName))

	requestBody := graphmodels.NewIntuneBrandingProfile()

	if data.LandingPageCustomizedImage != nil {
		image, tempFiles, err := sharedConstructors.ConstructImage(ctx, data.LandingPageCustomizedImage)
		if err != nil {
			return nil, fmt.Errorf("failed to construct landing page image: %w", err)
		}
		defer func() {
			for _, tempFile := range tempFiles {
				helpers.CleanupTempFile(ctx, tempFile)
			}
		}()
		requestBody.SetLandingPageCustomizedImage(image)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for landing page image resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s landing page image resource", ResourceName))

	return requestBody, nil
}
