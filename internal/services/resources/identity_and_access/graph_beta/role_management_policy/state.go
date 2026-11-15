package graphBetaRoleManagementPolicy

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// MapRemoteResourceStateToTerraform maps the remote role management policy to Terraform state.
//
// Only rules already tracked in data are mapped, so rules that are not configured never appear in state. When importAll
// is true (the first read after an import) every rule of a supported type is mapped instead.
func MapRemoteResourceStateToTerraform(
	ctx context.Context,
	data RoleManagementPolicyResourceModel,
	remoteResource graphmodels.UnifiedRoleManagementPolicyable,
	importAll bool,
) RoleManagementPolicyResourceModel {
	if remoteResource == nil {
		tflog.Debug(ctx, "Remote resource is nil")
		return data
	}

	tflog.Debug(ctx, "Starting to map remote state to Terraform state", map[string]any{
		"resourceId": remoteResource.GetId(),
	})

	data.ID = convert.GraphToFrameworkString(remoteResource.GetId())
	data.PolicyId = convert.GraphToFrameworkString(remoteResource.GetId())
	data.ScopeId = convert.GraphToFrameworkString(remoteResource.GetScopeId())
	data.ScopeType = convert.GraphToFrameworkString(remoteResource.GetScopeType())
	data.DisplayName = convert.GraphToFrameworkString(remoteResource.GetDisplayName())
	data.Description = convert.GraphToFrameworkString(remoteResource.GetDescription())
	data.IsOrganizationDefault = convert.GraphToFrameworkBool(remoteResource.GetIsOrganizationDefault())
	data.LastModifiedDateTime = convert.GraphToFrameworkTime(remoteResource.GetLastModifiedDateTime())

	remoteRules := rulesById(remoteResource)

	data.ExpirationRules = mapExpirationRules(data.ExpirationRules, remoteRules, importAll)
	data.EnablementRules = mapEnablementRules(ctx, data.EnablementRules, remoteRules, importAll)
	data.NotificationRules = mapNotificationRules(ctx, data.NotificationRules, remoteRules, importAll)
	data.AuthenticationContextRules = mapAuthenticationContextRules(data.AuthenticationContextRules, remoteRules, importAll)

	tflog.Debug(ctx, "Finished mapping remote state to Terraform state", map[string]any{
		"resourceId": data.ID.ValueString(),
	})

	return data
}

// tracked reports whether a rule should be mapped: every rule on import, otherwise only the rules already in state.
func tracked[T any](prior map[string]T, ruleId string, importAll bool) bool {
	if importAll {
		return true
	}
	_, ok := prior[ruleId]
	return ok
}

// finishRuleMap keeps an unconfigured rule type null after an import that found no rules of that type.
func finishRuleMap[T any](prior, result map[string]T, importAll bool) map[string]T {
	if prior == nil && !importAll {
		return nil
	}
	if prior == nil && len(result) == 0 {
		return nil
	}
	return result
}

func mapExpirationRules(
	prior map[string]ExpirationRuleModel,
	remoteRules map[string]graphmodels.UnifiedRoleManagementPolicyRuleable,
	importAll bool,
) map[string]ExpirationRuleModel {
	result := make(map[string]ExpirationRuleModel)
	for ruleId, remoteRule := range remoteRules {
		rule, ok := remoteRule.(graphmodels.UnifiedRoleManagementPolicyExpirationRuleable)
		if !ok || !tracked(prior, ruleId, importAll) {
			continue
		}
		model := ExpirationRuleModel{
			IsExpirationRequired: convert.GraphToFrameworkBool(rule.GetIsExpirationRequired()),
			MaximumDuration:      convert.GraphToFrameworkISODuration(rule.GetMaximumDuration()),
		}
		// Graph normalises durations, for example PT24H to P1D. Keep the configured spelling when it is equivalent.
		if priorRule, ok := prior[ruleId]; ok && !priorRule.MaximumDuration.IsNull() && !priorRule.MaximumDuration.IsUnknown() && !model.MaximumDuration.IsNull() {
			if durationsEquivalent(priorRule.MaximumDuration.ValueString(), model.MaximumDuration.ValueString()) {
				model.MaximumDuration = priorRule.MaximumDuration
			}
		}
		result[ruleId] = model
	}
	return finishRuleMap(prior, result, importAll)
}

func mapEnablementRules(
	ctx context.Context,
	prior map[string]EnablementRuleModel,
	remoteRules map[string]graphmodels.UnifiedRoleManagementPolicyRuleable,
	importAll bool,
) map[string]EnablementRuleModel {
	result := make(map[string]EnablementRuleModel)
	for ruleId, remoteRule := range remoteRules {
		rule, ok := remoteRule.(graphmodels.UnifiedRoleManagementPolicyEnablementRuleable)
		if !ok || !tracked(prior, ruleId, importAll) {
			continue
		}
		result[ruleId] = EnablementRuleModel{
			EnabledRules: stringSetPreserveEmpty(ctx, rule.GetEnabledRules()),
		}
	}
	return finishRuleMap(prior, result, importAll)
}

func mapNotificationRules(
	ctx context.Context,
	prior map[string]NotificationRuleModel,
	remoteRules map[string]graphmodels.UnifiedRoleManagementPolicyRuleable,
	importAll bool,
) map[string]NotificationRuleModel {
	result := make(map[string]NotificationRuleModel)
	for ruleId, remoteRule := range remoteRules {
		rule, ok := remoteRule.(graphmodels.UnifiedRoleManagementPolicyNotificationRuleable)
		if !ok || !tracked(prior, ruleId, importAll) {
			continue
		}
		result[ruleId] = NotificationRuleModel{
			NotificationLevel:          convert.GraphToFrameworkString(rule.GetNotificationLevel()),
			IsDefaultRecipientsEnabled: convert.GraphToFrameworkBool(rule.GetIsDefaultRecipientsEnabled()),
			NotificationRecipients:     stringSetPreserveEmpty(ctx, rule.GetNotificationRecipients()),
		}
	}
	return finishRuleMap(prior, result, importAll)
}

func mapAuthenticationContextRules(
	prior map[string]AuthenticationContextRuleModel,
	remoteRules map[string]graphmodels.UnifiedRoleManagementPolicyRuleable,
	importAll bool,
) map[string]AuthenticationContextRuleModel {
	result := make(map[string]AuthenticationContextRuleModel)
	for ruleId, remoteRule := range remoteRules {
		rule, ok := remoteRule.(graphmodels.UnifiedRoleManagementPolicyAuthenticationContextRuleable)
		if !ok || !tracked(prior, ruleId, importAll) {
			continue
		}
		result[ruleId] = AuthenticationContextRuleModel{
			IsEnabled:  convert.GraphToFrameworkBool(rule.GetIsEnabled()),
			ClaimValue: convert.GraphToFrameworkString(rule.GetClaimValue()),
		}
	}
	return finishRuleMap(prior, result, importAll)
}

// stringSetPreserveEmpty maps a Graph string collection to a set, mapping a missing collection to an empty set.
func stringSetPreserveEmpty(ctx context.Context, values []string) types.Set {
	if values == nil {
		values = []string{}
	}
	return convert.GraphToFrameworkStringSetPreserveEmpty(ctx, values)
}
