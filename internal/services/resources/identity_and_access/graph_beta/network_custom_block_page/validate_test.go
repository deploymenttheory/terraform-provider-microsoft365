package graphBetaNetworkCustomBlockPage

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestUnitResourceNetworkCustomBlockPage_BodyBoundaries(t *testing.T) {
	for _, tc := range []struct {
		name    string
		body    types.String
		invalid bool
	}{
		{"empty", types.StringValue(""), true},
		{"whitespace", types.StringValue(" \t\n "), true},
		{"one", types.StringValue("a"), false},
		{"maximum", types.StringValue(strings.Repeat("a", 1024)), false},
		{"too_long", types.StringValue(strings.Repeat("a", 1025)), true},
		{"japanese", types.StringValue(strings.Repeat("あ", 1024)), false},
		{"emoji_maximum", types.StringValue(strings.Repeat("😀", 512)), false},
		{"emoji_too_long", types.StringValue(strings.Repeat("😀", 513)), true},
		{"unknown", types.StringUnknown(), false},
		{"null", types.StringNull(), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var resp validator.StringResponse
			blockMessageValidator{}.ValidateString(context.Background(), validator.StringRequest{Path: path.Root("configuration").AtName("body"), ConfigValue: tc.body}, &resp)
			require.Equal(t, tc.invalid, resp.Diagnostics.HasError(), "%v", resp.Diagnostics)
		})
	}
}

func TestUnitResourceNetworkCustomBlockPage_ConstructRejectsInvalidConfiguration(t *testing.T) {
	for _, body := range []string{"", " \n ", strings.Repeat("😀", 513)} {
		model := testModel()
		model.Configuration = testConfiguration(body)
		_, err := constructResource(&model)
		require.Error(t, err)
	}
}
