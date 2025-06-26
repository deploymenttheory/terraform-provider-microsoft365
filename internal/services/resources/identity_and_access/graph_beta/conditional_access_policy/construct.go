package graphBetaConditionalAccessPolicy

import (
	"context"
	"fmt"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/microsoft/kiota-abstractions-go/serialization"
)

/*
NewConditionalAccessRoot() - https://github.com/microsoftgraph/msgraph-beta-sdk-go/blob/main/models/conditional_access_root.go
conditional_access_rule - https://github.com/microsoftgraph/msgraph-beta-sdk-go/blob/main/models/conditional_access_rule.go
NewConditionalAccessUsers() - https://github.com/microsoftgraph/msgraph-beta-sdk-go/blob/main/models/conditional_access_users.go
NewConditionalAccessFilter() - https://github.com/microsoftgraph/msgraph-beta-sdk-go/blob/main/models/conditional_access_filter.go
conditional_access_status - https://github.com/microsoftgraph/msgraph-beta-sdk-go/blob/main/models/conditional_access_status.go
NewConditionalAccessDevices() - https://github.com/microsoftgraph/msgraph-beta-sdk-go/blob/main/models/conditional_access_devices.go
NewConditionalAccessTemplate() - https://github.com/microsoftgraph/msgraph-beta-sdk-go/blob/main/models/conditional_access_template.go
NewConditionalAccessLocations() - https://github.com/microsoftgraph/msgraph-beta-sdk-go/blob/main/models/conditional_access_locations.go
NewConditionalAccessPlatforms() - https://github.com/microsoftgraph/msgraph-beta-sdk-go/blob/main/models/conditional_access_platforms.go
conditional_access_client_app - https://github.com/microsoftgraph/msgraph-beta-sdk-go/blob/main/models/conditional_access_client_app.go
conditional_access_conditions - https://github.com/microsoftgraph/msgraph-beta-sdk-go/blob/main/models/conditional_access_conditions.go
NewConditionalAccessApplications() - https://github.com/microsoftgraph/msgraph-beta-sdk-go/blob/main/models/conditional_access_applications.go
NewConditionalAccessConditionSet() - https://github.com/microsoftgraph/msgraph-beta-sdk-go/blob/main/models/conditional_access_condition_set.go
NewConditionalAccessDeviceStates() - https://github.com/microsoftgraph/msgraph-beta-sdk-go/blob/main/models/conditional_access_device_states.go
conditional_access_grant_control - https://github.com/microsoftgraph/msgraph-beta-sdk-go/blob/main/models/conditional_access_grant_control.go
NewConditionalAccessPolicyDetail() - https://github.com/microsoftgraph/msgraph-beta-sdk-go/blob/main/models/conditional_access_policy_detail.go
NewConditionalAccessGrantControls() - https://github.com/microsoftgraph/msgraph-beta-sdk-go/blob/main/models/conditional_access_grant_controls.go
NewConditionalAccessNetworkAccess() - https://github.com/microsoftgraph/msgraph-beta-sdk-go/blob/main/models/conditional_access_network_access.go
NewConditionalAccessRuleSatisfied() - https://github.com/microsoftgraph/msgraph-beta-sdk-go/blob/main/models/conditional_access_rule_satisfied.go
ParseConditionalAccessDevicePlatform - https://github.com/microsoftgraph/msgraph-beta-sdk-go/blob/main/models/conditional_access_device_platform.go
NewConditionalAccessSessionControl() - https://github.com/microsoftgraph/msgraph-beta-sdk-go/blob/main/models/conditional_access_session_control.go
NewConditionalAccessExternalTenants() - https://github.com/microsoftgraph/msgraph-beta-sdk-go/blob/main/models/conditional_access_external_tenants.go
NewConditionalAccessSessionControls() - https://github.com/microsoftgraph/msgraph-beta-sdk-go/blob/main/models/conditional_access_session_controls.go
ParseConditionalAccessTransferMethods - https://github.com/microsoftgraph/msgraph-beta-sdk-go/blob/main/models/conditional_access_transfer_methods.go
NewConditionalAccessClientApplications() - https://github.com/microsoftgraph/msgraph-beta-sdk-go/blob/main/models/conditional_access_client_applications.go
ParseConditionalAccessInsiderRiskLevels - https://github.com/microsoftgraph/msgraph-beta-sdk-go/blob/main/models/conditional_access_insider_risk_levels.go
NewConditionalAccessAllExternalTenants() - https://github.com/microsoftgraph/msgraph-beta-sdk-go/blob/main/models/conditional_access_all_external_tenants.go
NewConditionalAccessAuthenticationFlows() - https://github.com/microsoftgraph/msgraph-beta-sdk-go/blob/main/models/conditional_access_authentication_flows.go
*/

