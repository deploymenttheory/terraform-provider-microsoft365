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

// constructResourceForSDK converts the Terraform model to the SDK model for Create operations
func constructResourceForSDK(ctx context.Context, data *AuthenticationStrengthPolicyResourceModel) (graphmodels.AuthenticationStrengthPolicyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := graphmodels.NewAuthenticationStrengthPolicy()

	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)

	if !data.AllowedCombinations.IsNull() && !data.AllowedCombinations.IsUnknown() {
		var combinations []string
		if err := convert.FrameworkToGraphStringSet(ctx, data.AllowedCombinations, func(values []string) {
			combinations = values
		}); err != nil {
			return nil, fmt.Errorf("failed to convert allowedCombinations: %w", err)
		}

		authMethodModes := make([]graphmodels.AuthenticationMethodModes, 0, len(combinations))
		for _, combin := range combinations {
			parsedMode, err := graphmodels.ParseAuthenticationMethodModes(combin)
			if err != nil {
				return nil, fmt.Errorf("invalid authentication method mode: %s: %w", combin, err)
			}
			// ParseAuthenticationMethodModes returns *AuthenticationMethodModes as any
			if modePtr, ok := parsedMode.(*graphmodels.AuthenticationMethodModes); ok {
				authMethodModes = append(authMethodModes, *modePtr)
			} else {
				return nil, fmt.Errorf("failed to convert parsed mode to AuthenticationMethodModes: %s (got type %T)", combin, parsedMode)
			}
		}
		requestBody.SetAllowedCombinations(authMethodModes)
	}

	// Set combinationConfigurations
	if !data.CombinationConfigurations.IsNull() && !data.CombinationConfigurations.IsUnknown() {
		configs, err := constructCombinationConfigurationsForSDK(ctx, data)
		if err != nil {
			return nil, err
		}
		requestBody.SetCombinationConfigurations(configs)
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	return requestBody, nil
}

// constructCombinationConfigurationsForSDK converts combination configurations to SDK models
func constructCombinationConfigurationsForSDK(ctx context.Context, data *AuthenticationStrengthPolicyResourceModel) ([]graphmodels.AuthenticationCombinationConfigurationable, error) {
	var configurations []CombinationConfigurationModel
	diags := data.CombinationConfigurations.ElementsAs(ctx, &configurations, false)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to parse combination configurations: %v", diags)
	}

	result := make([]graphmodels.AuthenticationCombinationConfigurationable, 0, len(configurations))

	for _, config := range configurations {
		combin := config.AppliesToCombinations.ValueString()
		if combin == "" {
			return nil, fmt.Errorf("applies_to_combinations cannot be empty")
		}

		switch combin {
		case "fido2":
			tflog.Info(ctx, "Inferred odata_type as fido2CombinationConfiguration")
			fido2Config := graphmodels.NewFido2CombinationConfiguration()

			if err := setAppliesToCombinations(ctx, config, fido2Config); err != nil {
				return nil, err
			}

			if !config.AllowedAAGUIDs.IsNull() && !config.AllowedAAGUIDs.IsUnknown() {
				var aaguids []string
				if err := convert.FrameworkToGraphStringSet(ctx, config.AllowedAAGUIDs, func(values []string) {
					aaguids = values
				}); err != nil {
					return nil, fmt.Errorf("failed to convert allowedAAGUIDs: %w", err)
				}
				fido2Config.SetAllowedAAGUIDs(aaguids)
			}

			result = append(result, fido2Config)

		case "x509CertificateMultiFactor", "x509CertificateSingleFactor":
			tflog.Info(ctx, fmt.Sprintf("Inferred odata_type as x509CertificateCombinationConfiguration from: %s", combin))
			x509Config := graphmodels.NewX509CertificateCombinationConfiguration()

			if err := setAppliesToCombinations(ctx, config, x509Config); err != nil {
				return nil, err
			}

			if !config.AllowedIssuerSkis.IsNull() && !config.AllowedIssuerSkis.IsUnknown() {
				var skis []string
				if err := convert.FrameworkToGraphStringSet(ctx, config.AllowedIssuerSkis, func(values []string) {
					skis = values
				}); err != nil {
					return nil, fmt.Errorf("failed to convert allowedIssuerSkis: %w", err)
				}
				tflog.Debug(ctx, fmt.Sprintf("Setting allowedIssuerSkis for X509 config: %v", skis))
				x509Config.SetAllowedIssuerSkis(skis)

				// Automatically build allowedIssuers from allowedIssuerSkis by adding "CUSTOMIDENTIFIER:" prefix
				if len(skis) > 0 {
					issuers := make([]string, len(skis))
					for i, ski := range skis {
						issuers[i] = "CUSTOMIDENTIFIER:" + ski
					}
					tflog.Debug(ctx, fmt.Sprintf("Auto-generated allowedIssuers from SKIs for X509 config: %v", issuers))
					additionalData := map[string]any{
						"allowedIssuers": issuers,
					}
					x509Config.SetAdditionalData(additionalData)
				}
			}

			if !config.AllowedPolicyOIDs.IsNull() && !config.AllowedPolicyOIDs.IsUnknown() {
				var oids []string
				if err := convert.FrameworkToGraphStringSet(ctx, config.AllowedPolicyOIDs, func(values []string) {
					oids = values
				}); err != nil {
					return nil, fmt.Errorf("failed to convert allowedPolicyOIDs: %w", err)
				}
				tflog.Debug(ctx, fmt.Sprintf("Setting allowedPolicyOIDs for X509 config: %v", oids))
				x509Config.SetAllowedPolicyOIDs(oids)
			}

			tflog.Debug(ctx, fmt.Sprintf("Appending X509 config to result, current result length: %d", len(result)))
			result = append(result, x509Config)

		default:
			return nil, fmt.Errorf("invalid applies_to_combinations value '%s': must be 'fido2', 'x509CertificateMultiFactor', or 'x509CertificateSingleFactor'", combin)
		}
	}

	return result, nil
}

