package graphBetaDirectoryRoleAssignment

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
	NewDirectoryRoleAssignmentResource().Schema(ctx, resource.SchemaRequest{}, resp)
	legacy := func(exp tfsdk.State) func(context.Context, tfsdk.State) bool {
		var m DirectoryRoleAssignmentResourceModel
		exp.Get(ctx, &m)
		return func() func(context.Context, tfsdk.State) bool {
			data := m
			return func(ctx context.Context, state tfsdk.State) bool {
				var actual DirectoryRoleAssignmentResourceModel
				if state.Get(ctx, &actual).HasError() {
					return false
				}
				return actual.PrincipalID.Equal(data.PrincipalID) && actual.RoleDefinitionID.Equal(data.RoleDefinitionID) && actual.DirectoryScopeID.ValueString() == "/"
			}
		}()
	}
	shared := func(exp tfsdk.State) func(context.Context, tfsdk.State) bool {
		var m DirectoryRoleAssignmentResourceModel
		exp.Get(ctx, &m)
		_ = m
		return crud.AttributesMatchPredicate(exp, readbackMatch)
	}
	allow := regexp.MustCompile(`^(plan:)?directory_scope_id=`)
	predicateequivalence.Compare(t, resp.Schema, legacy, shared, predicateequivalence.Options{Reads: 1, Allow: func(name string) bool { return allow.String() != "^$" && allow.MatchString(name) }, Base: func(a map[string]tftypes.Value) { a["directory_scope_id"] = tftypes.NewValue(tftypes.String, "/") }, Known: []string{"directory_scope_id"}})
}
