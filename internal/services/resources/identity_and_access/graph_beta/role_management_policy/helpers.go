package graphBetaRoleManagementPolicy

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	msgraphbetasdk "github.com/microsoftgraph/msgraph-beta-sdk-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/microsoftgraph/msgraph-beta-sdk-go/policies"
)

const (
	scopeTypeDirectoryRole  = "DirectoryRole"
	scopeTypeGroup          = "Group"
	defaultDirectoryScopeId = "/"

	roleDefinitionIdRegex = `^([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}|member|owner)$`
	claimValueRegex       = `^c([1-9]|[1-9][0-9])$`

	odataTypeExpirationRule            = "#microsoft.graph.unifiedRoleManagementPolicyExpirationRule"
	odataTypeEnablementRule            = "#microsoft.graph.unifiedRoleManagementPolicyEnablementRule"
	odataTypeNotificationRule          = "#microsoft.graph.unifiedRoleManagementPolicyNotificationRule"
	odataTypeAuthenticationContextRule = "#microsoft.graph.unifiedRoleManagementPolicyAuthenticationContextRule"
)

var (
	enablementRuleValues    = []string{"MultiFactorAuthentication", "Justification", "Ticketing"}
	notificationLevelValues = []string{"All", "Critical", "None"}

	// isoDurationPattern matches the ISO 8601 durations Graph accepts for expiration rules: days and a time part.
	// Graph rejects week durations such as P2W.
	errInvalidDuration = errors.New("invalid ISO 8601 duration")
	errPolicyNotFound  = errors.New("role management policy not found")

	isoDurationPattern = regexp.MustCompile(`^P(?:(\d+)D)?(?:T(?:(\d+)H)?(?:(\d+)M)?(?:(\d+)S)?)?$`)
)

// parseDuration converts an ISO 8601 duration of days, hours, minutes and seconds into a time.Duration.
func parseDuration(value string) (time.Duration, error) {
	matches := isoDurationPattern.FindStringSubmatch(value)
	if matches == nil || value == "P" || strings.HasSuffix(value, "T") {
		return 0, fmt.Errorf("%w %q: use days, hours, minutes and seconds, for example P180D, PT8H or P1DT12H", errInvalidDuration, value)
	}

	units := []time.Duration{24 * time.Hour, time.Hour, time.Minute, time.Second}
	var total time.Duration
	for index, unit := range units {
		if matches[index+1] == "" {
			continue
		}
		count, err := strconv.Atoi(matches[index+1])
		if err != nil {
			return 0, fmt.Errorf("%w %q: %w", errInvalidDuration, value, err)
		}
		total += time.Duration(count) * unit
	}
	return total, nil
}

// durationsEquivalent reports whether two ISO 8601 durations describe the same length of time.
// Graph normalises durations on write, for example PT24H is returned as P1D and PT480M as PT8H.
func durationsEquivalent(left, right string) bool {
	if left == right {
		return true
	}
	leftDuration, err := parseDuration(left)
	if err != nil {
		return false
	}
	rightDuration, err := parseDuration(right)
	if err != nil {
		return false
	}
	return leftDuration == rightDuration
}

// ruleODataType returns the @odata.type of a policy rule, preferring the concrete SDK type.
func ruleODataType(rule graphmodels.UnifiedRoleManagementPolicyRuleable) string {
	switch rule.(type) {
	case graphmodels.UnifiedRoleManagementPolicyExpirationRuleable:
		return odataTypeExpirationRule
	case graphmodels.UnifiedRoleManagementPolicyEnablementRuleable:
		return odataTypeEnablementRule
	case graphmodels.UnifiedRoleManagementPolicyNotificationRuleable:
		return odataTypeNotificationRule
	case graphmodels.UnifiedRoleManagementPolicyAuthenticationContextRuleable:
		return odataTypeAuthenticationContextRule
	}
	if rule.GetOdataType() != nil {
		return *rule.GetOdataType()
	}
	return ""
}

// rulesById indexes the rules of a policy by rule ID.
func rulesById(policy graphmodels.UnifiedRoleManagementPolicyable) map[string]graphmodels.UnifiedRoleManagementPolicyRuleable {
	rules := make(map[string]graphmodels.UnifiedRoleManagementPolicyRuleable)
	if policy == nil {
		return rules
	}
	for _, rule := range policy.GetRules() {
		if rule == nil || rule.GetId() == nil {
			continue
		}
		rules[*rule.GetId()] = rule
	}
	return rules
}

// getPolicyWithRules reads a role management policy with its rules expanded.
//
// API Calls:
//   - GET /policies/roleManagementPolicies/{id}?$expand=rules
func getPolicyWithRules(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, policyId string) (graphmodels.UnifiedRoleManagementPolicyable, error) {
	config := &policies.RoleManagementPoliciesUnifiedRoleManagementPolicyItemRequestBuilderGetRequestConfiguration{
		QueryParameters: &policies.RoleManagementPoliciesUnifiedRoleManagementPolicyItemRequestBuilderGetQueryParameters{
			Expand: []string{"rules"},
		},
	}
	builder := client.Policies().RoleManagementPolicies().ByUnifiedRoleManagementPolicyId(policyId)
	return builder.Get(ctx, config) //nolint:wrapcheck // Keep the Kiota error type for shared Graph error handling.
}

