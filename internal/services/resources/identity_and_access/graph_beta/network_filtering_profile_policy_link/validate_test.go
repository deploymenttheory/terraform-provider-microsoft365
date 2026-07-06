package graphBetaNetworkFilteringProfilePolicyLink

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestValidateFilteringPolicyOnlyFields(t *testing.T) {
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
			diags := validateFilteringPolicyOnlyFields(tt.policyType, tt.priority, tt.loggingState)
			if len(diags) != tt.wantErrors {
				t.Fatalf("diagnostics count = %d, want %d: %#v", len(diags), tt.wantErrors, diags)
			}
		})
	}
}
