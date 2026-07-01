package graphBetaMacOSDepEnrollmentProfile

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

var (
	errOnboardingSettingsNil = errors.New("depOnboardingSettings response is nil")
	errNoAdeToken            = errors.New(
		"no Apple ADE/ABM (or ASM) DEP onboarding token was found on this tenant; " +
			"add an Apple token in Intune, or set dep_onboarding_settings_id explicitly",
	)
)

// resolveDepOnboardingSettingsId determines the depOnboardingSetting id (the Apple ABM/ASM
// ADE token) that owns macOS enrollment profiles.
//
// Resolution order:
//  1. If a value is provided in config/state, use it as-is (explicit escape hatch).
//  2. Otherwise list /deviceManagement/depOnboardingSettings and auto-select the Apple
//     ADE/ABM token (tokenType == dep) or ASM token (tokenType == appleSchoolManager).
//     Apple Configurator tokens (tokenType none) are ignored.
//
// If more than one Apple token exists, resolution is ambiguous and the caller must set
// dep_onboarding_settings_id explicitly.
func (r *MacOSDepEnrollmentProfileResource) resolveDepOnboardingSettingsId(
	ctx context.Context,
	provided types.String,
) (string, error) {
	if !provided.IsNull() && !provided.IsUnknown() && provided.ValueString() != "" {
		return provided.ValueString(), nil
	}

	settings, err := r.client.
		DeviceManagement().
		DepOnboardingSettings().
		Get(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("failed to GET /deviceManagement/depOnboardingSettings: %w", err)
	}
	if settings == nil {
		return "", errOnboardingSettingsNil
	}

	type candidate struct {
		id        string
		tokenName string
		tokenType string
	}
	var matches []candidate
	for _, s := range settings.GetValue() {
		if s == nil || s.GetId() == nil {
			continue
		}
		tt := s.GetTokenType()
		if tt == nil {
			continue
		}
		// Only Apple ADE/ABM (dep) and ASM tokens own macOS enrollment profiles.
		if *tt != graphmodels.DEP_DEPTOKENTYPE &&
			*tt != graphmodels.APPLESCHOOLMANAGER_DEPTOKENTYPE {
			continue
		}
		c := candidate{id: *s.GetId(), tokenType: tt.String()}
		if n := s.GetTokenName(); n != nil {
			c.tokenName = *n
		}
		matches = append(matches, c)
	}

	switch len(matches) {
	case 0:
		return "", errNoAdeToken
	case 1:
		return matches[0].id, nil
	default:
		var b strings.Builder
		for i, m := range matches {
			if i > 0 {
				b.WriteString(", ")
			}
			fmt.Fprintf(&b, "%s (%q, tokenType=%s)", m.id, m.tokenName, m.tokenType)
		}
		return "", fmt.Errorf(
			"%w: found %d Apple DEP tokens: [%s]",
			errAmbiguousAdeToken, len(matches), b.String(),
		)
	}
}

var errAmbiguousAdeToken = errors.New(
	"multiple Apple DEP tokens found; set dep_onboarding_settings_id explicitly to disambiguate",
)
