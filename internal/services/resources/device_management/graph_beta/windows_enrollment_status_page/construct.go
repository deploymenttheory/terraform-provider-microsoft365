package graphBetaWindowsEnrollmentStatusPage

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResource constructs a Windows10EnrollmentCompletionPageConfiguration from the Terraform model
func constructResource(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, plan *WindowsEnrollmentStatusPageResourceModel) (graphmodels.DeviceEnrollmentConfigurationable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s from plan", ResourceName))

	// Validate the request data
	if err := validateRequest(ctx, client, plan); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	requestBody := graphmodels.NewWindows10EnrollmentCompletionPageConfiguration()

	convert.FrameworkToGraphString(plan.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(plan.Description, requestBody.SetDescription)
	convert.FrameworkToGraphBool(plan.ShowInstallationProgress, requestBody.SetShowInstallationProgress)
	convert.FrameworkToGraphBool(plan.OnlyShowPageToDevicesProvisionedByOutOfBoxExperienceOobe, requestBody.SetDisableUserStatusTrackingAfterFirstUser) // using custom key name as dev was getting very confusing.
	convert.FrameworkToGraphBool(plan.AllowDeviceResetOnInstallFailure, requestBody.SetAllowDeviceResetOnInstallFailure)
	convert.FrameworkToGraphBool(plan.AllowLogCollectionOnInstallFailure, requestBody.SetAllowLogCollectionOnInstallFailure)
	convert.FrameworkToGraphString(plan.CustomErrorMessage, requestBody.SetCustomErrorMessage)
	convert.FrameworkToGraphInt32(plan.InstallProgressTimeoutInMinutes, requestBody.SetInstallProgressTimeoutInMinutes)
	convert.FrameworkToGraphBool(plan.AllowDeviceUseOnInstallFailure, requestBody.SetAllowDeviceUseOnInstallFailure)
	convert.FrameworkToGraphBool(plan.BlockDeviceUseUntilAllAppsAndProfilesAreInstalled, requestBody.SetBlockDeviceSetupRetryByUser) // using custom key name as dev was getting very confusing.
	convert.FrameworkToGraphBool(plan.OnlyFailSelectedBlockingAppsInTechnicianPhase, requestBody.SetAllowNonBlockingAppInstallation) // using custom key name as dev was getting very confusing.
	convert.FrameworkToGraphBool(plan.InstallQualityUpdates, requestBody.SetInstallQualityUpdates)

	// Set selected_mobile_app_ids to empty array by default so that we can add and remove apps in updates.
	BlockDeviceUseRequiredAppIds := []string{}
	requestBody.SetSelectedMobileAppIds(BlockDeviceUseRequiredAppIds)

	if !plan.SelectedMobileAppIds.IsNull() && !plan.SelectedMobileAppIds.IsUnknown() {
		if err := convert.FrameworkToGraphStringSet(ctx, plan.SelectedMobileAppIds, requestBody.SetSelectedMobileAppIds); err != nil {
			return nil, fmt.Errorf("failed to convert selected mobile app IDs: %v", err)
		}
	}

	if err := convert.FrameworkToGraphStringSet(ctx, plan.RoleScopeTagIds, requestBody.SetRoleScopeTagIds); err != nil {
		return nil, fmt.Errorf("failed to convert role scope tag IDs: %v", err)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s", ResourceName))

	return requestBody, nil
}
