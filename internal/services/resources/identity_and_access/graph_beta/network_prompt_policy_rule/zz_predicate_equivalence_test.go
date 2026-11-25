package graphBetaNetworkPromptPolicyRule

import (
	"context"
	"regexp"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud/predicateequivalence"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var (
	_ = path.Root
	_ = types.StringNull
	_ = tftypes.String
	_ = crud.RootPaths
)

func TestPredicateEquivalence(t *testing.T) {
	ctx := context.Background()
	resp := &resource.SchemaResponse{}
	NewNetworkPromptPolicyRuleResource().Schema(ctx, resource.SchemaRequest{}, resp)
	legacy := func(exp tfsdk.State) func(context.Context, tfsdk.State) bool {
		var m NetworkPromptPolicyRuleResourceModel
		exp.Get(ctx, &m)
		return promptPolicyRuleConsistencyPredicate(&m)
	}
	shared := func(exp tfsdk.State) func(context.Context, tfsdk.State) bool {
		var m NetworkPromptPolicyRuleResourceModel
		exp.Get(ctx, &m)
		_ = m
		return crud.AttributesMatchPredicate(exp, readbackMatch)
	}
	allow := regexp.MustCompile(`^$`)
	predicateequivalence.Compare(t, resp.Schema, legacy, shared, predicateequivalence.Options{Reads: 1, Allow: func(name string) bool { return allow.String() != "^$" && allow.MatchString(name) }, Base: func(a map[string]tftypes.Value) {
		l := a["conversation_schemes"]
		var elems []tftypes.Value
		_ = l.As(&elems)
		var obj map[string]tftypes.Value
		_ = elems[0].As(&obj)
		cp := map[string]tftypes.Value{}
		for k, v := range obj {
			cp[k] = v
		}
		cp["type"] = tftypes.NewValue(tftypes.String, schemeTypeCustom)
		cp["scheme_name"] = tftypes.NewValue(tftypes.String, nil)
		a["conversation_schemes"] = tftypes.NewValue(l.Type(), []tftypes.Value{tftypes.NewValue(elems[0].Type(), cp)})
	}})
}