// listPolicyAssignments lists the role management policy assignments that match an OData filter, following paging.
//
// API Calls:
//   - GET /policies/roleManagementPolicyAssignments?$filter=...
//
// Graph requires the filter to name scopeId and scopeType; it rejects unfiltered requests and filters on policyId.
func listPolicyAssignments(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, filter string) ([]graphmodels.UnifiedRoleManagementPolicyAssignmentable, error) {
	config := &policies.RoleManagementPolicyAssignmentsRequestBuilderGetRequestConfiguration{
		QueryParameters: &policies.RoleManagementPolicyAssignmentsRequestBuilderGetQueryParameters{
			Filter: &filter,
		},
	}

	response, err := client.Policies().RoleManagementPolicyAssignments().Get(ctx, config)
	if err != nil {
		return nil, err //nolint:wrapcheck // Keep the Kiota error type for shared Graph error handling.
	}

	var assignments []graphmodels.UnifiedRoleManagementPolicyAssignmentable
	for response != nil {
		assignments = append(assignments, response.GetValue()...)
		nextLink := response.GetOdataNextLink()
		if nextLink == nil || *nextLink == "" {
			break
		}
		response, err = client.Policies().RoleManagementPolicyAssignments().WithUrl(*nextLink).Get(ctx, nil)
		if err != nil {
			return nil, err //nolint:wrapcheck // Keep the Kiota error type for shared Graph error handling.
		}
	}
	return assignments, nil
}

// scopeFilter builds the OData filter for the assignments of a scope, optionally narrowed to one role definition.
func scopeFilter(scopeId, scopeType, roleDefinitionId string) string {
	filter := fmt.Sprintf("scopeId eq '%s' and scopeType eq '%s'", escapeODataString(scopeId), escapeODataString(scopeType))
	if roleDefinitionId != "" {
		filter += fmt.Sprintf(" and roleDefinitionId eq '%s'", escapeODataString(roleDefinitionId))
	}
	return filter
}

// escapeODataString escapes single quotes in an OData string literal.
func escapeODataString(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}

// resolvePolicyId returns the ID of the policy to adopt: policy_id when it is set, otherwise the policy assigned to
// role_definition_id at scope_id and scope_type.
func resolvePolicyId(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, data *RoleManagementPolicyResourceModel) (string, error) {
	if !data.PolicyId.IsNull() && !data.PolicyId.IsUnknown() && data.PolicyId.ValueString() != "" {
		return data.PolicyId.ValueString(), nil
	}

	roleDefinitionId := data.RoleDefinitionId.ValueString()
	scopeId, scopeType := policyScope(data)

	assignments, err := listPolicyAssignments(ctx, client, scopeFilter(scopeId, scopeType, roleDefinitionId))
	if err != nil {
		return "", fmt.Errorf("failed to look up the role management policy assignment for role definition %s at scope %s (%s): %w", roleDefinitionId, scopeId, scopeType, err)
	}

	policyIds := make(map[string]struct{})
	for _, assignment := range assignments {
		if assignment == nil || assignment.GetPolicyId() == nil {
			continue
		}
		if assignment.GetRoleDefinitionId() != nil && !strings.EqualFold(*assignment.GetRoleDefinitionId(), roleDefinitionId) {
			continue
		}
		policyIds[*assignment.GetPolicyId()] = struct{}{}
	}

	if len(policyIds) != 1 {
		scope := fmt.Sprintf("%s (%s)", scopeId, scopeType)
		return "", fmt.Errorf("%w: expected one policy assigned to role definition %s at scope %s, found %d", errPolicyNotFound, roleDefinitionId, scope, len(policyIds))
	}
	for policyId := range policyIds {
		return policyId, nil
	}
	return "", nil
}

// policyScope returns the configured scope, defaulting to the directory scope used by Microsoft Entra roles.
func policyScope(data *RoleManagementPolicyResourceModel) (string, string) {
	scopeId := defaultDirectoryScopeId
	if !data.ScopeId.IsNull() && !data.ScopeId.IsUnknown() && data.ScopeId.ValueString() != "" {
		scopeId = data.ScopeId.ValueString()
	}
	scopeType := scopeTypeDirectoryRole
	if !data.ScopeType.IsNull() && !data.ScopeType.IsUnknown() && data.ScopeType.ValueString() != "" {
		scopeType = data.ScopeType.ValueString()
	}
	return scopeId, scopeType
}

// lookupRoleDefinitionId finds the role definition a policy is assigned to by listing the assignments of its scope.
func lookupRoleDefinitionId(ctx context.Context, client *msgraphbetasdk.GraphServiceClient, policyId, scopeId, scopeType string) (string, error) {
	assignments, err := listPolicyAssignments(ctx, client, scopeFilter(scopeId, scopeType, ""))
	if err != nil {
		return "", err
	}
	for _, assignment := range assignments {
		if assignment == nil || assignment.GetPolicyId() == nil || assignment.GetRoleDefinitionId() == nil {
			continue
		}
		if *assignment.GetPolicyId() == policyId {
			return *assignment.GetRoleDefinitionId(), nil
		}
	}
	return "", nil
}
