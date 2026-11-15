package graphBetaIdentityAndAccessAuthenticationMethodsPolicy

import (
	"context"
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// Fixed IDs of the authentication method configurations this resource manages.
const (
	methodIDEmail               = "Email"
	methodIDTemporaryAccessPass = "TemporaryAccessPass"
	methodIDSms                 = "Sms"
	methodIDSoftwareOath        = "SoftwareOath"
	methodIDHardwareOath        = "HardwareOath"
)

// methodRequest is a PATCH body for one authentication method configuration.
type methodRequest struct {
	ID   string
	Body graphmodels.AuthenticationMethodConfigurationable
}

// requestBodies holds every PATCH request needed to apply the configuration.
//
// Policy is nil when no policy-level setting changes. Methods holds one entry per changed method
// configuration, each PATCHed at /policies/authenticationMethodsPolicy/authenticationMethodConfigurations/{id}.
type requestBodies struct {
	Policy  graphmodels.AuthenticationMethodsPolicyable
	Methods []methodRequest
}

// constructResource builds the PATCH requests for the authentication methods policy.
//
// A value is sent when it is known in the plan and, on update, differs from the prior state. prior is nil on
// create. Values that are not configured are unknown on create (and equal to the prior state on update), so
// unmanaged settings and method configurations are never written.
func constructResource(ctx context.Context, data *AuthenticationMethodsPolicyResourceModel, prior *AuthenticationMethodsPolicyResourceModel) (*requestBodies, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	if prior == nil {
		prior = &AuthenticationMethodsPolicyResourceModel{}
	}
	isCreate := prior.ID.IsNull() || prior.ID.IsUnknown()
	changed := func(planned, current attr.Value) bool {
		if planned.IsNull() || planned.IsUnknown() {
			return false
		}
		return isCreate || !planned.Equal(current)
	}

	bodies := &requestBodies{}

	policy, hasPolicyChanges, err := constructPolicy(ctx, data, prior, changed)
	if err != nil {
		return nil, err
	}
	if hasPolicyChanges {
		bodies.Policy = policy
	}

	if changed(data.Email, prior.Email) {
		body, err := constructEmail(ctx, data.Email)
		if err != nil {
			return nil, fmt.Errorf("email: %w", err)
		}
		bodies.Methods = append(bodies.Methods, methodRequest{ID: methodIDEmail, Body: body})
	}
	if changed(data.TemporaryAccessPass, prior.TemporaryAccessPass) {
		body, err := constructTemporaryAccessPass(ctx, data.TemporaryAccessPass)
		if err != nil {
			return nil, fmt.Errorf("temporary_access_pass: %w", err)
		}
		bodies.Methods = append(bodies.Methods, methodRequest{ID: methodIDTemporaryAccessPass, Body: body})
	}
	if changed(data.Sms, prior.Sms) {
		body, err := constructSms(ctx, data.Sms)
		if err != nil {
			return nil, fmt.Errorf("sms: %w", err)
		}
		bodies.Methods = append(bodies.Methods, methodRequest{ID: methodIDSms, Body: body})
	}
	if changed(data.SoftwareOath, prior.SoftwareOath) {
		configuration := graphmodels.NewSoftwareOathAuthenticationMethodConfiguration()
		if err := constructOath(ctx, data.SoftwareOath, configuration, configuration.SetIncludeTargets); err != nil {
			return nil, fmt.Errorf("software_oath: %w", err)
		}
		bodies.Methods = append(bodies.Methods, methodRequest{ID: methodIDSoftwareOath, Body: configuration})
	}
	if changed(data.HardwareOath, prior.HardwareOath) {
		configuration := graphmodels.NewHardwareOathAuthenticationMethodConfiguration()
		if err := constructOath(ctx, data.HardwareOath, configuration, configuration.SetIncludeTargets); err != nil {
			return nil, fmt.Errorf("hardware_oath: %w", err)
		}
		bodies.Methods = append(bodies.Methods, methodRequest{ID: methodIDHardwareOath, Body: configuration})
	}

	if bodies.Policy != nil {
		if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), bodies.Policy); err != nil {
			tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
		}
	}
	for _, method := range bodies.Methods {
		message := fmt.Sprintf("Final JSON to be sent to Graph API for %s method configuration %s", ResourceName, method.ID)
		if err := constructors.DebugLogGraphObject(ctx, message, method.Body); err != nil {
			tflog.Error(ctx, "Failed to debug log object", map[string]any{"error": err.Error()})
		}
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return bodies, nil
}

