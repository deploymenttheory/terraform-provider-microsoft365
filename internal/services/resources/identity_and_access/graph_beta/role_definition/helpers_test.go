package graphBetaIdentityAndAccessRoleDefinition

import (
	"context"
	"testing"

	"github.com/deploymenttheory/terraform-provider-microsoft365/internal/services/common/crud"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func schemaResponse(ctx context.Context, r *RoleDefinitionResource) schema.Schema {
	var resp resource.SchemaResponse
	r.Schema(ctx, resource.SchemaRequest{}, &resp)
	return resp.Schema
}

func timeoutsNull() timeouts.Value {
	return timeouts.Value{Object: types.ObjectNull(map[string]attr.Type{"create": types.StringType, "read": types.StringType, "update": types.StringType, "delete": types.StringType})}
}

func TestMapOptionalString(t *testing.T) {
	empty := ""
	value := "description"

	tests := []struct {
		name     string
		input    *string
		wantNull bool
		want     string
	}{
		{name: "nil returns null", input: nil, wantNull: true},
		{name: "empty string returns null", input: &empty, wantNull: true},
		{name: "value is kept", input: &value, want: "description"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := mapOptionalString(tt.input)
			if got.IsNull() != tt.wantNull {
				t.Fatalf("IsNull() = %v, want %v", got.IsNull(), tt.wantNull)
			}
			if !tt.wantNull && got.ValueString() != tt.want {
				t.Errorf("ValueString() = %q, want %q", got.ValueString(), tt.want)
			}
		})
	}
}

func TestRoleDefinitionReadbackIgnoresPermissionOrderAndCase(t *testing.T) {
	ctx := context.Background()
	schemaResp := schemaResponse(ctx, &RoleDefinitionResource{})
	stateFor := func(model RoleDefinitionResourceModel) tfsdk.State {
		state := tfsdk.State{Schema: schemaResp}
		if diags := state.Set(ctx, &model); diags.HasError() {
			t.Fatalf("set state: %v", diags)
		}
		return state
	}

	expected := *testModel()
	expected.ResourceScopes = types.SetNull(types.StringType)
	expected.Timeouts = timeoutsNull()
	expected.RolePermissions = append(testPermissions("microsoft.directory/groups/standard/read", "microsoft.directory/applications/standard/read"), testPermissions("microsoft.directory/devices/standard/read")...)
	current := expected
	current.ID = types.StringValue("33333333-3333-3333-3333-333333333333")
	current.RolePermissions = append(testPermissions("microsoft.directory/devices/standard/read"), testPermissions("Microsoft.Directory/Applications/Standard/Read", "microsoft.directory/groups/standard/read")...)

	predicate := crud.AttributesMatchPredicate(stateFor(expected), readbackMatch)
	if !predicate(ctx, stateFor(current)) {
		t.Error("expected predicate to accept reordered permissions that differ only in case")
	}
	current.RolePermissions = testPermissions("microsoft.directory/devices/standard/read")
	if predicate(ctx, stateFor(current)) {
		t.Error("expected predicate to reject different permissions")
	}
}

func TestRoleDefinitionConsistencyPredicate(t *testing.T) {
	ctx := context.Background()
	r := &RoleDefinitionResource{}
	schemaResp := schemaResponse(ctx, r)

	expected := testModel()
	expected.Description = types.StringValue("new")
	expected.ResourceScopes = types.SetNull(types.StringType)
	expected.Timeouts = timeoutsNull()

	stateFor := func(model RoleDefinitionResourceModel) tfsdk.State {
		state := tfsdk.State{Schema: schemaResp}
		if diags := state.Set(ctx, &model); diags.HasError() {
			t.Fatalf("set state: %v", diags)
		}
		return state
	}

	current := *expected
	current.ID = types.StringValue("33333333-3333-3333-3333-333333333333")
	current.ResourceScopes = types.SetNull(types.StringType)
	current.IsBuiltIn = types.BoolValue(false)
	current.IsPrivileged = types.BoolValue(false)
	current.Timeouts = timeoutsNull()

	predicate := crud.AttributesMatchPredicate(stateFor(*expected), readbackMatch)
	if !predicate(ctx, stateFor(current)) {
		t.Error("expected predicate to accept a read matching the plan")
	}

	stale := current
	stale.Description = types.StringValue("old")
	if predicate(ctx, stateFor(stale)) {
		t.Error("expected predicate to reject a stale description")
	}

	staleActions := current
	staleActions.RolePermissions = testPermissions("microsoft.directory/groups/standard/read")
	if predicate(ctx, stateFor(staleActions)) {
		t.Error("expected predicate to reject stale role permissions")
	}

	missingID := current
	missingID.ID = types.StringNull()
	if predicate(ctx, stateFor(missingID)) {
		t.Error("expected predicate to reject a read without an ID")
	}

	missingEnabled := current
	missingEnabled.IsEnabled = types.BoolNull()
	expected.IsEnabled = types.BoolValue(false)
	if crud.AttributesMatchPredicate(stateFor(*expected), readbackMatch)(ctx, stateFor(missingEnabled)) {
		t.Error("expected predicate to reject a read without is_enabled")
	}
}
