// REF: https://learn.microsoft.com/en-us/graph/api/resources/intune-enrollment-enrollmentprofile?view=graph-rest-beta
package graphBetaAppleConfiguratorEnrollmentPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// AppleConfiguratorEnrollmentPolicyResourceModel models a DEP (ABM/ASM) enrollment profile under a DEP onboarding setting
// Endpoint shape per reqs: POST /deviceManagement/depOnboardingSettings/{depId}/enrollmentProfiles
// Fields observed in reqs doc: displayName, description, requiresUserAuthentication,
// enableAuthenticationViaCompanyPortal, requireCompanyPortalOnSetupAssistantEnrolledDevices
// Computed from GET: id, configurationEndpointUrl
type AppleConfiguratorEnrollmentPolicyResourceModel struct {
	ID                                                  types.String   `tfsdk:"id"`
	DisplayName                                         types.String   `tfsdk:"display_name"`
	Description                                         types.String   `tfsdk:"description"`
	RequiresUserAuthentication                          types.Bool     `tfsdk:"requires_user_authentication"`
	EnableAuthenticationViaCompanyPortal                types.Bool     `tfsdk:"enable_authentication_via_company_portal"`
	RequireCompanyPortalOnSetupAssistantEnrolledDevices types.Bool     `tfsdk:"require_company_portal_on_setup_assistant_enrolled_devices"`
	ConfigurationEndpointUrl                            types.String   `tfsdk:"configuration_endpoint_url"`
	DepOnboardingSettingsId                             types.String   `tfsdk:"dep_onboarding_settings_id"` // Parent DEP onboarding settings id
	Timeouts                                            timeouts.Value `tfsdk:"timeouts"`
}