// constructPolicy builds the policy-level PATCH body and reports whether any policy-level setting changed.
func constructPolicy(
	ctx context.Context,
	data *AuthenticationMethodsPolicyResourceModel,
	prior *AuthenticationMethodsPolicyResourceModel,
	changed func(planned, current attr.Value) bool,
) (graphmodels.AuthenticationMethodsPolicyable, bool, error) {
	policy := graphmodels.NewAuthenticationMethodsPolicy()
	hasChanges := false

	if changed(data.PolicyMigrationState, prior.PolicyMigrationState) {
		if err := convert.FrameworkToGraphEnum(data.PolicyMigrationState, graphmodels.ParseAuthenticationMethodsPolicyMigrationState, policy.SetPolicyMigrationState); err != nil {
			return nil, false, fmt.Errorf("policy_migration_state: %w", err)
		}
		hasChanges = true
	}

	if changed(data.RegistrationEnforcement, prior.RegistrationEnforcement) {
		enforcement, err := constructRegistrationEnforcement(ctx, data.RegistrationEnforcement)
		if err != nil {
			return nil, false, fmt.Errorf("registration_enforcement: %w", err)
		}
		policy.SetRegistrationEnforcement(enforcement)
		hasChanges = true
	}

	if changed(data.ReportSuspiciousActivitySettings, prior.ReportSuspiciousActivitySettings) {
		settings, err := constructReportSuspiciousActivitySettings(ctx, data.ReportSuspiciousActivitySettings)
		if err != nil {
			return nil, false, fmt.Errorf("report_suspicious_activity_settings: %w", err)
		}
		policy.SetReportSuspiciousActivitySettings(settings)
		hasChanges = true
	}

	if changed(data.SystemCredentialPreferences, prior.SystemCredentialPreferences) {
		preferences, err := constructSystemCredentialPreferences(ctx, data.SystemCredentialPreferences)
		if err != nil {
			return nil, false, fmt.Errorf("system_credential_preferences: %w", err)
		}
		policy.SetSystemCredentialPreferences(preferences)
		hasChanges = true
	}

	return policy, hasChanges, nil
}

// constructRegistrationEnforcement builds the registrationEnforcement complex type.
func constructRegistrationEnforcement(ctx context.Context, value types.Object) (graphmodels.RegistrationEnforcementable, error) {
	var model RegistrationEnforcementModel
	if diags := value.As(ctx, &model, basetypes.ObjectAsOptions{}); diags.HasError() {
		return nil, readError("registration_enforcement", diags)
	}

	enforcement := graphmodels.NewRegistrationEnforcement()
	if model.AuthenticationMethodsRegistrationCampaign.IsNull() || model.AuthenticationMethodsRegistrationCampaign.IsUnknown() {
		return enforcement, nil
	}

	var campaignModel RegistrationCampaignModel
	if diags := model.AuthenticationMethodsRegistrationCampaign.As(ctx, &campaignModel, basetypes.ObjectAsOptions{}); diags.HasError() {
		return nil, readError("authentication_methods_registration_campaign", diags)
	}

	campaign := graphmodels.NewAuthenticationMethodsRegistrationCampaign()
	convert.FrameworkToGraphInt32(campaignModel.SnoozeDurationInDays, campaign.SetSnoozeDurationInDays)
	convert.FrameworkToGraphBool(campaignModel.EnforceRegistrationAfterAllowedSnoozes, campaign.SetEnforceRegistrationAfterAllowedSnoozes)
	if err := convert.FrameworkToGraphEnum(campaignModel.State, graphmodels.ParseAdvancedConfigState, campaign.SetState); err != nil {
		return nil, fmt.Errorf("state: %w", err)
	}

	if isKnown(campaignModel.IncludeTargets) {
		var targets []RegistrationCampaignIncludeTargetModel
		if diags := campaignModel.IncludeTargets.ElementsAs(ctx, &targets, false); diags.HasError() {
			return nil, readError("include_targets", diags)
		}
		includeTargets := make([]graphmodels.AuthenticationMethodsRegistrationCampaignIncludeTargetable, 0, len(targets))
		for _, target := range targets {
			includeTarget := graphmodels.NewAuthenticationMethodsRegistrationCampaignIncludeTarget()
			convert.FrameworkToGraphString(target.ID, includeTarget.SetId)
			if err := convert.FrameworkToGraphEnum(target.TargetType, graphmodels.ParseAuthenticationMethodTargetType, includeTarget.SetTargetType); err != nil {
				return nil, fmt.Errorf("include_targets.target_type: %w", err)
			}
			convert.FrameworkToGraphString(target.TargetedAuthenticationMethod, includeTarget.SetTargetedAuthenticationMethod)
			includeTargets = append(includeTargets, includeTarget)
		}
		campaign.SetIncludeTargets(includeTargets)
	}

	excludeTargets, err := constructExcludeTargets(ctx, campaignModel.ExcludeTargets)
	if err != nil {
		return nil, err
	}
	if excludeTargets != nil {
		campaign.SetExcludeTargets(excludeTargets)
	}

	enforcement.SetAuthenticationMethodsRegistrationCampaign(campaign)
	return enforcement, nil
}

