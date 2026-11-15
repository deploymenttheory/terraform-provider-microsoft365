package graphBetaRoleManagementPolicy

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/constructors"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/convert"
)

// constructResource resolves the policy to adopt, validates the configured rules against it and builds the PATCH body.
//
// API Calls:
//   - GET /policies/roleManagementPolicyAssignments?$filter=scopeId eq '{scopeId}' and scopeType eq '{scopeType}' and roleDefinitionId eq '{id}' (when policy_id is not set)
//   - GET /policies/roleManagementPolicies/{id}?$expand=rules
func constructResource(
	ctx context.Context,
	client *msgraphbetasdk.GraphServiceClient,
	data *RoleManagementPolicyResourceModel,
) (graphmodels.UnifiedRoleManagementPolicyable, error) {
	tflog.Debug(ctx, fmt.Sprintf("Constructing %s resource from model", ResourceName))

	policyId, err := resolvePolicyId(ctx, client, data)
	if err != nil {
		return nil, err
	}

	remotePolicy, err := getPolicyWithRules(ctx, client, policyId)
	if err != nil {
		return nil, fmt.Errorf("failed to read role management policy %s: %w", policyId, err)
	}

	if err := validateRules(data, remotePolicy); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	requestBody, err := constructRequestBody(ctx, data, policyId, remotePolicy)
	if err != nil {
		return nil, err
	}

	if err := constructors.DebugLogGraphObject(ctx, fmt.Sprintf("Final JSON to be sent to Graph API for resource %s", ResourceName), requestBody); err != nil {
		tflog.Error(ctx, "Failed to debug log object", map[string]any{
			"error": err.Error(),
		})
	}

	tflog.Debug(ctx, fmt.Sprintf("Finished constructing %s resource", ResourceName))

	return requestBody, nil
}

// constructRequestBody builds the policy PATCH body from the configured rules only. Rules that are not configured are
// left out so Graph leaves them unchanged. The policy is used for the notification and recipient types, which Graph
// requires on every notification rule update and which are fixed by the rule ID.
func constructRequestBody(
	ctx context.Context,
	data *RoleManagementPolicyResourceModel,
	policyId string,
	remotePolicy graphmodels.UnifiedRoleManagementPolicyable,
) (graphmodels.UnifiedRoleManagementPolicyable, error) {
	requestBody := graphmodels.NewUnifiedRoleManagementPolicy()
	requestBody.SetId(&policyId)

	remoteRules := rulesById(remotePolicy)
	rules := make([]graphmodels.UnifiedRoleManagementPolicyRuleable, 0)

	for _, ruleId := range sortedKeys(data.ExpirationRules) {
		rule, err := constructExpirationRule(ruleId, data.ExpirationRules[ruleId])
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}

	for _, ruleId := range sortedKeys(data.EnablementRules) {
		rule, err := constructEnablementRule(ctx, ruleId, data.EnablementRules[ruleId])
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}

	for _, ruleId := range sortedKeys(data.NotificationRules) {
		rule, err := constructNotificationRule(ctx, ruleId, data.NotificationRules[ruleId], remoteRules[ruleId])
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}

	for _, ruleId := range sortedKeys(data.AuthenticationContextRules) {
		rules = append(rules, constructAuthenticationContextRule(ruleId, data.AuthenticationContextRules[ruleId]))
	}

	requestBody.SetRules(rules)

	return requestBody, nil
}

// constructExpirationRule builds a unifiedRoleManagementPolicyExpirationRule.
func constructExpirationRule(ruleId string, data ExpirationRuleModel) (graphmodels.UnifiedRoleManagementPolicyRuleable, error) {
	rule := graphmodels.NewUnifiedRoleManagementPolicyExpirationRule()
	rule.SetId(&ruleId)
	convert.FrameworkToGraphBool(data.IsExpirationRequired, rule.SetIsExpirationRequired)
	if err := convert.FrameworkToGraphISODuration(data.MaximumDuration, rule.SetMaximumDuration); err != nil {
		return nil, fmt.Errorf("failed to set maximum_duration of expiration rule %s: %w", ruleId, err)
	}
	return rule, nil
}

// constructEnablementRule builds a unifiedRoleManagementPolicyEnablementRule. An empty set is sent as an empty array.
func constructEnablementRule(ctx context.Context, ruleId string, data EnablementRuleModel) (graphmodels.UnifiedRoleManagementPolicyRuleable, error) {
	rule := graphmodels.NewUnifiedRoleManagementPolicyEnablementRule()
	rule.SetId(&ruleId)
	if err := convert.FrameworkToGraphStringSetOmitUnset(ctx, data.EnabledRules, rule.SetEnabledRules); err != nil {
		return nil, fmt.Errorf("failed to set enabled_rules of enablement rule %s: %w", ruleId, err)
	}
	return rule, nil
}

// constructNotificationRule builds a unifiedRoleManagementPolicyNotificationRule, copying the notification and
// recipient types from the existing rule.
func constructNotificationRule(
	ctx context.Context,
	ruleId string,
	data NotificationRuleModel,
	remoteRule graphmodels.UnifiedRoleManagementPolicyRuleable,
) (graphmodels.UnifiedRoleManagementPolicyRuleable, error) {
	rule := graphmodels.NewUnifiedRoleManagementPolicyNotificationRule()
	rule.SetId(&ruleId)
	if existing, ok := remoteRule.(graphmodels.UnifiedRoleManagementPolicyNotificationRuleable); ok {
		rule.SetNotificationType(existing.GetNotificationType())
		rule.SetRecipientType(existing.GetRecipientType())
	}
	convert.FrameworkToGraphString(data.NotificationLevel, rule.SetNotificationLevel)
	convert.FrameworkToGraphBool(data.IsDefaultRecipientsEnabled, rule.SetIsDefaultRecipientsEnabled)
	if err := convert.FrameworkToGraphStringSetOmitUnset(ctx, data.NotificationRecipients, rule.SetNotificationRecipients); err != nil {
		return nil, fmt.Errorf("failed to set notification_recipients of notification rule %s: %w", ruleId, err)
	}
	return rule, nil
}

// constructAuthenticationContextRule builds a unifiedRoleManagementPolicyAuthenticationContextRule.
func constructAuthenticationContextRule(ruleId string, data AuthenticationContextRuleModel) graphmodels.UnifiedRoleManagementPolicyRuleable {
	rule := graphmodels.NewUnifiedRoleManagementPolicyAuthenticationContextRule()
	rule.SetId(&ruleId)
	convert.FrameworkToGraphBool(data.IsEnabled, rule.SetIsEnabled)
	convert.FrameworkToGraphString(data.ClaimValue, rule.SetClaimValue)
	return rule
}
