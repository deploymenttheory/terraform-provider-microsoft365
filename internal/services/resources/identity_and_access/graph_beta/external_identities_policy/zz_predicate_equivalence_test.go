package graphBetaExternalIdentitiesPolicy

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
	NewExternalIdentitiesPolicyResource().Schema(ctx, resource.SchemaRequest{}, resp)
	legacy := func(exp tfsdk.State) func(context.Context, tfsdk.State) bool {
		var m ExternalIdentitiesPolicyResourceModel
		exp.Get(ctx, &m)
		return stableConsistencyPredicate(&m)
	}
	shared := func(exp tfsdk.State) func(context.Context, tfsdk.State) bool {
		var m ExternalIdentitiesPolicyResourceModel
		exp.Get(ctx, &m)
		_ = m
		return crud.RequireStableReads(crud.AttributesMatchPredicate(exp, readbackMatch))
	}
	allow := regexp.MustCompile(`^(plan:)?id=`)
	predicateequivalence.Compare(t, resp.Schema, legacy, shared, predicateequivalence.Options{Reads: 5, Allow: func(name string) bool { return allow.String() != "^$" && allow.MatchString(name) }, Base: func(a map[string]tftypes.Value) { a["id"] = tftypes.NewValue(tftypes.String, "externalIdentityPolicy") }, Known: []string{"id"}})
}
