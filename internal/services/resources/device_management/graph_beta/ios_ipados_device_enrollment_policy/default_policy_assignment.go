package graphBetaIOSiPadOSDeviceEnrollmentPolicy

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/devicemanagement"

	errors "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/errors/kiota"
)

// extractDefaultIOSPolicyId derives the ios_ipados_device_enrollment_policy id from the raw id of
// the depOnboardingSettings' expanded defaultIosEnrollmentProfile, given the owning
// dep_onboarding_settings_id. The raw id ends in "{depOnboardingSettingsId}_{policyId}" - the same
// format used to address the enrollment profile for the setDefaultProfile action itself - with an
// optional legacy prefix in front (the macOS equivalent carries an "ECV2_" prefix), so everything
// up to and including "{depOnboardingSettingsId}_" is stripped.
func extractDefaultIOSPolicyId(rawProfileId string, depOnboardingSettingsId string) string {
	if _, after, found := strings.Cut(rawProfileId, depOnboardingSettingsId+"_"); found {
		return after
	}
	return rawProfileId
}

// setDefaultIOSProfile calls the setDefaultProfile action for the given iOS/iPadOS enrollment
// policy, making it the default iOS enrollment profile for the DEP token.
//
// API Call: POST /deviceManagement/depOnboardingSettings/{depOnboardingSettingsId}/enrollmentProfiles/{enrollmentProfileId}/setDefaultProfile
// Reference: https://learn.microsoft.com/en-us/graph/api/intune-enrollment-enrollmentprofile-setdefaultprofile?view=graph-rest-beta
func (r *IOSiPadOSDeviceEnrollmentPolicyResource) setDefaultIOSProfile(ctx context.Context, depOnboardingSettingsId, iosDeviceEnrollmentPolicyId string) error {
	enrollmentProfileId := fmt.Sprintf("%s_%s", depOnboardingSettingsId, iosDeviceEnrollmentPolicyId)

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

// setDefaultIOSProfileWithRetry calls setDefaultIOSProfile, retrying on transient errors
// (notably 404 "Profile not found") for a bounded period. This absorbs Microsoft Graph's
// eventual-consistency window between a newly created policy's `POST` response and the moment
// its enrollment profile becomes addressable by the setDefaultProfile action - observed in
// practice immediately after Create on the macOS equivalent.
func (r *IOSiPadOSDeviceEnrollmentPolicyResource) setDefaultIOSProfileWithRetry(ctx context.Context, depOnboardingSettingsId, iosDeviceEnrollmentPolicyId string) error {
	const (
		maxAttempts = 15
		retryDelay  = 4 * time.Second
	)

	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		lastErr = r.setDefaultIOSProfile(ctx, depOnboardingSettingsId, iosDeviceEnrollmentPolicyId)
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

// resolveIsDefaultPolicyAssignment determines whether the given iOS/iPadOS enrollment policy is
// currently the default iOS enrollment profile for its DEP token.
//
// API Call: GET /deviceManagement/depOnboardingSettings/{depOnboardingSettingsId}?$expand=defaultIosEnrollmentProfile
// Reference: https://learn.microsoft.com/en-us/graph/api/resources/intune-onboarding-deponboardingsetting?view=graph-rest-beta
func (r *IOSiPadOSDeviceEnrollmentPolicyResource) resolveIsDefaultPolicyAssignment(ctx context.Context, depOnboardingSettingsId, policyId string) (bool, error) {
	if depOnboardingSettingsId == "" || policyId == "" {
		return false, nil
	}

	settings, err := r.client.
		DeviceManagement().
		DepOnboardingSettings().
		ByDepOnboardingSettingId(depOnboardingSettingsId).
		Get(ctx, &devicemanagement.DepOnboardingSettingsDepOnboardingSettingItemRequestBuilderGetRequestConfiguration{
			QueryParameters: &devicemanagement.DepOnboardingSettingsDepOnboardingSettingItemRequestBuilderGetQueryParameters{
				Expand: []string{"defaultIosEnrollmentProfile"},
			},
		})
	if err != nil {
		return false, err
	}

	defaultProfile := settings.GetDefaultIosEnrollmentProfile()
	if defaultProfile == nil || defaultProfile.GetId() == nil {
		return false, nil
	}

	return extractDefaultIOSPolicyId(*defaultProfile.GetId(), depOnboardingSettingsId) == policyId, nil
}