// constructReportSuspiciousActivitySettings builds the reportSuspiciousActivitySettings complex type.
func constructReportSuspiciousActivitySettings(ctx context.Context, value types.Object) (graphmodels.ReportSuspiciousActivitySettingsable, error) {
	var model ReportSuspiciousActivitySettingsModel
	if diags := value.As(ctx, &model, basetypes.ObjectAsOptions{}); diags.HasError() {
		return nil, readError("report_suspicious_activity_settings", diags)
	}

	settings := graphmodels.NewReportSuspiciousActivitySettings()
	if err := convert.FrameworkToGraphEnum(model.State, graphmodels.ParseAdvancedConfigState, settings.SetState); err != nil {
		return nil, fmt.Errorf("state: %w", err)
	}
	convert.FrameworkToGraphInt32(model.VoiceReportingCode, settings.SetVoiceReportingCode)

	if !model.IncludeTarget.IsNull() && !model.IncludeTarget.IsUnknown() {
		var targetModel TargetModel
		if diags := model.IncludeTarget.As(ctx, &targetModel, basetypes.ObjectAsOptions{}); diags.HasError() {
			return nil, readError("include_target", diags)
		}
		includeTarget := graphmodels.NewIncludeTarget()
		convert.FrameworkToGraphString(targetModel.ID, includeTarget.SetId)
		if err := convert.FrameworkToGraphEnum(targetModel.TargetType, graphmodels.ParseAuthenticationMethodTargetType, includeTarget.SetTargetType); err != nil {
			return nil, fmt.Errorf("include_target.target_type: %w", err)
		}
		settings.SetIncludeTarget(includeTarget)
	}

	return settings, nil
}

// constructSystemCredentialPreferences builds the systemCredentialPreferences complex type.
func constructSystemCredentialPreferences(ctx context.Context, value types.Object) (graphmodels.SystemCredentialPreferencesable, error) {
	var model SystemCredentialPreferencesModel
	if diags := value.As(ctx, &model, basetypes.ObjectAsOptions{}); diags.HasError() {
		return nil, readError("system_credential_preferences", diags)
	}

	preferences := graphmodels.NewSystemCredentialPreferences()
	if err := convert.FrameworkToGraphEnum(model.State, graphmodels.ParseAdvancedConfigState, preferences.SetState); err != nil {
		return nil, fmt.Errorf("state: %w", err)
	}

	if isKnown(model.IncludeTargets) {
		var targets []TargetModel
		if diags := model.IncludeTargets.ElementsAs(ctx, &targets, false); diags.HasError() {
			return nil, readError("include_targets", diags)
		}
		includeTargets := make([]graphmodels.IncludeTargetable, 0, len(targets))
		for _, target := range targets {
			includeTarget := graphmodels.NewIncludeTarget()
			convert.FrameworkToGraphString(target.ID, includeTarget.SetId)
			if err := convert.FrameworkToGraphEnum(target.TargetType, graphmodels.ParseAuthenticationMethodTargetType, includeTarget.SetTargetType); err != nil {
				return nil, fmt.Errorf("include_targets.target_type: %w", err)
			}
			includeTargets = append(includeTargets, includeTarget)
		}
		preferences.SetIncludeTargets(includeTargets)
	}

	excludeTargets, err := constructExcludeTargets(ctx, model.ExcludeTargets)
	if err != nil {
		return nil, err
	}
	if excludeTargets != nil {
		preferences.SetExcludeTargets(excludeTargets)
	}

	return preferences, nil
}