// setAppliesToCombinations is a helper to set appliesToCombinations on any config type
func setAppliesToCombinations(ctx context.Context, config CombinationConfigurationModel, sdkConfig graphmodels.AuthenticationCombinationConfigurationable) error {
	combin := config.AppliesToCombinations.ValueString()
	if combin == "" {
		return fmt.Errorf("applies_to_combinations cannot be empty")
	}

	parsedMode, err := graphmodels.ParseAuthenticationMethodModes(combin)
	if err != nil {
		return fmt.Errorf("invalid authentication method mode: %s: %w", combin, err)
	}

	if modePtr, ok := parsedMode.(*graphmodels.AuthenticationMethodModes); ok {
		authMethodModes := []graphmodels.AuthenticationMethodModes{*modePtr}
		sdkConfig.SetAppliesToCombinations(authMethodModes)
		return nil
	}

	return fmt.Errorf("failed to convert parsed mode to AuthenticationMethodModes: %s (got type %T)", combin, parsedMode)
}

// constructAllowedCombinationsUpdateForSDK constructs the SDK request body for updating allowedCombinations
func constructAllowedCombinationsUpdateForSDK(ctx context.Context, plan *AuthenticationStrengthPolicyResourceModel) (*identity.ConditionalAccessAuthenticationStrengthPoliciesItemUpdateAllowedCombinationsPostRequestBody, error) {
	requestBody := identity.NewConditionalAccessAuthenticationStrengthPoliciesItemUpdateAllowedCombinationsPostRequestBody()

	var combinations []string
	if err := convert.FrameworkToGraphStringSet(ctx, plan.AllowedCombinations, func(values []string) {
		combinations = values
	}); err != nil {
		return nil, fmt.Errorf("failed to convert allowedCombinations: %w", err)
	}

	authMethodModes := make([]graphmodels.AuthenticationMethodModes, 0, len(combinations))
	for _, combin := range combinations {
		parsedMode, err := graphmodels.ParseAuthenticationMethodModes(combin)
		if err != nil {
			return nil, fmt.Errorf("invalid authentication method mode: %s: %w", combin, err)
		}
		if modePtr, ok := parsedMode.(*graphmodels.AuthenticationMethodModes); ok {
			authMethodModes = append(authMethodModes, *modePtr)
		} else {
			return nil, fmt.Errorf("failed to convert parsed mode to AuthenticationMethodModes: %s (got type %T)", combin, parsedMode)
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
	combin := config.AppliesToCombinations.ValueString()
	if combin == "" {
		return nil, fmt.Errorf("applies_to_combinations cannot be empty")
	}

	switch combin {
	case "fido2":
		tflog.Info(ctx, "Inferred odata_type as fido2CombinationConfiguration for update")
		fido2Config := graphmodels.NewFido2CombinationConfiguration()

		if err := setAppliesToCombinations(ctx, *config, fido2Config); err != nil {
			return nil, err
		}

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

	case "x509CertificateMultiFactor", "x509CertificateSingleFactor":
		tflog.Info(ctx, fmt.Sprintf("Inferred odata_type as x509CertificateCombinationConfiguration for update from: %s", combin))
		x509Config := graphmodels.NewX509CertificateCombinationConfiguration()

		if err := setAppliesToCombinations(ctx, *config, x509Config); err != nil {
			return nil, err
		}

		if !config.AllowedIssuerSkis.IsNull() && !config.AllowedIssuerSkis.IsUnknown() {
			var skis []string
			if err := convert.FrameworkToGraphStringSet(ctx, config.AllowedIssuerSkis, func(values []string) {
				skis = values
			}); err != nil {
				return nil, fmt.Errorf("failed to convert allowedIssuerSkis: %w", err)
			}
			x509Config.SetAllowedIssuerSkis(skis)

			// Automatically build allowedIssuers from allowedIssuerSkis by adding "CUSTOMIDENTIFIER:" prefix
			if len(skis) > 0 {
				issuers := make([]string, len(skis))
				for i, ski := range skis {
					issuers[i] = "CUSTOMIDENTIFIER:" + ski
				}
				tflog.Debug(ctx, fmt.Sprintf("Auto-generated allowedIssuers from SKIs for X509 config update: %v", issuers))
				additionalData := map[string]any{
					"allowedIssuers": issuers,
				}
				x509Config.SetAdditionalData(additionalData)
			}
		}

		if !config.AllowedPolicyOIDs.IsNull() && !config.AllowedPolicyOIDs.IsUnknown() {
			var oids []string
			if err := convert.FrameworkToGraphStringSet(ctx, config.AllowedPolicyOIDs, func(values []string) {
				oids = values
			}); err != nil {
				return nil, fmt.Errorf("failed to convert allowedPolicyOIDs: %w", err)
			}
			x509Config.SetAllowedPolicyOIDs(oids)
		}

		if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for update X509 combination configuration for resource %s", ResourceName), x509Config); err != nil {
			tflog.Error(ctx, "Failed to debug log object", map[string]any{
				"error": err.Error(),
			})
		}

		return x509Config, nil

	default:
		return nil, fmt.Errorf("invalid applies_to_combinations value '%s': must be 'fido2', 'x509CertificateMultiFactor', or 'x509CertificateSingleFactor'", combin)
	}
}

// combinationConfigurationsEqual checks if two combination configurations are equal
func combinationConfigurationsEqual(ctx context.Context, a, b *CombinationConfigurationModel) bool {
	return a.ODataType.Equal(b.ODataType) &&
		a.AppliesToCombinations.Equal(b.AppliesToCombinations) &&
		a.AllowedAAGUIDs.Equal(b.AllowedAAGUIDs) &&
		a.AllowedIssuerSkis.Equal(b.AllowedIssuerSkis) &&
		a.AllowedPolicyOIDs.Equal(b.AllowedPolicyOIDs)
	// Note: AllowedIssuers is computed from AllowedIssuerSkis, so no need to compare it
}
