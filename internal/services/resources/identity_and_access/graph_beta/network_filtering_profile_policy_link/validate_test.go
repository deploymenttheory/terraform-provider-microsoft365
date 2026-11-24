package graphBetaNetworkFilteringProfilePolicyLink

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestFilteringPolicyOnlyFieldsValidator(t *testing.T) {
	tests := []struct {
		name         string
		policyType   types.String
		priority     types.Int64
		loggingState types.String
		wantErrors   int
	}{
		{
			name:         "filtering policy allows both fields",
			policyType:   types.StringValue(policyTypeFiltering),
			priority:     types.Int64Value(100),
			loggingState: types.StringValue("enabled"),
			wantErrors:   0,
		},
		{
			name:         "web filtering policy with omitted fields",
			policyType:   types.StringValue(policyTypeWebFiltering),
			priority:     types.Int64Null(),
			loggingState: types.StringNull(),
			wantErrors:   0,
		},
		{
			name:         "prompt policy with omitted fields",
			policyType:   types.StringValue(policyTypePrompt),
			priority:     types.Int64Null(),
			loggingState: types.StringNull(),
			wantErrors:   0,
		},
		{
			name:         "content policy with omitted fields",
			policyType:   types.StringValue(policyTypeContent),
			priority:     types.Int64Null(),
			loggingState: types.StringNull(),
			wantErrors:   0,
		},
		{
			name:         "netskope dlp policy with omitted fields",
			policyType:   types.StringValue(policyTypeNetskopeDlp),
			priority:     types.Int64Null(),
			loggingState: types.StringNull(),
			wantErrors:   0,
		},
		{
			name:         "web filtering policy rejects priority",
			policyType:   types.StringValue(policyTypeWebFiltering),
			priority:     types.Int64Value(100),
			loggingState: types.StringNull(),
			wantErrors:   1,
		},
		{
			name:         "web filtering policy rejects logging state",
			policyType:   types.StringValue(policyTypeWebFiltering),
			priority:     types.Int64Null(),
			loggingState: types.StringValue("enabled"),
			wantErrors:   1,
		},
		{
			name:         "web filtering policy rejects both fields",
			policyType:   types.StringValue(policyTypeWebFiltering),
			priority:     types.Int64Value(100),
			loggingState: types.StringValue("enabled"),
			wantErrors:   2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			schemaResp := &resource.SchemaResponse{}
			(&NetworkFilteringProfilePolicyLinkResource{}).Schema(ctx, resource.SchemaRequest{}, schemaResp)
			objectType := schemaResp.Schema.Type().TerraformType(ctx).(tftypes.Object)
			attrs := map[string]tftypes.Value{}
			for name, typ := range objectType.AttributeTypes {
				attrs[name] = tftypes.NewValue(typ, nil)
			}
			for name, value := range map[string]interface {
				ToTerraformValue(context.Context) (tftypes.Value, error)
			}{"policy_type": tt.policyType, "priority": tt.priority, "logging_state": tt.loggingState} {
				tfValue, err := value.ToTerraformValue(ctx)
				if err != nil {
					t.Fatalf("convert %s: %v", name, err)
				}
				attrs[name] = tfValue
			}

			resp := &resource.ValidateConfigResponse{}
			filteringPolicyOnlyFieldsValidator{}.ValidateResource(ctx, resource.ValidateConfigRequest{Config: tfsdk.Config{Schema: schemaResp.Schema, Raw: tftypes.NewValue(objectType, attrs)}}, resp)
			if resp.Diagnostics.ErrorsCount() != tt.wantErrors {
				t.Fatalf("diagnostics count = %d, want %d: %#v", resp.Diagnostics.ErrorsCount(), tt.wantErrors, resp.Diagnostics)
			}
		})
	}
}
