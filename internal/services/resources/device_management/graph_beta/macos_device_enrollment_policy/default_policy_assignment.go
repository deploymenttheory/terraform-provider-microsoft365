package graphBetaMacOSDeviceEnrollmentPolicy

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

// defaultMacOSProfileIdPrefix is the prefix Graph puts on the legacy depMacOSEnrollmentProfile id
// returned by depOnboardingSettings?$expand=defaultMacOsEnrollmentProfile. The suffix is
// "{depOnboardingSettingsId}_{macosDeviceEnrollmentPolicyId}", matching the id format used to
// address the enrollment profile for the setDefaultProfile action itself.
const defaultMacOSProfileIdPrefix = "ECV2_"

// extractDefaultMacOSPolicyId derives the macos_device_enrollment_policy id from the raw id of the
// depOnboardingSettings' expanded defaultMacOsEnrollmentProfile, given the owning
// dep_onboarding_settings_id.
func extractDefaultMacOSPolicyId(rawProfileId string, depOnboardingSettingsId string) string {
	withoutPrefix := strings.TrimPrefix(rawProfileId, defaultMacOSProfileIdPrefix)
	return strings.TrimPrefix(withoutPrefix, fmt.Sprintf("%s_", depOnboardingSettingsId))
}

// setDefaultMacOSProfile calls the setDefaultProfile action for the given macOS enrollment
// policy, making it the default macOS enrollment profile for the DEP token.
//
// API Call: POST /deviceManagement/depOnboardingSettings/{depOnboardingSettingsId}/enrollmentProfiles/{enrollmentProfileId}/setDefaultProfile
// Reference: https://learn.microsoft.com/en-us/graph/api/intune-enrollment-enrollmentprofile-setdefaultprofile?view=graph-rest-beta
func (r *MacOSDeviceEnrollmentPolicyResource) setDefaultMacOSProfile(ctx context.Context, depOnboardingSettingsId, macosDeviceEnrollmentPolicyId string) error {
	enrollmentProfileId := fmt.Sprintf("%s_%s", depOnboardingSettingsId, macosDeviceEnrollmentPolicyId)

	tflog.Debug(ctx, fmt.Sprintf("Calling setDefaultProfile for dep_onboarding_settings_id %s with enrollment profile id %s", depOnboardingSettingsId, enrollmentProfileId))

	return r.client.
		DeviceManagement().
		DepOnboardingSettings().
		ByDepOnboardingSettingId(depOnboardingSettingsId).
		EnrollmentProfiles().
		ByEnrollmentProfileId(enrollmentProfileId).
		SetDefaultProfile().
		Post(ctx, nil)
}

// setDefaultMacOSProfileWithRetry calls setDefaultMacOSProfile, retrying on transient errors
// (notably 404 "Profile not found") for a bounded period. This absorbs Microsoft Graph's
// eventual-consistency window between a newly created policy's `POST` response and the moment
// its enrollment profile becomes addressable by the setDefaultProfile action - observed in
// practice immediately after Create.
func (r *MacOSDeviceEnrollmentPolicyResource) setDefaultMacOSProfileWithRetry(ctx context.Context, depOnboardingSettingsId, macosDeviceEnrollmentPolicyId string) error {
	const (
		maxAttempts = 15
		retryDelay  = 4 * time.Second
	)

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		lastErr = r.setDefaultMacOSProfile(ctx, depOnboardingSettingsId, macosDeviceEnrollmentPolicyId)
		if lastErr == nil {
			return nil
		}

		errorInfo := errors.GraphError(ctx, lastErr)
		if !errors.IsRetryableReadError(&errorInfo) || attempt == maxAttempts {
			return lastErr
		}

		tflog.Debug(ctx, fmt.Sprintf("setDefaultProfile attempt %d/%d failed (likely awaiting Graph propagation), retrying in %s: %s",
			attempt, maxAttempts, retryDelay, lastErr.Error()))

		select {
		case <-time.After(retryDelay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return lastErr
}

// resolveIsDefaultPolicyAssignment determines whether the given macOS enrollment policy is
// currently the default macOS enrollment profile for its DEP token.
//
// API Call: GET /deviceManagement/depOnboardingSettings/{depOnboardingSettingsId}?$expand=defaultMacOsEnrollmentProfile
// Reference: https://learn.microsoft.com/en-us/graph/api/resources/intune-onboarding-deponboardingsetting?view=graph-rest-beta
func (r *MacOSDeviceEnrollmentPolicyResource) resolveIsDefaultPolicyAssignment(ctx context.Context, depOnboardingSettingsId, policyId string) (bool, error) {
	if depOnboardingSettingsId == "" || policyId == "" {
		return false, nil
	}

	settings, err := r.client.
		DeviceManagement().
		DepOnboardingSettings().
		ByDepOnboardingSettingId(depOnboardingSettingsId).
		Get(ctx, &devicemanagement.DepOnboardingSettingsDepOnboardingSettingItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &devicemanagement.DepOnboardingSettingsDepOnboardingSettingItemRequestBuilderGetQueryParameters{
				Expand: []string{"defaultMacOsEnrollmentProfile"},
			},
		})
	if err != nil {
		return false, err
	}

	defaultProfile := settings.GetDefaultMacOsEnrollmentProfile()
	if defaultProfile == nil || defaultProfile.GetId() == nil {
		return false, nil
	}

	return extractDefaultMacOSPolicyId(*defaultProfile.GetId(), depOnboardingSettingsId) == policyId, nil
}
