package graphBetaRoleManagementPolicy

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/microsoft/kiota-abstractions-go/serialization"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/require"
)

func TestUnitParseDuration(t *testing.T) {
	valid := map[string]time.Duration{
		"P180D":       180 * 24 * time.Hour,
		"PT8H":        8 * time.Hour,
		"PT480M":      8 * time.Hour,
		"P1DT12H":     36 * time.Hour,
		"PT1H30M15S":  time.Hour + 30*time.Minute + 15*time.Second,
		"PT0S":        0,
		"P0D":         0,
		"P1DT0H0M30S": 24*time.Hour + 30*time.Second,
	}
	for value, want := range valid {
		got, err := parseDuration(value)
		require.NoError(t, err, value)
		require.Equal(t, want, got, value)
	}

	for _, value := range []string{"", "P", "PT", "P1DT", "P2W", "P1Y", "P1M", "8H", "PT1.5H", "P-1D"} {
		_, err := parseDuration(value)
		require.Error(t, err, value)
	}
}

func TestUnitDurationsEquivalent(t *testing.T) {
	require.True(t, durationsEquivalent("PT24H", "P1D"))
	require.True(t, durationsEquivalent("PT480M", "PT8H"))
	require.True(t, durationsEquivalent("P0D", "PT0S"))
	require.False(t, durationsEquivalent("P90D", "P60D"))
	require.False(t, durationsEquivalent("P2W", "P14D"))
}

func TestUnitScopeFilter(t *testing.T) {
	require.Equal(t, "scopeId eq '/' and scopeType eq 'DirectoryRole'", scopeFilter("/", "DirectoryRole", ""))
	require.Equal(t, "scopeId eq 'a''b' and scopeType eq 'Group' and roleDefinitionId eq 'member'", scopeFilter("a'b", "Group", "member"))
}

func TestUnitPolicyScopeDefaults(t *testing.T) {
	scopeId, scopeType := policyScope(&RoleManagementPolicyResourceModel{ScopeId: types.StringUnknown(), ScopeType: types.StringNull()})
	require.Equal(t, "/", scopeId)
	require.Equal(t, "DirectoryRole", scopeType)

	scopeId, scopeType = policyScope(&RoleManagementPolicyResourceModel{ScopeId: types.StringValue("group-id"), ScopeType: types.StringValue("Group")})
	require.Equal(t, "group-id", scopeId)
	require.Equal(t, "Group", scopeType)
}

func TestUnitValidateRules(t *testing.T) {
	expiration := graphmodels.NewUnifiedRoleManagementPolicyExpirationRule()
	expiration.SetId(new("Expiration_Admin_Assignment"))
	enablement := graphmodels.NewUnifiedRoleManagementPolicyEnablementRule()
	enablement.SetId(new("Enablement_Admin_Assignment"))
	policy := graphmodels.NewUnifiedRoleManagementPolicy()
	policy.SetId(new(testPolicyId))
	policy.SetRules([]graphmodels.UnifiedRoleManagementPolicyRuleable{expiration, enablement})

	valid := &RoleManagementPolicyResourceModel{
		ExpirationRules: map[string]ExpirationRuleModel{"Expiration_Admin_Assignment": {}},
		EnablementRules: map[string]EnablementRuleModel{"Enablement_Admin_Assignment": {}},
	}
	require.NoError(t, validateRules(valid, policy))

	invalid := &RoleManagementPolicyResourceModel{
		ExpirationRules:            map[string]ExpirationRuleModel{"Enablement_Admin_Assignment": {}},
		AuthenticationContextRules: map[string]AuthenticationContextRuleModel{"AuthenticationContext_EndUser_Assignment": {}},
	}
	err := validateRules(invalid, policy)
	require.Error(t, err)
	require.Contains(t, err.Error(), `expiration_rules: rule "Enablement_Admin_Assignment" is a unifiedRoleManagementPolicyEnablementRule, not a unifiedRoleManagementPolicyExpirationRule`)
	require.Contains(t, err.Error(), `authentication_context_rules: rule "AuthenticationContext_EndUser_Assignment" does not exist on the policy`)
	require.Contains(t, err.Error(), "Rules on the policy: Enablement_Admin_Assignment (unifiedRoleManagementPolicyEnablementRule), Expiration_Admin_Assignment (unifiedRoleManagementPolicyExpirationRule)")
}

