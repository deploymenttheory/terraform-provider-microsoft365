package graphBetaAppleConfiguratorEnrollmentPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// MapRemoteStateToTerraform maps the EnrollmentProfile to Terraform state.
// depId, when non-empty, will be stated to DepOnboardingSettingsId.
// Note: This function preserves the existing timeouts from the current state.
func MapRemoteStateToTerraform(ctx context.Context, data *AppleConfiguratorEnrollmentPolicyResourceModel, profile graphmodels.EnrollmentProfileable, depId string) {
	if profile == nil {
		tflog.Debug(ctx, "Remote enrollmentProfile is nil")
		return
	}

	tflog.Debug(ctx, "Starting to map remote enrollmentProfile to Terraform state", map[string]any{
		"resourceId": convert.GraphToFrameworkString(profile.GetId()),
	})

	data.ID = convert.GraphToFrameworkString(profile.GetId())
	data.DisplayName = convert.GraphToFrameworkString(profile.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(profile.GetDescription())
	data.RequiresUserAuthentication = convert.GraphToFrameworkBool(profile.GetRequiresUserAuthentication())
	data.EnableAuthenticationViaCompanyPortal = convert.GraphToFrameworkBool(profile.GetEnableAuthenticationViaCompanyPortal())
	data.RequireCompanyPortalOnSetupAssistantEnrolledDevices = convert.GraphToFrameworkBool(profile.GetRequireCompanyPortalOnSetupAssistantEnrolledDevices())
	data.ConfigurationEndpointUrl = convert.GraphToFrameworkString(profile.GetConfigurationEndpointUrl())
	if depId != "" {
		data.DepOnboardingSettingsId = types.StringValue(depId)
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished stating resource %s with id %s", ResourceName, data.ID.ValueString()))
}
