package graphBetaWindowsEnrollmentStatusPage

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the remote state to the Terraform state
func MapRemoteStateToTerraform(ctx context.Context, data *WindowsEnrollmentStatusPageResourceModel, remoteResource graphmodels.DeviceEnrollmentConfigurationable) {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Starting to map remote state to Terraform state for %s", ResourceName))

	// Cast to the specific type
	enrollmentConfig, ok := remoteResource.(graphmodels.Windows10EnrollmentCompletionPageConfigurationable)
	if !ok {
		tflog.Error(ctx, "Failed to cast remote resource to Windows10EnrollmentCompletionPageConfiguration")
		return
	}

	// Map basic fields using convert helpers
	data.ID = convert.GraphToFrameworkString(enrollmentConfig.GetId())
	data.DisplayName = convert.GraphToFrameworkString(enrollmentConfig.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(enrollmentConfig.GetDescription())

	// Map installation progress settings
	data.ShowInstallationProgress = convert.GraphToFrameworkBool(enrollmentConfig.GetShowInstallationProgress())
	data.BlockDeviceUseUntilAllAppsAndProfilesAreInstalled = convert.GraphToFrameworkBool(enrollmentConfig.GetBlockDeviceSetupRetryByUser())
	data.AllowDeviceResetOnInstallFailure = convert.GraphToFrameworkBool(enrollmentConfig.GetAllowDeviceResetOnInstallFailure())
	data.AllowLogCollectionOnInstallFailure = convert.GraphToFrameworkBool(enrollmentConfig.GetAllowLogCollectionOnInstallFailure())
	data.CustomErrorMessage = convert.GraphToFrameworkString(enrollmentConfig.GetCustomErrorMessage())
	data.InstallProgressTimeoutInMinutes = convert.GraphToFrameworkInt32(enrollmentConfig.GetInstallProgressTimeoutInMinutes())
	data.AllowDeviceUseOnInstallFailure = convert.GraphToFrameworkBool(enrollmentConfig.GetAllowDeviceUseOnInstallFailure())

	// Map autopilot and user tracking settings
	//data.TrackInstallProgressForAutopilotOnly = convert.GraphToFrameworkBool(enrollmentConfig.GetTrackInstallProgressForAutopilotOnly())
	data.OnlyShowPageToDevicesProvisionedByOutOfBoxExperienceOobe = convert.GraphToFrameworkBool(enrollmentConfig.GetDisableUserStatusTrackingAfterFirstUser())

	// Map app installation settings
	data.OnlyFailSelectedBlockingAppsInTechnicianPhase = convert.GraphToFrameworkBool(enrollmentConfig.GetAllowNonBlockingAppInstallation())
	data.InstallQualityUpdates = convert.GraphToFrameworkBool(enrollmentConfig.GetInstallQualityUpdates())

	// Map selected mobile app IDs
	data.SelectedMobileAppIds = convert.GraphToFrameworkStringSet(ctx, enrollmentConfig.GetSelectedMobileAppIds())

	// Map role scope tag IDs
	data.RoleScopeTagIds = convert.GraphToFrameworkStringSet(ctx, enrollmentConfig.GetRoleScopeTagIds())

	assignments := remoteResource.GetAssignments()

	if len(assignments) == 0 {
		data.Assignments = types.SetNull(WindowsEnrollmentStatusPageAssignmentType())
	} else {
		MapAssignmentsToTerraform(ctx, data, assignments)
		tflog.Debug(ctx, "Completed assignment stating process", map[string]interface{}{
			"resourceId": data.ID.ValueString(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))
}