func TestUnitMapRemoteResourceStateTracksOnlyConfiguredRules(t *testing.T) {
	ctx := context.Background()
	policy := remotePolicy()
	expiration := graphmodels.NewUnifiedRoleManagementPolicyExpirationRule()
	expiration.SetId(new("Expiration_EndUser_Assignment"))
	expiration.SetIsExpirationRequired(new(true))
	duration, err := serialization.ParseISODuration("P1D")
	require.NoError(t, err)
	expiration.SetMaximumDuration(duration)
	policy.SetRules(append(policy.GetRules(), expiration))

	prior := RoleManagementPolicyResourceModel{
		ExpirationRules: map[string]ExpirationRuleModel{
			"Expiration_EndUser_Assignment": {IsExpirationRequired: types.BoolValue(true), MaximumDuration: types.StringValue("PT24H")},
		},
	}
	mapped := MapRemoteResourceStateToTerraform(ctx, prior, policy, false)
	require.Len(t, mapped.ExpirationRules, 1)
	require.Equal(t, "PT24H", mapped.ExpirationRules["Expiration_EndUser_Assignment"].MaximumDuration.ValueString(), "an equivalent duration keeps the configured spelling")
	require.Nil(t, mapped.NotificationRules, "unconfigured rule types stay null")

	imported := MapRemoteResourceStateToTerraform(ctx, RoleManagementPolicyResourceModel{}, policy, true)
	require.Len(t, imported.ExpirationRules, 1)
	require.Equal(t, "P1D", imported.ExpirationRules["Expiration_EndUser_Assignment"].MaximumDuration.ValueString())
	require.Len(t, imported.NotificationRules, 1)
	require.Equal(t, 0, len(imported.NotificationRules["Notification_Approver_EndUser_Assignment"].NotificationRecipients.Elements()))
	require.Nil(t, imported.EnablementRules, "rule types the policy does not have stay null after import")
}

func TestUnitRulesMatch(t *testing.T) {
	expected := &RoleManagementPolicyResourceModel{
		ExpirationRules: map[string]ExpirationRuleModel{
			"Expiration_Admin_Assignment": {IsExpirationRequired: types.BoolValue(false), MaximumDuration: types.StringValue("PT24H")},
		},
		AuthenticationContextRules: map[string]AuthenticationContextRuleModel{
			"AuthenticationContext_EndUser_Assignment": {IsEnabled: types.BoolValue(true), ClaimValue: types.StringValue("c1")},
		},
	}

	// A stale replica still returns the values from before the PATCH.
	stale := &RoleManagementPolicyResourceModel{
		ExpirationRules:            map[string]ExpirationRuleModel{"Expiration_Admin_Assignment": {IsExpirationRequired: types.BoolValue(false), MaximumDuration: types.StringValue("P180D")}},
		AuthenticationContextRules: map[string]AuthenticationContextRuleModel{"AuthenticationContext_EndUser_Assignment": {IsEnabled: types.BoolValue(false), ClaimValue: types.StringNull()}},
	}
	require.False(t, rulesMatch(expected, stale))

	// A rule that was dropped from state because it is missing in Graph never matches.
	require.False(t, rulesMatch(expected, &RoleManagementPolicyResourceModel{}))

	fresh := &RoleManagementPolicyResourceModel{
		ExpirationRules:            map[string]ExpirationRuleModel{"Expiration_Admin_Assignment": {IsExpirationRequired: types.BoolValue(false), MaximumDuration: types.StringValue("P1D")}},
		AuthenticationContextRules: map[string]AuthenticationContextRuleModel{"AuthenticationContext_EndUser_Assignment": {IsEnabled: types.BoolValue(true), ClaimValue: types.StringValue("c1")}},
	}
	require.True(t, rulesMatch(expected, fresh))
}
