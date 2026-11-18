package schema

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestUnitDeviceConfigurationAssignmentFilter(t *testing.T) {
	ctx := context.Background()
	assignments := DeviceConfigurationWithAllGroupAssignmentsAndFilterSchema()
	filter := assignments.NestedObject.Attributes["filter_id"].(resourceschema.StringAttribute)
	require.True(t, filter.Optional)
	require.False(t, filter.Computed, "omitted filter IDs must remain null")
	require.Nil(t, filter.Default, "unfiltered assignments must not acquire a zero-GUID default")

	schema := resourceschema.Schema{Attributes: assignments.NestedObject.Attributes}
	objectType := schema.Type().(types.ObjectType)
	for _, tc := range []struct {
		name, filterType string
		id               types.String
		valid            bool
	}{
		{"unfiltered", "none", types.StringNull(), true},
		{"zero GUID", "none", types.StringValue("00000000-0000-0000-0000-000000000000"), false},
		{"include", "include", types.StringValue("11111111-1111-4111-8111-111111111111"), true},
		{"exclude", "exclude", types.StringValue("11111111-1111-4111-8111-111111111111"), true},
		{"missing include filter", "include", types.StringNull(), false},
		{"missing exclude filter", "exclude", types.StringNull(), false},
		{"invalid GUID", "include", types.StringValue("invalid"), false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			object := types.ObjectValueMust(objectType.AttrTypes, map[string]attr.Value{
				"type":        types.StringValue("groupAssignmentTarget"),
				"group_id":    types.StringValue("22222222-2222-4222-8222-222222222222"),
				"filter_type": types.StringValue(tc.filterType), "filter_id": tc.id,
			})
			raw, err := object.ToTerraformValue(ctx)
			require.NoError(t, err)
			valid := true
			for _, rule := range filter.Validators {
				var response validator.StringResponse
				rule.ValidateString(ctx, validator.StringRequest{
					Path: path.Root("filter_id"), PathExpression: path.MatchRoot("filter_id"),
					Config: tfsdk.Config{Schema: schema, Raw: raw}, ConfigValue: tc.id,
				}, &response)
				valid = valid && !response.Diagnostics.HasError()
			}
			require.Equal(t, tc.valid, valid)
		})
	}
}
