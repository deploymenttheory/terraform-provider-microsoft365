package resource

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestMutuallyExclusiveAttributes(t *testing.T) {
	t.Parallel()

	testSchema := schema.Schema{
		Attributes: map[string]schema.Attribute{
			"field1": schema.StringAttribute{
				Optional: true,
			},
			"field2": schema.SetAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}

	tests := []struct {
		name        string
		configValue map[string]tftypes.Value
		expectError bool
	}{
		{
			name: "neither field configured",
			configValue: map[string]tftypes.Value{
				"field1": tftypes.NewValue(tftypes.String, nil),
				"field2": tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil),
			},
			expectError: false,
		},
		{
			name: "only field1 configured",
			configValue: map[string]tftypes.Value{
				"field1": tftypes.NewValue(tftypes.String, "value1"),
				"field2": tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil),
			},
			expectError: false,
		},
		{
			name: "only field2 configured",
			configValue: map[string]tftypes.Value{
				"field1": tftypes.NewValue(tftypes.String, nil),
				"field2": tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{
					tftypes.NewValue(tftypes.String, "item1"),
				}),
			},
			expectError: false,
		},
		{
			name: "both fields configured",
			configValue: map[string]tftypes.Value{
				"field1": tftypes.NewValue(tftypes.String, "value1"),
				"field2": tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{
					tftypes.NewValue(tftypes.String, "item1"),
				}),
			},
			expectError: true,
		},
		{
			name: "field1 with empty string not considered configured",
			configValue: map[string]tftypes.Value{
				"field1": tftypes.NewValue(tftypes.String, ""),
				"field2": tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{
					tftypes.NewValue(tftypes.String, "item1"),
				}),
			},
			expectError: false,
		},
		{
			name: "field2 with empty set not considered configured",
			configValue: map[string]tftypes.Value{
				"field1": tftypes.NewValue(tftypes.String, "value1"),
				"field2": tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, []tftypes.Value{}),
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			validator := MutuallyExclusiveAttributes(
				[]path.Path{
					path.Root("field1"),
					path.Root("field2"),
				},
				[]string{
					"field1",
					"field2",
				},
			)

			req := resource.ValidateConfigRequest{
				Config: tfsdk.Config{
					Schema: testSchema,
					Raw:    tftypes.NewValue(testSchema.Type().TerraformType(context.Background()), tt.configValue),
				},
			}

			resp := &resource.ValidateConfigResponse{}

			validator.ValidateResource(context.Background(), req, resp)

			if tt.expectError && !resp.Diagnostics.HasError() {
				t.Fatal("expected error, got none")
			}

			if !tt.expectError && resp.Diagnostics.HasError() {
				t.Fatalf("expected no error, got: %v", resp.Diagnostics)
			}
		})
	}
}
