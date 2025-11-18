package graphBetaAuthenticationStrengthPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/identity"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructAllowedCombinationsUpdateForSDK constructs the SDK request body for updating allowedCombinations
func constructAllowedCombinationsUpdateForSDK(ctx context.Context, plan *AuthenticationStrengthPolicyResourceModel) (*identity.ConditionalAccessAuthenticationStrengthPoliciesItemUpdateAllowedCombinationsPostRequestBody, error) {
	requestBody := identity.NewConditionalAccessAuthenticationStrengthPoliciesItemUpdateAllowedCombinationsPostRequestBody()

	var combinations []string
	if err := convert.FrameworkToGraphStringSet(ctx, plan.AllowedCombinations, func(values []string) {
		combinations = values
	}); err != nil {
		return nil, fmt.Errorf("failed to convert allowedCombinations: %w", err)
	}

	// Convert strings to AuthenticationMethodModes enum
	authMethodModes := make([]graphmodels.AuthenticationMethodModes, 0, len(combinations))
	for _, combo := range combinations {
		parsedMode, err := graphmodels.ParseAuthenticationMethodModes(combo)
		if err != nil {
			return nil, fmt.Errorf("invalid authentication method mode: %s: %w", combo, err)
		}
		// ParseAuthenticationMethodModes returns *AuthenticationMethodModes as any
		if modePtr, ok := parsedMode.(*graphmodels.AuthenticationMethodModes); ok {
			authMethodModes = append(authMethodModes, *modePtr)
		} else {
			return nil, fmt.Errorf("failed to convert parsed mode to AuthenticationMethodModes: %s (got type %T)", combo, parsedMode)
		}
	}
	requestBody.SetAllowedCombinations(authMethodModes)

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for update allowed combinations for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}

// constructCombinationConfigurationUpdateForSDK constructs the SDK request body for updating a single combination configuration
func constructCombinationConfigurationUpdateForSDK(ctx context.Context, config *CombinationConfigurationModel) (graphmodels.AuthenticationCombinationConfigurationable, error) {
	odataType := config.ODataType.ValueString()

	switch odataType {
	case "#microsoft.graph.fido2CombinationConfiguration":
		fido2Config := graphmodels.NewFido2CombinationConfiguration()

		// Set appliesToCombinations
		if err := setAppliesToCombinations(ctx, *config, fido2Config); err != nil {
			return nil, err
		}

		// Set allowedAAGUIDs
		if !config.AllowedAAGUIDs.IsNull() && !config.AllowedAAGUIDs.IsUnknown() {
			var aaguids []string
			if err := convert.FrameworkToGraphStringSet(ctx, config.AllowedAAGUIDs, func(values []string) {
				aaguids = values
			}); err != nil {
				return nil, fmt.Errorf("failed to convert allowedAAGUIDs: %w", err)
			}
			fido2Config.SetAllowedAAGUIDs(aaguids)
		}

		if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for update FIDO2 combination configuration for resource %s", ResourceName), fido2Config); err != nil {
			tflog.Error(ctx, "Failed to debug log object", map[string]any{
				"error": err.Error(),
			})
		}

		return fido2Config, nil

	case "#microsoft.graph.x509CertificateCombinationConfiguration":
		x509Config := graphmodels.NewX509CertificateCombinationConfiguration()

		// Set appliesToCombinations
		if err := setAppliesToCombinations(ctx, *config, x509Config); err != nil {
			return nil, err
		}

		// Set allowedIssuerSkis
		if !config.AllowedIssuerSkis.IsNull() && !config.AllowedIssuerSkis.IsUnknown() {
			var skis []string
			if err := convert.FrameworkToGraphStringSet(ctx, config.AllowedIssuerSkis, func(values []string) {
				skis = values
			}); err != nil {
				return nil, fmt.Errorf("failed to convert allowedIssuerSkis: %w", err)
			}
			x509Config.SetAllowedIssuerSkis(skis)
		}

		// Set allowedPolicyOIDs
		if !config.AllowedPolicyOIDs.IsNull() && !config.AllowedPolicyOIDs.IsUnknown() {
			var oids []string
			if err := convert.FrameworkToGraphStringSet(ctx, config.AllowedPolicyOIDs, func(values []string) {
				oids = values
			}); err != nil {
				return nil, fmt.Errorf("failed to convert allowedPolicyOIDs: %w", err)
			}
			x509Config.SetAllowedPolicyOIDs(oids)
		}

		// Set allowedIssuers using AdditionalData (not officially supported by SDK)
		if !config.AllowedIssuers.IsNull() && !config.AllowedIssuers.IsUnknown() {
			var issuers []string
			if err := convert.FrameworkToGraphStringSet(ctx, config.AllowedIssuers, func(values []string) {
				issuers = values
			}); err != nil {
				return nil, fmt.Errorf("failed to convert allowedIssuers: %w", err)
			}
			if len(issuers) > 0 {
				additionalData := map[string]any{
					"allowedIssuers": issuers,
				}
				x509Config.SetAdditionalData(additionalData)
			}
		}

		if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for update X509 combination configuration for resource %s", ResourceName), x509Config); err != nil {
			tflog.Error(ctx, "Failed to debug log object", map[string]any{
				"error": err.Error(),
			})
		}

		return x509Config, nil

	default:
		return nil, fmt.Errorf("unsupported odata type: %s", odataType)
	}
}
