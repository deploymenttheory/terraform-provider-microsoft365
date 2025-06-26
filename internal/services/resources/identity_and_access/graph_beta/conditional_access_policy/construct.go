package graphBetaConditionalAccessPolicy

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// constructResource converts the Terraform resource model to a plain map for JSON marshaling
// Returns a map[string]interface{} that can be directly JSON marshaled by the HTTP client
func constructResource(ctx context.Context, data *ConditionalAccessPolicyResourceModel) (map[string]interface{}, error) {

	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := make(map[string]interface{})

	// Basic properties using convert helpers
	convert.FrameworkToGraphString(data.DisplayName, func(val *string) {
		if val != nil {
			requestBody["displayName"] = *val
		}
	})

	convert.FrameworkToGraphString(data.State, func(val *string) {
		if val != nil {
			requestBody["state"] = *val
		}
	})

	convert.FrameworkToGraphString(data.TemplateId, func(val *string) {
		if val != nil {
			requestBody["templateId"] = *val
		}
	})

	// Build conditions
	if data.Conditions != nil {
		conditions, err := constructConditions(ctx, data.Conditions)
		if err != nil {
			return nil, fmt.Errorf("failed to construct conditions: %w", err)
		}
		requestBody["conditions"] = conditions
	}

	// Build grant controls
	if data.GrantControls != nil {
		grantControls, err := constructGrantControls(ctx, data.GrantControls)
		if err != nil {
			return nil, fmt.Errorf("failed to construct grant controls: %w", err)
		}
		requestBody["grantControls"] = grantControls
	}

	// Build session controls
	if data.SessionControls != nil {
		sessionControls, err := constructSessionControls(ctx, data.SessionControls)
		if err != nil {
			return nil, fmt.Errorf("failed to construct session controls: %w", err)
		}
		requestBody["sessionControls"] = sessionControls
	}

	// Debug logging using plain JSON marshal
	if debugJSON, err := json.MarshalIndent(requestBody, "", "    "); err == nil {
		tflog.Debug(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), map[string]interface{}{
			"json": "\n" + string(debugJSON),
		})
	} else {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructConditions builds the conditions object using available Graph models
func constructConditions(ctx context.Context, data *ConditionalAccessConditions) (map[string]interface{}, error) {
	conditions := make(map[string]interface{})

	// Client app types using convert helper
	if err := convert.FrameworkToGraphStringSet(ctx, data.ClientAppTypes, func(values []string) {
		if len(values) > 0 {
			conditions["clientAppTypes"] = values
		}
	}); err != nil {
		return nil, fmt.Errorf("failed to convert client app types: %w", err)
	}

	// Applications
	if data.Applications != nil {
		applications := make(map[string]interface{})

		if err := convert.FrameworkToGraphStringSet(ctx, data.Applications.IncludeApplications, func(values []string) {
			if len(values) > 0 {
				applications["includeApplications"] = values
			}
		}); err != nil {
			return nil, fmt.Errorf("failed to convert include applications: %w", err)
		}

		if err := convert.FrameworkToGraphStringSet(ctx, data.Applications.ExcludeApplications, func(values []string) {
			// Always include excludeApplications if the field is configured (even if empty)
			if !data.Applications.ExcludeApplications.IsNull() {
				applications["excludeApplications"] = values
			}
		}); err != nil {
			return nil, fmt.Errorf("failed to convert exclude applications: %w", err)
		}

		if err := convert.FrameworkToGraphStringSet(ctx, data.Applications.IncludeUserActions, func(values []string) {
			// Always include includeUserActions if the field is configured (even if empty)
			if !data.Applications.IncludeUserActions.IsNull() {
				applications["includeUserActions"] = values
			}
		}); err != nil {
			return nil, fmt.Errorf("failed to convert include user actions: %w", err)
		}

		if err := convert.FrameworkToGraphStringSet(ctx, data.Applications.IncludeAuthenticationContextClassReferences, func(values []string) {
			// Always include includeAuthenticationContextClassReferences if the field is configured (even if empty)
			if !data.Applications.IncludeAuthenticationContextClassReferences.IsNull() {
				applications["includeAuthenticationContextClassReferences"] = values
			}
		}); err != nil {
			return nil, fmt.Errorf("failed to convert include auth context class refs: %w", err)
		}

		if data.Applications.ApplicationFilter != nil {
			appFilter := make(map[string]interface{})
			convert.FrameworkToGraphString(data.Applications.ApplicationFilter.Mode, func(value *string) {
				if value != nil {
					appFilter["mode"] = *value
				}
			})
			convert.FrameworkToGraphString(data.Applications.ApplicationFilter.Rule, func(value *string) {
				if value != nil {
					appFilter["rule"] = *value
				}
			})
			if len(appFilter) > 0 {
				applications["applicationFilter"] = appFilter
			}
		}

		if len(applications) > 0 {
			conditions["applications"] = applications
		}
	}

	// Users
	if data.Users != nil {
		users := make(map[string]interface{})

		if err := convert.FrameworkToGraphStringSet(ctx, data.Users.IncludeUsers, func(values []string) {
			if len(values) > 0 {
				users["includeUsers"] = values
			}
		}); err != nil {
			return nil, fmt.Errorf("failed to convert include users: %w", err)
		}

		if err := convert.FrameworkToGraphStringSet(ctx, data.Users.ExcludeUsers, func(values []string) {
			if len(values) > 0 {
				users["excludeUsers"] = values
			}
		}); err != nil {
			return nil, fmt.Errorf("failed to convert exclude users: %w", err)
		}

		if err := convert.FrameworkToGraphStringSet(ctx, data.Users.IncludeGroups, func(values []string) {
			if len(values) > 0 {
				users["includeGroups"] = values
			}
		}); err != nil {
			return nil, fmt.Errorf("failed to convert include groups: %w", err)
		}

		if err := convert.FrameworkToGraphStringSet(ctx, data.Users.ExcludeGroups, func(values []string) {
			if len(values) > 0 {
				users["excludeGroups"] = values
			}
		}); err != nil {
			return nil, fmt.Errorf("failed to convert exclude groups: %w", err)
		}

		if err := convert.FrameworkToGraphStringSet(ctx, data.Users.IncludeRoles, func(values []string) {
			if len(values) > 0 {
				users["includeRoles"] = values
			}
		}); err != nil {
			return nil, fmt.Errorf("failed to convert include roles: %w", err)
		}

		if err := convert.FrameworkToGraphStringSet(ctx, data.Users.ExcludeRoles, func(values []string) {
			if len(values) > 0 {
				users["excludeRoles"] = values
			}
		}); err != nil {
			return nil, fmt.Errorf("failed to convert exclude roles: %w", err)
		}

		if len(users) > 0 {
			conditions["users"] = users
		}
	}

	// Platforms
	if data.Platforms != nil {
		platforms := make(map[string]interface{})

		if err := convert.FrameworkToGraphStringSet(ctx, data.Platforms.IncludePlatforms, func(values []string) {
			if len(values) > 0 {
				platforms["includePlatforms"] = values
			}
		}); err != nil {
			return nil, fmt.Errorf("failed to convert include platforms: %w", err)
		}

		if err := convert.FrameworkToGraphStringSet(ctx, data.Platforms.ExcludePlatforms, func(values []string) {
			if len(values) > 0 {
				platforms["excludePlatforms"] = values
			}
		}); err != nil {
			return nil, fmt.Errorf("failed to convert exclude platforms: %w", err)
		}

		if len(platforms) > 0 {
			conditions["platforms"] = platforms
		}
	}

	// Locations
	if data.Locations != nil {
		locations := make(map[string]interface{})

		if err := convert.FrameworkToGraphStringSet(ctx, data.Locations.IncludeLocations, func(values []string) {
			if len(values) > 0 {
				locations["includeLocations"] = values
			}
		}); err != nil {
			return nil, fmt.Errorf("failed to convert include locations: %w", err)
		}

		if err := convert.FrameworkToGraphStringSet(ctx, data.Locations.ExcludeLocations, func(values []string) {
			if len(values) > 0 {
				locations["excludeLocations"] = values
			}
		}); err != nil {
			return nil, fmt.Errorf("failed to convert exclude locations: %w", err)
		}

		if len(locations) > 0 {
			conditions["locations"] = locations
		}
	}

	// Devices
	if data.Devices != nil {
		devices := make(map[string]interface{})

		if err := convert.FrameworkToGraphStringSet(ctx, data.Devices.IncludeDevices, func(values []string) {
			if len(values) > 0 {
				devices["includeDevices"] = values
			}
		}); err != nil {
			return nil, fmt.Errorf("failed to convert include devices: %w", err)
		}

		if err := convert.FrameworkToGraphStringSet(ctx, data.Devices.ExcludeDevices, func(values []string) {
			if len(values) > 0 {
				devices["excludeDevices"] = values
			}
		}); err != nil {
			return nil, fmt.Errorf("failed to convert exclude devices: %w", err)
		}

		if len(devices) > 0 {
			conditions["devices"] = devices
		}
	}

	// Risk levels
	if err := convert.FrameworkToGraphStringSet(ctx, data.SignInRiskLevels, func(values []string) {
		if len(values) > 0 {
			conditions["signInRiskLevels"] = values
		}
	}); err != nil {
		return nil, fmt.Errorf("failed to convert sign in risk levels: %w", err)
	}

	if err := convert.FrameworkToGraphStringSet(ctx, data.UserRiskLevels, func(values []string) {
		if len(values) > 0 {
			conditions["userRiskLevels"] = values
		}
	}); err != nil {
		return nil, fmt.Errorf("failed to convert user risk levels: %w", err)
	}

	return conditions, nil
}

// constructGrantControls builds the grant controls object
func constructGrantControls(ctx context.Context, data *ConditionalAccessGrantControls) (map[string]interface{}, error) {
	grantControls := make(map[string]interface{})

	convert.FrameworkToGraphString(data.Operator, func(value *string) {
		if value != nil {
			grantControls["operator"] = *value
		}
	})

	var builtInControls []string
	if err := convert.FrameworkToGraphStringSet(ctx, data.BuiltInControls, func(values []string) {
		builtInControls = values
	}); err != nil {
		return nil, fmt.Errorf("failed to convert built-in controls: %w", err)
	}
	if len(builtInControls) > 0 {
		grantControls["builtInControls"] = builtInControls
	}

	var customAuthFactors []string
	if err := convert.FrameworkToGraphStringSet(ctx, data.CustomAuthenticationFactors, func(values []string) {
		customAuthFactors = values
	}); err != nil {
		return nil, fmt.Errorf("failed to convert custom auth factors: %w", err)
	}
	// Always include customAuthenticationFactors if configured (even if empty)
	if !data.CustomAuthenticationFactors.IsNull() {
		grantControls["customAuthenticationFactors"] = customAuthFactors
	}

	var termsOfUse []string
	if err := convert.FrameworkToGraphStringSet(ctx, data.TermsOfUse, func(values []string) {
		termsOfUse = values
	}); err != nil {
		return nil, fmt.Errorf("failed to convert terms of use: %w", err)
	}
	// Always include termsOfUse if configured (even if empty)
	if !data.TermsOfUse.IsNull() {
		grantControls["termsOfUse"] = termsOfUse
	}

	if data.AuthenticationStrength != nil {
		authStrength := make(map[string]interface{})

		convert.FrameworkToGraphString(data.AuthenticationStrength.ID, func(value *string) {
			if value != nil {
				authStrength["id"] = *value
			}
		})

		convert.FrameworkToGraphString(data.AuthenticationStrength.DisplayName, func(value *string) {
			if value != nil {
				authStrength["displayName"] = *value
			}
		})

		convert.FrameworkToGraphString(data.AuthenticationStrength.Description, func(value *string) {
			if value != nil {
				authStrength["description"] = *value
			}
		})

		convert.FrameworkToGraphString(data.AuthenticationStrength.PolicyType, func(value *string) {
			if value != nil {
				authStrength["policyType"] = *value
			}
		})

		convert.FrameworkToGraphString(data.AuthenticationStrength.RequirementsSatisfied, func(value *string) {
			if value != nil {
				authStrength["requirementsSatisfied"] = *value
			}
		})

		var allowedCombinations []string
		if err := convert.FrameworkToGraphStringSet(ctx, data.AuthenticationStrength.AllowedCombinations, func(values []string) {
			allowedCombinations = values
		}); err != nil {
			return nil, fmt.Errorf("failed to convert allowed combinations: %w", err)
		}
		if len(allowedCombinations) > 0 {
			authStrength["allowedCombinations"] = allowedCombinations
		}

		if len(authStrength) > 0 {
			grantControls["authenticationStrength"] = authStrength
		}
	}

	return grantControls, nil
}

// constructSessionControls builds the session controls object
func constructSessionControls(ctx context.Context, data *ConditionalAccessSessionControls) (map[string]interface{}, error) {
	sessionControls := make(map[string]interface{})

	if data.ApplicationEnforcedRestrictions != nil {
		appEnforcedRestrictions := make(map[string]interface{})
		convert.FrameworkToGraphBool(data.ApplicationEnforcedRestrictions.IsEnabled, func(value *bool) {
			if value != nil {
				appEnforcedRestrictions["isEnabled"] = *value
			}
		})
		if len(appEnforcedRestrictions) > 0 {
			sessionControls["applicationEnforcedRestrictions"] = appEnforcedRestrictions
		}
	}

	if data.CloudAppSecurity != nil {
		cloudAppSecurity := make(map[string]interface{})
		convert.FrameworkToGraphBool(data.CloudAppSecurity.IsEnabled, func(value *bool) {
			if value != nil {
				cloudAppSecurity["isEnabled"] = *value
			}
		})
		convert.FrameworkToGraphString(data.CloudAppSecurity.CloudAppSecurityType, func(value *string) {
			if value != nil {
				cloudAppSecurity["cloudAppSecurityType"] = *value
			}
		})
		if len(cloudAppSecurity) > 0 {
			sessionControls["cloudAppSecurity"] = cloudAppSecurity
		}
	}

	if data.SignInFrequency != nil {
		signInFrequency := make(map[string]interface{})
		convert.FrameworkToGraphBool(data.SignInFrequency.IsEnabled, func(value *bool) {
			if value != nil {
				signInFrequency["isEnabled"] = *value
			}
		})
		convert.FrameworkToGraphString(data.SignInFrequency.Type, func(value *string) {
			if value != nil {
				signInFrequency["type"] = *value
			}
		})
		convert.FrameworkToGraphInt64(data.SignInFrequency.Value, func(value *int64) {
			if value != nil {
				signInFrequency["value"] = *value
			}
		})
		if len(signInFrequency) > 0 {
			sessionControls["signInFrequency"] = signInFrequency
		}
	}

	if data.PersistentBrowser != nil {
		persistentBrowser := make(map[string]interface{})
		convert.FrameworkToGraphBool(data.PersistentBrowser.IsEnabled, func(value *bool) {
			if value != nil {
				persistentBrowser["isEnabled"] = *value
			}
		})
		convert.FrameworkToGraphString(data.PersistentBrowser.Mode, func(value *string) {
			if value != nil {
				persistentBrowser["mode"] = *value
			}
		})
		if len(persistentBrowser) > 0 {
			sessionControls["persistentBrowser"] = persistentBrowser
		}
	}

	convert.FrameworkToGraphBool(data.DisableResilienceDefaults, func(value *bool) {
		if value != nil {
			sessionControls["disableResilienceDefaults"] = *value
		}
	})

	if data.ContinuousAccessEvaluation != nil {
		continuousAccessEval := make(map[string]interface{})
		convert.FrameworkToGraphString(data.ContinuousAccessEvaluation.Mode, func(value *string) {
			if value != nil {
				continuousAccessEval["mode"] = *value
			}
		})
		if len(continuousAccessEval) > 0 {
			sessionControls["continuousAccessEvaluation"] = continuousAccessEval
		}
	}

	if data.SecureSignInSession != nil {
		secureSignInSession := make(map[string]interface{})
		convert.FrameworkToGraphBool(data.SecureSignInSession.IsEnabled, func(value *bool) {
			if value != nil {
				secureSignInSession["isEnabled"] = *value
			}
		})
		if len(secureSignInSession) > 0 {
			sessionControls["secureSignInSession"] = secureSignInSession
		}
	}

	return sessionControls, nil
}
