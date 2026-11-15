package graphBetaRoleManagementPolicy

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
)

// errInvalidRules reports configured rules that do not exist on the policy or have another type.
var errInvalidRules = errors.New("invalid rules")

// isoDurationValidator validates that a string is an ISO 8601 duration Graph accepts for expiration rules.
type isoDurationValidator struct{}

var _ validator.String = isoDurationValidator{}

func (v isoDurationValidator) Description(ctx context.Context) string {
	return "value must be an ISO 8601 duration of days, hours, minutes and seconds, for example P180D, PT8H or P1DT12H"
}

func (v isoDurationValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

func (v isoDurationValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}
	if _, err := parseDuration(req.ConfigValue.ValueString()); err != nil {
		resp.Diagnostics.AddAttributeError(req.Path, "Invalid ISO 8601 duration", err.Error())
	}
}

// ValidateConfig checks that every authentication context rule sets claim_value exactly when it is enabled.
// Graph rejects an enabled rule without a claim value and clears the claim value when the rule is disabled.
func (r *RoleManagementPolicyResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var rules types.Map
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("authentication_context_rules"), &rules)...)
	if resp.Diagnostics.HasError() || rules.IsNull() || rules.IsUnknown() {
		return
	}

	var elements map[string]AuthenticationContextRuleModel
	resp.Diagnostics.Append(rules.ElementsAs(ctx, &elements, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	for ruleId, rule := range elements {
		if rule.IsEnabled.IsUnknown() || rule.ClaimValue.IsUnknown() {
			continue
		}
		attributePath := path.Root("authentication_context_rules").AtMapKey(ruleId).AtName("claim_value")
		if rule.IsEnabled.ValueBool() && rule.ClaimValue.IsNull() {
			detail := fmt.Sprintf("Authentication context rule %q is enabled, so claim_value must be set to an authentication context ID such as c1.", ruleId)
			resp.Diagnostics.AddAttributeError(attributePath, "Missing claim_value", detail)
		}
		if !rule.IsEnabled.ValueBool() && !rule.ClaimValue.IsNull() {
			detail := fmt.Sprintf("Authentication context rule %q is disabled, so claim_value must be omitted. Graph clears the claim value of a disabled rule.", ruleId)
			resp.Diagnostics.AddAttributeError(attributePath, "Unexpected claim_value", detail)
		}
	}
}

// validateRules checks that every configured rule exists on the policy with the type of the attribute that configures it.
func validateRules(data *RoleManagementPolicyResourceModel, policy graphmodels.UnifiedRoleManagementPolicyable) error {
	remoteRules := rulesById(policy)

	var problems []string
	check := func(attribute, expectedType string, ruleIds []string) {
		for _, ruleId := range ruleIds {
			rule, ok := remoteRules[ruleId]
			if !ok {
				problems = append(problems, fmt.Sprintf("%s: rule %q does not exist on the policy", attribute, ruleId))
				continue
			}
			if actualType := ruleODataType(rule); actualType != expectedType {
				actualName, expectedName := strings.TrimPrefix(actualType, "#microsoft.graph."), strings.TrimPrefix(expectedType, "#microsoft.graph.")
				problems = append(problems, fmt.Sprintf("%s: rule %q is a %s, not a %s", attribute, ruleId, actualName, expectedName))
			}
		}
	}

	check("expiration_rules", odataTypeExpirationRule, sortedKeys(data.ExpirationRules))
	check("enablement_rules", odataTypeEnablementRule, sortedKeys(data.EnablementRules))
	check("notification_rules", odataTypeNotificationRule, sortedKeys(data.NotificationRules))
	check("authentication_context_rules", odataTypeAuthenticationContextRule, sortedKeys(data.AuthenticationContextRules))

	if len(problems) > 0 {
		available := make([]string, 0, len(remoteRules))
		for ruleId, rule := range remoteRules {
			available = append(available, fmt.Sprintf("%s (%s)", ruleId, strings.TrimPrefix(ruleODataType(rule), "#microsoft.graph.")))
		}
		slices.Sort(available)
		summary := strings.Join(problems, "\n  - ")
		policyId := stringValue(policy.GetId())
		return fmt.Errorf("%w for role management policy %s:\n  - %s\nRules on the policy: %s", errInvalidRules, policyId, summary, strings.Join(available, ", "))
	}
	return nil
}

// sortedKeys returns the keys of a rule map in a stable order.
func sortedKeys[T any](rules map[string]T) []string {
	keys := make([]string, 0, len(rules))
	for key := range rules {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	return keys
}

func stringValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
