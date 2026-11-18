package graphBetaGroup

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestUnitGroupRequestDefaults(t *testing.T) {
	for _, tc := range []struct {
		name            string
		update, unified bool
		visibility      types.String
		role            types.Bool
		wantVisibility  *string
		wantRole        *bool
	}{
		{"create unknown fields", false, false, types.StringUnknown(), types.BoolUnknown(), new("Private"), new(false)},
		{"create null fields", false, false, types.StringNull(), types.BoolNull(), new("Private"), new(false)},
		{"create Microsoft 365 group", false, true, types.StringUnknown(), types.BoolUnknown(), new("Private"), new(false)},
		{"create explicit public group", false, false, types.StringValue("Public"), types.BoolValue(false), new("Public"), new(false)},
		{"create role assignable group", false, false, types.StringValue("Private"), types.BoolValue(true), new("Private"), new(true)},
		{"update null fields", true, false, types.StringNull(), types.BoolNull(), nil, nil},
		{"update unknown fields", true, false, types.StringUnknown(), types.BoolUnknown(), nil, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			data := GroupResourceModel{
				DisplayName:     types.StringValue("group-defaults-test"),
				MailNickname:    types.StringValue("group-defaults-test"),
				MailEnabled:     types.BoolValue(tc.unified),
				SecurityEnabled: types.BoolValue(!tc.unified),
				Visibility:      tc.visibility, IsAssignableToRole: tc.role,
				GroupTypes: types.SetNull(types.StringType),
			}
			if tc.unified {
				data.GroupTypes = types.SetValueMust(types.StringType, []attr.Value{types.StringValue("Unified")})
			}
			body, err := constructResource(context.Background(), &data, tc.update)
			require.NoError(t, err)
			require.Equal(t, tc.wantVisibility, body.GetVisibility())
			require.Equal(t, tc.wantRole, body.GetIsAssignableToRole())
		})
	}
}
