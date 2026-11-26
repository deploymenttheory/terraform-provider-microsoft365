package convert

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGraphToFrameworkBoolChoice(t *testing.T) {
	for _, test := range []struct {
		name, value, definition string
		want                    types.Bool
	}{
		{"false", "setting_0", "setting", types.BoolValue(false)},
		{"true", "setting_1", "setting", types.BoolValue(true)},
		{"unknown choice", "setting_2", "setting", types.BoolNull()},
		{"wrong setting", "other_1", "setting", types.BoolNull()},
		{"missing definition", "_1", "", types.BoolNull()},
		{"empty value", "", "setting", types.BoolNull()},
	} {
		t.Run(test.name, func(t *testing.T) {
			require.Equal(t, test.want, GraphToFrameworkBoolChoice(&test.value, test.definition))
		})
	}
	require.Equal(t, types.BoolNull(), GraphToFrameworkBoolChoice(nil, "setting"))
}
