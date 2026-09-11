package graphBetaNetworkCloudFirewallPolicyRule_test

import (
	"context"
	"testing"

	cloud "github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/resources/identity_and_access/graph_beta/network_cloud_firewall_policy_rule"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestUnitResourceNetworkCloudFirewallPolicyRule_06_SharedPortFormat(t *testing.T) {
	ctx := context.Background()
	var response resource.SchemaResponse
	cloud.NewNetworkCloudFirewallPolicyRuleResource().Schema(ctx, resource.SchemaRequest{}, &response)
	ports := response.Schema.Attributes["destinations"].(schema.SingleNestedAttribute).Attributes["ports"].(schema.SetAttribute)
	for _, tc := range []struct {
		value string
		valid bool
	}{
		{"0", true}, {"65535", true}, {"0-65535", true}, {"80-80", true}, {"443", true},
		{"65536", false}, {"65536-65536", false}, {"-1", false}, {"0443", false}, {"80-0443", false},
		{"443-80", false}, {"80-80-80", false}, {"*", false}, {"", false}, {" 443", false},
	} {
		t.Run(tc.value, func(t *testing.T) {
			value, diags := types.SetValueFrom(ctx, types.StringType, []string{tc.value})
			require.False(t, diags.HasError())
			var result validator.SetResponse
			for _, v := range ports.Validators {
				v.ValidateSet(ctx, validator.SetRequest{Path: path.Root("destinations").AtName("ports"), ConfigValue: value}, &result)
			}
			require.Equal(t, !tc.valid, result.Diagnostics.HasError(), result.Diagnostics)
		})
	}
}
