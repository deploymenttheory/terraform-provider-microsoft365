package crud

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestSetPartialStatePreservesKnownValuesAndResolvesNestedUnknowns(t *testing.T) {
	ctx := context.Background()
	nestedTypes := map[string]attr.Type{"name": types.StringType, "remote_status": types.StringType}
	nested := types.ObjectValueMust(nestedTypes, map[string]attr.Value{"name": types.StringValue("configured"), "remote_status": types.StringUnknown()})
	type model struct {
		ID      types.String `tfsdk:"id"`
		Remote  types.String `tfsdk:"remote"`
		Entries types.List   `tfsdk:"entries"`
		Tags    types.Set    `tfsdk:"tags"`
	}
	state := tfsdk.State{Schema: schema.Schema{Attributes: map[string]schema.Attribute{
		"id":      schema.StringAttribute{Computed: true},
		"remote":  schema.StringAttribute{Computed: true},
		"entries": schema.ListAttribute{Optional: true, ElementType: types.ObjectType{AttrTypes: nestedTypes}},
		"tags":    schema.SetAttribute{Optional: true, ElementType: types.StringType},
	}}}
	planned := model{ID: types.StringValue("allocated-resource"), Remote: types.StringUnknown(), Entries: types.ListValueMust(nested.Type(ctx), []attr.Value{nested}), Tags: types.SetValueMust(types.StringType, []attr.Value{types.StringValue("preserved")})}
	diags := SetPartialState(ctx, &state, &planned)
	require.False(t, diags.HasError(), diags)
	require.True(t, state.Raw.IsFullyKnown())
	var actual model
	require.False(t, state.Get(ctx, &actual).HasError())
	require.Equal(t, planned.ID, actual.ID)
	require.True(t, actual.Remote.IsNull())
	require.Equal(t, planned.Tags, actual.Tags)
	expected := types.ObjectValueMust(nestedTypes, map[string]attr.Value{"name": types.StringValue("configured"), "remote_status": types.StringNull()})
	require.Equal(t, types.ListValueMust(expected.Type(ctx), []attr.Value{expected}), actual.Entries)
	require.True(t, planned.Remote.IsUnknown(), "the caller's planned model must remain intact")
	original := state.Raw
	require.True(t, SetPartialState(ctx, &state, "invalid model").HasError())
	require.True(t, original.Equal(state.Raw), "encoding failure must leave the previous snapshot intact")
	require.False(t, SetPartialState(ctx, &state, &actual).HasError())
	require.True(t, original.Equal(state.Raw), "a known snapshot must be unchanged")
}
