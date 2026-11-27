package convert

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework-jsontypes/jsontypes"
	"github.com/stretchr/testify/require"
)

func TestUnitGraphToFrameworkNormalizedJSON(t *testing.T) {
	prior := jsontypes.NewNormalizedValue("{ \"enabled\": true, \"scopes\": [\"read\"] }")
	same, err := GraphToFrameworkNormalizedJSON(context.Background(), json.RawMessage(`{"scopes":["read"],"enabled":true}`), prior)
	require.NoError(t, err)
	require.Equal(t, prior.ValueString(), same.ValueString())
	changed, err := GraphToFrameworkNormalizedJSON(context.Background(), json.RawMessage(`{"scopes":["write"],"enabled":false}`), prior)
	require.NoError(t, err)
	require.NotEqual(t, prior.ValueString(), changed.ValueString())
	empty, err := GraphToFrameworkNormalizedJSON(context.Background(), nil, prior)
	require.NoError(t, err)
	require.True(t, empty.IsNull())
	_, err = GraphToFrameworkNormalizedJSON(context.Background(), json.RawMessage(`{"invalid":`), prior)
	require.Error(t, err)
}
