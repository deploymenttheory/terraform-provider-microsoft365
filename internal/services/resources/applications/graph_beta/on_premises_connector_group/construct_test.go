package graphBetaApplicationsOnPremisesConnectorGroup

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestConstructResourceSendsOnlyWritableFields(t *testing.T) {
	tests := []struct {
		name     string
		plan     *OnPremisesConnectorGroupResourceModel
		state    *OnPremisesConnectorGroupResourceModel
		expected string
	}{
		{
			name:     "create with region",
			plan:     &OnPremisesConnectorGroupResourceModel{Name: types.StringValue("unit-test-connector-group"), Region: types.StringValue("japan"), ConnectorGroupType: types.StringValue("applicationProxy"), IsDefault: types.BoolValue(true)},
			expected: `{"name":"unit-test-connector-group","region":"japan"}`,
		},
		{
			name:     "create without region",
			plan:     &OnPremisesConnectorGroupResourceModel{Name: types.StringValue("unit-test-connector-group"), Region: types.StringNull()},
			expected: `{"name":"unit-test-connector-group"}`,
		},
		{
			name:     "update omits unchanged region",
			plan:     &OnPremisesConnectorGroupResourceModel{Name: types.StringValue("unit-test-connector-group-renamed"), Region: types.StringValue("japan")},
			state:    &OnPremisesConnectorGroupResourceModel{Name: types.StringValue("unit-test-connector-group"), Region: types.StringValue("japan")},
			expected: `{"name":"unit-test-connector-group-renamed"}`,
		},
		{
			name:     "update sends changed region",
			plan:     &OnPremisesConnectorGroupResourceModel{Name: types.StringValue("unit-test-connector-group"), Region: types.StringValue("eur")},
			state:    &OnPremisesConnectorGroupResourceModel{Name: types.StringValue("unit-test-connector-group"), Region: types.StringValue("japan")},
			expected: `{"name":"unit-test-connector-group","region":"eur"}`,
		},
		{
			name:     "update sends region when state has none",
			plan:     &OnPremisesConnectorGroupResourceModel{Name: types.StringValue("unit-test-connector-group"), Region: types.StringValue("eur")},
			state:    &OnPremisesConnectorGroupResourceModel{Name: types.StringValue("unit-test-connector-group"), Region: types.StringNull()},
			expected: `{"name":"unit-test-connector-group","region":"eur"}`,
		},
		{
			name:     "update omits unknown region",
			plan:     &OnPremisesConnectorGroupResourceModel{Name: types.StringValue("unit-test-connector-group"), Region: types.StringUnknown()},
			state:    &OnPremisesConnectorGroupResourceModel{Name: types.StringValue("unit-test-connector-group"), Region: types.StringValue("japan")},
			expected: `{"name":"unit-test-connector-group"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, err := json.Marshal(constructResource(context.Background(), tt.plan, tt.state))
			if err != nil {
				t.Fatalf("marshal request body: %v", err)
			}
			if string(body) != tt.expected {
				t.Fatalf("request body = %s, expected %s", body, tt.expected)
			}
		})
	}
}
