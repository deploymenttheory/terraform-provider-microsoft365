package graphBetaDeviceConfigurationTemplatesJson

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	graphmodels "github.com/microsoftgraph/msgraph-beta-sdk-go/models"
	"github.com/stretchr/testify/require"
)

func TestUnitAssignmentFilterImport(t *testing.T) {
	ctx := context.Background()
	var response resource.SchemaResponse
	NewDeviceConfigurationTemplatesJsonResource().Schema(ctx, resource.SchemaRequest{}, &response)
	assignments := response.Schema.Attributes["assignments"].(schema.SetNestedAttribute)
	filter := assignments.NestedObject.Attributes["filter_id"].(schema.StringAttribute)
	require.True(t, filter.Optional)
	// Omitted filter IDs must remain null in both the plan and refreshed state.
	if filter.Computed || filter.Default != nil {
		t.Error("filter_id must be optional without a computed value or sentinel default")
	}

	targets := map[string]func() graphmodels.DeviceAndAppManagementAssignmentTargetable{
		"all devices": func() graphmodels.DeviceAndAppManagementAssignmentTargetable {
			return graphmodels.NewAllDevicesAssignmentTarget()
		},
		"all licensed users": func() graphmodels.DeviceAndAppManagementAssignmentTargetable {
			return graphmodels.NewAllLicensedUsersAssignmentTarget()
		},
		"group": func() graphmodels.DeviceAndAppManagementAssignmentTargetable {
			target := graphmodels.NewGroupAssignmentTarget()
			target.SetGroupId(new("22222222-2222-4222-8222-222222222222"))
			return target
		},
		"exclusion group": func() graphmodels.DeviceAndAppManagementAssignmentTargetable {
			target := graphmodels.NewExclusionGroupAssignmentTarget()
			target.SetGroupId(new("22222222-2222-4222-8222-222222222222"))
			return target
		},
	}
	for name, newTarget := range targets {
		t.Run(name, func(t *testing.T) {
			for _, tc := range []struct {
				name       string
				id         *string
				filterType graphmodels.DeviceAndAppManagementAssignmentFilterType
				want       types.String
			}{
				{"omitted", nil, graphmodels.NONE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE, types.StringNull()},
				{"empty", new(""), graphmodels.NONE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE, types.StringNull()},
				{"zero GUID", new("00000000-0000-0000-0000-000000000000"), graphmodels.NONE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE, types.StringNull()},
				{"include", new("11111111-1111-4111-8111-111111111111"), graphmodels.INCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE, types.StringValue("11111111-1111-4111-8111-111111111111")},
				{"exclude", new("11111111-1111-4111-8111-111111111111"), graphmodels.EXCLUDE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE, types.StringValue("11111111-1111-4111-8111-111111111111")},
			} {
				t.Run(tc.name, func(t *testing.T) {
					if name == "exclusion group" && tc.filterType != graphmodels.NONE_DEVICEANDAPPMANAGEMENTASSIGNMENTFILTERTYPE {
						t.Skip("Graph does not support filters on exclusion groups")
					}
					target := newTarget()
					target.SetDeviceAndAppManagementAssignmentFilterId(tc.id)
					target.SetDeviceAndAppManagementAssignmentFilterType(&tc.filterType)
					assignment := graphmodels.NewDeviceConfigurationAssignment()
					assignment.SetId(new("33333333-3333-4333-8333-333333333333"))
					assignment.SetTarget(target)
					var data DeviceConfigurationTemplatesJsonResourceModel
					diags := MapAssignmentsToTerraform(ctx, &data, []graphmodels.DeviceConfigurationAssignmentable{assignment})
					require.False(t, diags.HasError(), diags.Errors())
					require.Len(t, data.Assignments.Elements(), 1)
					object := data.Assignments.Elements()[0].(types.Object)
					got := object.Attributes()["filter_id"]
					require.True(t, got.Equal(tc.want), "filter_id = %s, want %s", got, tc.want)
					require.True(t, object.Attributes()["filter_type"].Equal(types.StringValue(tc.filterType.String())))

					// Generated import configuration uses these mapped attributes.
					raw, err := object.ToTerraformValue(ctx)
					require.NoError(t, err)
					config := tfsdk.Config{
						Schema: schema.Schema{Attributes: assignments.NestedObject.Attributes},
						Raw:    raw,
					}
					for _, rule := range filter.Validators {
						var result validator.StringResponse
						rule.ValidateString(ctx, validator.StringRequest{
							Path: path.Root("filter_id"), PathExpression: path.MatchRoot("filter_id"),
							Config: config, ConfigValue: got.(types.String),
						}, &result)
						require.False(t, result.Diagnostics.HasError(), result.Diagnostics.Errors())
					}
				})
			}
		})
	}
}
