package graphBetaRoleManagementPolicy

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	jsonserialization "github.com/microsoft/kiota-serialization-json-go"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/require"
)

const testPolicyId = "DirectoryRole_00000000-0000-0000-0000-000000000001_22222222-2222-2222-2222-222222222222"

// serializeBody serializes a request body the way the Graph SDK sends it and returns the rules keyed by ID.
func serializeBody(t *testing.T, body graphmodels.UnifiedRoleManagementPolicyable) (map[string]any, map[string]map[string]any) {
	t.Helper()
	writer := jsonserialization.NewJsonSerializationWriter()
	require.NoError(t, writer.WriteObjectValue("", body))
	content, err := writer.GetSerializedContent()
	require.NoError(t, err)

	var payload map[string]any
	require.NoError(t, json.Unmarshal(content, &payload))
	rules := map[string]map[string]any{}
	items, _ := payload["rules"].([]any)
	for _, item := range items {
		rule, _ := item.(map[string]any)
		id, _ := rule["id"].(string)
		rules[id] = rule
	}
	return payload, rules
}

// remotePolicy builds a policy holding one notification rule, as Graph returns it.
func remotePolicy() graphmodels.UnifiedRoleManagementPolicyable {
	notification := graphmodels.NewUnifiedRoleManagementPolicyNotificationRule()
	notification.SetId(new("Notification_Approver_EndUser_Assignment"))
	notification.SetNotificationType(new("Email"))
	notification.SetRecipientType(new("Approver"))
	policy := graphmodels.NewUnifiedRoleManagementPolicy()
	policy.SetId(new(testPolicyId))
	policy.SetRules([]graphmodels.UnifiedRoleManagementPolicyRuleable{notification})
	return policy
}

func TestUnitConstructRequestBodySendsOnlyConfiguredRules(t *testing.T) {
	ctx := context.Background()
	data := &RoleManagementPolicyResourceModel{
		ExpirationRules: map[string]ExpirationRuleModel{
			"Expiration_Admin_Assignment": {IsExpirationRequired: types.BoolValue(false), MaximumDuration: types.StringValue("P90D")},
		},
	}

	body, err := constructRequestBody(ctx, data, testPolicyId, remotePolicy())
	require.NoError(t, err)
	payload, rules := serializeBody(t, body)

	require.Equal(t, testPolicyId, payload["id"])
	require.Len(t, rules, 1)
	rule := rules["Expiration_Admin_Assignment"]
	require.Equal(t, "#microsoft.graph.unifiedRoleManagementPolicyExpirationRule", rule["@odata.type"])
	require.Equal(t, false, rule["isExpirationRequired"])
	require.Equal(t, "P90D", rule["maximumDuration"])
	require.NotContains(t, rule, "target")
}

func TestUnitConstructRequestBodySerializesDurationsAsConfigured(t *testing.T) {
	ctx := context.Background()
	for _, duration := range []string{"P90D", "P365D", "PT8H", "PT30M", "P1DT12H", "PT1H30M"} {
		t.Run(duration, func(t *testing.T) {
			data := &RoleManagementPolicyResourceModel{
				ExpirationRules: map[string]ExpirationRuleModel{
					"Expiration_EndUser_Assignment": {IsExpirationRequired: types.BoolValue(true), MaximumDuration: types.StringValue(duration)},
				},
			}
			body, err := constructRequestBody(ctx, data, testPolicyId, remotePolicy())
			require.NoError(t, err)
			_, rules := serializeBody(t, body)
			sent, ok := rules["Expiration_EndUser_Assignment"]["maximumDuration"].(string)
			require.True(t, ok)
			require.True(t, durationsEquivalent(duration, sent), "sent %s for %s", sent, duration)
		})
	}
}

func TestUnitConstructRequestBodyRuleTypes(t *testing.T) {
	ctx := context.Background()
	emptySet := types.SetValueMust(types.StringType, []attr.Value{})
	data := &RoleManagementPolicyResourceModel{
		EnablementRules: map[string]EnablementRuleModel{
			"Enablement_Admin_Eligibility":  {EnabledRules: emptySet},
			"Enablement_EndUser_Assignment": {EnabledRules: types.SetValueMust(types.StringType, []attr.Value{types.StringValue("MultiFactorAuthentication")})},
		},
		NotificationRules: map[string]NotificationRuleModel{
			"Notification_Approver_EndUser_Assignment": {NotificationLevel: types.StringValue("Critical"), IsDefaultRecipientsEnabled: types.BoolValue(false), NotificationRecipients: emptySet},
		},
		AuthenticationContextRules: map[string]AuthenticationContextRuleModel{
			"AuthenticationContext_EndUser_Assignment": {IsEnabled: types.BoolValue(false), ClaimValue: types.StringNull()},
		},
	}

	body, err := constructRequestBody(ctx, data, testPolicyId, remotePolicy())
	require.NoError(t, err)
	_, rules := serializeBody(t, body)
	require.Len(t, rules, 4)

	// An empty set is sent as an empty array so Graph clears the list.
	require.Equal(t, []any{}, rules["Enablement_Admin_Eligibility"]["enabledRules"])
	require.Equal(t, []any{"MultiFactorAuthentication"}, rules["Enablement_EndUser_Assignment"]["enabledRules"])

	// Graph requires the notification and recipient types, which come from the existing rule.
	notification := rules["Notification_Approver_EndUser_Assignment"]
	require.Equal(t, "#microsoft.graph.unifiedRoleManagementPolicyNotificationRule", notification["@odata.type"])
	require.Equal(t, "Email", notification["notificationType"])
	require.Equal(t, "Approver", notification["recipientType"])
	require.Equal(t, "Critical", notification["notificationLevel"])
	require.Equal(t, false, notification["isDefaultRecipientsEnabled"])
	require.Equal(t, []any{}, notification["notificationRecipients"])

	authenticationContext := rules["AuthenticationContext_EndUser_Assignment"]
	require.Equal(t, "#microsoft.graph.unifiedRoleManagementPolicyAuthenticationContextRule", authenticationContext["@odata.type"])
	require.Equal(t, false, authenticationContext["isEnabled"])
	require.NotContains(t, authenticationContext, "claimValue")
}

func TestUnitConstructRequestBodyWithoutRules(t *testing.T) {
	body, err := constructRequestBody(context.Background(), &RoleManagementPolicyResourceModel{}, testPolicyId, remotePolicy())
	require.NoError(t, err)
	payload, rules := serializeBody(t, body)
	require.Empty(t, rules)
	require.Equal(t, []any{}, payload["rules"])
}