// JSONRequestBody is a simple wrapper to make JSON data implement serialization.Parsable
type JSONRequestBody struct {
	data map[string]interface{}
}

// GetFieldDeserializers implements serialization.Parsable
func (j *JSONRequestBody) GetFieldDeserializers() map[string]func(serialization.ParseNode) error {
	return make(map[string]func(serialization.ParseNode) error)
}

// Serialize implements serialization.Parsable
func (j *JSONRequestBody) Serialize(writer serialization.SerializationWriter) error {
	if j.data == nil {
		return nil
	}

	// Serialize each field in the data map
	for key, value := range j.data {
		switch v := value.(type) {
		case string:
			err := writer.WriteStringValue(key, &v)
			if err != nil {
				return err
			}
		case map[string]interface{}:
			// For nested objects, we need to serialize them recursively
			nested := &JSONRequestBody{data: v}
			err := writer.WriteObjectValue(key, nested)
			if err != nil {
				return err
			}
		case []string:
			err := writer.WriteCollectionOfStringValues(key, v)
			if err != nil {
				return err
			}
		case bool:
			err := writer.WriteBoolValue(key, &v)
			if err != nil {
				return err
			}
		case int64:
			err := writer.WriteInt64Value(key, &v)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// constructResource converts the Terraform resource model to the Microsoft Graph API model
// Since there's no direct ConditionalAccessPolicy constructor, we need to build it from components
func constructResource(ctx context.Context, data *ConditionalAccessPolicyResourceModel) (serialization.Parsable, error) {

	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	requestBody := make(map[string]interface{})

	// Basic properties
	if !data.DisplayName.IsNull() && !data.DisplayName.IsUnknown() {
		requestBody["displayName"] = data.DisplayName.ValueString()
	}

	if !data.State.IsNull() && !data.State.IsUnknown() {
		requestBody["state"] = data.State.ValueString()
	}

	if !data.TemplateId.IsNull() && !data.TemplateId.IsUnknown() {
		requestBody["templateId"] = data.TemplateId.ValueString()
	}

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

	resource := &JSONRequestBody{data: requestBody}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), resource); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]interface{}{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return resource, nil
}

// constructConditions builds the conditions object using available Graph models
func constructConditions(ctx context.Context, data *ConditionalAccessConditions) (map[string]interface{}, error) {
	conditions := make(map[string]interface{})

	// Client app types
	var clientAppTypes []string
	if err := convert.FrameworkToGraphStringSet(ctx, data.ClientAppTypes, func(values []string) {
		clientAppTypes = values
	}); err != nil {
		return nil, fmt.Errorf("failed to convert client app types: %w", err)
	}
	if len(clientAppTypes) > 0 {
		conditions["clientAppTypes"] = clientAppTypes
	}

	// Applications
	if data.Applications != nil {
		applications := make(map[string]interface{})

		var includeApplications []string
		if err := convert.FrameworkToGraphStringSet(ctx, data.Applications.IncludeApplications, func(values []string) {
			includeApplications = values
		}); err != nil {
			return nil, fmt.Errorf("failed to convert include applications: %w", err)
		}
		if len(includeApplications) > 0 {
			applications["includeApplications"] = includeApplications
		}

		var excludeApplications []string
		if err := convert.FrameworkToGraphStringSet(ctx, data.Applications.ExcludeApplications, func(values []string) {
			excludeApplications = values
		}); err != nil {
			return nil, fmt.Errorf("failed to convert exclude applications: %w", err)
		}
		if len(excludeApplications) > 0 {
			applications["excludeApplications"] = excludeApplications
		}

		var includeUserActions []string
		if err := convert.FrameworkToGraphStringSet(ctx, data.Applications.IncludeUserActions, func(values []string) {
			includeUserActions = values
		}); err != nil {
			return nil, fmt.Errorf("failed to convert include user actions: %w", err)
		}
		if len(includeUserActions) > 0 {
			applications["includeUserActions"] = includeUserActions
		}

		var includeAuthContextClassRefs []string
		if err := convert.FrameworkToGraphStringSet(ctx, data.Applications.IncludeAuthenticationContextClassReferences, func(values []string) {
			includeAuthContextClassRefs = values
		}); err != nil {
			return nil, fmt.Errorf("failed to convert include auth context class refs: %w", err)
		}
		if len(includeAuthContextClassRefs) > 0 {
			applications["includeAuthenticationContextClassReferences"] = includeAuthContextClassRefs
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

		var includeUsers []string
		if err := convert.FrameworkToGraphStringSet(ctx, data.Users.IncludeUsers, func(values []string) {
			includeUsers = values
		}); err != nil {
			return nil, fmt.Errorf("failed to convert include users: %w", err)
		}
		if len(includeUsers) > 0 {
			users["includeUsers"] = includeUsers
		}

		var excludeUsers []string
		if err := convert.FrameworkToGraphStringSet(ctx, data.Users.ExcludeUsers, func(values []string) {
			excludeUsers = values
		}); err != nil {
			return nil, fmt.Errorf("failed to convert exclude users: %w", err)
		}
		if len(excludeUsers) > 0 {
			users["excludeUsers"] = excludeUsers
		}

		var includeGroups []string
		if err := convert.FrameworkToGraphStringSet(ctx, data.Users.IncludeGroups, func(values []string) {
			includeGroups = values
		}); err != nil {
			return nil, fmt.Errorf("failed to convert include groups: %w", err)
		}
		if len(includeGroups) > 0 {
			users["includeGroups"] = includeGroups
		}

		var excludeGroups []string
		if err := convert.FrameworkToGraphStringSet(ctx, data.Users.ExcludeGroups, func(values []string) {
			excludeGroups = values
		}); err != nil {
			return nil, fmt.Errorf("failed to convert exclude groups: %w", err)
		}
		if len(excludeGroups) > 0 {
			users["excludeGroups"] = excludeGroups
		}

		var includeRoles []string
		if err := convert.FrameworkToGraphStringSet(ctx, data.Users.IncludeRoles, func(values []string) {
			includeRoles = values
		}); err != nil {
			return nil, fmt.Errorf("failed to convert include roles: %w", err)
		}
		if len(includeRoles) > 0 {
			users["includeRoles"] = includeRoles
		}

		var excludeRoles []string
		if err := convert.FrameworkToGraphStringSet(ctx, data.Users.ExcludeRoles, func(values []string) {
			excludeRoles = values
		}); err != nil {
			return nil, fmt.Errorf("failed to convert exclude roles: %w", err)
		}
		if len(excludeRoles) > 0 {
			users["excludeRoles"] = excludeRoles
		}

		if len(users) > 0 {
			conditions["users"] = users
		}
	}

	// Platforms
	if data.Platforms != nil {
		platforms := make(map[string]interface{})

		var includePlatforms []string
		if err := convert.FrameworkToGraphStringSet(ctx, data.Platforms.IncludePlatforms, func(values []string) {
			includePlatforms = values
		}); err != nil {
			return nil, fmt.Errorf("failed to convert include platforms: %w", err)
		}
		if len(includePlatforms) > 0 {
			platforms["includePlatforms"] = includePlatforms
		}

		var excludePlatforms []string
		if err := convert.FrameworkToGraphStringSet(ctx, data.Platforms.ExcludePlatforms, func(values []string) {
			excludePlatforms = values
		}); err != nil {
			return nil, fmt.Errorf("failed to convert exclude platforms: %w", err)
		}
		if len(excludePlatforms) > 0 {
			platforms["excludePlatforms"] = excludePlatforms
		}

		if len(platforms) > 0 {
			conditions["platforms"] = platforms
		}
	}

	// Locations
	if data.Locations != nil {
		locations := make(map[string]interface{})

		var includeLocations []string
		if err := convert.FrameworkToGraphStringSet(ctx, data.Locations.IncludeLocations, func(values []string) {
			includeLocations = values
		}); err != nil {
			return nil, fmt.Errorf("failed to convert include locations: %w", err)
		}
		if len(includeLocations) > 0 {
			locations["includeLocations"] = includeLocations
		}

		var excludeLocations []string
		if err := convert.FrameworkToGraphStringSet(ctx, data.Locations.ExcludeLocations, func(values []string) {
			excludeLocations = values
		}); err != nil {
			return nil, fmt.Errorf("failed to convert exclude locations: %w", err)
		}
		if len(excludeLocations) > 0 {
			locations["excludeLocations"] = excludeLocations
		}

		if len(locations) > 0 {
			conditions["locations"] = locations
		}
	}

	// Devices
	if data.Devices != nil {
		devices := make(map[string]interface{})

		var includeDevices []string
		if err := convert.FrameworkToGraphStringSet(ctx, data.Devices.IncludeDevices, func(values []string) {
			includeDevices = values
		}); err != nil {
			return nil, fmt.Errorf("failed to convert include devices: %w", err)
		}
		if len(includeDevices) > 0 {
			devices["includeDevices"] = includeDevices
		}

		var excludeDevices []string
		if err := convert.FrameworkToGraphStringSet(ctx, data.Devices.ExcludeDevices, func(values []string) {
			excludeDevices = values
		}); err != nil {
			return nil, fmt.Errorf("failed to convert exclude devices: %w", err)
		}
		if len(excludeDevices) > 0 {
			devices["excludeDevices"] = excludeDevices
		}

		var includeDeviceStates []string
		if err := convert.FrameworkToGraphStringSet(ctx, data.Devices.IncludeDeviceStates, func(values []string) {
			includeDeviceStates = values
		}); err != nil {
			return nil, fmt.Errorf("failed to convert include device states: %w", err)
		}
		if len(includeDeviceStates) > 0 {
			devices["includeDeviceStates"] = includeDeviceStates
		}

		var excludeDeviceStates []string
		if err := convert.FrameworkToGraphStringSet(ctx, data.Devices.ExcludeDeviceStates, func(values []string) {
			excludeDeviceStates = values
		}); err != nil {
			return nil, fmt.Errorf("failed to convert exclude device states: %w", err)
		}
		if len(excludeDeviceStates) > 0 {
			devices["excludeDeviceStates"] = excludeDeviceStates
		}

		if data.Devices.DeviceFilter != nil {
			deviceFilter := make(map[string]interface{})
			convert.FrameworkToGraphString(data.Devices.DeviceFilter.Mode, func(value *string) {
				if value != nil {
					deviceFilter["mode"] = *value
				}
			})
			convert.FrameworkToGraphString(data.Devices.DeviceFilter.Rule, func(value *string) {
				if value != nil {
					deviceFilter["rule"] = *value
				}
			})
			if len(deviceFilter) > 0 {
				devices["deviceFilter"] = deviceFilter
			}
		}

		if len(devices) > 0 {
			conditions["devices"] = devices
		}
	}

	// Sign-in risk levels
	var signInRiskLevels []string
	if err := convert.FrameworkToGraphStringSet(ctx, data.SignInRiskLevels, func(values []string) {
		signInRiskLevels = values
	}); err != nil {
		return nil, fmt.Errorf("failed to convert sign-in risk levels: %w", err)
	}
	if len(signInRiskLevels) > 0 {
		conditions["signInRiskLevels"] = signInRiskLevels
	}

	// User risk levels
	var userRiskLevels []string
	if err := convert.FrameworkToGraphStringSet(ctx, data.UserRiskLevels, func(values []string) {
		userRiskLevels = values
	}); err != nil {
		return nil, fmt.Errorf("failed to convert user risk levels: %w", err)
	}
	if len(userRiskLevels) > 0 {
		conditions["userRiskLevels"] = userRiskLevels
	}

	// Service principal risk levels
	var servicePrincipalRiskLevels []string
	if err := convert.FrameworkToGraphStringSet(ctx, data.ServicePrincipalRiskLevels, func(values []string) {
		servicePrincipalRiskLevels = values
	}); err != nil {
		return nil, fmt.Errorf("failed to convert service principal risk levels: %w", err)
	}
	if len(servicePrincipalRiskLevels) > 0 {
		conditions["servicePrincipalRiskLevels"] = servicePrincipalRiskLevels
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
	if len(customAuthFactors) > 0 {
		grantControls["customAuthenticationFactors"] = customAuthFactors
	}

	var termsOfUse []string
	if err := convert.FrameworkToGraphStringSet(ctx, data.TermsOfUse, func(values []string) {
		termsOfUse = values
	}); err != nil {
		return nil, fmt.Errorf("failed to convert terms of use: %w", err)
	}
	if len(termsOfUse) > 0 {
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
