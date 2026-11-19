package graphBetaServicePrincipalClaimsPolicyJSON

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const emptyPolicy = `{"includeBasicClaimSet":true,"includeApplicationIdInIssuer":false,"claims":[]}`

func TestUnitClaimsPolicyJSONSemanticComparison(t *testing.T) {
	for _, tc := range []struct {
		name, left, right string
		equal             bool
	}{
		{name: "format and key order", left: emptyPolicy, right: `{ "claims":[],"includeApplicationIdInIssuer":false,"includeBasicClaimSet":true }`, equal: true},
		{name: "server metadata and nullable root", left: emptyPolicy, right: `{"id":"server-id","@odata.context":"metadata","includeBasicClaimSet":true,"includeApplicationIdInIssuer":false,"audienceOverride":null,"groupFilter":null,"claims":[]}`, equal: true},
		{name: "version remains meaningful", left: strings.Replace(emptyPolicy, `"claims":[]`, `"claims":[],"version":"one"`, 1), right: strings.Replace(emptyPolicy, `"claims":[]`, `"claims":[],"version":"two"`, 1)},
		{name: "root flag changes", left: emptyPolicy, right: strings.Replace(emptyPolicy, `"includeBasicClaimSet":true`, `"includeBasicClaimSet":false`, 1)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			equal, diags := NewPolicyJSONValue(tc.left).StringSemanticEquals(context.Background(), NewPolicyJSONValue(tc.right))
			require.False(t, diags.HasError())
			require.Equal(t, tc.equal, equal)
		})
	}
	base := `{"includeBasicClaimSet":true,"includeApplicationIdInIssuer":false,"claims":[{"@odata.type":"#microsoft.graph.samlNameIdClaim","nameIdFormat":"persistent","configurations":[{"attribute":{"@odata.type":"#microsoft.graph.sourcedAttribute","id":"userprincipalname","source":"user"},"transformations":[{"@odata.type":"#microsoft.graph.trimTransformation","value":" "},{"@odata.type":"#microsoft.graph.toLowercaseTransformation"}]}]}]}`
	added := strings.Replace(base, `"nameIdFormat":`, `"serviceProviderNameQualifier":null,"nameIdFormat":`, 1)
	added = strings.Replace(added, `"value":" "`, `"input":null,"value":" "`, 1)
	equal, _ := NewPolicyJSONValue(base).StringSemanticEquals(context.Background(), NewPolicyJSONValue(added))
	require.True(t, equal, "observed service-added nulls should not cause drift")
	changed := strings.Replace(base, "#microsoft.graph.sourcedAttribute", "#microsoft.graph.valueBasedAttribute", 1)
	equal, _ = NewPolicyJSONValue(base).StringSemanticEquals(context.Background(), NewPolicyJSONValue(changed))
	require.False(t, equal, "nested discriminators must remain meaningful")
	var object map[string]any
	require.NoError(t, json.Unmarshal([]byte(base), &object))
	transforms := object["claims"].([]any)[0].(map[string]any)["configurations"].([]any)[0].(map[string]any)["transformations"].([]any)
	transforms[0], transforms[1] = transforms[1], transforms[0]
	reversed, err := json.Marshal(object)
	require.NoError(t, err)
	equal, _ = NewPolicyJSONValue(base).StringSemanticEquals(context.Background(), NewPolicyJSONValue(string(reversed)))
	require.False(t, equal, "ordered transformations cannot be sorted or ignored")
}

func TestUnitClaimsPolicyJSONLosslessConstruction(t *testing.T) {
	raw := `{"id":"readonly","@odata.context":"metadata","includeBasicClaimSet":false,"includeApplicationIdInIssuer":false,"claims":[],"futureField":{"number":9007199254740993,"@odata.type":"#future.type"},"groupFilter":{"value":"owned-group"}}`
	body, err := constructResource(NewPolicyJSONValue(raw))
	require.NoError(t, err)
	require.NotContains(t, string(body), `"id":"readonly"`)
	require.NotContains(t, string(body), "@odata.context")
	require.Contains(t, string(body), "9007199254740993")
	require.Contains(t, string(body), "#future.type")
	require.Contains(t, string(body), "owned-group")
}

func TestUnitClaimsPolicyJSONValidation(t *testing.T) {
	for _, raw := range []string{"null", "[]", `{}`, `{"includeBasicClaimSet":"true","includeApplicationIdInIssuer":false,"claims":[]}`, emptyPolicy + ` {}`, `{"includeBasicClaimSet":true,"includeApplicationIdInIssuer":false,"claims":null}`} {
		_, err := decodePolicyJSON([]byte(raw))
		require.Error(t, err, raw)
	}
}
