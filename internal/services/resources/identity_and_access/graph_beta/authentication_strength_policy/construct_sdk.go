package graphBetaAuthenticationStrengthPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// constructResourceForSDK converts the Terraform model to the SDK model for Create operations
func constructResourceForSDK(ctx context.Context, data *AuthenticationStrengthPolicyResourceModel) (graphmodels.AuthenticationStrengthPolicyable, error) {
	requestBody := graphmodels.NewAuthenticationStrengthPolicy()

	// Set basic fields using helpers
	convert.FrameworkToGraphString(data.DisplayName, requestBody.SetDisplayName)
	convert.FrameworkToGraphString(data.Description, requestBody.SetDescription)

	// Set allowedCombinations
	if !data.AllowedCombinations.IsNull() && !data.AllowedCombinations.IsUnknown() {
		var combinations []string
		if err := convert.FrameworkToGraphStringSet(ctx, data.AllowedCombinations, func(values []string) {
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
		odataType := config.ODataType.ValueString()

		switch odataType {
		case "#microsoft.graph.fido2CombinationConfiguration":
			fido2Config := graphmodels.NewFido2CombinationConfiguration()

			// Set appliesToCombinations
			if err := setAppliesToCombinations(ctx, config, fido2Config); err != nil {
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

			result = append(result, fido2Config)

		case "#microsoft.graph.x509CertificateCombinationConfiguration":
			x509Config := graphmodels.NewX509CertificateCombinationConfiguration()

			// Set appliesToCombinations
			if err := setAppliesToCombinations(ctx, config, x509Config); err != nil {
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

			result = append(result, x509Config)

		default:
			return nil, fmt.Errorf("unsupported odata type: %s", odataType)
		}
	}

	return result, nil
}

// setAppliesToCombinations is a helper to set appliesToCombinations on any config type
func setAppliesToCombinations(ctx context.Context, config CombinationConfigurationModel, sdkConfig graphmodels.AuthenticationCombinationConfigurationable) error {
	var combinations []string
	if err := convert.FrameworkToGraphStringSet(ctx, config.AppliesToCombinations, func(values []string) {
		combinations = values
	}); err != nil {
		return fmt.Errorf("failed to convert appliesToCombinations: %w", err)
	}

	// Convert to AuthenticationMethodModes enum
	authMethodModes := make([]graphmodels.AuthenticationMethodModes, 0, len(combinations))
	for _, combo := range combinations {
		parsedMode, err := graphmodels.ParseAuthenticationMethodModes(combo)
		if err != nil {
			return fmt.Errorf("invalid authentication method mode: %s: %w", combo, err)
		}
		// ParseAuthenticationMethodModes returns *AuthenticationMethodModes as any
		if modePtr, ok := parsedMode.(*graphmodels.AuthenticationMethodModes); ok {
			authMethodModes = append(authMethodModes, *modePtr)
		} else {
			return fmt.Errorf("failed to convert parsed mode to AuthenticationMethodModes: %s (got type %T)", combo, parsedMode)
		}
	}
	sdkConfig.SetAppliesToCombinations(authMethodModes)
	return nil
}