// constructEmail builds the emailAuthenticationMethodConfiguration PATCH body.
func constructEmail(ctx context.Context, value types.Object) (graphmodels.AuthenticationMethodConfigurationable, error) {
	var model EmailMethodModel
	if diags := value.As(ctx, &model, basetypes.ObjectAsOptions{}); diags.HasError() {
		return nil, readError("email", diags)
	}

	configuration := graphmodels.NewEmailAuthenticationMethodConfiguration()
	if err := setMethodCommon(ctx, configuration, model.State, model.ExcludeTargets); err != nil {
		return nil, err
	}
	if err := convert.FrameworkToGraphEnum(model.AllowExternalIdToUseEmailOtp, graphmodels.ParseExternalEmailOtpState, configuration.SetAllowExternalIdToUseEmailOtp); err != nil {
		return nil, fmt.Errorf("allow_external_id_to_use_email_otp: %w", err)
	}
	includeTargets, err := constructMethodTargets(ctx, model.IncludeTargets)
	if err != nil {
		return nil, err
	}
	if includeTargets != nil {
		configuration.SetIncludeTargets(includeTargets)
	}

	return configuration, nil
}

// constructTemporaryAccessPass builds the temporaryAccessPassAuthenticationMethodConfiguration PATCH body.
func constructTemporaryAccessPass(ctx context.Context, value types.Object) (graphmodels.AuthenticationMethodConfigurationable, error) {
	var model TemporaryAccessPassMethodModel
	if diags := value.As(ctx, &model, basetypes.ObjectAsOptions{}); diags.HasError() {
		return nil, readError("temporary_access_pass", diags)
	}

	configuration := graphmodels.NewTemporaryAccessPassAuthenticationMethodConfiguration()
	if err := setMethodCommon(ctx, configuration, model.State, model.ExcludeTargets); err != nil {
		return nil, err
	}
	convert.FrameworkToGraphInt32(model.DefaultLifetimeInMinutes, configuration.SetDefaultLifetimeInMinutes)
	convert.FrameworkToGraphInt32(model.DefaultLength, configuration.SetDefaultLength)
	convert.FrameworkToGraphInt32(model.MinimumLifetimeInMinutes, configuration.SetMinimumLifetimeInMinutes)
	convert.FrameworkToGraphInt32(model.MaximumLifetimeInMinutes, configuration.SetMaximumLifetimeInMinutes)
	convert.FrameworkToGraphBool(model.IsUsableOnce, configuration.SetIsUsableOnce)
	includeTargets, err := constructMethodTargets(ctx, model.IncludeTargets)
	if err != nil {
		return nil, err
	}
	if includeTargets != nil {
		configuration.SetIncludeTargets(includeTargets)
	}

	return configuration, nil
}

// constructSms builds the smsAuthenticationMethodConfiguration PATCH body.
func constructSms(ctx context.Context, value types.Object) (graphmodels.AuthenticationMethodConfigurationable, error) {
	var model SmsMethodModel
	if diags := value.As(ctx, &model, basetypes.ObjectAsOptions{}); diags.HasError() {
		return nil, readError("sms", diags)
	}

	configuration := graphmodels.NewSmsAuthenticationMethodConfiguration()
	if err := setMethodCommon(ctx, configuration, model.State, model.ExcludeTargets); err != nil {
		return nil, err
	}

	if isKnown(model.IncludeTargets) {
		var targets []SmsMethodTargetModel
		if diags := model.IncludeTargets.ElementsAs(ctx, &targets, false); diags.HasError() {
			return nil, readError("include_targets", diags)
		}
		includeTargets := make([]graphmodels.SmsAuthenticationMethodTargetable, 0, len(targets))
		for _, target := range targets {
			includeTarget := graphmodels.NewSmsAuthenticationMethodTarget()
			convert.FrameworkToGraphString(target.ID, includeTarget.SetId)
			if err := convert.FrameworkToGraphEnum(target.TargetType, graphmodels.ParseAuthenticationMethodTargetType, includeTarget.SetTargetType); err != nil {
				return nil, fmt.Errorf("include_targets.target_type: %w", err)
			}
			convert.FrameworkToGraphBool(target.IsRegistrationRequired, includeTarget.SetIsRegistrationRequired)
			convert.FrameworkToGraphBool(target.IsUsableForSignIn, includeTarget.SetIsUsableForSignIn)
			includeTargets = append(includeTargets, includeTarget)
		}
		configuration.SetIncludeTargets(includeTargets)
	}

	return configuration, nil
}

// constructOath fills a softwareOath or hardwareOath method configuration PATCH body.
func constructOath(
	ctx context.Context,
	value types.Object,
	configuration graphmodels.AuthenticationMethodConfigurationable,
	setIncludeTargets func([]graphmodels.AuthenticationMethodTargetable),
) error {
	var model OathMethodModel
	if diags := value.As(ctx, &model, basetypes.ObjectAsOptions{}); diags.HasError() {
		return readError("method configuration", diags)
	}

	if err := setMethodCommon(ctx, configuration, model.State, model.ExcludeTargets); err != nil {
		return err
	}
	includeTargets, err := constructMethodTargets(ctx, model.IncludeTargets)
	if err != nil {
		return err
	}
	if includeTargets != nil {
		setIncludeTargets(includeTargets)
	}

	return nil
}

// setMethodCommon sets the state and excludeTargets shared by every method configuration.
func setMethodCommon(ctx context.Context, configuration graphmodels.AuthenticationMethodConfigurationable, state types.String, excludeTargets types.Set) error {
	if err := convert.FrameworkToGraphEnum(state, graphmodels.ParseAuthenticationMethodState, configuration.SetState); err != nil {
		return fmt.Errorf("state: %w", err)
	}
	targets, err := constructExcludeTargets(ctx, excludeTargets)
	if err != nil {
		return err
	}
	if targets != nil {
		configuration.SetExcludeTargets(targets)
	}
	return nil
}

// constructMethodTargets builds authenticationMethodTarget values. It returns nil when the set is not known,
// so an unset collection is omitted from the request instead of being sent as null.
func constructMethodTargets(ctx context.Context, value types.Set) ([]graphmodels.AuthenticationMethodTargetable, error) {
	if !isKnown(value) {
		return nil, nil
	}
	var targets []MethodTargetModel
	if diags := value.ElementsAs(ctx, &targets, false); diags.HasError() {
		return nil, readError("include_targets", diags)
	}
	result := make([]graphmodels.AuthenticationMethodTargetable, 0, len(targets))
	for _, target := range targets {
		includeTarget := graphmodels.NewAuthenticationMethodTarget()
		convert.FrameworkToGraphString(target.ID, includeTarget.SetId)
		if err := convert.FrameworkToGraphEnum(target.TargetType, graphmodels.ParseAuthenticationMethodTargetType, includeTarget.SetTargetType); err != nil {
			return nil, fmt.Errorf("include_targets.target_type: %w", err)
		}
		convert.FrameworkToGraphBool(target.IsRegistrationRequired, includeTarget.SetIsRegistrationRequired)
		result = append(result, includeTarget)
	}
	return result, nil
}

// constructExcludeTargets builds excludeTarget values. It returns nil when the set is not known, so an unset
// collection is omitted from the request instead of being sent as null.
func constructExcludeTargets(ctx context.Context, value types.Set) ([]graphmodels.ExcludeTargetable, error) {
	if !isKnown(value) {
		return nil, nil
	}
	var targets []TargetModel
	if diags := value.ElementsAs(ctx, &targets, false); diags.HasError() {
		return nil, readError("exclude_targets", diags)
	}
	result := make([]graphmodels.ExcludeTargetable, 0, len(targets))
	for _, target := range targets {
		excludeTarget := graphmodels.NewExcludeTarget()
		convert.FrameworkToGraphString(target.ID, excludeTarget.SetId)
		if err := convert.FrameworkToGraphEnum(target.TargetType, graphmodels.ParseAuthenticationMethodTargetType, excludeTarget.SetTargetType); err != nil {
			return nil, fmt.Errorf("exclude_targets.target_type: %w", err)
		}
		result = append(result, excludeTarget)
	}
	return result, nil
}

// errReadValue is returned when a Terraform value cannot be read into its model.
var errReadValue = errors.New("could not read Terraform value")

// readError wraps the diagnostics from reading the named Terraform value.
func readError(name string, diags diag.Diagnostics) error {
	return fmt.Errorf("%w %s: %v", errReadValue, name, diags)
}

// isKnown reports whether a value is neither null nor unknown.
func isKnown(value attr.Value) bool {
	return !value.IsNull() && !value.IsUnknown()
}
